package main

import (
	"fmt"
	"time"

	"github.com/irrisdev/go-coinlore/coinlore"
	"github.com/irrisdev/tickrfetch/services"
)

func main() {


	client := coinlore.NewClient("https://api.coinlore.net")
	
	_ = services.NewCoinService(client, nil)
	
	start := time.Now()

	fmt.Println(time.Since(start).Seconds())

	
}
