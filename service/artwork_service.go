package service

import (
	"errors"
	"pln/models"
	"pln/repo"

	"gorm.io/gorm"
)

type ArtworkService interface {
	CreateArtwork(req *models.ArtworkCreateRequest) (*models.ArtworkResponse, error)
	GetArtwork(id uint) (*models.ArtworkResponse, error)
	GetByHash(hash string) (*models.Artwork, error)
	GetArtworks(page, pageSize int, filters map[string]any) ([]models.ArtworkResponse, int64, error)
	GetRandomArtworks(limit int, filters map[string]any) ([]models.ArtworkResponse, error)

	UpdateArtwork(id uint, req *models.ArtworkUpdateRequest) (*models.ArtworkResponse, error)
	DeleteArtwork(id uint) error

	IncrementViews(id uint) error
	IncrementLikes(id uint) error
	DecrementLikes(id uint) error
	IncrementBookmarks(id uint) error
	DecrementBookmarks(id uint) error
}

type artworkService struct {
	repo repo.ArtworkRepo
}

func NewArtworkService(repo repo.ArtworkRepo) ArtworkService {
	return &artworkService{repo: repo}
}

func (s *artworkService) CreateArtwork(req *models.ArtworkCreateRequest) (*models.ArtworkResponse, error) {
	artwork := &models.Artwork{
		URL:          req.URL,
		ThumbnailURL: req.ThumbnailURL,
		Hash:         req.Hash,
		FileID:       req.FileID,
		Views:        0,
		Likes:        0,
		Bookmarks:    0,
	}

	// 设置 tags，如果为空则设置为空数组
	tags := req.Tags
	if len(tags) == 0 {
		tags = []string{}
	}
	if err := artwork.SetTags(tags); err != nil {
		return nil, err
	}

	if err := s.repo.Create(artwork); err != nil {
		return nil, err
	}

	resp := artwork.ToResponse()
	return &resp, nil
}

func (s *artworkService) GetArtwork(id uint) (*models.ArtworkResponse, error) {
	artwork, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("artwork not found")
		}
		return nil, err
	}

	// 增加浏览次数
	_ = s.repo.IncrementViews(id)

	resp := artwork.ToResponse()
	return &resp, nil
}

func (s *artworkService) GetByHash(hash string) (*models.Artwork, error) {
	var artwork models.Artwork
	if err := s.repo.GetByHash(hash, &artwork); err != nil {
		return nil, err
	}
	return &artwork, nil
}

func (s *artworkService) GetArtworks(page, pageSize int, filters map[string]any) ([]models.ArtworkResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	artworks, total, err := s.repo.GetAll(offset, pageSize, filters)
	if err != nil {
		return nil, 0, err
	}

	var responses []models.ArtworkResponse
	for _, artwork := range artworks {
		responses = append(responses, artwork.ToResponse())
	}

	return responses, total, nil
}

func (s *artworkService) GetRandomArtworks(limit int, filters map[string]any) ([]models.ArtworkResponse, error) {
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	artworks, err := s.repo.GetRandom(limit, filters)
	if err != nil {
		return nil, err
	}

	var responses []models.ArtworkResponse
	for _, artwork := range artworks {
		responses = append(responses, artwork.ToResponse())
	}

	return responses, nil
}

func (s *artworkService) UpdateArtwork(id uint, req *models.ArtworkUpdateRequest) (*models.ArtworkResponse, error) {
	artwork := &models.Artwork{}

	if req.URL != "" {
		artwork.URL = req.URL
	}
	if len(req.Tags) > 0 {
		artwork.SetTags(req.Tags)
	}

	if err := s.repo.Update(id, artwork); err != nil {
		return nil, err
	}

	// 重新获取更新后的数据
	updated, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	resp := updated.ToResponse()
	return &resp, nil
}

func (s *artworkService) DeleteArtwork(id uint) error {
	return s.repo.Delete(id)
}

func (s *artworkService) IncrementViews(id uint) error {
	return s.repo.IncrementViews(id)
}

func (s *artworkService) IncrementLikes(id uint) error {
	artwork, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if artwork == nil {
		return errors.New("artwork not found")
	}

	return s.repo.IncrementLikes(id)
}

func (s *artworkService) DecrementLikes(id uint) error {
	artwork, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if artwork == nil {
		return errors.New("artwork not found")
	}

	return s.repo.DecrementLikes(id)
}

func (s *artworkService) IncrementBookmarks(id uint) error {
	artwork, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if artwork == nil {
		return errors.New("artwork not found")
	}

	return s.repo.IncrementBookmarks(id)
}

func (s *artworkService) DecrementBookmarks(id uint) error {
	artwork, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if artwork == nil {
		return errors.New("artwork not found")
	}

	return s.repo.DecrementBookmarks(id)
}
