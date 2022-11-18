#!/bin/bash

admin_create_result=`rosa create admin --cluster=client-0`
#admin_create_result="INFO: Admin account has been added to cluster 'log-analysis'. INFO: Please securely store this generated password. If you lose this password you can delete and recreate the cluster admin user. INFO: To login, run the following command: oc login https://api.log-analysis.f670.p1.openshiftapps.com:6443 --username cluster-admin --password Iqauq-rTW2D-2SWih-sVtgG INFO: It may take up to a minute for the account to become active."
echo "admin_create_result: $admin_create_result"
#info_head="INFO: Admin account has been added to cluster 'log-analysis'. INFO: Please securely store this generated password. If you lose this password you can delete and recreate the cluster admin user. INFO: To login, run the following command: "
#admin_create_result=admin_create_result#$info_head
#echo $admin_create_result
#info_tail=" INFO: It may take up to a minute for the account to become active."
#admin_login_cmd=`admin_create_result%$info_tail`
#echo $admin_login_cmd >> admin_login_cmd.txt
#admin_login_result=`$admin_login_cmd`
#echo $admin_login_result