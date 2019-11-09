package main

import (
	"context"
	"go-grpc-example/config"
	pb "go-grpc-example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

type SearchService struct {
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Printf("接受到client 的请求内容：%v", r.GetRequest())
	return &pb.SearchResponse{Response: "这是服务端返回内容：Server"}, nil
}

func main() {
	tc, e := credentials.NewServerTLSFromFile(config.ServerPemPath, config.ServerKeyPath)
	if e != nil {
		log.Printf("credentials.NewServerTLSFromFile err:%v", e.Error())
	}

	server := grpc.NewServer(grpc.Creds(tc))
	pb.RegisterSearchServiceServer(server, &SearchService{})

	listener, e := net.Listen("tcp", config.PORT)

	if e != nil {
		log.Printf("net.Listen err:%s", e.Error())
	}

	server.Serve(listener)
}
