package events

import (
	"context"
	"fmt"
	"log"
	"time"
	"townwatch/base"
)

func LoadInit(b *base.Base) {

	count, errC := b.Queries.CountEvents(context.Background())
	if errC != nil {
		log.Fatal(errC)
	}

	if count > 0 {
		return
	}

	count, err := fetchAndStoreEvents(b, context.Background(), time.Now().Add(-time.Duration(24*7)*time.Hour).UTC(), time.Now().UTC())
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err := fmt.Errorf("no police reports found within 7 days")
		log.Fatal(err)
	}

}
