#!/bin/bash

python3 scripts/pythonScripts/sucDeplYamlProcess.py $1
oc apply -f openshift-deploy.yaml