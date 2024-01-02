package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/ackieeee/grpc-sample/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50093", "the address to connect to")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRecipienterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	request := pb.GetRecipientsRequest{
		RecipientIds: []int64{1, 2, 4, 3},
	}
	res, err := c.GetRecipients(ctx, &request)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	recipients := res.GetRecipient()
	for _, recipient := range recipients {
		log.Printf("res: %#v", recipient.RepientId)
	}
}
