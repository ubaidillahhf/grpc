package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	pb "github.com/ubaidillahhf/grpc/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChatServiceServer
}

// 1. Unary RPC: Simple request-response
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received Unary request from: %s", req.Name)
	response := &pb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s! Welcome to gRPC!", req.Name),
	}
	return response, nil
}

// 2. Server Streaming RPC: Send multiple responses to client
func (s *server) GetServerStream(req *pb.StreamRequest, stream pb.ChatService_GetServerStreamServer) error {
	log.Printf("Server Streaming started for: %s", req.Message)
	
	messages := []string{
		"First update: Processing your request...",
		"Second update: Fetching data...",
		"Third update: Almost done...",
		"Fourth update: Finalizing...",
		"Fifth update: Complete!",
	}
	
	for i, msg := range messages {
		response := &pb.StreamResponse{
			Message: fmt.Sprintf("[%d] %s", i+1, msg),
		}
		if err := stream.Send(response); err != nil {
			return err
		}
		log.Printf("Sent: %s", response.Message)
		time.Sleep(1 * time.Second)
	}
	
	return nil
}

// 3. Client Streaming RPC: Receive multiple requests, send one response
func (s *server) GetClientStream(stream pb.ChatService_GetClientStreamServer) error {
	log.Println("Client Streaming started")
	
	var messages []string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			response := &pb.StreamResponse{
				Message: fmt.Sprintf("Received %d messages: %v", len(messages), messages),
			}
			log.Printf("Client streaming completed. Total messages: %d", len(messages))
			return stream.SendAndClose(response)
		}
		if err != nil {
			return err
		}
		log.Printf("Received from client: %s", req.Message)
		messages = append(messages, req.Message)
	}
}

// 4. Bidirectional Streaming RPC: Both send and receive multiple messages
func (s *server) GetBidirectionalStream(stream pb.ChatService_GetBidirectionalStreamServer) error {
	log.Println("Bidirectional Streaming started")
	
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("Client closed the stream")
			return nil
		}
		if err != nil {
			return err
		}
		
		log.Printf("Received: %s", req.Message)
		
		response := &pb.StreamResponse{
			Message: fmt.Sprintf("Echo: %s (received at %s)", req.Message, time.Now().Format("15:04:05")),
		}
		
		if err := stream.Send(response); err != nil {
			return err
		}
		log.Printf("Sent: %s", response.Message)
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, &server{})
	
	log.Println("🚀 gRPC Server is running on port :50051")
	log.Println("Waiting for client connections...")
	
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
