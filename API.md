# API Документация - Система консольного казино

## Описание

Система консольного казино с регистрацией пользователей, пополнением баланса и игрой на спинах.

## Эндпоинты

### 1. Регистрация пользователя

**POST** `/api/v1/register`

**Тело запроса:**
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

**Ответ (201 Created):**
```json
{
  "id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "balance": 0
}
```

### 2. Вход в систему

**POST** `/api/v1/login`

**Тело запроса:**
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**Ответ (200 OK):**
```json
{
  "id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "balance": 100.50
}
```

### 3. Пополнение баланса

**POST** `/api/v1/balance/deposit?user_id=1`

**Тело запроса:**
```json
{
  "amount": 100.50
}
```

**Ответ (200 OK):**
```json
{
  "balance": 100.50
}
```

### 4. Игра на спинах

**POST** `/api/v1/spin?user_id=1`

**Тело запроса:**
```json
{
  "bet_amount": 10.00
}
```

**Ответ (200 OK):**
```json
{
  "reel1": 7,
  "reel2": 7,
  "reel3": 7,
  "is_win": true,
  "win_amount": 100.00,
  "balance": 190.00
}
```

## Правила игры на спинах

### Символы и вероятности

- **0**: 0.5% (джекпот символ)
- **1-3**: 5% каждый
- **4-6**: 10% каждый
- **7-9**: 20% каждый

### Выигрышные комбинации

#### Три одинаковых символа (джекпот):
- Три нуля: **x1000** от ставки
- Три единицы/двойки/тройки: **x50** от ставки
- Три четверки/пятерки/шестерки: **x20** от ставки
- Три семерки/восьмерки/девятки: **x10** от ставки

#### Два одинаковых символа:
- Два нуля: **x10** от ставки
- Две единицы/двойки/тройки: **x3** от ставки
- Две четверки/пятерки/шестерки: **x2** от ставки
- Две семерки/восьмерки/девятки: **x1.5** от ставки

#### Последовательность:
- 0-1-2 или 7-8-9: **x5** от ставки

### RTP (Return to Player)

Система настроена на RTP ~95%, что соответствует стандартам классических казино.

## Примеры использования

### Регистрация и пополнение баланса

```bash
# Регистрация
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"player1","email":"player1@example.com","password":"secret123"}'

# Пополнение баланса на 1000 рублей
curl -X POST "http://localhost:8080/api/v1/balance/deposit?user_id=1" \
  -H "Content-Type: application/json" \
  -d '{"amount":1000}'

# Игра на спинах со ставкой 10 рублей
curl -X POST "http://localhost:8080/api/v1/spin?user_id=1" \
  -H "Content-Type: application/json" \
  -d '{"bet_amount":10}'
```

