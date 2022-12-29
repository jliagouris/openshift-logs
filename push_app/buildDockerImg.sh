#!/bin/bash

docker build -t push_app .

echo 'docker build finished'

docker tag push_app jingyusu/push_app:latest

docker login

docker push jingyusu/push_app:latest
