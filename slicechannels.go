package main

import (
	"sync"
)

type valuation struct {
	address string
	town    string
	value   int
}

func main() {
}

func filterValue(valuations []valuation, price int) []valuation {
	//filter out houses under a certain value
	var filteredValues []valuation

	for k := range valuations {
		if valuations[k].value > price {
			filteredValues = append(filteredValues, valuations[k])
		}
	}
	return filteredValues
}

func filterInChunks(s []valuation, streets []string, price int, x int) []valuation {
	//takes a slice and splits into smaller chunks, runs the filter on each chunk and then joins them all together
	c := splitSplices(s)
	c1 := filterSplices(c, price)
	c2 := filterSplices(c, price)

	var valuations []valuation
	for v := range merge(c1, c2) {
		valuations = append(valuations, v...)
	}
	return valuations
}

func splitSplices(v []valuation) <-chan []valuation {
	//splits the splice into smaller slices and fans out
	out := make(chan []valuation)
	threads := 4
	chunk := (len(v) + 1/threads)
	go func() {
		for i := 0; i < len(v); i += chunk {
			s := v[i : i+chunk]
			out <- s
		}
		close(out)
	}()
	return out
}
func filterSplices(in <-chan []valuation, price int) <-chan []valuation {
	//filters the slices it recieves in the chan's, filters and fans back in
	out := make(chan []valuation)
	go func() {
		for v := range in {
			out <- filterValue(v, price)
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan []valuation) <-chan []valuation {
	//receives the slices from the multiple channels and merges them into one
	var wg sync.WaitGroup
	out := make(chan []valuation)

	output := func(c <-chan []valuation) {
		for v := range c {
			out <- v
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
