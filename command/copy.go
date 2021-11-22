package command

import (
	"context"
	"log"
	"runtime"

	"github.com/go-redis/redis/v8"
)

type Copier struct {
	Ctx          context.Context
	SourceClient *redis.Client
	TargetClient *redis.Client
	Pattern      string
}

func (c *Copier) Copy() error {
	keys, err := c.SourceClient.Keys(c.Ctx, c.Pattern).Result()
	if err != nil {
		return err
	}
	chunkSize := runtime.NumCPU()
	chunks := chunks(keys, chunkSize)

	copyOp := func(keys []string, ch chan error) {
		for _, key := range keys {
			log.Printf("copying %s", key)
			data, err := c.SourceClient.Get(c.Ctx, key).Result()
			if err != nil {
				log.Fatal("can't get source data")
				ch <- err
			}
			ttl, err := c.SourceClient.TTL(c.Ctx, key).Result()
			log.Printf("ttl: %d", ttl)
			if err != nil {
				log.Fatal("can't get ttl for key " + key)
				ch <- err
			}

			if ttl < 0 {
				ttl = 0
			}
			log.Printf("start writing to new source: %s", key)
			err = c.TargetClient.Set(c.Ctx, key, data, ttl).Err()

			if err != nil {
				log.Fatalf("Set data to new source failed: %v", err)
				ch <- err
			}
		}
	}

	errCh := make(chan error)
	for _, chunkedKeys := range chunks {
		go copyOp(chunkedKeys, errCh)
	}
	err = <-errCh
	log.Println("copy data success for ", c.Pattern)
	return err
}

func chunks(k []string, size int) [][]string {
	var prev, i int

	if len(k) == 0 {
		return nil
	}
	div := make([][]string, (len(k)+size-1)/size)
	till := len(k) - size
	for prev < till {
		next := prev + size
		div[i] = k[prev:next]
		prev = next
		i++
	}
	div[i] = k[prev:]
	return div
}
