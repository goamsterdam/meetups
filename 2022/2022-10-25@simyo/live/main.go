package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/mjl-/bstore"

	_ "embed"
)

type Msg struct {
	ID uint64 `json:"-"`

	Timestamp time.Time `bstore:"index,default now" json:"timestamp"`

	Text string `json:"message"`

	Email    string `bstore:"-" json:"email"`
	EmailMD5 string `json:"email_md5"`
}

var (
	RequestTotal atomic.Int64

	//go:embed index.html
	indexHtml []byte
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write(indexHtml); err != nil {
		panic(err)
	}
}

func sent(db *bstore.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg Msg
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			panic(err)
		}

		m := md5.Sum([]byte(msg.Email))
		msg.EmailMD5 = hex.EncodeToString(m[:])

		if err := db.Insert(&msg); err != nil {
			panic(err)
		}
	}
}

func list(db *bstore.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msgs, err := bstore.QueryDB[Msg](db).SortDesc("Timestamp").Limit(10).List()
		if err != nil {
			panic(err)
		}

		if len(msgs) == 0 {
			msgs = []Msg{
				{Text: "No messages yet"},
			}
		}

		if err := json.NewEncoder(w).Encode(msgs); err != nil {
			panic(err)
		}
	}
}

func metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]any{
		"request_total": RequestTotal.Load(),
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	db, err := bstore.Open("msgs.db", nil, Msg{})
	if err != nil {
		panic(err)
	}

	statsWrapper := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			RequestTotal.Add(1)

			h(w, r)
		}
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	http.HandleFunc("/", statsWrapper(index))
	http.HandleFunc("/metrics", statsWrapper(metrics))
	http.HandleFunc("/sent", statsWrapper(sent(db)))
	http.HandleFunc("/list", statsWrapper(list(db)))
	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			time.Sleep(time.Second * 5)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				panic(err)
			}
		}()
	})

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}
}
