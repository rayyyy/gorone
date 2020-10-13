package main

import (
	"fmt"
	"gorone/db"
	"gorone/lib/utils"
	"gorone/models"
	"time"

	"github.com/adjust/rmq/v3"
)

// MaxRetryCount リトライ回数
// TODO: リトライ処理
const MaxRetryCount = 5

// 制限時間で強制終了されるようにしたほうが良いか 30分とかで
// その場合は紐付けテーブルも30分で無効にするとか

func main() {
	// workerがどのリクエストを対応するか判断する処理
	ip := setup()
	req := models.FindRequest()
	if req == nil {
		panic("Request Not Found")
		// ECSのDesired Countを最適にしないといけない
	}
	models.NewAssignment(int(req.ID), ip)

	// redisコネクト
	errorChannel := make(chan error)
	connection, err := rmq.OpenConnection("gorone redis", "tcp", "host.docker.internal:6379", 0, errorChannel)
	if err != nil {
		panic(err)
	}
	taskQueue, err := connection.OpenQueue(req.RedisTagName)
	if err != nil {
		panic(err)
	}

	taskQueue.SetPushQueue(taskQueue) // リトライ用
	taskQueue.StartConsuming(1, time.Second*1)
	taskConsumer := &Consumer{}
	taskQueue.AddConsumer("consumer-"+ip, taskConsumer)

	// taskQueue.ReturnAllRejected() // rejectをreadyに戻す

	c := rmq.NewCleaner(connection)
	for _ = range time.Tick(time.Minute) {
		// TODO: リクエストの進行状況watchする終わってたらECS-1しを解除し、consuming stopする
		if req.IsDone() {
			taskQueue.StopConsuming()
		}
		c.Clean() // 定期的なゴミ掃除
	}
}

type Consumer struct {
	name string
}

func (consumer *Consumer) Consume(delivery rmq.Delivery) {
	fmt.Println(delivery.Payload())
	db := db.DbManager()
	result := models.CalcResult{KeyName: delivery.Payload()}
	db.FirstOrInit(&result, models.CalcResult{KeyName: result.KeyName})
	result.Result = 100
	db.Save(&result)
	delivery.Ack()
}

// init workerのIPアドレスを返す
func setup() string {
	db.Init()
	ip, err := utils.GetPublicIP()
	if err != nil {
		panic("ipが取得できませんでした")
	}
	fmt.Printf("Worker PublicIP: %v\n", ip)
	return ip
}
