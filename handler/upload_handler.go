package handler

import (
	"path/filepath"

	"pln/conf"
	"pln/models"
	"pln/service"

	"github.com/Yuelioi/gkit/web/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UploadHandler struct {
	fileService *service.FileService
	cfg         *conf.AppConfig
}

func NewUploadHandler(fileService *service.FileService, cfg *conf.AppConfig) *UploadHandler {
	return &UploadHandler{
		fileService: fileService,
		cfg:         cfg,
	}
}

// @Summary 上传文件到 CDN
// @Description 上传图片或文件到 CDN
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "要上传的文件"
// @Success 200 {object} response.Response{data=models.FileInfo} "上传成功"
// @Router /files [post]
func (h *UploadHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest("没有找到文件").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	// 验证文件大小（默认限制 5MB）
	const maxFileSize = 5 * 1024 * 1024
	if file.Size > maxFileSize {
		response.BadRequest("文件大小超过限制（最大 5MB）").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	// 验证文件类型
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".svg":  true,
		".pdf":  true,
		".zip":  true,
		".mp4":  true,
		".webm": true,
	}

	ext := filepath.Ext(file.Filename)
	if !allowedExts[ext] {
		response.BadRequest("不支持的文件类型").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	// 打开上传文件
	src, err := file.Open()
	if err != nil {
		log.Error().Err(err).Msg("打开上传文件失败")
		response.InternalError("上传失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}
	defer src.Close()

	// 上传到三方服务
	fileInfo, err := h.fileService.UploadFile(src, file.Filename, h.cfg.FileServer.CDNSpaceID)
	if err != nil {
		log.Error().Err(err).Msg("上传到三方服务失败")
		response.InternalError(err.Error()).
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	response.OK().WithData(fileInfo).
		WithRequestID(c.GetString("request_id")).
		GJSON(c)
}

// Delete 删除文件
// @Summary 删除文件
// @Description 从三方服务删除文件
// @Tags Upload
// @Produce json
// @Param path query string true "文件路径"
// @Success 200 {object} response.Response{data=models.DeleteFileResponse} "删除成功"
// @Router /files [delete]
func (h *UploadHandler) Delete(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		response.BadRequest("文件路径不能为空").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	deleted, err := h.fileService.DeleteFile(
		h.cfg.FileServer.AppID,
		h.cfg.FileServer.CDNSpaceID,
		path,
	)
	if err != nil {
		log.Error().Err(err).Msg("删除文件失败")
		response.InternalError(err.Error()).
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	result := models.DeleteFileResponse{
		Path:    path,
		Deleted: deleted,
	}

	response.OK().WithData(result).
		WithRequestID(c.GetString("request_id")).
		GJSON(c)
}

// RegisterRoutes 注册路由
func (h *UploadHandler) RegisterRoutes(api *gin.RouterGroup) {
	files := api.Group("/files")
	{
		files.POST("", h.UploadFile)
		files.DELETE("", h.Delete)
	}
}
