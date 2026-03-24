```mermaid
graph TD
    subgraph Cluster ["Kubernetes Cluster"]
        
        %% Общий неймспейс, проходящий через все ноды
        subgraph R_NS ["Namespace: gitlab-runner"]
            direction TB
            subgraph R_DEV [" "]
                R1_1[Runner Pod 1]
                R1_2[Runner Pod 2]
            end
            subgraph R_RUN [" "]
                R2_1[Runner Pod 3]
                R2_2[Runner Pod 4]
            end
            subgraph R_OSN [" "]
                R3_1[Runner Pod 5]
                R3_2[Runner Pod 6]
            end
        end

        subgraph Node_DEV ["Node: dev"]
            subgraph OS_DEV ["Host OS Level"]
                GitLab[GitLab Service]
                Registry((GitLab Registry))
            end
            R_DEV
        end

        subgraph Node_RUNNER ["Node: runner"]
            R_RUN
            Ingress[NGINX Ingress Controller]
        end

        subgraph Node_OSNOVA ["Node: osnova"]
            subgraph AppNS ["Namespace: chat-app"]
                Pods[Go-Chat Pods]
                ESO[External Secrets Operator]
                K8sSecret[K8s Secret]
            end
            R_OSN
            DB[(Postgres 15)]
        end
        
    end

    subgraph External ["Security"]
        Vault[(HashiCorp Vault)]
    end

    %% Связи
    User((User)) -->|chat.local| Ingress
    Ingress -->|Route| Pods
    
    GitLab <-->|Jobs / API| R_NS
    R_NS -->|Push Image| Registry
    R_NS -->|Deploy| AppNS

    Vault -->|Provide Secrets| ESO
    ESO -.->|Sync| K8sSecret
    K8sSecret -.->|Inject ENV| Pods
    Pods --> DB

    %% Стилизация
    style OS_DEV fill:#fff,stroke:#fc6d26,stroke-width:2px
    style GitLab fill:#fc6d26,color:#fff
    style R_NS fill:#fff,stroke:#fca326,stroke-width:3px,stroke-dasharray: 10 5
    style R1_1 fill:#fca326,color:#fff; style R1_2 fill:#fca326,color:#fff
    style R2_1 fill:#fca326,color:#fff; style R2_2 fill:#fca326,color:#fff
    style R3_1 fill:#fca326,color:#fff; style R3_2 fill:#fca326,color:#fff
    style Ingress fill:#00a6ed,color:#fff
    style Vault fill:#f5d142,stroke:#333
    style Pods fill:#61dafb,stroke:#333
    style Node_DEV fill:#f9f9f9,stroke:#333,stroke-dasharray: 5 5
    style Node_RUNNER fill:#f9f9f9,stroke:#333,stroke-dasharray: 5 5
    style Node_OSNOVA fill:#f9f9f9,stroke:#333,stroke-dasharray: 5 5
    style R_DEV fill:none,stroke:none
    style R_RUN fill:none,stroke:none
    style R_OSN fill:none,stroke:none
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

