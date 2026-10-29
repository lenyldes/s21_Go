# Высоконагруженный обратный прокси Nginx, балансировка нагрузки и микросервисная архитектура (Nginx Reverse Proxy, Load Balancing & Microservices)

[![Go Version](https://img.shields.io/badge/Go-1.25%20%7C%201.26-00ADD8?style=flat&logo=go&logoColor=white)](https://golang.org/)
[![Nginx](https://img.shields.io/badge/Nginx-1.30.6-009639?style=flat&logo=nginx&logoColor=white)](https://nginx.org/)
[![gRPC](https://img.shields.io/badge/gRPC-Protobuf-244f5c?style=flat&logo=grpc&logoColor=white)](https://grpc.io/)
[![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL%20(pgx%20v5)-336791?style=flat&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Clean Architecture](https://img.shields.io/badge/Architecture-Clean%20Architecture-blue?style=flat)](#архитектура-и-сетевая-топология)
[![Uber Fx](https://img.shields.io/badge/DI-Uber%20Fx-black?style=flat)](https://github.com/uber-go/fx)
[![Docker Compose](https://img.shields.io/badge/Docker-Compose%20(8%20Services)-2496ED?style=flat&logo=docker&logoColor=white)](#оркестрация-и-развертывание-в-docker-compose)
[![Status Completed](https://img.shields.io/badge/Status-Completed-success?style=flat)](#)

Сделал отказоустойчивую распределённую инфраструктуру для микросервисного интернет‑магазина за обратным прокси‑сервером Nginx. Настроил интеллектуальную балансировку нагрузки между пулом Read‑Only реплик и Master‑узлом с весовым распределением (2:1:1), разделил трафик чтения и записи (CQRS‑паттерн на уровне прокси), внедрил многоуровневое кеширование HTTP‑ответов (`proxy_cache`) с дедупликацией нагрузки, сквозную компрессию контента GZIP и защищённый HTTPS/TLS‑канал с локальным центром сертификации. Разбил архитектуру сервисов на два изолированных Go‑микросервиса (Shop REST API и Auth gRPC Service) на базе Clean Architecture и DI‑контейнера Uber Fx. Сервисы общаются по высокопроизводительному протоколу gRPC/Protobuf, используют изолированные базы данных PostgreSQL с разграничением прав на уровне ролей и упакованы в согласованный стек из 8 контейнеров в Docker Compose.

---

## О проекте

В этом проекте я решал комплекс системных и инфраструктурных задач для высоконагруженного бэкенда: нужно было грамотно распределить сетевой трафик, защитить ядро системы от перегрузок, сократить время отклика клиентов и обеспечить высокую доступность (High Availability) сервисов:

- **Обратный прокси и маршрутизация:** Сделал единую входную точку (API Gateway / Ingress) на базе Nginx 1.30 — она маршрутизирует запросы к REST API магазина, интерактивной документации Swagger UI, панели администрирования баз данных pgAdmin и раздаче статических ресурсов.
- **Взвешенная балансировка нагрузки (2:1:1) и разделение чтения/записи:** Сегрегировал сетевой трафик по HTTP‑глаголам с помощью директивы `map`. Мутирующие запросы (`POST`, `PUT`, `DELETE`, `PATCH`) отправляю строго на главный Master‑экземпляр (`shopapp`). Запросы на чтение (`GET`) балансирую по алгоритму Weighted Round‑Robin между тремя экземплярами бэкенда (`shopapp` с весом 2, `shopapp_read1` с весом 1, `shopapp_read2` с весом 1).
- **Кеширование ответов на стороне прокси:** Настроил общую зону памяти `proxy_cache` на 10 ГБ с временем жизни кеша 1 минута для успешных ответов (HTTP 200). Добавил диагностические заголовки `X-Cache-Status` (`HIT` / `MISS`) и `X-Upstream-Addr`, чтобы удобно отслеживать распределение запросов.
- **Потоковая GZIP‑компрессия:** Реализовал динамическое сжатие текстовых форматов (`JSON`, `HTML`, `CSS`, `JS`, `XML`) с уровнем сжатия 6, при этом исключил бинарные медиафайлы (`PNG`, `JPEG`) — так не тратим процессорные мощности впустую.
- **Локальный HTTPS/TLS:** Поднял самоподписанный X.509 сертификат для домена `be3.shop` с расширением SAN (Subject Alternative Name) через OpenSSL. Поддерживаю одновременную работу портов 80 (HTTP) и 443 (HTTPS).
- **Межсервисный протокол gRPC / Protobuf:** Выделил отдельный микросервис аутентификации `Auth` — он предоставляет быстрый RPC‑интерфейс для регистрации, выпуска и ротации двухтокенных связок JWT (Access/Refresh), сброса паролей и аутентификации запросов клиентов.
- **Безопасность баз данных и Read‑Only пользователи:** Использую два независимых инстанса PostgreSQL (`shopdb` и `authdb`). Создал роль `shop_reader` с правами только на чтение (`SELECT`) — её используют реплики чтения `shopapp_read1` и `shopapp_read2`.
- **Чистая архитектура и Dependency Injection:** Спроектировал Go‑сервисы по принципам Clean Architecture (слои Domain, Web/API, Repository/Data), внедрил зависимости и управление жизненным циклом через Uber Fx.

---

## Ключевые навыки и освоенные концепции

### 1. Nginx Reverse Proxy и интеллектуальная маршрутизация

- **Проксирование запросов (`proxy_pass`):** Перенаправляю входящие HTTP‑запросы во внутреннюю виртуальную сеть Docker Compose к нужным контейнерам — при этом их внутренние порты не видны снаружи.
- **Преобразование путей и редиректы:**
  - Префикс `/api/` прозрачно подменяю на `/api/v1/` приложения;
  - Запрос `/api/v1` возвращает HTTP‑редирект `302 Found` на интерактивный Swagger UI (`/swagger/index.html`);
  - Статический путь `/` настроил на прямую раздачу файлов `index.html` и `image.png` из файлового хранилища `/data/www`;
  - Маршрут `/admin/` настроил на прозрачное проксирование веб‑интерфейса pgAdmin — корректно передаю заголовок пути `X-Script-Name /admin`.
- **Встроенный мониторинг состояния (`stub_status`):** Реализовал защищённый эндпоинт `/status`, который отдаёт метрики Nginx в реальном времени: число активных соединений, количество принятых и обработанных подключений и запросов, а также статус очередей чтения и записи.

### 2. Балансировка нагрузки и CQRS на уровне прокси (Read/Write Splitting)

- **Динамический маппинг методов (`map $request_method`):** Я настроил маппинг так, чтобы целевая группа бэкендов выбиралась по семантике запроса:
  ```nginx
  map $request_method $shop_backend {
      GET     shop_read;
      default shop_write;
  }
  ```
- **Пул записи (`upstream shop_write`):** Все методы, которые меняют состояние данных (`POST`, `PUT`, `DELETE`, `PATCH`), я направил строго на основной инстанс `shopapp`, подключённый к БД с правами на запись. Так я исключил попытки записи в Read‑Only реплики и предотвратил рассинхронизацию данных.
- **Взвешенный пул чтения (`upstream shop_read`):** Запросы на чтение (`GET`) я распределил по трём узлам с весовыми коэффициентами:
  - `shopapp:8080 weight=2` (основной инстанс, берёт 50 % всех GET‑запросов);
  - `shopapp_read1:8080 weight=1` (первая реплика чтения, берёт 25 % GET‑запросов);
  - `shopapp_read2:8080 weight=1` (вторая реплика чтения, берёт 25 % GET‑запросов);
  - Итоговая пропорция распределения: **2:1:1** (50 % : 25 % : 25 %).

### 3. Кеширование запросов на уровне прокси (`proxy_cache`)

- **Настройка хранилища кеша (`proxy_cache_path`):** Я реализовал двухуровневую файловую структуру `levels=1:2`, выделил 10 МБ разделяемой памяти для зоны ключей `my_cache:10m` (она вмещает примерно 80 000 ключей), ограничил максимальный объём на диске 10 ГБ (`max_size=10g`) и настроил автоматическую очистку неиспользуемых записей через 60 минут (`inactive=60m`).
- **Правила инвалидации и валидности (`proxy_cache_valid`):** Успешные ответы с кодом `200 OK` я кеширую ровно на 1 минуту (`1m`). Ошибочные ответы и запросы, требующие авторизации с невалидным токеном (HTTP 401), не кешируются — так система мгновенно реагирует на изменение прав доступа.
- **Отладочные заголовки трассировки:**
  - `add_header X-Cache-Status $upstream_cache_status always;` — этот заголовок сообщает клиенту статус кеша (`MISS` — запрос ушёл на бэкенд, `HIT` — ответ сразу из памяти или диска Nginx);
  - `add_header X-Upstream-Addr $upstream_addr always;` — фиксирует внутренний IP/порт бэкенда, который обработал запрос (при `HIT` заголовок будет пустым, потому что бэкенд не задействован).

### 4. Потоковая GZIP‑компрессия трафика

- **Оптимизация полосы пропускания:** Я снизил сетевую нагрузку и ускорил передачу данных клиентам за счёт сжатия ответов на лету.
- **Тонкая настройка директив:**
  - `gzip_comp_level 6` — это даёт оптимальный баланс между степенью сжатия и нагрузкой на CPU сервера;
  - `gzip_min_length 1100` — отсекает микроответы меньше 1,1 КБ: их сжатие невыгодно из‑за накладных расходов на заголовки;
  - `gzip_proxied any` — включает сжатие для всех проксированных запросов;
  - `gzip_vary on` — добавляет HTTP‑заголовок `Vary: Accept-Encoding`, чтобы корректно работали промежуточные CDN и браузерные кеши;
  - `gzip_types` — я оставил строгий белый список MIME‑типов (`text/*`, `application/json`, `application/javascript`, `text/xml` и др.). Изображения (`image/png`, `image/jpeg`) не сжимаю: у них и так высокая степень встроенной энтропии.

### 5. Локальный HTTPS/TLS и криптографическая безопасность

- **Генерация самоподписанного сертификата X.509:** Я написал автоматизированный скрипт `generate-cert.sh` на OpenSSL:
  - 2048‑битный RSA‑ключ (`rsa:2048`);
  - хэш‑алгоритм SHA‑256;
  - расширение `subjectAltName = DNS:be3.shop` — оно нужно, чтобы домен корректно валидировался современными браузерами и HTTP‑клиентами;
  - срок действия сертификата — 365 дней;
  - права на закрытый ключ я ограничил командой `chmod 600 be3.shop.key`.
- **Локальный DNS‑маппинг:** Настройте маршрутизацию доменного имени `be3.shop` на локальный интерфейс петли `127.0.0.1` через системный файл `/etc/hosts`.

### 6. Микросервисная архитектура и межсервисный gRPC (Protobuf)

- **Декомпозиция ответственности:** Я выделил отдельный сервис `authapp`, который полностью отвечает за криптографическую обработку паролей (bcrypt), формирование JWT‑токенов (HS256) и аутентификацию пользователей.
- **Высокоскоростной транспорт gRPC:** Сервис магазина (`shopapp`) взаимодействует с `authapp` по бинарному протоколу HTTP/2 через сгенерированные Protobuf‑контракты (`auth.proto`). Это даёт минимальные задержки и строгую типизацию сообщений.
- **Централизованная верификация токенов:** Middleware `UserAuthenticator` в `shopapp` для каждого защищённого запроса валидирует Bearer‑токен через gRPC‑вызов `ValidateToken` к сервису Auth — так я гарантирую актуальность сессии.

### 7. Изоляция данных в PostgreSQL и Read‑Only реплики

- **Разделение баз данных:** Я создал независимые базы данных `shop` и `shop_auth` на отдельных инстансах СУБД, чтобы соблюсти принцип Database‑per‑Service.
- **Принцип наименьших привилегий (Least Privilege):**
  - Главный инстанс `shopapp` работает под суперпользователем `postgres` с полными правами DDL/DML;
  - Инстансы `shopapp_read1` и `shopapp_read2` подключаются под специально созданной учётной записью `shop_reader`, у которой есть только право `SELECT` на таблицы схемы `public`. Любая попытка выполнить `INSERT`/`UPDATE`/`DELETE` на репликах чтения приведёт к ошибке СУБД.

### 8. Чистая архитектура (Clean Architecture) и DI (Uber Fx)

- **Многослойный дизайн бэкенда Shop:**
  - `domain`: чистые бизнес‑сущности (`Product`, `Client`, `Supplier`, `Address`, `Image`) и интерфейсы репозиториев (`ProductRepository`, `SupplierRepository` и др.);
  - `web`: DTO с JSON‑тегами, HTTP‑хэндлеры на базе `net/http` (с роутингом Go 1.22+), DTO‑мапперы, валидация входных структур (`DisallowUnknownFields`);
  - `repository`: слой данных на `pgxpool.Pool` (PostgreSQL), DAO‑структуры, SQL‑мапперы и транзакционные операции (например, атомарное обновление адресов).
- **Граф зависимостей Uber Fx:** Я реализовал полную декларативную инициализацию пула БД, gRPC‑клиентов, репозиториев, хэндлеров и запуск HTTP‑сервера через `fx.Provide` и `fx.Invoke`.

---

## Архитектура и сетевая топология

### Схема взаимодействия компонентов

```text
                               Клиентские запросы
                        (Браузер, cURL, Swagger UI)
                                     │
                    ┌────────────────┴────────────────┐
                    │  HTTPS (443) / HTTP (80)        │
                    ▼                                 ▼
         ┌─────────────────────────────────────────────────────────┐
         │                   NGINX Reverse Proxy                   │
         │                  (Container: nginx)                     │
         │   TLS Termination / GZIP / Caching / Map Routing        │
         └───────┬──────────────┬──────────────┬─────────────┬─────┘
                 │              │              │             │
        GET /    │     /status  │     /admin/  │             │ /api/v1/
   (Static files)│ (stub_status)│ (pgAdmin Web)│             │
                 ▼              ▼              ▼             ▼
          ┌─────────────┐ ┌───────────┐ ┌─────────────┐ ┌─────────────────────────┐
          │  /data/www  │ │   Nginx   │ │   pgAdmin   │ │   Load Balancer (Map)   │
          │ index.html  │ │  Metrics  │ │ (pgadmin)   │ │ GET -> shop_read (2:1:1)│
          │ image.png   │ └───────────┘ └──────┬──────┘ │ Other -> shop_write (1) │
          └─────────────┘                      │        └────────────┬────────────┘
                                               │                     │
                        ┌──────────────────────┘                     │
                        │                                            │
                        ▼                                            ▼
                 ┌──────────────┐                       ┌─────────────────────────┐
                 │  PostgreSQL  │◄──────────────────────┤  shopapp (Port 8080)    │ (Weight: 2)
                 │   (shopdb)   │   Read/Write          │  (Master: RW Instance)  │
                 │              │                       └────────────┬────────────┘
                 │              │◄──────────────────────┐            │
                 │              │   Read-Only (SELECT)  │            │
                 │              │                       ├────────────┼────────────┐
                 │              │◄────────┐             │            │            │
                 └──────────────┘         │             │            │            │
                                          │             ▼            ▼            ▼
                                   ┌──────┴──────────────┐   ┌───────┴─────────────┐
                                   │ shopapp_read1 (8080)│   │ shopapp_read2 (8080)│
                                   │ (Weight: 1, RO)     │   │ (Weight: 1, RO)     │
                                   └──────────┬──────────┘   └───────┬─────────────┘
                                              │                      │
                                              │    gRPC :50051       │
                                              └──────────┬───────────┘
                                                         │
                                                         ▼
                                              ┌─────────────────────┐
                                              │       authapp       │
                                              │ (Auth gRPC Service) │
                                              └──────────┬──────────┘
                                                         │
                                                         ▼
                                              ┌─────────────────────┐
                                              │     PostgreSQL      │
                                              │      (authdb)       │
                                              └─────────────────────┘
```

---

### Стек технологий

| Технология / Компонент | Версия / Спецификация | Назначение в проекте |
| :--- | :--- | :--- |
| **Go** | `1.25.6` (Shop) / `1.26.8` (Auth) | Основной язык микросервисов, стандартная библиотека `net/http` |
| **Nginx** | `1.30.6-alpine` | Обратный прокси, балансировщик нагрузки 2:1:1, кеширующий сервер, SSL/TLS, GZIP |
| **gRPC** | `v1.84.0` | Высокопроизводительный бинарный межсервисный транспорт |
| **Protocol Buffers** | `v1.36.12` | Схема контрактов авторизации (`api/auth.proto`) |
| **PostgreSQL** | `18.1-alpine` | Реляционные базы данных (`shopdb` и `authdb`) |
| **pgx** | `v5.11.0` | Высокоэффективный Go-драйвер и пул соединений `pgxpool` |
| **Uber Fx** | `v1.24.0` | Фреймворк внедрения зависимостей (Dependency Injection) |
| **golang-jwt/jwt** | `v5.3.1` | Выпуск, верификация и подпись токенов JWT (HS256) |
| **golang.org/x/crypto** | `v0.58.0` | Криптографическое хэширование паролей (bcrypt) |
| **Swagger / Swaggo** | `v2.0.2` | Автоматическая генерация и отображение интерактивной OpenAPI-документации |
| **pgAdmin 4** | `dpage/pgadmin4` | Веб-интерфейс администрирования PostgreSQL |
| **Docker & Compose** | `Compose Spec` | Оркестрация стека из 8 независимых контейнеров |
| **OpenSSL** | `x509 / RSA 2048` | Генерация самоподписанного TLS-сертификата с поддержкой SAN |

---

## Маршрутизация Nginx и спецификация API

### 1. Маршрутизация Nginx

| Путь (Location) | Назначение / Поведение | Проксирование / Директива |
| :--- | :--- | :--- |
| `/` | Раздача статического контента | `root /data/www` (`index.html`, `image.png`) |
| `/status` | Статистика соединений и состояния Nginx | Директива `stub_status` |
| `/admin/` | Веб-интерфейс администрирования pgAdmin | `proxy_pass http://pgadmin/` с `X-Script-Name /admin` |
| `/api/v1` | Точное совпадение: редирект на документацию Swagger | HTTP `302 Found` -> `/swagger/index.html` |
| `/swagger/` | Статика и интерфейс Swagger UI | `proxy_pass http://shopapp:8080` |
| `/api/v1/` | Полнофункциональный REST API магазина | `proxy_pass http://$shop_backend;` (балансировка `GET`, кеш `1m`, заголовки отладки) |
| `/api/` | Короткий псевдоним API | `proxy_pass http://shopapp:8080/api/v1/;` |

---

### 2. Спецификация REST API Shop Service

#### Аутентификация и сессии (через gRPC-вызовы к Auth Service)

| Метод | Путь | Доступ | Назначение |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/register` | Публичный | Регистрация нового пользователя (login, email, password, profile) |
| `POST` | `/api/v1/auth` | Публичный | Аутентификация, выдача пары `access_token` и `refresh_token` |
| `POST` | `/api/v1/refresh` | Публичный | Ротация и обновление токенов по действующему Refresh-токену |
| `POST` | `/api/v1/reset` | Публичный | Запрос на сброс пароля пользователя |
| `POST` | `/api/v1/change_password` | Bearer JWT | Смена текущего пароля на новый |

#### Товары (Products)

| Метод | Путь | Доступ | Назначение | Балансировка / Кеш |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/products` | Bearer JWT | Список доступных товаров с пагинацией (`page`, `size`) | Да (2:1:1), Кеш 1m |
| `GET` | `/api/v1/products/{id}` | Bearer JWT | Детальная информация о товаре по UUID | Да (2:1:1), Кеш 1m |
| `POST` | `/api/v1/products` | Bearer JWT | Создание новой позиции товара в каталоге | Master only (`shopapp`) |
| `PATCH` | `/api/v1/products/{id}/stock` | Bearer JWT | Уменьшение складского остатка товара | Master only (`shopapp`) |
| `DELETE` | `/api/v1/products/{id}` | Bearer JWT | Удаление товара по UUID | Master only (`shopapp`) |

#### Адреса (Addresses)

| Метод | Путь | Доступ | Назначение | Балансировка / Кеш |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/addresses` | Bearer JWT | Список адресов с пагинацией | Да (2:1:1), Кеш 1m |
| `POST` | `/api/v1/addresses` | Bearer JWT | Добавление нового адреса (country, city, street) | Master only (`shopapp`) |
| `PUT` | `/api/v1/addresses/{id}` | Bearer JWT | Обновление информации об адресе | Master only (`shopapp`) |
| `DELETE` | `/api/v1/addresses/{id}` | Bearer JWT | Безопасное удаление адреса (если нет внешних связей) | Master only (`shopapp`) |

#### Клиенты (Clients)

| Метод | Путь | Доступ | Назначение | Балансировка / Кеш |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/clients` | Bearer JWT | Список зарегистрированных клиентов с пагинацией | Да (2:1:1), Кеш 1m |
| `POST` | `/api/v1/clients` | Bearer JWT | Создание профиля клиента с привязкой к адресу | Master only (`shopapp`) |
| `PUT` | `/api/v1/clients/{id}/address`| Bearer JWT | Атомарная транзакционная смена адреса клиента | Master only (`shopapp`) |
| `DELETE` | `/api/v1/clients/{id}` | Bearer JWT | Удаление клиента по UUID | Master only (`shopapp`) |

#### Поставщики (Suppliers)

| Метод | Путь | Доступ | Назначение | Балансировка / Кеш |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/suppliers` | Bearer JWT | Список поставщиков с пагинацией | Да (2:1:1), Кеш 1m |
| `GET` | `/api/v1/suppliers/{id}` | Bearer JWT | Получение данных поставщика по UUID | Да (2:1:1), Кеш 1m |
| `POST` | `/api/v1/suppliers` | Bearer JWT | Регистрация нового поставщика | Master only (`shopapp`) |
| `PUT` | `/api/v1/suppliers/{id}/address`| Bearer JWT| Транзакционное обновление адреса поставщика | Master only (`shopapp`) |
| `DELETE` | `/api/v1/suppliers/{id}` | Bearer JWT | Удаление поставщика | Master only (`shopapp`) |

#### Изображения (Images)

| Метод | Путь | Доступ | Назначение | Балансировка / Кеш |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/images/{id}` | Bearer JWT | Скачивание бинарного файла изображения по ID | Да (2:1:1), Кеш 1m |
| `GET` | `/api/v1/products/{id}/image`| Bearer JWT | Скачивание бинарного изображения, привязанного к товару | Да (2:1:1), Кеш 1m |
| `POST` | `/api/v1/products/{id}/image`| Bearer JWT | Загрузка и привязка нового изображения к товару | Master only (`shopapp`) |
| `PUT` | `/api/v1/images/{id}` | Bearer JWT | Замена бинарных данных существующего изображения | Master only (`shopapp`) |
| `DELETE` | `/api/v1/images/{id}` | Bearer JWT | Удаление изображения | Master only (`shopapp`) |

---

### 3. Спецификация gRPC сервиса Auth (`AuthService`)

| RPC Метод | Входное сообщение | Выходное сообщение | Описание |
| :--- | :--- | :--- | :--- |
| `Signup` | `SignupRequest` | `TokensResponse` | Регистрация нового пользователя, хэширование пароля через bcrypt, генерация JWT |
| `Signin` | `SigninRequest` | `TokensResponse` | Проверка учетных данных пользователя и выдача пары Access/Refresh токенов |
| `Refresh` | `RefreshTokensRequest` | `TokensResponse` | Ротация токенов: проверка refresh-токена и выдача новой пары |
| `ValidateToken` | `AccessTokenRequest` | `EmptyMessage` | Быстрая верификация подписи и срока действия access-токена для middleware |
| `ResetPassword` | `ResetPasswordRequest` | `EmptyMessage` | Сброс пароля пользователя с генерацией временного пароля |
| `ChangePassword` | `ChangePasswordRequest` | `EmptyMessage` | Проверка старого пароля и установка нового |
| `User` | `UserRequest` | `UserResponse` | Получение атрибутов профиля пользователя по ID или логину |

---

## Структура проекта и реализованные модули

| Файл / Каталог | Назначение и функционал | Реализованные паттерны и концепции |
| :--- | :--- | :--- |
| `src/docker-compose.yml` | Декларативная оркестрация 8 сервисов стека | Healthchecks, Volumes, зависимости `condition: service_healthy`, изолированные сети |
| `src/Makefile` | Сценарии автоматизации сборки и запуска | Таргеты `all`, `run`, `clean`, `rebuild`, `genproto`, `genswag` |
| `src/.env.template` | Шаблон переменных окружения сервисов | Параметры БД, порты, секретные ключи JWT и время жизни токенов |
| `src/nginx/default.conf` | Конфигурация Nginx: реверс-прокси, SSL, балансировка | `upstream`, `map $request_method`, `proxy_cache`, `gzip`, `stub_status` |
| `src/nginx/generate-cert.sh` | Генерация самоподписанного сертификата X.509 | OpenSSL, RSA 2048, SHA-256, Subject Alternative Name (SAN) |
| `src/nginx/html/` | Каталог статических файлов веб-сервера | `index.html` и статическое изображение `image.png` |
| `src/api/auth.proto` | Контракт Protocol Buffers сервиса авторизации | Описание сообщений и интерфейса `AuthService` для генерации gRPC-кода |
| `src/shop/cmd/main.go` | Точка входа в сервис магазина | Инициализация Uber Fx, подключение gRPC-клиента, маршрутизация HTTP, Swagger |
| `src/shop/internal/domain/` | Доменные модели и контракты интерфейсов | Чистые модели `model.go`, интерфейсы `*_repository.go` без внешних зависимостей |
| `src/shop/internal/web/` | Слой представления и REST API | DTO, хэндлеры `handler/`, мапперы `mapper/`, middleware аутентификации |
| `src/shop/internal/repository/`| Слой работы с базой данных PostgreSQL | Инициализация `pgxpool`, DAO, SQL-мапперы, скрипты `schema.sql` и `readonly.sql` |
| `src/auth/cmd/server/server.go`| Точка входа в gRPC-сервис авторизации | Uber Fx модули `data`, `auth`, `shared`, запуск `grpc.Server` на TCP-порту |
| `src/auth/internal/auth/` | Бизнес-логика авторизации и JWT | Провайдер JWT, кеш токенов, валидация сессий, хэширование паролей |
| `src/auth/internal/data/` | Репозиторий пользователей сервиса Auth | PostgreSQL операции для пользователей, ролей и хэшей |
| `src/test/test.md` | Инструкция и скрипты верификации инфраструктуры | JS-скрипты проверки балансировки (2:1:1) и кеширования (MISS/HIT) |

---

## Инструкция по сборке, настройке и запуску

### Предварительные требования

- **Docker** и **Docker Compose** (рекомендуется Docker Desktop или Docker Engine 24+);
- **OpenSSL** (для генерации TLS-сертификатов);
- Утилита `curl` и современный браузер (для функционального тестирования).

---

### Шаг 1. Настройка локального DNS

Для корректной маршрутизации по доменному имени `be3.shop` добавьте запись в системный файл `/etc/hosts`:

```bash
# Добавление сопоставления be3.shop -> 127.0.0.1
echo "127.0.0.1 be3.shop" | sudo tee -a /etc/hosts
```

---

### Шаг 2. Генерация SSL/TLS-сертификата

Сгенерируйте самоподписанный X.509 сертификат и закрытый ключ для домена `be3.shop`:

```bash
cd BE3_Nginx_Pro_ID_1269323-2-develop/src
./nginx/generate-cert.sh
```

*Скрипт создаст файлы `be3.shop.crt` и `be3.shop.key` в директории `src/nginx/certs` с корректными правами доступа (`chmod 600`).*

---

### Шаг 3. Подготовка переменных окружения

Скопируйте шаблон переменных окружения:

```bash
cp .env.template .env
```

*Файл `.env` содержит базовые настройки для баз данных `shop` и `shop_auth`, учетные данные pgAdmin и параметры генерации токенов JWT.*

---

### Шаг 4. Сборка и запуск контейнеров

Запустите мультиконтейнерный стек с помощью команды `make run` или напрямую через Docker Compose:

```bash
docker compose up -d --build
```

#### Запускаемые контейнеры:
1. `nginx` — прокси-сервер на портах `80` и `443`;
2. `shopapp` — основной Master Go-бэкенд магазина (чтение и запись, порт `8080`);
3. `shopapp_read1` — первая реплика чтения магазина (порт `8081`);
4. `shopapp_read2` — вторая реплика чтения магазина (порт `8082`);
5. `authapp` — микросервис авторизации на gRPC (порт `50051`);
6. `shopdb` — база данных PostgreSQL магазина (порт `5432`);
7. `authdb` — база данных PostgreSQL сервиса авторизации;
8. `pgadmin` — веб-панель управления PostgreSQL.

---

### Шаг 5. Проверка статуса сервисов

Убедитесь, что все контейнеры успешно стартовали и находятся в статусе `Up` (а базы данных — в `healthy`):

```bash
docker compose ps
```

---

## Проверка и тестирование функционала

### 1. Проверка статического контента и HTTPS

Откройте в браузере или запросите через `curl` главную страницу:
```bash
curl -k https://be3.shop/
```
**Результат:** Возвращается HTML-страница `index.html` с тегом `<img src="/image.png">`. Картинка доступна по адресу `https://be3.shop/image.png`.

---

### 2. Мониторинг метрик Nginx (`stub_status`)

Запросите внутреннюю статистику Nginx:
```bash
curl -k https://be3.shop/status
```
**Пример ответа:**
```text
Active connections: 2 
server accepts handled requests
 48 48 112 
Reading: 0 Writing: 1 Waiting: 1 
```

---

### 3. Интерактивная документация Swagger UI

Откройте в браузере адрес:
```text
https://be3.shop/api/v1
```
*Nginx автоматически выполнит 302-редирект на `https://be3.shop/swagger/index.html`, где доступен интерактивный интерфейс OpenAPI со всеми маршрутами, описанием схем запросов и кнопкой авторизации.*

---

### 4. Панель управления pgAdmin

Откройте панель управления СУБД:
```text
https://be3.shop/admin/
```
- **Email:** `admin@example.com` (из `.env`);
- **Пароль:** `pgadmin` (из `.env`).

---

### 5. Тестирование балансировки нагрузки (2:1:1)

Для подтверждения весового распределения GET-запросов между тремя инстансами (`shopapp`, `shopapp_read1`, `shopapp_read2`) в пропорции **2:1:1** выполните скрипт в консоли разработчика браузера (F12 на странице `https://be3.shop`):

```javascript
const counts = {};

for (let i = 0; i < 40; i++) {
  const response = await fetch('/api/v1/products');
  const backend = response.headers.get('X-Upstream-Addr');
  counts[backend] = (counts[backend] || 0) + 1;
}

console.table(counts);
```

**Ожидаемый результат распределения из 40 запросов:**
- Инстанс `shopapp` (Master, вес 2): ~**20 запросов (50%)**;
- Инстанс `shopapp_read1` (Реплика, вес 1): ~**10 запросов (25%)**;
- Инстанс `shopapp_read2` (Реплика, вес 1): ~**10 запросов (25%)**.

*Без авторизационного токена API возвращает `401 Unauthorized` — эти ответы не кешируются, что позволяет наглядно отследить выбор бэкенда для каждого запроса по заголовку `X-Upstream-Addr`.*

---

### 6. Тестирование кеширования Nginx (`proxy_cache`)

Получите токен доступа (см. сценарий cURL ниже) и выполните в консоли браузера два последовательных запроса к одному и тому же URL:

```javascript
{
  const token = prompt('Вставьте access_token');
  const url = `/api/v1/products?demo=${crypto.randomUUID()}`;
  const results = [];

  for (let i = 1; i <= 2; i++) {
    const response = await fetch(url, {
      headers: { Authorization: `Bearer ${token}` },
      cache: 'no-store'
    });
    await response.text();

    results.push({
      Запрос: i,
      HTTP: response.status,
      Кеш: response.headers.get('X-Cache-Status'),
      Бэкенд: response.headers.get('X-Upstream-Addr')
    });
  }

  console.table(results);
}
```

**Результат в консоли:**

| Запрос | HTTP | Кеш (`X-Cache-Status`) | Бэкенд (`X-Upstream-Addr`) | Пояснение |
| :--- | :--- | :--- | :--- | :--- |
| **1** | `200` | **`MISS`** | `172.x.x.x:8080` | Запрос передан на один из бэкендов, ответ записан в кеш Nginx |
| **2** | `200` | **`HIT`** | `null` | Ответ мгновенно отдан из кеша Nginx без обращения к бэкенду |

---

### 7. Тестирование GZIP-сжатия

Проверьте наличие сжатия для текстовых ответов:
```bash
curl -k -I -H "Accept-Encoding: gzip" https://be3.shop/swagger/doc.json
```
**В заголовках ответа присутствует:**
```http
Content-Encoding: gzip
Vary: Accept-Encoding
```

При запросе статического изображения компрессия отсутствует, как и заложено в конфигурации:
```bash
curl -k -I -H "Accept-Encoding: gzip" https://be3.shop/image.png
```
*Заголовок `Content-Encoding: gzip` отсутствует, бинарные данные отдаются в исходном виде.*

---

### 8. Сценарий работы с API через cURL

#### 1. Регистрация нового пользователя
```bash
curl -k -X POST https://be3.shop/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "login": "manager1",
    "password": "Password123!",
    "email": "manager1@shop.com",
    "name": "Ivan",
    "surname": "Ivanov",
    "phonenumber": "+79991234567"
  }'
```

#### 2. Аутентификация и получение JWT-токенов
```bash
curl -k -X POST https://be3.shop/api/v1/auth \
  -H "Content-Type: application/json" \
  -d '{
    "login": "manager1",
    "password": "Password123!"
  }'
```
**Ответ:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 3. Добавление адреса (Master Only)
```bash
curl -k -X POST https://be3.shop/api/v1/addresses \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "country": "Россия",
    "city": "Москва",
    "street": "Тверская, 1"
  }'
```

#### 4. Создание поставщика (Master Only)
```bash
curl -k -X POST https://be3.shop/api/v1/suppliers \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "ООO Главпоставка",
    "address_id": "<ADDRESS_UUID>",
    "phone_number": "+74951234567"
  }'
```

#### 5. Создание нового товара (Master Only)
```bash
curl -k -X POST https://be3.shop/api/v1/products \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ноутбук Pro 16",
    "category": "Электроника",
    "price": 149990.00,
    "available_stock": 25,
    "supplier_id": "<SUPPLIER_UUID>"
  }'
```

#### 6. Получение списка товаров (балансировка 2:1:1 + кеш)
```bash
curl -k -i -X GET https://be3.shop/api/v1/products \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

---

### Остановка и очистка окружения

Для остановки всех контейнеров:
```bash
docker compose down
```

Для полной остановки с удалением томов баз данных:
```bash
docker compose down -v
# или: make clean
```
