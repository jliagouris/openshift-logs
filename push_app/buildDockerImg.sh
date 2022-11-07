#!/bin/bash

docker build -t prometheus-app .

echo 'docker build finished'

docker tag prometheus-app jingyusu/adhoc-app:1.0.0

docker login

docker push jingyusu/adhoc-app:1.0.0