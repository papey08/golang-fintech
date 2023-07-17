package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"homework10/internal/model/users"
	"homework10/internal/ports/grpc/pb"
)

// UserToUserResponse converts users.User to *UserResponse
func UserToUserResponse(u users.User) *pb.UserResponse {
	return &pb.UserResponse{
		Id:       u.ID,
		Nickname: u.Nickname,
		Email:    u.Email,
	}
}

func (s *Server) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.UserResponse, error) {
	gotUser, getErr := s.App.GetUserByID(ctx, req.Id)
	if getErr != nil {
		return nil, ErrorToGRPCError(getErr)
	}

	return UserToUserResponse(gotUser), nil
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	createdUser, createErr := s.App.CreateUser(ctx, req.Nickname, req.Email)
	if createErr != nil {
		return nil, ErrorToGRPCError(createErr)
	}

	return UserToUserResponse(createdUser), nil
}

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	updatedUser, updateErr := s.App.UpdateUser(ctx, req.Id, req.Nickname, req.Email)
	if updateErr != nil {
		return nil, ErrorToGRPCError(updateErr)
	}

	return UserToUserResponse(updatedUser), nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*empty.Empty, error) {
	return &empty.Empty{}, ErrorToGRPCError(s.App.DeleteUser(ctx, req.Id))
}
