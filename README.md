# Serverless Platform

A serverless computing platform inspired by AWS Lambda. This project offers a functions-as-a-service (FaaS) implementation where code is executed in response to events. The platform integrates Docker for function isolation, Kubernetes for orchestration, and RabbitMQ for event-driven architecture and asynchronous task processing.

## Prerequisites

- Kubernetes cluster
- Docker
- Helm (for RabbitMQ deployment)
- Go 1.17 or newer

## Setup

1. **Kubernetes Cluster**: Set up a Kubernetes cluster
2. **RabbitMQ**: Deploy RabbitMQ on your cluster using its Helm chart:
    ```bash
    helm repo add bitnami https://charts.bitnami.com/bitnami
    helm install my-rabbitmq bitnami/rabbitmq
    ```

3. **Build & Push Docker Images**:
    ```bash
    cd faas-controller
    docker build -t yourdockerhubusername/faas-controller .
    docker push yourdockerhubusername/faas-controller

    cd ../faas-executor
    docker build -t yourdockerhubusername/faas-executor .
    docker push yourdockerhubusername/faas-executor
    ```

4. **Deploy to Kubernetes**:
    ```bash
    kubectl apply -f k8s-manifests/controller-deployment.yaml
    kubectl apply -f k8s-manifests/executor-deployment.yaml
    ```

## Usage

1. **Deploy a Function**:
    ```bash
    curl -X POST http://controller-service-endpoint/function -d "@path-to-your-function.zip"
    ```

2. **Invoke a Function**:
    ```bash
    curl -X POST http://controller-service-endpoint/function/invoke -d '{"function_id": "your_function_id"}'
    ```

Replace `controller-service-endpoint` with the actual endpoint for your FaaS controller service. The exact details and additional API endpoints will depend on your implementation.

## License

This project is licensed under the MIT License.
