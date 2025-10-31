package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"steam-backend/config"
	"steam-backend/controllers"
	"steam-backend/middleware"
	"steam-backend/repositories"
	"steam-backend/services"
)

func main() {
	cfg := config.LoadConfig()

	db, err := config.InitDb(cfg)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	r := gin.Default()

	userRepo, err_usr := repositories.NewUserRepository(db)
	if err_usr != nil {
		log.Fatalf("Create UserRepository failed: %w", err_usr)
		return
	}
	appRepo, err_app := repositories.NewAppRepository(db)
	if err_app != nil {
		log.Fatalf("Create UserRepository failed: %w", err_app)
		return
	}
	friendRepo, err_friend := repositories.NewFriendRepository(db)
	if err_friend != nil {
		log.Fatalf("Create UserRepository failed: %w", err_friend)
		return
	}
	wishlistRepo, err_wishlist := repositories.NewWishlistRepository(db)
	if err_wishlist != nil {
		log.Fatalf("Create UserRepository failed: %w", err_wishlist)
		return
	}

	userService := services.NewUserService(userRepo, *cfg)
	appService := services.NewAPPService(appRepo)
	friendService := services.NewFriendService(friendRepo)
	wishlistService := services.NewWishlistService(wishlistRepo, appRepo)

	userController := controllers.NewUserController(userService)
	appController := controllers.NewAppController(appService)
	friendController := controllers.NewFriendController(friendService)
	wishlistController := controllers.NewWishlistController(wishlistService)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "backend service is running"})
	})

	api := r.Group("/api")
	{
		userRoutes := api.Group("/user")
		{
			userRoutes.POST("/join", userController.Register)
			userRoutes.POST("/login", userController.Login)
			userRoutes.GET("/available", userController.CheckUsernameAvailable)
			userRoutes.GET("/search", userController.SearchUsers)
			userRoutes.GET("/:id", userController.GetUserByID)

			authUserRoutes := userRoutes.Group("/")
			authUserRoutes.Use(middleware.AuthMiddleware(cfg))
			{
				authUserRoutes.GET("/info", userController.GetUserInfo)
			}

		}

		appRoutes := api.Group("/app")
		{
			appRoutes.GET("/recommendations", appController.GetRecommendations)
			appRoutes.GET("/specials", appController.GetSpecials)
			appRoutes.GET("/search/suggestions", appController.GetSearchSuggestions)
			appRoutes.GET("/:id", appController.GetAppByID)
		}

		friendRoutes := api.Group("/friend")
		friendRoutes.Use(middleware.AuthMiddleware(cfg))
		{
			friendRoutes.GET("/num", friendController.GetFriendCount)
			friendRoutes.GET("/list", friendController.GetFriendList)
			friendRoutes.POST("/invite", friendController.SendInvitation)
			friendRoutes.POST("/invite/accept/:id", friendController.AcceptInvitation)
			friendRoutes.POST("/invite/refuse/:id", friendController.RefuseInvitation)
			friendRoutes.GET("/invite/list/received", friendController.GetReceivedInvitations)
			friendRoutes.GET("/invite/list/sent", friendController.GetSentInvitations)
			friendRoutes.POST("/delete", friendController.RemoveFriend)
			friendRoutes.GET("/check", friendController.CheckFriendship)
		}

		wishlistRoutes := api.Group("/wishlist")
		wishlistRoutes.Use(middleware.AuthMiddleware(cfg))
		{
			wishlistRoutes.GET("/size", wishlistController.GetWishlistSize)
			wishlistRoutes.GET("", wishlistController.GetWishlist)
			wishlistRoutes.POST("", wishlistController.AddToWishlist)
			wishlistRoutes.DELETE("/:appId", wishlistController.RemoveFromWishlist)
			wishlistRoutes.GET("/check", wishlistController.IsInWishlist)
			wishlistRoutes.POST("/sort", wishlistController.SortWishlist)
		}
	}

	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Server starting on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server :%v", err)
	}
}
