# Parameter definitions

This file provides descriptions for all the configurable parameters in the application.

## Global parameters

_`github`_

Required when cloning from Vivek's (private) repository to deploy the benchmark applications. Not used when cloning from public repositories.

* `username` : GitHub username
* `access_token` : GitHub access token. For more info on how to create an access token for your GitHub account check [here](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token).

_`benchmark_application`_

Contains choice of benchmark application and the parameters necessary to run the benchmark application itself.

* `application_choice` : Valid values are ["online-boutique", "hotel-reservation"]. Default value is "online-boutique". 

  * `online-boutique` : Deploys the Online Boutique microservices demo application with Locust load generator.
  * `hotel-reservation` : Deploys the Hotel Reservation microservices demo application with wrk2 load generator. _If deploying this application, wrk2 load generator parameters (below) have to be specified._

* `wrk2`: The following parameters correspond to the wrk2 load generator (not the application)

  * `num_threads` : Number of threads to use
  * `num_connections` : Number of HTTP connections to keep open
  * `benchmark_duration` : Duration of the test
  * `total_rps` : Number of constant HTTP requests per second (to simulate sustained load)


## Cloud provider parameters

_`cloud_providers`_

The parameters specific to each cloud provider in a multi-cloud benchmark setup. At this time, benchmark deployment is supported only on Chameleon Cloud. AWS integration will be done soon.

* `cc`
  
    Contains parameters required for instance creation on Chemeleon Cloud.

  * `lease_name` : Lease name for creating reservations for instances on Chemeleon Cloud.
  * `ssh_key_pair_name` : Name for public and private key pair to SSH into the created instances. The public key is stored on Chemeleon Cloud and the private key will be stored locally.
  * `instance_name` : A valid name for the cloud instances (single name is enough for bulk creation of instances)
  * `instance_flavor` : Machine type to provision. At this time, only "baremetal" type instances are supported.
  * `default_instance_user` : Linux username when we login into the created instances. Default value is "cc" when using a Chameleon Cloud supported Linux image.
  * `cluster` : Contains parameters required to create multiple k8s clusters with desired cardinality.
    * `num_clusters` : Total number of k8s clusters to be created
    * `num_nodes_per_cluster` : Number of cloud instances per cluster to be allocated


* `aws`

    Contains parameters required for instance creation on Amazon Web Services.


* `gcp`

    Contains parameters required for instance creation on Google Cloud Platform.


* `azure` 

    Contains parameters required for instance creation on Azure cloud.


* `nerc`
  
    Contains parameters required for instance creation on New England Research Cloud.