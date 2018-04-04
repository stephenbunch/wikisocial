package app

import (
	"context"

	"github.com/gocql/gocql"
	pb "github.com/stephenbunch/wikitribe/server/_proto"
)

type user struct {
	id   gocql.UUID
	name string
}

// The WikiTribeService is the main service for interacting with WikiTribe
// objects.
type WikiTribeService struct {
	db *gocql.Session
}

// NewWikiTribeService creates a new instance of WikiTribeService.
func NewWikiTribeService(db *gocql.Session) *WikiTribeService {
	return &WikiTribeService{db: db}
}

// CreateUser creates a new user.
func (s *WikiTribeService) CreateUser(ctx context.Context,
	req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	id, err := gocql.RandomUUID()
	if err != nil {
		return nil, err
	}
	if err = s.db.Query(`INSERT INTO users (id, name) VALUES (?, ?)`,
		id, req.Name).Exec(); err != nil {
		return nil, err
	}
	return &pb.UserResponse{Id: id.String(), Name: req.Name}, nil
}

// GetUser gets a user.
func (s *WikiTribeService) GetUser(ctx context.Context,
	req *pb.GetUserRequest) (*pb.UserResponse, error) {
	var id gocql.UUID
	var name string
	uuid, err := gocql.ParseUUID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = s.db.Query(`SELECT id, name FROM users WHERE id = ?`,
		uuid).Consistency(gocql.One).Scan(&id, &name); err != nil {
		return nil, err
	}
	return &pb.UserResponse{Id: id.String(), Name: name}, nil
}

// ListUsers returns all users.
func (s *WikiTribeService) ListUsers(req *pb.ListUsersRequest,
	stream pb.WikiTribeService_ListUsersServer) error {
	iter := s.db.Query(`SELECT id, name FROM users`).Iter()
	var id gocql.UUID
	var name string
	for iter.Scan(&id, &name) {
		stream.Send(&pb.UserResponse{Id: id.String(), Name: name})
	}
	if err := iter.Close(); err != nil {
		return err
	}
	return nil
}
