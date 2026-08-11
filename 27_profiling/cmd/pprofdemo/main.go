package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"

	"go_examples/27_profiling/fib"
)

func main() {
	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		n := uint(35)
		if raw := r.URL.Query().Get("n"); raw != "" {
			parsed, err := strconv.ParseUint(raw, 10, 32)
			if err != nil || parsed > 45 {
				http.Error(w, "n must be an integer from 0 to 45", http.StatusBadRequest)
				return
			}
			n = uint(parsed)
		}
		fmt.Fprintln(w, fib.Slow(n))
	})

	log.Println("app:   http://127.0.0.1:6060/work?n=35")
	log.Println("pprof: http://127.0.0.1:6060/debug/pprof/")
	log.Fatal(http.ListenAndServe("127.0.0.1:6060", nil))
}
