package boards

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *handler) HandleListBoard(ctx *gin.Context) {
	ah, ok := common.GetAuthHeader(ctx)
	if !ok {
		return
	}

	collection := h.mongoDb.Database("photoflux").Collection("boards")
	filter := bson.D{{Key: "user_id", Value: ah.User.String()}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		common.EmitError(ctx, ListBoardError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not get the list of boards: %s", err.Error())))
		return
	}
	defer cursor.Close(ctx)
	var boards []model.Board
	if err = cursor.All(ctx, &boards); err != nil {
		common.EmitError(ctx, ListBoardError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not get the list of boards: %s", err.Error())))
		return
	}

	resp := public.ListBoardsResponse{
		Data: make([]public.BoardData, 0, len(boards)),
	}
	for _, board := range boards {
		resp.Data = append(resp.Data, BoardToItem(board, h.apiPaths))
	}
	ctx.JSON(http.StatusOK, &resp)
}
