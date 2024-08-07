/*
 * Turing Minimal Openapi Spec for SDK
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// EnsemblingResources struct for EnsemblingResources
type EnsemblingResources struct {
	DriverCpuRequest *string `json:"driver_cpu_request,omitempty"`
	DriverMemoryRequest *string `json:"driver_memory_request,omitempty"`
	ExecutorReplica *int32 `json:"executor_replica,omitempty"`
	ExecutorCpuRequest *string `json:"executor_cpu_request,omitempty"`
	ExecutorMemoryRequest *string `json:"executor_memory_request,omitempty"`
}

// NewEnsemblingResources instantiates a new EnsemblingResources object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewEnsemblingResources() *EnsemblingResources {
	this := EnsemblingResources{}
	return &this
}

// NewEnsemblingResourcesWithDefaults instantiates a new EnsemblingResources object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewEnsemblingResourcesWithDefaults() *EnsemblingResources {
	this := EnsemblingResources{}
	return &this
}

// GetDriverCpuRequest returns the DriverCpuRequest field value if set, zero value otherwise.
func (o *EnsemblingResources) GetDriverCpuRequest() string {
	if o == nil || o.DriverCpuRequest == nil {
		var ret string
		return ret
	}
	return *o.DriverCpuRequest
}

// GetDriverCpuRequestOk returns a tuple with the DriverCpuRequest field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EnsemblingResources) GetDriverCpuRequestOk() (*string, bool) {
	if o == nil || o.DriverCpuRequest == nil {
		return nil, false
	}
	return o.DriverCpuRequest, true
}

// HasDriverCpuRequest returns a boolean if a field has been set.
func (o *EnsemblingResources) HasDriverCpuRequest() bool {
	if o != nil && o.DriverCpuRequest != nil {
		return true
	}

	return false
}

// SetDriverCpuRequest gets a reference to the given string and assigns it to the DriverCpuRequest field.
func (o *EnsemblingResources) SetDriverCpuRequest(v string) {
	o.DriverCpuRequest = &v
}

// GetDriverMemoryRequest returns the DriverMemoryRequest field value if set, zero value otherwise.
func (o *EnsemblingResources) GetDriverMemoryRequest() string {
	if o == nil || o.DriverMemoryRequest == nil {
		var ret string
		return ret
	}
	return *o.DriverMemoryRequest
}

// GetDriverMemoryRequestOk returns a tuple with the DriverMemoryRequest field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EnsemblingResources) GetDriverMemoryRequestOk() (*string, bool) {
	if o == nil || o.DriverMemoryRequest == nil {
		return nil, false
	}
	return o.DriverMemoryRequest, true
}

// HasDriverMemoryRequest returns a boolean if a field has been set.
func (o *EnsemblingResources) HasDriverMemoryRequest() bool {
	if o != nil && o.DriverMemoryRequest != nil {
		return true
	}

	return false
}

// SetDriverMemoryRequest gets a reference to the given string and assigns it to the DriverMemoryRequest field.
func (o *EnsemblingResources) SetDriverMemoryRequest(v string) {
	o.DriverMemoryRequest = &v
}

// GetExecutorReplica returns the ExecutorReplica field value if set, zero value otherwise.
func (o *EnsemblingResources) GetExecutorReplica() int32 {
	if o == nil || o.ExecutorReplica == nil {
		var ret int32
		return ret
	}
	return *o.ExecutorReplica
}

// GetExecutorReplicaOk returns a tuple with the ExecutorReplica field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EnsemblingResources) GetExecutorReplicaOk() (*int32, bool) {
	if o == nil || o.ExecutorReplica == nil {
		return nil, false
	}
	return o.ExecutorReplica, true
}

// HasExecutorReplica returns a boolean if a field has been set.
func (o *EnsemblingResources) HasExecutorReplica() bool {
	if o != nil && o.ExecutorReplica != nil {
		return true
	}

	return false
}

// SetExecutorReplica gets a reference to the given int32 and assigns it to the ExecutorReplica field.
func (o *EnsemblingResources) SetExecutorReplica(v int32) {
	o.ExecutorReplica = &v
}

// GetExecutorCpuRequest returns the ExecutorCpuRequest field value if set, zero value otherwise.
func (o *EnsemblingResources) GetExecutorCpuRequest() string {
	if o == nil || o.ExecutorCpuRequest == nil {
		var ret string
		return ret
	}
	return *o.ExecutorCpuRequest
}

// GetExecutorCpuRequestOk returns a tuple with the ExecutorCpuRequest field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EnsemblingResources) GetExecutorCpuRequestOk() (*string, bool) {
	if o == nil || o.ExecutorCpuRequest == nil {
		return nil, false
	}
	return o.ExecutorCpuRequest, true
}

// HasExecutorCpuRequest returns a boolean if a field has been set.
func (o *EnsemblingResources) HasExecutorCpuRequest() bool {
	if o != nil && o.ExecutorCpuRequest != nil {
		return true
	}

	return false
}

// SetExecutorCpuRequest gets a reference to the given string and assigns it to the ExecutorCpuRequest field.
func (o *EnsemblingResources) SetExecutorCpuRequest(v string) {
	o.ExecutorCpuRequest = &v
}

// GetExecutorMemoryRequest returns the ExecutorMemoryRequest field value if set, zero value otherwise.
func (o *EnsemblingResources) GetExecutorMemoryRequest() string {
	if o == nil || o.ExecutorMemoryRequest == nil {
		var ret string
		return ret
	}
	return *o.ExecutorMemoryRequest
}

// GetExecutorMemoryRequestOk returns a tuple with the ExecutorMemoryRequest field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EnsemblingResources) GetExecutorMemoryRequestOk() (*string, bool) {
	if o == nil || o.ExecutorMemoryRequest == nil {
		return nil, false
	}
	return o.ExecutorMemoryRequest, true
}

// HasExecutorMemoryRequest returns a boolean if a field has been set.
func (o *EnsemblingResources) HasExecutorMemoryRequest() bool {
	if o != nil && o.ExecutorMemoryRequest != nil {
		return true
	}

	return false
}

// SetExecutorMemoryRequest gets a reference to the given string and assigns it to the ExecutorMemoryRequest field.
func (o *EnsemblingResources) SetExecutorMemoryRequest(v string) {
	o.ExecutorMemoryRequest = &v
}

func (o EnsemblingResources) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.DriverCpuRequest != nil {
		toSerialize["driver_cpu_request"] = o.DriverCpuRequest
	}
	if o.DriverMemoryRequest != nil {
		toSerialize["driver_memory_request"] = o.DriverMemoryRequest
	}
	if o.ExecutorReplica != nil {
		toSerialize["executor_replica"] = o.ExecutorReplica
	}
	if o.ExecutorCpuRequest != nil {
		toSerialize["executor_cpu_request"] = o.ExecutorCpuRequest
	}
	if o.ExecutorMemoryRequest != nil {
		toSerialize["executor_memory_request"] = o.ExecutorMemoryRequest
	}
	return json.Marshal(toSerialize)
}

type NullableEnsemblingResources struct {
	value *EnsemblingResources
	isSet bool
}

func (v NullableEnsemblingResources) Get() *EnsemblingResources {
	return v.value
}

func (v *NullableEnsemblingResources) Set(val *EnsemblingResources) {
	v.value = val
	v.isSet = true
}

func (v NullableEnsemblingResources) IsSet() bool {
	return v.isSet
}

func (v *NullableEnsemblingResources) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableEnsemblingResources(val *EnsemblingResources) *NullableEnsemblingResources {
	return &NullableEnsemblingResources{value: val, isSet: true}
}

func (v NullableEnsemblingResources) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableEnsemblingResources) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


