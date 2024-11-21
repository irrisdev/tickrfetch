package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/irrisdev/go-coinlore/coinlore"
	"github.com/irrisdev/tickrfetch/logger"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

func NewCoinService(client *coinlore.Client, limiter *rate.Limiter) *CoinService {
	return &CoinService{
		client:  client,
		limiter: limiter,
	}
}

func (s *CoinService) GetCoin(symbol string) (*coinlore.Coin, error) {
	symbol = strings.ToUpper(symbol)
	symbol = strings.ReplaceAll(symbol, " ", "")

	return nil, fmt.Errorf("$%s not found", symbol)
}

// func (s *CoinService) FetchCoin(symbol string) error {

// 	coin, err := s.client.GetCoin(symbol)
// 	if err != nil {
// 		logger.Error("failed to fetch coin: ", err)
// 		return fmt.Errorf("failed to fetch coin: %v", err)
// 	}
// 	logger.Info("fetched coin: ", coin.Name)
// 	return nil
// }

func (s *CoinService) Start(c chan string) {

}

func (s *CoinService) fetchHotCoins() error {
	coins, err := s.client.GetCoins(0, 100)
	if err != nil {
		logger.Error("Failed to fetch coins: ", err)
		return fmt.Errorf("failed to fetch coins: %v", err)
	}

	for _, coin := range *coins {
		s.cache.Store(coin.Symbol, &coin)
	}

	logger.Info("fetched top coins: ", len(*coins))

	return nil
}

func (s *CoinService) fetchSymbols() error {

	global, err := s.client.GetGlobal()
	if err != nil {
		logger.Error("Failed to global stats: ", err)
		return fmt.Errorf("failed to fetch global stats: %v", err)
	}

	const limit = 100
	var results []coinlore.Coin

	for start := 0; start < global.CoinsCount; start += limit {
		coins, err := s.client.GetCoins(start, limit)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"start": start,
				"error": err,
			}).Warn("Retrying failed fetch coins request")

			retries := 3
			for attempt := 1; attempt <= retries; retries++ {
				time.Sleep(time.Second * time.Duration(attempt))
				coins, err = s.client.GetCoins(start, limit)
				if err == nil {
					break
				}
				logger.WithField("attempt", attempt).Warn("Fetch coins retry failed")
			}

			if err != nil {
				logger.WithField("start", start).Error("Failed to fetch coins after retries")
				continue
			}

		}
		results = append(results, (*coins)...)
		logger.WithField("start", start).Info("Fetched batch successfully")

		logger.Info("Sleeping for 1s before next request")
		time.Sleep(1 * time.Second)
	}

	for _, coin := range results {
		s.symbols.Store(coin.Symbol, coin.ID)
	}

	if len(results) != global.CoinsCount {

		logger.WithFields(logrus.Fields{
			"total":   global.CoinsCount,
			"fetched": len(results),
		}).Warn("Partially fetched all cold coins")

	} else {
		logger.WithField("total_coins", len(results)).Info("Fetched cold coins successfully")
	}

	return nil
}
