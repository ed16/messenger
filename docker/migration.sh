#!/bin/sh

set -e

# Variables
DB_HOST=${POSTGRES_HOST}
DB_PORT=${POSTGRES_PORT}
DB_USER=${POSTGRES_USER}
DB_PASS=${POSTGRES_PASSWORD}
DB_NAME=${POSTGRES_DB}

# Wait for the database to be ready using goose
until goose -dir /migrations postgres "host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASS dbname=$DB_NAME sslmode=disable" status &>/dev/null; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - running goose migrations..."

# Run the goose migrations
goose -dir /migrations postgres "host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASS dbname=$DB_NAME sslmode=disable" up
echo "--- Messenger is now ready for you. Enjoy and have fun! ---"