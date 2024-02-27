package routes

import (
	"Recommendation-System/envconfig"
	api "Recommendation-System/external/api"
	"Recommendation-System/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Register(
	router *gin.Engine,
	logger *zap.Logger,
	db *repository.GormDatabase,
	env *envconfig.Env,
) *gin.Engine {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Create Repo instances
	userRepo := repository.NewUserRepo(db)

	// Register handlers for no authentication API
	// authHandler := NewAuthHandler(logger, userRepo)
	noAuthRouters := router.Group("")
	// noAuthRouters.POST("/api/auth/register", authHandler.register)
	// noAuthRouters.POST("/api/auth/login", authHandler.login)
	// noAuthRouters.POST("/api/auth/refresh", authHandler.refresh)

	// authRouters := router.Group("", authenticate(userRepo, logger))
	// router group to add middle ware for authentication
	userHandler := NewUserHandler(logger, userRepo, api.NewMapUtilities(env.GOOGLE_MAP_API_KEY), api.NewEventsSearcher(env.TICKET_MASTER_API_KEY))

	noAuthRouters.GET("/api/user/events", userHandler.listEvents)
	noAuthRouters.POST("/api/user/events/like", userHandler.likeEvent)
	noAuthRouters.POST("/api/user/events/dislike", userHandler.dislikeEvent)
	noAuthRouters.GET("/api/user/events/recommend", userHandler.recommendEvents)
	return router
}
