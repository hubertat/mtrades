package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "sort"
)


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

	var tradesNext bool
	var trade Trade

	trades := []Trade{}

	// pętla w której szukamy transakcji
	for scanner.Scan() {
		
		if tradesNext {
			trade, err = NewTrade(scanner.Text())
			if err != nil {
				log.Print(err)
			} else {
				trades = append(trades, trade)
			}
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
	var dealSlice []*Deal
	var newDeal *Deal

	for _, trade := range(trades) {

		_, present := deals[trade.Stock]

		if !present {
			newDeal = NewDeal(trade)
			newDeal.Search(&trades)
			deals[trade.Stock] = *newDeal
			dealSlice = append(dealSlice, newDeal)
		}
	}

	fmt.Printf("Pogrupowałem je w %d grup (akcjami)", len(deals))

	// zobaczmy jak nasze deal'e!
		
	var allRounds []*RoundTrade

	for _, deal := range(deals) {
		deal.SortTrades()
		deal.PrintAll()
		allRounds = append(allRounds, deal.RoundTrades()...)
	}

	sort.Slice(allRounds, func(i, j int) bool {
		return allRounds[i].PercentResult() <  allRounds[j].PercentResult()
	})

	fmt.Printf("Transakcje sprzedaży z bilansem:\nakcje\t\t\tilość\tkwota\t\tdni (od 1st)\twynik (na dzień)\n")

	for _, round := range(allRounds) {
		round.PrintDetailedLine()
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