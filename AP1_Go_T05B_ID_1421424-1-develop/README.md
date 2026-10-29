# Веб-сервис «Крестики-нолики» с JWT‑аутентификацией и PostgreSQL (Tic‑Tac‑Toe REST API, JWT & PostgreSQL)

[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Clean Architecture](https://img.shields.io/badge/Architecture-Clean%20Architecture-blue?style=flat)](#архитектура-и-стек-технологий)
[![JWT Auth](https://img.shields.io/badge/Auth-JWT%20(Access%20%26%20Refresh)-red?style=flat)](#jwt-аутентификация-и-безопасность)
[![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL%20(pgx%20v5)-336791?style=flat&logo=postgresql&logoColor=white)](#архитектура-и-стек-технологий)
[![Uber Fx](https://img.shields.io/badge/DI-Uber%20Fx-black?style=flat)](https://github.com/uber-go/fx)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker&logoColor=white)](#инструкция-по-сборке-и-запуску)
[![Status Completed](https://img.shields.io/badge/Status-Completed-success?style=flat)](#)

Сделал модульный бэкенд для «Крестиков‑ноликов» — надёжный, с чёткими границами ответственности. Реализовал одиночный режим против непобедимого Minimax и мультиплеер между пользователями (PvP). Спроектировал сервис по принципам Clean Architecture: разнёс логику по слоям, подключил зависимости через Uber Fx, интегрировал PostgreSQL на pgx/v5, настроил JWT‑безопасность с ротацией токенов, добавил аудит партий и расчёт таблицы лидеров, а потом упаковал всё в Docker Compose.

---

## О проекте

Это мой серверный веб‑сервис — на нём я потренировал чистую архитектуру и связку сервисов. Вот что реализовал:

- **Управление сессиями и безопасность.** Сделал Bearer JWT с разделением на короткоживущие Access‑токены и долгоживущие Refresh‑токены. Добавил безопасную ротацию токенов и middleware для защиты приватных эндпоинтов.
- **Персистентность и надёжность данных.** Ушёл от in‑memory хранилищ к PostgreSQL. Настроил пул соединений `pgxpool.Pool`, автоматическую инициализацию схемы данных и защиту от состояний гонки.
- **Игровой движок.** Поддержал два режима: одиночный (`vs_computer`) с Minimax для оптимальных ходов и мультиплеер (`vs_human`) с комнатами, списком сессий и присоединением оппонентов.
- **Аналитика и история.** Настроил фиксацию хронологии игр с таймстемпами, выборку архива завершённых партий и расчёт таблицы лидеров по коэффициенту побед (Win Ratio).
- **Инфраструктура и развёртывание.** Подготовил двухконтейнерный стек в Docker Compose (приложение + PostgreSQL) и добавил Healthcheck для проверки готовности БД.

---

## Ключевые навыки и освоенные концепции

### 1. JWT‑аутентификация и безопасность

#### Двухтокенная схема (Access & Refresh Tokens)
- **Access Token.** Выдаю на 15 минут, подписываю HS256. В токене храню `sub` (UUID пользователя), `exp` (срок действия) и `token_type: "access"`.
- **Refresh Token.** Делаю на 7 дней, добавляю `jti` для уникальности. Использую только для продления сессии без повторного логина.

#### Безопасное хранение и ротация токенов (Token Rotation)
- При выдаче Refresh‑токена сохраняю в БД его хэш SHA‑256 (`crypto/sha256`). Даже если базу утечёт, сессии останутся в безопасности.
- При обновлении токенов атомарно заменяю старый хэш на новый (`Replace`) — это аннулирует предыдущий токен и защищает от Replay Attacks.

#### Авторизационный Middleware (`UserAuthenticator`)
- Перехватываю заголовок `Authorization: Bearer <token>`, проверяю формат и подпись токена.
- Извлекаю UUID из claims и безопасно передаю в контекст запроса (`context.WithValue`) для использования хэндлерами.
- Если токена нет или он невалиден — возвращаю `401 Unauthorized` и заголовок `WWW-Authenticate: Bearer`.

### 2. Реляционная СУБД PostgreSQL и `pgx/v5`

- **Пул соединений (`pgxpool.Pool`).** Настроил `MinConns = 3` и `MaxConns = 10` — этого хватает для устойчивости к пиковым нагрузкам.
- **Автоматическая инициализация схемы (DDL).** При старте сервиса создаю таблицы (`users`, `games`, `refresh_tokens`) и внешние ключи (`REFERENCES users(uuid)`), проверяю доступность БД через `Ping` с таймаутом.
- **Нативные типы PostgreSQL.** Хранил игровое поле 3×3 как одномерный массив `TEXT[]` и настроил трансляцию между 2D‑матрицей и 1D‑массивом на уровне мапперов.
- **Параметризованные запросы.** Во всех SQL‑запросах использую плейсхолдеры `$1, $2, ...` — это надёжно защищает от SQL‑инъекций.

### 3. Таблица лидеров и статистика (Leaderboard Analytics)

#### Расчёт коэффициента побед (Win Ratio)
- Анализирую все завершённые игры пользователя с фазами `win` и `draw`.
- Считаю по формуле:
  $$\text{Win Ratio} = \frac{\text{Wins}}{\text{Losses} + \text{Draws}}$$
- Обрабатываю граничные случаи: если знаменатель равен нулю, приравниваю его к 1 — так избегаю деления на ноль.

#### Многофакторная сортировка и пагинация
- Сначала сортирую игроков по убыванию `WinRatio`.
- Если рейтинги совпадают — сортирую по возрастанию `UserUUID`, чтобы порядок был детерминированным.
- Ограничиваю выборку топ‑N игроков и обогащаю данные логинами через `UserService`.

### 4. Игровой движок, правила и алгоритм Minimax

#### Управление игровыми фазами (`GameState`)
- `ready` — жду подключения второго игрока (в `vs_human`).
- `move` — партия активна, сейчас ход конкретного игрока (`PlayerUUID`).
- `win` — игра завершена победой участника.
- `draw` — ничья: все 9 клеток заполнены, победителя нет.

#### Строгая валидация ходов
- Проверяю, что ходит именно активный игрок (совпадение UUID).
- Слежу, чтобы менялся ровно один элемент доски.
- Не даю перезаписывать уже занятые клетки.
- Контролирую символы: `X` — для создателя игры, `O` — для второго игрока или компьютера.

#### ИИ на базе Minimax
- Построил рекурсивное дерево исходов и добавил эвристическую оценку: +1 за выигрыш компьютера, −1 за победу человека, 0 за ничью. Так ИИ всегда отвечает оптимальным ходом.

### 5. Архитектура Clean Architecture и DI (Uber Fx)

#### Слоистая структура
- `domain`: чистые сущности, модели и интерфейсы репозиториев/сервисов без внешних зависимостей.
- `datasource`: инфраструктура БД, репозитории, JWT‑провайдер и реализация бизнес‑логики.
- `web`: контроллеры REST API, JSON‑сериализация, роутинг и middleware аутентификации.
- `di`: точка входа, сборка графа зависимостей через `fx.Provide` и запуск HTTP‑сервера через `fx.Invoke`.

#### Изоляция и инкапсуляция
- Логику спрятал в пакетах `internal`. Снаружи оставил только фабричные конструкторы в `provide.go`.

---

## Архитектура и стек технологий

### Архитектурные слои

```text
       ┌────────────────────────────────────────────────────────┐
       │                   DI Layer (di/cmd)                    │
       │     Компоновка Uber Fx, конфигурация, запуск сервера   │
       └──────────────┬──────────────────────────┬──────────────┘
                      │                          │
       ┌──────────────▼──────────┐    ┌──────────▼──────────────┐
       │     Web Layer (web)     │    │ Datasource (datasource) │
       │  REST API, Middleware,  │    │  PostgreSQL (pgx v5),   │
       │      DTO, Мапперы       │    │ JWT Provider, Репозитории│
       └──────────────┬──────────┘    └──────────┬──────────────┘
                      │                          │
                      └───────────► ◄────────────┘
                                   │
                      ┌─────────────▼────────────┐
                      │   Domain Layer (domain)  │
                      │  Бизнес-модели, правила, │
                      │       интерфейсы         │
                      └──────────────────────────┘
```

### Стек технологий

| Технология / Библиотека | Версия | Назначение в проекте |
| :--- | :--- | :--- |
| **Go** | `1.25.6` | Базовый язык и стандартная библиотека (`net/http`, `context`, `crypto/sha256`) |
| **PostgreSQL** | `18.1-Alpine` | Основная реляционная СУБД |
| **pgx** | `v5.11.0` | Драйвер PostgreSQL и пул соединений `pgxpool` |
| **golang-jwt/jwt** | `v5.3.1` | Генерация, подпись HS256 и валидация JWT‑токенов |
| **Uber Fx** | `v1.24.0` | DI‑контейнер и управление жизненным циклом |
| **Uber Dig** | `v1.19.0` | Ядро графа зависимостей для Uber Fx |
| **google/uuid** | `v1.6.0` | Генерация криптографически стойких UUID |
| **Docker & Docker Compose** | — | Контейнеризация и оркестрация стека БД + приложение |

---

## Спецификация REST API

### Аутентификация и пользователи

| Метод | Путь | Доступ | Назначение |
| :--- | :--- | :--- | :--- |
| `POST` | `/signup` | Публичный | Регистрация нового пользователя по логину и паролю |
| `POST` | `/login` | Публичный | Аутентификация, получение пары Access & Refresh токенов |
| `POST` | `/auth/access` | Публичный | Обновление пары токенов по действующему Refresh‑токену |
| `POST` | `/auth/refresh` | Публичный | Ротация Refresh‑токена с выдачей новых Access и Refresh токенов |
| `GET` | `/users/me` | Bearer JWT | Получение UUID и логина текущего авторизованного пользователя |
| `GET` | `/users/{id}` | Bearer JWT | Получение публичной информации о пользователе по UUID |

### Управление играми и таблица лидеров

| Метод | Путь | Доступ | Назначение |
| :--- | :--- | :--- | :--- |
| `POST` | `/game` | Bearer JWT | Создание новой игры (`vs_computer` или `vs_human`) |
| `GET` | `/games/available` | Bearer JWT | Список доступных игр со статусом `ready` |
| `GET` | `/games/completed` | Bearer JWT | Архив завершённых игр текущего пользователя (`win`, `draw`) |
| `GET` | `/games/leaderboard?n=10` | Bearer JWT | Таблица лидеров: топ‑N игроков по Win Ratio |
| `GET` | `/game/{id}` | Bearer JWT | Получение актуального состояния игры по UUID |
| `POST` | `/game/{id}` | Bearer JWT | Совершение хода игроком с обновлением доски |
| `POST` | `/game/{id}/join` | Bearer JWT | Присоединение второго игрока к открытой игре (`vs_human`) |

---

## Структура проекта и реализованные модули

| Файл / Каталог | Слой | Описание и функционал | Реализованные Go‑паттерны и конструкции |
| :--- | :--- | :--- | :--- |
| `domain/model.go` | Domain | Доменные сущности: `Game`, `Board`, `User`, `GameState`, DTO авторизации | Пользовательские типы, константы фаз (`GamePhase`), модели данных |
| `domain/repository.go` | Domain | Контракты персистентности данных | Интерфейсы `GameRepository`, `UserRepository`, `RefreshTokenRepository` |
| `domain/service.go` | Domain | Контракты доменных служб и провайдеров | Интерфейсы `GameService`, `UserService`, `AuthService`, `JwtProvider` |
| `datasource/provide.go` | Datasource | Фасад слоя данных для DI‑контейнера | Фабричные функции `NewGameRepository`, `NewPostgresPool`, `NewAuthService` и др. |
| `datasource/internal/db.go` | Datasource | Инициализация пула соединений PostgreSQL | `pgxpool.ParseConfig`, `pgxpool.NewWithConfig`, DDL‑миграции `CREATE TABLE` |
| `datasource/internal/jwt_provider.go` | Datasource | Реализация провайдера JWT | Подпись и парсинг токенов HS256 (`jwt.NewWithClaims`, `jwt.NewParser`) |
| `datasource/internal/auth_service.go` | Datasource | Бизнес‑логика аутентификации и сессий | Регистрация, вход, валидация и ротация refresh‑токенов |
| `datasource/internal/refresh_token_repository.go` | Datasource | Персистентность refresh‑токенов | Хэширование SHA‑256 (`sha256.Sum256`), атомарная замена `Replace` |
| `datasource/internal/user_service.go` | Datasource | Сервис работы с учётными записями | Реализация интерфейса `domain.UserService` |
| `datasource/internal/game_service.go` | Datasource | Игровые правила и расчёт ходов ИИ | Алгоритм Minimax, валидация ходов человека, проверка линий победы |
| `datasource/internal/repository.go` | Datasource | Реализации репозиториев на базе PostgreSQL | Запросы через `pgxpool.Pool`, сканирование строк, расчёт Win Ratio |
| `datasource/internal/model.go` | Datasource | Персистентные структуры таблиц БД | Структуры `GameEntity`, `UserEntity`, срез `BoardEntity [9]string` |
| `datasource/internal/mapper.go` | Datasource | Трансляция данных между БД и доменом | Конвертация 2D‑матрицы `[3][3]` в 1D‑массив `[9]` и обратно |
| `web/provide.go` | Web | Публичные фасады веб‑контроллеров для DI | Фабрики `NewGameHandler`, `NewAuthHandler`, `NewUserAuthenticator` |
| `web/internal/user_authenticator.go` | Web | Авторизационный HTTP Middleware | Проверка Bearer токена, валидация claims, проброс UUID через Context |
| `web/internal/auth_handler.go` | Web | Контроллер эндпоинтов авторизации | Обработка `/signup`, `/login`, `/auth/access`, `/auth/refresh` |
| `web/internal/user_handler.go` | Web | Контроллер профилей пользователей | Обработка `/users/me` и `/users/{id}` |
| `web/internal/game_handler.go` | Web | Контроллер игровых сценариев | Обработка `/game`, `/games/available`, `/games/completed`, `/games/leaderboard` |
| `web/internal/model.go` | Web | Сетевые DTO запросов и ответов | Структуры с `json`-тегами: `GameResponse`, `TopPlayerResponse` и др. |
| `web/internal/mapper.go` | Web | Маппинг между DTO и доменными типами | Чистые функции `toWeb` и `toDomain` |
| `di/cmd/main.go` | DI | Точка входа в приложение | Инициализация `fx.New`, сборка зависимостей, запуск `http.ListenAndServe` |
| `docker-compose.yaml` | DevOps | Декларативная оркестрация сервисов | Конфигурация PostgreSQL 18.1, Go‑сервиса, Healthcheck и Volumes |

---

## Инструкция по сборке и запуску

### Вариант 1. Запуск через Docker Compose (рекомендуемый)

Я собрал всё в Docker Compose, чтобы не ставить локально Go и PostgreSQL. Достаточно Docker с Compose.

1. Перейдите в каталог `src`:
   ```bash
   cd src
   ```

2. Запустите сервисы в фоне:
   ```bash
   docker compose up --build -d
   ```

3. Проверьте статус контейнеров:
   ```bash
   docker compose ps
   ```

   Контейнер БД `tictactoe_db` станет `healthy`, после чего запустится `tictactoe_app`. Сервис будет доступен по адресу: `http://localhost:8000`.

4. Для остановки контейнеров выполните:
   ```bash
   docker compose down
   ```

---

### Вариант 2. Локальный запуск

#### Требования
- Установленный Go версии 1.25 или выше.
- Работающий экземпляр PostgreSQL (порт по умолчанию: `5432`).

#### Шаги запуска

1. Задайте строку подключения к базе данных через переменную окружения (если параметры отличаются от значений по умолчанию):
   ```bash
   export DATABASE_URL="postgres://postgres:postgres@localhost:5432/tictactoe?sslmode=disable"
   export JWT_SECRET="your-secure-secret-key"
   ```

2. Перейдите в каталог точки входа и запустите приложение:
   ```bash
   cd src/di/cmd
   go run main.go
   ```

При первом запуске сервис автоматически создаст необходимые таблицы в указанной базе данных PostgreSQL.

---

## Примеры использования API через cURL

### 1. Регистрация нового пользователя
```bash
curl -X POST http://localhost:8000/signup \
  -H "Content-Type: application/json" \
  -d '{"login": "player1", "password": "securepassword"}'
```
**Ответ:**
```json
{"uuid": "3fa85f64-5717-4562-b3fc-2c963f66afa6"}
```

### 2. Вход в систему и получение токенов
```bash
curl -X POST http://localhost:8000/login \
  -H "Content-Type: application/json" \
  -d '{"login": "player1", "password": "securepassword"}'
```
**Ответ:**
```json
{
  "type": "Bearer",
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 3. Получение информации о текущем пользователе
```bash
curl -X GET http://localhost:8000/users/me \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

### 4. Создание новой игры против компьютера
```bash
curl -X POST http://localhost:8000/game \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"game_mode": "vs_computer"}'
```
**Ответ:**
```json
{
  "uuid": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "board": [
    ["", "", ""],
    ["", "", ""],
    ["", "", ""]
  ],
  "game_over": false,
  "phase": "move",
  "player_uuid": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "player_x": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "player_o": "00000000-0000-0000-0000-000000000000",
  "created_at": "2026-10-29T10:00:00Z"
}
```

### 5. Совершение хода игроком
```bash
curl -X POST http://localhost:8000/game/a1b2c3d4-e5f6-7890-abcd-ef1234567890 \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "board": [
      ["X", "", ""],
      ["", "", ""],
      ["", "", ""]
    ]
  }'
```
*Я настроил сервер так, что он валидирует ход игрока, а затем ИИ автоматически рассчитывает ответный ход по алгоритму Minimax и возвращает обновлённое состояние доски.*

### 6. Запрос таблицы лидеров (Top-5)
```bash
curl -X GET "http://localhost:8000/games/leaderboard?n=5" \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```
**Ответ:**
```json
[
  {
    "uuid": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "login": "player1",
    "win_ratio": 0.75
  }
]
```

### 7. Запрос завершённых игр текущего пользователя
```bash
curl -X GET http://localhost:8000/games/completed \
  -H "Authorization: Bearer <ACCESS_TOKEN}"
```

### 8. Обновление токенов доступа (Token Refresh)
```bash
curl -X POST http://localhost:8000/auth/access \
  -H "Content-Type: application/json" \
  -d '{"refreshToken": "<REFRESH_TOKEN>"}'
```
