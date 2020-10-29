package main

import (
	"fmt"
	"time"
)

type RoundTrade struct {
	Sold			Trade
	AverageCost		float64
	LastBuyTime		time.Time
	FirstBuyTime	time.Time
}

func (rt *RoundTrade) Result() float64 {
	return rt.Sold.CashFlow() - rt.AverageCost * float64(rt.Sold.Count)
}

func (rt *RoundTrade) FromLastBuy() time.Duration {
	return rt.Sold.Time.Sub(rt.LastBuyTime)
}

func (rt *RoundTrade) FromFirstBuy() time.Duration {
	return rt.Sold.Time.Sub(rt.FirstBuyTime)
}

func (rt *RoundTrade) LastBuyDays() float64 {
	return float64(rt.FromLastBuy().Hours()) / 24
}

func (rt *RoundTrade) FirstBuyDays() float64  {
	return float64(rt.FromFirstBuy().Hours()) / 24
}

func (rt *RoundTrade) PrintLine() {
	fmt.Printf("%d\t%s [%.2f]\t%.2f\t%.2f (%.2f)\t%.2f %%\t%.2f %%/d\n", 
		rt.Sold.Count, 
		rt.Sold.Currency, 
		rt.Sold.Price,
		rt.Result(), 
		rt.LastBuyDays(), 
		rt.FirstBuyDays(), 
		rt.PercentResult(),
		rt.DayReturnRate())
}

func (rt *RoundTrade) PrintDetailedLine() {

	fmt.Printf("%20s\t%d\t%s %4.2f\t%3.2f (%3.2f)\t%3.2f %%\t%2.2f %%/d\n", 
		rt.Sold.Stock,
		rt.Sold.Count, 
		rt.Sold.Currency, 
		rt.Result(), 
		rt.LastBuyDays(), 
		rt.FirstBuyDays(), 
		rt.PercentResult(),
		rt.DayReturnRate())
}

func (rt *RoundTrade) PercentResult() float64 {
	return (rt.Sold.AverageCost() / rt.AverageCost - 1) * 100
}

func (rt *RoundTrade) DayReturnRate() float64 {
	return rt.PercentResult() / rt.FirstBuyDays()
}