#!/bin/bash

cluster_info=`rosa describe cluster -c rosa-client`
#cluster_ready=`python3 pythonScripts/parseClusterState.py $cluster_info`
cluster_ready="installing"
echo $cluster_ready
while [ "$cluster_ready" != "true" ]
do
    echo "cluster is not ready, check after 5m"
    sleep 300
    cluster_info=`rosa describe cluster -c rosa-client`
    cluster_ready=`python3 pythonScripts/parseClusterState.py $cluster_info`
done