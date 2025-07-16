# ДЗ: Масштабируемая подсистема диалогов

Цель: реализовать масштабируюемую подсистему диалогов путем горизонтального масштабирования хранилищ на запись с помощью шардинга.

Шардинг будем реализовывать с помощью утилиты Citus.

### План мероприятий (настройка)

1. Добавляем Citus + shard1 и shard2 в docker compose и запускаем его
```
docker compose up -d
```
2. Настройка Citus (master)
 - Заходим в postgres_master
```
docker exec -it postgres_master
psql -U postgres -h localhost -d postgres
```
- Установливаем расширение citus
```
CREATE EXTENSION IF NOT EXISTS citus;
```
- Устанавливаем координатор (мастер нода которая распределяет данные по шардам)
```
SELECT citus_set_coordinator_host('postgres_master');
```
- Добавляем ноды (worker nodes)
```
SELECT * FROM citus_add_node('postgres_shard1', 5432);
SELECT * FROM citus_add_node('postgres_shard2', 5432);
```

3. проверяем на мастере список активных воркер нод
```
SELECT * FROM citus_get_active_worker_nodes();
```
ссылка на скриншот списка активных воркер нод
https://github.com/proweb-zone/otus_social_network/tree/main/load_test/lesson_5/list-shards.jpg

4. Проводим миграцию таблиц
Выходим из docker и проводим миграцию
```
make migration-up
```

5. Создаем распределительную таблицу на мастере (postgres_master)
- Заходим в postgres_master
```
docker exec -it postgres_master
psql -U postgres -h localhost -d postgres
SELECT create_distributed_table('dialog', 'user_id_recipient');
```
Важно чтобы на воркер нодах таблицы были удалены.

Шардировать таблицу будем по полю user_id_recipient (так же можно по id), главное чтобы данные поля были primary key.

- проверка распределительных таблиц
```
SELECT * FROM citus_tables;
```

Ссылка на скриншот настройки распределительной таблицы
https://github.com/proweb-zone/otus_social_network/tree/main/load_test/lesson_5/sharding-table.jpg

6. Проверка результата
- После создания правил распределения таблиц (правила шардирования), у Вас на шардах (worker nodes) должна создаться таблица dialog.
- Далее через Postman регистрируем в приложении двух пользователей и авторизуем их
- Пишем друг другу сообщение и проверяем данные в таблице dialog на всех worker nodes. Данные у вас должны распределиться равномерно по всем шардам.
