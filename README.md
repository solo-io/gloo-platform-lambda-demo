

### Requirements

* aws cli
* go >= 1.18
* jq

### Deployment

Deploy the lambda to AWS:

```bash
./deploy/aws/deploy.sh
```

Deploy Gloo Platform resources:

```bash
kubectl apply -k ./deploy/kustomize/overlays/demo
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