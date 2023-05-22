package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/amanakin/shortener/internal/handler/grpc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := api.NewShortenerClient(conn)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		args := strings.Split(text, " ")
		if len(args) != 2 {
			fmt.Printf("invalid command: %q\n", text)
			continue
		}

		switch args[0] {
		case "shorten":
			resp, err := client.Shorten(context.Background(), &api.ShortenRequest{
				Url: args[1],
			})
			if err != nil {
				fmt.Printf("error: %v\n", err)
				continue
			}
			fmt.Printf("original: %v\nshortened: %v\ncreated: %v\n",
				resp.Original, resp.Shortened, resp.Created)
		case "resolve":
			resp, err := client.Resolve(context.Background(), &api.ResolveRequest{
				Shortened: args[1],
			})
			if err != nil {
				fmt.Printf("error: %v\n", err)
				continue
			}
			fmt.Printf("resolved: %v\n", resp.Original)
		default:
			fmt.Printf("invalid command: %q\n", text)
		}
	}
}
