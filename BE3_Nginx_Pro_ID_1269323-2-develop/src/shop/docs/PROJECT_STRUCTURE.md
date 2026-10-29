# Структура проекта

```text
src/shop/
├── go.mod                          # Модуль shop, версия Go и зависимости
├── go.sum                          # Контрольные суммы зависимостей
├── cmd/
│   └── main.go                     # Точка входа: настройка зависимостей через Fx, gRPC-клиента AUTH и маршрутов HTTP-сервера
├── api/
│   ├── auth.pb.go                  # Сгенерированные protobuf-сообщения AUTH
│   └── auth_grpc.pb.go             # Сгенерированный gRPC-клиент и серверный контракт AUTH
├── docs/
│   ├── PROJECT_STRUCTURE.md        # Дерево проекта и назначение файлов
│   ├── openapi/                    # Сгенерированные Swagger-файлы: docs.go, swagger.json, swagger.yaml
│   └── api/
│       ├── address.md              # API адресов: список с пагинацией, независимое создание, замена общих полей и удаление без ссылок
│       ├── auth.md                 # Регистрация, вход и восстановление пароля через AUTH: поля, примеры запросов, ответы и ошибки
│       ├── client.md               # API клиентов: Bearer-аутентификация, операции, примеры запросов и ошибки
│       ├── image.md                # API изображений: Bearer-аутентификация, добавление, изменение, удаление и скачивание по ID изображения или товара
│       ├── product.md              # API товаров: Bearer-аутентификация, добавление, уменьшение остатка, получение и удаление
│       └── supplier.md             # API поставщиков: Bearer-аутентификация, добавление, удаление, смена адреса, получение по ID и список
└── internal/                       # Внутренние пакеты приложения
    ├── domain/                     # Модели, интерфейсы и правила магазина
    │   ├── model.go                # Client, Product, Supplier, Image, Address
    │   ├── address_repository.go   # Интерфейс списка с пагинацией, создания, изменения и удаления адресов; различение отсутствия и использования
    │   ├── client_repository.go    # Интерфейс операций с клиентами, включая смену адреса
    │   ├── image_repository.go     # Интерфейс операций с изображениями, включая получение по ID изображения и товара
    │   ├── product_repository.go   # Интерфейс операций с товарами, включая получение по ID и список с пагинацией
    │   └── supplier_repository.go  # Интерфейс операций с поставщиками, включая получение по ID и список с пагинацией
    ├── web/                        # HTTP-слой
    │   ├── dto/                    # Входные и выходные структуры API, валидация
    │   │   ├── address.go          # DTO создания адреса (в том числе при смене адреса владельца), замены полей и пагинации, ответ и валидация
    │   │   ├── auth.go             # SignupRequest/Response, SigninRequest/Response и ResetPasswordRequest; проверка логина и пароля
    │   │   ├── client.go           # DTO добавления, поиска и пагинации клиентов, валидация и ответ
    │   │   ├── error.go            # JSON-ответ с ошибкой
    │   │   ├── image.go            # DTO добавления и изменения, валидация, ответ с ID и бинарный ответ файла
    │   │   ├── product.go          # DTO добавления, уменьшения остатка и пагинации товаров, валидация и ответ
    │   │   └── supplier.go         # DTO добавления и пагинации поставщиков, валидация и ответ
    │   ├── middleware/             # Обёртки HTTP-обработчиков: аутентификация и управление кешированием
    │   │   ├── cache.go            # NoStore: Cache-Control: no-store для ответов; сейчас не подключён в main.go
    │   │   └── user_authenticator.go # UserAuthenticator: Bearer-заголовок, проверка токена в AUTH, JSON-ошибки 401/500
    │   ├── mapper/                 # Преобразование DTO в модели
    │   │   ├── address.go          # Создание адреса с серверным UUID, замена полей с сохранением UUID и преобразование ответа
    │   │   ├── client.go           # ClientToDomain и ClientToWeb: DTO и domain.Client
    │   │   ├── image.go            # ImageToDomain, UpdateImageToDomain, ImageToWeb и ImageFileToWeb
    │   │   ├── product.go          # ProductToDomain и ProductToWeb: DTO и domain.Product
    │   │   └── supplier.go         # SupplierToDomain и SupplierToWeb: DTO и domain.Supplier
    │   └── handler/                # HTTP-обработчики
    │       ├── address.go          # Список с пагинацией, создание, изменение и удаление адресов; ответы 404 и 409 при удалении
    │       ├── auth.go             # Signup, Signin и ResetPassword: вызовы AUTH по gRPC, ответы с ID, токенами или без тела
    │       ├── client.go           # HTTP-обработчики клиентов, включая смену адреса и общий GET
    │       ├── image.go            # Обработчики изображений; оба GET используют общую функцию writeImageFile
    │       ├── product.go          # HTTP-обработчики товаров, включая GetByID и List
    │       ├── request.go          # Общий разбор JSON-тела с проверкой отсутствия лишних значений и мусора
    │       ├── response.go         # Отправка ошибок в JSON
    │       └── supplier.go         # HTTP-обработчики поставщиков, включая GetByID и List
    └── repository/                 # PostgreSQL, запросы, DAO и мапперы
        ├── address.go              # Список с пагинацией, создание, изменение и условное удаление адресов; общая транзакция смены адреса сущности
        ├── client.go               # Операции с клиентами в PostgreSQL; смена адреса в транзакции
        ├── image.go                # Операции с изображениями; GetByProductID использует JOIN и различает отсутствие товара и изображения
        ├── product.go              # Операции с товарами в PostgreSQL; список доступных товаров с пагинацией
        ├── supplier.go             # Операции с поставщиками в PostgreSQL; получение по ID и список с пагинацией
        ├── dao/
        │   ├── address.go          # Структура записи таблицы address
        │   ├── client.go           # Структура записи таблицы client
        │   ├── image.go            # Структура записи таблицы images
        │   ├── product.go          # Структура записи таблицы product
        │   └── supplier.go         # Структура записи таблицы supplier
        ├── mapper/
        │   ├── address.go          # AddressToDAO и AddressToDomain
        │   ├── client.go           # ClientToDAO и ClientToDomain
        │   ├── image.go            # ImageToDAO и ImageToDomain: доменная модель и запись БД
        │   ├── product.go          # ProductToDAO и ProductToDomain
        │   └── supplier.go         # SupplierToDAO и SupplierToDomain
        ├── db.go                   # Настройка пула и проверка подключения к PostgreSQL
        └── schema.sql              # Таблицы БД; Docker выполняет при первой инициализации
```
