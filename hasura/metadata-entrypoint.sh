#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

apt update && apt install -y curl

curl -L https://github.com/hasura/graphql-engine/releases/download/v2.0.4/cli-hasura-linux-amd64 -o /usr/local/bin/hasura
chmod +x /usr/local/bin/hasura

mkdir -p /root/.hasura/cli-ext/v2.0.4/
curl -L -o /root/.hasura/cli-ext/v2.0.4/cli-ext https://github.com/hasura/graphql-engine/releases/download/v2.0.4/cli-ext-linux-amd64
chmod +x /root/.hasura/cli-ext/v2.0.4/cli-ext

cd /hasura

hasura metadata apply --endpoint "${HASURA_GRAPHQL_ENDPOINT}" --admin-secret "${HASURA_GRAPHQL_ADMIN_SECRET}"