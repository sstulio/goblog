#!/bin/bash
export CGO_ENABLED=0

echo "Build productservice binary file"
go build -o productservice

echo "Building product service docker image"
docker build -t productservice .

echo "Deleting previous productservice deployment"
kubectl delete deployment productservice

echo "Creting productservice deployment"
kubectl create -f kubernetes/deployments/productservice.yaml


