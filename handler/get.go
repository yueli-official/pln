package handler

import (
	"strconv"

	"github.com/Yuelioi/gkit/web/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

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

// RandomArtworks 随机获取作品
// @Summary 随机获取作品
// @Description 随机获取指定数量的作品
// @Tags Artwork
// @Produce json
// @Param limit query int false "数量" default(10)
// @Success 200 {object} response.Response{data=[]models.ArtworkResponse} "获取成功"
// @Router /artworks/random [get]
func (h *ArtworkHandler) RandomArtworks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// 限制最大数量为 100
	if limit < 1 || limit > 100 {
		limit = 10
	}

	filters := make(map[string]any)
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
