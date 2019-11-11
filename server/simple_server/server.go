package main

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go-grpc-example/config"
	pb "go-grpc-example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"runtime/debug"
)

type SearchService struct {
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Printf("接受到client 的请求内容：%v", r.GetRequest())
	return &pb.SearchResponse{Response: "this response come from Server"}, nil
}

func main() {
	tc, e := credentials.NewServerTLSFromFile(config.ServerPemPath, config.ServerKeyPath)
	if e != nil {
		log.Printf("credentials.NewServerTLSFromFile err:%v", e.Error())
	}

	opts := []grpc.ServerOption{
		grpc.Creds(tc),
		grpc_middleware.WithUnaryServerChain(
			RecoverInterceptor,
			LoggingInterceptor,
		),
	}

	server := grpc.NewServer(opts...)
	pb.RegisterSearchServiceServer(server, &SearchService{})

	listener, e := net.Listen("tcp", config.PORT)

	if e != nil {
		log.Printf("net.Listen err:%s", e.Error())
	}

	server.Serve(listener)
}

//gRPC 拦截器

/*
	RPC 方法的入参出参的日志输出

		ctx context.Context：请求上下文
		req interface{}：RPC 方法的请求参数
		info *UnaryServerInfo：RPC 方法的所有信息
		handler UnaryHandler：RPC 方法本身
*/
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Printf("gRPC 方法： %s,  %v\n", info.FullMethod, req)
	resp, err = handler(ctx, req)
	log.Printf("gRPC 方法结果： %s,  %v\n", info.FullMethod, resp)
	return resp, err
}

/*
	RPC 方法的异常保护和日志输出
*/
func RecoverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "调用出错信息：%v", e)
		}
	}()

	return handler(ctx, req)
}
