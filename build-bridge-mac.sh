#!/bin/bash

# Check if the SECRET variable is passed
if [ -z "$SECRET" ]; then
  echo "The SECRET environment variable is not defined. Define it before running the script."
  exit 1
fi

# Image name
IMAGE_NAME="vottun/qubic-bridge/mac-arm64"

# Path and parameters passed as arguments
ARGS="--yaml /otp/q/properties/bot-int.yaml -p 2116 --env 0 -L TRACE -s"

# Build the Docker image with the SECRET environment variable
docker build -f ./dockerfile.bridge.arm64 -t $IMAGE_NAME --build-arg SECRET=$SECRET .

# Print build status
if [ $? -eq 0 ]; then
  echo "Image $IMAGE_NAME built successfully."
else
  echo "Error building the image."
  exit 1
fi

docker run -p 2116:2116 -p 2126:2126 -p 50051:50051 --name qubic-bridge-dev --network bridge-network-dev --restart unless-stopped -e SECRET=$SECRET -d vottun/qubic-bridge/mac-arm64 
