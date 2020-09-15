package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"strings"
)

var ctx = context.Background()

/*
word = bear
b
be
bea
bear$
*/
func prefixes(word string) []string {
	n := len(word)
	if n == 0 {
		return []string{}
	}
	parts := make([]string, 0, len(word))
	for i := 1; i <= n; i++ {
		w := word[:i]
		if i == n {
			parts = append(parts, w)
			w += "$"
		}
		parts = append(parts, w)
	}
	return parts
}

func main() {

	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Load dict file: words.txt
	f, err := os.Open("words.txt")
	if err != nil {
		fmt.Println("err:", err)
		panic("failed to open file")
	}
	defer f.Close()

	// TODO: This takes about 24s to load.  Cannot be parallelized as redis is single threaded.
	// Try pipelining.
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		w := scanner.Text()
		w = strings.TrimSuffix(w, "\n")
		plist := prefixes(w)
		zlist := make([]*redis.Z, 0, len(plist))
		for _, p := range plist {
			zlist = append(zlist, &redis.Z{
				Score:  0,
				Member: p,
			})
		}
		intCmd := r.ZAdd(ctx, "prefix_list", zlist...)
		_, err := intCmd.Result()
		if err != nil {
			fmt.Printf("zadd failed- err:%s word:%s", err, w)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("failed to scan: ", err)
	}
}
