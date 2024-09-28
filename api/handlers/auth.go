package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jasurxaydarov/book_shop_api_getway/genproto/product_service"
	"github.com/saidamir98/udevs_pkg/logger"
)

func (h *Handler) CreateAuth(ctx *gin.Context) {

	log := logger.NewLogger("", logger.LevelDebug)
	var req product_service.AuthorUpdateReq

	ctx.BindJSON(&req)

	resp, err := h.service.GetProductSevice().CreateAuth(context.Background(), &req)

	if err != nil {

		log.Debug("errrrr ")

		log.Debug(err.Error())
		ctx.JSON(500, err)
		return
	}

	ctx.JSON(201, resp)
}


func (h *Handler) GetAuthById(ctx *gin.Context) {

	var req product_service.GetByIdReq

	req.Id = ctx.Param("id")

	resp, err := h.service.GetProductSevice().GetAuth(context.Background(), &req)

	if err != nil {
		ctx.JSON(500, err)
		return
	}

	ctx.JSON(201, resp)

}


