package main

import (
	"context"
	"log"
	"net"
	"slices"

	pb "github.com/ackieeee/grpc-sample/sample"
	"google.golang.org/grpc"
)

var (
  recipients = []*pb.Recipient{
    {
      RepientId: 1,
      Email: "test1@test.email.com",
    },
    {
      RepientId: 2,
      Email: "test2@test.email.com",
    },
    {
      RepientId: 3,
      Email: "test3@test.email.com",
    },
  }
)

type Recipient1Server struct {
  pb.UnimplementedRecipienterServer
}

func (s *Recipient1Server) GetRecipients(ctx context.Context, in *pb.GetRecipientsRequest) (*pb.GetRecipientsResponse, error) {
  targets := []*pb.Recipient{}
  for _, recipient := range recipients {
    recipientIds := in.GetRecipientIds()
    if slices.Index[int64](recipientIds, recipient.RepientId) > -1 {
      targets = append(targets, recipient)
    }
  }
  return *pb.GetRecipientsResponse{
    Recipient: targets,
  }, nil
}

func main() {
  lis, err := net.Listen("tcp", ":50090")
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	s := grpc.NewServer()
  server := &Recipient1Server{}
	pb.RegisterRecipienterServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: %v", err)
	}
}
