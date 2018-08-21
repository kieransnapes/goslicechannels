# goslicechannels
processes a slice  in multiple goroutines by splitting the slice into smaller parts

My examples takes a houses value (valuation) and filters based on a value.
The idea is that we can split the filtering process into smaller chunks, filter
the chunks and then merge the results. Ideally this will handle large datasets
faster.
