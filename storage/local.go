package storage

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pln/conf"
	"pln/models"

	"github.com/nfnt/resize"
	"github.com/rs/zerolog/log"
)

// variantDef defines a variant to generate (thumbnail, preview, etc.)
type variantDef struct {
	suffix  string // e.g. "_thumbnail", "_preview"
	width   uint
	height  uint
	quality int
	enabled bool
}

// LocalUploader implements Uploader interface for local file storage
type LocalUploader struct {
	storagePath string
	urlPrefix   string // URL prefix for accessing files, e.g. "/api/v1/files"
	variants    []variantDef
}

func NewLocalUploader(storagePath, urlPrefix string, thumbnail, preview conf.ThumbnailOption) *LocalUploader {
	os.MkdirAll(storagePath, 0755)

	var variants []variantDef
	if thumbnail.Enabled {
		variants = append(variants, variantDef{
			suffix: "_thumbnail", width: uint(thumbnail.Width), height: uint(thumbnail.Height),
			quality: thumbnail.Quality, enabled: true,
		})
	}
	if preview.Enabled {
		variants = append(variants, variantDef{
			suffix: "_preview", width: uint(preview.Width), height: uint(preview.Height),
			quality: preview.Quality, enabled: true,
		})
	}

	return &LocalUploader{
		storagePath: storagePath,
		urlPrefix:   urlPrefix,
		variants:    variants,
	}
}

func (l *LocalUploader) Upload(ctx context.Context, file io.Reader, filename string, options map[string]any) (*UploadResponse, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	hash := sha256.Sum256(data)
	id := hex.EncodeToString(hash[:])
	newFilename := id + ext

	// Save original
	origPath := filepath.Join(l.storagePath, newFilename)
	if err := os.WriteFile(origPath, data, 0644); err != nil {
		return nil, fmt.Errorf("写入文件失败: %w", err)
	}

	// Generate all variants
	l.ensureVariants(id, ext, data)

	return &UploadResponse{
		FileID: id,
		JobID:  id,
		URL:    l.urlPrefix + "/" + newFilename,
		Status: "completed",
	}, nil
}

// ensureVariants generates missing variants for a given image
func (l *LocalUploader) ensureVariants(id, ext string, data []byte) {
	for _, v := range l.variants {
		outExt := ext
		if outExt == ".webp" {
			outExt = ".jpg"
		}
		variantPath := filepath.Join(l.storagePath, id+v.suffix+outExt)

		// Skip if already exists
		if _, err := os.Stat(variantPath); err == nil {
			continue
		}

		if err := generateResized(bytes.NewReader(data), variantPath, v.width, v.height, v.quality, outExt); err != nil {
			log.Warn().Err(err).Str("variant", v.suffix).Str("file_id", id).Msg("生成变体失败，跳过")
		}
	}
}

func generateResized(r io.Reader, outPath string, width, height uint, quality int, ext string) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return fmt.Errorf("解码图片失败: %w", err)
	}

	resized := resize.Thumbnail(width, height, img, resize.Lanczos3)

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	switch ext {
	case ".png":
		return png.Encode(out, resized)
	case ".gif":
		return gif.Encode(out, resized, nil)
	default:
		return jpeg.Encode(out, resized, &jpeg.Options{Quality: quality})
	}
}

func (l *LocalUploader) Delete(ctx context.Context, fileID string) error {
	// Delete all files matching this fileID (original + all variants)
	matches, _ := filepath.Glob(filepath.Join(l.storagePath, fileID+"*"))
	if len(matches) == 0 {
		return fmt.Errorf("文件不存在")
	}
	for _, m := range matches {
		os.Remove(m)
	}
	return nil
}

func (l *LocalUploader) GetJobProgress(ctx context.Context, jobID string) (*JobProgressResponse, error) {
	return &JobProgressResponse{
		JobID:          jobID,
		Status:         "task.completed",
		TotalTasks:     1,
		CompletedTasks: 1,
		CreatedAt:      time.Now().Unix(),
	}, nil
}

func (l *LocalUploader) GetFileInfo(fileID string) (*models.FileInfo, error) {
	origPath := l.findOriginal(fileID)
	if origPath == "" {
		return nil, fmt.Errorf("文件不存在: %s", fileID)
	}

	ext := filepath.Ext(origPath)
	filename := filepath.Base(origPath)

	stat, err := os.Stat(origPath)
	if err != nil {
		return nil, err
	}

	info := &models.FileInfo{
		Name:         filename,
		OriginalName: filename,
		Size:         stat.Size(),
		Path:         filename,
		AccessURL:    l.urlPrefix + "/" + filename,
		StorageType:  "local",
		FileType:     models.FileTypeImage,
		FileID:       fileID,
		UploadTime:   stat.ModTime().Unix(),
	}

	// Check all variants
	variantExt := ext
	if variantExt == ".webp" {
		variantExt = ".jpg"
	}
	for _, v := range l.variants {
		vFilename := fileID + v.suffix + variantExt
		vPath := filepath.Join(l.storagePath, vFilename)
		if _, err := os.Stat(vPath); err == nil {
			info.Variants = append(info.Variants, models.ResourceVariant{
				Type:      strings.TrimPrefix(v.suffix, "_"),
				Path:      vFilename,
				AccessURL: l.urlPrefix + "/" + vFilename,
				CreatedAt: time.Now(),
			})
		}
	}

	return info, nil
}

// findOriginal finds the original file for a fileID (excluding variant files)
func (l *LocalUploader) findOriginal(fileID string) string {
	matches, _ := filepath.Glob(filepath.Join(l.storagePath, fileID+".*"))
	for _, m := range matches {
		base := filepath.Base(m)
		if !strings.Contains(base, "_thumbnail") && !strings.Contains(base, "_preview") {
			return m
		}
	}
	return ""
}

// ScannedFile represents a file found during directory scan
type ScannedFile struct {
	FileID       string
	URL          string
	ThumbnailURL string
	PreviewURL   string
}

// ScanFiles scans the storage directory, auto-generates missing variants,
// and returns all original image files for DB import.
func (l *LocalUploader) ScanFiles() []ScannedFile {
	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	}

	entries, err := os.ReadDir(l.storagePath)
	if err != nil {
		log.Warn().Err(err).Msg("扫描存储目录失败")
		return nil
	}

	var files []ScannedFile
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.Contains(name, "_thumbnail") || strings.Contains(name, "_preview") {
			continue
		}

		ext := strings.ToLower(filepath.Ext(name))
		if !imageExts[ext] {
			continue
		}

		fileID := strings.TrimSuffix(name, ext)

		// Auto-generate missing variants
		origPath := filepath.Join(l.storagePath, name)
		data, err := os.ReadFile(origPath)
		if err != nil {
			log.Warn().Err(err).Str("file", name).Msg("读取文件失败，跳过")
			continue
		}
		l.ensureVariants(fileID, ext, data)

		// Build ScannedFile
		sf := ScannedFile{
			FileID: fileID,
			URL:    l.urlPrefix + "/" + name,
		}

		variantExt := ext
		if variantExt == ".webp" {
			variantExt = ".jpg"
		}
		thumbName := fileID + "_thumbnail" + variantExt
		if _, err := os.Stat(filepath.Join(l.storagePath, thumbName)); err == nil {
			sf.ThumbnailURL = l.urlPrefix + "/" + thumbName
		}
		previewName := fileID + "_preview" + variantExt
		if _, err := os.Stat(filepath.Join(l.storagePath, previewName)); err == nil {
			sf.PreviewURL = l.urlPrefix + "/" + previewName
		}

		files = append(files, sf)
	}

	return files
}
