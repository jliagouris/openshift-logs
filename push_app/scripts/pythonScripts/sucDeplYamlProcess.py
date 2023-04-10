import yaml
import os
import sys

URL_IDX, TOKEN_IDX, CLIENT_ID_IDX = 0, 1, 2

print(os. getcwd())
with open(sys.path[0] + "/../../openshift-deploy.yaml", "r") as yamlfile:
    ocDeplConfig = yaml.load(yamlfile, Loader=yaml.FullLoader)
    print("Read successful")
    yamlfile.close()
deplTemplateSpecContainerEnvVar = ocDeplConfig['spec']['template']['spec']['containers'][0]['env']
with open(sys.path[0] + "/../../credentials/prom_url-" + sys.argv[1] + "-" + str(sys.argv[2]) + ".txt", "r") as urlfile:
    prom_url = urlfile.read()
    urlfile.close()
deplTemplateSpecContainerEnvVar[URL_IDX]['value'] = prom_url[:-1]

with open(sys.path[0] + "/../../credentials/token-" + sys.argv[1] + "-" + str(sys.argv[2]) + ".txt", "r") as tokenfile:
    token = tokenfile.read()
    tokenfile.close()
deplTemplateSpecContainerEnvVar[TOKEN_IDX]['value'] = token[:-1]

deplTemplateSpecContainerEnvVar[CLIENT_ID_IDX]['value'] = sys.argv[2]

with open(sys.path[0] + "/../../openshift-deploy.yaml", "w") as yamlfile:
    ocDeplConfig = yaml.dump(ocDeplConfig, yamlfile)
    print(ocDeplConfig)
    print("yaml file write successful")
    yamlfile.close()




