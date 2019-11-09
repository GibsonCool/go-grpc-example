package simple_server

import (
	"context"
	"go-grpc-example/config"
	pb "go-grpc-example/proto"
	"google.golang.org/grpc"
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
	server := grpc.NewServer()
	pb.RegisterSearchServiceServer(server, &SearchService{})

	listener, e := net.Listen("tcp", config.PORT)

	if e != nil {
		log.Printf("net.Listen err:%s", e.Error())
	}

	server.Serve(listener)
}
