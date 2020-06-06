package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "strconv"
    "time"
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

func (tr *Trade) CashFlow() float64 {
	var provision float64

	if tr.Currency == "PLN" {
		provision = 0.0039
	} else {
		provision = 0.0029
	}

	amount := float64(tr.Count) * tr.Price

	if tr.Buy {
		return - amount * (1 + provision)
	} else {
		return amount * (1 - provision)
	}
}

type RoundTrade struct {
	Buy, Sell		*Trade
	AverageCost		float64
}



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
	fmt.Printf("Deal akcji: %s\tgiełda: %s\twaluta: %s\n", dl.Stock, dl.Exchange, dl.Currency)
	fmt.Printf("Posiadane akcje: %d\n", dl.countStock())
	fmt.Printf("Wynik finansowy: %.2f %s\n\n", dl.sumCashFlow(), dl.Currency)
}

func (dl *Deal) PrintRowIfSold() (float64, string) {
	if dl.countStock() == 0 {
		fmt.Printf("%.2f\t\t%s\t\t%s\n",  dl.sumCashFlow(), dl.Currency, dl.Stock)
		return dl.sumCashFlow(), dl.Currency
	}
	
	return 0, dl.Currency
}

func main() {
	fmt.Println("mtrades")
	fmt.Println("szukam pliku wyeksportowanego z eMaklera do pliku csv")

	file, err := os.Open("./eMakler_historia_zlecen.Csv")
	if err != nil {
		log.Println("Nie udało mi się otworzyć pliku, czy znajduje się w tym samym folderze?")
		log.Fatal(err)
	}

	defer file.Close()

	fmt.Printf("Plik znaleziony, czytam zawartość...\n\n")

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	if !strings.Contains(scanner.Text(), "mBank S.A.") {
		log.Fatal("Nie rozpoznano pierwszej linijki pliku, czy to właściwy plik?")
	}

	// zmienna, która powie nam że następna linijka to nazwisko
	var nameNext bool
	// pętla w której szukamy imienia w pliku
	for scanner.Scan() {
		if nameNext {
			fmt.Printf("Cześć %s\n\n", scanner.Text())
			break
		}

		if !nameNext && strings.Contains(scanner.Text(), "nazwisko") {
			nameNext = true
		}
	}

	var tradesNext bool
	trades := []Trade{}

	// pętla w której szukamy transakcji
	for scanner.Scan() {
		
		if tradesNext {
			trades = append(trades, NewTrade(scanner.Text()))
		}

		if strings.Contains(scanner.Text(), "Stan;Walor;") {
			tradesNext = true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Błąd w trakcie czytania pliku :( ")
        log.Fatal(err)
    }
    
	file.Close()

	fmt.Printf("Znalazłem %d transakcji.\n", len(trades))

	// tutaj grupujemy transakcje w deal'e
	deals := map[string]Deal{}
	var newDeal *Deal

	for _, trade := range(trades) {

		_, present := deals[trade.Stock]

		if !present {
			newDeal = NewDeal(trade)
			newDeal.Search(&trades)
			deals[trade.Stock] = *newDeal
		}
	}

	fmt.Printf("Pogrupowałem je w %d grup (akcjami)", len(deals))

	// zobaczmy jak nasze deal'e!
	for _, deal := range(deals) {
		deal.PrintAll()
	}

	fmt.Printf("Deale zakończone (stan akcji == 0)\n\nWynik\t\twwaluta\t\takcja\n")
	finishedCashFlow := map[string]float64{}
	var sum float64
	var curr string
	for _, deal := range(deals) {
		sum, curr = deal.PrintRowIfSold()
		finishedCashFlow[curr] += sum
	}

	fmt.Printf("\nWynik zamkniętych transakcji:\n")
	for currency, score := range(finishedCashFlow) {
		fmt.Printf("%.2f\t%s\n", score, currency)
	}
	
    fmt.Printf("\n\nkoniec#!")
    fmt.Scanln()
}