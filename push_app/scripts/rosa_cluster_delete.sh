#!/bin/bash

rosa delete cluster -c rosa-client
client_info=`rosa describe cluster -c rosa-client`
while [ "$client_info" != "ready" ]
do
    echo "cluster is $cluster_state, check after 5m"
    sleep 300
    cluster_info=`rosa describe cluster -c rosa-client`
    cluster_state=`python3 pythonScripts/parseClusterState.py $cluster_info`
done
