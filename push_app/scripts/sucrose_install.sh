#!/bin/bash

python3 pythonScripts/sucDeplYamlProcess.py
oc apply -f ../openshift-deploy.yaml