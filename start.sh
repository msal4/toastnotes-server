#!/bin/bash
# Use this script to start the notes server, should be used after starting the
# database service or at the same time.

source .env

./wait-for-it.sh $POSTGRES_HOST:$POSTGRES_PORT
./toastnotes
