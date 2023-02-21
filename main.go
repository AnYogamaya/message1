package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
)

type Details struct {
	Firstname string
	Lastname  string
}

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
	var data Details
	var wg sync.WaitGroup
	var totalErrors uint64
	t := client.Topic(topicID)
	n := 1
	data.Firstname = "maya"
	data.Lastname = "mishra"
	empdata, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error", err)
	}

	fmt.Println(empdata)

	for i := 0; i < n; i++ {
		result := t.Publish(ctx, &pubsub.Message{
			Data: []byte(empdata),
		})

		wg.Add(1)
		go func(i int, res *pubsub.PublishResult) {
			defer wg.Done()

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
