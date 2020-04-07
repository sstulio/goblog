#!/bin/bash
export GOOS=linux
export CGO_ENABLED=0

echo "Build accountservice binary file"
go build -o accountservice-linux-amd64;

export GOOS=darwin

echo "Building account service docker image"
docker build -t accountservice .

echo "Deleting previous accountservice deployment"
kubectl delete deployment accountservice

echo "Creting accountservice deployment"
kubectl create -f kubernetes/deployments/accountservice.yaml


