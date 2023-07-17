package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/adapters/user_repo"
	"homework9/internal/app"
	"homework9/internal/model/errs"
	pb "homework9/internal/ports/grpc"
	"net"
	"testing"
	"time"
)

func GRPCClientInit(t *testing.T) (pb.AdServiceClient, context.Context) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		err := lis.Close()
		if err != nil {
			return
		}
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	pb.RegisterAdServiceServer(srv, &pb.Server{
		AdServiceServer: nil,
		App:             app.NewApp(adrepo.New(), user_repo.New()),
	})
	go func() {
		err := srv.Serve(lis)
		assert.NoError(t, err)
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	t.Cleanup(func() {
		err = conn.Close()
		if err != nil {
			return
		}
	})
	return pb.NewAdServiceClient(conn), ctx
}

func TestGRPCCreateAd(t *testing.T) {
	client, ctx := GRPCClientInit(t)

	createdUser, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Nickname: "papey08",
		Email:    "email@golang.com",
	})
	assert.NoError(t, err)

	createdAd, err := client.CreateAd(ctx, &pb.CreateAdRequest{
		Title:    "hello",
		Text:     "world",
		AuthorId: createdUser.Id,
	})
	assert.NoError(t, err)
	assert.Zero(t, createdAd.Id)
	assert.Equal(t, createdAd.Title, "hello")
	assert.Equal(t, createdAd.Text, "world")
	assert.Equal(t, createdAd.AuthorId, createdUser.Id)
	assert.False(t, createdAd.Published)
}

func TestGRPCChangeAdStatus(t *testing.T) {
	client, ctx := GRPCClientInit(t)

	createdUser, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Nickname: "papey08",
		Email:    "email@golang.com",
	})
	assert.NoError(t, err)

	createdAd, err := client.CreateAd(ctx, &pb.CreateAdRequest{
		Title:    "hello",
		Text:     "world",
		AuthorId: createdUser.Id,
	})
	assert.NoError(t, err)

	response, err := client.ChangeAdStatus(ctx, &pb.ChangeAdStatusRequest{
		Id:        createdAd.Id,
		AuthorId:  createdUser.Id,
		Published: true,
	})
	assert.NoError(t, err)
	assert.True(t, response.Published)

	response, err = client.ChangeAdStatus(ctx, &pb.ChangeAdStatusRequest{
		Id:        createdAd.Id,
		AuthorId:  createdUser.Id,
		Published: false,
	})
	assert.NoError(t, err)
	assert.False(t, response.Published)

	response, err = client.ChangeAdStatus(ctx, &pb.ChangeAdStatusRequest{
		Id:        createdAd.Id,
		AuthorId:  createdUser.Id,
		Published: false,
	})
	assert.NoError(t, err)
	assert.False(t, response.Published)
}

func TestGRPCUpdateAd(t *testing.T) {
	client, ctx := GRPCClientInit(t)

	createdUser, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Nickname: "papey08",
		Email:    "email@golang.com",
	})
	assert.NoError(t, err)

	createdAd, err := client.CreateAd(ctx, &pb.CreateAdRequest{
		Title:    "hello",
		Text:     "world",
		AuthorId: createdUser.Id,
	})
	assert.NoError(t, err)

	response, err := client.UpdateAd(ctx, &pb.UpdateAdRequest{
		Id:       createdAd.Id,
		Title:    "привет",
		Text:     "мир",
		AuthorId: createdUser.Id,
	})
	assert.NoError(t, err)
	assert.Equal(t, response.Title, "привет")
	assert.Equal(t, response.Text, "мир")
}

func TestGRPCGetAdByID(t *testing.T) {
	client, ctx := GRPCClientInit(t)

	createdUser, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Nickname: "papey08",
		Email:    "email@golang.com",
	})
	assert.NoError(t, err)

	createdAd, err := client.CreateAd(ctx, &pb.CreateAdRequest{
		Title:    "hello",
		Text:     "world",
		AuthorId: createdUser.Id,
	})
	assert.NoError(t, err)

	response, err := client.GetAdByID(ctx, &pb.GetAdByIDRequest{Id: createdAd.Id})
	assert.NoError(t, err)
	assert.Equal(t, response.Id, createdAd.Id)
	assert.Equal(t, response.Text, createdAd.Text)
	assert.Equal(t, response.Title, createdAd.Title)
	assert.Equal(t, response.Published, createdAd.Published)
	assert.Equal(t, response.AuthorId, createdAd.AuthorId)
}

func TestGRPCGetAdsList(t *testing.T) {
	client, ctx := GRPCClientInit(t)

	createdUser, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Nickname: "papey08",
		Email:    "email@golang.com",
	})
	assert.NoError(t, err)

	createdAd, err := client.CreateAd(ctx, &pb.CreateAdRequest{
		Title:    "hello",
		Text:     "world",
		AuthorId: createdUser.Id,
	})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &pb.ChangeAdStatusRequest{
		Id:        createdAd.Id,
		AuthorId:  createdAd.AuthorId,
		Published: true,
	})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &pb.CreateAdRequest{
		Title:    "best cat",
		Text:     "not for sale",
		AuthorId: createdUser.Id,
	})
	assert.NoError(t, err)

	ads, err := client.GetAdsList(ctx, &pb.FilterRequest{})
	assert.NoError(t, err)
	assert.Len(t, ads.List, 1)
}

func TestGRPCSearchAds(t *testing.T) {
	client, ctx := GRPCClientInit(t)

	createdUser, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Nickname: "papey08",
		Email:    "email@golang.com",
	})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &pb.CreateAdRequest{
		Title:    "hello",
		Text:     "world",
		AuthorId: createdUser.Id,
	})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &pb.CreateAdRequest{
		Title:    "best cat",
		Text:     "not for sale",
		AuthorId: createdUser.Id,
	})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &pb.CreateAdRequest{
		Title:    "hello world",
		Text:     "привет мир",
		AuthorId: createdUser.Id,
	})
	assert.NoError(t, err)

	ads, err := client.SearchAds(ctx, &pb.SearchAdsRequest{Pattern: "hello"})
	assert.NoError(t, err)
	assert.Len(t, ads.List, 2)
}

func TestGRPCDeleteAd(t *testing.T) {
	client, ctx := GRPCClientInit(t)

	createdUser, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Nickname: "papey08",
		Email:    "email@golang.com",
	})
	assert.NoError(t, err)

	createdAd, err := client.CreateAd(ctx, &pb.CreateAdRequest{
		Title:    "hello",
		Text:     "world",
		AuthorId: createdUser.Id,
	})
	assert.NoError(t, err)

	_, err = client.DeleteAd(ctx, &pb.DeleteAdRequest{
		Id:       createdAd.Id,
		AuthorId: createdUser.Id,
	})
	assert.NoError(t, err)

	_, err = client.GetAdByID(ctx, &pb.GetAdByIDRequest{Id: createdAd.Id})
	assert.Error(t, err, errs.AdNotExist)
}

func TestGRPCGetUserByID(t *testing.T) {
	client, ctx := GRPCClientInit(t)

	createdUser, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Nickname: "papey08",
		Email:    "email@golang.com",
	})
	assert.NoError(t, err)

	gotUser, err := client.GetUserByID(ctx, &pb.GetUserByIDRequest{Id: createdUser.Id})
	assert.NoError(t, err)
	assert.Equal(t, createdUser.Id, gotUser.Id)
	assert.Equal(t, createdUser.Nickname, gotUser.Nickname)
	assert.Equal(t, createdUser.Email, gotUser.Email)
}

func TestGRPCCreateUser(t *testing.T) {
	client, ctx := GRPCClientInit(t)

	res, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Nickname: "papey08",
		Email:    "email@golang.com",
	})
	assert.NoError(t, err)
	assert.Zero(t, res.Id)
	assert.Equal(t, res.Nickname, "papey08")
	assert.Equal(t, res.Email, "email@golang.com")
}

func TestGRPCUpdateUser(t *testing.T) {
	client, ctx := GRPCClientInit(t)

	createdUser, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Nickname: "papey08",
		Email:    "email@golang.com",
	})
	assert.NoError(t, err)

	updatedUser, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:       createdUser.Id,
		Nickname: createdUser.Nickname,
		Email:    "email@golang.ru",
	})
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.Id, createdUser.Id)
	assert.Equal(t, updatedUser.Nickname, createdUser.Nickname)
	assert.Equal(t, updatedUser.Email, "email@golang.ru")
}

func TestGRPCDeleteUser(t *testing.T) {
	client, ctx := GRPCClientInit(t)

	createdUser, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Nickname: "papey08",
		Email:    "email@golang.com",
	})
	assert.NoError(t, err)

	_, err = client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: createdUser.Id})
	assert.NoError(t, err)

	_, err = client.GetUserByID(ctx, &pb.GetUserByIDRequest{Id: createdUser.Id})
	assert.Error(t, err, errs.UserNotExist)
}
