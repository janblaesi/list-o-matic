package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	jwt "github.com/appleboy/gin-jwt"
)

func main() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsConfig))

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "List-O-Matic",
		Key:         []byte("very secret"),
		IdentityKey: "id",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					"id":       v.Username,
					"is_admin": v.IsAdmin,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				Username: claims["id"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			// logik hier
			if username == "admin" && password == "admin" {
				return &User{
					Username: username,
					IsAdmin:  true,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
	})

	if err != nil {
		log.Fatal("Failed to initialize JWT subsystem: " + err.Error())
		return
	}

	if err = authMiddleware.MiddlewareInit(); err != nil {
		log.Fatal("Failed to initialize JWT subsystem: " + err.Error())
		return
	}

	router.POST("/login", authMiddleware.LoginHandler)

	protected := router.Group("/protected")
	public := router.Group("/public")
	protected.Use(authMiddleware.MiddlewareFunc())
	{
		setupRoutes(public, protected)
	}

	initDb()

	if err = router.Run(); err != nil {
		log.Fatal("Failed to start Web Server: " + err.Error())
		return
	}
}
