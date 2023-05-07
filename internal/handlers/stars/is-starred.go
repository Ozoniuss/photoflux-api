package stars

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleIsPhotoStarred(ctx *gin.Context) {
	var req public.StarRequest
	err := ctx.BindJSON(&req)

	if err != nil {
		common.EmitError(ctx, IsPhotoStarredError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
		return
	}
	isStar := false
	var star []model.Star
	err = h.db.WithContext(ctx).Clauses(clause.Returning{}).Table("stars").
		Where("photo_id = ?", req.PhotoId).
		Where("user_id = ?", req.UserId).
		Find(&star).
		Error

	if err != nil {
		common.EmitError(ctx, IsPhotoStarredError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not get star data: %s", err.Error())))
		return
	}

	if len(star) == 1 {
		// Star didn't exist, therefore we create it
		isStar = true
	}

	resp := public.StarResponse{
		IsStar: isStar,
	}

	ctx.JSON(http.StatusCreated, &resp)
}
