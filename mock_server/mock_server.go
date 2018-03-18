package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	count := parseIntQueryParam(r, "num", 10)
	min := parseIntQueryParam(r, "min", 1)
	max := parseIntQueryParam(r, "max", 100)
	if max < min {
		fmt.Fprintln(w, "Error: max < min")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	log.Printf("count = %v, min = %v, max = %v\n", count, min, max)
	w.WriteHeader(http.StatusOK)
	for i := 0; i < count; i++ {
		fmt.Fprintf(w, "%d\n", rand.Intn(max-min+1)+min)
	}
}

func parseIntQueryParam(r *http.Request, name string, defaultValue int) int {
	val, err := strconv.Atoi(r.FormValue(name))
	if err != nil {
		val = defaultValue
	}
	return val
}

func main() {
	http.HandleFunc("/integers/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
