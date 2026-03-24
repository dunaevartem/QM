```mermaid
graph TD
    subgraph "Public Internet"
        User((User)) -->|chat.local| Ingress[NGINX Ingress Controller]
    end

    subgraph "Kubernetes Cluster (3 Nodes)"
        subgraph Node1 ["Node 1 (System)"]
            Runner[GitLab Runner Pod]
            Ingress
        end

        subgraph Node2 ["Node 2 (Apps)"]
            subgraph AppNS ["Namespace: chat-app"]
                Pods[Go-Chat Pods]
                ESO[External Secrets Operator]
                K8sSecret[K8s Secret]
            end
        end

        subgraph Node3 ["Node 3 (Data)"]
            DB[(Postgres 15)]
            PV[Persistent Volume]
        end
    end

    subgraph "External Security & CI/CD"
        Vault[(HashiCorp Vault)]
        GitLab[GitLab Project]
        Registry((GitLab Registry))
    end

    %% Трафик пользователя
    Ingress -->|Route to Port 80| Pods

    %% CI/CD Потоки
    GitLab <-->|Jobs| Runner
    Runner -->|1. Build & Test| Runner
    Runner -->|2. Kaniko Push| Registry
    Runner -->|3. Helm Deploy| AppNS

    %% Работа с секретами
    Vault -->|secret/chat-app| ESO
    ESO -.->|Sync| K8sSecret
    K8sSecret -.->|Inject ENV| Pods

    %% Данные
    Pods --> DB
    DB --- PV

    %% Стилизация
    style Ingress fill:#00a6ed,color:#fff,stroke:#333
    style Runner fill:#fca326,color:#fff,stroke:#333
    style Vault fill:#f5d142,stroke:#333
    style Pods fill:#61dafb,stroke:#333
    style Node1 fill:#fafafa,stroke:#999,stroke-dasharray: 5 5
    style Node2 fill:#fafafa,stroke:#999,stroke-dasharray: 5 5
    style Node3 fill:#fafafa,stroke:#999,stroke-dasharray: 5 5
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

