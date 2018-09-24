# OpenTracing Demo using k8s and istio
Base code forked from [Tracing HTTP request latency in Go with OpenTracing](https://medium.com/@YuriShkuro/tracing-http-request-latency-in-go-with-opentracing-7cc1282a100a#.r6h4w2n9o)

## Prerequisites
- minikube
- istio

## Added
- Added Dockerfile
- Added k8s deployment using istio
- B3HTTPHeader 

## Usage
If you use minikube with vm driver, first follow [Local Kubernetes setup on macOS with minikube on VirtualBox and local Docker registry](https://gist.github.com/kevin-smets/b91a34cea662d0c523968472a81788f7)

Install istio with ```tracing: enable```
*Note :* if you use minikune don't forget to change istio ingresscontroller type to NodePort 

Build code
```bash
docker build -t opentracing-go-nethttp-demo:1.0 . -f Dockerfile
```

Push to registry
```bash
docker tag opentracing-go-nethttp-demo:1.0  localhost:5000/opentracing-go-nethttp-demo:1.2
docker push localhost:5000/opentracing-go-nethttp-demo
```

Deploy on minikube cluster
```bash
kubectl apply -f k8s/deployment.yaml
```

Try it
```bash
curl -HHost:tracing.192.168.99.100.xip.io http://tracing.192.168.99.100.xip.io:31380/gettime
```

See traces
```bash
kubectl port-forward service/jaeger-query 16686 -n istio-system
# Check on browser : http://127.0.0.1:16686/trace
```
