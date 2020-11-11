#!/bin/bash

source .env

./wait-for-it.sh $POSTGRES_HOST:$POSTGRES_PORT
./toastnotes
