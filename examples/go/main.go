package main

import (
	"context"
	"fmt"
	clientpkg "github.com/skibid991/httpresponder.com/sdk/go"
	"log"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	ctx := context.Background()
	client := clientpkg.New()
	r, err := client.Decode(ctx, "SGVsbG8sIHdvcmxkIQ==", "")
	must(err)
	fmt.Println(r.Text())
	r, err = client.Status(ctx, 418)
	must(err)
	fmt.Println("Requested status:", r.StatusCode)
}
