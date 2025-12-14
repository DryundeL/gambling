# Domain-Driven Design (DDD) Архитектура

## 📋 Содержание

1. [Что такое DDD?](#что-такое-ddd)
2. [Структура проекта](#структура-проекта)
3. [Слои архитектуры](#слои-архитектуры)
4. [Схема архитектуры](#схема-архитектуры)
5. [Детальное описание компонентов](#детальное-описание-компонентов)
6. [Принципы и паттерны](#принципы-и-паттерны)
7. [Поток данных](#поток-данных)

---

## Что такое DDD?

**Domain-Driven Design (DDD)** — это подход к разработке программного обеспечения, который фокусируется на моделировании бизнес-логики в соответствии с доменной областью (предметной областью) приложения.

### Основные принципы DDD:

1. **Фокус на домене** — бизнес-логика находится в центре архитектуры
2. **Разделение ответственности** — четкое разделение на слои
3. **Независимость от инфраструктуры** — доменная логика не зависит от технических деталей
4. **Явное моделирование** — бизнес-правила выражены явно в коде

---

## Структура проекта

```
internal/
├── domain/                    # ДОМЕННЫЙ СЛОЙ (ядро бизнес-логики)
│   ├── user/
│   │   ├── entity.go          # Сущность User
│   │   ├── errors.go          # Доменные ошибки
│   │   ├── repository.go      # Интерфейс репозитория (порт)
│   │   └── value_object.go    # Value Objects
│   ├── transaction/
│   │   ├── entity.go          # Сущность Transaction
│   │   └── repository.go      # Интерфейс репозитория
│   └── spin/
│       ├── entity.go          # Сущность SpinResult
│       ├── repository.go      # Интерфейс репозитория
│       └── service.go         # Доменный сервис для логики игры
│
├── application/               # APPLICATION СЛОЙ (use cases)
│   └── use_case/
│       ├── auth/
│       │   ├── register.go    # Use case регистрации
│       │   └── login.go       # Use case входа
│       ├── balance/
│       │   └── deposit.go     # Use case пополнения баланса
│       └── spin/
│           └── spin.go        # Use case выполнения спина
│
├── infrastructure/            # INFRASTRUCTURE СЛОЙ (технические детали)
│   ├── repository/
│   │   ├── user_repository.go      # Реализация репозитория User
│   │   ├── transaction_repository.go
│   │   └── spin_repository.go
│   └── database/
│       └── pgsql/
│           ├── pgsql.go       # Подключение к БД
│           └── migrations.go  # Миграции
│
└── interfaces/                # INTERFACES СЛОЙ (внешние интерфейсы)
    └── http/
        ├── handlers/          # HTTP handlers
        ├── router/            # Маршрутизация
        └── middleware/        # Middleware
```

---

## Слои архитектуры

### 🎯 1. Domain Layer (Доменный слой)

**Назначение:** Содержит чистую бизнес-логику без зависимостей от внешних технологий.

#### Компоненты:

- **Entities (Сущности)** — объекты с уникальной идентичностью
  - `User` — пользователь системы
  - `Transaction` — транзакция пользователя
  - `SpinResult` — результат игры

- **Value Objects** — объекты без идентичности, определяемые только значениями
  - `Credentials` — учетные данные пользователя

- **Domain Services** — сервисы для логики, не принадлежащей конкретной сущности
  - `SpinService` — логика генерации символов и расчета выигрыша

- **Repository Interfaces (Ports)** — интерфейсы для работы с данными
  - `user.Repository`
  - `transaction.Repository`
  - `spin.Repository`

**Принцип:** Доменный слой НЕ зависит от других слоев!

---

### 📱 2. Application Layer (Слой приложения)

**Назначение:** Координирует выполнение бизнес-операций (use cases).

#### Компоненты:

- **Use Cases** — конкретные бизнес-операции
  - `RegisterUseCase` — регистрация пользователя
  - `LoginUseCase` — вход пользователя
  - `DepositUseCase` — пополнение баланса
  - `SpinUseCase` — выполнение спина

- **Commands** — входные данные для use cases
  - `RegisterCommand`, `LoginCommand`, `DepositCommand`, `SpinCommand`

- **Results** — выходные данные use cases
  - `RegisterResult`, `LoginResult`, `DepositResult`, `SpinResult`

**Принцип:** Application слой зависит только от Domain слоя!

---

### 🔧 3. Infrastructure Layer (Слой инфраструктуры)

**Назначение:** Реализует технические детали (БД, внешние сервисы, файловая система).

#### Компоненты:

- **Repository Implementations (Adapters)** — реализации интерфейсов из Domain
  - `UserRepository` — реализация через GORM
  - `TransactionRepository`
  - `SpinRepository`

- **Database** — подключение к БД и миграции
  - `pgsql.Storage` — подключение к PostgreSQL

**Принцип:** Infrastructure реализует интерфейсы из Domain!

---

### 🌐 4. Interfaces Layer (Слой интерфейсов)

**Назначение:** Предоставляет точки входа в приложение (HTTP, CLI, gRPC и т.д.).

#### Компоненты:

- **HTTP Handlers** — обработчики HTTP запросов
  - `AuthHandler` — обработка запросов аутентификации
  - `BalanceHandler` — обработка запросов баланса
  - `SpinHandler` — обработка запросов игры

- **Router** — маршрутизация HTTP запросов

- **Middleware** — промежуточное ПО (логирование, аутентификация и т.д.)

**Принцип:** Interfaces преобразует внешние запросы в команды Application слоя!

---

## Схема архитектуры

```
┌─────────────────────────────────────────────────────────────┐
│                    INTERFACES LAYER                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ AuthHandler  │  │BalanceHandler│  │ SpinHandler  │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                 │                 │               │
│         └─────────────────┴─────────────────┘               │
│                            │                                 │
│                    ┌───────▼────────┐                        │
│                    │     Router     │                        │
│                    └───────┬────────┘                        │
└────────────────────────────┼─────────────────────────────────┘
                             │
                             │ Commands/Results
                             │
┌────────────────────────────▼─────────────────────────────────┐
│                  APPLICATION LAYER                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │RegisterUseCase│ │LoginUseCase  │ │DepositUseCase │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                 │                 │               │
│  ┌──────▼─────────────────┴─────────────────▼───────┐      │
│  │              SpinUseCase                          │      │
│  └──────────────────────┬────────────────────────────┘      │
└─────────────────────────┼────────────────────────────────────┘
                           │
                           │ Uses Domain Services & Entities
                           │
┌──────────────────────────▼────────────────────────────────────┐
│                    DOMAIN LAYER                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │    User      │  │ Transaction  │  │ SpinResult    │      │
│  │   Entity     │  │   Entity     │  │   Entity     │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                 │                 │               │
│  ┌──────▼─────────────────┴─────────────────▼───────┐      │
│  │            SpinService (Domain Service)           │      │
│  └───────────────────────────────────────────────────┘      │
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │user.Repository│ │transaction.  │ │spin.Repository│      │
│  │  (Interface)  │  │Repository     │  │  (Interface)  │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
└─────────┼─────────────────┼─────────────────┼───────────────┘
          │                 │                 │
          │                 │                 │
          │  Implements     │  Implements     │  Implements
          │                 │                 │
┌─────────▼─────────────────▼─────────────────▼───────────────┐
│              INFRASTRUCTURE LAYER                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │UserRepository│ │TransactionRepo│ │SpinRepository │      │
│  │(GORM Adapter)│ │(GORM Adapter) │ │(GORM Adapter) │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                 │                 │               │
│         └─────────────────┴─────────────────┘               │
│                            │                                 │
│                    ┌───────▼────────┐                        │
│                    │  PostgreSQL    │                        │
│                    │   Database     │                        │
│                    └────────────────┘                        │
└───────────────────────────────────────────────────────────────┘
```

---

## Детальное описание компонентов

### Domain Layer

#### Entity: User

```go
// internal/domain/user/entity.go
type User struct {
    ID           uint
    Username     string
    Email        string
    PasswordHash string
    Balance      float64
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

// Доменные методы инкапсулируют бизнес-логику
func (u *User) Deposit(amount float64) error
func (u *User) Withdraw(amount float64) error
func (u *User) AddWin(amount float64) error
```

**Почему это Entity?**
- Имеет уникальную идентичность (ID)
- Может изменяться со временем
- Содержит бизнес-логику (методы Deposit, Withdraw)

#### Domain Service: SpinService

```go
// internal/domain/spin/service.go
type Service struct {
    rng *rand.Rand
}

func (s *Service) GenerateSymbol() int
func (s *Service) CalculateWin(reel1, reel2, reel3 int, betAmount float64) float64
```

**Почему это Domain Service?**
- Логика не принадлежит конкретной сущности
- Используется несколькими use cases
- Содержит сложную бизнес-логику (вероятности, расчеты)

#### Repository Interface (Port)

```go
// internal/domain/user/repository.go
type Repository interface {
    Create(user *User) error
    GetByID(id uint) (*User, error)
    GetByUsername(username string) (*User, error)
    // ...
}
```

**Почему это Port?**
- Определяет контракт для работы с данными
- Не зависит от реализации (GORM, MongoDB, файлы и т.д.)
- Позволяет легко менять инфраструктуру

---

### Application Layer

#### Use Case: RegisterUseCase

```go
// internal/application/use_case/auth/register.go
type RegisterUseCase struct {
    userRepo user.Repository
}

func (uc *RegisterUseCase) Execute(cmd RegisterCommand) (*RegisterResult, error) {
    // 1. Валидация бизнес-правил
    // 2. Использование доменных сущностей
    // 3. Координация операций
    // 4. Сохранение через репозиторий
}
```

**Почему это Use Case?**
- Представляет одну бизнес-операцию
- Координирует работу доменных объектов
- Не содержит технических деталей

---

### Infrastructure Layer

#### Repository Implementation (Adapter)

```go
// internal/infrastructure/repository/user_repository.go
type UserRepository struct {
    db *gorm.DB  // Техническая деталь
}

func (r *UserRepository) Create(u *user.User) error {
    // Преобразование доменной сущности в модель БД
    dbUser := toDBModel(u)
    return r.db.Create(dbUser).Error
}
```

**Почему это Adapter?**
- Адаптирует доменный интерфейс к GORM
- Скрывает технические детали БД от домена
- Можно заменить на другую реализацию без изменения домена

---

### Interfaces Layer

#### HTTP Handler

```go
// internal/interfaces/http/handlers/auth_handler.go
type AuthHandler struct {
    registerUseCase *auth.RegisterUseCase
    loginUseCase    *auth.LoginUseCase
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    // 1. Парсинг HTTP запроса
    // 2. Преобразование в Command
    // 3. Вызов Use Case
    // 4. Преобразование Result в HTTP ответ
}
```

**Почему это Adapter?**
- Адаптирует HTTP протокол к Application слою
- Обрабатывает HTTP-специфичные детали (коды статусов, заголовки)
- Можно заменить на gRPC, GraphQL без изменения Application слоя

---

## Принципы и паттерны

### 1. Dependency Inversion Principle (DIP)

**Принцип:** Модули высокого уровня не должны зависеть от модулей низкого уровня. Оба должны зависеть от абстракций.

**В нашем проекте:**
- Domain определяет интерфейсы (абстракции)
- Infrastructure реализует эти интерфейсы
- Application использует интерфейсы, а не конкретные реализации

### 2. Ports & Adapters (Hexagonal Architecture)

**Порты (Ports):** Интерфейсы в Domain слое
- `user.Repository`
- `transaction.Repository`
- `spin.Repository`

**Адаптеры (Adapters):** Реализации в Infrastructure и Interfaces слоях
- `infrastructure/repository/UserRepository` (GORM)
- `interfaces/http/handlers/AuthHandler` (HTTP)

### 3. Separation of Concerns

Каждый слой имеет четкую ответственность:
- **Domain** — бизнес-логика
- **Application** — координация операций
- **Infrastructure** — технические детали
- **Interfaces** — внешние интерфейсы

### 4. Single Responsibility Principle (SRP)

Каждый компонент отвечает за одну вещь:
- `User` — управление балансом пользователя
- `SpinService` — логика игры
- `RegisterUseCase` — регистрация пользователя

---

## Поток данных

### Пример: Регистрация пользователя

```
1. HTTP Request
   POST /api/v1/register
   {
     "username": "john",
     "email": "john@example.com",
     "password": "secret123"
   }
   │
   ▼
2. AuthHandler.Register()
   - Парсит HTTP запрос
   - Создает RegisterCommand
   │
   ▼
3. RegisterUseCase.Execute()
   - Проверяет существование пользователя
   - Хеширует пароль
   - Создает доменную сущность User
   │
   ▼
4. user.Repository.Create()
   - Вызывается через интерфейс
   │
   ▼
5. infrastructure/repository/UserRepository.Create()
   - Преобразует User в DBUser
   - Сохраняет в PostgreSQL через GORM
   │
   ▼
6. PostgreSQL Database
   - Сохраняет данные
   │
   ▼
7. Возврат результата (обратный путь)
   - UserRepository → RegisterUseCase → AuthHandler → HTTP Response
```

### Пример: Выполнение спина

```
1. HTTP Request
   POST /api/v1/spin?user_id=1
   {
     "bet_amount": 10.0
   }
   │
   ▼
2. SpinHandler.Spin()
   - Парсит запрос
   - Создает SpinCommand
   │
   ▼
3. SpinUseCase.Execute()
   - Получает User через user.Repository
   - Вызывает user.Withdraw() (доменная логика)
   - Создает Transaction через transaction.Repository
   - Использует spin.Service для генерации символов
   - Использует spin.Service для расчета выигрыша
   - Если есть выигрыш: вызывает user.AddWin()
   - Сохраняет SpinResult через spin.Repository
   │
   ▼
4. Domain Services & Entities
   - User.Withdraw() проверяет баланс
   - SpinService.GenerateSymbol() генерирует символы
   - SpinService.CalculateWin() вычисляет выигрыш
   │
   ▼
5. Infrastructure Repositories
   - Сохраняют данные в PostgreSQL
   │
   ▼
6. HTTP Response
   {
     "reel1": 7,
     "reel2": 7,
     "reel3": 8,
     "is_win": false,
     "win_amount": 0,
     "balance": 90.0
   }
```

---

## Преимущества DDD архитектуры

### ✅ 1. Тестируемость
- Доменная логика легко тестируется без БД
- Use cases можно тестировать с моками репозиториев

### ✅ 2. Независимость от инфраструктуры
- Можно заменить PostgreSQL на MongoDB без изменения домена
- Можно добавить gRPC API без изменения Application слоя

### ✅ 3. Явная бизнес-логика
- Бизнес-правила видны в коде
- Легко понять, что делает приложение

### ✅ 4. Масштабируемость
- Легко добавлять новые use cases
- Легко добавлять новые интерфейсы (CLI, gRPC)

### ✅ 5. Поддерживаемость
- Четкое разделение ответственности
- Изменения в одном слое не влияют на другие

---

## Заключение

DDD архитектура обеспечивает:
- **Чистую бизнес-логику** в центре приложения
- **Гибкость** в выборе технологий
- **Тестируемость** всех компонентов
- **Масштабируемость** для роста приложения

Эта архитектура особенно полезна для сложных бизнес-приложений, где бизнес-логика является ключевым активом.

