#!/bin/bash

docker build -t prometheus-app .

echo 'docker build finished'

docker tag prometheus-app jingyusu/prometheus-app:1.0.1

docker login

docker push jingyusu/prometheus-app:1.0.1