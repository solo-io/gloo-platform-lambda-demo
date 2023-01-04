### gloo-platform-lambda-demo

This repo contains AWS Lambda code and Gloo resources to demonstrate AWS Lambda capability in Gloo Platform.

## Solution Overview

This demo produces the following configuration:

TODO diagram

## Requirements

* aws cli
* go >= 1.18
* jq

## Deployment

First, ensure your local AWS cli environment is healthy, then deploy the lambda to AWS:

> the script does some brute-force crud operations; some error output is expected

```bash
./deploy/aws/deploy.sh
```

Deploy Gloo Platform resources (edit overlays for your environment first!):

```bash
kubectl apply -k ./deploy/kustomize/overlays/demo
```

Test the /lambda route:

```bash
curl 'http://****.amazonaws.com/lambda/echo?input=solo' -i

HTTP/1.1 200 OK
date: Thu, 29 Dec 2022 14:56:09 GMT
content-type: application/json
content-length: 57
x-amzn-requestid: ****
x-amzn-remapped-content-length: 0
x-amz-executed-version: $LATEST
x-envoy-upstream-service-time: 20
server: istio-envoy
 
{"output": "solo"}
```
## Resource Description

TODO

## Local Development

Invoke the demo lambda locally by passing the `-event` param:

```bash
go build && ./gloo-platform-lambda-demo -event test/event-gloo-request-apigw-response.json
```

```
INFO[0000] AWS lambda runtime not detected, invoking with test file test/event-gloo-request-apigw-response.json 
{"statusCode":200,"headers":null,"multiValueHeaders":null,"body":"anannab"}
```