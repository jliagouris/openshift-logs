import yaml
import os
import sys

print(os. getcwd())
with open(sys.path[0] + "/../../openshift-deploy.yaml", "r") as yamlfile:
    ocDeplConfig = yaml.load(yamlfile, Loader=yaml.FullLoader)
    print("Read successful")
    yamlfile.close()
deplTemplateSpecContainerEnvVar = ocDeplConfig['spec']['template']['spec']['containers'][0]['env']
URL_IDX, TOKEN_IDX = 0, 1
with open(sys.path[0] + "/../../credentials/prom_url.txt", "r") as urlfile:
    prom_url = urlfile.read()
    urlfile.close()
deplTemplateSpecContainerEnvVar[URL_IDX]['value'] = prom_url[:-1]
with open(sys.path[0] + "/../../credentials/token.txt", "r") as tokenfile:
    token = tokenfile.read()
    tokenfile.close()
deplTemplateSpecContainerEnvVar[TOKEN_IDX]['value'] = token[:-1]
with open(sys.path[0] + "/../../openshift-deploy.yaml", "w") as yamlfile:
    ocDeplConfig = yaml.dump(ocDeplConfig, yamlfile)
    print(ocDeplConfig)
    print("yaml file write successful")
    yamlfile.close()




