package controller

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"

	"nofelet/config"
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

	credentials := generateCoTurnCredentials(c.Config)

	err := conn.WriteJSON(credentials)
	if err != nil {
		c.Logger.Error("generate coturn credentials", slog.Any("error", err))
	}
}

func generateCoTurnCredentials(cfg *config.Config) view.TURNConfig {
	// Логин (username) - это временная метка в минутах
	login := fmt.Sprintf("%d", time.Now().Add(24*time.Hour).Unix())
	login = fmt.Sprintf("%s:%s", login, usernamePrefix)

	// Генерируем временный пароль с помощью HMAC-SHA1 хэша от логина и общего секрета
	h := hmac.New(sha1.New, []byte(cfg.CoTURN.SharedSecret))
	h.Write([]byte(login))
	sha1Hash := h.Sum(nil)

	// Пароль должен быть закодирован в Base64
	password := base64.StdEncoding.EncodeToString(sha1Hash)

	// Формируем структуру ответа для клиента
	return view.TURNConfig{
		ICEServers: []view.ICEServer{
			{
				// STUN-сервер
				URLs: fmt.Sprintf("stun:%s:3478", cfg.CoTURN.TurnServerIP),
			},
			{
				// TURN-сервер
				URLs:       fmt.Sprintf("turn:%s:3478", cfg.CoTURN.TurnServerIP),
				Username:   login,
				Credential: password,
			},
		},
	}
}
