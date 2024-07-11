package server

import (
	"context"
	"testing"

	pb "github.com/brotherlogic/auth/proto"
)

func TestLogin_Success(t *testing.T) {
	ctx := context.Background()

	s := NewServer()
	s.userMap["test"] = "test"

	resp, err := s.GetToken(ctx, &pb.GetTokenRequest{User: "test", Password: "test"})
	if err != nil {
		t.Errorf("Unable to get token: %v", err)
	}

	_, err = s.Authenticate(ctx, &pb.AuthenticateRequest{Token: resp.GetToken()})
	if err != nil {
		t.Errorf("Unable to authenticate: %v", err)
	}
}

func TestLogin_Failure(t *testing.T) {
	ctx := context.Background()

	s := NewServer()

	_, err := s.GetToken(ctx, &pb.GetTokenRequest{User: "test", Password: "test"})
	if err == nil {
		t.Errorf("Token returned: %v", err)
	}
}
