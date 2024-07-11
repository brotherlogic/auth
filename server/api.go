package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brotherlogic/auth/proto"
)

type Server struct {
	token   *pb.AuthToken
	userMap map[string]string
}

func NewServer() *Server {
	return &Server{
		token:   &pb.AuthToken{},
		userMap: make(map[string]string),
	}
}

func (s *Server) GetToken(ctx context.Context, req *pb.GetTokenRequest) (*pb.GetTokenResponse, error) {
	if s.userMap[req.GetUser()] != req.GetPassword() {
		return nil, status.Errorf(codes.PermissionDenied, "Incorrect password for %v", req.GetUser())
	}

	s.recycleToken()

	return &pb.GetTokenResponse{Token: s.token.GetToken()}, nil
}

func (s *Server) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	// Recycle before we start in order to ensure we're not authenticating a stale token
	s.recycleToken()

	if req.GetToken() != s.token.GetToken() {
		return nil, status.Errorf(codes.PermissionDenied, "Bad token")
	}

	return &pb.AuthenticateResponse{}, nil
}
