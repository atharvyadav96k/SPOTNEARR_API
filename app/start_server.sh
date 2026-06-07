#!/bin/bash

# Exit immediately if any command fails
set -e

echo "Starting infrastructure containers (Postgres & Redis)..."
cd /home/atharv96k/SPOTNEARR_API/app
/usr/bin/docker compose up -d

echo "Waiting for Postgres database network interface to be ready..."
until /usr/bin/docker exec $(/usr/bin/docker ps -qf "name=postgres") pg_isready -U admin -d spotnearr >/dev/null 2>&1; do
    echo "Postgres is still initializing internal processes... waiting 2s"
    sleep 2
done

echo "Infrastructure is fully online!"
echo "Starting Go API Server directly..."
cd /home/atharv96k/SPOTNEARR_API/app/cmd

/usr/local/go/bin/go run main.go