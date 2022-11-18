#!/bin/bash

docker build -t adhoc-app-0 .

echo 'docker build finished'

docker tag adhoc-app-0 jingyusu/adhoc-app-0:demo

docker login

docker push jingyusu/adhoc-app-0:demo