package cluster

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	apisparkv1beta2 "github.com/GoogleCloudPlatform/spark-on-k8s-operator/pkg/apis/sparkoperator.k8s.io/v1beta2"
	sparkclient "github.com/GoogleCloudPlatform/spark-on-k8s-operator/pkg/client/clientset/versioned"
	sparkoperatorv1beta2 "github.com/GoogleCloudPlatform/spark-on-k8s-operator/pkg/client/clientset/versioned/typed/sparkoperator.k8s.io/v1beta2" //nolint
	mlpcluster "github.com/caraml-dev/mlp/api/pkg/cluster"
	"github.com/caraml-dev/turing/api/turing/cluster/labeller"
	"github.com/pkg/errors"
	networkingv1beta1 "istio.io/client-go/pkg/clientset/versioned/typed/networking/v1beta1"
	apiappsv1 "k8s.io/api/apps/v1"
	apibatchv1 "k8s.io/api/batch/v1"
	apicorev1 "k8s.io/api/core/v1"
	apipolicyv1 "k8s.io/api/policy/v1"
	apirbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	policyv1cfg "k8s.io/client-go/applyconfigurations/policy/v1"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	batchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	policyv1 "k8s.io/client-go/kubernetes/typed/policy/v1"
	rbacv1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	rest "k8s.io/client-go/rest"
	"knative.dev/pkg/kmp"
	knservingv1 "knative.dev/serving/pkg/apis/serving/v1"
	knservingclientset "knative.dev/serving/pkg/client/clientset/versioned"
	knservingclient "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"

	// Load required auth plugin
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/caraml-dev/turing/api/turing/config"
)

var ErrNamespaceAlreadyExists = errors.New("namespace already exists")

// clusterConfig Model cluster authentication settings
type clusterConfig struct {
	// Use Kubernetes service account in cluster config
	InClusterConfig bool

	// Credentials Interface to access cluster
	Credentials mlpcluster.Credentials
}

// Controller defines the operations supported by the cluster controller
type Controller interface {
	// Namespace
	CreateNamespace(ctx context.Context, name string) error

	// Knative Service
	GetKnativeServiceURL(ctx context.Context, svcName string, namespace string) string
	DeployKnativeService(ctx context.Context, svc *KnativeService) error
	DeleteKnativeService(ctx context.Context, svcName string, namespace string, ignoreNotFound bool) error
	GetKnativeServiceDesiredReplicas(ctx context.Context, svcName string, namespace string) (int, error)

	// Istio VirtualService
	ApplyIstioVirtualService(ctx context.Context, routerEndpoint *VirtualService) error
	DeleteIstioVirtualService(ctx context.Context, svcName string, namespace string) error

	// StatefulSet
	DeleteKubernetesStatefulSet(ctx context.Context, name string, namespace string, ignoreNotFound bool) error

	// Service
	DeployKubernetesService(ctx context.Context, svc *KubernetesService) error
	DeleteKubernetesService(ctx context.Context, svcName string, namespace string, ignoreNotFound bool) error

	// ConfigMap
	ApplyConfigMap(ctx context.Context, namespace string, configMap *ConfigMap) error
	DeleteConfigMap(ctx context.Context, name string, namespace string, ignoreNotFound bool) error

	// Secret
	CreateSecret(ctx context.Context, secret *Secret) error
	DeleteSecret(ctx context.Context, secretName string, namespace string, ignoreNotFound bool) error

	// PVC
	DeletePVCs(ctx context.Context, listOptions metav1.ListOptions, namespace string, ignoreNotFound bool) error

	// Pod
	ListPods(ctx context.Context, namespace string, labelSelector string) (*apicorev1.PodList, error)
	ListPodLogs(ctx context.Context, namespace string,
		podName string, opts *apicorev1.PodLogOptions) (io.ReadCloser, error)

	// Job
	CreateJob(ctx context.Context, namespace string, job Job) (*apibatchv1.Job, error)
	GetJob(ctx context.Context, namespace string, jobName string) (*apibatchv1.Job, error)
	DeleteJob(ctx context.Context, namespace string, jobName string) error

	// Service Account
	CreateServiceAccount(ctx context.Context, namespace string,
		serviceAccount *ServiceAccount) (*apicorev1.ServiceAccount, error)
	CreateRole(ctx context.Context, namespace string, role *Role) (*apirbacv1.Role, error)
	CreateRoleBinding(ctx context.Context, namespace string, roleBinding *RoleBinding) (*apirbacv1.RoleBinding, error)

	// Spark Application
	CreateSparkApplication(ctx context.Context, namespace string,
		request *CreateSparkRequest) (*apisparkv1beta2.SparkApplication, error)
	GetSparkApplication(ctx context.Context, namespace, appName string) (*apisparkv1beta2.SparkApplication, error)
	DeleteSparkApplication(ctx context.Context, namespace, appName string) error

	// PodDisruptionBudget
	ApplyPodDisruptionBudget(ctx context.Context, pdb PodDisruptionBudget) (*apipolicyv1.PodDisruptionBudget, error)
	DeletePodDisruptionBudget(ctx context.Context, namespace, pdbName string) error
}

// controller implements the Controller interface
type controller struct {
	knServingClient  knservingclient.ServingV1Interface
	k8sCoreClient    corev1.CoreV1Interface
	k8sAppsClient    appsv1.AppsV1Interface
	k8sBatchClient   batchv1.BatchV1Interface
	k8sRBACClient    rbacv1.RbacV1Interface
	k8sPolicyClient  policyv1.PolicyV1Interface
	k8sSparkOperator sparkoperatorv1beta2.SparkoperatorV1beta2Interface
	istioClient      networkingv1beta1.NetworkingV1beta1Interface
}

// newController initializes a new cluster controller with the given cluster config
func newController(clusterCfg clusterConfig) (Controller, error) {
	var cfg *rest.Config
	var err error
	if clusterCfg.InClusterConfig {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		cfg, err = clusterCfg.Credentials.ToRestConfig()
		if err != nil {
			return nil, err
		}
	}

	// Init the knative serving client
	knsClientSet, err := knservingclientset.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	// Init the k8s clientset
	k8sClientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	istioClientSet, err := networkingv1beta1.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	sparkClient, err := sparkclient.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &controller{
		knServingClient:  knsClientSet.ServingV1(),
		k8sCoreClient:    k8sClientset.CoreV1(),
		k8sAppsClient:    k8sClientset.AppsV1(),
		k8sBatchClient:   k8sClientset.BatchV1(),
		k8sRBACClient:    k8sClientset.RbacV1(),
		k8sPolicyClient:  k8sClientset.PolicyV1(),
		k8sSparkOperator: sparkClient.SparkoperatorV1beta2(),
		istioClient:      istioClientSet,
	}, nil
}

// InitClusterControllers takes in the app config and initializes one cluster controller
// per environment and returns a map where the
// key is the env name and the value is the corresponding controller.
func InitClusterControllers(
	cfg *config.Config,
	environmentClusterMap map[string]*mlpcluster.K8sConfig,
) (map[string]Controller, error) {
	// For each supported environment, init controller
	controllers := make(map[string]Controller)
	for envName, k := range environmentClusterMap {
		clusterCfg := clusterConfig{}

		if cfg.ClusterConfig.InClusterConfig {
			clusterCfg.InClusterConfig = true
		} else {
			creds := mlpcluster.NewK8sClusterCreds(k)
			clusterCfg.Credentials = creds
		}

		ctl, err := newController(clusterCfg)
		if err != nil {
			return nil, errors.Wrap(err, "unable to initialize cluster controller")
		}

		controllers[envName] = ctl
	}

	return controllers, nil
}

// CreateNamespace creates a namespace. If the namespace already exists, will throw an error.
func (c *controller) CreateNamespace(ctx context.Context, name string) error {
	_, err := c.k8sCoreClient.Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		return ErrNamespaceAlreadyExists
	}
	ns := apicorev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labeller.BuildNamespaceLabels(labeller.KubernetesLabelsRequest{}),
		},
	}
	_, err = c.k8sCoreClient.Namespaces().Create(ctx, &ns, metav1.CreateOptions{})
	return err
}

// ApplyConfigMap creates a config map in the namespace given the configuration if not exists.
// If the config map already exists, ApplyConfigMap will update the configuration with the given
// data.
func (c *controller) ApplyConfigMap(ctx context.Context, namespace string, configMap *ConfigMap) error {
	cm := apicorev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMap.Name,
			Namespace: namespace,
			Labels:    configMap.Labels,
		},
		Data: map[string]string{
			configMap.FileName: configMap.Data,
		},
	}
	_, err := c.k8sCoreClient.ConfigMaps(namespace).Get(ctx, cm.Name, metav1.GetOptions{})
	if err == nil {
		// exists, we update instead
		_, err = c.k8sCoreClient.ConfigMaps(namespace).Update(ctx, &cm, metav1.UpdateOptions{})
		return err
	}
	_, err = c.k8sCoreClient.ConfigMaps(namespace).Create(ctx, &cm, metav1.CreateOptions{})
	return err
}

// DeleteConfigMap deletes a configmap if exists.
func (c *controller) DeleteConfigMap(ctx context.Context, name string, namespace string, ignoreNotFound bool) error {
	_, err := c.k8sCoreClient.ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if ignoreNotFound {
			return nil
		}
		return err
	}
	return c.k8sCoreClient.ConfigMaps(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// Deploy creates / updates a Kubernetes/Knative service with the given specs
func (c *controller) DeployKnativeService(ctx context.Context, svcConf *KnativeService) error {
	var existingSvc *knservingv1.Service
	var err error

	// Build the deployment specs
	desiredSvc, err := svcConf.BuildKnativeServiceConfig()
	if err != nil {
		return err
	}

	// Init knative ServicesGetter
	services := c.knServingClient.Services(svcConf.Namespace)

	// Check if service already exists. If exists, update it. If not, create.
	existingSvc, err = services.Get(ctx, svcConf.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Create new service
			_, err = services.Create(ctx, desiredSvc, metav1.CreateOptions{})
		} else {
			// Unexpected error, return it
			return err
		}
	} else {
		// Check for differences between current and new specs
		if !knServiceSemanticEquals(desiredSvc, existingSvc) {
			_, err = kmp.SafeDiff(
				desiredSvc.Spec.ConfigurationSpec,
				existingSvc.Spec.ConfigurationSpec,
			)
			if err != nil {
				return fmt.Errorf("Failed to diff Knative Service: %v", err)
			}
			// Update the existing service with the new config
			existingSvc.Spec.ConfigurationSpec = desiredSvc.Spec.ConfigurationSpec
			existingSvc.ObjectMeta.Labels = desiredSvc.ObjectMeta.Labels
			_, err = services.Update(ctx, existingSvc, metav1.UpdateOptions{})
		}
	}

	if err != nil {
		return err
	}

	// Wait until service ready and return any errors
	return c.waitKnativeServiceReady(ctx, svcConf.Name, svcConf.Namespace)
}

// Delete removes the Kubernetes/Knative service and all related artifacts
func (c *controller) DeleteKnativeService(
	ctx context.Context,
	svcName string,
	namespace string,
	ignoreNotFound bool,
) error {
	// Init knative ServicesGetter
	services := c.knServingClient.Services(namespace)

	// Get the service
	_, err := services.Get(ctx, svcName, metav1.GetOptions{})
	if err != nil {
		if ignoreNotFound {
			return nil
		}
		return err
	}

	// Delete the service
	return services.Delete(ctx, svcName, metav1.DeleteOptions{})
}

// GetKnativeServiceDesiredReplicas determines the current desired number of replicas for the default revision
func (c *controller) GetKnativeServiceDesiredReplicas(
	ctx context.Context,
	svcName string,
	namespace string,
) (int, error) {
	// Init knative revisions getter
	revisions := c.knServingClient.Revisions(namespace)

	rev, err := revisions.Get(ctx, getDefaultRevisionName(svcName), metav1.GetOptions{})
	if err != nil {
		return 0, err
	}
	if rev.Status.DesiredReplicas == nil {
		return 0, fmt.Errorf("Desired Replicas for %s/%s is not set", namespace, svcName)
	}
	return int(*rev.Status.DesiredReplicas), nil
}

// DeployKubernetesService deploys a kubernetes service and stateful set
func (c *controller) DeployKubernetesService(
	ctx context.Context,
	svcConf *KubernetesService,
) error {
	desiredStatefulSet, desiredSvc := svcConf.BuildKubernetesServiceConfig()

	// Deploy stateful set
	statefulSets := c.k8sAppsClient.StatefulSets(svcConf.Namespace)
	// Check if stateful set already exists. If exists, update it. If not, create.
	var existingStatefulSet *apiappsv1.StatefulSet
	var err error
	existingStatefulSet, err = statefulSets.Get(ctx, svcConf.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Create new deployment
			_, err = statefulSets.Create(ctx, desiredStatefulSet, metav1.CreateOptions{})
		} else {
			// Unexpected error, return it
			return err
		}
	} else {
		// Check for differences between current and new specs
		if !k8sStatefulSetSemanticEquals(existingStatefulSet, existingStatefulSet) {
			// Update the existing service with the new config
			existingStatefulSet.Spec.Template = desiredStatefulSet.Spec.Template
			existingStatefulSet.ObjectMeta.Labels = desiredStatefulSet.ObjectMeta.Labels
			_, err = statefulSets.Update(ctx, existingStatefulSet, metav1.UpdateOptions{})
		}
	}
	if err != nil {
		return err
	}

	// Deploy Service
	services := c.k8sCoreClient.Services(svcConf.Namespace)
	// Check if service already exists. If exists, update it. If not, create.
	var existingSvc *apicorev1.Service
	existingSvc, err = services.Get(ctx, svcConf.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Create new service
			_, err = services.Create(ctx, desiredSvc, metav1.CreateOptions{})
		} else {
			// Unexpected error, return it
			return err
		}
	} else {
		// Check for differences between current and new specs
		if !k8sServiceSemanticEquals(desiredSvc, existingSvc) {
			_, err = services.Update(ctx, desiredSvc, metav1.UpdateOptions{})
		}
	}
	if err != nil {
		return err
	}

	// Wait until deployment ready and return any errors
	return c.waitStatefulSetReady(ctx, svcConf.Name, svcConf.Namespace)
}

// DeleteKubernetesStatefulSet deletes a stateful set
func (c *controller) DeleteKubernetesStatefulSet(
	ctx context.Context,
	name string,
	namespace string,
	ignoreNotFound bool,
) error {
	statefulSets := c.k8sAppsClient.StatefulSets(namespace)
	_, err := statefulSets.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if ignoreNotFound {
			return nil
		}
		return err
	}
	return statefulSets.Delete(ctx, name, metav1.DeleteOptions{})
}

// DeleteKubernetesService deletes a kubernetes service
func (c *controller) DeleteKubernetesService(
	ctx context.Context,
	svcName string,
	namespace string,
	ignoreNotFound bool,
) error {
	services := c.k8sCoreClient.Services(namespace)
	_, err := services.Get(ctx, svcName, metav1.GetOptions{})
	if err != nil {
		if ignoreNotFound {
			return nil
		}
		return err
	}
	return services.Delete(ctx, svcName, metav1.DeleteOptions{})
}

// ApplyIstioVirtualService creates a virtual service if not exists, if exists, updates the
// existing service with the new configuration
func (c *controller) ApplyIstioVirtualService(ctx context.Context, routerEndpoint *VirtualService) error {
	vservices := c.istioClient.VirtualServices(routerEndpoint.Namespace)
	existingVsvc, err := vservices.Get(ctx, routerEndpoint.Name, metav1.GetOptions{})
	if err == nil {
		// patch
		existingVsvc.Spec.Http = routerEndpoint.BuildVirtualService().Spec.Http
		_, err := vservices.Update(ctx, existingVsvc, metav1.UpdateOptions{})
		return err
	}
	_, err = vservices.Create(ctx, routerEndpoint.BuildVirtualService(), metav1.CreateOptions{})
	return err
}

// DeleteIstioVirtualService deletes an istio virtual service.
func (c *controller) DeleteIstioVirtualService(
	ctx context.Context,
	svcName string,
	namespace string,
) error {
	vservices := c.istioClient.VirtualServices(namespace)
	_, err := vservices.Get(ctx, svcName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("unable to retrieve virtual service %s: %s", svcName, err.Error())
	}
	return vservices.Delete(ctx, svcName, metav1.DeleteOptions{})
}

// CreateSecret creates a secret. If the secret already exists, the existing secret will be updated.
func (c *controller) CreateSecret(ctx context.Context, secret *Secret) error {
	secrets := c.k8sCoreClient.Secrets(secret.Namespace)
	_, err := secrets.Get(ctx, secret.Name, metav1.GetOptions{})
	if err == nil {
		_, err = secrets.Update(ctx, secret.BuildSecret(), metav1.UpdateOptions{})
		return err
	}
	_, err = secrets.Create(ctx, secret.BuildSecret(), metav1.CreateOptions{})
	return err
}

// DeleteSecret deletes a secret
func (c *controller) DeleteSecret(ctx context.Context, secretName string, namespace string, ignoreNotFound bool) error {
	secrets := c.k8sCoreClient.Secrets(namespace)
	_, err := secrets.Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		if ignoreNotFound {
			return nil
		}
		return fmt.Errorf("unable to get secret with name %s: %s", secretName, err.Error())
	}
	return secrets.Delete(ctx, secretName, metav1.DeleteOptions{})
}

// GetKnativeServiceURL returns the URL at which the given service is accessible, if found.
// Else, an empty string is returned.
func (c *controller) GetKnativeServiceURL(ctx context.Context, svcName string, namespace string) string {
	var url string

	// Init knative ServicesGetter
	services := c.knServingClient.Services(namespace)

	// Get the service
	svc, err := services.Get(ctx, svcName, metav1.GetOptions{})

	// If the service is ready, get its url
	if err == nil && svc.IsReady() {
		url = svc.Status.URL.String()
	}

	return url
}

// DeletePVCs deletes all PVCs specified by the given list options in the given namespace.
func (c *controller) DeletePVCs(
	ctx context.Context,
	listOptions metav1.ListOptions,
	namespace string,
	ignoreNotFound bool,
) error {
	_, err := c.k8sCoreClient.PersistentVolumeClaims(namespace).List(ctx, listOptions)
	if err != nil {
		if ignoreNotFound {
			return nil
		}
		return err
	}
	return c.k8sCoreClient.PersistentVolumeClaims(namespace).DeleteCollection(ctx, metav1.DeleteOptions{}, listOptions)
}

func (c *controller) ListPods(ctx context.Context, namespace string, labelSelector string) (*apicorev1.PodList, error) {
	return c.k8sCoreClient.Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
}

func (c *controller) ListPodLogs(
	ctx context.Context,
	namespace string,
	podName string,
	opts *apicorev1.PodLogOptions,
) (io.ReadCloser, error) {
	return c.k8sCoreClient.Pods(namespace).GetLogs(podName, opts).Stream(ctx)
}

// CreateJob creates a Kubernetes job
func (c *controller) CreateJob(ctx context.Context, namespace string, job Job) (*apibatchv1.Job, error) {
	return c.k8sBatchClient.Jobs(namespace).Create(ctx, job.Build(), metav1.CreateOptions{})
}

// GetJob gets the Kubernetes job
func (c *controller) GetJob(ctx context.Context, namespace, jobName string) (*apibatchv1.Job, error) {
	return c.k8sBatchClient.Jobs(namespace).Get(ctx, jobName, metav1.GetOptions{})
}

// DeleteJob deletes the Kubernetes job
func (c *controller) DeleteJob(ctx context.Context, namespace, jobName string) error {
	return c.k8sBatchClient.Jobs(namespace).Delete(ctx, jobName, metav1.DeleteOptions{})
}

func (c *controller) CreateServiceAccount(
	ctx context.Context,
	namespace string,
	serviceAccount *ServiceAccount,
) (*apicorev1.ServiceAccount, error) {
	sa, err := c.k8sCoreClient.ServiceAccounts(namespace).Get(ctx, serviceAccount.Name, metav1.GetOptions{})
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return nil, errors.Errorf(
				"failed getting status of driver service account %s in namespace %s",
				serviceAccount.Name,
				namespace,
			)
		}

		saCfg := serviceAccount.BuildServiceAccount()
		sa, err = c.k8sCoreClient.ServiceAccounts(namespace).Create(ctx, saCfg, metav1.CreateOptions{})

		if err != nil {
			return nil, errors.Errorf(
				"failed creating driver service account %s in namespace %s", serviceAccount.Name, namespace,
			)
		}
	}

	return sa, nil
}

func (c *controller) CreateRole(
	ctx context.Context,
	namespace string,
	r *Role,
) (*apirbacv1.Role, error) {
	role, err := c.k8sRBACClient.Roles(namespace).Get(ctx, r.Name, metav1.GetOptions{})
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return nil, errors.Errorf(
				"failed getting status of driver role %s in namespace %s: %s",
				r.Name,
				namespace,
				err.Error(),
			)
		}

		roleCfg := r.BuildRole()
		role, err = c.k8sRBACClient.Roles(namespace).Create(ctx, roleCfg, metav1.CreateOptions{})

		if err != nil {
			return nil, errors.Errorf(
				"failed creating driver roles %s in namespace %s: %s",
				r.Name,
				namespace,
				err.Error(),
			)
		}
	}

	return role, nil
}

func (c *controller) CreateRoleBinding(
	ctx context.Context,
	namespace string,
	roleBinding *RoleBinding,
) (*apirbacv1.RoleBinding, error) {
	rb, err := c.k8sRBACClient.RoleBindings(namespace).Get(ctx, roleBinding.Name, metav1.GetOptions{})
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return nil, errors.Errorf(
				"failed getting status of driver rolebinding %s in namespace %s: %s",
				roleBinding.Name,
				namespace,
				err.Error(),
			)
		}

		rbConfig := roleBinding.BuildRoleBinding()
		rb, err = c.k8sRBACClient.RoleBindings(namespace).Create(ctx, rbConfig, metav1.CreateOptions{})

		if err != nil {
			return nil, errors.Errorf(
				"failed creating driver roles binding %s in namespace %s: %s",
				roleBinding.Name,
				namespace,
				err.Error(),
			)
		}
	}

	return rb, nil
}

func (c *controller) CreateSparkApplication(
	ctx context.Context,
	namespace string,
	request *CreateSparkRequest,
) (*apisparkv1beta2.SparkApplication, error) {
	s, err := createSparkRequest(request)
	if err != nil {
		return nil, err
	}

	return c.k8sSparkOperator.SparkApplications(namespace).Create(ctx, s, metav1.CreateOptions{})
}

func (c *controller) GetSparkApplication(
	ctx context.Context,
	namespace string,
	appName string,
) (*apisparkv1beta2.SparkApplication, error) {
	return c.k8sSparkOperator.SparkApplications(namespace).Get(ctx, appName, metav1.GetOptions{})
}

func (c *controller) DeleteSparkApplication(ctx context.Context, namespace, appName string) error {
	return c.k8sSparkOperator.SparkApplications(namespace).Delete(ctx, appName, metav1.DeleteOptions{})
}

func (c *controller) ApplyPodDisruptionBudget(
	ctx context.Context,
	pdb PodDisruptionBudget,
) (*apipolicyv1.PodDisruptionBudget, error) {
	pdbSpec, err := pdb.BuildPDBSpec()
	if err != nil {
		return nil, err
	}

	pdbCfg := policyv1cfg.PodDisruptionBudget(pdb.Name, pdb.Namespace)
	pdbCfg.WithLabels(pdb.Labels)
	pdbCfg.WithSpec(pdbSpec)

	pdbObj, err := c.k8sPolicyClient.PodDisruptionBudgets(pdb.Namespace).Apply(
		ctx,
		pdbCfg,
		metav1.ApplyOptions{
			FieldManager: "application/apply-patch",
		},
	)
	if err != nil {
		return nil, err
	}

	return pdbObj, nil
}

func (c *controller) DeletePodDisruptionBudget(ctx context.Context, namespace, pdbName string) error {
	_, err := c.k8sPolicyClient.PodDisruptionBudgets(namespace).Get(ctx, pdbName, metav1.GetOptions{})
	if err != nil {
		// If pdb not found, it might be already deleted or not exist at all. So we just return nil here.
		if k8serrors.IsNotFound(err) {
			return nil
		}

		return fmt.Errorf(
			"failed getting status of pod disruption budget %s in namespace %s: %w",
			pdbName,
			namespace,
			err,
		)
	}
	return c.k8sPolicyClient.PodDisruptionBudgets(namespace).Delete(ctx, pdbName, metav1.DeleteOptions{})
}

// waitKnativeServiceReady waits for the given knative service to become ready, until the
// default timeout
func (c *controller) waitKnativeServiceReady(
	ctx context.Context,
	svcName string,
	namespace string,
) error {
	// Init ticker to check status every second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Init knative ServicesGetter
	services := c.knServingClient.Services(namespace)

	for {
		select {
		case <-ctx.Done():
			terminationMessage := c.getKnativePodTerminationMessage(context.Background(), svcName, namespace)
			if terminationMessage == "" {
				// Pod was not created (as with invalid image names), get status messages from the knative service.
				svc, err := services.Get(ctx, svcName, metav1.GetOptions{})
				if err != nil {
					terminationMessage = err.Error()
				} else {
					terminationMessage = getKnServiceStatusMessages(svc)
				}
			}
			return fmt.Errorf("timeout waiting for service %s to be ready:\n%s", svcName, terminationMessage)
		case <-ticker.C:
			svc, err := services.Get(ctx, svcName, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("unable to get service status for %s: %v", svcName, err)
			}

			if svc.IsReady() {
				// Service is completely ready
				return nil
			}
		}
	}
}

// getKnativePodTerminationMessage retrieves the termination message of the user container
// in the pod, which will be returned for logging as a part of the deployment failure error.
func (c *controller) getKnativePodTerminationMessage(ctx context.Context, svcName string, namespace string) string {
	labelSelector := KnativeServiceLabelKey + "=" + svcName
	podList, err := c.ListPods(ctx, namespace, labelSelector)
	if err != nil {
		return err.Error()
	}

	var terminationMessage string
	for _, pod := range podList.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Name == KnativeUserContainerName {
				if containerStatus.LastTerminationState.Terminated != nil {
					terminationMessage = containerStatus.LastTerminationState.Terminated.Message
					break
				}
			}
		}
	}
	return terminationMessage
}

// waitStatefulSetReady waits for the given k8s stateful set to become ready, until the
// default timeout
func (c *controller) waitStatefulSetReady(
	ctx context.Context,
	statefulSetName string,
	namespace string,
) error {
	// Init ticker to check status every second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Init stateful set getter
	statefulSets := c.k8sAppsClient.StatefulSets(namespace)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for stateful set %s to be ready", statefulSetName)
		case <-ticker.C:
			statefulSet, err := statefulSets.Get(ctx, statefulSetName, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("unable to get stateful set status for %s: %v", statefulSetName, err)
			}

			if statefulSetReady(statefulSet) {
				// Service is completely ready
				return nil
			}
		}
	}
}

func statefulSetReady(statefulSet *apiappsv1.StatefulSet) bool {
	if statefulSet.Generation <= statefulSet.Status.ObservedGeneration {
		if statefulSet.Spec.Replicas != nil {
			// Account for replica surge during updates
			return statefulSet.Status.ReadyReplicas == *statefulSet.Spec.Replicas &&
				statefulSet.Status.Replicas == *statefulSet.Spec.Replicas
		}
	}
	return false
}

func getKnServiceStatusMessages(svc *knservingv1.Service) string {
	logs := []string{}
	conditions := svc.Status.GetConditions()
	for _, cond := range conditions {
		logs = append(logs, fmt.Sprintf("Type: %s, Status: %t. %s", cond.Type, cond.IsTrue(), cond.GetMessage()))
	}
	return strings.Join(logs, "\n")
}

func knServiceSemanticEquals(desiredService, service *knservingv1.Service) bool {
	return equality.Semantic.DeepEqual(
		desiredService.Spec.ConfigurationSpec,
		service.Spec.ConfigurationSpec) &&
		equality.Semantic.DeepEqual(desiredService.ObjectMeta.Labels, service.ObjectMeta.Labels)
}

func k8sStatefulSetSemanticEquals(desiredStatefulSet, statefulSet *apiappsv1.StatefulSet) bool {
	return equality.Semantic.DeepEqual(desiredStatefulSet.Spec.Template, statefulSet.Spec.Template) &&
		equality.Semantic.DeepEqual(desiredStatefulSet.ObjectMeta.Labels, statefulSet.ObjectMeta.Labels) &&
		desiredStatefulSet.Spec.Replicas == statefulSet.Spec.Replicas
}

func k8sServiceSemanticEquals(desiredService, service *apicorev1.Service) bool {
	return equality.Semantic.DeepEqual(desiredService.Spec.Ports, service.Spec.Ports) &&
		equality.Semantic.DeepEqual(desiredService.ObjectMeta.Labels, service.ObjectMeta.Labels) &&
		equality.Semantic.DeepEqual(desiredService.Spec.Selector, service.Spec.Selector)
}
