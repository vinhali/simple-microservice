# Jaeger Deployment
This readme provides a simple guide on deploying Jaeger, an open-source, end-to-end distributed tracing system.

## Deployment YAML
Use the following Kubernetes YAML configuration to deploy Jaeger:

```bash
kubectl create -f jaeger/deploy.yaml -n <NAMESPACE>
```

# Prometheus Deployment

Get Helm Repository

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
```

Install

```bash
helm install <RELEASE_NAME> prometheus-community/kube-prometheus-stack -f kube-stack/values.yaml -n <NAMESPACE>
```

Addtional

```bash
kubectl create -f kube-stack/cadvisor.yaml -n <NAMESPACE>
```
