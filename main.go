package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "io/ioutil"
    "strings"
    "sort"
    "time"
)


func main() {
	fmt.Printf("##\n#mtrades\n##\n")
	fmt.Println("szukam pliku wyeksportowanego z eMaklera do pliku csv")

    files, err := ioutil.ReadDir(".")
    if err != nil {
    	log.Println("Nieudany odczyt folderu:")
        log.Fatal(err)
    }

    var csvs []string
    for _, file := range files {
        if strings.HasSuffix(strings.ToLower(file.Name()), ".csv") {
        	csvs = append(csvs, file.Name())
        }
    }

    var whichCsv int

    if len(csvs) == 0 {
    	log.Fatal("Nie odnaleziono żadnego pliku .csv w tym folderze :( ")
    }

    if len(csvs) > 1 {
    	fmt.Println("Odnaleziono więcej niż jeden plik .csv")
    	fmt.Println("Wybierz który otworzyć, wpisz odpowiednią liczbę i wciśnij [enter]:")
    	for count, name := range(csvs) {
    		fmt.Printf("%d - %s\n", count, name)
    	}
    	n, err := fmt.Scan(&whichCsv)
    	if n != 1 || err != nil {
    		log.Fatal("Błąd odczytywania wybranego numeru.")
    	}
    	if whichCsv >= len(csvs) {
    		log.Fatal("Wpisano zły numer, spróbuj ponownie uruchomić program.")
    	}
    }

    fmt.Printf("Otwieram plik: %s\n\n", csvs[whichCsv])

	file, err := os.Open("./" + csvs[whichCsv])
	if err != nil {
		log.Println("Nie udało mi się otworzyć pliku:")
		log.Fatal(err)
	}

	defer file.Close()

	fmt.Printf("Plik otwarty, czytam zawartość...\n\n")

	scanner := bufio.NewScanner(file)

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

	fmt.Printf("\nPosiadane akcje, wybierz nr do symulacji:\n0 - zakończ\n")
	iter := 0
	unfinishedDeals := map[int]Deal{}
	for _, deal := range(deals) {
		if amount := deal.countStock(); amount > 0 {
			iter++
			fmt.Printf("%d - %s [%d]\n", iter, deal.Stock, amount)
			unfinishedDeals[iter] = deal
		}
	}

	var whichDeal int
	n, err := fmt.Scan(&whichDeal)
	if err != nil || n != 1 {
		log.Fatal("Błąd odczytywania wybranego numeru.")
	}

	if whichDeal > 0 && whichDeal <= iter {
		deal := unfinishedDeals[whichDeal]
		amount := deal.countStock()
		fmt.Printf("Wybrano: %s, liczba akcji: %d\nWpisz docelowy kurs sprzedaży: \n", deal.Stock, amount)
		var targetPrice float64
		n, err = fmt.Scan(&targetPrice)
		if n != 1 || err != nil {
			log.Fatal("Błąd odczytywania ceny/kursu!")
		}
		simulatedSell := Trade{}
		simulatedSell.Realised = true
		simulatedSell.Stock = deal.Stock
		simulatedSell.Currency = deal.Currency
		simulatedSell.Exchange = deal.Exchange
		simulatedSell.Count = amount
		simulatedSell.Price = targetPrice
		simulatedSell.Time = time.Now()
		deal.Trades = append(deal.Trades, simulatedSell)
		fmt.Printf("\nDodano symulowaną transkakcję sprzedaży, wynik:\n")
		deal.PrintAll()

	}
	
    fmt.Printf("\n\nkoniec#!")
    fmt.Scanln()
}