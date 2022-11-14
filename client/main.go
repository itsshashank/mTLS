package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	timeserver "github.com/itsshashank/gtimenow/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	// Load the client certificate and its key
	clientCert, err := tls.LoadX509KeyPair("/workspaces/mtlscert/cert/client-cert2.pem",
		"/workspaces/mtlscert/cert/client-key2.pem")
	if err != nil {
		log.Fatalf("Failed to load client certificate and key. %s.", err)
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

	// Create the TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	// Create a new TLS credentials based on the TLS configuration
	cred := credentials.NewTLS(tlsConfig)
	log.Println(cred)

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(cred))
	// conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Panic(err)
	}

	client := timeserver.NewTimeServerClient(conn)
	stream, err := client.TimeNow(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	// ans := "1"
	// for {
	// 	stream.Send(&timeserver.Response{Word: ans})
	// 	in, err := client.TimeNow(context.Background())
	// 	if err == io.EOF {
	// 		return
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("Failed to receive a note : %v", err)
	// 	}
	// 	in.RecvMsg(ans)
	// 	log.Printf("Got message %v", ans)
	// 	time.Sleep(10 * time.Second)
	// }
	var max int32
	ctx := stream.Context()
	done := make(chan bool)

	// first goroutine sends random increasing numbers to stream
	// and closes it after 10 iterations
	go func() {
		for i := 1; i <= 10; i++ {
			// generates random number and sends it to stream
			rnd := int32(rand.Intn(i))
			req := timeserver.Response{Word: int64(rnd)}
			if err := stream.Send(&req); err != nil {
				log.Fatalf("can not send %v", err)
			}
			log.Printf("%d sent", req.Word)
			time.Sleep(time.Millisecond * 200)
		}
		if err := stream.CloseSend(); err != nil {
			log.Println(err)
		}
	}()

	// second goroutine receives data from stream
	// and saves result in max variable
	//
	// if stream is finished it closes done channel
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(done)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			max = int32(resp.GetWord())
			log.Printf("new max %d received", max)
		}
	}()

	// third goroutine closes done channel
	// if context is done
	go func() {
		<-ctx.Done()
		if err := ctx.Err(); err != nil {
			log.Println(err)
		}
		close(done)
	}()

	<-done
	log.Printf("finished with max=%d", max)
}
