package main

import (
	"context"
	"fmt"
	"net/http/httptest"

	simplev1 "example.com/connect-go-bug/gen"
	"example.com/connect-go-bug/gen/simplev1connect"
	"github.com/bufbuild/connect-go"
)

type SimpleServer struct {
	simplev1connect.UnimplementedSimpleServiceHandler
}

func (SimpleServer) Unary(context.Context, *connect.Request[simplev1.Empty]) (*connect.Response[simplev1.Empty], error) {
	response := connect.NewResponse(&simplev1.Empty{})
	response.Trailer()["lowercase"] = []string{"test"}
	return response, nil
}

func (SimpleServer) Stream(_ context.Context, _ *connect.Request[simplev1.Empty], stream *connect.ServerStream[simplev1.Empty]) error {
	stream.ResponseTrailer()["lowercase"] = []string{"test"}
	return nil
}

func main() {
	_, handler := simplev1connect.NewSimpleServiceHandler(&SimpleServer{})

	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	testClient := simplev1connect.NewSimpleServiceClient(testServer.Client(), testServer.URL)
	stream, err := testClient.Stream(context.Background(), connect.NewRequest(&simplev1.Empty{}))
	if err != nil {
		panic(err)
	}
	stream.Receive()
	response, err := testClient.Unary(context.Background(), connect.NewRequest(&simplev1.Empty{}))
	if err != nil {
		panic(err)
	}

	fmt.Println("Unary:")
	fmt.Printf("- Trailer().Get(\"lowercase\"): %s\n", response.Trailer().Get("lowercase"))
	fmt.Printf("- Trailer()[\"Lowercase\"][0]: %v\n", response.Trailer()["Lowercase"][0])

	fmt.Println("Streaming:")
	fmt.Printf("- Trailer().Get(\"lowercase\"): %s\n", stream.ResponseTrailer().Get("lowercase"))
	fmt.Printf("- Trailer()[\"lowercase\"][0]: %v\n", stream.ResponseTrailer()["lowercase"][0])
}
