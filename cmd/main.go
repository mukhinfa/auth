package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	desc "github.com/mukhinfa/auth/pkg/user/v1"
)

const (
	grpcPort = 50051
	dbDSN    = "host=localhost port=54321 dbname=note user=note-user password=note-password sslmode=disable"
)

type server struct {
	desc.UnimplementedUserServiceV1Server
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Println(color.RedString("Get user request:"), fmt.Sprintf("%+v", req))

	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	builderSelectOne := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.Id}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var id int64
	var name string
	var email string
	var role desc.Role
	var createdAt time.Time
	var updatedAt time.Time

	err = pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to select user: %w", err)
	}

	return &desc.GetResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      role,
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt),
	}, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Println(color.RedString("Create user request:"), fmt.Sprintf("%+v", req))

	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	builderInsert := sq.Insert("users").
		Columns("name", "email", "role").
		Values(req.Name, req.Email, req.Role).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var id int64
	err = pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &desc.CreateResponse{Id: id}, nil
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
