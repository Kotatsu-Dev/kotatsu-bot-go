package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"rr/kotatsutgbot/config"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
)

func signData(data string) string {
	hmacSecret := []byte(config.GetConfig().AUTH_SECRET)
	mac := hmac.New(sha256.New, hmacSecret)
	mac.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func verifyData(data, signature string) bool {
	expectedSig := signData(data)
	return hmac.Equal([]byte(expectedSig), []byte(signature))
}

func CreateSessionCookie(userID string, validFor time.Duration) string {
	expiry := time.Now().Add(validFor).Unix()
	data := fmt.Sprintf("%d:%s", expiry, userID)
	sig := signData(data)
	return fmt.Sprintf("%s:%s", data, sig)
}

func CheckIsAdmin(userId int64) bool {
	b, err := bot.New(config.GetConfig().CONFIG_BOT_TOKEN)
	if err != nil {
		return false
	}

	admins, err := b.GetChatAdministrators(context.TODO(), &bot.GetChatAdministratorsParams{
		ChatID: config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
	})
	if err != nil {
		return false
	}

	for _, admin := range admins {
		if admin.Owner != nil {
			if admin.Owner.User.ID == userId {
				return true
			}
		} else if admin.Administrator != nil {
			if admin.Administrator.User.ID == userId {
				return true
			}
		}
	}

	return false
}

func ParseAndVerifySessionCookie(cookieValue string) (userID int64, isValid bool) {
	parts := strings.Split(cookieValue, ":")
	if len(parts) != 3 {
		return 0, false
	}

	expiryStr, userIDStr, signature := parts[0], parts[1], parts[2]
	data := fmt.Sprintf("%s:%s", expiryStr, userIDStr)

	if !verifyData(data, signature) {
		return 0, false
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, false
	}

	expiry, err := strconv.ParseInt(expiryStr, 10, 64)
	if err != nil {
		return 0, false
	}

	if time.Now().Unix() > expiry {
		return 0, false
	}

	return userID, true
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid session cookie",
			})
			return
		}

		userID, isValid := ParseAndVerifySessionCookie(cookie)
		if !isValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired session",
			})
			return
		}

		if !CheckIsAdmin(userID) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "you shall not pass",
			})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
