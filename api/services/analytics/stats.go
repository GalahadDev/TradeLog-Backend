package analytics

import (
	"math"
	"samll-trading-back/api/domains"

	"github.com/shopspring/decimal"
)

type TradingStats struct {
	// --- Rendimiento General ---
	TotalNetProfit decimal.Decimal `json:"total_net_profit"`
	GrossProfit    decimal.Decimal `json:"gross_profit"`
	GrossLoss      decimal.Decimal `json:"gross_loss"`
	ProfitFactor   decimal.Decimal `json:"profit_factor"`
	RecoveryFactor decimal.Decimal `json:"recovery_factor"`
	SharpeRatio    decimal.Decimal `json:"sharpe_ratio"`
	ExpectedPayoff decimal.Decimal `json:"expected_payoff"`

	// --- Gastos ---
	TotalCommissions decimal.Decimal `json:"total_commissions"`

	// --- Actividad ---
	TotalTrades  int             `json:"total_trades"`
	AvgTradeSize decimal.Decimal `json:"avg_trade_size"`

	// --- Drawdown ---
	MaxDrawdown decimal.Decimal `json:"max_drawdown"`

	// --- Promedios y Extremos ---
	AvgWin      decimal.Decimal `json:"avg_win"`
	AvgLoss     decimal.Decimal `json:"avg_loss"`
	LargestWin  decimal.Decimal `json:"largest_win"`
	LargestLoss decimal.Decimal `json:"largest_loss"`

	// --- Tasas de Éxito ---
	WinRate      decimal.Decimal `json:"win_rate"`
	LossRate     decimal.Decimal `json:"loss_rate"`
	LongWinRate  decimal.Decimal `json:"long_win_rate"`
	ShortWinRate decimal.Decimal `json:"short_win_rate"`

	// --- Rachas (Consecutive) ---
	MaxConsecutiveWins   int             `json:"max_consecutive_wins"`
	MaxConsecutiveLosses int             `json:"max_consecutive_losses"`
	MaxConsecutiveProfit decimal.Decimal `json:"max_consecutive_profit_usd"`
	MaxConsecutiveLoss   decimal.Decimal `json:"max_consecutive_loss_usd"`
}

func CalculateStats(trades []domains.Trade) TradingStats {
	stats := TradingStats{
		TotalNetProfit:       decimal.Zero,
		GrossProfit:          decimal.Zero,
		GrossLoss:            decimal.Zero,
		MaxDrawdown:          decimal.Zero,
		TotalCommissions:     decimal.Zero,
		MaxConsecutiveProfit: decimal.Zero,
		MaxConsecutiveLoss:   decimal.Zero,
	}

	if len(trades) == 0 {
		return stats
	}

	stats.TotalTrades = len(trades)

	// Variables auxiliares
	wins := 0
	winSum := decimal.Zero
	lossCount := 0
	lossSum := decimal.Zero
	totalSize := decimal.Zero

	longs := 0
	longWins := 0
	shorts := 0
	shortWins := 0

	// Para Drawdown
	peakBalance := decimal.Zero
	currentBalance := decimal.Zero

	// Para Sharpe Ratio (Almacenamos los PnLs)
	var pnlHistory []float64

	// Para Rachas (Consecutive)
	currentWinStreak := 0
	currentLossStreak := 0
	currentWinStreakMoney := decimal.Zero
	currentLossStreakMoney := decimal.Zero

	for _, t := range trades {
		pnl := t.PnL
		pnlHistory = append(pnlHistory, pnl.InexactFloat64())

		// Acumular Tamaño (Lots/Size)
		totalSize = totalSize.Add(t.Size)

		// 1. Acumulados Generales
		stats.TotalNetProfit = stats.TotalNetProfit.Add(pnl)

		// 2. Análisis Win/Loss y Rachas
		if pnl.GreaterThan(decimal.Zero) {
			// Es WIN
			stats.GrossProfit = stats.GrossProfit.Add(pnl)
			wins++
			winSum = winSum.Add(pnl)

			if pnl.GreaterThan(stats.LargestWin) {
				stats.LargestWin = pnl
			}

			// Manejo de Rachas Ganadoras
			currentWinStreak++
			currentWinStreakMoney = currentWinStreakMoney.Add(pnl)

			// Reset de Rachas Perdedoras
			currentLossStreak = 0
			currentLossStreakMoney = decimal.Zero

			if currentWinStreak > stats.MaxConsecutiveWins {
				stats.MaxConsecutiveWins = currentWinStreak
			}
			if currentWinStreakMoney.GreaterThan(stats.MaxConsecutiveProfit) {
				stats.MaxConsecutiveProfit = currentWinStreakMoney
			}

			// Direction Analysis
			if t.Direction == "LONG" {
				longWins++
			} else {
				shortWins++
			}

		} else {
			// Es LOSS (o Break Even negativo)
			stats.GrossLoss = stats.GrossLoss.Add(pnl)
			lossCount++
			lossSum = lossSum.Add(pnl)

			if pnl.LessThan(stats.LargestLoss) {
				stats.LargestLoss = pnl
			}

			// Manejo de Rachas Perdedoras
			currentLossStreak++
			currentLossStreakMoney = currentLossStreakMoney.Add(pnl)

			// Reset de Rachas Ganadoras
			currentWinStreak = 0
			currentWinStreakMoney = decimal.Zero

			if currentLossStreak > stats.MaxConsecutiveLosses {
				stats.MaxConsecutiveLosses = currentLossStreak
			}
			// Como es dinero negativo, buscamos el "menor" número
			if currentLossStreakMoney.LessThan(stats.MaxConsecutiveLoss) {
				stats.MaxConsecutiveLoss = currentLossStreakMoney
			}
		}

		if t.Direction == "LONG" {
			longs++
		} else {
			shorts++
		}

		// 3. Cálculo de Drawdown
		currentBalance = currentBalance.Add(pnl)
		if currentBalance.GreaterThan(peakBalance) {
			peakBalance = currentBalance
		}
		drawdown := peakBalance.Sub(currentBalance)
		if drawdown.GreaterThan(stats.MaxDrawdown) {
			stats.MaxDrawdown = drawdown
		}

		// Acumular Comisiones
		stats.TotalCommissions = stats.TotalCommissions.Add(t.Commission)
	}

	// Restar comisiones del Net Profit
	stats.TotalNetProfit = stats.TotalNetProfit.Sub(stats.TotalCommissions)

	// --- CÁLCULOS FINALES ---

	// Profit Factor
	if !stats.GrossLoss.IsZero() {
		stats.ProfitFactor = stats.GrossProfit.Div(stats.GrossLoss.Abs())
	} else {
		stats.ProfitFactor = stats.GrossProfit
	}

	// Expectancy (Expected Payoff): Total Net Profit / Count
	stats.ExpectedPayoff = stats.TotalNetProfit.Div(decimal.NewFromInt(int64(stats.TotalTrades)))

	// Avg Trade Size
	stats.AvgTradeSize = totalSize.Div(decimal.NewFromInt(int64(stats.TotalTrades)))

	// Recovery Factor: Net Profit / Max Drawdown
	if !stats.MaxDrawdown.IsZero() {
		stats.RecoveryFactor = stats.TotalNetProfit.Div(stats.MaxDrawdown)
	} else {
		// Si no hay drawdown, el factor es teóricamente infinito o igual al profit
		stats.RecoveryFactor = decimal.Zero
		if stats.TotalNetProfit.GreaterThan(decimal.Zero) {
			stats.RecoveryFactor = stats.TotalNetProfit
		}
	}

	// Win/Loss Rates
	totalTradesDec := decimal.NewFromInt(int64(stats.TotalTrades))
	stats.WinRate = decimal.NewFromInt(int64(wins)).Div(totalTradesDec).Mul(decimal.NewFromInt(100))
	stats.LossRate = decimal.NewFromInt(int64(lossCount)).Div(totalTradesDec).Mul(decimal.NewFromInt(100))

	// Averages
	if wins > 0 {
		stats.AvgWin = winSum.Div(decimal.NewFromInt(int64(wins)))
	}
	if lossCount > 0 {
		stats.AvgLoss = lossSum.Div(decimal.NewFromInt(int64(lossCount)))
	}

	// Direction Win Rates
	if longs > 0 {
		stats.LongWinRate = decimal.NewFromInt(int64(longWins)).Div(decimal.NewFromInt(int64(longs))).Mul(decimal.NewFromInt(100))
	}
	if shorts > 0 {
		stats.ShortWinRate = decimal.NewFromInt(int64(shortWins)).Div(decimal.NewFromInt(int64(shorts))).Mul(decimal.NewFromInt(100))
	}

	// Sharpe Ratio (Simplificado: Promedio / Desviación Estándar)
	if len(pnlHistory) > 1 {
		stdDev := calculateStdDev(pnlHistory)
		if stdDev > 0 {
			// Usamos Expected Payoff como el promedio por trade
			avgReturn, _ := stats.ExpectedPayoff.Float64()
			stats.SharpeRatio = decimal.NewFromFloat(avgReturn / stdDev)
		}
	}

	return stats
}

// Helper para Desviación Estándar
func calculateStdDev(data []float64) float64 {
	var sum, mean, sd float64
	for _, num := range data {
		sum += num
	}
	mean = sum / float64(len(data))

	for _, num := range data {
		sd += math.Pow(num-mean, 2)
	}
	sd = math.Sqrt(sd / float64(len(data)))
	return sd
}
