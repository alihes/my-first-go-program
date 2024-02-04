package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/freshman-tech/news-demo-starter-files/news"
	"github.com/joho/godotenv"
)

var test int= 0
var tpl = template.Must(template.ParseFiles("index.html"))
var newsapi *news.Client

type Search struct {
	Var1	   int
	Query      string
	NextPage   int
	TotalPages int
	Results    *news.Results
}

func incr(w http.ResponseWriter, r *http.Request) {
	test++
	a := map[string]interface{}{
        "Var1": test,
    }
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}
func decr(w http.ResponseWriter, r *http.Request) {
	test--
	a := map[string]interface{}{
        "Var1": test,
    }
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}


func indexHandler(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("<h1>Hello World! hmm</h1>"))
	// tpl.Execute(w, nil)
	a := map[string]interface{}{
        "Var1": test,
    }
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)

}

func searchHandler(newsapi *news.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		// fmt.Println("Search Query is: ", searchQuery)

		// fmt.Println("page is: ", page)

		results, err := newsapi.FetchEverything(searchQuery, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// fmt.Printf("%+v",results)
		nextPage, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		search := &Search{
			Var1: test,
			Query:      searchQuery,
			NextPage:   nextPage,
			TotalPages: int(math.Ceil(float64(results.TotalResults) / float64(newsapi.PageSize))),
			Results:    results,
		}


		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)

	}
}

func main() {
	fmt.Println("start")
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	apikey := os.Getenv("NEWS_API_KEY")
	if apikey == "" {
		log.Fatal("Env: apikey must be set")
	}

	myClient := &http.Client{Timeout: 10 * time.Second}
	newsapi := news.NewClient(myClient, apikey, 20)

	fs := http.FileServer(http.Dir("assets"))

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/search", searchHandler(newsapi))
	mux.HandleFunc("/incr", incr)
	mux.HandleFunc("/decr", decr)
	http.ListenAndServe(":"+port, mux)
}
