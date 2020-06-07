package main

import (
	"fmt"
	"sort"
	"time"
)


type Deal struct {
	Stock, Exchange, Currency		string
	Trades				[]Trade
}

func NewDeal(firstTrade Trade) (dl *Deal) {
	dl = &Deal{}
	dl.Stock = firstTrade.Stock
	dl.Currency = firstTrade.Currency
	dl.Exchange = firstTrade.Exchange

	return
}

func (dl *Deal) Search(tradeSlice *[]Trade) error {
	if len(*tradeSlice) == 0 {
		return fmt.Errorf("Deal Search failed: provided an empty slice")
	}

	for _, trade := range(*tradeSlice) {
		if trade.Stock == dl.Stock {
			if trade.Exchange != dl.Exchange {
				return fmt.Errorf("Deal Search failed: Exchange not matched!")
			}

			if trade.Currency != dl.Currency {
				return fmt.Errorf("Deal Search failed: Currency not matched!")
			}

			dl.Trades = append(dl.Trades, trade)
		}
	}

	return nil
}

func (dl *Deal) countStock() (count int) {
	for _, trade := range(dl.Trades) {

		if trade.Realised {
			if trade.Buy {
				count += trade.Count
			} else {
				count -= trade.Count
			}
		}
		
	}

	return
}

func (dl *Deal) sumCashFlow() (sum float64) {
	for _, trade := range(dl.Trades) {
		if trade.Realised {
			sum += trade.CashFlow()
		}
	}

	return
}

func (dl *Deal) PrintAll() {
	fmt.Printf("\n#\t\t\t#\nDeal akcji: %s\tgiełda: %s\twaluta: %s\n", dl.Stock, dl.Exchange, dl.Currency)
	fmt.Printf("Posiadane akcje: %d\n", dl.countStock())
	fmt.Printf("Wynik finansowy: %.2f %s\n\n", dl.sumCashFlow(), dl.Currency)

	dl.SortTrades()
	fmt.Printf("Wszystkie transakcje (%d)):\n", len(dl.Trades))
	for _, trade := range(dl.Trades) {
		trade.PrintDetailedRow()
	}

	
	rounds := dl.RoundTrades()
	fmt.Printf("Transakcje sprzedaży z wynikami (%d):\n", len(rounds))
	
	for _, round := range(rounds) {
		round.PrintLine()
	}
	fmt.Printf("\n#\t#\n")
}

func (dl *Deal) PrintRowIfSold() (float64, string) {
	if dl.countStock() == 0 {
		fmt.Printf("%.2f\t\t%s\t\t%s\n",  dl.sumCashFlow(), dl.Currency, dl.Stock)
		return dl.sumCashFlow(), dl.Currency
	}
	
	return 0, dl.Currency
}

func (dl *Deal) SortTrades() {
	sort.Slice(dl.Trades, func(i, j int) bool {
		return dl.Trades[i].Time.Before(dl.Trades[j].Time)
	})
	
}

func (dl *Deal) RoundTrades() (rounds []*RoundTrade) {
	dl.SortTrades()

	var cost float64
	var count int
	var roundTrade *RoundTrade
	var lastBuyTime, firstBuyTime time.Time
	var boughtAlready bool	

	for _, trade := range dl.Trades {
		if trade.Realised {

			if trade.Buy {
				lastBuyTime = trade.Time

				if !boughtAlready {
					firstBuyTime = trade.Time
					boughtAlready = true
				}
				
				if count == 0 {

					
					count = trade.Count
					cost = trade.AverageCost()

				} else {
					
					cost = (trade.AverageCost() * float64(trade.Count) + cost * float64(count)) / float64(trade.Count + count)
					count += trade.Count

				}
			} else {
				if count == 0 {
					return
				}

				roundTrade = &RoundTrade{}
				roundTrade.AverageCost = cost
				roundTrade.Sold = trade
				roundTrade.LastBuyTime = lastBuyTime
				roundTrade.FirstBuyTime = firstBuyTime
				rounds = append(rounds, roundTrade)

				count -= trade.Count
			}
		}
	}

	return
}