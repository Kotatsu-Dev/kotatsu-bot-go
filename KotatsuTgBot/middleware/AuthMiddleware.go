package middleware

import (
	"context"
	"net/http"
	"rr/kotatsutgbot/config"
	"strconv"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/golang-jwt/jwt/v5"
)

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

func CheckIsMember(userId int64) bool {
	b, err := bot.New(config.GetConfig().CONFIG_BOT_TOKEN)
	if err != nil {
		return false
	}

	user, err := b.GetChatMember(context.TODO(), &bot.GetChatMemberParams{
		ChatID: config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
		UserID: userId,
	})

	if err != nil {
		return false
	}

	if user.Type == models.ChatMemberTypeLeft || user.Type == models.ChatMemberTypeBanned || user.Type == models.ChatMemberTypeRestricted {
		return false
	}
	return true
}

var k, _ = keyfunc.NewDefaultCtx(context.TODO(), []string{"https://oauth.telegram.org/.well-known/jwks.json"})

func ParseAndVerifySession(session string, keyfunc jwt.Keyfunc) (userID int64, isValid bool) {
	data, err := jwt.Parse(session, keyfunc)
	if err != nil {
		return 0, false
	}

	userIDStr := data.Claims.(jwt.MapClaims)["id"].(string)

	userID, err = strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, false
	}

	return userID, true
}

func ExchangeToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization must be Bearer [token]"})
			c.Abort()
			return
		}

		userID, isValid := ParseAndVerifySession(parts[1], k.Keyfunc)
		if !isValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired session",
			})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":  strconv.FormatInt(userID, 10),
			"exp": time.Now().Add(24 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString([]byte(config.GetConfig().AUTH_SECRET))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "failed to create token",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"token": tokenString,
		})
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer [token]"})
			c.Abort()
			return
		}

		userID, isValid := ParseAndVerifySession(parts[1], func(t *jwt.Token) (any, error) {
			return []byte(config.GetConfig().AUTH_SECRET), nil
		})
		if !isValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired session",
			})
			return
		}

		if !CheckIsMember(userID) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "you shall not pass",
			})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
