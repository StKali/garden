package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stkali/errors"
	db "github.com/stkali/garden/db/sqlc"
	"github.com/stkali/garden/util"
	"net/http"
	"time"
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
	User        UserResponse
	AccessToken string
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
	user, err := s.store.GetUser(context.Background(), req.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errResponse(err))
		return
	}
	// verify password
	if err = util.VerifyPassword(req.Password, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(errors.New("username or password error")))
		return
	}
	// create token
	token, err := s.maker.CreateToken(req.Username, time.Second*60*60*24*7)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(errors.New("failed to create user token")))
	}

	res := &LoginResponse{
		User: UserResponse{
			Username:          user.Username,
			FullName:          user.FullName,
			Email:             user.Email,
			PasswordChangedAt: user.PasswordChangedAt,
			CreatedAt:         user.CreatedAt,
		},
		AccessToken: token,
	}
	ctx.JSON(http.StatusOK, res)
}

// Start http server on address
func (s *Server) Start(address string) {
	err := s.engine.Run(address)
	util.CheckError("cannot start HTTP server, err:", err)
}

// errResponse wrap err to gin.H
func errResponse(err error) gin.H {
	return gin.H{"err": err.Error()}
}
