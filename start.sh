#!/bin/sh
if [[ $1 == "--build" || $1 == "-b" ]]
then
  . .env
  go build -o app .
fi

# get the db url from env vars
URL=$(echo $DATABASE_URL | grep -o "@.*\/")
# truncate the '@' at the start and the '/' at the end
URL=${URL#?}
URL=${URL%?}
# wait for the db to start
./wait-for-it.sh $URL
# start the app
./app
