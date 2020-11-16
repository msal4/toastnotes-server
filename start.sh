#!/bin/bash
# get the db url from env vars
URL=$(echo $DATABASE_URL | grep -o "@.*\/")
# truncate the '@' at the start and the '/' at the end
URL=${URL:1:-1}
# wait for the db to start
./wait-for-it.sh $URL
# start the app
./app