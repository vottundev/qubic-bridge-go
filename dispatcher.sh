#!/bin/bash

export BRIDGE_SECRET='s8NvcyW9o6Fz!0AcSaUz#JY1uLK}gG2'

ARGS="--yaml=/Users/alexlopez/dev/golang/vottundev/vottun-qubic-bridge-go/properties/qubic-dispatcher-dev.yaml -l=TRACE -s=true --launch=dispatcher --grpc-server-port=50051"

LOG_DIR="/var/log/qubic-bridge"
LOG_FILE="${LOG_DIR}/qubic-bridge.log"

if [[ ! -d "$LOG_DIR" ]]; then
    sudo mkdir -p "$LOG_DIR"
    sudo chown $USER "$LOG_DIR"
fi

./qubic-dispatcher $ARGS >> "$LOG_FILE" 2>&1 &

echo "application running in background. Log: $LOG_FILE"