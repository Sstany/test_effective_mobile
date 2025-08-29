# Сервис аггрегации онлайн-подписок пользователей

## Запуск сервиса

`docker compose up -d --build `

Последующие запуски можно выполнять с помощью `docker compose up -d`

Полностью остановить сервис вместе с удалением данныз можно с помощью `docker compose down -v`

Сервис доступен по адресу: `http://localhost:8080`

## Конфигурация проекта

Дефолтные значения переменных заданы в .env файле

Внутри `docker-compose.yaml` задается строка подключения к базе в виде перменной окружения `DATABASE_CONNECTION_STRING`

## Тесты

В файле `coverage.out` находится покрытие бизнес-логики тестами. Сгенерировать покрытие можно с помощью команды `make gen-coverage` .

Посмотреть покрытие можно с помощью команды `make show-coverage`

```
go tool cover -func ./coverage.out
subscription-service/internal/app/usecase/subscription.go:31:   NewSubscription 100.0%
subscription-service/internal/app/usecase/subscription.go:47:   Create          87.5%
subscription-service/internal/app/usecase/subscription.go:73:   Read            83.3%
subscription-service/internal/app/usecase/subscription.go:85:   Update          50.0%
subscription-service/internal/app/usecase/subscription.go:108:  update          37.5%
subscription-service/internal/app/usecase/subscription.go:138:  Delete          83.3%
subscription-service/internal/app/usecase/subscription.go:150:  List            75.0%
subscription-service/internal/app/usecase/subscription.go:159:  Sum             100.0%
total:                                                          (statements)    66.1%
```

## Структура папок проекта

```
├── api -> Содержит openapi файл с описанием контрактов и схем сервиса.
├── cmd
│   └── core -> Точка входа при запуске сервиса.
├── internal
│   ├── adapter -> Реализации интерфейсов из repo.
│   │   ├── db -> Адаптер к базе.
│   │   └── repo -> Адаптеры репозиторного слоя.
│   │       └── mock -> Моковые реализации адаптеров репозиторного слоя.
│   ├── app
│   │   ├── entity -> Сущности бизнес-логики.
│   │   └── usecase -> Бизнес-логика.
│   ├── config
│   ├── controller
│   │   └── http
│   │       └── gen -> Содержит в себе сгенерированные из openapi структуры и сервер.
│   ├── migrations -> Миграции и код для встраивания миграций в бинарный файл сборки.
│   ├── pkg
│   │   └── utils
│   └── port -> Описание интерфейсов для связи с внешними системами.
└── pkg
    └── client -> Сгенерированный из openapi клиент к сервису.
```
