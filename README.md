# B-ASS (Banek-as-a-Service)

Очень простое апи для получения анекдотов с baneks.site и baneks.ru

## Фичи

- [x] Случайный анекдот из двух источников (api/baneks/random).
- [x] Round Robin Балансировщик для /random, чтобы не наглеть.
- [x] Получение конкретного анекдота по slug из banek.site (api/baneks/:slug)
- [x] Получение случайных мемов (api/memes/random?year={year})

## Технологии

- [GoLang](https://go.dev/)
- [Echo](https://echo.labstack.com/)

## Использование

- Установите пакеты с помощью команды:

```sh
$ go mod tidy
```

Запустите:

```ssh
go run cmd/main.go
```

Или соберите:
```ssh
go build cmd/main.go
```

### Зачем вы разработали этот проект?

Нужны были анекдоты. Анекдоты теперь есть.
