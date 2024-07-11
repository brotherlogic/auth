package server

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"

	pb "github.com/brotherlogic/auth/proto"
)

func generateUUID() string {
	id := uuid.New()
	return id.String()
}

func generateRandomNumber() string {
	return fmt.Sprintf("%v", rand.Int63())
}

func (s *Server) recycleToken() {

	// Only recycle token every 24 hours
	if time.Since(time.Unix(s.token.GetCreationDate(), 0)) < time.Hour*24 {
		return
	}

	s.token = &pb.AuthToken{
		CreationDate: time.Now().Unix(),
		Token:        fmt.Sprintf("%v-%v", generateUUID(), generateRandomNumber()),
	}
}
