package orders

import (
	"fmt"
	"log"

	"github.com/luminishion/wildberries-l0/orders/storage"

	"github.com/luminishion/wildberries-l0/config"
)

var st *storage.Storage

func Connect() {
	s, err := storage.New(config.DatabaseURL)
	if err != nil {
		log.Fatal("orders connect: ", err)
	}

	st = s
}

func List() []string {
	list, err := st.List()
	if err != nil {
		log.Fatal("orders list: ", err)
	}

	return list
}

func Add(orderJson []byte) {
	if err := st.Add(orderJson); err != nil {
		if err == storage.ErrBadData {
			fmt.Println("orders add bad data")
			return
		}

		log.Fatal("orders add: ", err)
	}
}

func Get(id string) []byte {
	orderJson, err := st.Get(id)
	if err != nil {
		log.Fatal("orders get: ", err)
	}

	return orderJson
}

func Close() {
	st.Close()
}
