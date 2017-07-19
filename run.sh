#!/usr/bin/env bash

go run main.go

echo "Find running go trace services"
sudo netstat -plnt | grep trace


echo "Kill go trace servicesk"
for pid in $(sudo netstat -plnt | grep trace | awk '{split($7,a,"/"); print a[1]}'); do sudo kill -9 $pid; done