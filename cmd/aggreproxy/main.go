package main

import "Aggregator/cmd/aggreproxy/proxy"

func main() {
	p := proxy.New()
	p.Run()
}
