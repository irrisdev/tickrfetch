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

func (s *CoinService) GetCoin(symbol string) *coinlore.Coin {
	symbol = strings.ToUpper(symbol)
	symbol = strings.ReplaceAll(symbol, " ", "")

	if coin, ok := s.cache.Load(symbol); ok {
		return coin.(*coinlore.Coin)
	}

	if symbolID, ok := s.symbols.Load(symbol); ok {
		coin, err := s.client.GetCoin(symbolID.(int))
		if err != nil {
			logger.WithFields(logrus.Fields{
				"symbol":    symbol,
				"symbol_id": symbolID,
				"error":     err,
			}).Error("Failed to fetch coin by symbol")
			return nil
		}
		return coin
	}
	return nil
}

func (s *CoinService) Start(c chan ServiceState, coldcoins bool) {
	logger.Info("Starting coin service")
	logger.Info("Prefetching coins and symbols")

	err := s.fetchHotCoins()
	if err != nil {
		logger.Error("Failed to prefetch hot coins: ", err)
		c <- Stopped
		return
	}
	c <- Running

	go s.fetchCycle(c)

	s.symbols.Store("prefetched", true)

	if coldcoins {
		err = s.fetchSymbols()
		if err != nil {
			logger.Error("Failed to prefetch symbols: ", err)
			s.symbols.Store("prefetched", false)
		}
	}
}

func (s *CoinService) fetchCycle(c chan ServiceState) {
	hotSchedule := time.NewTicker(1 * time.Second)
	coldSchedule := time.NewTicker(24 * time.Hour)
	defer hotSchedule.Stop()
	defer coldSchedule.Stop()

	for {
		select {
		case <-hotSchedule.C:
			go s.fetchHotCoins()
		case <-coldSchedule.C:
			go s.fetchSymbols()
		case _, ok := <-c:
			if !ok {
				logger.Info("Channel closed, stopping service")
				return
			}
		}
	}

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

	s.cache.Store("lastfetch", time.Now())

	if len(*coins) != 100 {
		logger.WithFields(logrus.Fields{
			"total":   100,
			"fetched": len(*coins),
		}).Warn("Partially fetched hot coins")
	} else {
		logger.Info("Fetched hot coins successfully")
	}

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

	s.symbols.Store("lastfetch", time.Now())

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

func (s *CoinService) Prefetched() bool {
	if val, ok := s.symbols.Load("prefetched"); ok {
		return val.(bool)
	}
	return false
}
