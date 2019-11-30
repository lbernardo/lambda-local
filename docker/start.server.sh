#!/bin/bash
if [ -z $HOST ]
then
    HOST="0.0.0.0"
fi
if [ -z $PORT ]
then
    PORT="3000"
fi
lambda-local start --host $HOST --port $PORT --volume $VOLUME_APP