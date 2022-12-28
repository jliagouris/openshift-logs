#!/bin/bash

echo 'wait 10 min for cluster to initialize'
sleep 600
rosa create admin --cluster=rosa-client > createAdmin.txt
admin_login_cmd=`python3 pythonScripts/parseAdminLogin.py`
echo $admin_login_cmd
echo 'wait 15 min for id server to get ready'
sleep 900
admin_login_output=`$admin_login_cmd`
echo $admin_login_output
echo "Create Sucrose project"
oc new-project sucrose-app --description="Prometheus scrape app with secrecy: Sucrose"
oc project sucrose-app
oc adm policy add-cluster-role-to-user cluster-reader system:serviceaccount:sucrose-app:default
echo 'wait 1 min for cluster to initialize new role'
sleep 60
service_acc_info=`oc describe sa default`
echo $service_acc_info
token_name=`python3 pythonScripts/parseSA.py $service_acc_info`
echo $token_name
secret_info=`oc describe secret $token_name`
token=`python3 pythonScripts/parseToken.py $secret_info`
echo 'token:'
echo $token
echo $token > ../credentials/token.txt
echo '\n'
echo 'prometheus url:\n'
prom_url=`oc -n openshift-monitoring get routes | grep prometheus-k8s-openshift-monitoring | awk '{print $2}'`
echo $prom_url
echo $prom_url > ../credentials/prom_url.txt