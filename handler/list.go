package handler

import (
	"strconv"

	"github.com/Yuelioi/gkit/web/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ListArtworks 获取作品列表
// @Summary 获取作品列表
// @Description 分页获取作品列表，支持过滤
// @Tags Artwork
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=[]models.ArtworkResponse} "获取成功"
// @Router /artworks [get]
func (h *ArtworkHandler) ListArtworks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filters := make(map[string]any)

	if tags := c.QueryArray("tags"); len(tags) > 0 {
		filters["tags"] = tags
	}

	artworks, total, err := h.service.GetArtworks(page, pageSize, filters)
	if err != nil {
		log.Error().Err(err).Msg("获取作品列表失败")
		response.InternalError("获取作品列表失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	response.Page(artworks, total, page, pageSize).
		WithRequestID(c.GetString("request_id")).
		GJSON(c)
}
