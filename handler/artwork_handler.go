package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"path/filepath"
	"time"

	"pln/conf"
	"pln/models"
	"pln/service"

	"github.com/Yuelioi/gkit/web/response"
	"github.com/corona10/goimagehash"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ArtworkHandler struct {
	service     service.ArtworkService
	fileService *service.FileService
	cfg         *conf.AppConfig
}

func NewArtworkHandler(service service.ArtworkService, fileService *service.FileService, cfg *conf.AppConfig) *ArtworkHandler {
	return &ArtworkHandler{service: service, fileService: fileService, cfg: cfg}
}

// @Summary 上传文件并创建作品
// @Description 上传图片到 CDN 并同时创建艺术作品记录
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "要上传的文件"
// @Success 201 {object} response.Response{data=models.ArtworkResponse}
// @Router /artworks/upload [post]
func (h *ArtworkHandler) UploadAndCreateArtwork(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := c.GetString("request_id")
	logger := log.Ctx(ctx).With().Str("component", "ArtworkHandler").Logger()
	logger.Info().Msg("开始上传作品")

	file, err := c.FormFile("file")
	if err != nil {
		logger.Warn().Err(err).Msg("没有找到文件")
		response.BadRequest("没有找到文件").
			WithRequestID(requestID).
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
		logger.Warn().Str("filename", file.Filename).Msg("不支持的文件格式")
		response.BadRequest("只支持图片格式").
			WithRequestID(requestID).
			GJSON(c)
		return
	}

	// ============ 步骤 0：计算文件 Hash（提前检测重复）============
	src0, err := file.Open()
	if err != nil {
		logger.Error().Err(err).Msg("打开上传文件失败")
		response.InternalError("上传失败").
			WithRequestID(requestID).
			GJSON(c)
		return
	}
	defer src0.Close()

	// 计算 MD5 Hash
	hash, err := h.calculateFileHash(src0)
	if err != nil {
		logger.Error().Err(err).Msg("计算文件 Hash 失败")
		response.InternalError("计算文件 Hash 失败").
			WithRequestID(requestID).
			GJSON(c)
		return
	}

	logger.Debug().Str("hash", hash).Msg("文件 Hash 计算完成")

	// 检查 MD5 Hash 是否已存在
	existingArtwork, err := h.service.GetByHash(hash)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error().Err(err).Msg("查询 Hash 失败")
		response.InternalError("查询失败").
			WithRequestID(requestID).
			GJSON(c)
		return
	}

	// Hash 已存在，返回冲突错误（409 Conflict）
	if existingArtwork != nil {
		logger.Info().Str("hash", hash).Uint("artwork_id", existingArtwork.ID).Msg("文件已存在")
		response.Conflict("图片已存在").
			WithRequestID(requestID).
			GJSON(c)
		return
	}

	// ============ 步骤 0.5：计算 pHash（感知哈希，检测相似图片）============
	src0_2, err := file.Open()
	if err != nil {
		logger.Error().Err(err).Msg("打开上传文件失败")
		response.InternalError("上传失败").
			WithRequestID(requestID).
			GJSON(c)
		return
	}
	defer src0_2.Close()

	pHash, err := h.calculatePHash(src0_2)
	if err != nil {
		logger.Warn().Err(err).Msg("计算 pHash 失败，继续处理")
		// pHash 失败不中断流程，只记录日志
	} else {
		logger.Debug().Int64("phash", pHash).Msg("pHash 计算完成")

		// 查询相似的图片（使用汉明距离）
		similarArtworks, err := h.service.GetByPHashSimilarity(pHash, 5) // 汉明距离阈值为 5
		if err != nil {
			logger.Warn().Err(err).Msg("查询相似图片失败")
		} else if len(similarArtworks) > 0 {
			logger.Info().
				Int64("phash", pHash).
				Int("similar_count", len(similarArtworks)).
				Uint("similar_artwork_id", similarArtworks[0].ID).
				Msg("发现相似图片")
			response.Conflict("图片过于相似，已存在").
				WithRequestID(requestID).
				GJSON(c)
			return
		}
	}

	// ============ 步骤 1：上传本地存储（异步）============
	src, err := file.Open()
	if err != nil {
		logger.Error().Err(err).Msg("打开上传文件失败")
		response.InternalError("上传失败").
			WithRequestID(requestID).
			GJSON(c)
		return
	}
	defer src.Close()

	logger.Debug().Msg("开始上传到本地存储")
	localUploadResp, err := h.fileService.UploadLocalFile(ctx, src, file.Filename)
	if err != nil {
		logger.Error().Err(err).Msg("上传到本地存储失败")
		response.InternalError("上传到本地存储失败").
			WithRequestID(requestID).
			GJSON(c)
		return
	}

	// ============ 步骤 3：返回异步上传信息，客户端需轮询查询状态 ============

	// 保存待处理的上传任务（稍后异步完成后入库）
	uploadTask := &models.UploadTask{
		Hash:              hash,
		PHash:             pHash, // 保存 pHash
		FileID:            localUploadResp.FileID,
		JobID:             localUploadResp.JobID,
		StatusURL:         localUploadResp.StatusURL,
		Status:            models.UploadStatusPending,
		CreatedAt:         time.Now().Unix(),
		LastStatusCheckAt: time.Now().Unix(),
	}

	if err := h.PollUploadJobStatus(ctx, uploadTask, requestID); err != nil {
		logger.Error().Err(err).Msg("保存上传任务失败")
		response.InternalError("保存上传任务失败").
			WithRequestID(requestID).
			GJSON(c)
		return
	}

	logger.Info().Str("upload_id", uploadTask.ID).Msg("上传任务已创建，等待异步完成")

	info, err := h.fileService.GetFileInfo(localUploadResp.FileID)
	if err != nil {
		response.BusinessError("获取文件信息失败").GJSON(nil)
		return
	}
	thumbnail := ""

	for _, v := range info.Variants {
		if v.Type == "thumbnail" {
			thumbnail = v.AccessURL
		}
	}

	artworkResp, err := h.service.CreateArtwork(&models.ArtworkCreateRequest{
		FileID:       localUploadResp.FileID,
		URL:          h.cfg.FileServer.BaseURL + info.AccessURL,
		Hash:         hash,
		PHash:        pHash, // 保存 pHash
		ThumbnailURL: h.cfg.FileServer.BaseURL + thumbnail,
		Tags:         []string{},
	})

	if err != nil {
		response.InternalError("创建条目失败").GJSON(c)
		return
	}

	response.OK().WithData(artworkResp).GJSON(c)
}

// cleanupLocalFile 清理本地上传的文件
func (h *ArtworkHandler) cleanupLocalFile(ctx context.Context, localFileID string) {
	if localFileID == "" {
		return
	}

	logger := log.Ctx(ctx).With().Str("file_id", localFileID).Logger()
	if _, err := h.fileService.DeleteFileByFileID(ctx, localFileID); err != nil {
		logger.Error().Err(err).Msg("清理本地文件失败")
	}
}

// cleanupUploadedFiles 清理本地和CDN上传的文件
func (h *ArtworkHandler) cleanupUploadedFiles(ctx context.Context, localFileID, cdnFileID string) {
	logger := log.Ctx(ctx)

	if localFileID != "" {
		if _, err := h.fileService.DeleteFileByFileID(ctx, localFileID); err != nil {
			logger.Error().Err(err).Str("file_id", localFileID).Msg("清理本地文件失败")
		}
	}

	if cdnFileID != "" {
		if _, err := h.fileService.DeleteFileByFileID(ctx, cdnFileID); err != nil {
			logger.Error().Err(err).Str("file_id", cdnFileID).Msg("清理CDN文件失败")
		}
	}
}

func (h *ArtworkHandler) calculateFileHash(file io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (h *ArtworkHandler) calculatePHash(src io.Reader) (int64, error) {
	// 解码图片
	img, _, err := image.Decode(src)
	if err != nil {
		return 0, fmt.Errorf("解码图片失败: %w", err)
	}

	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return 0, fmt.Errorf("计算 pHash 失败: %w", err)
	}

	return int64(hash.GetHash()), nil
}

func (h *ArtworkHandler) PollUploadJobStatus(ctx context.Context, uploadTask *models.UploadTask, requestID string) error {
	logger := log.Ctx(ctx).With().Str("component", "ArtworkHandler").Str("upload_id", uploadTask.ID).Logger()

	maxRetries := 20
	retryInterval := 1 * time.Second
	retryCount := 0

	for retryCount < maxRetries {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 获取本地存储和CDN上传进度
		localProgress, localErr := h.fileService.GetJobProgress(ctx, uploadTask.JobID)

		// 判断是否有错误
		if localErr != nil && !errors.Is(localErr, context.DeadlineExceeded) {
			logger.Error().Err(localErr).Msg("获取本地存储进度失败")
			return fmt.Errorf("获取本地存储进度失败: %w", localErr)
		}

		// 判断两个任务是否都完成

		if localProgress != nil && localProgress.Status == "task.completed" {
			return nil
		} else {
			logger.Info().Str("progress", localProgress.Status)

		}

		// 未完成，等待后重试
		retryCount++
		logger.Debug().Int("retry_count", retryCount).Msg("等待中...")
		time.Sleep(retryInterval)
	}

	logger.Error().Int("max_retries", maxRetries).Msg("轮询超时")

	return fmt.Errorf("上传任务轮询超时，超过 %d 次尝试", maxRetries)
}
