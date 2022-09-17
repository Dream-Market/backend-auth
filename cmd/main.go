package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Dream-Market/backend-auth/pkg/config"
	"github.com/Dream-Market/backend-auth/pkg/db"
	"github.com/Dream-Market/backend-auth/pkg/pb"
	"github.com/Dream-Market/backend-auth/pkg/services"
	"github.com/Dream-Market/backend-auth/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c)

	jwt := utils.JwtWrapper{
		SecretKey: c.JWTSecretKey,
		Issuer:    "auth-srvc",
	}

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Auth Srvc on", c.Port)

	s := services.Server{
		H:   h,
		Jwt: jwt,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
