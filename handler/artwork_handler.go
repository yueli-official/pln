package handler

import (
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"pln/conf"
	"pln/models"
	"pln/service"

	"github.com/Yuelioi/gkit/web/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ArtworkHandler struct {
	service     service.ArtworkService
	fileService *service.FileService
	cfg         *conf.AppConfig
}

func NewArtworkHandler(service service.ArtworkService, fileService *service.FileService, cfg *conf.AppConfig) *ArtworkHandler {
	return &ArtworkHandler{service: service, fileService: fileService, cfg: cfg}
}

// CreateArtwork 创建作品
// @Summary 创建作品
// @Description 创建新的艺术作品
// @Tags Artwork
// @Accept json
// @Produce json
// @Param body body models.ArtworkCreateRequest true "作品信息"
// @Success 201 {object} response.Response{data=models.ArtworkResponse} "创建成功"
// @Router /artworks [post]
func (h *ArtworkHandler) CreateArtwork(c *gin.Context) {
	var req models.ArtworkCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(err.Error()).
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	artwork, err := h.service.CreateArtwork(&req)
	if err != nil {
		log.Error().Err(err).Msg("创建作品失败")
		response.InternalError("创建作品失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	response.Created(artwork).
		WithRequestID(c.GetString("request_id")).
		GJSON(c)
}

// GetArtwork 获取单个作品
// @Summary 获取作品详情
// @Description 获取指定ID的作品详情
// @Tags Artwork
// @Produce json
// @Param id path int true "作品ID"
// @Success 200 {object} response.Response{data=models.ArtworkResponse} "获取成功"
// @Router /artworks/{id} [get]
func (h *ArtworkHandler) GetArtwork(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest("invalid artwork id").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	artwork, err := h.service.GetArtwork(uint(id))
	if err != nil {
		response.NotFound("artwork not found").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	response.OK().WithData(artwork).
		WithRequestID(c.GetString("request_id")).
		GJSON(c)
}

// ListArtworks 获取作品列表
// @Summary 获取作品列表
// @Description 分页获取作品列表，支持过滤
// @Tags Artwork
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param category query string false "分类"
// @Param artist query string false "艺术家"
// @Success 200 {object} response.Response{data=[]models.ArtworkResponse} "获取成功"
// @Router /artworks [get]
func (h *ArtworkHandler) ListArtworks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filters := make(map[string]interface{})
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if artist := c.Query("artist"); artist != "" {
		filters["artist"] = artist
	}
	if tags := c.QueryArray("tags"); len(tags) > 0 {
		filters["tags"] = tags
	}

	artworks, total, err := h.service.GetArtworks(page, pageSize, filters)
	if err != nil {
		log.Error().Err(err).Msg("获取作品列表失败")
		response.InternalError("获取作品列表失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	response.Page(artworks, total, page, pageSize).
		WithRequestID(c.GetString("request_id")).
		GJSON(c)
}

// RandomArtworks 随机获取作品
// @Summary 随机获取作品
// @Description 随机获取指定数量的作品
// @Tags Artwork
// @Produce json
// @Param limit query int false "数量" default(10)
// @Param category query string false "分类"
// @Param artist query string false "艺术家"
// @Success 200 {object} response.Response{data=[]models.ArtworkResponse} "获取成功"
// @Router /artworks/random [get]
func (h *ArtworkHandler) RandomArtworks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// 限制最大数量为 100
	if limit < 1 || limit > 100 {
		limit = 10
	}

	filters := make(map[string]interface{})
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if artist := c.Query("artist"); artist != "" {
		filters["artist"] = artist
	}
	if tags := c.QueryArray("tags"); len(tags) > 0 {
		filters["tags"] = tags
	}

	artworks, err := h.service.GetRandomArtworks(limit, filters)
	if err != nil {
		log.Error().Err(err).Msg("获取随机作品失败")
		response.InternalError("获取随机作品失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	response.OK().WithData(artworks).
		WithRequestID(c.GetString("request_id")).
		GJSON(c)
}

// UpdateArtwork 更新作品
// @Summary 更新作品
// @Description 更新指定ID的作品信息
// @Tags Artwork
// @Accept json
// @Produce json
// @Param id path int true "作品ID"
// @Param body body models.ArtworkUpdateRequest true "作品信息"
// @Success 200 {object} response.Response{data=models.ArtworkResponse} "更新成功"
// @Router /artworks/{id} [put]
func (h *ArtworkHandler) UpdateArtwork(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest("invalid artwork id").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	var req models.ArtworkUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(err.Error()).
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	artwork, err := h.service.UpdateArtwork(uint(id), &req)
	if err != nil {
		log.Error().Err(err).Msg("更新作品失败")
		response.InternalError("更新作品失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	response.OK().WithData(artwork).
		WithRequestID(c.GetString("request_id")).
		GJSON(c)
}

// DeleteArtwork 删除作品
// @Summary 删除作品
// @Description 删除指定ID的作品
// @Tags Artwork
// @Param id path int true "作品ID"
// @Success 204
// @Router /artworks/{id} [delete]
func (h *ArtworkHandler) DeleteArtwork(c *gin.Context) {
	//  解析 artwork ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest("invalid artwork id").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	//  获取作品
	artwork, err := h.service.GetArtwork(uint(id))
	if err != nil {
		log.Error().Err(err).Msg("获取作品失败")
		response.NotFound("artwork not found").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	// 删除本地文件
	if _, err := h.fileService.DeleteFile(h.cfg.FileServer.AppID, h.cfg.FileServer.LocalSpaceID, artwork.LocalPath); err != nil {
		log.Error().Err(err).Msg("删除本地文件失败")
		response.InternalError("删除本地文件失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	// 删除云端文件
	if _, err := h.fileService.DeleteFile(h.cfg.FileServer.AppID, h.cfg.FileServer.CDNSpaceID, artwork.CDNPath); err != nil {
		log.Error().Err(err).Msg("删除云端文件失败")
		response.InternalError("删除云端文件失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	//  删除数据库记录
	if err := h.service.DeleteArtwork(uint(id)); err != nil {
		log.Error().Err(err).Msg("删除作品失败")
		response.InternalError("删除作品失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	// 6️⃣ 返回 204 No Content
	response.NoContent().
		WithRequestID(c.GetString("request_id")).
		GJSON(c)
}

// GetArtworksByCategory 按分类获取作品
// @Summary 按分类获取作品
// @Description 获取指定分类的作品列表
// @Tags Artwork
// @Produce json
// @Param category path string true "分类"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /artworks/category/{category} [get]
func (h *ArtworkHandler) GetArtworksByCategory(c *gin.Context) {
	category := c.Param("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	artworks, total, err := h.service.GetArtworksByCategory(category, page, pageSize)
	if err != nil {
		log.Error().Err(err).Msg("获取分类作品失败")
		response.InternalError("获取分类作品失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	resp := gin.H{
		"data":      artworks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	response.OK().WithData(resp).
		WithRequestID(c.GetString("request_id")).
		GJSON(c)
}

// IncrementLikes 增加点赞数
// @Summary 增加作品点赞数
// @Description 增加作品的点赞数
// @Tags Artwork
// @Produce json
// @Param id path int true "作品ID"
// @Success 200 {object} response.Response "操作成功"
// @Router /artworks/{id}/like [post]
func (h *ArtworkHandler) IncrementLikes(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest("invalid artwork id").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	if err := h.service.IncrementLikes(uint(id)); err != nil {
		log.Error().Err(err).Msg("点赞失败")
		response.InternalError("点赞失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	response.OK().
		WithRequestID(c.GetString("request_id")).
		GJSON(c)
}

// DecrementLikes 减少点赞数
// @Summary 取消点赞
// @Description 减少作品的点赞数
// @Tags Artwork
// @Produce json
// @Param id path int true "作品ID"
// @Success 200 {object} response.Response "操作成功"
// @Router /artworks/{id}/unlike [post]
func (h *ArtworkHandler) DecrementLikes(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest("invalid artwork id").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	if err := h.service.DecrementLikes(uint(id)); err != nil {
		log.Error().Err(err).Msg("取消点赞失败")
		response.InternalError("取消点赞失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	response.OK().
		WithRequestID(c.GetString("request_id")).
		GJSON(c)
}

// IncrementBookmarks 增加收藏数
// @Summary 增加收藏数
// @Description 增加作品的收藏数
// @Tags Artwork
// @Produce json
// @Param id path int true "作品ID"
// @Success 200 {object} response.Response "操作成功"
// @Router /artworks/{id}/bookmark [post]
func (h *ArtworkHandler) IncrementBookmarks(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest("invalid artwork id").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	if err := h.service.IncrementBookmarks(uint(id)); err != nil {
		log.Error().Err(err).Msg("收藏失败")
		response.InternalError("收藏失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	response.OK().
		WithRequestID(c.GetString("request_id")).
		GJSON(c)
}

// DecrementBookmarks 减少收藏数
// @Summary 取消收藏
// @Description 减少作品的收藏数
// @Tags Artwork
// @Produce json
// @Param id path int true "作品ID"
// @Success 200 {object} response.Response "操作成功"
// @Router /artworks/{id}/unbookmark [post]
func (h *ArtworkHandler) DecrementBookmarks(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest("invalid artwork id").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	if err := h.service.DecrementBookmarks(uint(id)); err != nil {
		log.Error().Err(err).Msg("取消收藏失败")
		response.InternalError("取消收藏失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	response.OK().
		WithRequestID(c.GetString("request_id")).
		GJSON(c)
}

// UploadAndCreateArtwork 上传文件并创建作品
// @Summary 上传文件并创建作品
// @Description 上传图片到 CDN 并同时创建艺术作品记录
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "要上传的文件"
// @Param title formData string true "作品标题"
// @Param artist formData string true "艺术家名称"
// @Param description formData string false "作品描述"
// @Param category formData string false "分类"
// @Param avatar_url formData string false "头像URL"
// @Success 201 {object} map[string]interface{}
// @Router /artworks/upload [post]
func (h *ArtworkHandler) UploadAndCreateArtwork(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest("没有找到文件").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	// 验证必填字段
	title := c.PostForm("title")
	artist := c.PostForm("artist")

	if title == "" || artist == "" {
		response.BadRequest("标题和艺术家名称不能为空").
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
	}

	ext := filepath.Ext(file.Filename)
	if !allowedExts[ext] {
		response.BadRequest("只支持图片格式").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	// 打开原始文件
	src, err := file.Open()
	if err != nil {
		log.Error().Err(err).Msg("打开上传文件失败")
		response.InternalError("上传失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}
	defer src.Close()

	// 上传原图到本地存储
	localFileInfo, err := h.fileService.UploadFile(src, file.Filename, h.cfg.FileServer.LocalSpaceID)
	if err != nil {
		log.Error().Err(err).Msg("上传到本地存储失败")
		response.InternalError("上传到本地存储失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	// 重新打开文件进行压缩和CDN上传
	src2, err := file.Open()
	if err != nil {
		log.Error().Err(err).Msg("打开上传文件失败")
		response.InternalError("上传失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}
	defer src2.Close()

	// 如果文件超过5MB，需要压缩并转为JPG
	const maxFileSize = 5 * 1024 * 1024
	var fileToUpload io.Reader
	var uploadFilename string

	if file.Size > maxFileSize {
		compressedFile, err := h.compressImageToJPEG(src2)
		if err != nil {
			log.Error().Err(err).Msg("图片压缩失败")
			response.InternalError("图片压缩失败").
				WithRequestID(c.GetString("request_id")).
				GJSON(c)
			return
		}
		defer compressedFile.Close()

		fileToUpload = compressedFile
		// 生成新的文件名，转为JPG格式
		uploadFilename = fmt.Sprintf("%s_compressed.jpg", strings.TrimSuffix(file.Filename, ext))
	} else {
		src3, err := file.Open()
		if err != nil {
			log.Error().Err(err).Msg("打开上传文件失败")
			response.InternalError("上传失败").
				WithRequestID(c.GetString("request_id")).
				GJSON(c)
			return
		}
		defer src3.Close()
		fileToUpload = src3
		uploadFilename = file.Filename
	}

	// 上传到 CDN
	fileInfo, err := h.fileService.UploadFile(fileToUpload, uploadFilename, h.cfg.FileServer.CDNSpaceID)
	if err != nil {
		log.Error().Err(err).Msg("上传到 CDN 失败")
		response.InternalError("上传到 CDN 失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	// 提取缩略图URL
	thumbnailURL := ""
	if fileInfo.Variants != nil {
		for _, variant := range fileInfo.Variants {
			if variant.Type == "thumbnail" {
				thumbnailURL = variant.AccessURL
				break
			}
		}
	}

	// 创建作品
	description := c.PostForm("description")
	category := c.PostForm("category")
	avatarURL := c.PostForm("avatar_url")

	artworkReq := &models.ArtworkCreateRequest{
		Title:        title,
		Description:  description,
		Artist:       artist,
		AvatarURL:    avatarURL,
		URL:          fileInfo.AccessURL,
		LocalPath:    localFileInfo.Path,
		CDNPath:      fileInfo.Path,
		CDNURL:       fileInfo.AccessURL,
		ThumbnailURL: thumbnailURL,
		Tags:         c.PostFormArray("tags[]"),
		Category:     category,
		IsPublished:  true,
		Hash:         fileInfo.Metadata.Hash,
	}

	artwork, err := h.service.CreateArtwork(artworkReq)
	if err != nil {
		log.Error().Err(err).Msg("创建作品失败")
		response.InternalError("创建作品失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	resp := gin.H{
		"artwork": artwork,
		"upload":  fileInfo,
	}

	response.Created(resp).
		WithRequestID(c.GetString("request_id")).
		GJSON(c)
}

// compressImageToJPEG 压缩图片为JPG格式，确保不超过5MB
func (h *ArtworkHandler) compressImageToJPEG(src io.Reader) (*os.File, error) {
	// 读取原始图片
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("解码图片失败: %w", err)
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "compressed_*.jpg")
	if err != nil {
		return nil, fmt.Errorf("创建临时文件失败: %w", err)
	}

	// 如果压缩失败需要删除临时文件
	defer func() {
		if err != nil {
			os.Remove(tmpFile.Name())
		}
	}()

	// 逐步降低质量直到文件大小满足要求，初始质量设为90
	quality := 90
	const maxFileSize = 5 * 1024 * 1024

	for quality >= 50 {
		tmpFile.Seek(0, 0)
		tmpFile.Truncate(0)

		err = jpeg.Encode(tmpFile, img, &jpeg.Options{Quality: quality})
		if err != nil {
			return nil, fmt.Errorf("编码图片失败: %w", err)
		}

		// 检查文件大小
		fi, err := tmpFile.Stat()
		if err != nil {
			return nil, fmt.Errorf("获取文件信息失败: %w", err)
		}

		if fi.Size() <= maxFileSize {
			// 大小满足要求，重置文件指针
			tmpFile.Seek(0, 0)
			return tmpFile, nil
		}

		quality -= 10
	}

	return nil, fmt.Errorf("无法将图片压缩到5MB以内")
}
