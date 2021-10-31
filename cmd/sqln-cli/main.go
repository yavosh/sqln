// package main is a simple shell for the sqln api
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	pb "github.com/yavosh/sqln/proto"
	"google.golang.org/grpc"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	var (
		endpoint string
	)

	flag.StringVar(&endpoint, "endpoint", "localhost:5051", "connect to db")
	flag.Parse()

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error connecting to sqln @ %s %v", endpoint, err)
		return
	}

	fmt.Printf("Connected to %q\n", endpoint)

	client := pb.NewQueryClient(conn)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.ToUpper(strings.TrimRight(cmd, "\n"))

		if strings.TrimSpace(cmd) == "" {
			continue
		}

		if cmd == ".q" {
			fmt.Println("bye")
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		req := pb.QueryRequest{Query: cmd}

		result, err := client.ExecuteQuery(ctx, &req)

		if err != nil {
			fmt.Println("err:", err)
			continue
		}

		fmt.Printf("Result: %+v\n", result.Columns)
		fmt.Printf("Result: %+v\n", result.Rows)
	}
}
