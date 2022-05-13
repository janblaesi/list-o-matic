//     __    _      __        ____        __  ___      __  _
//    / /   (_)____/ /_      / __ \      /  |/  /___ _/ /_(_)____
//   / /   / / ___/ __/_____/ / / /_____/ /|_/ / __ `/ __/ / ___/
//  / /___/ (__  ) /_/_____/ /_/ /_____/ /  / / /_/ / /_/ / /__
// /_____/_/____/\__/      \____/     /_/  /_/\__,_/\__/_/\___/
//
// Copyright 2021-2022 Jan Blaesi
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files
// (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO
// THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
// CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// User represents a user that may log in to the system
type User struct {
	// The name of the user
	Username string `json:"username"`

	// The SHA-256 password hash of the user
	PasswordHash string `json:"password_hash"`

	// Flag, if the user is an admin
	IsAdmin bool `json:"is_admin"`
}

// Login represents arguments provided in a login form
type Login struct {
	// The username provided in the login form
	Username string `form:"username" json:"username" binding:"required"`

	// The clear-text password provided in the login form
	Password string `form:"password" json:"password" binding:"required"`
}

// The database of users
var users []User

// The middleware used for authentication
var authMiddleware *jwt.GinJWTMiddleware

func authSetup() error {
	var err error

	if err := parseJsonFromFile(&users, cfg.Database.UsersPath); err != nil {
		return err
	}

	// For authentication we use JSON Web Tokens, in this case the implementation by GitHub user appleboy
	authMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "List-O-Matic",
		Key:         []byte(cfg.Authentication.Secret),
		Timeout:     time.Duration(cfg.Authentication.TimeoutSeconds) * time.Second,
		IdentityKey: "id",

		// The PayloadFunc dumps the claims into the JSON Web Token
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					"id":       v.Username,
					"is_admin": v.IsAdmin,
				}
			}
			return jwt.MapClaims{}
		},

		// The IdentityHandler extracts the claims from the JSON Web Token and returns
		// the currently logged in user
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			for _, user := range users {
				if user.Username == claims["id"].(string) {
					return &user
				}
			}

			return nil
		},

		// The Authenticator is called when a user tries to log in
		// This will check if the user provided correct credentials
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			// Iterate over all users and check the password hash if the correct user was found
			for _, user := range users {
				if user.Username == loginVals.Username {
					hasher := sha256.New()
					hasher.Write([]byte(loginVals.Password))
					if user.PasswordHash == hex.EncodeToString(hasher.Sum(nil)) {
						return &user, nil
					}
				}
			}

			return nil, jwt.ErrFailedAuthentication
		},
	})

	if err != nil {
		return err
	}

	if err = authMiddleware.MiddlewareInit(); err != nil {
		return err
	}

	return nil
}
