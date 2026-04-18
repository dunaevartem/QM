# Go-K8s-Chat
Современный масштабируемый чат на WebSockets с использованием Golang, Kubernetes и безопасным управлением секретами через HashiCorp Vault.
# Архитектура
* Backend: Go (Gin Framework + Gorilla WebSocket).
* Frontend: Статический HTML/JS (встроен в образ).
* Database: PostgreSQL 15 (StatefulSet с PersistentVolume).
* Secrets: HashiCorp Vault + External Secrets Operator.
* CI/CD: GitLab CI (Kaniko для сборки, Helm для деплоя).
# Технологический стек
* Runtime: golang:1.24-alpine
* Orchestration: Kubernetes 1.29+
* Package Manager: Helm 3
* Ingress: NGINX Ingress Controller (доступ по chat.local)
# Управление секретами (Vault)
В проекте реализован принцип Zero Secrets in Git. Все чувствительные данные (пароли БД, ключи реестра) хранятся в Vault и доставляются в кластер через ExternalSecrets.

__Настройка Vault__

Путь к секрету базы: secret/chat-app (поле password).

Путь к секрету реестра: secret/registry (поля user, password).

Роль в Vault: chat-app-role (привязана к ServiceAccount default).

# CI/CD Пайплайн
Автоматизация описана в .gitlab-ci.yml и включает 5 стадий:

* Build: Компиляция бинарного файла chat-server.
* Test: Запуск unit-тестов из директории /test.
* Push: Сборка Docker-образа через Kaniko (без привилегий root) и пуш в GitLab Registry.
* Scan: 
* Deploy: Деплой приложения в Kubernetes через Helm.

# Развертывание (Helm)
__Предварительные требования__

* Установленный ingress-nginx.
* Настроенный External Secrets Operator.
* Прописанный IP кластера в /etc/hosts:
```
<node-ip> chat.local
```
* Установка чарта
```
helm upgrade --install chat ./charts/chat-app \
  --namespace chat-app \
  --create-namespace \
  --set image.tag=latest \
  --set registry.url=$CI_REGISTRY
```
# Конфигурация (Values.yaml)

* replicaCount	Количество подов приложения	1
* service.port	Внутренний порт сервиса	80
* ingress.host	Домен для доступа к чату	chat.local
* postgres.secretName	Имя создаваемого секрета для БД	postgres-secret

# Разработка

Для локального запуска (требуется установленный Postgres):
```
export DATABASE_URL="postgres://user:pass@localhost:5432/chat_db"
go run ./cmd/main.go
```
Приложение будет доступно на порту :8080.

