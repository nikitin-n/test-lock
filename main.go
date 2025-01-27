// Найти решение без утечки памяти? Можно ли сделать решение более компактным?
// Подойдет ли https://pkg.go.dev/oya.to/namedlocker ?
// Или https://pkg.go.dev/github.com/go-auxiliaries/shrinking-map/pkg/safe-map ?

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	handler := Handler{
		lock:       &sync.Mutex{},
		inProgress: make(map[string]struct{}), // тут потенциальная утечка памяти, так как map не уменьшается
	}

	handler.Handler("some key")
}

type Handler struct {
	lock       *sync.Mutex
	inProgress map[string]struct{}
}

func (h *Handler) Handler(key string) {
	h.lock.Lock()

	if _, ok := h.inProgress[key]; ok {
		// Если такой ключ обрабатывается, то выводим лог и выходим
		fmt.Println("Already in progress")

		return
	} else {
		h.inProgress[key] = struct{}{}
		defer func() {
			h.lock.Lock()
			delete(h.inProgress, key)
			h.lock.Unlock()
		}()
	}

	defer h.lock.Unlock()

	// Вызываем метод только если этот key уже не находится в обработке
	doWork(key)
}

func doWork(key string) {
	time.Sleep(time.Second)
}
