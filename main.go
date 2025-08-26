package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Tweet represents a simplified tweet structure
// Twitter used to use API v1.1, which returned tweets in this format.
type Tweet struct {
	Text string `json:"text"`
}

func tweetStreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	tracks := strings.Split(r.FormValue("track"), ",")
	if len(tracks) == 0 {
		http.Error(w, "Missing 'track' parameter", http.StatusBadRequest)
		return
	}

	log.Println("Simulating tweets for topics:", tracks)

	// Set headers for streaming response
	w.Header().Set("Content-Type", "application/json")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	for {
		// Delay to simulate real-time tweet arrival
		time.Sleep(time.Duration(1+rand.IntN(3)) * time.Second)

		// Pick random track keyword
		word := tracks[rand.IntN(len(tracks))]

		tweet := Tweet{
			Text: "Someone just mentioned " + word,
		}
		_ = json.NewEncoder(w).Encode(tweet)
		flusher.Flush()
	}
}

func main() {
	// Define the server address flag
	addr := flag.String("addr", ":8080", "HTTP server address")
	flag.Parse()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Register the handler for the mock Twitter streaming endpoint
	http.HandleFunc("/1.1/statuses/filter.json", tweetStreamHandler)

	server := &http.Server{
		Addr: *addr,
	}

	log.Println("Starting server on", *addr)
	// Start the server in a separate goroutine, so we can shut it down gracefully
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal("ListenAndServe:", err)
			}
		}
	}()

	// Wait for termination signal
	<-signalChan
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown Failed:", err)
	}
	log.Println("Server exited properly")
}
