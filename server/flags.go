package main

import "flag"

var inet string
var port int
var verbosity string

func defineConfigFlags() {
	flag.StringVar(&inet, "interface", "0.0.0.0", "interface to listen to")
	flag.StringVar(&verbosity, "verbosity", "info", "verbosity levels: debug, info, warn, error, panic, fatal")
	flag.IntVar(&port, "port", 8000, "Port to listen to")
}
