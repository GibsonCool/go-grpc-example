package main

import (
	"context"
	"go-grpc-example/config"
	pb "go-grpc-example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

func main() {
	tc, e := credentials.NewClientTLSFromFile(config.ServerPemPath, "go-grpc-example")
	if e != nil {
		log.Printf("credentials.NewServerTLSFromFile err:%v", e.Error())
	}

	conn, e := grpc.Dial(config.PORT, grpc.WithTransportCredentials(tc))

	if e != nil {
		log.Printf("grpc.Dial err:%v", e.Error())
	}

	defer conn.Close()

	client := pb.NewSearchServiceClient(conn)
	resp, e := client.Search(context.Background(), &pb.SearchRequest{Request: "gRPC  testing"})
	if e != nil {
		log.Printf("client.Search err:%v", e.Error())
	}

	log.Printf("resp:  %s", resp.GetResponse())
}
