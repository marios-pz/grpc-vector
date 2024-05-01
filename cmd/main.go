// gRPC Client
package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/marios-pz/grpc-vector/pb"
)

func main() {
	serverAddr := flag.String("server", "localhost:8080", "host:port")
	flag.Parse()

	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("could not connect to grpc server")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client := pb.NewVectorClient(conn)

	bucketA := []int64{1, 3, 5, 7, 9}
	bucketB := []int64{2, 4, 6, 8, 10}

	inputVector := &pb.VectorInput{X: bucketA, Y: bucketB, R: 2}

	res, err := client.InnerProduct(ctx, inputVector)
	if err != nil {
		log.Fatalln("could not send inner product")
	}
	log.Println(res)

	res1, err := client.AverageValues(ctx, inputVector)
	if err != nil {
		log.Fatalln("could not send average values")
	}
	log.Println(res1)

	res1, err = client.ScalarVectorProduct(ctx, inputVector)
	if err != nil {
		log.Fatalln("could not send scalar values")
	}
	log.Println(res1)
}
