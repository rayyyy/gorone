package main

import (
	"fmt"
	"time"

	"github.com/adjust/rmq"
)

func main() {
	connection := rmq.OpenConnection("gorone redis", "tcp", "host.docker.internal:6379", 0)
	taskQueue := connection.OpenQueue("calc")
	taskQueue.SetPushQueue(taskQueue) // リトライ用
	c := rmq.NewCleaner(connection)

	taskQueue.StartConsuming(1, time.Second*1)
	taskConsumer := &Consumer{}
	taskQueue.AddConsumer("calc consumer", taskConsumer)

	// taskQueue.ReturnAllRejected() // rejectをreadyに戻す

	for _ = range time.Tick(time.Second) {
		c.Clean() // 定期的なゴミ掃除
	}
}

type Consumer struct {
	name string
}

func (consumer *Consumer) Consume(delivery rmq.Delivery) {
	fmt.Println("%+v", delivery.Payload())
	delivery.Ack()
}
