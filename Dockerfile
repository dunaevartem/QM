# Используем легкий образ только для запуска
FROM alpine:3.19

# Библиотеки для работы Go и SSL (нужны для внешних API)
RUN apk add --no-cache ca-certificates libc6-compat

# Копируем УЖЕ СКОМПИЛИРОВАННЫЙ бинарник (из артефактов GitLab)
COPY chat-server /app

# Копируем статику для фронтенда чата
COPY static /static

# Выставляем порт чата
EXPOSE 8080

# Запускаем приложение
ENTRYPOINT ["/app"]
