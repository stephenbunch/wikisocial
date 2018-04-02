package app

import (
	"context"

	pb "github.com/stephenbunch/wikitribe/server/_proto"
)

type PostService struct{}

func NewPostService() *PostService {
	return &PostService{}
}

func (s *PostService) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.Post, error) {
	return &pb.Post{Id: "12345", Message: "hello world"}, nil
}
