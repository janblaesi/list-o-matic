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

// List-o-matic is a talking list management system
package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Start by loading the configuration, fail if it does not exist
	if err := cfgLoad(); err != nil {
		log.Fatal("Failed to load configuration file.")
	}

	// Then try to load the pseudo database of talking lists
	if err := setupDatabase(); err != nil {
		log.Print("Could not load an existing database, creating a new one.")
	}

	// Release mode, comment out this line when developing
	gin.SetMode(gin.ReleaseMode)

	// Setup gin-gonic Library
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Modern browsers use CORS preflighting for requests to ensure higher
	// security.
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AddAllowHeaders("Content-Type")
	router.Use(cors.New(corsConfig))

	// Setup authentication middleware
	if err := authSetup(); err != nil {
		log.Fatal("Setting up authentication subsystem failed.")
	}
	router.POST("/login", authMiddleware.LoginHandler)

	protected := router.Group("/protected")
	public := router.Group("/public")
	protected.Use(authMiddleware.MiddlewareFunc())
	{
		setupRoutes(public, protected)
	}

	if err := router.Run(); err != nil {
		log.Fatal("Failed to start web server")
	}
}
