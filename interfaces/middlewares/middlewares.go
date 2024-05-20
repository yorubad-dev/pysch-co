package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/KingDaemonX/ddd-template/domain/repository/infrastructures/auth"
	"github.com/KingDaemonX/ddd-template/domain/repository/interfaces/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Middleware struct {
	slg     *slog.Logger
	ratelimter *rate.Limiter
}

func NewMiddleware(slg *slog.Logger) *Middleware {
	return &Middleware{
		slg:     slg,
		ratelimter: rate.NewLimiter(10, 5),
	}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.ValidateToken(c.Request, m.slg)
		if err != nil {
			m.slg.Error("Unauthorized access", "error context", "Invalid token", "function", "AuthMiddleware")
			response := response.NewResponse(http.StatusUnauthorized, "Unauthorized access", "Invalid token provided")
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}
		c.Next()
	}
}

func (m *Middleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			m.slg.Error("Preflight request", "error context", "Preflight request detected", "function", "CORSMiddleware")
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func (m *Middleware) PriviledgeCheckMiddleware(td auth.TokenInterface, rd auth.RedisInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user has the right to access the resource
		uid := c.Param("uid")

		ad, err := td.ExtractMetadata(c.Request)
		if err != nil {
			m.slg.Error("Error extracting metadata", "error context", err, "function", "PriviledgeCheckMiddleware")
			response := response.NewResponse(http.StatusUnauthorized, "Unauthorized access", err.Error())
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		// , err := rd.FetchAuth(ad.TokenUuid)
		if uid != ad.UID {
			m.slg.Error("Unauthorized access", "error context", "possible security breach and CRSF in place on user with the following uid", "uid", ad.UID, "function", "PriviledgeCheckMiddleware")
			response := response.NewResponse(http.StatusUnauthorized, "Unauthorized access", "You are not allowed to access this resource")
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		if uid == ad.UID && ad.Role != "admin" {
			m.slg.Error("Unauthorized access", "error context", "User is not allowed to access this resource", "function", "PriviledgeCheckMiddleware")
			response := response.NewResponse(http.StatusUnauthorized, "Unauthorized access", "You are not allowed to access this resource")
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		c.Next()
	}
}

func (m *Middleware) RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.ratelimter.Allow() {
			m.slg.Error("Rate limit exceeded", "error context", "Rate limit exceeded", "function", "RateLimiter")
			response := response.NewResponse(http.StatusTooManyRequests, "Rate limit exceeded", "You have exceeded the rate limit")
			c.JSON(http.StatusTooManyRequests, response)
			c.Abort()
			return
		}
		c.Next()
	}
}
