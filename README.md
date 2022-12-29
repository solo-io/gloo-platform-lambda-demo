

### Requirements

* aws cli
* go >= 1.18
* jq

### Deployment

Deploy the lambda to AWS:

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
 
solo
```

### Run Local

Invoke the demo lambda locally by passing the `-event` param:

```bash
go build && ./gloo-platform-lambda-demo -event test/event-gloo-request-apigw-response.json
```

```
INFO[0000] AWS lambda runtime not detected, invoking with test file test/event-gloo-request-apigw-response.json 
{"statusCode":200,"headers":null,"multiValueHeaders":null,"body":"anannab"}
```