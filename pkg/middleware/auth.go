package middleware

import (
	"net/http"
	"strings"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	repo   *repository.Repository
	Logger *zap.Logger
}

type contextKey string

const authUserKey contextKey = "authUser"

type AuthUser struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Name   string
	Role   string
}

func NewAuthMiddleware(repo *repository.Repository, log *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		repo:   repo,
		Logger: log,
	}
}

func GetAuthUser(c *gin.Context) (*AuthUser, bool) {
	user, ok := c.Get(authUserKey)
	return user.(*AuthUser), ok
}

func (m *AuthMiddleware) SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// get session from Authorization
		authHeader := c.GetHeader("Authorization")

		// validae
		if authHeader == "" {
			utils.ResponseFailed(c, http.StatusUnauthorized, "unauthorized", "token is empty please login")
			c.Abort()
			return
		}

		// validate token type
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ResponseFailed(c, http.StatusUnauthorized, "unauthorized", "token type invalid")
			c.Abort()
			return
		}

		// parse string to uuid
		sessionID, err := uuid.Parse(parts[1])
		if err != nil {
			utils.ResponseFailed(c, http.StatusUnauthorized, "unauthorized", "token invalid")
			c.Abort()
			return
		}

		// validate token
		sess, err := m.repo.AuthRepository.ValidateSession(ctx, sessionID)
		if err != nil {
			m.Logger.Error("failed to validate session", zap.Error(err))
			utils.ResponseFailed(c, http.StatusUnauthorized, "unauthorized", "token invalid or inactive")
			c.Abort()
			return
		}

		authUser := &AuthUser{
			ID:     sess.SessionID,
			UserID: sess.UserID,
			Role:   sess.Role,
		}

		c.Set(authUserKey, authUser)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		user, ok := GetAuthUser(c)
		if !ok {
			utils.ResponseFailed(c, http.StatusUnauthorized, "unauthorized", nil)
			c.Abort()
			return
		}

		if _, exists := allowed[user.Role]; !exists {
			utils.ResponseFailed(c, http.StatusForbidden, "forbidden", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

func (m *AuthMiddleware) CheckPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := GetAuthUser(c)
		if !ok {
			utils.ResponseFailed(c, http.StatusUnauthorized, "unauthorized", nil)
			c.Abort()
			return
		}

		userData, err := m.repo.UserRepository.GetUserByID(c.Request.Context(), user.UserID)
		if err != nil {
			m.Logger.Error("failed to get user by id", zap.Error(err))
			utils.ResponseFailed(c, http.StatusInternalServerError, "internal server error", nil)
			c.Abort()
			return
		}

		allowed, ok := userData.Permissions[permission]
		if !ok || !allowed {
			utils.ResponseFailed(c, http.StatusForbidden, "forbidden to access "+permission, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
