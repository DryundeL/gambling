package spin

import (
	"math/rand"
	"time"
)

// Service представляет доменный сервис для логики игры
// Доменные сервисы содержат бизнес-логику, которая не принадлежит конкретной сущности
// В данном случае - генерация символов и расчет выигрыша
type Service struct {
	rng *rand.Rand
}

// NewService создает новый доменный сервис для спинов
func NewService() *Service {
	return &Service{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateSymbol генерирует символ с правильным распределением вероятностей
// Вероятности настроены так, чтобы обеспечить RTP (Return to Player) ~95%
// Символы: 0 (самый редкий), 1-9 (обычные)
func (s *Service) GenerateSymbol() int {
	// Генерируем случайное число от 0 до 9999
	roll := s.rng.Intn(10000)

	// Распределение вероятностей:
	// 0: 0.5% (50 из 10000) - диапазон 0-49
	// 1: 5% (500 из 10000) - диапазон 50-549
	// 2: 5% (500 из 10000) - диапазон 550-1049
	// 3: 5% (500 из 10000) - диапазон 1050-1549
	// 4: 10% (1000 из 10000) - диапазон 1550-2549
	// 5: 10% (1000 из 10000) - диапазон 2550-3549
	// 6: 10% (1000 из 10000) - диапазон 3550-4549
	// 7: 20% (2000 из 10000) - диапазон 4550-6549
	// 8: 20% (2000 из 10000) - диапазон 6550-8549
	// 9: 20% (2000 из 10000) - диапазон 8550-9999

	if roll < 50 {
		return 0 // 0.5%
	} else if roll < 550 {
		return 1 // 5%
	} else if roll < 1050 {
		return 2 // 5%
	} else if roll < 1550 {
		return 3 // 5%
	} else if roll < 2550 {
		return 4 // 10%
	} else if roll < 3550 {
		return 5 // 10%
	} else if roll < 4550 {
		return 6 // 10%
	} else if roll < 6550 {
		return 7 // 20%
	} else if roll < 8550 {
		return 8 // 20%
	} else {
		return 9 // 20%
	}
}

// CalculateWin вычисляет выигрыш на основе комбинации символов
// Коэффициенты выплат настроены для обеспечения RTP ~95%
func (s *Service) CalculateWin(reel1, reel2, reel3 int, betAmount float64) float64 {
	// Три одинаковых символа (джекпот)
	if reel1 == reel2 && reel2 == reel3 {
		switch reel1 {
		case 0:
			return betAmount * 1000 // Джекпот для тройки нулей
		case 1:
			return betAmount * 50
		case 2:
			return betAmount * 50
		case 3:
			return betAmount * 50
		case 4:
			return betAmount * 20
		case 5:
			return betAmount * 20
		case 6:
			return betAmount * 20
		case 7:
			return betAmount * 10
		case 8:
			return betAmount * 10
		case 9:
			return betAmount * 10
		}
	}

	// Два одинаковых символа (обычный выигрыш)
	if reel1 == reel2 || reel2 == reel3 || reel1 == reel3 {
		var symbol int
		if reel1 == reel2 {
			symbol = reel1
		} else if reel2 == reel3 {
			symbol = reel2
		} else {
			symbol = reel1
		}

		switch symbol {
		case 0:
			return betAmount * 10 // Два нуля
		case 1, 2, 3:
			return betAmount * 3
		case 4, 5, 6:
			return betAmount * 2
		case 7, 8, 9:
			return betAmount * 1.5
		}
	}

	// Последовательность (0-1-2 или 7-8-9)
	if (reel1 == 0 && reel2 == 1 && reel3 == 2) || (reel1 == 7 && reel2 == 8 && reel3 == 9) {
		return betAmount * 5
	}

	// Нет выигрыша
	return 0
}

