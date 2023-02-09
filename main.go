// From DataWire Example service as base

package main

// NOTE: VERY WIP, DOES NOT WORK YET

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	envoyAuthV3 "github.com/datawire/ambassador/v2/pkg/api/envoy/service/auth/v3"

	"github.com/datawire/dlib/dhttp"

	server "authorizer/internal"
)

func main() {

	p := &server.AuthService{}

	grpcHandler := grpc.NewServer()
	envoyAuthV3.RegisterAuthorizationServer(grpcHandler, p)

	sc := &dhttp.ServerConfig{
		Handler: grpcHandler,
	}

	fmt.Print("starting...")
	log.Fatal(sc.ListenAndServe(context.Background(), ":3000"))
}
