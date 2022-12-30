package routes

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/luminishion/wildberries-l0/config"
	"github.com/luminishion/wildberries-l0/orders"
	"github.com/luminishion/wildberries-l0/orders/storage"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/stan.go"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr() string {
	n := rand.Intn(16) + 8
	buf := make([]rune, n)

	for i := range buf {
		buf[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(buf)
}

func postGenerate2(c *gin.Context) {
	d := []byte(randStr())
	orders.Add(d)
	c.JSON(http.StatusOK, gin.H{"data": d})
}

func postGenerate(c *gin.Context) {
	var id string
	for {
		id = randStr()
		if orders.Get(id) == nil {
			break
		}
	}

	sc, _ := stan.Connect(config.NatsClusterID, "test", stan.NatsURL(config.NatsURL))
	defer sc.Close()

	order := storage.Order{
		Id:          id,
		TrackNumber: randStr(),
		Entry:       randStr(),

		Delivery: storage.Delivery{
			Name:    randStr(),
			Phone:   randStr(),
			Zip:     randStr(),
			City:    randStr(),
			Address: randStr(),
			Region:  randStr(),
			Email:   randStr(),
		},
		Payment: storage.Payment{
			Transaction:  randStr(),
			RequestId:    randStr(),
			Currency:     randStr(),
			Provider:     randStr(),
			Amount:       rand.Intn(1337),
			PaymentDt:    rand.Intn(1337),
			Bank:         randStr(),
			DeliveryCost: rand.Intn(1337),
			GoodsTotal:   rand.Intn(1337),
			CustomFee:    rand.Intn(1337),
		},
		Items: []storage.Item{},

		Locale:           randStr(),
		InternalSignture: randStr(),
		CustomerId:       randStr(),
		DeliveryService:  randStr(),
		Shardkey:         randStr(),
		SmId:             rand.Intn(1337),
		DateCreated:      randStr(),
		OofShard:         randStr(),
	}

	for i := 0; i < rand.Intn(3)+1; i++ {
		item := storage.Item{
			OrderId: order.Id,

			ChrtId:      rand.Intn(1337),
			TrackNumber: randStr(),
			Price:       rand.Intn(1337),
			Rid:         randStr(),
			Name:        randStr(),
			Sale:        rand.Intn(1337),
			Size:        randStr(),
			TotalPrice:  rand.Intn(1337),
			NmId:        rand.Intn(1337),
			Brand:       randStr(),
			Status:      rand.Intn(1337),
		}
		order.Items = append(order.Items, item)
	}

	orderJson, _ := json.Marshal(order)
	sc.Publish(config.NatsChannel, orderJson)

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func getGenerate(c *gin.Context) {
	c.HTML(http.StatusOK, "generate.html", gin.H{})
}

func Generate(r *gin.Engine) {
	r.POST("/generate", postGenerate)
	r.POST("/generate2", postGenerate2)

	r.GET("/generate", getGenerate)
}
