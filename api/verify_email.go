package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	db "simplebank/db/sqlc"
)

type verifyEmailRequest struct {
	EmailID    int64  `form:"email_id"`
	SecretCode string `form:"secret_code"`
}

type verifyEmailResponse struct {
	IsVerified bool `json:"is_verified"`
}

func (server *Server) verifyEmail(ctx *gin.Context) {
	var req verifyEmailRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	txResult, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    req.EmailID,
		SecretCode: req.SecretCode,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, verifyEmailResponse{
		IsVerified: txResult.User.IsEmailVerified,
	})
}
