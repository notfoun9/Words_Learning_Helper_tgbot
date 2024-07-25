package event_consumer

import (
	"log"
	"telegram-bot/events"
	"time"
)

type Consumer struct {
	fetcher          events.Fetcher
	processor        events.Processor
	updatesBatchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:          fetcher,
		processor:        processor,
		updatesBatchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.updatesBatchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}
		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.HandleEvents(gotEvents); err != nil {
			log.Print(err)
			continue
		}
	}
}

func (c *Consumer) HandleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("cant handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
