package handlers

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/jasurxaydarov/book_shop_api_getway/genproto/user_service"
	"github.com/jasurxaydarov/book_shop_api_getway/mail"
	"github.com/jasurxaydarov/book_shop_api_getway/pkg/helpers"
	"github.com/jasurxaydarov/book_shop_api_getway/token"

	"github.com/saidamir98/udevs_pkg/logger"
)

func (h *Handler) GetUserById(ctx *gin.Context) {

	var req user_service.GetByIdReq

	req.Id = ctx.Param("id")

	resp, err := h.service.GetUserSevice().GetUser(context.Background(), &req)

	if err != nil {
		ctx.JSON(500, err)
		return
	}

	ctx.JSON(201, resp)

}

func (h *Handler) GetUsers(ctx *gin.Context) {

	var req user_service.GetListReq

	ctx.BindJSON(&req)

	resp, err := h.service.GetUserSevice().GetUsers(context.Background(), &req)

	if err != nil {
		ctx.JSON(500, err)
		return
	}

	ctx.JSON(201, resp)

}

func (h *Handler) CheckUser(ctx *gin.Context) {

	var reqBody user_service.CheckUser

	err := ctx.BindJSON(&reqBody)

	if err != nil {

		ctx.JSON(400, err.Error())
		return
	}

	isExists, err := h.service.GetUserSevice().CheckExists(context.Background(), &user_service.Common{
		TableName:  "users",
		ColumnName: "email",
		Expvalue:   reqBody.Email,
	})

	if err != nil {

		ctx.JSON(500, err)
		return
	}

	if isExists.IsExists {
		ctx.JSON(201, user_service.CheckExists{
			IsExists: isExists.IsExists,
			Status:   "sign-in",
		})
		return
	}

	otp := user_service.OtpData{
		Otp:   mail.GenerateOtp(6),
		Email: reqBody.Email,
	}

	otpdataB, err := json.Marshal(otp)

	if err != nil {

		ctx.JSON(500, err)
		return
	}

	err = h.cache.Set(ctx, reqBody.Email, string(otpdataB), 120)

	err = mail.SendMail([]string{reqBody.Email}, otp.Otp)

	if err != nil {
		h.log.Error("errrr on Send mail", logger.Error(err))
		ctx.JSON(500, err)
		return
	}

	ctx.JSON(201, struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "registr",
		Message: "we sent otp",
	})

}

func (h *Handler) SignUp(ctx *gin.Context) {

	var otpData user_service.OtpData

	var reqBody user_service.UserCreateReq

	err := ctx.BindJSON(&reqBody)

	if err != nil {
		h.log.Error("errrr on ShouldBindJSON", logger.Error(err))
		ctx.JSON(500, err.Error())
		return
	}

	otpSData, err := h.cache.GetDell(ctx, reqBody.Email)

	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	if otpSData == "" {
		ctx.JSON(201, "otp is expired")
		return
	}
	err = json.Unmarshal([]byte(otpSData), &otpData)

	if otpData.Otp != reqBody.Otp {

		ctx.JSON(405, "incorrect otp")
		return
	}

	reqBody.Password, err = helpers.HashPassword(reqBody.Password)

	claim, err := h.service.GetUserSevice().CreateUser(context.Background(), &reqBody)

	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	accessToken, err := token.GenerateJWT(*&user_service.Clamis{UserId: claim.UserId, UserRole: claim.UserRole})

	if err != nil {
		ctx.JSON(201, "registreted")
		return
	}

	ctx.JSON(201, accessToken)

}
/////////////////////////////////////////////////

func (h *Handler) SigIn(ctx *gin.Context) {

	var reqBody user_service.UserLogIn

	err := ctx.BindJSON(&reqBody)

	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	claim, err := h.service.GetUserSevice().UserLogin(ctx, &reqBody)


	if err != nil {
		if err.Error() == "password is incorrect" {
			ctx.JSON(405, err.Error())
			return
		}
		ctx.JSON(500, err.Error())
		return
	}

	accessToken, err := token.GenerateJWT(*claim)

	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(201, accessToken)

}

func (h *Handler) UpdateUser(ctx *gin.Context) {

	req := user_service.UserUpdateReq{}
	tokenString := ctx.GetHeader("authorization")

	if tokenString == "" {
		ctx.JSON(401, gin.H{"error": "authorization token not provided"})
		ctx.Abort()
	}

	claim, err := token.ParseJWT(tokenString)
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		ctx.Abort()
	}

	ctx.BindJSON(&req)
	req.UserId = claim.UserId

	pp,err:=helpers.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		ctx.Abort()
	}

	req.Password=pp

	res, err := h.service.GetUserSevice().UpdateUser(context.Background(), &req)

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201,res)

}

func (h *Handler) DeleteUser(ctx *gin.Context) {

	req := user_service.DeleteReq{}
	tokenString := ctx.GetHeader("authorization")

	if tokenString == "" {
		ctx.JSON(401, gin.H{"error": "authorization token not provided"})
		ctx.Abort()
	}

	claim, err := token.ParseJWT(tokenString)
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		ctx.Abort()
	}

	ctx.BindJSON(&req)
	req.Id = claim.UserId

	_, err = h.service.GetUserSevice().DeleteUser(context.Background(), &req)

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201,"succssfully deleted")

}

func (h *Handler) AdmDeleteUser(ctx *gin.Context) {

	req := user_service.DeleteReq{}


	ctx.BindJSON(&req)

	_, err := h.service.GetUserSevice().DeleteUser(context.Background(), &req)

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201,"succssfully deleted")

}

func (h *Handler) AdmUpdateUser(ctx *gin.Context) {

	req := user_service.UserUpdateReq{}
	
	ctx.BindJSON(&req)

	pp,err:=helpers.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		ctx.Abort()
	}

	req.Password=pp

	res, err := h.service.GetUserSevice().UpdateUser(context.Background(), &req)

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201,res)

}