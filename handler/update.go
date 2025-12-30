package handler

import (
	"pln/models"
	"strconv"

	"github.com/Yuelioi/gkit/web/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

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
