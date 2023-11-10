#!/bin/bash

# Define here a bash script that will be executing the migration on the database(es).
# Don't hardcode any credentials here. Use the credentials which are supplied through 
# TaskDefinition (see environment and secrets list from ./deployment/task-definition.json.tpl)

# This script should always exit with 0 if migration succeeded, and with non-zero if it fails
# defaultSchema is created by flyway implicitly before first migration script execution
flyway migrate \
  -defaultSchema=${FLYWAY_DB_DEFAULT_SCHEMA} \
  -user=${FLYWAY_DB_USERNAME} \
  -password=${FLYWAY_DB_PASSWORD} \
  -url=jdbc:postgresql://${FLYWAY_DB_HOST}:${FLYWAY_DB_PORT}/${FLYWAY_DB_DATABASE}?${FLYWAY_DB_PARAMS} \
  -connectRetries=5 \
  -locations='filesystem:/migrations'
