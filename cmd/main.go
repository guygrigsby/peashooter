package main

import (
	"context"
	"os"
	"strconv"

	"github.com/guygrigsby/peashooter/client"
	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_CONCURRENCY = 1000
)

func main() {
	var (
		concurrency int
		log         = logrus.New()
	)
	c, err := strconv.ParseInt(os.Getenv("PEASHOOTER_CONCURRENCY"), 10, 16)
	if err != nil {

		concurrency = DEFAULT_CONCURRENCY
	} else {
		concurrency = int(c)
	}
	client := client.New(context.Background(), "https://prolifewhistleblower.com/anonymous-form/", log.WithField("service", "client"))
	client.Forever(concurrency)
}
