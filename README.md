# go-musthave-metrics-tpl

Шаблон репозитория для трека «Сервер сбора метрик и алертинга».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-metrics-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).


## 

### iter 2 

metricstest -test.v -test.run="^TestIteration2[AB]*$" \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent

### iter 3

metricstest -test.v -test.run="^TestIteration3[AB]*$" \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server

### iter 4

SERVER_PORT=8085
          TEMP_FILE=./temp
          metricstest -test.v -test.run="^TestIteration4$" \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=$SERVER_PORT \
            -source-path=.


### iter 5

SERVER_PORT=8081
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=./temp
          metricstest -test.v -test.run="^TestIteration5$" \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=$SERVER_PORT \
            -source-path=.


### iter 6

SERVER_PORT=8081
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=./temp
          metricstest -test.v -test.run="^TestIteration6$" \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=$SERVER_PORT \
            -source-path=.

### iter 7

SERVER_PORT=8081
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=./temp
          metricstest -test.v -test.run="^TestIteration7$" \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=$SERVER_PORT \
            -source-path=.

### iter 8

SERVER_PORT=8080
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=./temp
          metricstest -test.v -test.run="^TestIteration8$" \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=$SERVER_PORT \
            -source-path=.