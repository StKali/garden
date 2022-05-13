package api

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stkali/errors"
	"github.com/stkali/log"
	"net/http"
	"strings"

	db "github.com/stkali/garden/db/sqlc"
	"github.com/stkali/garden/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationTypeBasic  = "basic"
	authorizationPayloadKey = "authorization_payload"
)

// bearer authentication
// header item
// Authorization bearer token

// basic
// Authorization Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==

// authMiddleware return gin authenticattion middleware function
func authMiddleware(maker token.Maker, store db.Querier) gin.HandlerFunc {

	// return auth middleware handle
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}
		authorizationType := strings.ToLower(fields[0])
		payload := new(token.Payload)
		switch authorizationType {
		// basic auth
		case authorizationTypeBasic:
			user, werr := verifyBasicText(ctx, store, fields[1])
			if werr != nil {
				ctx.AbortWithStatusJSON(werr.code, errResponse(werr.message))
				return
			}
			info, werr := makeSession(ctx, user, maker, store)
			if werr != nil {
				ctx.AbortWithStatusJSON(werr.code, werr.message)
				return
			}
			// reset auth
			ctx.Set(authorizationHeaderKey, "Bearer "+info.AccessToken)
			payload = &token.Payload{
				ID:        info.SessionID,
				Username:  info.Username,
				IssuedAt:  info.CreatedAt,
				ExpiredAt: info.AccessExpireAt,
			}
		// token
		case authorizationTypeBearer:
			var err error
			payload, err = maker.VerifyToken(fields[1])
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
				return
			}
		default:
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func verifyBasicText(ctx *gin.Context, store db.Querier, text string) (*db.User, *Error) {
	peer, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		log.Errorf("decode base64 failed, content:%s, err: %s", text, err)
		return nil, newError(http.StatusUnauthorized, "invalid basic auth fields")
	}
	fields := strings.Split(string(peer), ":")
	if len(fields) != 2 {
		log.Errorf("failed to parse username and password, fields: %s", string(peer))
		return nil, newError(http.StatusUnauthorized, "invalid basic auth fields")
	}
	user, werr := getAndVerifyUser(ctx, store, fields[0], fields[1])
	if werr != nil {
		return nil, werr
	}
	return user, nil
}
