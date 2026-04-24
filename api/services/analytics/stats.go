package analytics

import (
	"math"
	"samll-trading-back/api/domains"

	"github.com/shopspring/decimal"
)

type TradingStats struct {
	// Rendimiento General
	TotalNetProfit decimal.Decimal `json:"total_net_profit"`
	GrossProfit    decimal.Decimal `json:"gross_profit"`
	GrossLoss      decimal.Decimal `json:"gross_loss"`
	ProfitFactor   decimal.Decimal `json:"profit_factor"`
	RecoveryFactor decimal.Decimal `json:"recovery_factor"`
	SharpeRatio    decimal.Decimal `json:"sharpe_ratio"`
	ExpectedPayoff decimal.Decimal `json:"expected_payoff"`

	// Gastos
	TotalCommissions decimal.Decimal `json:"total_commissions"`

	// Actividad
	TotalTrades  int             `json:"total_trades"`
	AvgTradeSize decimal.Decimal `json:"avg_trade_size"`

	// Drawdown
	MaxDrawdown decimal.Decimal `json:"max_drawdown"`

	// Promedios y Extremos
	AvgWin      decimal.Decimal `json:"avg_win"`
	AvgLoss     decimal.Decimal `json:"avg_loss"`
	LargestWin  decimal.Decimal `json:"largest_win"`
	LargestLoss decimal.Decimal `json:"largest_loss"`

	// Tasas de Éxito
	WinRate      decimal.Decimal `json:"win_rate"`
	LossRate     decimal.Decimal `json:"loss_rate"`
	LongWinRate  decimal.Decimal `json:"long_win_rate"`
	ShortWinRate decimal.Decimal `json:"short_win_rate"`

	// Rachas Consecutivas
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

	// Para Sharpe Ratio (pre-alocado con capacidad exacta para evitar re-allocaciones)
	pnlHistory := make([]decimal.Decimal, 0, len(trades))

	// Para Rachas Consecutivas
	currentWinStreak := 0
	currentLossStreak := 0
	currentWinStreakMoney := decimal.Zero
	currentLossStreakMoney := decimal.Zero

	for _, t := range trades {
		pnl := t.PnL
		pnlHistory = append(pnlHistory, pnl)

		// Acumular tamaño de posición (lotes/size)
		totalSize = totalSize.Add(t.Size)

		// 1. Acumulados Generales
		stats.TotalNetProfit = stats.TotalNetProfit.Add(pnl)

		// 2. Análisis de ganancias/pérdidas y rachas
		if pnl.GreaterThan(decimal.Zero) {
			// Es ganancia
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

			// Análisis por dirección
			if t.Direction == "LONG" {
				longWins++
			} else {
				shortWins++
			}

		} else {
			// Es pérdida (o break even negativo)
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

	// Factor de ganancia
	if !stats.GrossLoss.IsZero() {
		stats.ProfitFactor = stats.GrossProfit.Div(stats.GrossLoss.Abs())
	} else {
		stats.ProfitFactor = stats.GrossProfit
	}

	// Expectativa (ganancia esperada): Beneficio neto total / Cantidad de trades
	stats.ExpectedPayoff = stats.TotalNetProfit.Div(decimal.NewFromInt(int64(stats.TotalTrades)))

	// Tamaño promedio de posición
	stats.AvgTradeSize = totalSize.Div(decimal.NewFromInt(int64(stats.TotalTrades)))

	// Factor de recuperación: Beneficio neto / Drawdown máximo
	if !stats.MaxDrawdown.IsZero() {
		stats.RecoveryFactor = stats.TotalNetProfit.Div(stats.MaxDrawdown)
	} else {
		// Si no hay drawdown, el factor es teóricamente infinito o igual al profit
		stats.RecoveryFactor = decimal.Zero
		if stats.TotalNetProfit.GreaterThan(decimal.Zero) {
			stats.RecoveryFactor = stats.TotalNetProfit
		}
	}

	// Tasas de acierto/fallo
	totalTradesDec := decimal.NewFromInt(int64(stats.TotalTrades))
	stats.WinRate = decimal.NewFromInt(int64(wins)).Div(totalTradesDec).Mul(decimal.NewFromInt(100))
	stats.LossRate = decimal.NewFromInt(int64(lossCount)).Div(totalTradesDec).Mul(decimal.NewFromInt(100))

	// Promedios
	if wins > 0 {
		stats.AvgWin = winSum.Div(decimal.NewFromInt(int64(wins)))
	}
	if lossCount > 0 {
		stats.AvgLoss = lossSum.Div(decimal.NewFromInt(int64(lossCount)))
	}

	// Tasas de acierto por dirección
	if longs > 0 {
		stats.LongWinRate = decimal.NewFromInt(int64(longWins)).Div(decimal.NewFromInt(int64(longs))).Mul(decimal.NewFromInt(100))
	}
	if shorts > 0 {
		stats.ShortWinRate = decimal.NewFromInt(int64(shortWins)).Div(decimal.NewFromInt(int64(shorts))).Mul(decimal.NewFromInt(100))
	}

	// Sharpe Ratio (Simplificado: Promedio / Desviación Estándar)
	// La media y varianza se calculan en decimal; solo se convierte a float64 para math.Sqrt.
	if len(pnlHistory) > 1 {
		stdDev := calculateStdDev(pnlHistory)
		if stdDev.IsPositive() {
			stats.SharpeRatio = stats.ExpectedPayoff.Div(stdDev)
		}
	}

	return stats
}

// calculateStdDev calcula la desviación estándar poblacional en decimal.
// La media y varianza se calculan con precisión decimal; solo el Sqrt final usa float64.
func calculateStdDev(data []decimal.Decimal) decimal.Decimal {
	n := decimal.NewFromInt(int64(len(data)))

	sum := decimal.Zero
	for _, v := range data {
		sum = sum.Add(v)
	}
	mean := sum.Div(n)

	variance := decimal.Zero
	for _, v := range data {
		diff := v.Sub(mean)
		variance = variance.Add(diff.Mul(diff))
	}
	variance = variance.Div(n)

	varianceF64, _ := variance.Float64()
	return decimal.NewFromFloat(math.Sqrt(varianceF64))
}
