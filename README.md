# Портфолио проектов на Go (Golang Portfolio)

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go&logoColor=white)
![Architecture](https://img.shields.io/badge/Architecture-Clean%20Architecture%20%2F%20DDD-blue?style=flat)
![Concurrency](https://img.shields.io/badge/Concurrency-Goroutines%20%26%20Channels-orange?style=flat)
![UI](https://img.shields.io/badge/UI-Fyne%20v2%20(GUI)%20%7C%20Termbox%20(TUI)-green?style=flat)
![Status](https://img.shields.io/badge/Status-Completed-success?style=flat)

Репозиторий включает 5 законченных проектов разного уровня масштаба и сложности, моделирующих реальные задачи системного программирования, backend-разработки и создания пользовательских инструментов:

1. **Алгоритмические настольные приложения (GUI):** генерация идеальных лабиринтов (алгоритм Эллера), поиск кратчайшего маршрута.
2. **Игровая логика и терминальные интерфейсы (TUI):** пошаговая игра в жанре Roguelike с процедурной генерацией подземелий, туманом войны (Field of View), искусственным интеллектом противников и модульной боевой системой на `termbox-go`.
3. **Фундаментальные основы и базовые алгоритмы:** потоковая обработка данных (`bufio`), парсинг строк, математические расчеты, частотный анализ текста, алгоритмические операции над множествами и работа со временем (`time.Time`).
4. **Конкурентное и многопоточное программирование (Concurrency):** организация параллельных вычислений на горутинах, синхронизация через `sync.WaitGroup`, потоковая передача через каналы, паттерны Pipeline и Generator, управление контекстами (`context.Context`) и корректное завершение процессов (Graceful Shutdown).
5. **Серверная архитектура и внедрение зависимостей (DI):** веб-сервис на базе чистой многослойной архитектуры (Clean Architecture / Onion Architecture) с DI-контейнером Uber Fx / Dig, структурированным логированием Zap и изоляцией доменных сущностей от транспортного уровня и хранилища.

---

## Навигация по проектам

| # | Проект | Направление | Основные технологии и концепции | Документация |
| :- | :--- | :--- | :--- | :--- |
| 01 | **A1: Maze** | Desktop GUI Application | Fyne v2, алгоритм Эллера, клеточные автоматы, поиск пути BFS, Clean Architecture | [README.md](./A1_Maze_Go_ID_1391788-Team_TL_maryamer_96897b51_c272_4712-1-develop/README.md) |
| 02 | **AP1: Rogue Game** | Terminal Roguelike Game | Termbox-go, процедурная генерация уровней, FoW, AI противников, инвентарь, сохранения | [README.md](./AP1_Go_P01.ID_1375361-Team_TL_tyananai.bbb98297_c374_40a6-1-develop/README.md) |
| 03 | **AP1: Go Basics** | Algorithms & Data Structures | CLI-утилиты, потоковый ввод `bufio`, сортировка `sort.Slice`, карты `map`, `time.Time`, тесты | [README.md](./AP1_Go_T01.ID_1375359-1-develop/README.md) |
| 04 | **AP1: Concurrency** | Multithreading & Pipelines | Горутины, каналы, `sync.WaitGroup`, конвейер простых чисел, `context.Context`, Graceful Shutdown | [README.md](./AP1_Go_T02.ID_1375360-2-develop/README.md) |
| 05 | **AP1: Clean Architecture** | Web Service & DI | Clean Architecture, DI Uber Fx / Dig, Uber Zap, REST API, In-memory Storage, DTO-мапперы | [README.md](./AP1_Go_T03.ID_1364928-2-develop/README.md) |

---

## Структура репозитория

```text
.
├── README.md                                                        # Общая документация портфолио
├── A1_Maze_Go_ID_1391788.../                                        # Проект 01: Генератор лабиринтов и пещер (Fyne GUI)
│   ├── README.md
│   └── src/ (cmd, domain, application, infrastructure, ui, Makefile)
├── AP1_Go_P01.ID_1375361.../                                        # Проект 02: Терминальный Roguelike (Termbox-go)
│   ├── README.md
│   └── src/ (cmd/game, internal/domain, internal/datalayer, internal/representation)
├── AP1_Go_T01.ID_1375359.../                                        # Проект 03: Базовые алгоритмы и структуры данных
│   ├── README.md
│   └── src/ (01: калькулятор, 02: частотный анализ, 03: пересечения, 04: учет визитов)
├── AP1_Go_T02.ID_1375360.../                                        # Проект 04: Конкурентность и многопоточность
│   ├── README.md
│   └── src/ (ex01: горутины/воркеры, ex02: каналы/pipeline, ex03: context/shutdown)
└── AP1_Go_T03.ID_1364928.../                                        # Проект 05: Веб-сервис «Крестики-нолики» (Clean Architecture & DI)
    ├── README.md
    └── src/ (domain, datasource, web, di/cmd)
```

Каждый проект содержит:
- `README.md` — подробное руководство с описанием архитектуры, компонентов и инструкции по запуску;
- `src/` — исходный код на Go с конфигурационными файлами `go.mod` / `go.sum` и сценариями сборки.

---

## Требования и инструкции по сборке и запуску

### Системные требования
- **Go:** версия 1.22 или выше (рекомендуется Go 1.24+);
- **Системные библиотеки для GUI (проект Maze):** поддержка OpenGL, компилятор C (gcc) и библиотеки разработки для Fyne (на Linux: `libgl1-mesa-dev`, `xorg-dev`).

### Запуск проектов

1. **Клонирование репозитория:**
   ```bash
   git clone <url-репозитория>
   cd <каталог-репозитория>
   ```

2. **Запуск графического приложения Maze (Проект 01):**
   ```bash
   cd A1_Maze_Go_ID_1391788-Team_TL_maryamer_96897b51_c272_4712-1-develop/src
   make install
   ./build/maze
   ```

3. **Запуск терминальной игры Rogue Game (Проект 02):**
   ```bash
   cd AP1_Go_P01.ID_1375361-Team_TL_tyananai.bbb98297_c374_40a6-1-develop/src
   go run ./cmd/game/main.go
   ```

4. **Запуск задач по многопоточности (Проект 04):**
   ```bash
   cd AP1_Go_T02.ID_1375360-2-develop/src
   go run ./ex01/main.go -N 5 -M 100
   go run ./ex02/main.go -K 1 -N 50
   go run ./ex03/main.go -K 2
   ```

5. **Запуск веб-сервиса на Clean Architecture и Uber Fx (Проект 05):**
   ```bash
   cd AP1_Go_T03.ID_1364928-2-develop/src/di/cmd
   go run main.go
   # Сервис доступен в браузере по адресу: http://localhost:8080
   ```
