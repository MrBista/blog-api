#!/bin/sh

echo "Waiting for MySQL to be ready..."

# Tunggu sampai port MySQL terbuka
until nc -z -v -w30 blog-api-db 3306
do
  echo "Waiting for database connection..."
  sleep 5
done

echo "Database is up! Starting app..."
./app
