#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

export PGPASSWORD=${POSTGRES_PASSWORD}

echo "Waiting for database to be ready..."
timeout 30 bash -c 'until pg_isready -h db -p 5432 -U ${POSTGRES_USER}; do sleep 1; done'

# Check if RESET_DB is true
if [ "${RESET_DB}" != "true" ]; then
    echo "RESET_DB is false. Skipping setup."
    exit 0
fi

echo "Creating bdjuno database if it does not exist..."
if ! psql -h db -U ${POSTGRES_USER} -tc "SELECT 1 FROM pg_database WHERE datname = 'bdjuno'" | grep -q 1; then
    psql -h db -U ${POSTGRES_USER} -c "CREATE DATABASE bdjuno"
fi

echo "Creating bdjuno user if it does not exist..."
if ! psql -h db -U ${POSTGRES_USER} -tc "SELECT 1 FROM pg_roles WHERE rolname = '${BDJUNO_USER}'" | grep -q 1; then
    psql -h db -U ${POSTGRES_USER} -d bdjuno -c "CREATE USER ${BDJUNO_USER} WITH ENCRYPTED PASSWORD '${BDJUNO_PASSWORD}'"
fi

echo "Granting privileges to bdjuno user..."
psql -h db -U ${POSTGRES_USER} -d bdjuno -c "GRANT ALL PRIVILEGES ON DATABASE bdjuno TO ${BDJUNO_USER}"

echo "Granting privileges on public schema..."
psql -h db -U ${POSTGRES_USER} -d bdjuno -c "GRANT USAGE, CREATE ON SCHEMA public TO ${BDJUNO_USER}"
psql -h db -U ${POSTGRES_USER} -d bdjuno -c "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO ${BDJUNO_USER}"

echo "Running SQL schema files..."
if [ -d "/schema" ] && [ "$(ls -A /schema/*.sql 2>/dev/null)" ]; then
    for file in /schema/*.sql; do
        echo "Running $file..."
        PGPASSWORD=${BDJUNO_PASSWORD} psql -v ON_ERROR_STOP=1 -h db -U ${BDJUNO_USER} -d bdjuno -f "$file"
    done
else
    echo "No SQL files found in /schema"
    exit 1
fi

echo "Setup completed!"
