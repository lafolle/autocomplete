package server

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strings"
)

var ctx = context.Background()

const PREFIX_LIST_KEY = "prefix_list"
const PREFIX_DELIMITER = "$"
const READ_BATCH_SIZE = 50
const MAX_SUGGESTIONS_COUNT = 10

// In strict mode only those results are returned whose prefix matches _prefix_
const STRICT_MODE = true

var r *redis.Client

func init() {
	r = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

// If err != nil then len of suggestions is 0.
func Autocomplete(prefix string) ([]string, error) {
	suggestions := make([]string, 0, MAX_SUGGESTIONS_COUNT)
	// Get rank of prefix
	rank, err := r.ZRank(ctx, PREFIX_LIST_KEY, prefix).Result()
	if err == redis.Nil {
		return suggestions, nil
	}

	if err != nil {
		return suggestions, fmt.Errorf("zadd failed- err:%s word:%s", err, prefix)
	}

	// Get at max MAX_SUGGESTIONS_COUNT suggestions.
	got := 0
LOOP:
	for got < MAX_SUGGESTIONS_COUNT {
		plist, err := r.ZRange(ctx, PREFIX_LIST_KEY, rank+1, rank+READ_BATCH_SIZE).Result()
		if err != nil {
			fmt.Println("zrange error:", err)
			continue
		}
		for _, p := range plist {
			if got == MAX_SUGGESTIONS_COUNT {
				break LOOP
			}
			if strings.HasSuffix(p, PREFIX_DELIMITER) {
				if STRICT_MODE && !strings.HasPrefix(p, prefix) {
					break LOOP
				}
				suggestions = append(suggestions, p[:len(p)-1])
				got += 1
			}
		}
		rank += READ_BATCH_SIZE
	}

	return suggestions, nil
}
