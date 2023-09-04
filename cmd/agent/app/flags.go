package app

import (
	"flag"
)

var flagReportInterval int
var flagPoolInterval int
var flagServerAddr string

func parseFlags() {
	flag.IntVar(&flagPoolInterval, "p", 2, "frequency pool metrics in sec")
	flag.IntVar(&flagReportInterval, "r", 10, "frequency send metrics to server in sec")
	flag.StringVar(&flagServerAddr, "a", ":8080", "address and port to run server")
	flag.Parse()
}
