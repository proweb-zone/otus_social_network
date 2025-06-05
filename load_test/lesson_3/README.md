# Репликация - практическое применение

### Настройка репликации (sync, async)
1. Добавляем в docker-compose дополнительные сервисы slave и async_slave

2. Заходим в master(docker) и добавим роль "replicator", для репликации
```
    # после
    psql -U postgres -d otus
```

3. Заходим в master(docker) и делаем backup базы данных
```
docker exec -it postgres_master bash
mkdir /pgslave
pg_basebackup -h postgres_master -D /pgslave -U replicator -v -P --wal-method=stream
```

4. Выгружаем дамп в корень нашего проекта и создаем файл "standby.signal"
```
docker cp postgres_master:/pgslave pg_data/pgslave/
touch pg_data/pgslave/standby.signal
```

5. Проделываем со вторым slave ту же саму операцию
```
docker cp postgres_master:/pgslave pg_data/pgasyncslave/
touch pg_data/pgasyncslave/standby.signal
```

### Нагрузочное тестирование (K6)

При проведении нагрузочное тестирования, на все тесты должны быть написаны сценарии (K6).

#### Стратегия сценария

- Выбираем два EndPoints /user/get/{id} и /user/search на чтение
- Отправляем 1000 запросов на выбранные EndPoints в течении 4 минут

### Нагрузочное тестирование (без реплик)

1. Настраиваем наш проект на работу без репликации (только один сервер postgres в стеке)
2. Запускаем тесты и проводим нагрузочное тестирование
```
# Результат тестирования находиться по пути
load_test/lesson_3/graphics/
```

### Нагрузочное тестирование (sync)

1. Настраиваем наш класте на работу в режиме master - slave - slave
2. Запускаем тесты и проводим нагрузочное тестирование
```
# Результат тестирования находиться по пути
load_test/lesson_3/graphics/
```

####Создаем кворумную репликацию Postgresql

### На Master

1.  Делаем  дамп roles и schema (если их не делать то при копировании бд на slave произойдет ошибка при подписке на master)
```
pg_dumpall -U postgres -r -h postgres_master -f /var/lib/postgresql/data/roles.dmp
pg_dump -U postgres -Fc -h postgres_master -f /var/lib/postgresql/data/schema.dmp -s postgres
```

2. Заходим в конфиг postgresql.conf и меняем следующие параметры:
   ```
   wal_level = logical
   ```

3. Перезапускаем сервер(master)
4. заходим в psql и создаем публикацию
```
GRANT CONNECT ON DATABASE  otus TO replicator;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO replicator;
create publication pg_pub for table users;
```

### На SLAVE

1. На всех slave заходим в psql и создаем подписку на master
```
CREATE SUBSCRIPTION pg_sub CONNECTION 'host=postgres_master port=5432 user=replicator password=pass dbname=otus' PUBLICATION pg_pub;
```
