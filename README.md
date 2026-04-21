## Go-Chat

Современный масштабируемый чат на WebSockets, построенный на микросервисной архитектуре с использованием Golang и Kubernetes. Особое внимание уделено безопасности и автоматизации доставки (CI/CD).

## Архитектура
*   **Backend:** Go 1.24 (Gin Framework + Gorilla WebSocket).
*   **Frontend:** Статический HTML/JS (встроен в Docker-образ).
*   **Database:** PostgreSQL 15 (StatefulSet с PersistentVolume для хранения данных).
*   **Secrets:** HashiCorp Vault + External Secrets Operator (автоматическая доставка секретов в K8s).
*   **CI/CD:** GitLab CI (Kaniko для безопасной сборки, Helm для деплоя).

## Технологический стек
*   **Runtime:** `golang:1.24-alpine` (минимальный размер и поверхность атаки).
*   **Orchestration:** Kubernetes 1.29+.
*   **Package Manager:** Helm 3.
*   **Ingress:** NGINX Ingress Controller (доступ по адресу `chat.local`).

## Безопасность (DevSecOps)
Проект придерживается подходов **DevSecOps** и обеспечивает защиту на всех уровнях:

1.  **Zero Secrets in Git:** Чувствительные данные (пароли БД, ключи реестра) никогда не попадают в репозиторий.
2.  **HashiCorp Vault:** Используется как централизованное доверенное хранилище секретов.
3.  **Container Scanning (Trivy):** В пайплайн интегрирован сканер `aquasec/trivy`. Он проводит аудит образа на наличие CVE. Пайплайн принудительно останавливается (`exit-code 1`), если найдены уязвимости уровня **CRITICAL**.
4.  **Rootless Build (Kaniko):** Сборка образов происходит без использования Docker-in-Docker и без привилегий root, что исключает возможность атаки на хостовую систему GitLab Runner.
5.  **External Secrets Operator:** Автоматически синхронизирует секреты из Vault в нативные `Secret` объекты Kubernetes.

## Настройка Vault
*   **Путь к секрету базы:** `secret/chat-app` (поле `password`).
*   **Путь к секрету реестра:** `secret/registry` (поля `user`, `password`).
*   **Роль в Vault:** `chat-app-role` (привязана к ServiceAccount `default`).

## CI/CD Пайплайн
Автоматизация описана в `.gitlab-ci.yml` и включает 5 стадий:

*   **Build:** Компиляция бинарного файла приложения.
*   **Test:** Запуск unit-тестов из директории `/test`.
*   **Push:** Сборка Docker-образа через Kaniko и пуш в GitLab Registry.
*   **Scan:** Глубокий аудит безопасности образа с помощью **Trivy** (блокировка при CRITICAL).
*   **Deploy:** Автоматический деплой приложения в Kubernetes через Helm.

## Развертывание (Helm)
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
## Конфигурация (Values.yaml)

| Параметр | Описание | Значение |
| :--- | :--- | :--- |
| `replicaCount` | Количество подов приложения | `2` |
| `service.port` | Внутренний порт сервиса | `80` |
| `ingress.host` | Домен для доступа к чату | `chat.local` |
| `postgres.secretName` | Имя секрета для подключения к БД | `postgres-secret` |

## Разработка

Для локального запуска (требуется установленный Postgres):
```
export DATABASE_URL="postgres://user:pass@localhost:5432/chat_db"
go run ./cmd/main.go
```
Приложение будет доступно на порту :8080.

## Лицензия

Данный проект распространяется под лицензией **MIT**. Подробности в файле [LICENSE](LICENSE).

