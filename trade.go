package main

import (
	"time"
	"strconv"
	"strings"
	"fmt"
)

type Trade struct {
	Realised		bool
	Cancelled		bool
	Exchange		string
	Stock			string
	Buy				bool
	Count			int
	Price			float64
	Currency		string
	Time			time.Time
}

func NewTrade(input string) (tr Trade) {
	inSlice := strings.Split(input, ";")
	
	tr.checkStates(inSlice[0])
	
	tr.Stock = inSlice[1]
	tr.Exchange = inSlice[2]
	
	if (inSlice[3] == "K") {
		tr.Buy = true
	}

	tr.Count, _ = strconv.Atoi(strings.Replace(inSlice[4], " ", "", -1))
	tr.Price, _ = strconv.ParseFloat(strings.Replace(inSlice[5], ",", ".", 1), 64)
	tr.Currency = inSlice[6]

	tr.Time, _ = time.Parse("02.01.2006 15:04:05", inSlice[len(inSlice)-1])

	return
}

func (tr *Trade) checkStates(input string) {
	if strings.Contains(input, "Zrealizowane") {
		tr.Realised = true
		return
	}

	if strings.Contains(input, "Anulowane") {
		tr.Cancelled = true
		return
	}
}

func (tr *Trade) getProvision() float64 {
	if tr.Currency == "PLN" {
		return 0.0039
	}

	return 0.0029
}

func (tr *Trade) CashFlow() float64 {
	
	var provision = tr.getProvision()

	amount := float64(tr.Count) * tr.Price

	if tr.Buy {
		return - amount * (1 + provision)
	} else {
		return amount * (1 - provision)
	}
}

func (tr *Trade) TransactionTypeString() string {
	if tr.Realised {
		if tr.Buy {
			return "[+]"
		}
		return "[-]"
	}

	return "[ ]"
}

func (tr *Trade) PrintDetailedRow() {
	fmt.Printf("%s\t%d\t%s %.4f (%.2f)\n", tr.TransactionTypeString(), tr.Count, tr.Currency, tr.Price, tr.CashFlow())
}

func (tr *Trade) AverageCost() float64 {
	if tr.Buy {
		return tr.Price * (1 + tr.getProvision())
	}

	return tr.Price * (1 - tr.getProvision())
}