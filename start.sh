#!/bin/sh

set -e

echo "setting up app dependencies..."
export DB_USER=$(cat /run/secrets/db_user)
export DB_PASSWORD=$(cat /run/secrets/db_password)
export TOKEN_SYMMETRIC_KEY=$(cat /run/secrets/token)
make migrateup

echo "start the app"
exec "$@"
