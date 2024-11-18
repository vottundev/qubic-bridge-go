#!/bin/bash

export SECRET='s8NvcyW9o6Fz!0AcSaUz#JY1uLK}gG2'

# Nombre de la imagen
IMAGE_NAME="vottun/bridge-dispatcher/linux-amd64"

docker buildx create --use --name linuxBuilder

# build docker image
docker buildx build -f ./dockerfile.dispatcher.arm64 --load --platform linux/amd64/v2 -t $IMAGE_NAME --build-arg SECRET=$SECRET .

# print build status
if [ $? -eq 0 ]; then
  echo "Image $IMAGE_NAME build succesfully."
else
  echo "failed building image."
  exit 1
fi

docker run --name bridge-dispatcher-dev --network bridge-network-dev -e SECRET=$SECRET -d vottun/bridge-dispatcher/mac-arm64