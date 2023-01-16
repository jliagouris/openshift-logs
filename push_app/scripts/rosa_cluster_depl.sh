#!/bin/bash

client_id=$1
echo "creating rosa-client-$client_id"
rosa create account-roles --mode auto --yes
rosa create cluster --cluster-name rosa-client-$client_id --sts --mode auto --yes

# wait until cluster is successfully deployed
cluster_info=`rosa describe cluster -c rosa-client-$client_id`
cluster_state=`python3 scripts/pythonScripts/parseClusterState.py $cluster_info`
while [ "$cluster_state" != "ready" ]
do
    echo "cluster $1 is $cluster_state, check after 5m"
    sleep 300
    cluster_info=`rosa describe cluster -c rosa-client-$client_id`
    cluster_state=`python3 scripts/pythonScripts/parseClusterState.py $cluster_info`
done

echo "Cluster installation finished"
echo 'wait 10 min for cluster to initialize'
sleep 600
rosa create admin --cluster=rosa-client-$client_id > scripts/createAdmin-rosa-$client_id.txt
admin_login_cmd=`python3 scripts/pythonScripts/parseAdminLogin.py rosa $client_id`
echo $admin_login_cmd
echo 'wait 15 min for id server to get ready'
sleep 900
admin_login_output=`$admin_login_cmd`
echo $admin_login_output

./scripts/rosa_cluster_setup.sh rosa $client_id

