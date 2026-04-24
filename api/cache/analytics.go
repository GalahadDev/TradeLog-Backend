// Package cache provee un caché TTL en memoria para resultados de analíticas.
// Evita recalcular TradingStats en cada request cuando los trades no han cambiado.
// Las entradas se invalidan de forma inmediata ante cualquier mutación de trades.
package cache

import (
	"strings"
	"sync"
	"time"

	"samll-trading-back/api/services/analytics"
)

const statsTTL = 5 * time.Minute

type statsEntry struct {
	stats     analytics.TradingStats
	expiresAt time.Time
}

// StatsCache es un caché TTL con seguridad para goroutines, indexado por "accountID:startDate:endDate".
type StatsCache struct {
	mu    sync.RWMutex
	store map[string]statsEntry
}

// Stats es la instancia global del caché utilizada por los handlers.
var Stats = &StatsCache{store: make(map[string]statsEntry)}

// Key construye la clave canónica del caché para una solicitud de estadísticas.
func Key(accountID, startDate, endDate string) string {
	return accountID + ":" + startDate + ":" + endDate
}

// Get retorna las estadísticas cacheadas y true si existe una entrada válida (no expirada).
func (c *StatsCache) Get(key string) (analytics.TradingStats, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.store[key]
	if !ok || time.Now().After(e.expiresAt) {
		return analytics.TradingStats{}, false
	}
	return e.stats, true
}

// Set almacena las estadísticas bajo la clave dada con el TTL por defecto.
func (c *StatsCache) Set(key string, stats analytics.TradingStats) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = statsEntry{stats: stats, expiresAt: time.Now().Add(statsTTL)}
}

// Invalidate elimina todas las entradas cacheadas que pertenezcan al accountID dado.
// Se llama cuando se crea, actualiza o elimina un trade de esa cuenta.
func (c *StatsCache) Invalidate(accountID string) {
	prefix := accountID + ":"
	c.mu.Lock()
	defer c.mu.Unlock()
	for k := range c.store {
		if strings.HasPrefix(k, prefix) {
			delete(c.store, k)
		}
	}
}
