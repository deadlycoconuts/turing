import Mustache from "mustache";
const moment = require("moment");

export const createLogImageBuilderUrl = (template, cluster, namespace, jobName, startTime) => {
    const endTime = moment(startTime, "YYYY-MM-DDTHH:mm.SSZ").add(1, "hour") // we assume the image builder finished after 1 hour
    const data = {
        cluster_name: cluster,
        namespace_name: namespace,
        job_name: jobName,
        start_time: startTime,
        end_time: endTime.toISOString(),
    };

    return Mustache.render(template, data);
}

export const createLogUrl = (template, cluster, namespace, podPrefixName, startTime) => {
    const data = {
        cluster_name: cluster,
        namespace_name: namespace,
        pod_prefix_name: podPrefixName,
        start_time: startTime,
    };

    return Mustache.render(template, data);
}