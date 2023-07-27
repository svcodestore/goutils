package grpc

import (
	"context"
	"encoding/json"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type CommonResp struct {
	Result *structpb.Struct

	GetResult func() *structpb.Struct
}

const (
	AuthTokenKey string = "authorization"
)

func ToRpcStruct(data any) *structpb.Struct {
	b, e := json.Marshal(data)
	if e != nil {
		return nil
	}
	var m map[string]any
	e = json.Unmarshal(b, &m)

	if e != nil {
		return nil
	}
	s, err := structpb.NewStruct(m)
	if err != nil {
		return nil
	}

	return s
}

func RunServer(addr string, register func(grpcServer *grpc.Server), auth func(token string) bool) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 10 * time.Second,
		}),
		grpc.UnaryInterceptor(func(ctx context.Context,
			req interface{},
			info *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler) (resp interface{}, err error) {

			md, ok := metadata.FromIncomingContext(ctx)
			authToken := md.Get(AuthTokenKey)
			if !ok || (len(authToken) < 1) || (len(authToken) > 0 && auth(authToken[0])) {
				return resp, status.Error(codes.Unauthenticated, "")
			}

			return handler(ctx, req)
		}),
	)

	register(grpcServer)

	return grpcServer.Serve(listener)
}

func NewClient(addr string, authToken string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			func(ctx context.Context,
				method string,
				req,
				reply interface{},
				cc *grpc.ClientConn,
				invoker grpc.UnaryInvoker,
				opts ...grpc.CallOption) error {

				if authToken != "" {
					md := metadata.New(map[string]string{
						AuthTokenKey: authToken,
					})

					ctx = metadata.NewOutgoingContext(ctx, md)
				}

				err := invoker(ctx, method, req, reply, cc, opts...)

				if err != nil {
					return err
				}

				if resp, ok := reply.(*CommonResp); ok {
					resp.Result = ToRpcStruct(map[string]any{
						"data": reply.(*CommonResp).GetResult().AsMap()["data"],
					})
				}

				return nil
			}))

	if err != nil {
		return nil, err
	}

	return conn, nil
}
