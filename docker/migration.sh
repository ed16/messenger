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
    status VARCHAR(1),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    password_hash CHAR(60) -- CHAR(60) is suitable for bcrypt hashes
  );"
  echo "Table 'users' created."
else
  echo "Table 'users' already exists."
fi

# Create the 'contacts' table if it does not exist
TABLE_EXISTS=$(psql -h $DB_HOST -U $DB_USER -d $DB_NAME -tAc "SELECT 1 FROM information_schema.tables WHERE table_name='contacts'")

if [ "$TABLE_EXISTS" != "1" ]; then
  echo "Table 'contacts' does not exist. Creating..."
  psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "
  CREATE TABLE contacts (
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    contact_user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, contact_user_id)
  );"

  # Create index for quick retrieval of all contacts for a specific user
  psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "
  CREATE INDEX idx_user_id ON contacts (user_id);"

  echo "Table 'contacts' and index created."
else
  echo "Table 'contacts' already exists."
fi

# Upsert the admin user record
USERNAME="admin"
PASSWORD_HASH="\$2a\$10\$p7X62PHGUAGFnhdBDLFjs.ufDZY.59FbWlrBi1PxG4OKlHEb.lTVO"
STATUS="1"

echo "Upserting user $USERNAME..."
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "
INSERT INTO users (username, password_hash, status) 
VALUES ('$USERNAME', '$PASSWORD_HASH', '$STATUS') 
ON CONFLICT (username) 
DO 
  UPDATE SET password_hash = EXCLUDED.password_hash, status = EXCLUDED.status;"
echo "User $USERNAME upserted."
