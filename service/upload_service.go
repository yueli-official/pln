package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"pln/conf"
	"pln/models"

	"github.com/rs/zerolog/log"
)

type FileService struct {
	cfg *conf.AppConfig
	cli *http.Client
}

func NewFileService(cfg *conf.AppConfig) *FileService {
	return &FileService{
		cfg: cfg,
		cli: &http.Client{},
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
func (fs *FileService) UploadFile(file io.Reader, fileName string, spaceID string) (*models.FileInfo, error) {
	url := fs.cfg.FileServer.BaseURL + "/api/v1/files"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加表单字段
	writer.WriteField("app_id", fs.cfg.FileServer.AppID)
	writer.WriteField("space_id", spaceID)

	log.Info().
		Str("app_id", fs.cfg.FileServer.AppID).
		Str("space_id", spaceID).
		Str("file_name", fileName).
		Msg("上传文件")

	// 添加缩略图配置
	options := UploadOptions{
		Thumbnail: fs.cfg.ThumbnailConfig,
	}
	optionsJSON, _ := json.Marshal(options)
	writer.WriteField("options", string(optionsJSON))

	// 添加文件
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		log.Error().Err(err).Str("file_name", fileName).Msg("创建表单字段失败")
		return nil, fmt.Errorf("添加文件失败: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		log.Error().Err(err).Str("file_name", fileName).Msg("复制文件内容失败")
		return nil, fmt.Errorf("复制文件内容失败: %w", err)
	}

	writer.Close()

	// 发送请求
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("创建请求失败")
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-App-ID", fs.cfg.FileServer.AppID)
	req.Header.Set("X-API-Key", fs.cfg.FileServer.APIKey)

	// 添加超时设置
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := fs.cli.Do(req)
	if err != nil {
		log.Error().
			Err(err).
			Str("url", url).
			Str("file_name", fileName).
			Str("space_id", spaceID).
			Msg("上传请求失败")

		if errors.Is(err, context.DeadlineExceeded) {
			log.Error().Msg("上传超时")
		}
		return nil, fmt.Errorf("上传失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取完整响应体，避免 EOF 错误
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().
			Err(err).
			Str("file_name", fileName).
			Msg("读取响应失败")
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	// 先检查状态码
	if resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status_code", resp.StatusCode).
			Str("response", string(respData)).
			Str("file_name", fileName).
			Str("space_id", spaceID).
			Msg("服务器返回错误")
		return nil, fmt.Errorf("上传状态码错误 %d: %s", resp.StatusCode, string(respData))
	}

	//  解析 JSON 响应
	var result ThirdPartyResponse
	if err := json.Unmarshal(respData, &result); err != nil {
		log.Error().
			Err(err).
			Str("response", string(respData)).
			Str("file_name", fileName).
			Msg("解析响应 JSON 失败")
		return nil, fmt.Errorf("解析响应 JSON 失败: %w", err)
	}

	// 检查业务逻辑错误
	if result.Code != 0 {
		log.Error().
			Int("code", result.Code).
			Str("message", result.Message).
			Str("file_name", fileName).
			Str("space_id", spaceID).
			Msg("服务器业务错误")
		return nil, fmt.Errorf("上传失败: code=%d, message=%s", result.Code, result.Message)
	}

	// 检查返回数据是否为空
	if result.Data == nil {
		log.Error().
			Str("file_name", fileName).
			Str("space_id", spaceID).
			Msg("服务器返回数据为空")
		return nil, fmt.Errorf("服务器返回数据为空")
	}

	// 解析返回的 FileInfo
	fileInfoData, _ := json.Marshal(result.Data)
	var fileInfo models.FileInfo
	if err := json.Unmarshal(fileInfoData, &fileInfo); err != nil {
		log.Error().
			Err(err).
			Str("file_name", fileName).
			Msg("解析文件信息失败")
		return nil, fmt.Errorf("解析文件信息失败: %w", err)
	}

	// 构建完整的访问 URL
	fileInfo.AccessURL = fs.cfg.FileServer.BaseURL + fileInfo.AccessURL

	// 处理缩略图 URL
	for i, v := range fileInfo.Variants {
		if v.Type == "thumbnail" {
			fileInfo.Variants[i].AccessURL = fs.cfg.FileServer.BaseURL + v.AccessURL
		}
	}

	log.Info().
		Str("file_name", fileName).
		Str("space_id", spaceID).
		Str("access_url", fileInfo.AccessURL).
		Msg("文件上传成功")

	return &fileInfo, nil
}

// ============ 删除文件 ============

// 删除文件
func (fs *FileService) DeleteFile(appID string, spaceID string, path string) (bool, error) {
	url := fmt.Sprintf("%s?app_id=%s&space_id=%s&path=%s",
		fs.cfg.FileServer.BaseURL+"/api/v1/files", appID, spaceID, path)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return false, err
	}

	log.Info().
		Str("app_id", fs.cfg.FileServer.AppID).
		Str("space_id", spaceID).
		Str("path", path).
		Msg("删除文件")

	req.Header.Set("X-App-ID", appID)
	req.Header.Set("X-API-Key", fs.cfg.FileServer.APIKey)

	resp, err := fs.cli.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result ThirdPartyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	if result.Code != 0 {
		return false, fmt.Errorf("delete failed: %s", result.Message)
	}

	// 解析删除响应
	deleteResp, _ := json.Marshal(result.Data)
	var deleteResult models.DeleteFileResponse
	if err := json.Unmarshal(deleteResp, &deleteResult); err != nil {
		return false, err
	}

	return deleteResult.Deleted, nil
}

// ============ 获取文件信息 ============

// GetFileInfo 获取文件信息
func (fs *FileService) GetFileInfo(appID string, spaceID string, path string) (*models.FileInfo, error) {
	url := fmt.Sprintf("%s/api/v1/files/info?app_id=%s&space_id=%s&path=%s",
		fs.cfg.FileServer.BaseURL, appID, spaceID, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("app_id", fs.cfg.FileServer.AppID).
		Str("space_id", spaceID).
		Str("path", path).
		Msg("获取文件信息")

	req.Header.Set("X-App-ID", appID)
	req.Header.Set("X-API-Key", fs.cfg.FileServer.APIKey)

	resp, err := fs.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ThirdPartyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("get file info failed: %s", result.Message)
	}

	fileInfoData, _ := json.Marshal(result.Data)
	var fileInfo models.FileInfo
	if err := json.Unmarshal(fileInfoData, &fileInfo); err != nil {
		return nil, err
	}

	return &fileInfo, nil
}
