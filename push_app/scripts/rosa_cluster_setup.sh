#!/bin/bash

cluster_provider=$1
client_id=$2
echo "Create Sucrose project"
oc new-project sucrose-app --description="Prometheus scrape app with secrecy: Sucrose"
oc project sucrose-app
oc adm policy add-cluster-role-to-user cluster-reader system:serviceaccount:sucrose-app:default
echo 'wait 1 min for cluster to initialize new role'
sleep 60
service_acc_info=`oc describe sa default`
echo $service_acc_info
token_name=`python3 scripts/pythonScripts/parseSA.py $service_acc_info`
echo $token_name
secret_info=`oc describe secret $token_name`
token=`python3 scripts/pythonScripts/parseToken.py $secret_info`
echo 'token:'
echo $token
echo $token > credentials/token-$cluster_provider-$client_id.txt
echo '\n'
echo 'prometheus url:\n'
prom_url=`oc -n openshift-monitoring get routes | grep prometheus-k8s-openshift-monitoring | awk '{print $2}'`
echo $prom_url
echo $prom_url > credentials/prom_url-$cluster_provider-$client_id.txt

./scripts/sucrose_install.sh