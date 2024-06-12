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

# Create the 'users' table if it does not exist
TABLE_EXISTS=$(psql -h $DB_HOST -U $DB_USER -d $DB_NAME -tc "SELECT 1 FROM information_schema.tables WHERE table_name = 'users'" | tr -d '[:space:]')

if [ "$TABLE_EXISTS" != "1" ]; then
  echo "Table 'users' does not exist. Creating..."
  psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "
  CREATE TYPE user_status AS ENUM ('active', 'deleted', 'blocked');
  CREATE TABLE users (
    user_id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    status user_status, -- Use the enum type here
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
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    contact_user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
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

# Create the 'media' table if it does not exist
TABLE_EXISTS=$(psql -h $DB_HOST -U $DB_USER -d $DB_NAME -tAc "SELECT 1 FROM information_schema.tables WHERE table_name='media'")

if [ "$TABLE_EXISTS" != "1" ]; then
  echo "Table 'media' does not exist. Creating..."
  psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "
  CREATE TABLE media (
    media_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    file_path VARCHAR(260) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
  );"

  # Create index for quick retrieval of all media for a specific user
  psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "
  CREATE INDEX idx_user_id ON media (user_id);"

  echo "Table 'media' and index created."
else
  echo "Table 'media' already exists."
fi

# Create the 'messages' table if it does not exist
TABLE_EXISTS=$(psql -h $DB_HOST -U $DB_USER -d $DB_NAME -tAc "SELECT 1 FROM information_schema.tables WHERE table_name='messages'")

if [ "$TABLE_EXISTS" != "1" ]; then
  echo "Table 'messages' does not exist. Creating..."
  psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "
  CREATE TYPE message_status AS ENUM ('sent', 'received', 'read');
  CREATE TABLE messages (
      message_id BIGSERIAL PRIMARY KEY,
      sender_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
      recipient_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
      content VARCHAR(1000) NOT NULL,
      created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
      status message_status
      media_id BIGINT,
      CONSTRAINT fk_media FOREIGN KEY (media_id) REFERENCES media(media_id) ON DELETE RESTRICT
  );

  # Create index for quick retrieval of all messages for a specific user
  psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "
  CREATE INDEX idx_sender_id ON messages (sender_id);"

  echo "Table 'messages' and index created."
else
  echo "Table 'messages' already exists."
fi

# Upsert the admin user record
USERNAME="admin"
PASSWORD_HASH="\$2a\$10\$p7X62PHGUAGFnhdBDLFjs.ufDZY.59FbWlrBi1PxG4OKlHEb.lTVO"
STATUS="active"

echo "Upserting user $USERNAME..."
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "
INSERT INTO users (username, password_hash, status) 
VALUES ('$USERNAME', '$PASSWORD_HASH', '$STATUS') 
ON CONFLICT (username) 
DO 
  UPDATE SET password_hash = EXCLUDED.password_hash, status = EXCLUDED.status;"
echo "User $USERNAME upserted."
