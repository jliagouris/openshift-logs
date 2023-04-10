#!/bin/bash

python3 scripts/pythonScripts/sucDeplYamlProcess.py $1 $2
oc apply -f openshift-deploy.yaml