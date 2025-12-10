package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"pln/conf"
	"pln/models"
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

// UploadFile 上传文件到三方服务
func (fs *FileService) UploadFile(file io.Reader, fileName string, spaceId string) (*models.FileInfo, error) {
	url := fs.cfg.FileServer.UploadURL

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加表单字段
	writer.WriteField("app_id", fs.cfg.FileServer.AppID)
	writer.WriteField("space_id", spaceId)

	// 添加缩略图配置
	options := UploadOptions{
		Thumbnail: fs.cfg.ThumbnailConfig,
	}
	optionsJSON, _ := json.Marshal(options)
	writer.WriteField("options", string(optionsJSON))

	// 添加文件
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}
	writer.Close()

	// 发送请求
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-App-ID", fs.cfg.FileServer.AppID)
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
		return nil, fmt.Errorf("upload failed: %s", result.Message)
	}

	// 解析返回的 FileInfo
	fileInfoData, _ := json.Marshal(result.Data)
	var fileInfo models.FileInfo
	if err := json.Unmarshal(fileInfoData, &fileInfo); err != nil {
		return nil, err
	}

	fileInfo.AccessURL = fs.cfg.FileServer.BaseURL + fileInfo.AccessURL

	for i, v := range fileInfo.Variants {
		if v.Type == "thumbnail" {
			fileInfo.Variants[i].AccessURL = fs.cfg.FileServer.BaseURL + v.AccessURL
		}
	}

	return &fileInfo, nil
}

// ============ 删除文件 ============

// DeleteFile 删除文件从三方服务
func (fs *FileService) DeleteFile(appID string, spaceID string, path string) (bool, error) {
	url := fmt.Sprintf("%s?app_id=%s&space_id=%s&path=%s",
		fs.cfg.FileServer.DeleteURL, appID, spaceID, path)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return false, err
	}

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
