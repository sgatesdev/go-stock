#!/bin/bash
#builds and compresses docker image for release

version=$1
docker build . -t go-stock:latest -t go-stock:$version --platform linux/amd64
echo "Saving docker image to ./docker_images/go-stock_$version.tar.gz..."
docker save go-stock:$version | gzip > ./docker_images/go-stock_$version.tar.gz
echo "Done!"

# cache version number
touch ./docker_images/version.txt
echo "Version: $version, Timestamp: $(date)" > ./docker_images/version.txt

# docker rmi $(docker images -q 'go-stock') -fhow