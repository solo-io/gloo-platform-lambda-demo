#!/bin/bash

NAME="gloo-platform-lambda-demo"

aws_identity=$(aws sts get-caller-identity)
if [ "$?" -ne "0" ]; then
  echo "failed to execute \"aws sts get-caller-identity\"; please make sure your aws cli environment is set up"
  exit -1
fi
accountId=$(echo $aws_identity | jq -r ".Account")

rm -rf .build
mkdir -p .build

set -euo pipefail

GOOS=linux GOARCH=amd64 go build -o .build/main

zip .build/lambda.zip .build/main

set +euo pipefail

aws iam create-role --role-name "${NAME}" --tags "Key=created-by,Value=benji_lilley" "Key=team,Value=product" "Key=purpose,Value=product-development" --assume-role-policy-document '{"Version": "2012-10-17","Statement": [{ "Effect": "Allow", "Principal": {"Service": "lambda.amazonaws.com"}, "Action": "sts:AssumeRole"}]}' || true
aws iam attach-role-policy --role-name "${NAME}" --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole || true

aws lambda create-function \
    --function-name "${NAME}" \
    --runtime go1.x \
    --zip-file "fileb://.build/lambda.zip" \
    --handler ".build/main" \
    --role "arn:aws:iam::${accountId}:role/${NAME}" \
    --tags "created-by=benji_lilley,team=product,purpose=product-development" \
    || true

aws lambda update-function-code \
    --function-name "${NAME}" \
    --zip-file "fileb://.build/lambda.zip"