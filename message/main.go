package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
)

func main() {

	err2 := publishThatScales()
	if err2 != nil {
		fmt.Println(err2)
	}

}

func publishThatScales() error {
	projectID := "key-rider-338615"
	topicID := "new_topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()
	
	var wg sync.WaitGroup
	var totalErrors uint64
	t := client.Topic(topicID)
	n := 1
	for i := 0; i < n; i++ {
		result := t.Publish(ctx, &pubsub.Message{
			Data: []byte("hello world"),
		})

		wg.Add(1)
		go func(i int, res *pubsub.PublishResult) {
			defer wg.Done()
			// The Get method blocks until a server-generated ID or
			// an error is returned for the published message.
			id, err := res.Get(ctx)
			if err != nil {
				// Error handling code can be added here.
				fmt.Println("Failed to publish", err)
				atomic.AddUint64(&totalErrors, 1)
				return
			}
			fmt.Println("Published message  msg ID\n", id)
		}(i, result)
	}

	wg.Wait()

	if totalErrors > 0 {
		return fmt.Errorf("%d of %d messages did not publish successfully", totalErrors, n)
	}
	return nil
}
