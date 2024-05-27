#!/bin/bash

set -e

# Variables
DB_HOST=${POSTGRES_HOST}
DB_PORT=${POSTGRES_PORT}
DB_USER=${POSTGRES_USER}
DB_PASS=${POSTGRES_PASSWORD}
DB_NAME=${POSTGRES_DB}

export PGPASSWORD=$DB_PASS

# Check if the database exists
DB_EXISTS=$(psql -h $DB_HOST -U $DB_USER -d $DB_NAME -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | tr -d '[:space:]')

# Create the database if it does not exist
if [ "$DB_EXISTS" != "1" ]; then
  echo "Database $DB_NAME does not exist. Creating..."
  createdb -h $DB_HOST -U $DB_USER $DB_NAME
  echo "Database $DB_NAME created."
else
  echo "Database $DB_NAME already exists."
fi

# Check if the 'users' table exists
TABLE_EXISTS=$(psql -h $DB_HOST -U $DB_USER -d $DB_NAME -tc "SELECT 1 FROM information_schema.tables WHERE table_name = 'users'" | tr -d '[:space:]')

# Create the 'users' table if it does not exist
if [ "$TABLE_EXISTS" != "1" ]; then
  echo "Table 'users' does not exist. Creating..."
  psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "
  CREATE TABLE users (
    user_id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    password_hash CHAR(60) -- CHAR(60) is suitable for bcrypt hashes
  );"
  echo "Table 'users' created."
else
  echo "Table 'users' already exists."
fi
