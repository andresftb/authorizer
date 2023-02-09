package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"

	envoyCoreV3 "github.com/datawire/ambassador/v2/pkg/api/envoy/config/core/v3"
	envoyAuthV3 "github.com/datawire/ambassador/v2/pkg/api/envoy/service/auth/v3"
	envoyType "github.com/datawire/ambassador/v2/pkg/api/envoy/type/v3"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type AuthService struct{}

func (s *AuthService) Check(ctx context.Context, req *envoyAuthV3.CheckRequest) (*envoyAuthV3.CheckResponse, error) {

	requestURI, err := url.ParseRequestURI(req.GetAttributes().GetRequest().GetHttp().GetPath())
	if err != nil {
		fmt.Println("=> ERROR", err)
		return &envoyAuthV3.CheckResponse{
			Status: &status.Status{Code: int32(code.Code_UNKNOWN)},
			HttpResponse: &envoyAuthV3.CheckResponse_DeniedResponse{
				DeniedResponse: &envoyAuthV3.DeniedHttpResponse{
					Status: &envoyType.HttpStatus{Code: http.StatusInternalServerError},
					Headers: []*envoyCoreV3.HeaderValueOption{
						{Header: &envoyCoreV3.HeaderValue{Key: "Content-Type", Value: "application/json"}},
					},
					Body: `{"msg": "internal server error"}`,
				},
			},
		}, nil
	}

	res := req.GetAttributes().GetRequest().GetHttp().GetHeaders()
	Host := req.GetAttributes().GetRequest().GetHttp().GetHost()

	tokenString := res["authorization"]
	tokenString_clean := ""
	if tokenString != "" {
		tokenString_clean = tokenString[7:]
	}

	conf := &firebase.Config{
		ServiceAccountID: os.Getenv("SERVICE_ACCOUNT_ID"),
		ProjectID:        os.Getenv("PROJECT_ID"),
	}
	app, err := firebase.NewApp(context.Background(), conf)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	verifyIDToken(context.Background(), app, tokenString_clean, Host, requestURI.Path)

	return &envoyAuthV3.CheckResponse{
		Status: &status.Status{Code: int32(code.Code_OK)},
		HttpResponse: &envoyAuthV3.CheckResponse_OkResponse{
			OkResponse: &envoyAuthV3.OkHttpResponse{
				Headers: []*envoyCoreV3.HeaderValueOption{},
			},
		},
	}, nil

}

func verifyIDToken(ctx context.Context, app *firebase.App, idToken string, Host string, Path string) *auth.Token {
	// [START verify_id_token_golang]
	client, err := app.Auth(ctx)
	if err != nil {
		fmt.Printf("error getting Auth client: %v\n", err)
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		fmt.Printf("Host: %s, Path: %s, Error verifying ID token: %v\n", Host, Path, err)
	} else {
		fmt.Printf("Host: %s, Path: %s, Verified ID token: %v\n", Host, Path, token)
	}

	// [END verify_id_token_golang]

	return token
}
