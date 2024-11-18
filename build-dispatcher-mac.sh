#!/bin/bash

# Verifica si se ha pasado la variable SECRET
if [ -z "$SECRET" ]; then
  echo "La variable de entorno SECRET no está definida. Defínela antes de ejecutar el script."
  exit 1
fi

# Nombre de la imagen
IMAGE_NAME="vottun/bridge-dispatcher/mac-arm64"

# Ruta y parámetros que se pasan como argumentos
ARGS="--yaml /otp/q/properties/bot-int.yaml -p 2116 --env 0 -L TRACE -s"

# Establece el directorio del proyecto (modifica esta ruta si es necesario)
# PROJECT_ROOT="./"

# Cambiar al directorio raíz del proyecto
# cd $PROJECT_ROOT || { echo "No se pudo cambiar al directorio del proyecto"; exit 1; }

# Construir la imagen Docker con la variable de entorno SECRET
docker build -f ./dockerfile.dispatcher.arm64 -t $IMAGE_NAME --build-arg SECRET=$SECRET .

# Imprimir estado de la construcción
if [ $? -eq 0 ]; then
  echo "Imagen $IMAGE_NAME construida con éxito."
else
  echo "Error al construir la imagen."
  exit 1
fi

docker run --name bridge-dispatcher-dev --network bridge-network-dev --restart unless-stopped -e SECRET=$SECRET -d vottun/bridge-dispatcher/mac-arm64