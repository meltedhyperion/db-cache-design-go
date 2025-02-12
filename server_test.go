package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestHandleGetUserConncurent(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(http.HandlerFunc(s.handleGetUser))
	nreq := 1000

	wg := &sync.WaitGroup{}
	startTime := time.Now()
	for i := 0; i < nreq; i++ {
		wg.Add(1)
		go func(i int) {
			id := i % 100
			url := fmt.Sprintf("%s/?id=%d", ts.URL, id)
			resp, err := http.Get(url)
			if err != nil {
				t.Error(err)
			}
			user := &User{}
			if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
				t.Error(err)
			}
			// fmt.Printf("request %d: %+v\n", i+1, user)
			wg.Done()
		}(i)

		time.Sleep(1 * time.Millisecond) // i've kept it low since my system is powerful. One can adjust as per their system :)
	}
	wg.Wait()
	endTime := time.Now()

	fmt.Println("Time taken for 1000 requests (async):", endTime.Sub(startTime))
	fmt.Println("DB hits:", s.dbHit)
}

func TestHandleGetUserSync(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(http.HandlerFunc(s.handleGetUser))
	nreq := 1000

	startTime := time.Now()
	for i := 0; i < nreq; i++ {
		id := i % 100
		url := fmt.Sprintf("%s/?id=%d", ts.URL, id)
		resp, err := http.Get(url)
		if err != nil {
			t.Error(err)
		}
		user := &User{}
		if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
			t.Error(err)
		}
		// fmt.Printf("request %d: %+v\n", i+1, user)

		time.Sleep(1 * time.Millisecond) // i've kept it low since my system is powerful. One can adjust as per their system :)
	}

	endTime := time.Now()

	fmt.Println("Time taken for 1000 requests (sync):", endTime.Sub(startTime))
	fmt.Println("DB hits:", s.dbHit)
}
