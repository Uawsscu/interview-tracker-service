package routers

import (
	"interview-tracker/internal/adapters/handlers"
	"interview-tracker/internal/adapters/repositories"
	"interview-tracker/internal/config"
	"interview-tracker/internal/middleware"
	"interview-tracker/internal/usecases"

	"github.com/gin-gonic/gin"
)

func Card(r *gin.RouterGroup) {
	repo := repositories.NewCardRepo(config.DB)
	uc := usecases.NewCardUsecase(repo)
	h := handlers.NewCardHandler(uc)

	// ทั้งหมดอยู่ใต้ /interview-tracker/authen (ติด Authn อยู่แล้ว)
	g := r.Group("/authen")
	{
		// list/detail (view)
		g.GET("/cards", middleware.Authorize("card_view"), h.List)
		g.GET("/cards/:id", middleware.Authorize("card_view"), h.Detail)

		// create/update/status (edit)
		g.POST("/cards", middleware.Authorize("card_add"), h.Create)
		g.PATCH("/cards/:id", middleware.Authorize("card_edit"), h.Update)
		g.PATCH("/cards/:id/status", middleware.Authorize("card_edit"), h.UpdateStatus)

		// comments
		g.POST("/cards/:id/comments", middleware.Authorize("comment_add"), h.AddComment)
		g.PATCH("/cards/comments/:commentId", middleware.Authorize("comment_edit"), h.UpdateComment)
		g.GET("/cards/:id/comments", middleware.Authorize("comment_view"), h.ListComments)

		// History
		g.POST("/cards/:id/keep", middleware.Authorize("card_edit"), h.Keep)
		g.GET("/cards/:id/history", middleware.Authorize("card_view"), h.ListHistory)
	}
}
