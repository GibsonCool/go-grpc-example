package main

import (
	"context"
	"go-grpc-example/config"
	pb "go-grpc-example/proto"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	conn, err := grpc.Dial(config.PortStream, grpc.WithInsecure())

	if err != nil {
		log.Printf("stream grpc.Dial err:%v", err.Error())
	}

	defer conn.Close()

	client := pb.NewStreamServiceClient(conn)

	err = printLists(client, &pb.ReqStream{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: List", Value: 111}})
	if err != nil {
		log.Fatalf("printLists.err: %v", err)
	}

	err = printRecord(client, &pb.ReqStream{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Record", Value: 222}})
	if err != nil {
		log.Fatalf("printRecord.err: %v", err)
	}

	err = printRoute(client, &pb.ReqStream{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Route", Value: 333}})
	if err != nil {
		log.Fatalf("printRoute.err: %v", err)
	}

}

func printLists(client pb.StreamServiceClient, r *pb.ReqStream) error {
	stream, e := client.List(context.Background(), r)
	if e != nil {
		return e
	}

	for {
		respStream, e := stream.Recv()
		if e == io.EOF {
			break
		}
		if e != nil {
			return e
		}

		log.Printf("resp pt.name:%s,   pt.value:%d", respStream.Pt.Name, respStream.Pt.Value)
	}
	return nil
}

func printRecord(client pb.StreamServiceClient, r *pb.ReqStream) error {
	stream, e := client.Record(context.Background())
	if e != nil {
		return e
	}

	for n := 0; n < 8; n++ {
		r.Pt.Value = r.Pt.Value + int32(n)
		e := stream.Send(r)
		if e != nil {
			return e
		}
	}

	resp, e := stream.CloseAndRecv()
	if e != nil {
		return e
	}

	log.Printf("resp:pt.name:%s,  pt.value:%d", resp.Pt.Name, resp.Pt.Value)

	return nil
}

func printRoute(client pb.StreamServiceClient, r *pb.ReqStream) error {
	routeClient, e := client.Route(context.Background())
	if e != nil {
		return nil
	}

	for n := 0; n < 6; n++ {
		r.Pt.Value = r.Pt.Value + int32(n)
		e := routeClient.Send(r)
		if e != nil {
			return e
		}

		respStream, e := routeClient.Recv()
		if e == io.EOF {
			return nil
		}
		if e != nil {
			return e
		}

		log.Printf("接受到服务端发送的内容  pt.name:%v,  pt.value:%d", respStream.Pt.Name, respStream.Pt.Value)
	}

	routeClient.CloseSend()
	return nil
}
