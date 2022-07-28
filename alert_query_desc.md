1. [ExtremelyHighIndividualControlPlaneCPU](https://github.com/openshift/cluster-kube-apiserver-operator/blob/master/bindata/assets/alerts/cpu-utilization.yaml)  
    Query Rule:
    ```
    100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[1m])) * 100) > 90 
    AND 
    on (instance) label_replace( kube_node_role{role="master"}, "instance", "$1", "node", "(.+)" )
    ```
    Description: The average busy time of CPU among all instances for the past 1 minute is larger than 0.9s/1s, and the instance is a master instance.  
    _on(instance)_ matches the instance label of time series on the left and right. *label_replace()* copies the "node" label to "instance" label of kube_node_role time series.

    For: 5m(warning), 1h(critical)

2. [KubePodNotReady](https://github.com/openshift/cluster-monitoring-operator/blob/aefc8fc5fc61c943dc1ca24b8c151940ae5f8f1c/assets/control-plane/prometheus-rule.yaml#L440-L449)  
   Query Rule:
   ```
   sum by (namespace, pod) (
          max by(namespace, pod) (
            kube_pod_status_phase{namespace=~"(openshift-.*|kube-.*|default|logging)",job="kube-state-metrics", phase=~"Pending|Unknown"}
          ) 
          * on(namespace, pod) group_left(owner_kind) 
          topk by(namespace, pod) (
            1, max by(namespace, pod, owner_kind) (kube_pod_owner{owner_kind!="Job"})
          )
    ) > 0
   ```
   Description: 
   Time series *kube_pod_status_phase{}* has value 0 or 1, depending on if there is a pod in a specific status  
   *kube_pod_status_phase{namespace=~"(openshift-.*|kube-.*|default|logging)",job="kube-state-metrics", phase=~"Pending|Unknown"}* filters out the series that matches the selectors.  
   max by(namespace, pod) (kube_pod_status_phase{namespace=~"(openshift-.*|kube-.*|default|logging)",job="kube-state-metrics", phase=~"Pending|Unknown"}) returns series of all pods that are in pending/unknown state. (kube_pod_status_phase == 1)  
   *topk by(namespace, pod) (1, max by(namespace, pod, owner_kind) (kube_pod_owner{owner_kind!="Job"}))* Gets the maximum kube_pod_owner series that is not owned by "job". *on(namespace, pod) group_left(owner_kind)* is similer to left join in SQL, (I think this *group_left(owner_kind)* is not necessary here).  
   In summary, this PromQL can be translated to plain language as: The number of pending/unknown pods whose owner is not "Job" is greater than 0.  
   For: 15m


3. [etcdDatabaseQuotaLowSpace](https://github.com/openshift/runbooks/blob/master/alerts/cluster-etcd-operator/etcdDatabaseQuotaLowSpace.md)  
   Query Rule:
   ```
   (etcd_mvcc_db_total_size_in_bytes / etcd_server_quota_backend_bytes) * 100 > 95
   ```
   Description: This alert fires when the total existing DB size exceeds 95% of the maximum DB quota. The consumed space is in Prometheus represented by the metric etcd_mvcc_db_total_size_in_bytes, and the DB quota size is defined by etcd_server_quota_backend_bytes. ([etcdDatabaseQuotaLowSpace](https://github.com/openshift/runbooks/blob/master/alerts/cluster-etcd-operator/etcdDatabaseQuotaLowSpace.md))

   For: 10m