package console

import (
	"bufio"
	"fmt"
	"gambling/internal/model"
	"gambling/internal/service"
	"os"
	"strconv"
	"strings"
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

	// Отображаем результат
	fmt.Println("╔═══════════════════════════════════════╗")
	fmt.Printf("║         [%d] [%d] [%d]          ║\n", result.Reel1, result.Reel2, result.Reel3)
	fmt.Println("╚═══════════════════════════════════════╝")

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

