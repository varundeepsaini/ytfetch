#!/bin/bash

# Load environment variables from .env
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
else
    echo ".env file not found!"
    exit 1
fi

# Validate required environment variables
if [[ -z "$DB_NAME" || -z "$DB_USER" || -z "$DB_HOST" || -z "$DB_PORT" ]]; then
    echo "Missing one or more required environment variables."
    exit 1
fi

# Confirmation prompt
read -p "Are you sure you want to delete and recreate the database '$DB_NAME'? (yes/no): " confirm

if [[ "$confirm" != "yes" ]]; then
    echo "Operation cancelled."
    exit 0
fi

# Run MySQL command
echo "Dropping and recreating '$DB_NAME'..."

# Handle empty password case with conditional
if [[ -z "$DB_PASSWORD" ]]; then
    mysql -u"$DB_USER" -h "$DB_HOST" -P "$DB_PORT" -e "DROP DATABASE IF EXISTS \`$DB_NAME\`; CREATE DATABASE \`$DB_NAME\`;"
    mysql -u"$DB_USER" -h "$DB_HOST" -P "$DB_PORT" "$DB_NAME" < internal/models/schema.sql
else
    mysql -u"$DB_USER" -p"$DB_PASSWORD" -h "$DB_HOST" -P "$DB_PORT" -e "DROP DATABASE IF EXISTS \`$DB_NAME\`; CREATE DATABASE \`$DB_NAME\`;"
    mysql -u"$DB_USER" -p"$DB_PASSWORD" -h "$DB_HOST" -P "$DB_PORT" "$DB_NAME" < internal/models/schema.sql
fi

# Check result
if [[ $? -eq 0 ]]; then
    echo "Database '$DB_NAME' recreated successfully."
else
    echo "Failed to recreate the database."
    exit 1
fi

go run cmd/server/main.go