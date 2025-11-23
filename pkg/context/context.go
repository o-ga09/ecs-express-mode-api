package context

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CtxUserKey string
type CtxGinKey string
type CtxRequestIDKey string

const GinCtx CtxGinKey = "ginCtx"
const USERID CtxUserKey = "userID"
const REQUESTID CtxRequestIDKey = "requestId"

func GetCtxFromUser(ctx context.Context) string {
	userID, ok := ctx.Value(USERID).(string)
	if !ok {
		return ""
	}
	return userID
}

func SetCtxFromUser(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, USERID, userID)
}

func SetRequestID(ctx context.Context) context.Context {
	reqID := GetRequestID(ctx)
	if reqID != "" {
		return ctx
	}
	return context.WithValue(ctx, REQUESTID, uuid.NewString())
}

func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(REQUESTID).(string); ok {
		return reqID
	}
	return ""
}

func SetCtxGinCtx(ctx context.Context, c *gin.Context) context.Context {
	ctx = context.WithValue(ctx, GinCtx, c)
	return ctx
}

func GetCtxGinCtx(ctx context.Context) *gin.Context {
	if c, ok := ctx.Value(GinCtx).(*gin.Context); ok {
		return c
	}
	return nil
}
