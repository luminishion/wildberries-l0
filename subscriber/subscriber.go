package subscriber

import (
	"log"

	"github.com/luminishion/wildberries-l0/config"
	"github.com/luminishion/wildberries-l0/orders"

	"github.com/nats-io/stan.go"
)

var conn stan.Conn
var sub stan.Subscription

func handler(m *stan.Msg) {
	data := m.Data
	orders.Add(data)

	m.Ack()
}

func RunNats() {
	sc, err := stan.Connect(config.NatsClusterID, config.NatsClientID, stan.NatsURL(config.NatsURL))
	if err != nil {
		log.Fatal("subscriber run connect: ", err)
	}

	su, err := sc.Subscribe(config.NatsChannel, handler, stan.SetManualAckMode())
	if err != nil {
		log.Fatal("subscriber run subscribe: ", err)
	}

	sub = su
	conn = sc
}

func StopNats() {
	sub.Unsubscribe()
	conn.Close()
}
