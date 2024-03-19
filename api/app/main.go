package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/net/http2"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var pingInterval = time.Second * 10

// HandleLiveLog returns an http.HandlerFunc that processes an http.Request
// to server sent event.
func HandleLiveLog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		h := w.Header()
		h.Set("Content-Type", "text/event-stream")
		h.Set("Cache-Control", "no-cache")
		h.Set("Connection", "keep-alive")
		h.Set("X-Accel-Buffering", "no")
		h.Set("Access-Control-Allow-Origin", "wails://wails.localhost:34115")
		h.Set("Access-Control-Allow-Headers", "Content-Type")
		h.Set("Access-Control-Allow-Credentials", "true")

		f, ok := w.(http.Flusher)
		if !ok {
			return
		}

		io.WriteString(w, ": ping\n\n")
		f.Flush()
		n := 0
	L:

		for {
			logrus.Info("looping: ", n)
			select {
			case <-ctx.Done():
				logrus.Info("events: stream cancelled")
				break L
			case <-time.After(pingInterval):
				io.WriteString(w, ": heart beat\n\n")
				f.Flush()
			default:
				logrus.Info("events: ping")
				event, err := formatServerSentEvent("ping", time.Now().Format(time.RFC1123))
				if err != nil {
					fmt.Println(err)
					break
				}
				fmt.Fprintf(w, event)
				f.Flush()
				time.Sleep(1 * time.Second)
				if n%2 == 0 {
					logrus.Info("events: pong")
					event, err := formatServerSentEvent("pong", string(time.Now().UnixNano()))
					if err != nil {
						fmt.Println(err)
						break
					}
					fmt.Fprintf(w, event)
					f.Flush()
				}
				n = n + 1
			}
		}

		io.WriteString(w, "event: error\ndata: eof\n\n")
		f.Flush()

		logrus.Info("events: stream closed")

	}
}

func main() {

	r := chi.NewRouter()
	r.Get("/livelog", HandleLiveLog())

	srv := &http.Server{
		Addr:        ":4344",
		ReadTimeout: 75 * time.Second,
		Handler:     r,
	}
	http2.ConfigureServer(srv, &http2.Server{})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(c)

		select {
		case <-ctx.Done():
		case <-c:
			cancel()
		}
	}()
	var g errgroup.Group

	g.Go(func() error {
		return srv.ListenAndServeTLS("example.com+5.pem", "example.com+5-key.pem")
	})

	g.Go(func() error {
		<-ctx.Done()
		return srv.Shutdown(ctx)
	})

	if err := g.Wait(); err != nil {
		logrus.WithError(err).Info("shutdown server")
	}
}

// formatServerSentEvent takes name of an event and any kind of data and transforms
// into a server sent event payload structure.
// Data is sent as a json object, { "data": <your_data> }.
//
// Example:
//
//	Input:
//		event="price-update"
//		data=date
//	Output:
//		event: price-update\n
//		data: date"\n\n
func formatServerSentEvent(event string, data any) (string, error) {
	//m := map[string]any{
	//	"data": data,
	//}

	//buff := bytes.NewBuffer([]byte{})

	//encoder := json.NewEncoder(buff)

	//err := encoder.Encode(m)
	//if err != nil {
	//	return "", err
	//}

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("event: %s\n", event))
	sb.WriteString(fmt.Sprintf("data: %v\n\n", data.(string)))

	return sb.String(), nil
}
