# homework_extra

## Описание

Данный проект представляет собой сервис объявлений, состоящий из двух основных 
сущностей: [`ad`](https://github.com/papey08/golang-fintech/blob/homework/extra/homework_extra/internal/model/ads/ads.go)
(объявление) и [`user`](https://github.com/papey08/golang-fintech/blob/homework/extra/homework_extra/internal/model/users/users.go)
(пользователь).

### Структура проекта

```text
├── cmd
│   └── server
│       └── main.go // точка входа в приложение
│
├── doc
│   └── coverage.html // отчёт о покрытии тестами в формате html
│
├── internal
│   ├── adapters // слой БД
│   │   ├── adrepo // хранилище объявлений
│   │   └── user_repo // хранилище пользователей
│   │
│   ├── app // слой бизнес-логики (usecase)
│   │   ├── adrepo_mocks
│   │   ├── app.go // интерфейс приложения
│   │   ├── my_app.go // реализация интерфейса приложения
│   │   ├── my_app_test.go
│   │   └── user_repo_mocks
│   │
│   ├── model // слой сущностей (entities)
│   │   ├── ads // описание объявления
│   │   ├── errs
│   │   ├── filter
│   │   └── users // описание пользователя
│   │
│   └── ports // сетевой слой (infrastructure)
│       ├── grpc // gRPC-сервер
│       └── httpgin // HTTP-сервер на gin
│
├── migrations
│   └── adrepo_init.sql // скрипт для конфигурации adrepo
│
├── config.yml // файл с конфигами
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

### Стркуктура `ad`

* `ID` — целочисленный идентификатор
* `Title` — заголовок (до 100 символов)
* `Text` — текст объявления (до 500 символов)
* `AuthorID` — ID автора (пользователя, который разместил объявление)
* `Published` — опубликовано ли объявление (`true` или `false`)
* `CreationDate` — дата создания
* `UpdateDate` — дата последнего изменения

### Структура `user`

* `ID` — целочисленный идентификатор
* `Nickname` — имя пользователя
* `Email` — электронная почта пользователя

### Бизнес-логика

Реализован CRUD для пользователя (изменять можно `Nickname` и `Email`), CRUD 
для объявления (изменять можно `Text`, `Title` и `Published`), причём изменять 
и удалять объявления может только автор, а при изменении любого из полей 
обновляется поле `UpdateDate`. Кроме того можно найти несколько объявлений по 
префиксу заголовка или текста, а также отфильтровать все объявления, используя 
один или несколько доступных фильтров:

* только опубликованные
* по автору
* по дате создания

### Поддерживаемые методы

* [HTTP](https://github.com/papey08/golang-fintech/blob/homework/extra/homework_extra/internal/ports/httpgin/router.go)
* [gRPC](https://github.com/papey08/golang-fintech/blob/homework/extra/homework_extra/internal/ports/grpc/pb/service.proto)

### Используемые технологии

* go 1.19
* PostgreSQL — хранение объявлений
* MongoDB — хранение пользователей
* [Gin Web Framework](https://github.com/gin-gonic/gin)
* gRPC
* Docker

## Запуск

### С помощью Docker

```shell
$ docker-compose up
```

### Локально

Самостоятельно сконфигурировать базы данных PostgreSQL ([скрипт для конфигурации adrepo](https://github.com/papey08/golang-fintech/blob/homework/extra/homework_extra/migrations/adrepo_init.sql))
и MongoDB, изменить файл *[config.yml](https://github.com/papey08/golang-fintech/blob/homework/extra/homework_extra/config.yml)*,
после чего выполнить команды:

```shell
$ go mod download
$ go run cmd/server/main.go
```
