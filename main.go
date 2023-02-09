// From DataWire Example service as base

package main

// NOTE: VERY WIP, DOES NOT WORK YET

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	envoyAuthV3 "github.com/datawire/ambassador/v2/pkg/api/envoy/service/auth/v3"

	"github.com/andresftb/authorizer/internal/server"
	"github.com/datawire/dlib/dhttp"
)

func main() {
	test := &server.AuthService{}

	grpcHandler := grpc.NewServer()
	envoyAuthV3.RegisterAuthorizationServer(grpcHandler, test)

	sc := &dhttp.ServerConfig{
		Handler: grpcHandler,
	}

	fmt.Print("starting...")
	log.Fatal(sc.ListenAndServe(context.Background(), ":3000"))
}
