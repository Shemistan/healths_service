package internal

import (
	"net/http"
	"time"
)

type Runner struct {
	isRunning bool
	urls      []string
}

func NewRunner(isRunning bool) *Runner {
	return &Runner{isRunning: isRunning}
}

func (r *Runner) Add(arr []string) {
	for _, url := range arr {
		if !isExist(r.urls, url) {
			r.urls = append(r.urls, url)
		}
	}
}

func (r *Runner) Stop() {
	r.isRunning = false
}

func (r *Runner) Start(result chan string, timeout time.Duration) {
	r.isRunning = true
	for {
		if r.isRunning {
			for _, url := range r.urls {
				res, _ := http.Get(url)
				if res.StatusCode != 200 {
					result <- url + " : " + res.Status + "\n"
				}
				result <- url + " : " + res.Status + "\n"
			}
			time.Sleep(time.Second * timeout)
		} else {
			break
		}
	}
}

func (r *Runner) Remove(str string) {
	for i, item := range r.urls {
		if item == str {
			r.urls[i] = r.urls[len(r.urls)-1]
			break
		}
	}
	r.urls = r.urls[:len(r.urls)-1]
}

func (r *Runner) Export() []string {
	return r.urls

}

func isExist(urls []string, str string) bool {
	for _, item := range urls {
		if item == str {
			return true
		}
	}
	return false
}
