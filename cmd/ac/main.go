package main

import (
	"fmt"
	"net/http"

	acm "autocomplete/autocomplete"
)

func ac(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		return
	}

	params := r.URL.Query()

	qvalue, present := params["q"]
	if !present || len(qvalue) == 0 {
		return
	}

	suggestions, err := acm.Autocomplete(qvalue[0])
	if err != nil {
		fmt.Println("failed to return a value:", err)
		return
	}

	fmt.Fprint(w, suggestions)
}

func main() {

	// http://localhost;6666/ac?q=alsdkjf
	http.HandleFunc("/ac", ac)

	println("listening on localhost:6666")
	http.ListenAndServe(":6666", nil)
}
