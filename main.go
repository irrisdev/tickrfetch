package main

import (
	"github.com/irrisdev/go-coinlore/coinlore"
	"github.com/irrisdev/tickrfetch/logger"
	"github.com/irrisdev/tickrfetch/services"
)

func main() {

	client := coinlore.NewClient("https://api.coinlore.net")

	s := services.NewCoinService(client, nil)

	state := make(chan services.ServiceState)

	s.Start(state)

	for msg := range state {
		switch msg {
		case services.Running:
			logger.Info("Service is running")
		case services.Stopped:
			close(state)
		}
	}

}
