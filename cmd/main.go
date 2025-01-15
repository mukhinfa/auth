package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/mukhinfa/auth/pkg/user/v1"
)

const (
	grpcPort = 50051
)

type server struct {
	desc.UnimplementedUserServiceV1Server
}

func (s *server) Get(_ context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Println(color.RedString("Get user request:"), fmt.Sprintf("%+v", req))

	return &desc.GetResponse{
		Id:        req.Id,
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      desc.Role(gofakeit.Number(1, 3)), //nolint:gosec // The range of values is specified
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	}, nil
}

func (s *server) Create(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Println(color.RedString("Create user request:"), fmt.Sprintf("%+v", req))

	return &desc.CreateResponse{
		Id: int64(gofakeit.Number(1, 100)),
	}, nil
}

func (s *server) Update(_ context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Println(color.RedString("Update user request:"), fmt.Sprintf("%+v", req))
	return &emptypb.Empty{}, nil
}

func (s *server) Delete(_ context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Println(color.RedString("Delete user request:"), fmt.Sprintf("%+v", req))
	return &emptypb.Empty{}, nil
}

func main() {
	log.Println(color.GreenString("Starting server..."))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserServiceV1Server(s, &server{})

	log.Println(color.GreenString("Server started at %s", lis.Addr().String()))

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
