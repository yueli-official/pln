package repo

import (
	"pln/models"

	"gorm.io/gorm"
)

type ArtworkRepo interface {
	Create(artwork *models.Artwork) error
	GetByID(id uint) (*models.Artwork, error)
	GetByHash(hash string, artwork *models.Artwork) error
	GetAll(offset, limit int, filters map[string]any) ([]models.Artwork, int64, error)
	GetRandom(limit int, filters map[string]any) ([]models.Artwork, error)
	Update(id uint, artwork *models.Artwork) error
	Delete(id uint) error
	IncrementViews(id uint) error
	IncrementLikes(id uint) error
	DecrementLikes(id uint) error
	IncrementBookmarks(id uint) error
	DecrementBookmarks(id uint) error
}

type artworkRepo struct {
	db *gorm.DB
}

func NewArtworkRepo(db *gorm.DB) ArtworkRepo {
	return &artworkRepo{db: db}
}

func (r *artworkRepo) Create(artwork *models.Artwork) error {
	return r.db.Create(artwork).Error
}

func (r *artworkRepo) GetByID(id uint) (*models.Artwork, error) {
	var artwork models.Artwork
	err := r.db.Where("id = ?", id).First(&artwork).Error
	if err != nil {
		return nil, err
	}
	return &artwork, nil
}

func (r *artworkRepo) GetByHash(hash string, artwork *models.Artwork) error {
	return r.db.Where("hash = ?", hash).First(artwork).Error
}

func (r *artworkRepo) GetAll(offset, limit int, filters map[string]any) ([]models.Artwork, int64, error) {
	var artworks []models.Artwork
	var total int64

	query := r.db.Model(&models.Artwork{})

	if tags, ok := filters["tags"]; ok {
		query = query.Where("tags LIKE ?", "%"+tags.(string)+"%")
	}

	// 计算总数
	if err := query.Model(&models.Artwork{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&artworks).Error; err != nil {
		return nil, 0, err
	}

	return artworks, total, nil
}

func (r *artworkRepo) GetRandom(limit int, filters map[string]any) ([]models.Artwork, error) {
	var artworks []models.Artwork

	query := r.db.Model(&models.Artwork{})

	if tags, ok := filters["tags"]; ok {
		query = query.Where("tags LIKE ?", "%"+tags.(string)+"%")
	}

	// 随机排序并限制数量
	if err := query.Order("RANDOM()").Limit(limit).Find(&artworks).Error; err != nil {
		return nil, err
	}

	return artworks, nil
}

func (r *artworkRepo) Update(id uint, artwork *models.Artwork) error {
	return r.db.Model(&models.Artwork{}).Where("id = ?", id).Updates(artwork).Error
}

func (r *artworkRepo) Delete(id uint) error {
	return r.db.Delete(&models.Artwork{}, id).Error
}

func (r *artworkRepo) IncrementViews(id uint) error {
	return r.db.Model(&models.Artwork{}).Where("id = ?", id).Update("views", gorm.Expr("views + ?", 1)).Error
}

// IncrementLikes 增加点赞数
func (r *artworkRepo) IncrementLikes(id uint) error {
	return r.db.Model(&models.Artwork{}).Where("id = ?", id).Update("likes", gorm.Expr("likes + ?", 1)).Error
}

// DecrementLikes 减少点赞数
func (r *artworkRepo) DecrementLikes(id uint) error {
	return r.db.Model(&models.Artwork{}).Where("id = ?", id).Update("likes", gorm.Expr("CASE WHEN likes > 0 THEN likes - 1 ELSE 0 END")).Error
}

// IncrementBookmarks 增加收藏数
func (r *artworkRepo) IncrementBookmarks(id uint) error {
	return r.db.Model(&models.Artwork{}).Where("id = ?", id).Update("bookmarks", gorm.Expr("bookmarks + ?", 1)).Error
}

// DecrementBookmarks 减少收藏数
func (r *artworkRepo) DecrementBookmarks(id uint) error {
	return r.db.Model(&models.Artwork{}).Where("id = ?", id).Update("bookmarks", gorm.Expr("CASE WHEN bookmarks > 0 THEN bookmarks - 1 ELSE 0 END")).Error
}
