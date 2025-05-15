## otus_social_network (Документация разработчика)

### Описание

Социальная сеть - домшнее задание онлайн школы OTUS.

### Стек
Golang+Postgresql+Docker

###  Установка

1) Установите и запустите docker на вашем устройстве.

2) Скачайте проект с удаленного репозитария
```
git clone https://github.com/proweb-zone/otus_social_network.git
```

3) Перейдите в проект
```
cd otus_social_network
```
4) Запустите проект в docker с помощью docker-compose
```
docker compose up -d
```
5) Проведите миграцию (необходимо при первичном запуске)
```
docker exec -it app_socnet go run /otus/app/cmd/migration/main.go -action up
```

### Пользовательская документация

Пользовтельская документация находиться в репозитарии ввиде (postman коллекции). Загрузите коллекцию в Postman для изучения.
