#!/bin/bash

oc new-project secrecy-app --description="Prometheus scrape app with secrecy"
oc project secrecy-app
oc adm policy add-cluster-role-to-user cluster-reader system:serviceaccount:secrecy-app:default
oc describe sa default
//TODO: get secret name description
oc describe secret <secret name>
//TODO: get token from secret description
//TODO: set env var for app deployment (prometheus route & token)