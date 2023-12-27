# Web Banking Application 🚀

This repository contains a simple microservice built in Go Lang simulating a bank transfer with metrics and traces.

<details>
  <summary>🎉 Read the key points</summary>
<br>

🏦 Go Microservice: Banking Simplicity

- Welcome to the no-nonsense zone! This Go program means business, simulating a banking application using the Back-End-to-Frontend (BFF) pattern. It's like the efficient, no-frills cousin of banking microservices.

🛠️ What's Inside?

- Cut to the chase: this microservice handles fund transfers between predefined bank accounts. Simple, straightforward, and no gimmicks.

🔧 Tech Essentials:

- Powered by Go Lang and wrapped in a Docker image, it's the real deal. The Back-End-to-Frontend (BFF) design pattern keeps things neat and tidy. No magic tricks, just solid coding.

🌐 Where's the Action?

- Catch the action on port 8081 via HTTP server. Need access? Head straight to the /auth route. No velvet ropes, just straightforward functionality.

⚙️ Ready to Roll?

- Think of this Docker image as your trusty tool. It gets the job done without unnecessary frills. Efficient, effective, and ready to roll.

- No flashy lights, no confetti, just a reliable banking microservice doing its thing. Let's get down to business!

📸 Screenshot:

Check out this sneak peek into the action:

<img src="https://github.com/vinhali/simple-microservice/blob/main/files/screen.png?raw=true" alt="Banking Microservice" style="border-radius: 10px;">

</details>

<details>
  <summary>🏷️ Release Notes</summary>
  <br>

- v1.1.0 (stable)

</details>

## 👷🏻‍♂️ Architecture

</details>

<details>
  <summary>View Diagram</summary>
  <br>

  <img src="https://github.com/vinhali/simple-microservice/blob/main/files/web-banking.png?raw=true"/>

</details>

## 🏝️ Frontend

  <details>
    <summary>Read more</summary>
<br>
The frontend of the web banking application is built using Go Lang. It provides a user-friendly interface for customers to interact with their accounts, make transactions, and manage their finances.

### 🐳 Docker Image

You can find the Docker image for the frontend on [Docker Hub](https://hub.docker.com/r/luisvinhali/web-banking-frontend). To pull the image, use the following command:

```bash
docker pull luisvinhali/web-banking-frontend:latest
```

This will expose the frontend on `http://<ENDPOINT>:8080/auth`.

  </details>

## 🚢 Backend

<details>
  <summary>Read more</summary>
<br>
The backend of the web banking application is built using Go Lang. It handles business logic and serves as the API for the frontend.

### 🐳 Docker Image

You can find the Docker image for the backend on [Docker Hub](https://hub.docker.com/r/luisvinhali/web-banking-backend). To pull the image, use the following command:

```bash
docker pull luisvinhali/web-banking-backend:latest
```

This will expose the backend on `http://<ENDPOINT>:8081/transfer` and `http://<ENDPOINT>:8081/output`.

</details>

## 🔗 Connecting Frontend and Backend

<details>
  <summary>Read More</summary>
<br>

By default, the frontend is configured to communicate with the backend at `http://<ENDPOINT>:8081/transfer.` Ensure that the backend container is running and accessible for seamless integration.

Feel free to explore and customize the code according to your requirements. For additional details, refer to the documentation in the respective GitHub repositories:

- **Frontend Repository:** [luisvinhali/web-banking-frontend](https://hub.docker.com/r/luisvinhali/web-banking-frontend)
- **Backend Repository:** [luisvinhali/web-banking-backend](https://hub.docker.com/r/luisvinhali/web-banking-backend)

</details>

## 🏃 How to Run

<details>
  <summary>Read the step by step</summary>

### Using Docker Run:

Frontend:

```bash
docker run -p 8080:8080 web-banking-frontend:1.0
```

Backend:

```bash
docker run -p 8081:8081 web-banking-backend:1.0
```

###  Using Kubernetes:

Apply the following YAML manifest:

```bash
kubectl apply -f deploy.yaml
```

  <details>
    <summary>Use this file</summary>

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-banking-frontend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: web-banking-frontend
  template:
    metadata:
      labels:
        app: web-banking-frontend
    spec:
      containers:
      - name: web-banking-frontend
        image: web-banking-frontend:1.0
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: web-banking-frontend-service
spec:
  selector:
    app: web-banking-frontend
  ports:
    - protocol: TCP
      port: 8080
      targetPort: 8080
  type: NodePort
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-banking-backend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: web-banking-backend
  template:
    metadata:
      labels:
        app: web-banking-backend
    spec:
      containers:
      - name: web-banking-backend
        image: web-banking-backend:1.0
        ports:
        - containerPort: 8081
---
apiVersion: v1
kind: Service
metadata:
  name: web-banking-backend-service
spec:
  selector:
    app: web-banking-backend
  ports:
    - protocol: TCP
      port: 8081
      targetPort: 8081
  type: NodePort
```

  </details>

</details>

## ⚙️ Technical description

<details>
  <summary>Read More</summary>
<br>

Here is a concise technical overview of the sections of the code:

Backend

| Code | Description |
| ---- | ----------- |
| `accountBalance`, `transactionStatus`, `transactionAmount` | Defines Prometheus metrics to monitor account balance, transaction status, and transferred amount. |
| `tracer = otel.Tracer("transfer")` | Initializes the OpenTelemetry tracer for transaction tracing. |
| `Jaeger DNS` | http://jaeger-service.jaeger.svc.cluster.local:14268/api/traces |
| `otelsetup.SetupOTelSDK` | Configures the OpenTelemetry SDK for the banking service, including Prometheus and Jaeger. |
| `/metrics` | Exposes Prometheus metrics through the `/metrics` endpoint. |
| `transactionHandler()` | Handles transfer requests, generating metrics and tracing. |

Frontend

| Code | Description |
| ---- | ----------- |
| `authTotalCounter`, `authSuccessCounter`, `authFailureCounter` | Defines Prometheus metrics to monitor authentication attempts and outcomes. |
| `tracer = otel.Tracer("auth")` | Initializes the OpenTelemetry tracer for authentication tracing. |
| `/metrics` | Exposes Prometheus metrics through the `/metrics` endpoint. |
| `frontHandler()` | Handles authentication requests, generating metrics and tracing. |

Common

| Libs | Description |
| ---- | ----------- |
| `go.opentelemetry.io/otel and derivatives` | Imports the OpenTelemetry library for code instrumentation. |
| `github.com/vinhali/simple-microservice/otelsetup` |Initialize the JDK for OpenTelemetry. |
| `github.com/prometheus/client_golang/prometheus` | Imports the Prometheus library for metrics and monitoring. |
| `github.com/prometheus/client_golang/prometheus/promhttp` | Imports the Prometheus library for exposing metrics via HTTP. |

Prometheus

| Metric Name | Description |
| ----------- | ----------- |
| `web_banking_account_balance` | Current balance of an account. |
| `web_banking_transaction_status` | Transaction status (success or failure). |
| `web_banking_transaction_amount` | Amount transferred in a transaction. |
| `web_banking_auth_total` | Total number of auth attempts. |
| `web_banking_auth_success_total` | Total number of successful auths. |
| `web_banking_auth_failure_total` | Total number of failed auths. |

Jaeger (service name is `digital-bank` in version `1.0.0`)

| Key | Description |
| ----------- | ----------- |
| `/auth` | Context |
| `/transfer` | Context |
| `auth.success` | Parent Context |
| `auth.failure` | Parent Context |
| `auth.info` | Attribute |
| `transfer.forward` | Attribute |

Workflow

| Action | Description |
| ------------ | ----------- |
| `Request /auth` | ✉️ Successfully sent a call for authentication. |
| `Context auth.success sets the attribute auth.info` | ✅ Authentication completed successfully. |
| `Context auth.failure sets the attribute auth.info` | ❌ Authentication failed. |
| `Request /transfer` | 📤 No specific message, becomes a context as a child of /auth.. |
| `Context /transfer sets the attribute transfer.forward` | ❓ Server not found or ✅ request sent successfully. |

</details>
