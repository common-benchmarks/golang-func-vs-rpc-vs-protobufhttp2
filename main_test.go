package main

import (
	pb "github.com/common-benchmarks/golang-func-vs-rpc-vs-protobufhttp2/protobufs"
	"log"
	"net"
	"os"
	"testing"

	"github.com/valyala/gorpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	globalString  string
	greeterClient pb.GreeterClient
)

func tempFuncStr() string {
	return "2006"
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type grpcServer struct{}

func (g *grpcServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func TestMain(m *testing.M) {
	s := &gorpc.Server{
		Addr: ":55555",
		Handler: func(clientAddr string, request interface{}) interface{} {
			//			log.Printf("Obtained request %+v from the client %s\n", request, clientAddr)
			return tempFuncStr()
		},
	}
	go func() {
		if err := s.Serve(); err != nil {
			log.Fatalf("Cannot start rpc server: %s", err.Error())
		}
	}()

	//https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_server/main.go
	lis, err := net.Listen("tcp", ":55556")
	checkError(err)

	s2 := grpc.NewServer()
	pb.RegisterGreeterServer(s2, &grpcServer{})
	go func() {
		err = s2.Serve(lis)
		checkError(err)
	}()

	//https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_client/main.go
	address := "localhost:55556"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	checkError(err)
	defer conn.Close()

	greeterClient = pb.NewGreeterClient(conn)

	os.Exit(m.Run())
}

func BenchmarkProtobufRpcCall(b *testing.B) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("ERROR: %+v", r)
		}
	}()

	var g string

	for n := 0; n < b.N; n++ {
		r, err := greeterClient.SayHello(context.Background(), &pb.HelloRequest{Name: tempFuncStr()})
		checkError(err)

		g = r.Message
	}

	globalString = g
}

func BenchmarkRpcCall(b *testing.B) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("ERROR: %+v", r)
		}
	}()

	var g string

	c := &gorpc.Client{
		// TCP address of the server.
		Addr: "localhost:55555",
	}
	c.Start()

	for n := 0; n < b.N; n++ {
		resp, err := c.Call("foobar")
		checkError(err)

		g = resp.(string)
	}

	// log.Printf("Client stats: %+v", c.Stats)
	globalString = g
}

func BenchmarkNormalFunction(b *testing.B) {
	var g string
	for n := 0; n < b.N; n++ {
		g = tempFuncStr()
	}
	globalString = g
}
