package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"slices"
	"sync"

	pb "github.com/ackieeee/grpc-sample/sample"
	"google.golang.org/grpc"
)

var (
	recipients = []*pb.Recipient{
		{
			RepientId: 1,
			Email:     "test1@test.email.com",
		},
		{
			RepientId: 2,
			Email:     "test2@test.email.com",
		},
		{
			RepientId: 3,
			Email:     "test3@test.email.com",
		},
	}
	servers = []struct {
		server pb.RecipienterServer
		port   string
	}{
		{
			server: &Recipient1Server{},
			port:   ":50090",
		},
		{
			server: &Recipient2Server{},
			port:   ":50093",
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
		if slices.Contains[[]int64, int64](recipientIds, recipient.RepientId) {
			targets = append(targets, recipient)
		}
	}
	return &pb.GetRecipientsResponse{
		Recipient: targets,
	}, nil
}

type Recipient2Server struct {
	pb.UnimplementedRecipienterServer
}

func (s *Recipient2Server) GetRecipients(ctx context.Context, in *pb.GetRecipientsRequest) (*pb.GetRecipientsResponse, error) {
	return &pb.GetRecipientsResponse{
		Recipient: []*pb.Recipient{
			{
				RepientId: 100,
				Email:     "test",
			},
			{
				RepientId: 200,
				Email:     "test2",
			},
		},
	}, nil
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	for _, server := range servers {
		sv := server
		go func() {
			fmt.Println("start listen ", sv.port)
			lis, err := net.Listen("tcp", sv.port)
			if err != nil {
				log.Fatal("failed to listen: %v", err)
			}
			s := grpc.NewServer()
			pb.RegisterRecipienterServer(s, server.server)
			log.Printf("server listening at %v", lis.Addr())
			if err := s.Serve(lis); err != nil {
				log.Fatal("failed to serve: %v", err)
				wg.Done()
			}
		}()
	}
	wg.Wait()
}
