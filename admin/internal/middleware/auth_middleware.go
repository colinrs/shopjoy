// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"net/http"
	"strings"

	"github.com/colinrs/shopjoy/admin/internal/domain/adminuser"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/httpy"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type Claims struct {
	UserID   int64 `json:"user_id"`
	TenantID int64 `json:"tenant_id"`
	Type     int   `json:"type"`
	jwt.RegisteredClaims
}

func NewAuthMiddleware(jwtSecret string, db *gorm.DB, adminUserRepo adminuser.Repository) rest.Middleware {
	secret := []byte(jwtSecret)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httpy.ResultCtx(r, w, nil, code.ErrUnauthorized)
				return
			}

			// Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				httpy.ResultCtx(r, w, nil, code.ErrTokenInvalid)
				return
			}

			tokenString := parts[1]
			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return secret, nil
			})

			if err != nil || !token.Valid {
				httpy.ResultCtx(r, w, nil, code.ErrTokenExpired)
				return
			}

			// Set user info to context
			ctx := r.Context()
			ctx = contextx.SetUserID(ctx, claims.UserID)
			ctx = contextx.SetTenantID(ctx, claims.TenantID)
			ctx = contextx.SetUserType(ctx, claims.Type)

			// Fetch username from database for audit logging
			if adminUserRepo != nil && db != nil {
				adminUser, err := adminUserRepo.FindByID(ctx, db, claims.UserID)
				if err == nil && adminUser != nil {
					ctx = contextx.SetUserName(ctx, adminUser.Username)
				}
			}

			next(w, r.WithContext(ctx))
		}
	}
}
