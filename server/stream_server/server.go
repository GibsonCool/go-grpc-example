package main

import (
	"go-grpc-example/config"
	pb "go-grpc-example/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type StreamService struct {
}

func (StreamService) List(r *pb.ReqStream, stream pb.StreamService_ListServer) error {
	for n := 0; n < 6; n++ {
		err := stream.Send(&pb.RespStream{
			Pt: &pb.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(n),
			},
		})
		log.Printf("test:%d\n", n)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (StreamService) Record(stream pb.StreamService_RecordServer) error {
	for {
		reqStream, e := stream.Recv()
		if e == io.EOF {
			return stream.SendAndClose(&pb.RespStream{Pt: &pb.StreamPoint{Name: "服务端接收完毕，请求关闭", Value: int32(666)}})
		}

		if e != nil {
			return e
		}
		log.Printf("stream.Recv pt.name:%v,  pt.value:%d", reqStream.Pt.Name, reqStream.Pt.Value)
	}
	return nil
}

func (StreamService) Route(stream pb.StreamService_RouteServer) error {
	n := 0
	for {
		e := stream.Send(&pb.RespStream{
			Pt: &pb.StreamPoint{
				Name:  "route 服务端 发送的内容",
				Value: int32(n),
			},
		})
		if e != nil {
			return e
		}

		reqStream, e := stream.Recv()
		if e == io.EOF {
			return nil
		}
		if e != nil {
			return e
		}

		n++
		log.Printf("pt.name:%v,  pt.value:%d", reqStream.Pt.Name, reqStream.Pt.Value)
	}

	return nil
}

func main() {
	// 创建rpc服务
	server := grpc.NewServer()
	// 注册要提供的方法
	pb.RegisterStreamServiceServer(server, &StreamService{})

	// 开启监听
	listener, e := net.Listen("tcp", config.PortStream)
	if e != nil {
		log.Printf("stream net listen err:%v", e.Error())
	}

	server.Serve(listener)
}
