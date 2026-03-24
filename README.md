mermaid
```
graph TD
    subgraph "External Access"
        User((User)) --> Ingress[NGINX Ingress: chat.local]
    end

    subgraph "Kubernetes Cluster (Namespace: chat-app)"
        Ingress --> Service[K8s Service: Port 80]
        
        subgraph "Application Layer"
            Service --> Pods[Go-Chat Pods: Gin + WebSocket]
            ESO[External Secrets Operator] -.->|Creates| K8sSecret[K8s Secret]
            K8sSecret -.->|Injects ENV| Pods
        end

        subgraph "Data Layer"
            Pods --> DB[(PostgreSQL 15)]
            DB --- PV[Persistent Volume]
        end
    end

    subgraph "Security & Secrets"
        Vault[(HashiCorp Vault)]
        Vault -->|secret/chat-app| ESO
        Vault -->|secret/registry| ESO
    end

    subgraph "CI/CD Pipeline (GitLab)"
        Code[Git Repo] --> Build[Build: Go Binary]
        Build --> Test[Test: Unit Tests]
        Test --> Push[Push: Kaniko Image]
        Push --> Deploy[Deploy: Helm Upgrade]
        
        Push -.-> Registry((GitLab Registry))
        Deploy -.-> Pods
    end

    style Vault fill:#f5d142,stroke:#333,stroke-width:2px
    style Pods fill:#61dafb,stroke:#333,stroke-width:2px
    style DB fill:#336791,stroke:#f2f2f2,color:#fff
```
# Go-K8s-Chat
Современный масштабируемый чат на WebSockets с использованием Golang, Kubernetes и безопасным управлением секретами через HashiCorp Vault.
# Архитектура
Backend: Go (Gin Framework + Gorilla WebSocket).
Frontend: Статический HTML/JS (встроен в образ).
Database: PostgreSQL 15 (StatefulSet с PersistentVolume).
Secrets: HashiCorp Vault + External Secrets Operator.
CI/CD: GitLab CI (Kaniko для сборки, Helm для деплоя).
# Технологический стек
Runtime: golang:1.24-alpine
Orchestration: Kubernetes 1.29+
Package Manager: Helm 3
Ingress: NGINX Ingress Controller (доступ по chat.local)
# Управление секретами (Vault)
В проекте реализован принцип Zero Secrets in Git. Все чувствительные данные (пароли БД, ключи реестра) хранятся в Vault и доставляются в кластер через ExternalSecrets.
Настройка Vault
Путь к секрету базы: secret/chat-app (поле password).
Путь к секрету реестра: secret/registry (поля user, password).
Роль в Vault: chat-app-role (привязана к ServiceAccount default).
# CI/CD Пайплайн
Автоматизация описана в .gitlab-ci.yml и включает 5 стадий:
Build: Компиляция бинарного файла chat-server.
Test: Запуск unit-тестов из директории /test.
Push: Сборка Docker-образа через Kaniko (без привилегий root) и пуш в GitLab Registry.
Pre-deploy: Подготовка окружения.
Deploy: Деплой приложения в Kubernetes через Helm.
# Развертывание (Helm)
Предварительные требования
Установленный ingress-nginx.
Настроенный External Secrets Operator.
Прописанный IP кластера в /etc/hosts:
```
<node-ip> chat.local
```
Установка чарта
```
helm upgrade --install chat ./charts/chat-app \
  --namespace chat-app \
  --create-namespace \
  --set image.tag=latest \
  --set registry.url=$CI_REGISTRY
```
# Конфигурация (Values.yaml)
Параметр	Описание	Значение по умолчанию
replicaCount	Количество подов приложения	1
service.port	Внутренний порт сервиса	80
ingress.host	Домен для доступа к чату	chat.local
postgres.secretName	Имя создаваемого секрета для БД	postgres-secret
# Разработка
Для локального запуска (требуется установленный Postgres):
```
export DATABASE_URL="postgres://user:pass@localhost:5432/chat_db"
go run ./cmd/main.go
```
Приложение будет доступно на порту :8080.

