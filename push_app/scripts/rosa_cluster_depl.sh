#!/bin/bash

rosa create account-roles --mode auto --yes
rosa create cluster --cluster-name rosa-client --sts --mode auto --yes

# wait until cluster is successfully deployed
cluster_info=`rosa describe cluster -c rosa-client`
cluster_state=`python3 scripts/pythonScripts/parseClusterState.py $cluster_info`
while [ "$cluster_state" != "ready" ]
do
    echo "cluster is $cluster_state, check after 5m"
    sleep 300
    cluster_info=`rosa describe cluster -c rosa-client`
    cluster_state=`python3 scripts/pythonScripts/parseClusterState.py $cluster_info`
done

echo "Cluster installation finished"

./rosa_cluster_setup.sh

