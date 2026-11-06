# Подготовка.
FROM golang:alpine AS preper
ARG NAME
WORKDIR /usr/src/${NAME}
# Создание директорий.
RUN mkdir /etc/${NAME}
RUN mkdir /var/log/${NAME}
# Копирование файлов.
COPY . .
COPY ./config/prod.yaml /etc/${NAME}/prod.yaml
RUN go mod download

# Сборка.
FROM preper AS building
RUN go build -v -o /usr/local/bin/${NAME} ./cmd/main.go

# Релиз.
FROM alpine AS release
ARG NAME
# Настройки.
RUN apk add tzdata
RUN ln -s /usr/share/zoneinfo/Europe/Moscow /etc/localtime
# Создание директорий.
RUN mkdir /var/log/${NAME}
RUN mkdir /etc/${NAME}
# Копирование файлов.
COPY ./config/prod.yaml /etc/${NAME}/prod.yaml
COPY --from=building /usr/local/bin/${NAME} /usr/local/bin/${NAME}
COPY ./docker-entrypoint.sh /usr/local/bin
# Создание точки входа.
RUN echo "exec ${NAME} 1>> /var/log/${NAME}/info.log" >> /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT [ "docker-entrypoint.sh" ]
