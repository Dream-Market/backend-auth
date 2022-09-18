package services

import (
	"context"
	"net/http"

	"github.com/Dream-Market/backend-auth/pkg/db"
	"github.com/Dream-Market/backend-auth/pkg/models"
	"github.com/Dream-Market/backend-auth/pkg/pb"
	"github.com/Dream-Market/backend-auth/pkg/utils"
	"github.com/Dream-Market/backend-auth/pkg/validation"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	H   db.Handler
	Jwt utils.JwtWrapper
}

func (s *Server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	var user models.User

	if _, err := s.H.UserHandler.FindByEmailOrPhone(req.Email, req.Phone); err == nil {
		return &pb.RegisterUserResponse{
			Status: http.StatusConflict,
			Error:  "E-Mail already exists",
		}, nil
	}

	var ok bool
	if user.Email, ok = validation.ValidateEmail(req.Email); !ok {
		return &pb.RegisterUserResponse{
			Status: http.StatusBadRequest,
			Error:  "E-Mail is invalid",
		}, nil
	}

	if user.Phone, ok = validation.ValidatePhone(req.Phone); !ok {
		return &pb.RegisterUserResponse{
			Status: http.StatusBadRequest,
			Error:  "Phone is invalid",
		}, nil
	}

	if !validation.ValidatePassword(req.Password) {
		return &pb.RegisterUserResponse{
			Status: http.StatusBadRequest,
			Error:  "Password is invalid",
		}, nil
	}

	user.Password = utils.HashPassword(req.Password)

	user, err := s.H.UserHandler.Create(user)
	if err != nil {
		return &pb.RegisterUserResponse{
			Status: http.StatusInternalServerError,
		}, err
	}

	return &pb.RegisterUserResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (resp *pb.LoginUserResponse, err error) {
	var user models.User

	if user, err = s.H.UserHandler.FindByEmailOrPhone(req.Login, req.Login); err != nil {
		return &pb.LoginUserResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)

	if !match {
		return &pb.LoginUserResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	session, err := s.H.SessionHandler.Create(models.Session{UserId: user.Id})
	if err != nil {
		return &pb.LoginUserResponse{
			Status: http.StatusInternalServerError,
			Error:  "Session creation failed",
		}, err
	}

	token, _ := s.Jwt.GenerateToken(session)

	return &pb.LoginUserResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (s *Server) ValidateUser(ctx context.Context, req *pb.ValidateUserRequest) (*pb.ValidateUserResponse, error) {
	claims, err := s.Jwt.ValidateToken(req.Token)

	if err != nil {
		return &pb.ValidateUserResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	var user models.User

	if active, err := s.H.SessionHandler.IsActive(claims.SessionId); !active {
		return &pb.ValidateUserResponse{
			Status: http.StatusUnauthorized,
		}, err
	}

	return &pb.ValidateUserResponse{
		Status: http.StatusOK,
		UserId: user.Id,
	}, nil
}
