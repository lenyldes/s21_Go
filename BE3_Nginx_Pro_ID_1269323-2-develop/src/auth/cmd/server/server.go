package main

import (
	"auth/api"
	"auth/internal/auth"
	"auth/internal/data"
	"auth/internal/shared"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func grpcServer() *grpc.Server {
	return grpc.NewServer()
}

func Run(authService *auth.AuthService, config *shared.Config) {
	lis, err := net.Listen("tcp", ":"+config.AppPort)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterAuthServiceServer(grpcServer, authService)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

}

func main() {
	fx.New(
		data.Module(),
		auth.Module(),
		shared.Module(),
		fx.Invoke(Run),
		// fx.NopLogger,
	)
}
