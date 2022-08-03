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

2. [HighOverallControlPlaneCPU](https://github.com/openshift/cluster-kube-apiserver-operator/blob/master/bindata/assets/alerts/cpu-utilization.yaml)
   QueryRule:
   ```
   sum(
    100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[1m])) * 100) 
    and 
    on(instance) label_replace(kube_node_role{role="master"}, "instance", "$1", "node", "(.+)")) / count(kube_node_role{role="master"}
   ) > 60
   ```
   Description: More than 60% of ControlPlane(master) CPUs are busy. CPU utilization across all three control plane nodes is higher than two control plane nodes can sustain; a single control plane node outage may cause a cascading failure; increase available CPU.  
   Why 60%: Given three control plane nodes, the overall CPU utilization may only be about 2/3 of all available capacity. This is because if a single control plane node fails, the remaining two must handle the load of the cluster in order to be HA. If the cluster is using more than 2/3 of all capacity, if one control plane node fails, the remaining two are likely to fail when they take the load. To fix this, increase the CPU and memory on your control plane nodes.

   For: 10m

3. [KubePodNotReady](https://github.com/openshift/cluster-monitoring-operator/blob/aefc8fc5fc61c943dc1ca24b8c151940ae5f8f1c/assets/control-plane/prometheus-rule.yaml#L440-L449)  
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

   Example metric entries:
   ```
    kube_pod_status_phase{cluster="moc/smaug", container="kube-rbac-proxy-main", endpoint="https-main", job="kube-state-metrics", namespace="acme-operator", phase="Failed", pod="openshift-acme-54496d7cc7-hvfgf", prometheus="openshift-monitoring/k8s", service="kube-state-metrics"} = 1

    kube_pod_owner{cluster="moc/smaug", container="kube-rbac-proxy-main", endpoint="https-main", job="kube-state-metrics", namespace="acme-operator", owner_is_controller="true", owner_kind="ReplicaSet", owner_name="openshift-acme-54496d7cc7", pod="openshift-acme-54496d7cc7-hvfgf", prometheus="openshift-monitoring/k8s", service="kube-state-metrics"} = 1
   ```
   The query
   ```
   max by(namespace, pod) (
        kube_pod_status_phase{namespace=~"(openshift-.*|kube-.*|default|logging)",job="kube-state-metrics", phase=~"Pending|Unknown"}
    )
   ```
   is equivalent to the following SQL query:
   ```sql
   SELECT MAX(value) AS value, namespace, pod
   FROM kube_pod_status_phase
   WHERE namespace =~ "(openshift-.*|kube-.*|default|logging)" AND job = "kube-state-metrics" AND phase =~ "Pending|Unknown"
   GROUP BY namespace, pod;
   ```
   Gives: 
   ```
   {namespace="openshift-authentication", pod="oauth-openshift-84db7cff56-sfdrm"} = 1
   ```
   Similarly, 
   ```
   topk by(namespace, pod) (
        1, max by(namespace, pod, owner_kind) (kube_pod_owner{owner_kind!="Job"})
    )
   ```
   is equivalent to
   ```sql
   SELECT *
   FROM (
            SELECT MAX(value) AS value, namespace, pod, ownerkind
            FROM kube_pod_owner
            WHERE owner_kind != "Job"
            GROUP BY namespace, pod, ownerkind
        )
   GROUP BY namespace, pod
   ORDER BY value DESC
   LIMIT 1
   ```
   Returns
   ```
   {namespace="openshift-authentication", owner_kind="ReplicaSet", pod="oauth-openshift-84db7cff56-sfdrm"} = 1
   ```
   *on(namespace, pod) group_left(owner_kind)* is like left join in SQL, on label namespace and pod. Attribute owner_kind attaches the owner_kind label to the multiplied result. So 
   ```
    max by(namespace, pod) (
        kube_pod_status_phase{namespace=~"(openshift-.*|kube-.*|default|logging)",job="kube-state-metrics", phase=~"Pending|Unknown"}
    ) 
    * on(namespace, pod) group_left(owner_kind) 
        topk by(namespace, pod) (
        1, max by(namespace, pod, owner_kind) (kube_pod_owner{owner_kind!="Job"})
    )
   ```
   is equivalent to:
   ```sql
   SELECT (T1.value * T2.value) AS value, T1.namespace, T1.pod, T2.owner_kind
   FROM (
           SELECT MAX(value) AS value, namespace, pod
           FROM kube_pod_status_phase
           WHERE namespace =~ "(openshift-.*|kube-.*|default|logging)" AND job = "kube-state-metrics" AND phase =~ "Pending|Unknown"
           GROUP BY namespace, pod
        ) AS T1, 
        (
           SELECT *
           FROM (
                   SELECT MAX(value) AS value, namespace, pod, ownerkind
                   FROM kube_pod_owner
                   WHERE owner_kind != "Job"
                   GROUP BY namespace, pod, ownerkind
                )
           GROUP BY namespace, pod
           ORDER BY value DESC
           LIMIT 1
        ) AS T2
   WHERE T1.namespace = T2.namespace AND T1.pod = T2.pod
   ```
   result is:
   ```
    {namespace="openshift-authentication", owner_kind="ReplicaSet", pod="oauth-openshift-84db7cff56-f92kv"} = 1
   ```
   Sum over all filtered namespaces and pods, if the result is larger than 0, that means some pod is not ready.

   Pseudo SQL Equivalent:
    ```sql
    SELECT SUM(T3.value) AS sum, T3.namespace, T3.pod, T3.Time
    FROM (
        SELECT (T1.value * T2.value) AS value, T1.namespace, T1.pod, T2.owner_kind
        FROM (
                SELECT MAX(value) AS value, namespace, pod
                FROM kube_pod_status_phase
                WHERE namespace =~ "(openshift-.*|kube-.*|default|logging)" AND job = "kube-state-metrics" AND phase =~ "Pending|Unknown"
                GROUP BY namespace, pod
             ) AS T1, 
             (
                SELECT *
                FROM (
                        SELECT MAX(value) AS value, namespace, pod, ownerkind
                        FROM kube_pod_owner
                        WHERE owner_kind != "Job"
                        GROUP BY namespace, pod, ownerkind
                    )
                GROUP BY namespace, pod
                ORDER BY value DESC
                LIMIT 1
             ) AS T2
        WHERE T1.namespace = T2.namespace AND T1.pod = T2.pod
    ) AS T3
    WHERE sum > 0
    GROUP BY T3.namespace, T3.pod;
    ```

4. [etcdDatabaseQuotaLowSpace](https://github.com/openshift/runbooks/blob/master/alerts/cluster-etcd-operator/etcdDatabaseQuotaLowSpace.md)  
   Query Rule:
   ```
   (etcd_mvcc_db_total_size_in_bytes / etcd_server_quota_backend_bytes) * 100 > 95
   ```
   Description: This alert fires when the total existing DB size exceeds 95% of the maximum DB quota. The consumed space is in Prometheus represented by the metric etcd_mvcc_db_total_size_in_bytes, and the DB quota size is defined by etcd_server_quota_backend_bytes. ([etcdDatabaseQuotaLowSpace](https://github.com/openshift/runbooks/blob/master/alerts/cluster-etcd-operator/etcdDatabaseQuotaLowSpace.md))

   For: 10m

5. [KubePodCrashLooping](https://github.com/openshift/cluster-monitoring-operator/blob/aefc8fc5fc61c943dc1ca24b8c151940ae5f8f1c/assets/control-plane/prometheus-rule.yaml#L440-L449)
   Query Rule:
   ```
   rate(kube_pod_container_status_restarts_total{namespace=~"(openshift-.*|kube-.*|default|logging)",job="kube-state-metrics"}[10m]) * 60 * 5 > 0
   ```
   Description:  
   The pod container restarts more than 50% of the time during the last 10 minutes  

   For: 15m

6. [KubeDaemonSetNotScheduled](https://github.com/openshift/cluster-monitoring-operator/blob/aefc8fc5fc61c943dc1ca24b8c151940ae5f8f1c/assets/control-plane/prometheus-rule.yaml#L186-L195)
   Query Rule:
   ```
    kube_daemonset_status_desired_number_scheduled{namespace=~"(openshift-.*|kube-.*|default|logging)",job="kube-state-metrics"}
    -
    kube_daemonset_status_current_number_scheduled{namespace=~"(openshift-.*|kube-.*|default|logging)",job="kube-state-metrics"} 
    > 0
   ```
   Description:
   Desired number of daemonset larger than scheduled, means some are not scheduled.

   For: 10m

7. [KubeJobFailed](https://github.com/openshift/cluster-monitoring-operator/blob/aefc8fc5fc61c943dc1ca24b8c151940ae5f8f1c/assets/control-plane/prometheus-rule.yaml#L186-L195)
   Query Rule:
   ```
   kube_job_failed{namespace=~"(openshift-.*|kube-.*|default|logging)",job="kube-state-metrics"}  > 0
   ```
   Description: This one is self-explainatory.

   For: 15m