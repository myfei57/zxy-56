package main

import (
	"flag"
	"os"
)

type config struct {
	addr string
	data string
}

func loadConfig() config {
	cfg := config{addr: ":18056", data: "./data"}
	flag.StringVar(&cfg.addr, "addr", cfg.addr, "listen address")
	flag.StringVar(&cfg.data, "data", cfg.data, "data directory")
	flag.Parse()
	if env := os.Getenv("LABVENT_ADDR"); env != "" {
		cfg.addr = env
	}
	if env := os.Getenv("LABVENT_DATA"); env != "" {
		cfg.data = env
	}
	return cfg
}
