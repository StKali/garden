package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	sterr "errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/stkali/garden/db/sqlc"
	"github.com/stkali/garden/token"
	"github.com/stkali/garden/util"
	"github.com/stkali/log"
)

var InternalError = sterr.New("occurred error, please try again later")

type CreateUserRequest struct {
	Username        string `json:"username" binding:"required,alphanum"`
	FullName        string `json:"full_name" binding:"required,alphanum"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Email           string `json:"email"`
}

// CreateUser user register view
func (s *Server) CreateUser(ctx *gin.Context) {

	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, errResponse("password != verify password"))
		return
	}
	passwordHash, err := util.HashPassword(req.Password)
	if err != nil {
		log.Errorf("failed to hash password, password: %s err: %s", req.Password, err)
		ctx.JSON(http.StatusInternalServerError, errResponse(InternalError.Error()))
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: passwordHash,
		FullName:       req.FullName,
		Email:          req.Email,
	}
	user, err := s.store.CreateUser(context.Background(), arg)
	if err != nil {
		log.Errorf("failed to create user, err: %s", err)
		ctx.JSON(http.StatusInternalServerError, errResponse(InternalError))
		return
	}
	info, werr := makeSession(ctx, &user, s.maker, s.store)
	if werr != nil {
		ctx.JSON(werr.code, errResponse(werr.message))
		return
	}
	ctx.JSON(http.StatusOK, info)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	SessionID       uuid.UUID `json:"session_id"`
	AccessToken     string    `json:"access_token"`
	AccessExpireAt  time.Time `json:"access_expire_at"`
	RefreshToken    string    `json:"refresh_token"`
	RefreshExpireAt time.Time `json:"refresh_expire_at"`
	Username        string    `json:"username"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email"`
	CreatedAt       time.Time `json:"created_at"`
}

type Error struct {
	code    int
	message any
}

func makeSession(ctx *gin.Context, user *db.User, maker token.Maker, store db.Querier) (*UserInfo, *Error) {
	// create access token
	accessToken, accessPayload, err := maker.CreateToken(user.Username, setting.TokenDuration)
	if err != nil {
		log.Errorf("failed to create user token, err: %s", err)
		return nil, newError(http.StatusInternalServerError, InternalError)
	}
	// create refresh token
	refreshToken, refreshPayload, err := maker.CreateToken(user.Username, setting.RefreshTokenDuration)
	if err != nil {
		log.Errorf("failed to create user refresh token, err: %s", err)
		return nil, newError(http.StatusInternalServerError, InternalError)
	}
	// create session params to save session to database
	createSessionParams := db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     refreshPayload.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    refreshPayload.ExpiredAt,
	}
	// save user ssession to database
	session, err := store.CreateSession(ctx, createSessionParams)
	if err != nil {
		log.Errorf("failed to save session in storage, err: %s", err)
		return nil, newError(http.StatusInternalServerError, InternalError)
	}

	info := &UserInfo{
		SessionID:       session.ID,
		AccessToken:     accessToken,
		AccessExpireAt:  accessPayload.ExpiredAt,
		RefreshToken:    refreshToken,
		RefreshExpireAt: refreshPayload.ExpiredAt,
		Username:        user.Username,
		FullName:        user.FullName,
		Email:           user.Email,
		CreatedAt:       user.CreatedAt,
	}
	return info, nil
}

func newError(code int, message any) *Error {
	return &Error{code: code, message: message}
}

func getAndVerifyUser(ctx *gin.Context, store db.Querier, username, password string) (*db.User, *Error) {
	// ensure has registered
	user, err := store.GetUser(ctx, username)
	if err != nil {
		return nil, newError(http.StatusNotFound, fmt.Sprintf("not found user named: %q", username))
	}
	// verify password
	if err = util.VerifyPassword(password, user.HashedPassword); err != nil {
		return nil, newError(http.StatusUnauthorized, "invalid username or password")
	}
	return &user, nil
}

// Login handle view
func (s *Server) Login(ctx *gin.Context) {

	var req LoginRequest
	// get arguments
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	// verify user status
	user, werr := getAndVerifyUser(ctx, s.store, req.Username, req.Password)
	if werr != nil {
		ctx.JSON(werr.code, errResponse(werr.message))
		return
	}
	// make session
	info, werr := makeSession(ctx, user, s.maker, s.store)
	if werr != nil {
		ctx.JSON(werr.code, errResponse(werr.message))
		return
	}
	ctx.JSON(http.StatusOK, info)
}

// Home view
func (s *Server) Home(ctx *gin.Context) {

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	ctx.JSON(http.StatusOK, gin.H{
		"info":      fmt.Sprintf("welcome to %s home", payload.Username),
		"username":  payload.Username,
		"expire at": payload.ExpiredAt,
	})
}

// errResponse wrap err to gin.H
func errResponse(err any) gin.H {
	var errString string
	switch err.(type) {
	case error:
		errString = err.(error).Error()
	case string:
		errString = err.(string)
	default:
		errString = fmt.Sprintf("%s", err)
	}
	return gin.H{"err": errString}
}
