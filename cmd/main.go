package main

import (
	"io"
	"os"
	"strconv"
	"time"

	"pln/conf"
	_ "pln/docs"
	"pln/handler"
	"pln/models"
	"pln/repo"
	"pln/service"

	"github.com/gin-contrib/cors"
	"github.com/rs/zerolog/log"

	"github.com/Yuelioi/gkit/log/zerologx"
	"github.com/Yuelioi/gkit/web/gin/middleware/apikey"
	"github.com/Yuelioi/gkit/web/gin/middleware/log/gzero"
	"github.com/Yuelioi/gkit/web/gin/middleware/ratelimit"
	"github.com/Yuelioi/gkit/web/gin/server"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var configPath = "./config.yaml"

func main() {
	// 初始化日志
	logger := zerologx.Default()
	log.Logger = logger

	// 加载配置
	if err := conf.LoadConfig(configPath); err != nil {
		log.Fatal().Err(err).Msg("配置加载失败")
	}

	// 初始化数据库
	db, err := gorm.Open(sqlite.Open(conf.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("数据库连接失败")
	}

	// 自动迁移
	if err := db.AutoMigrate(models.Artwork{}); err != nil {
		log.Fatal().Err(err).Msg("数据库迁移失败")
	}

	key := conf.InitAPIKey()

	// 禁用默认 Gin 输出
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	// 初始化上传服务和处理器
	uploadService := service.NewFileService(
		conf.Config,
	)
	uploadHandler := handler.NewUploadHandler(uploadService, conf.Config)

	// 初始化仓储、服务和处理器
	artworkRepo := repo.NewArtworkRepo(db)
	artworkService := service.NewArtworkService(artworkRepo)
	artworkHandler := handler.NewArtworkHandler(artworkService, uploadService, conf.Config)

	// 自定义 CORS 配置
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	corsConfig.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"X-API-Key",
		"Authorization",
	}
	corsConfig.ExposeHeaders = []string{"Content-Length", "X-API-Key"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 24 * time.Hour

	// 输出
	logger.Info().Str("app_id", conf.Config.FileServer.AppID).Str("file_server", conf.Config.FileServer.BaseURL).Msg("文件服务注册完毕")

	// 配置默认端口
	port := conf.Config.Server.Port
	if port == 0 {
		port = 9000
	}
	addr := ":" + strconv.Itoa(port)

	// 服务器配置
	cfg := server.ServerConfig{
		Addr:      addr,
		Logger:    logger,
		Mode:      os.Getenv("APP_MODE"),
		APIPrefix: "/api/v1",
		Middlewares: []gin.HandlerFunc{
			gzero.Default(logger),
			gzero.GinRecovery(logger),
			ratelimit.Default(),
			cors.New(corsConfig),
		},
		EnableCORS: false,
		SPAPath:    "./frontend/dist",
	}

	logger.Info().Msg("后台访问密码为: " + key)

	// 启动服务器
	err = server.Start(cfg, func(api *gin.RouterGroup) {

		// 公开路由
		public := api.Group("/")

		public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		{
			public.GET("/artworks", artworkHandler.ListArtworks)
			public.GET("/artworks/random", artworkHandler.RandomArtworks)
			public.GET("/artworks/:id", artworkHandler.GetArtwork)
			public.GET("/artworks/category/:category", artworkHandler.GetArtworksByCategory)

			public.POST("/artworks/:id/like", artworkHandler.IncrementLikes)
			public.POST("/artworks/:id/unlike", artworkHandler.DecrementLikes)
			public.POST("/artworks/:id/bookmark", artworkHandler.IncrementBookmarks)
			public.POST("/artworks/:id/unbookmark", artworkHandler.DecrementBookmarks)

			public.POST("/artworks/upload", artworkHandler.UploadAndCreateArtwork)

			public.POST("/upload", uploadHandler.UploadFile)

		}

		// 需要认证的路由
		auth := api.Group("/", apikey.NewBuilder().WithScheme("X-API-Key").WithValidator(conf.IsValidAPIKey).Handler())
		{
			auth.POST("/artworks", artworkHandler.CreateArtwork)
			auth.PUT("/artworks/:id", artworkHandler.UpdateArtwork)
			auth.DELETE("/artworks/:id", artworkHandler.DeleteArtwork)
		}
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("服务启动失败")
	}
}
