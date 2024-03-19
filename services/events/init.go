package events

import (
	"context"
	"fmt"
	"log"
	"time"
	"townwatch/base"
)

func LoadInit(b *base.Base) {
	fetchedCount, err := FetchAndStoreTorontoEvents(b, context.Background(), time.Now().Add(-time.Duration(24*7)*time.Hour).UTC(), time.Now().UTC())
	if err != nil {
		log.Fatal(err)
	}
	if fetchedCount == 0 {
		err := fmt.Errorf("no police reports found within 7 days")
		log.Fatal(err)
	}
	fmt.Println("==========================")
	fmt.Println(fmt.Sprintf("Initial police reports added | count: %v", fetchedCount))
	fmt.Println("==========================")
}
