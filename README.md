# Enterprise Cloud-Native DevOps & LLMOps Blueprint

Welcome to the **Enterprise Cloud-Native DevOps & LLMOps Blueprint**. This project is a comprehensive reference architecture demonstrating a zero-trust, fully automated, and highly observable GitOps pipeline for modern Kubernetes deployments. It also features a **"Secure, Air-Gapped Enterprise LLMOps Platform"** that allows running Large Language Models (LLMs) completely in-house without exposing sensitive corporate data.

## 🏛 Architecture Overview

At the heart of this blueprint lies a highly secure, automated GitOps flow running on **K3s**. It leverages industry-standard tools to handle everything from CI/CD to Zero Trust networking.

```mermaid
graph TD
    %% Developers & Code
    Developer((Developer)) -->|Git Push| GitHub[GitHub Repository]
    GitHub -->|Trigger| GHA[GitHub Actions CI]
    GHA -->|Build & Push| DockerHub[(Docker Hub)]

    %% GitOps
    GitHub -->|Monitor| ArgoCD[ArgoCD GitOps]
    ArgoCD -->|Sync State| K3s((K3s Kubernetes Cluster))

    %% Infrastructure & Security inside K3s
    subgraph K3s Cluster
        %% Security
        Vault[HashiCorp Vault / OpenBao<br/>Internal PKI & Secrets]
        CertManager[Cert-Manager]
        ESO[External Secrets Operator]

        %% Networking
        IstioIngress[Istio IngressGateway<br/>mTLS Edge Router]
        IstioMesh[Istio Service Mesh]

        %% Observability
        Prometheus[Prometheus & Grafana]

        %% Applications & AI
        subgraph Microservices
            Frontend[Go Frontend]
            Backend[Go Backend API]
            Redis[(Redis Cache)]
        end
        
        subgraph AI Stack (LLMOps)
            WebUI[Open WebUI<br/>Internal ChatGPT]
            Ollama[Ollama Engine<br/>Llama 3.2 Model]
        end

        %% Connections within cluster
        Vault -->|Issues Certs| CertManager
        Vault -->|Syncs Secrets| ESO
        CertManager -->|Injects TLS| IstioIngress
        ESO -->|Injects Passwords| Redis
        
        IstioIngress -->|Routes Traffic| Frontend
        IstioIngress -->|Routes Traffic| WebUI
        Frontend -->|GRPC / REST| Backend
        Backend -->|Reads/Writes| Redis
        WebUI -->|API Requests| Ollama
        
        IstioMesh -.->|Enforces mTLS| Microservices
        IstioMesh -.->|Enforces mTLS| WebUI
        Prometheus -.->|Scrapes Metrics| Microservices
        Prometheus -.->|Scrapes Metrics| Ollama
    end

    %% External Traffic
    User((End User)) -->|HTTPS / TLS| IstioIngress
```

## 🚀 Key Features & Technologies

### 1. True GitOps with ArgoCD
- **No Manual Kubectl:** The entire cluster state (from microservices to the Vault PKI engine) is defined in this repository and continuously reconciled by **ArgoCD**.
- **Self-Healing:** Any manual changes made directly to the cluster are automatically overwritten by ArgoCD to maintain the desired state defined in Git.

### 2. Zero Trust Security & Automated PKI
- **HashiCorp Vault (OpenBao):** Acts as the central secret store and the internal Root Certificate Authority (PKI).
- **Cert-Manager:** Automatically requests, renews, and injects X.509 TLS certificates from Vault into the Istio IngressGateway.
- **External Secrets Operator (ESO):** Fetches database credentials from Vault and mounts them securely into application pods without exposing them in Git.

### 3. Service Mesh & Edge Routing
- **Istio:** Provides advanced traffic management, telemetry, and enforces mutual TLS (mTLS) between all microservices. The **Istio IngressGateway** serves as the single entry point, securing traffic with the Vault-issued certificates.

### 4. Comprehensive Observability
- **Prometheus & Grafana:** Automatically scrapes metrics from K3s nodes, Istio proxies, and custom application endpoints, visualizing them in real-time dashboards.

### 5. Infrastructure as Code (IaC) & DevSecOps CI/CD
- **Ansible:** Used to bootstrap the underlying Ubuntu VMs and install the base K3s cluster.
- **GitHub Actions (CI/CD):** Automates the testing, building, and pushing of multi-arch Docker images on every push to the `main` branch. 
- **Shift-Left Security (Trivy):** Before any image is pushed, the pipeline runs **Aqua Trivy** to scan the Docker images for OS and library vulnerabilities.
- **Automated GitOps Commit:** Once the CI/CD pipeline passes, GitHub Actions automatically updates the Kustomize manifests with the new image tag (`${{ github.sha }}`) and pushes the changes back to the repository. ArgoCD detects this commit and instantly deploys the new version.

### 6. Secure, Air-Gapped Enterprise LLMOps Platform
- **In-House LLMOps (Data Privacy):** Runs state-of-the-art open-source LLMs (like Llama 3 via Ollama) completely within the internal Kubernetes network (Air-Gapped). Sensitive corporate data never leaves the internal infrastructure, eliminating the privacy risks of public APIs like OpenAI.
- **Open WebUI:** Provides a sleek, ChatGPT-like interface connected directly to the internal LLM, making AI accessible to internal employees.
- **Zero-Trust AI:** The AI interfaces and APIs are secured behind the Istio Ingress Gateway, utilizing mTLS and Vault-issued certificates to ensure enterprise-grade security.
## 📦 How to Run

*Since this is a reference architecture, it assumes you have a running K3s cluster. The entire state is managed via the `gitops/` directory.*

1. **Bootstrap ArgoCD:**
   ```bash
   kubectl create namespace argocd
   kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
   kubectl apply -f gitops/bootstrap-argocd.yaml
   ```
2. **Watch the Magic:** ArgoCD will automatically read the repository and begin spinning up Vault, Istio, Prometheus, and the Microservices.

## 🎯 Motivation
This repository is the culmination of deep architectural research into Cloud-Native ecosystems. It serves as a personal sandbox and a public reference for engineers looking to implement enterprise-grade Kubernetes infrastructures.
