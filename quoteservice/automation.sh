#!/bin/bash
export GOOS=linux
export CGO_ENABLED=0

echo "Build quoteservice binary file"
go build -o quoteservice-linux-amd64;

export GOOS=darwin

echo "Building quote service docker image"
docker build -t quoteservice .

echo "Deleting previous quoteservice deployment"
kubectl delete deployment quoteservice

echo "Creting quoteservice deployment"
kubectl create -f kubernetes/deployments/quoteservice.yaml


