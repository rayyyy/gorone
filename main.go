package main

import (
	"fmt"
	"gorone/db"
	"gorone/models"
	"time"

	"github.com/adjust/rmq"
)

const MAX_COUNT = 10

var nowCount = 0

func main() {
	db.Init()
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
	nowCount = nowCount + 1
	if nowCount >= MAX_COUNT {
		panic("over max count 10")
	}
	fmt.Println(delivery.Payload())
	db := db.DbManager()
	result := models.CalcResult{KeyName: delivery.Payload()}
	db.FirstOrInit(&result, models.CalcResult{KeyName: result.KeyName})
	result.Result = 100
	db.Save(&result)
	delivery.Ack()
}
