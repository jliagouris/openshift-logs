#!/bin/bash

rosa create account-roles --mode auto --yes
rosa create cluster --cluster-name rosa-client --sts --mode auto --yes

# wait until cluster is successfully deployed
cluster_info=`rosa describe cluster -c rosa-client`
cluster_ready=`python3 pythonScripts/parseClusterState.py $cluster_info`
while [ "$cluster_ready" != "true" ]
do
    echo "cluster is not ready, check after 5m"
    sleep 300
    cluster_info=`rosa describe cluster -c rosa-client`
    cluster_ready=`python3 pythonScripts/parseClusterState.py $cluster_info`
done

echo "Cluster installation finished"

#./cluster_setup.sh

