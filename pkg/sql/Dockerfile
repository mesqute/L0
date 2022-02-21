
FROM postgres:14-alpine

# копирование файла для использования встроенного механизма инициализации postgre контейнера
COPY postgres_init_sql.sql /docker-entrypoint-initdb.d/postgres_init_sql.sql
