#!/bin/bash
export CGO_ENABLED=0

echo "Build userservice binary file"
go build -o userservice

echo "Building user service docker image"
docker build -t userservice .

echo "Deleting previous userservice deployment"
kubectl delete deployment userservice

echo "Creting userservice deployment"
kubectl create -f kubernetes/deployments/userservice.yaml


