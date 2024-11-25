package services

import (
	"sync"

	"github.com/irrisdev/go-coinlore/coinlore"
	"golang.org/x/time/rate"
)

type CoinService struct {
	client  *coinlore.Client
	limiter *rate.Limiter
	cache   sync.Map
	symbols sync.Map
}

type ServiceState string

const (
	Running ServiceState = "Running"
	Stopped ServiceState = "Stopped"
)
