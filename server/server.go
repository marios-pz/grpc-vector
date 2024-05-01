package main

import (
	"context"
	"errors"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/marios-pz/grpc-vector/pb"
)

type server struct {
	pb.UnimplementedVectorServer
}

func (s *server) InnerProduct(ctx context.Context, in *pb.VectorInput) (*pb.VectorProductResult, error) {
	if in.X == nil || in.Y == nil || len(in.X) != len(in.Y) {
		return nil, errors.New("x or y was not initialized or lengths do not match")
	}

	var product pb.VectorProductResult
	for i := 0; i < len(in.X); i++ {
		product.Result += in.X[i] * in.Y[i]
	}

	return &product, nil
}

func (s *server) AverageValues(ctx context.Context, in *pb.VectorInput) (*pb.VectorResult, error) {
	if in.X == nil || in.Y == nil || len(in.X) != len(in.Y) {
		return nil, status.Error(codes.InvalidArgument, "x or y was not initialized or lengths do not match")
	}

	var (
		result pb.VectorResult
		sumX   int64 = 0
		sumY   int64 = 0
	)

	for _, num := range in.X {
		sumX += int64(num)
	}

	result.A = float32(sumX) / float32(len(in.X))

	for _, num := range in.Y {
		sumY += int64(num)
	}

	result.B = float32(sumY) / float32(len(in.Y))

	return &result, nil
}

func (s *server) ScalarVectorProduct(ctx context.Context, in *pb.VectorInput) (*pb.VectorResult, error) {
	if in.X == nil || in.Y == nil || len(in.X) != len(in.Y) {
		return nil, status.Error(codes.InvalidArgument, "x or y was not initialized or lengths do not match")
	}
	if in.R <= 0 {
		return nil, status.Error(codes.InvalidArgument, "r needs to be > 0")
	}

	result := &pb.VectorResult{}

	for i := 0; i < len(in.X); i++ {
		sumXY := in.X[i] + in.Y[i]
		result.A += in.R * float32(sumXY)
		result.B += in.R * float32(sumXY)
	}

	return result, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("could not create listener")
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterVectorServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalln("failed to serve", err)
	}
}
