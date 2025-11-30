package response

import "PetTrack/infra/02-handler/request"

func GetPageResponse(req request.PageInfo, count, pages int64) map[string]interface{} {
	return map[string]interface{}{
		"page":       req.Page,
		"size":       req.Size,
		"total":      count,
		"totalPages": pages,
	}
}
