#!/bin/bash

# Use docker build -t push_app-1 . if not on apple m1 chip
docker buildx build --platform linux/amd64 -t push_app .

echo 'docker build finished'

docker tag push_app jingyusu/push_app:latest

docker login

docker push jingyusu/push_app:latest
