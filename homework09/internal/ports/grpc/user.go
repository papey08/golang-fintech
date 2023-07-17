package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"homework9/internal/model/users"
)

// UserToUserResponse converts users.User to *UserResponse
func UserToUserResponse(u users.User) *UserResponse {
	return &UserResponse{
		Id:       u.ID,
		Nickname: u.Nickname,
		Email:    u.Email,
	}
}

func (s *Server) GetUserByID(ctx context.Context, req *GetUserByIDRequest) (*UserResponse, error) {
	gotUser, getErr := s.App.GetUserByID(ctx, req.Id)
	if getErr != nil {
		return nil, ErrorToGRPCError(getErr)
	}

	return UserToUserResponse(gotUser), nil
}

func (s *Server) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	createdUser, createErr := s.App.CreateUser(ctx, req.Nickname, req.Email)
	if createErr != nil {
		return nil, ErrorToGRPCError(createErr)
	}

	return UserToUserResponse(createdUser), nil
}

func (s *Server) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UserResponse, error) {
	updatedUser, updateErr := s.App.UpdateUser(ctx, req.Id, req.Nickname, req.Email)
	if updateErr != nil {
		return nil, ErrorToGRPCError(updateErr)
	}

	return UserToUserResponse(updatedUser), nil
}

func (s *Server) DeleteUser(ctx context.Context, req *DeleteUserRequest) (*empty.Empty, error) {
	return &empty.Empty{}, ErrorToGRPCError(s.App.DeleteUser(ctx, req.Id))
}
