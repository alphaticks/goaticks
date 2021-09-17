package alphaticks

import (
	"fmt"
	"gitlab.com/tachikoma.ai/tickobjects/market"
	"io"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client, err := NewClient("-6946936700949124468", "4KGMaKcRm28ekc_yLkffjw==", DB_LIVE)
	if err != nil {
		t.Fatal(err)
	}
	qs := NewQuery()
	from := time.Now().Add(-time.Hour * 24 * 30)
	to := from.Add(10 * time.Second)
	qs.WithFrom(from)
	qs.WithTo(to)
	qs.WithSelector("BestBid(orderbook)")
	qs.WithTags(map[string]string{
		"exchange": "^binance$",
		"symbol":   "^SOLUSDT$",
	})
	q, err := client.Query(qs)
	if err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	for q.Next() {
		tick, obj, _ := q.Read()
		sec := int64(tick / 1000)
		nsec := (int64(tick) - sec*1000) * 1000000
		bb := obj.(*market.BestBid)
		ts := time.Unix(sec, nsec)
		fmt.Println(ts, bb.Price, bb.Quantity)
	}
	fmt.Println(time.Now().Sub(start))
	if q.Err() != io.EOF {
		t.Fatal(q.Err())
	}
}

func TestBestBid(t *testing.T) {
	client, err := NewClient("-6946936700949124468", "4KGMaKcRm28ekc_yLkffjw==", DB_LIVE)
	if err != nil {
		t.Fatal(err)
	}
	qs := NewQuery()
	from := time.Now()
	to := from.Add(10 * time.Second)
	qs.WithFrom(from)
	qs.WithTo(to)
	qs.WithSelector("BestBid(orderbook)")
	//qs.WithTimeout(100 * time.Millisecond)
	qs.WithTags(map[string]string{
		"base":  "^RUNE$",
		"quote": "^USDT$",
		"type":  "CRSPOT",
	})
	qs.WithStreaming(false)
	q, err := client.Query(qs)
	if err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	for q.Next() {
		tick, obj, _ := q.Read()
		tags := q.Tags()
		sec := int64(tick / 1000)
		nsec := (int64(tick) - sec*1000) * 1000000
		bb := obj.(*market.BestBid)
		ts := time.Unix(sec, nsec)
		fmt.Println(ts, tags["exchange"], bb.Price, bb.Quantity)
	}
	fmt.Println(time.Now().Sub(start))
	if q.Err() != io.EOF {
		t.Fatal(q.Err())
	}
}
