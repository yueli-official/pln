package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"pln/models"
	"time"

	"github.com/Yuelioi/gkit/web/response"
)

// Third Party API Responses
type UploadResponse struct {
	FileID      string `json:"file_id"`
	JobID       string `json:"job_id"`
	URL         string `json:"url"`
	StorageType string `json:"storage_type"`
	Message     string `json:"message"`
	Status      string `json:"status"`
	StatusURL   string `json:"status_url"`
}

type JobProgressResponse struct {
	JobID          string `json:"job_id"`
	Status         string `json:"status"`
	TotalTasks     int    `json:"total_tasks"`
	CompletedTasks int    `json:"completed_tasks"`
	FailedTasks    int    `json:"failed_tasks"`
	ErrorMsg       string `json:"error_msg,omitempty"`
	CreatedAt      int64  `json:"created_at"`
}

type DeleteResponse struct {
	FileID  string `json:"file_id"`
	Deleted bool   `json:"deleted"`
}

// Uploader defines the interface for third-party storage
type Uploader interface {
	Upload(ctx context.Context, file io.Reader, filename string, options map[string]any) (*UploadResponse, error)
	Delete(ctx context.Context, fileID string) error
	GetJobProgress(ctx context.Context, jobID string) (*JobProgressResponse, error)
	GetFileInfo(fileID string) (*models.FileInfo, error)
}

// ThirdPartyUploader implements Uploader interface
type ThirdPartyUploader struct {
	baseURL    string
	appKey     string
	appID      string
	spaceID    string
	httpClient *http.Client
}

// NewThirdPartyUploader creates a new third-party uploader instance
func NewThirdPartyUploader(baseURL, appID, spaceID, appKey string) *ThirdPartyUploader {
	return &ThirdPartyUploader{
		baseURL: baseURL,
		appID:   appID,
		spaceID: spaceID,
		appKey:  appKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Upload uploads file to third-party storage
func (t *ThirdPartyUploader) Upload(ctx context.Context, file io.Reader, filename string, options map[string]any) (*UploadResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	if err := writer.WriteField("app_id", t.appID); err != nil {
		return nil, fmt.Errorf("写入app_id失败: %w", err)
	}

	if err := writer.WriteField("space_id", t.spaceID); err != nil {
		return nil, fmt.Errorf("写入space_id失败: %w", err)
	}

	// Add options if provided
	if len(options) > 0 {
		optionsJSON, err := json.Marshal(options)
		if err != nil {
			return nil, fmt.Errorf("序列化options失败: %w", err)
		}
		if err := writer.WriteField("options", string(optionsJSON)); err != nil {
			return nil, fmt.Errorf("写入options失败: %w", err)
		}
	}

	// Add file field
	fw, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("创建表单文件失败: %w", err)
	}

	if _, err := io.Copy(fw, file); err != nil {
		return nil, fmt.Errorf("复制文件内容失败: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("关闭writer失败: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/v1/files", t.baseURL), body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("上传失败，状态码: %d，响应: %s", resp.StatusCode, string(respBody))
	}

	uploadResp, err := DecodeResponseBody[UploadResponse](resp.Body)

	return uploadResp, nil
}

// Delete deletes a file from third-party storage by file_id
func (t *ThirdPartyUploader) Delete(ctx context.Context, fileID string) error {
	url := fmt.Sprintf("%s/api/v1/files/%s", t.baseURL, fileID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("创建删除请求失败: %w", err)
	}

	req.Header.Set("X-App-ID", t.appID)
	req.Header.Set("X-API-Key", t.appKey)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return errors.New("文件不存在")
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("删除失败，状态码: %d，响应: %s", resp.StatusCode, string(respBody))
	}

	deleteResp, err := DecodeResponseBody[DeleteResponse](resp.Body)
	if err != nil {
		return fmt.Errorf("解析删除响应失败: %w", err)
	}

	if !deleteResp.Deleted {
		return errors.New("文件不存在")
	}

	return nil
}

// GetJobProgress queries the upload progress by job_id
// This polls the StatusURL returned from Upload response
func (t *ThirdPartyUploader) GetJobProgress(ctx context.Context, jobID string) (*JobProgressResponse, error) {
	url := fmt.Sprintf("%s/api/v1/jobs/%s", t.baseURL, jobID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建进度查询请求失败: %w", err)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("查询任务进度失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("查询进度失败，状态码: %d，响应: %s", resp.StatusCode, string(respBody))
	}

	progressResp, err := DecodeResponseBody[JobProgressResponse](resp.Body)

	return progressResp, nil
}

// GetFileInfo 获取文件信息
func (t *ThirdPartyUploader) GetFileInfo(path string) (*models.FileInfo, error) {

	url := fmt.Sprintf("%s/api/v1/files/%s",
		t.baseURL, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return DecodeResponseBody[models.FileInfo](resp.Body)

}

// Helper function to convert map to options for easy usage

func DecodeResponseBody[T any](body io.Reader) (*T, error) {
	var resp response.Response

	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("请求失败: %s", resp.Message)
	}

	dataBytes, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("序列化 Data 失败: %w", err)
	}

	var result T
	if err := json.Unmarshal(dataBytes, &result); err != nil {
		return nil, fmt.Errorf("反序列化 Data 失败: %w", err)
	}

	return &result, nil
}
