package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/ubaidillahhf/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	
	client := pb.NewChatServiceClient(conn)
	
	fmt.Println("🎯 gRPC Client Connected!")
	fmt.Println("=" + string(make([]byte, 50)) + "=")
	
	// 1. Unary RPC Demo
	fmt.Println("\n1️⃣  UNARY RPC (Simple Request-Response)")
	fmt.Println("   Like asking: 'What's your name?' and getting one answer")
	unaryRPC(client)
	
	time.Sleep(2 * time.Second)
	
	// 2. Server Streaming RPC Demo
	fmt.Println("\n2️⃣  SERVER STREAMING RPC (One Request, Multiple Responses)")
	fmt.Println("   Like subscribing to weather updates - you ask once, get many updates")
	serverStreamingRPC(client)
	
	time.Sleep(2 * time.Second)
	
	// 3. Client Streaming RPC Demo
	fmt.Println("\n3️⃣  CLIENT STREAMING RPC (Multiple Requests, One Response)")
	fmt.Println("   Like uploading photos - send many, get one confirmation")
	clientStreamingRPC(client)
	
	time.Sleep(2 * time.Second)
	
	// 4. Bidirectional Streaming RPC Demo
	fmt.Println("\n4️⃣  BIDIRECTIONAL STREAMING RPC (Multiple Requests & Responses)")
	fmt.Println("   Like a live chat - both sides can talk anytime")
	bidirectionalStreamingRPC(client)
	
	fmt.Println("\n✅ All demos completed!")
}

func unaryRPC(client pb.ChatServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	req := &pb.HelloRequest{Name: "Alice"}
	res, err := client.SayHello(ctx, req)
	if err != nil {
		log.Fatalf("Unary RPC failed: %v", err)
	}
	
	fmt.Printf("   📨 Sent: Hello from %s\n", req.Name)
	fmt.Printf("   📬 Received: %s\n", res.Message)
}

func serverStreamingRPC(client pb.ChatServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	req := &pb.StreamRequest{Message: "Give me updates"}
	stream, err := client.GetServerStream(ctx, req)
	if err != nil {
		log.Fatalf("Server Streaming failed: %v", err)
	}
	
	fmt.Printf("   📨 Sent request: %s\n", req.Message)
	fmt.Println("   📬 Receiving stream:")
	
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}
		fmt.Printf("      ➜ %s\n", res.Message)
	}
}

func clientStreamingRPC(client pb.ChatServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	stream, err := client.GetClientStream(ctx)
	if err != nil {
		log.Fatalf("Client Streaming failed: %v", err)
	}
	
	messages := []string{"Message 1", "Message 2", "Message 3", "Message 4", "Message 5"}
	
	fmt.Println("   📨 Sending multiple messages:")
	for _, msg := range messages {
		req := &pb.StreamRequest{Message: msg}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error sending: %v", err)
		}
		fmt.Printf("      ➜ Sent: %s\n", msg)
		time.Sleep(500 * time.Millisecond)
	}
	
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error closing stream: %v", err)
	}
	
	fmt.Printf("   📬 Received final response: %s\n", res.Message)
}

func bidirectionalStreamingRPC(client pb.ChatServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	stream, err := client.GetBidirectionalStream(ctx)
	if err != nil {
		log.Fatalf("Bidirectional Streaming failed: %v", err)
	}
	
	messages := []string{"Hello", "How are you?", "This is cool!", "Goodbye"}
	
	waitc := make(chan struct{})
	
	// Goroutine to receive messages
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Error receiving: %v", err)
			}
			fmt.Printf("   📬 Received: %s\n", res.Message)
		}
	}()
	
	// Send messages
	fmt.Println("   📨 Sending messages:")
	for _, msg := range messages {
		req := &pb.StreamRequest{Message: msg}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error sending: %v", err)
		}
		fmt.Printf("      ➜ Sent: %s\n", msg)
		time.Sleep(1 * time.Second)
	}
	
	stream.CloseSend()
	<-waitc
}
