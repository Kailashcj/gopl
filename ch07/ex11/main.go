// ch07/ex11 は、クライアントによる商品の作成、読み出し、更新、削除が可能な Web サーバです。
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	db := database{
		"shoes": 50,
		"socks": 5,
	}

	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)

	http.HandleFunc("/create", db.create)
	http.HandleFunc("/read", db.read)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)

	fmt.Println("Listening at http://localhost:8000")
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if item == "" {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintln(w, "item has no name")
	} else if _, ok := db[item]; ok {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "item already exists: %q\n", item)
	} else {
		priceStr := req.URL.Query().Get("price")
		price, err := strconv.ParseFloat(priceStr, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(w, "invalid price: %q\n", priceStr)
			return
		}
		db[item] = dollars(price)
		fmt.Fprintf(w, "created: %s: %s\n", item, db[item])
	}
}

func (db database) read(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if _, ok := db[item]; ok {
		priceStr := req.URL.Query().Get("price")
		price, err := strconv.ParseFloat(priceStr, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(w, "invalid price: %q\n", priceStr)
			return
		}
		db[item] = dollars(price)
		fmt.Fprintf(w, "updated: %s: %s\n", item, db[item])
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if _, ok := db[item]; ok {
		delete(db, item)
		fmt.Fprintf(w, "deleted: %s\n", item)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}
