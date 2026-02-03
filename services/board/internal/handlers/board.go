// Package handlers — HTTP-обработчики Board Service.
package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	ws "github.com/vokzal-tech/board-service/internal/websocket"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// isAllowedOrigin возвращает true, если origin разрешён: либо allowAllInDev включён (для разработки),
// либо origin совпадает с одним из записей в allowed (точное совпадение).
func isAllowedOrigin(origin string, allowed []string, allowAllInDev bool) bool {
	if allowAllInDev {
		return true
	}
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return false
	}
	for _, a := range allowed {
		if strings.TrimSpace(a) == origin {
			return true
		}
	}
	return false
}

// BoardHandler обрабатывает HTTP-запросы к API табло.
type BoardHandler struct {
	db       *gorm.DB
	hub      *ws.Hub
	logger   *zap.Logger
	upgrader websocket.Upgrader
}

// NewBoardHandler создаёт новый BoardHandler. allowedOrigins — whitelist origins для WebSocket CheckOrigin;
// allowAllOriginsInDev при true отключает проверку (удобно для разработки).
func NewBoardHandler(db *gorm.DB, hub *ws.Hub, logger *zap.Logger, allowedOrigins []string, allowAllOriginsInDev bool) *BoardHandler {
	h := &BoardHandler{
		db:     db,
		hub:    hub,
		logger: logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				return isAllowedOrigin(origin, allowedOrigins, allowAllOriginsInDev)
			},
		},
	}
	return h
}

// HandleWebSocket обрабатывает WebSocket-подключение для табло.
func (h *BoardHandler) HandleWebSocket(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade to WebSocket", zap.Error(err))
		return
	}

	ws.ServeWs(h.hub, conn)
}

// GetPublicBoard возвращает данные для общего табло.
func (h *BoardHandler) GetPublicBoard(c *gin.Context) {
	date := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	// Запрос рейсов из БД
	var trips []map[string]interface{}
	query := `
		SELECT 
			t.id, t.date, t.status, t.delay_minutes, t.platform,
			s.departure_time, r.name as route_name
		FROM trips t
		JOIN schedules s ON s.id = t.schedule_id
		JOIN routes r ON r.id = s.route_id
		WHERE t.date = ?
		ORDER BY s.departure_time ASC
	`

	rows, err := h.db.Raw(query, date).Rows()
	if err != nil {
		h.logger.Error("Failed to query trips", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trips"})
		return
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var trip map[string]interface{}
		_ = h.db.ScanRows(rows, &trip)
		trips = append(trips, trip)
	}

	c.JSON(http.StatusOK, gin.H{"data": trips})
}

// GetPlatformBoard возвращает данные для перронного табло.
func (h *BoardHandler) GetPlatformBoard(c *gin.Context) {
	platform := c.Param("platform")
	date := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	var trips []map[string]interface{}
	query := `
		SELECT 
			t.id, t.date, t.status, t.delay_minutes,
			s.departure_time, r.name as route_name,
			COUNT(tk.id) as total_tickets,
			COUNT(bm.id) as boarded_count
		FROM trips t
		JOIN schedules s ON s.id = t.schedule_id
		JOIN routes r ON r.id = s.route_id
		LEFT JOIN tickets tk ON tk.trip_id = t.id AND tk.status = 'active'
		LEFT JOIN boarding_marks bm ON bm.ticket_id = tk.id
		WHERE t.date = ? AND t.platform = ?
		GROUP BY t.id, s.departure_time, r.name
		ORDER BY s.departure_time ASC
	`

	rows, err := h.db.Raw(query, date, platform).Rows()
	if err != nil {
		h.logger.Error("Failed to query platform trips", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trips"})
		return
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var trip map[string]interface{}
		_ = h.db.ScanRows(rows, &trip)
		trips = append(trips, trip)
	}

	c.JSON(http.StatusOK, gin.H{"data": trips})
}

// GetWebSocketStats возвращает статистику WebSocket-соединений.
func (h *BoardHandler) GetWebSocketStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"connected_clients": h.hub.GetClientCount(),
	})
}
