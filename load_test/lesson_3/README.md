# Репликация - практическое применение

### ReplicationRoutingDataSource

Для реализации репликации (в рамках тестового задания) я написал обработку переключения между slave на Golang. Скрипт переключает на slave и проверяет доступность сервера, если реплика не доступно происходит переключение на другую.

-- Скрипт находиться по пути app/internal/db/postgres/replicationRoutingDataSource.go

### Настройка репликации (sync, async)
1. Добавляем в docker-compose дополнительные сервисы slave и async_slave

2. Заходим в master(docker) и добавим роль "replicator", для репликации
```
    psql -U postgres -d otus
    create role replicator with login replication password 'pass';
```

1. Добавляем запись в конфиг /var/lib/postgres/data/pg_hba.conf на master
```
host    replication     replicator      172.33.0.3/24           md5
host    replication     replicator      172.33.0.4/24           md5
host    replication     replicator      172.33.0.5/24           md5
host    replication     replicator      172.33.0.6/24           md5
```

1. Добавляем запись в конфиг /var/lib/postgres/data/postgresql.conf на master
```
ssl = off
wal_level = replica
max_wal_senders = 4

synchronous_commit = on
synchronous_standby_names = 'FIRST 1 (pgslave, pgasyncslave)'
```

После чего перезапускаем сервер (master)

2. Заходим в master(docker) и делаем backup базы данных
```
docker exec -it postgres_master bash
mkdir /pgslave
pg_basebackup -h postgres_master -D /pgslave -U replicator -v -P --wal-method=stream
```

1. Выгружаем дамп в корень нашего проекта и создаем файл "standby.signal"
```
docker cp postgres_master:/pgslave pg_data/pgslave/
touch pg_data/pgslave/standby.signal
```

1. Проделываем со всеми репликами (slave) ту же саму операцию
```
docker cp postgres_master:/pgslave pg_data/pgasyncslave/
touch pg_data/pgasyncslave/standby.signal
```

### Нагрузочное тестирование (K6)

При проведении нагрузочное тестирования, на все тесты должны быть написаны сценарии (K6).

#### Стратегия сценария

- Выбираем два EndPoints /user/get/{id} и /user/search на чтение
- Отправляем 500 запросов на выбранные EndPoints в течении 4 минут

### Нагрузочное тестирование (без реплик)

1. Настраиваем наш проект на работу без репликации (только один сервер postgres в стеке)
2. Запускаем тесты и проводим нагрузочное тестирование
```
# Результат тестирования находиться по пути
load_test/lesson_3/graphics/none_replication_500.html
```

### Нагрузочное тестирование (sync)

1. Настраиваем наш класте на работу в режиме master - slave - slave
2. Запускаем тесты и проводим нагрузочное тестирование
```
# Результат тестирования находиться по пути
load_test/lesson_3/graphics/test_repication_500.html
```

 ###  Результат (sync) репликации
  По итогам проведения нагрузочного тестирования видно, что отказоустойчивость системы повысилась, но latency незначительно возрасло - этот over head вызван постоянным переключением между slave и проверки его доступности.


## Создаем кворумную(логическую) репликацию Postgresql

#### На Master

1.  Делаем  дамп roles и schema (если их не делать то при копировании бд на slave произойдет ошибка при подписке на master)
```
pg_dumpall -U postgres -r -h postgres_master -f /var/lib/postgresql/data/roles.dmp
pg_dump -U postgres -Fc -h postgres_master -f /var/lib/postgresql/data/schema.dmp -s postgres
```

2. Заходим в конфиг postgresql.conf и меняем следующие параметры:
   ```
   wal_level = logical
   synchronous_commit = off
   ```

3. Перезапускаем сервер(master)
4. заходим в psql и создаем публикацию
```
GRANT CONNECT ON DATABASE  otus TO replicator;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO replicator;
create publication otus_pub for table users;
```
-- Скриншот созданных публикаций /load_test/lesson_3/scripts/public_list.png

### На SLAVE

1. На всех slave заходим в psql и создаем подписку на master
```
# Название подписки 'otus_sub' - на каждом slave должно быть уникальным
CREATE SUBSCRIPTION otus_sub CONNECTION 'host=postgres_master port=5432 user=replicator password=pass dbname=otus' PUBLICATION otus_pub;
```
-- Скриншот созданных подписок /load_test/lesson_3/scripts/subscruption_list.png

2. На каждом slave заходим в конфиг postgresql.conf и включаем логическую репликацию
    ```
      wal_level = logical
    ```
3. Перезапускаем сервера и если все выполнено правильно,то логическая репликация должна заработать
    -- Скриншот логов о правильной настройки репликации /load_test/lesson_3/scripts/start_logic_replication.png

### Справочная информация

- Вывести список публикаций (на master)
```
SELECT * FROM pg_publication;
```
-  Вывести список подписок (на slave)
```
SELECT * FROM pg_stat_subscription;
```

### Результат

1. При логической репликации, если произойдет сбой на реплицируемой сервере, то часть данных может потеряться, потому что логическая репликация работает в режиме synchronous_commit = off
2. Вторая проблема в доставке данных может возникнут в консистенции данных, если в репликаторе данные с таким идентификатором уже есть то возникнет конфликт

- Скриншот по нагрузочному тестировани на запись находиться по адресу /load_test/lesson_3/scripts/load_test_write.html

 - Скриншот ошибки репликации находиться по адресу /load_test/lesson_3/scripts/error_replication.png
