#!/bin/bash

rosa create account-roles --mode auto --yes
rosa create cluster --cluster-name rosa-client --sts --mode auto --yes

