package controller

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"

	"nofelet/internal/domain/signaling/controller/view"
)

const usernamePrefix = "nofelet_user"

// GetCoTURNCredentials /turn-credentials/generate генерит временные креды для CoTURN
func (c *Controller) GetCoTURNCredentials(ctx *gin.Context) {
	conn, sErr := NewWebSocket(ctx, c.Logger)
	if sErr != nil {
		c.Logger.Error("socket creation", slog.Any("err", sErr))
	}
	defer conn.Close()

	// CoTURN по умолчанию использует интервал в 60 секунд для временной метки.
	// Имя пользователя = Текущий_Timestamp_в_минутах : Префикс_пользователя

	timestamp := time.Now().Unix()
	// Логин (username) - это временная метка в минутах
	login := fmt.Sprintf("%d", timestamp/60)

	// Если нужен префикс, можно добавить (необязательно, зависит от вашей настройки CoTURN)
	login = fmt.Sprintf("%s:%s", login, usernamePrefix)

	// Генерируем временный пароль с помощью HMAC-SHA1 хэша от логина и общего секрета
	h := hmac.New(sha1.New, []byte(c.Config.CoTURN.SharedSecret))
	h.Write([]byte(login))
	sha1_hash := h.Sum(nil)

	// Пароль должен быть закодирован в Base64
	password := base64.StdEncoding.EncodeToString(sha1_hash)

	// Формируем структуру ответа для клиента
	config := view.TURNConfig{
		ICEServers: []view.ICEServer{
			{
				// STUN-сервер
				URLs: fmt.Sprintf("stun:%s:3478", c.Config.CoTURN.TurnServerIP),
			},
			{
				// TURN-сервер
				URLs:       fmt.Sprintf("turn:%s:3478", c.Config.CoTURN.TurnServerIP),
				Username:   login,
				Credential: password,
			},
		},
	}

	err := conn.WriteJSON(config)
	if err != nil {
		c.Logger.Error("coturn credentials", slog.Any("err", err))
	}
}
