package interceptors

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"google.golang.org/api/option"
	"log"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var firebaseApp *firebase.App

func initFirebaseApp() (*firebase.App, error) {

	opt := option.WithCredentialsFile("chatapplictiontest-firebase-adminsdk-ohx08-4ac75a61e3.json")
	var err error
	firebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
		//log.Fatal(fmt.Errorf("error firebase initializing app: %v", err))
	}
	return firebaseApp, nil
}

func getFirebaseConn() (*firebase.App, error) {
	if firebaseApp == nil {
		return initFirebaseApp()
	}
	return firebaseApp, nil
}

func verifyUserToken(ctx context.Context, firebaseApp *firebase.App, idToken string) (*auth.Token, error) {
	client, err := firebaseApp.Auth(ctx)
	if err != nil {
		//log.Fatalf("error getting Auth client: %v\n", err)
		return nil, err
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Printf("error verifying ID token: %v\n", err)
		return nil, err

	}
	log.Printf("Verified ID token: %v\n", token)
	return token, nil
}

func AuthHandlingInterceptor(ctx context.Context) (context.Context, error) {

	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "no auth header in context")
	}
	fmt.Printf("metadata: %v\n", meta)
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	f, err := getFirebaseConn()
	if err != nil {
		return nil, status.Error(codes.Internal, "can't connect to firebase")
	}
	userToken, err := verifyUserToken(ctx, f, token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}
	fmt.Printf("User validation success %v\n", userToken)

	// WARNING: in production define your own type to avoid context collisions
	newCtx := context.WithValue(ctx, "isAuthenticated", true)
	return newCtx, nil
}
