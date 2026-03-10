# Go Calculator (K8s Project)

Веб-приложение для вычислений, написанное на Go, развернутое в Kubernetes-кластере с использованием современного стека автоматизации.

# Стек технологий

* Backend: Go 1.24 (Framework Gin Gonic)
* Runtime: containerd + nerdctl (без использования Docker Docker Engine)
* Orchestration: Kubernetes v1.35.1
* Deployment: Helm 3
* CI/CD: GitLab CI + Kaniko (rootless image build)

# Архитектура и трафик

Проект демонстрирует полный путь запроса от клиента до контейнера в изолированной сети.

```mermaid
graph TD
    User((Пользователь)) -->|1. http://kurtonic.local| Hosts[Файл /etc/hosts]
    Hosts -->|2. IP: 192.168.1.50| IngressSVC[Ingress Controller]
    IngressSVC -->|3. Host: kurtonic.local| IngressRule[Ingress Resource]
    IngressRule -->|4. Backend Routing| GoService[Service: go-app-service]
    GoService -->|5. Load Balancing| Pod1[Pod: replica-1]
    GoService -->|5. Load Balancing| Pod2[Pod: replica-2]

    subgraph "K8s Cluster"
        IngressRule
        GoService
        Pod1
        Pod2
    end
```

# CI/CD Pipeline

__Автоматизация реализована через .gitlab-ci.yaml и включает 4 стадии:__
* Build: Сборка Go-бинарника (артефакт app).
* Test: Запуск модульных тестов (go test ./...).
* Push: Сборка Docker-образа через Kaniko. Это позволяет собирать образы внутри K8s без привилегированного доступа и Docker-in-Docker.
* Deploy: Деплой в кластер через Helm. Тег образа автоматически обновляется на основе короткого SHA коммита.

# Установка и запуск

Локально (Development)
```
go run main.go
```

Приложение будет доступно на http://localhost:8080.

В кластере (Production-ready)

# Настройка containerd:

Убедитесь, что в /etc/containerd/config.toml прописан config_path для работы с вашим реестром по HTTP/HTTPS.

# Деплой через Helm:

```
helm upgrade --install go-app ./charts/go-app \
  --set image.repository=$CI_REGISTRY_IMAGE \
  --set image.tag=latest
```

# Доступ:

Добавьте запись в /etc/hosts:

192.168.1.50 kurtonic.local


# Структура репозитория

* __main.go__ — Серверная логика и API (Gin).
* __templates/__ — UI шаблоны (index, result, error).
* __charts/go-app/__ — Helm-чарт (Deployment на 2 реплики, ClusterIP Service, Ingress).
* __Dockerfile__ — Легковесный образ на базе Alpine.
* __test/__ — Модульные тесты приложения.
* __myapp.service__ — (Legacy) Конфиг для запуска через systemd.

# Мониторинг и отладка

__Посмотреть логи всех реплик сразу:__
```
kubectl logs -l app.kubernetes.io/name=go-app -f
```

__Проверить статус здоровья:__
```
curl http://kurtonic.local

```

