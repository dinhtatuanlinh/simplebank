package gapi

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

func GrpcLogger(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", duration).
		Str("start_time", startTime.String()).
		Msg("received a gRPC request")
	return result, err
}

func HttpLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		duration := time.Since(startTime)

		logger := log.Info()
		if ctx.Writer.Status() != http.StatusOK {
			logger = log.Error()
			err := ctx.Errors
			fmt.Print(err)
			for _, ginErr := range ctx.Errors {
				logger.Str("error", ginErr.Error())
			}
		}

		logger.Str("protocol", "http").
			Str("method", ctx.Request.Method).
			Str("path", ctx.Request.URL.Path).
			Int("status_code", ctx.Writer.Status()).
			Str("status_text", http.StatusText(ctx.Writer.Status())).
			Dur("duration", duration).
			Str("start_time", startTime.String()).
			Msg("received a http request")
	}
}
