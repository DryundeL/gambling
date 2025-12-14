package console

import (
	"bufio"
	"fmt"
	"gambling/internal/model"
	"gambling/internal/service"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Console представляет консольный интерфейс для игры
type Console struct {
	authService    *service.AuthService
	balanceService *service.BalanceService
	spinService    *service.SpinService
	scanner        *bufio.Scanner
	currentUser    *model.User
}

// NewConsole создает новый экземпляр консольного интерфейса
func NewConsole(
	authService *service.AuthService,
	balanceService *service.BalanceService,
	spinService *service.SpinService,
) *Console {
	return &Console{
		authService:    authService,
		balanceService: balanceService,
		spinService:    spinService,
		scanner:        bufio.NewScanner(os.Stdin),
	}
}

// Run запускает консольное приложение
func (c *Console) Run() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║     Добро пожаловать в Казино! 🎰     ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	for {
		if c.currentUser == nil {
			c.showAuthMenu()
		} else {
			c.showMainMenu()
		}
	}
}

// showAuthMenu показывает меню аутентификации
func (c *Console) showAuthMenu() {
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("1. Регистрация")
	fmt.Println("2. Вход")
	fmt.Println("3. Выход")
	fmt.Println("═══════════════════════════════════════")
	fmt.Print("Выберите действие: ")

	c.scanner.Scan()
	choice := strings.TrimSpace(c.scanner.Text())

	switch choice {
	case "1":
		c.register()
	case "2":
		c.login()
	case "3":
		fmt.Println("До свидания!")
		os.Exit(0)
	default:
		fmt.Println("❌ Неверный выбор. Попробуйте снова.")
		fmt.Println()
	}
}

// showMainMenu показывает главное меню игры
func (c *Console) showMainMenu() {
	fmt.Println()
	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("👤 Пользователь: %s\n", c.currentUser.Username)
	fmt.Printf("💰 Баланс: %.2f ₽\n", c.currentUser.Balance)
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("1. Пополнить баланс")
	fmt.Println("2. Играть в спинах")
	fmt.Println("3. Выйти из аккаунта")
	fmt.Println("4. Выход из программы")
	fmt.Println("═══════════════════════════════════════")
	fmt.Print("Выберите действие: ")

	c.scanner.Scan()
	choice := strings.TrimSpace(c.scanner.Text())

	switch choice {
	case "1":
		c.deposit()
	case "2":
		c.playSpin()
	case "3":
		c.currentUser = nil
		fmt.Println("✅ Вы вышли из аккаунта")
		fmt.Println()
	case "4":
		fmt.Println("До свидания!")
		os.Exit(0)
	default:
		fmt.Println("❌ Неверный выбор. Попробуйте снова.")
	}
}

// register обрабатывает регистрацию
func (c *Console) register() {
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📝 РЕГИСТРАЦИЯ")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	fmt.Print("Имя пользователя: ")
	c.scanner.Scan()
	username := strings.TrimSpace(c.scanner.Text())

	fmt.Print("Email: ")
	c.scanner.Scan()
	email := strings.TrimSpace(c.scanner.Text())

	fmt.Print("Пароль: ")
	c.scanner.Scan()
	password := strings.TrimSpace(c.scanner.Text())

	if username == "" || email == "" || password == "" {
		fmt.Println("❌ Все поля обязательны для заполнения!")
		fmt.Println()
		return
	}

	user, err := c.authService.Register(username, email, password)
	if err != nil {
		if err == service.ErrUserExists {
			fmt.Println("❌ Пользователь с таким именем или email уже существует!")
		} else {
			fmt.Printf("❌ Ошибка при регистрации: %v\n", err)
		}
		fmt.Println()
		return
	}

	c.currentUser = user
	fmt.Printf("✅ Регистрация успешна! Добро пожаловать, %s!\n", user.Username)
	fmt.Println()
}

// login обрабатывает вход
func (c *Console) login() {
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🔐 ВХОД")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	fmt.Print("Имя пользователя: ")
	c.scanner.Scan()
	username := strings.TrimSpace(c.scanner.Text())

	fmt.Print("Пароль: ")
	c.scanner.Scan()
	password := strings.TrimSpace(c.scanner.Text())

	if username == "" || password == "" {
		fmt.Println("❌ Имя пользователя и пароль обязательны!")
		fmt.Println()
		return
	}

	user, err := c.authService.Login(username, password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			fmt.Println("❌ Неверные учетные данные!")
		} else {
			fmt.Printf("❌ Ошибка при входе: %v\n", err)
		}
		fmt.Println()
		return
	}

	c.currentUser = user
	fmt.Printf("✅ Вход выполнен! Добро пожаловать, %s!\n", user.Username)
	fmt.Printf("💰 Ваш баланс: %.2f ₽\n", user.Balance)
	fmt.Println()
}

// deposit обрабатывает пополнение баланса
func (c *Console) deposit() {
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💳 ПОПОЛНЕНИЕ БАЛАНСА")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Текущий баланс: %.2f ₽\n", c.currentUser.Balance)
	fmt.Print("Введите сумму для пополнения: ")

	c.scanner.Scan()
	amountStr := strings.TrimSpace(c.scanner.Text())

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		fmt.Println("❌ Неверная сумма!")
		fmt.Println()
		return
	}

	user, err := c.balanceService.Deposit(c.currentUser.ID, amount)
	if err != nil {
		fmt.Printf("❌ Ошибка при пополнении: %v\n", err)
		fmt.Println()
		return
	}

	c.currentUser = user
	fmt.Printf("✅ Баланс успешно пополнен на %.2f ₽\n", amount)
	fmt.Printf("💰 Новый баланс: %.2f ₽\n", user.Balance)
	fmt.Println()
}

// playSpin обрабатывает игру на спинах
func (c *Console) playSpin() {
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎰 ИГРА НА СПИНАХ")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Текущий баланс: %.2f ₽\n", c.currentUser.Balance)
	fmt.Print("Введите сумму ставки: ")

	c.scanner.Scan()
	betStr := strings.TrimSpace(c.scanner.Text())

	betAmount, err := strconv.ParseFloat(betStr, 64)
	if err != nil || betAmount <= 0 {
		fmt.Println("❌ Неверная сумма ставки!")
		fmt.Println()
		return
	}

	if betAmount > c.currentUser.Balance {
		fmt.Println("❌ Недостаточно средств на балансе!")
		fmt.Println()
		return
	}

	fmt.Println()
	fmt.Println("🎰 Крутим барабаны...")
	fmt.Println()

	// Сначала выполняем спин (получаем результат)
	result, err := c.spinService.Spin(c.currentUser.ID, betAmount)
	if err != nil {
		if err == service.ErrInsufficientFunds {
			fmt.Println("❌ Недостаточно средств!")
		} else {
			fmt.Printf("❌ Ошибка при игре: %v\n", err)
		}
		fmt.Println()
		return
	}

	// Обновляем баланс пользователя
	c.currentUser.Balance = result.Balance

	// Показываем анимацию вращения барабанов
	c.animateSpin(result.Reel1, result.Reel2, result.Reel3)

	if result.IsWin {
		fmt.Printf("🎉 ВЫИГРЫШ! Вы выиграли %.2f ₽\n", result.WinAmount)
	} else {
		fmt.Println("😔 Не повезло, попробуйте еще раз!")
	}

	fmt.Printf("💰 Ваш баланс: %.2f ₽\n", result.Balance)
	fmt.Println()

	// Показываем правила выигрыша
	c.showWinRules()
}

// animateSpin показывает анимацию вращения барабанов с постепенным замедлением
func (c *Console) animateSpin(reel1, reel2, reel3 int) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Количество оборотов для каждого барабана
	spins1 := 15 + rng.Intn(10) // 15-24 оборота
	spins2 := 15 + rng.Intn(10)
	spins3 := 15 + rng.Intn(10)

	fmt.Println("╔═══════════════════════════════════════╗")

	// Первый барабан с замедлением
	c.spinReel1(rng, spins1, reel1)
	time.Sleep(400 * time.Millisecond)

	// Второй барабан с замедлением
	c.spinReel2(rng, spins2, reel1, reel2)
	time.Sleep(400 * time.Millisecond)

	// Третий барабан с замедлением
	c.spinReel3(rng, spins3, reel1, reel2, reel3)

	fmt.Println()
	fmt.Println("╚═══════════════════════════════════════╝")
}

// spinReel1 вращает первый барабан
func (c *Console) spinReel1(rng *rand.Rand, totalSpins, finalSymbol int) {
	fastSpins := totalSpins - 5
	slowSpins := 5

	// Быстрое вращение
	for i := 0; i < fastSpins; i++ {
		symbol := rng.Intn(10)
		fmt.Printf("\r║         [%d] [ ] [ ]          ║", symbol)
		time.Sleep(50 * time.Millisecond)
	}

	// Замедление перед остановкой
	delays := []time.Duration{100, 150, 200, 250, 300}
	for i := 0; i < slowSpins; i++ {
		symbol := rng.Intn(10)
		fmt.Printf("\r║         [%d] [ ] [ ]          ║", symbol)
		if i < len(delays) {
			time.Sleep(delays[i])
		} else {
			time.Sleep(300 * time.Millisecond)
		}
	}

	// Финальный символ
	fmt.Printf("\r║         [%d] [ ] [ ]          ║", finalSymbol)
}

// spinReel2 вращает второй барабан
func (c *Console) spinReel2(rng *rand.Rand, totalSpins, reel1, finalSymbol int) {
	fastSpins := totalSpins - 5
	slowSpins := 5

	// Быстрое вращение
	for i := 0; i < fastSpins; i++ {
		symbol := rng.Intn(10)
		fmt.Printf("\r║         [%d] [%d] [ ]          ║", reel1, symbol)
		time.Sleep(50 * time.Millisecond)
	}

	// Замедление перед остановкой
	delays := []time.Duration{100, 150, 200, 250, 300}
	for i := 0; i < slowSpins; i++ {
		symbol := rng.Intn(10)
		fmt.Printf("\r║         [%d] [%d] [ ]          ║", reel1, symbol)
		if i < len(delays) {
			time.Sleep(delays[i])
		} else {
			time.Sleep(300 * time.Millisecond)
		}
	}

	// Финальный символ
	fmt.Printf("\r║         [%d] [%d] [ ]          ║", reel1, finalSymbol)
}

// spinReel3 вращает третий барабан
func (c *Console) spinReel3(rng *rand.Rand, totalSpins, reel1, reel2, finalSymbol int) {
	fastSpins := totalSpins - 5
	slowSpins := 5

	// Быстрое вращение
	for i := 0; i < fastSpins; i++ {
		symbol := rng.Intn(10)
		fmt.Printf("\r║         [%d] [%d] [%d]          ║", reel1, reel2, symbol)
		time.Sleep(50 * time.Millisecond)
	}

	// Замедление перед остановкой
	delays := []time.Duration{100, 150, 200, 250, 300}
	for i := 0; i < slowSpins; i++ {
		symbol := rng.Intn(10)
		fmt.Printf("\r║         [%d] [%d] [%d]          ║", reel1, reel2, symbol)
		if i < len(delays) {
			time.Sleep(delays[i])
		} else {
			time.Sleep(300 * time.Millisecond)
		}
	}

	// Финальный символ
	fmt.Printf("\r║         [%d] [%d] [%d]          ║", reel1, reel2, finalSymbol)
}

// showWinRules показывает правила выигрыша
func (c *Console) showWinRules() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📋 ПРАВИЛА ВЫИГРЫША:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Три одинаковых:")
	fmt.Println("  • Три нуля: x1000")
	fmt.Println("  • Три 1-3: x50")
	fmt.Println("  • Три 4-6: x20")
	fmt.Println("  • Три 7-9: x10")
	fmt.Println()
	fmt.Println("Два одинаковых:")
	fmt.Println("  • Два нуля: x10")
	fmt.Println("  • Две 1-3: x3")
	fmt.Println("  • Две 4-6: x2")
	fmt.Println("  • Две 7-9: x1.5")
	fmt.Println()
	fmt.Println("Последовательность (0-1-2 или 7-8-9): x5")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
}
