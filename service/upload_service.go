package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"pln/conf"
	"pln/models"
	"pln/repo"
	"pln/storage"

	"github.com/rs/zerolog/log"
)

type FileService struct {
	cfg      *conf.AppConfig
	cli      *http.Client
	repo     repo.ArtworkRepo
	uploader storage.Uploader
}

func NewFileService(cfg *conf.AppConfig, repo repo.ArtworkRepo, uploader storage.Uploader) *FileService {
	return &FileService{
		cfg:      cfg,
		cli:      &http.Client{Timeout: 30 * time.Second},
		repo:     repo,
		uploader: uploader,
	}
}

// ============ API 响应结构 ============

type ThirdPartyResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type UploadOptions struct {
	Thumbnail conf.ThumbnailOption `json:"thumbnail"`
}

// ============ 上传文件 ============

// UploadLocalFile 上传文件到本地存储空间
func (fs *FileService) UploadLocalFile(ctx context.Context, file io.Reader, fileName string) (*storage.UploadResponse, error) {
	logger := log.Ctx(ctx).With().
		Str("component", "FileService").
		Str("file_name", fileName).
		Str("space_type", "local").
		Logger()

	logger.Debug().Msg("开始上传到本地存储")

	options := buildThumbnailOptions(fs.cfg.ThumbnailConfig.Width, fs.cfg.ThumbnailConfig.Height, fs.cfg.ThumbnailConfig.Mode, fs.cfg.ThumbnailConfig.Quality)

	resp, err := fs.uploader.Upload(context.Background(), file, fileName, options)

	if err != nil {
		logger.Error().Err(err).Msg("上传到本地存储失败")
		return nil, fmt.Errorf("上传到本地存储失败: %w", err)
	}

	logger.Debug().
		Str("file_id", resp.FileID).
		Str("status_url", resp.StatusURL).
		Msg("本地文件上传成功")

	return resp, nil
}

func buildThumbnailOptions(width, height int, mode string, quality int) map[string]any {
	return map[string]any{
		"thumbnail": map[string]any{
			"enabled": true,
			"width":   width,
			"height":  height,
			"mode":    mode,
			"quality": quality,
		},
	}
}

// ============ 删除文件 ============

// DeleteFile 删除文件（通过 artworkID 查询并删除相关文件）
func (fs *FileService) DeleteFile(ctx context.Context, artworkID uint) (bool, error) {
	logger := log.Ctx(ctx).With().
		Str("app_id", fs.cfg.FileServer.AppID).
		Uint("artwork_id", artworkID).
		Str("component", "FileService").
		Logger()

	logger.Info().Msg("开始删除文件")

	// Get artwork record
	artwork, err := fs.repo.GetByID(artworkID)
	if err != nil {
		logger.Error().Err(err).Msg("获取文件信息失败")
		return false, fmt.Errorf("获取文件信息失败: %w", err)
	}

	if artwork == nil {
		logger.Warn().Msg("文件记录不存在")
		return false, fmt.Errorf("文件记录不存在")
	}

	// Delete from third-party storage
	err = fs.uploader.Delete(ctx, artwork.FileID)
	if err != nil {
		if strings.Contains(err.Error(), "文件不存在") {
			logger.Warn().Err(err).Msg("文件不存在")
			// File doesn't exist on remote, but we still delete the DB record
			if err := fs.repo.Delete(artworkID); err != nil {
				logger.Error().Err(err).Msg("删除数据库记录失败")
				return false, fmt.Errorf("删除数据库记录失败: %w", err)
			}
			return false, nil
		}

		logger.Error().Err(err).Msg("删除文件失败")
		return false, fmt.Errorf("删除文件失败: %w", err)
	}

	logger.Debug().Msg("远程文件删除成功")

	return true, nil
}

// DeleteFileByFileID 按文件ID直接删除文件（不查询数据库）
func (fs *FileService) DeleteFileByFileID(ctx context.Context, fileID string) (bool, error) {
	logger := log.Ctx(ctx).With().
		Str("file_id", fileID).
		Str("component", "FileService").
		Logger()

	logger.Debug().Msg("开始删除文件")

	if fileID == "" {
		logger.Warn().Msg("文件ID不能为空")
		return false, fmt.Errorf("文件ID不能为空")
	}

	// Delete from third-party storage using uploader
	err := fs.uploader.Delete(ctx, fileID)
	if err != nil {
		// Check if error is because file doesn't exist
		if strings.Contains(err.Error(), "文件不存在") {
			logger.Warn().Msg("文件不存在")
			return false, nil
		}

		logger.Error().Err(err).Msg("删除文件失败")
		return false, fmt.Errorf("删除文件失败: %w", err)
	}

	logger.Debug().Msg("文件删除成功")
	return true, nil
}

// ============ 获取文件信息 ============

// GetFileInfo 获取文件信息
func (fs *FileService) GetFileInfo(fileID string) (*models.FileInfo, error) {

	log.Info().
		Str("file_id", fileID).
		Msg("获取文件信息")

	return fs.uploader.GetFileInfo(fileID)

}

//  查询上传进度
func (fs *FileService) GetJobProgress(ctx context.Context, jobID string) (*storage.JobProgressResponse, error) {
	logger := log.Ctx(ctx).With().
		Str("job_id", jobID).
		Str("component", "FileService").
		Logger()

	logger.Debug().Msg("开始查询上传进度")

	if jobID == "" {
		return nil, fmt.Errorf("任务ID不能为空")
	}

	progress, err := fs.uploader.GetJobProgress(ctx, jobID)
	if err != nil {
		logger.Error().Err(err).Msg("查询上传进度失败")
		return nil, fmt.Errorf("查询上传进度失败: %w", err)
	}

	logger.Debug().
		Int("total_tasks", progress.TotalTasks).
		Int("completed_tasks", progress.CompletedTasks).
		Str("status", progress.Status).
		Msg("获取上传进度成功")

	return progress, nil
}
