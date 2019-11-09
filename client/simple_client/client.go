package simple_client

import (
	"context"
	"go-grpc-example/config"
	pb "go-grpc-example/proto"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, e := grpc.Dial(config.PORT, grpc.WithInsecure())

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
