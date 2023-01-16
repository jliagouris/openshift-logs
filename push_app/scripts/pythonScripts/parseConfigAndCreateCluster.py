import yaml
import sys
import subprocess

with open(sys.path[0] + "/../../rosa_config.yaml", "r") as yamlfile:
    ocDeplConfig = yaml.load(yamlfile, Loader=yaml.FullLoader)
    print("Machine Config Read successful")
    yamlfile.close()

print(ocDeplConfig)

command = "rosa create cluster --cluster-name rosa-client-" + str(sys.argv[1])

for key, value in ocDeplConfig.items():
    command += " "
    command += "--"
    command += key
    command += " "
    command += str(value)

command += " --sts --mode auto --yes"
print(command)
commandArgs = command.split(" ")
print(commandArgs)
subprocess.run(commandArgs)
