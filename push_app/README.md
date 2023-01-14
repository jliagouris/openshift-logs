# Sucrose - Client-side Application for Secure Cross-Site Log Analytics using Secrecy

This is the client application that lives on the OpenShift servers of RedHat clients.

Its duty is to periodically or on demand collect, preprocess and send desensitified logs data shares to secrecy servers for further calculation.

Its development relies on Golang v1.19, confluent-kafka-go v1.9 and go-yaml v3, make sure they are installed before running the app.

## Assumptions made:
1. All kafka producers have the same timeouts.
2. All producers across clients push to the same topic.

## Service workflow
To support different scenarios, the application supports two service models: periodic push
(prototype up and running), and ad-hoc pull (under development). Graphic illustrations of workflows
of the two models are as follows: 

Periodic Push:
![image](pictures/push.png)

Ad-hoc Pull:
![image](pictures/pull.png)

### Components:
- Data Source: Component that sends queries to Prometheus, recieves response, and sends structured response to components downstream of the pipeline.
- Parser: Parses the response and extracts predefined arguments, in order to reduce workload of downstream components.


## Configurable Parameters:
The application supports a number of configurable parameters, all defined in config.yaml. For 
details, please look at: [Configurable Parameters](docs/config_param.md).

## Cluster & Application Deployment
### Prerequesites:
- General: 
  
  1. Make sure that [Openshift CLI](https://docs.openshift.com/container-platform/4.8/cli_reference/openshift_cli/getting-started-cli.html) is installed.
  2. Install [Docker](https://www.docker.com/), preferably Docker Desktop, which comes with GUI.
  3. This project runs successfully on Linux(Ubuntu) and MacOS. Success of deployment and running on Windows is not guaranteed (most likely it won't).
- ROSA (Redhat Openshift on AWS) cluster:
  
  1. Have [ROSA CLI](https://docs.openshift.com/rosa/rosa_cli/rosa-get-started-cli.html) installed.
  2. Have [AWS CLI 2](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) installed and configured.
  3. Create service linked roles for Elastic Load Balancer(ELB) on your AWS account:
    - Check if the role exists:
     ```Bash
     aws iam get-role --role-name "AWSServiceRoleForElasticLoadBalancing"
     ```
    - If not, run:
     ```Bash
     aws iam create-service-linked-role --aws-service-name "elasticloadbalancing.amazonaws.com"
     ```
  4. Make sure that you have enough resource quota on your account: 
     ```Bash
     rosa verify quota
     ```
     If not, apply for quota increase in your AWS console.
