#!/bin/bash

version=$1

echo "Loading docker container for go-stock_$version..."
docker load --input go-stock_$version.tar.gz
echo "Done!"
echo "Starting docker container for go-stock_$version..."
docker run -d --env-file ./docker.env --net=host --name=go-stock go-stock:$version
echo "Done!"