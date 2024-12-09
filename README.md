# Meower post service

Сервис ответственный за работу с сообщениями пользователей. 

Предоставляет доступ к сообщениям посредством grpc c интерфейсом и сообщениями описанным в [api](https://github.com/Karzoug/meower-api/tree/main/proto/post). События сервиса рассылаются в брокер [outbox сервисом](https://github.com/Karzoug/meower-post-outbox).

### Стек
- Основной язык: go
- База данных: mongoDB
- Брокер: kafka
- Наблюдаемость: opentelemetry, jaeger, prometheus
- Контейнеры: docker, docker compose

## Дальнейшее развитие

- [ ] дополнительные поля для сообщений: ссылки, картинки, видео,
- [ ] кастомные метрики,
- [ ] кеширование сообщений,