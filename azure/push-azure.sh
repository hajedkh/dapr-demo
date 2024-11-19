#!/bin/sh

cd articles-ms-java
docker build --platform=linux/amd64 -t daprdemo.azurecr.io/articles-ms .
cd ../orders-ms-golang
docker build --platform=linux/amd64 -t daprdemo.azurecr.io/orders-ms .
cd ../payment-ms-python
docker build --platform=linux/amd64 -t daprdemo.azurecr.io/payment-ms .
cd ..
docker push daprdemo.azurecr.io/articles-ms
docker push daprdemo.azurecr.io/orders-ms
docker push daprdemo.azurecr.io/payment-ms