package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"reflect"

	timeserver "github.com/itsshashank/gtimenow/model"
	"github.com/itsshashank/gtimenow/server/timer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Println("Unary Int invocked")
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("couldn't parse incoming context metadata")
		}
		log.Println("Meta data", md)
		h, err := handler(ctx, req)
		return h, err
	}
}

type EdgeServerStream struct {
	grpc.ServerStream
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := &EdgeServerStream{
			ServerStream: ss,
		}
		return handler(srv, wrapper)
	}
}
func (e *EdgeServerStream) RecvMsg(m interface{}) error {
	// Here we can perform additional logic on the received message, such as
	// validation
	log.Printf("intercepted server stream message, type: %s", reflect.TypeOf(m).String())
	if err := e.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	return nil
}

func (e *EdgeServerStream) SendMsg(m interface{}) error {
	// Here we can perform additional logic on the received message, such as
	// validation

	log.Printf("intercepted send server stream message, type: %s", reflect.TypeOf(m).String())
	if err := e.ServerStream.SendMsg(m); err != nil {
		return err
	}
	return nil
}

func main() {

	// Load the server certificate and its key
	// RSA
	// serverCert, err := tls.LoadX509KeyPair("/workspaces/mtlscert/cert/server-cert.pem",
	// 	"/workspaces/mtlscert/cert/server-key.pem")
	// ECDSA
	serverCert, err := tls.LoadX509KeyPair("/workspaces/mtlscert/CA/server.pem",
		"/workspaces/mtlscert/CA/server.key")
	if err != nil {
		log.Fatalf("Failed to load server certificate and key. %s.", err)
	}

	// Load the CA certificate
	trustedCert, err := ioutil.ReadFile("/workspaces/mtlscert/cert/ca-cert.pem")
	if err != nil {
		log.Fatalf("Failed to load trusted certificate. %s.", err)
	}
	// Put the CA certificate to certificate pool
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(trustedCert) {
		log.Fatalf("Failed to append trusted certificate to certificate pool. %s.", err)
	}
	trustedCert, err = ioutil.ReadFile("/workspaces/mtlscert/cert/ca-cert2.pem")
	if err != nil {
		log.Fatalf("Failed to load trusted certificate. %s.", err)
	}
	clientPool := x509.NewCertPool()
	if !clientPool.AppendCertsFromPEM(trustedCert) {
		log.Fatalf("Failed to append trusted certificate to certificate pool. %s.", err)
	}
	// Create the TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		RootCAs:      certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientPool,
	}

	// Create a new TLS credentials based on the TLS configuration
	cred := credentials.NewTLS(tlsConfig)

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panic(err)
	}
	server := grpc.NewServer(grpc.Creds(cred),
		grpc.UnaryInterceptor(UnaryServerInterceptor()),
		grpc.StreamInterceptor(StreamServerInterceptor()),
	)
	timeserver.RegisterTimeServerServer(server, &timer.TimeServer{})
	reflection.Register(server)
	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
