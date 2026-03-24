```mermaid
graph TD
    subgraph Server_DEV ["Standalone Server: DEV"]
        GitLab[GitLab Service]
        Registry((GitLab Registry))
    end

    subgraph K8s_Cluster ["Kubernetes Cluster"]
        subgraph Node_RUNNER ["Node: RUNNER"]
            subgraph R_NS ["Namespace: gitlab-runner"]
                Runner[GitLab Runner Pod]
            end
            Ingress[NGINX Ingress Controller]
        end

        subgraph Node_OSNOVA ["Node: OSNOVA"]
            subgraph AppNS ["Namespace: chat-app"]
                Pods[Go-Chat Pods]
                ESO[External Secrets Operator]
                K8sSecret[K8s Secret]
            end
            DB[(Postgres 15)]
        end
    end

    subgraph Security ["Security"]
        Vault[(HashiCorp Vault)]
    end

    %% Трафик пользователя
    User((User)) -->|chat.local| Ingress
    Ingress -->|Forward| Pods

    %% CI/CD Потоки (Вне -> Внутри)
    GitLab <-->|Jobs / API| Runner
    Runner -->|1. Build & Test| Runner
    Runner -->|2. Kaniko Push| Registry
    Runner -->|3. Helm Deploy| AppNS

    %% Работа с секретами и данными
    Vault -->|Provide Secrets| ESO
    ESO -.->|Sync| K8sSecret
    K8sSecret -.->|Inject ENV| Pods
    Pods --> DB

    %% Стилизация
    style Server_DEV fill:#f1f1f1,stroke:#fc6d26,stroke-width:2px
    style K8s_Cluster fill:#f8f9ff,stroke:#326ce5,stroke-width:2px
    style GitLab fill:#fc6d26,color:#fff
    style Runner fill:#fca326,color:#fff
    style Ingress fill:#00a6ed,color:#fff
    style Vault fill:#f5d142,stroke:#333
    style Pods fill:#61dafb,stroke:#333
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

