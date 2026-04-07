package main

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"pln/conf"
	_ "pln/docs"
	"pln/handler"
	"pln/models"
	"pln/repo"
	"pln/service"
	"pln/storage"

	"github.com/gin-contrib/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Yuelioi/gkit/log/zerologx"
	"github.com/Yuelioi/gkit/log/zerologx/adapter/gormzerolog"
	"github.com/Yuelioi/gkit/web/gin/middleware/apikey"
	"github.com/Yuelioi/gkit/web/gin/middleware/log/gzero"
	"github.com/Yuelioi/gkit/web/gin/middleware/requestid"

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
	zerolog.DefaultContextLogger = &logger

	// 加载配置
	if err := conf.LoadConfig(configPath); err != nil {
		log.Fatal().Err(err).Msg("配置加载失败")
	}

	dbPath := conf.GetDSN()

	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(dbPath), 0755)
		if err != nil {
			panic(err)
		}
	}

	// 初始化数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: gormzerolog.New(),
	})
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

	// 初始化仓储、服务和处理器
	artworkRepo := repo.NewArtworkRepo(db)
	artworkService := service.NewArtworkService(artworkRepo)

	uploader := storage.NewLocalUploader(
		conf.Config.FileServer.StoragePath,
		"/api/v1/files",
		conf.Config.ThumbnailConfig,
		conf.Config.PreviewConfig,
	)

	uploadService := service.NewFileService(
		conf.Config, artworkRepo, uploader,
	)
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

	// 扫描本地目录，导入未入库的图片，并补全已有记录缺失的 URL
	scannedFiles := uploader.ScanFiles()
	imported, updated := 0, 0
	for _, sf := range scannedFiles {
		var existing models.Artwork
		if err := db.Where("file_id = ?", sf.FileID).First(&existing).Error; err == nil {
			// 已存在，检查是否需要补全缺失的 URL
			updates := map[string]any{}
			if existing.ThumbnailURL == "" && sf.ThumbnailURL != "" {
				updates["thumbnail_url"] = sf.ThumbnailURL
			}
			if existing.PreviewURL == "" && sf.PreviewURL != "" {
				updates["preview_url"] = sf.PreviewURL
			}
			if len(updates) > 0 {
				db.Model(&existing).Updates(updates)
				updated++
			}
			continue
		}
		artwork := &models.Artwork{
			FileID:       sf.FileID,
			Hash:         sf.FileID,
			URL:          sf.URL,
			ThumbnailURL: sf.ThumbnailURL,
			PreviewURL:   sf.PreviewURL,
		}
		artwork.SetTags([]string{})
		if err := db.Create(artwork).Error; err != nil {
			logger.Warn().Err(err).Str("file_id", sf.FileID).Msg("导入文件失败")
			continue
		}
		imported++
	}
	if imported > 0 || updated > 0 {
		logger.Info().Int("imported", imported).Int("updated", updated).Int("total_scanned", len(scannedFiles)).Msg("扫描目录完成")
	}

	// 输出
	logger.Info().Str("storage_path", conf.Config.FileServer.StoragePath).Msg("本地文件存储服务注册完毕")

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
			// ratelimit.Default(),
			cors.New(corsConfig),
			requestid.RequestID(),
			gzero.RequestIDMiddleware(),
		},
		EnableCORS: false,
		SPAPath:    "./frontend/dist",
	}

	logger.Info().Msg("后台访问密码为: " + key)

	// 启动服务器
	err = server.Start(cfg, func(api *gin.RouterGroup) {

		// 公开路由
		public := api.Group("/")

		// 本地文件静态服务
		public.Static("/files", conf.Config.FileServer.StoragePath)

		public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		{
			public.GET("/artworks", artworkHandler.ListArtworks)
			public.GET("/artworks/random", artworkHandler.RandomArtworks)
			public.GET("/artworks/:id", artworkHandler.GetArtwork)

			public.POST("/artworks/:id/like", artworkHandler.IncrementLikes)
			public.POST("/artworks/:id/unlike", artworkHandler.DecrementLikes)
			public.POST("/artworks/:id/bookmark", artworkHandler.IncrementBookmarks)
			public.POST("/artworks/:id/unbookmark", artworkHandler.DecrementBookmarks)

			public.POST("/artworks/upload", artworkHandler.UploadAndCreateArtwork)

		}

		// 需要认证的路由
		auth := api.Group("/", apikey.NewBuilder().WithScheme("X-API-Key").WithValidator(conf.IsValidAPIKey).Handler())
		{
			auth.PUT("/artworks/:id", artworkHandler.UpdateArtwork)
			auth.DELETE("/artworks/:id", artworkHandler.DeleteArtwork)
		}
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("服务启动失败")
	}
}
