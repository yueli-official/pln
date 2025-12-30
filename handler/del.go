package handler

import (
	"strconv"

	"github.com/Yuelioi/gkit/web/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// DeleteArtwork 删除作品
// @Summary 删除作品
// @Description 删除指定ID的作品
// @Tags Artwork
// @Param id path int true "作品ID"
// @Success 204
// @Router /artworks/{id} [delete]
func (h *ArtworkHandler) DeleteArtwork(c *gin.Context) {
	//  解析 artwork ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest("invalid artwork id").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	//  获取作品

	// 删除本地文件
	if _, err := h.fileService.DeleteFile(c.Request.Context(), uint(id)); err != nil {
		log.Error().Err(err).Msg("删除本地文件失败")
		response.InternalError("删除本地文件失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	//  删除数据库记录
	if err := h.service.DeleteArtwork(uint(id)); err != nil {
		log.Error().Err(err).Msg("删除作品失败")
		response.InternalError("删除作品失败").
			WithRequestID(c.GetString("request_id")).
			GJSON(c)
		return
	}

	// 6️⃣ 返回 204 No Content
	response.NoContent().
		WithRequestID(c.GetString("request_id")).
		GJSON(c)
}
