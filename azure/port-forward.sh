#!/bin/bash

k port-forward service/articles-ms 8080:8080 &
k port-forward service/payment-ms 5000:5000 &
k port-forward service/orders-ms 8081:8081 &