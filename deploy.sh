#!/bin/sh

lambdaArn="$1"

go build
zip handler.zip ./slurpCwsImages-lambda

echo "Deploying slurpCwsImages to $1"
aws lambda update-function-code \
--function-name $lambdaArn \
--zip-file fileb://handler.zip
