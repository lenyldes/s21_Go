-- Пользователь дополнительных бэкендов может только читать таблицы магазина.
\getenv db_name POSTGRES_DB
\getenv reader_password SHOP_READONLY_PASSWORD

CREATE USER shop_reader WITH PASSWORD :'reader_password';
GRANT CONNECT ON DATABASE :"db_name" TO shop_reader;
GRANT USAGE ON SCHEMA public TO shop_reader;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO shop_reader;
