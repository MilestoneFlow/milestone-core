#!/bin/bash

export $(grep -v '^#' ../.env | xargs)
POSTGRES_CONTAINER_NAME="milestone_postgres"

echo "Drop existing test database $POSTGRES_TEST_DB_NAME"
docker exec -e PGPASSWORD=$POSTGRES_DB_PASSWORD $POSTGRES_CONTAINER_NAME psql -U $POSTGRES_DB_USER -c "DROP DATABASE IF EXISTS $POSTGRES_TEST_DB_NAME"

echo "Creating test database $POSTGRES_TEST_DB_NAME"
docker exec -e PGPASSWORD=$POSTGRES_DB_PASSWORD $POSTGRES_CONTAINER_NAME psql -U $POSTGRES_DB_USER -c "CREATE DATABASE $POSTGRES_TEST_DB_NAME"

if [ $? -eq 0 ]; then
  echo "Test database $POSTGRES_TEST_DB_NAME created successfully"
else
  echo "Failed to create test database $POSTGRES_TEST_DB_NAME"
fi

echo "Apply migrations to test database $POSTGRES_TEST_DB_NAME"
POSTGRES_URI="postgresql://$POSTGRES_DB_USER:$POSTGRES_DB_PASSWORD@$POSTGRES_SERVER_HOST:$POSTGRES_SERVER_PORT/$POSTGRES_TEST_DB_NAME"
goose -dir ../migrations postgres "$POSTGRES_URI" up