package client

import (
	"context"
	"fmt"
	"log"
	"os/user"

	pb "github.com/brotherlogic/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	authClient pb.AuthServiceClient
	cToken     string
	conn       *grpc.ClientConn
}

func NewAuthInterceptor(ctx context.Context) (*AuthInterceptor, error) {
	conn, err := grpc.NewClient("kclust1:30776", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	log.Printf("Dialled auth")

	return &AuthInterceptor{
		authClient: pb.NewAuthServiceClient(conn),
		conn:       conn,
	}, nil
}

func (a *AuthInterceptor) newToken(ctx context.Context) error {
	token, err := a.authClient.GetToken(ctx, &pb.GetTokenRequest{})
	if err != nil {
		return err
	}
	a.cToken = token.GetToken()
	return nil
}

func (a *AuthInterceptor) AuthIntercept(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	user, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("unable to get user for auth: %w", err)
	}
	_, err = a.authClient.Authenticate(ctx, &pb.AuthenticateRequest{User: user.Username, Token: a.cToken})

	switch status.Code(err) {
	case codes.OK:
		return handler(ctx, req)
	case codes.PermissionDenied:
		// Request a new token
		err = a.newToken(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to get new token: %w", err)
		}
		_, err = a.authClient.Authenticate(ctx, &pb.AuthenticateRequest{User: user.Username, Token: a.cToken})
		if err != nil {
			return nil, fmt.Errorf("unable to auth with new token: %w", err)
		}
		return handler(ctx, req)
	default:
		return nil, fmt.Errorf("authentication error: %w", err)
	}
}
