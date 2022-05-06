package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stkali/errors"
	db "github.com/stkali/garden/db/sqlc"
	"github.com/stkali/garden/token"
	"github.com/stkali/garden/util"
)

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
		ctx.JSON(http.StatusBadRequest, errResponse(errors.New("password != verify password")))
		return
	}

	passwordHash, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(errors.New("failed to hash password, err: %v", err)))
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
		ctx.JSON(http.StatusInternalServerError, errResponse(errors.New("failed to create user, err: %v", err)))
		return
	}

	res := UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, res)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type LoginResponse struct {
	SessionID   uuid.UUID            `json:"session_id"`
	User        UserResponse      `json:"user"`
	AccessToken string            `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	AccessExpireAt time.Time  `json:"access_expire_at"`
	RefreshExpireAt time.Time `json:"refresh_expire_at"`
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
	// ensure has registered
	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errResponse(err))
		return
	}
	// verify password
	if err = util.VerifyPassword(req.Password, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(errors.New("username or password error, err: %s", err)))
		return
	}
	// create access token
	token, payload, err := s.maker.CreateToken(req.Username, setting.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(errors.New("failed to create user token, err: %s", err)))
	}

	// create refresh token
	refreshToken, refreshPayload, err := s.maker.CreateToken(req.Username, setting.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(errors.New("failed to create user refresh token, err: %s", err)))
		return
	}
	// create session params to save session to database
	createSessionParams := db.CreateSessionParams{
		ID: refreshPayload.ID,
		Username: refreshPayload.Username,
		RefreshToken: refreshToken,
		UserAgent: ctx.Request.UserAgent(),
		ClientIp: ctx.ClientIP(),
		ExpiresAt: refreshPayload.ExpiredAt,
	}
	// save user ssession to database
	session, err := s.store.CreateSession(ctx, createSessionParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(errors.New("failed to store session to db, err: %s", err)))
		return
	}
	// make login success response
	res := &LoginResponse{
		User: UserResponse{
			Username:          user.Username,
			FullName:          user.FullName,
			Email:             user.Email,
			PasswordChangedAt: user.PasswordChangedAt,
			CreatedAt:         user.CreatedAt,
		},
		AccessToken: token,
		AccessExpireAt: payload.ExpiredAt,
		RefreshToken: refreshToken,
		RefreshExpireAt: refreshPayload.ExpiredAt,
		SessionID: session.ID,
	}
	ctx.JSON(http.StatusOK, res)
}


// Home view
func (s *Server) Home(ctx *gin.Context) {

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	ctx.JSON(http.StatusOK, gin.H{
		"info": fmt.Sprintf("welcome to %s home", payload.Username),
		"username": payload.Username,
		"expire at": payload.ExpiredAt,
	})
}

// errResponse wrap err to gin.H
func errResponse(err error) gin.H {
	return gin.H{"err": err.Error()}
}
