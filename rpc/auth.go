package rpc

import (
	"context"
	db "github.com/stkali/garden/db/sqlc"
	"github.com/stkali/garden/pb"
	"github.com/stkali/garden/token"
	"github.com/stkali/garden/util"
	"github.com/stkali/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var InternalErrorString = "occurred error, please try again later"

type Error struct {
	code    codes.Code
	message any
}

func makeSession(ctx context.Context, user *db.User, maker token.Maker, store db.Querier) (*pb.UserInfo, *Error) {
	// create access token
	accessToken, accessPayload, err := maker.CreateToken(user.Username, setting.TokenDuration)
	if err != nil {
		log.Errorf("failed to create user token, err: %s", err)
		return nil, &Error{codes.Internal, InternalErrorString}
	}
	// create refresh token
	refreshToken, refreshPayload, err := maker.CreateToken(user.Username, setting.RefreshTokenDuration)
	if err != nil {
		log.Errorf("failed to create user refresh token, err: %s", err)
		return nil, &Error{codes.Internal, InternalErrorString}
	}
	
	// create session params to save session to database
	meta := metaFromCtx(ctx)
	createSessionParams := db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     refreshPayload.Username,
		RefreshToken: refreshToken,
		UserAgent:    meta.UserAgent,
		ClientIp:     meta.ClientIP,
		ExpiresAt:    refreshPayload.ExpiredAt,
	}

	// save user ssession to database
	session, err := store.CreateSession(ctx, createSessionParams)
	if err != nil {
		log.Errorf("failed to save session in storage, err: %s", err)
		return nil, &Error{codes.Internal, InternalErrorString}
	}

	info := &pb.UserInfo{
		SessionID:       session.ID.String(),
		AccessToken:     accessToken,
		AccessExpireAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:    refreshToken,
		RefreshExpireAt: timestamppb.New(refreshPayload.ExpiredAt),
		Username:        user.Username,
		FullName:        user.FullName,
		Email:           user.Email,
		CreatedAt:       timestamppb.New(user.CreatedAt),
	}
	return info, nil
}

func (s *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.UserInfo, error) {

	if request.Password != request.ConfirmPassword {
		return nil, status.Errorf(codes.Internal, "password != confirm password")
	}

	passwordHash, err := util.HashPassword(request.Password)
	if err != nil {
		log.Errorf("failed to hash password, password: %s err: %s", request.Password, err)
		return nil, status.Errorf(codes.Internal, InternalErrorString)
	}
	// prepare save user to storage
	arg := db.CreateUserParams{
		Username:       request.Username,
		HashedPassword: passwordHash,
		FullName:       request.FullName,
		Email:          request.Email,
	}
	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		log.Errorf("failed to create user, err: %s", err)
		return nil, status.Errorf(codes.Internal, InternalErrorString)
	}
	info, werr := makeSession(ctx, &user, s.maker, s.store)
	if werr != nil {
		return nil, status.Errorf(werr.code, "%s", werr.message)
	}

	return info, nil
}

func (s *Server) Login(ctx context.Context, request *pb.LoginRequest) (*pb.UserInfo, error) {

	// ensure has registered
	user, err := s.store.GetUser(ctx, request.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "not registered user: %s", request.Username)
	}

	// verify password
	if err = util.VerifyPassword(request.Password, user.HashedPassword); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "username or password error")
	}
	info, werr := makeSession(ctx, &user, s.maker, s.store)
	if werr != nil {
		return nil, status.Errorf(werr.code, "%s", werr.message)
	}
	return info, nil
}
