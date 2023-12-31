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

const (
	  defaultName = "sample_name"
)

var (
	addr = flag.String("addr", "localhost:50061", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
