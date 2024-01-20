#!/bin/sh

set -e

echo "setting up app dependencies..."
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv ./migrate /usr/bin/migrate
make migrateup

echo "starting app..."
./app
