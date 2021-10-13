// List-O-Matic Talking List Management System
// Copyright (C) 2021 Jan Blaesi
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
