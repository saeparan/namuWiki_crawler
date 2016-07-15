package main

import (
	"log"
	"math/rand"
	"github.com/PuerkitoBio/goquery"
    "github.com/mauidude/go-readability"
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
    "time"
    "sync"
)

type Article struct {
    Id            bson.ObjectId `bson:"_id,omitempty"`
    BoardID       string        `bson:"BoardID"`
    Title         string        `bson:"Title"`
    Content       string        `bson:"Content"`
    NickName      string        `bson:"NickName"`
    RegDate       time.Time     `bson:"RegDate"`
    RegDateString string
    Hit           int    `bson:"Hit"`
    Vote          int    `bson:"Vote"`
    IP            string `bson:"IP"`
}
var wg sync.WaitGroup

func main() {
    for {
        wg.Add(1)
        go execute()
        wg.Wait()        
    }
}

func execute() {
    defer wg.Done()

    url := fetchUrl()
    article, err := crawl(url)
    if err != nil {
        log.Println(err)
    }
    log.Println("Insert")
    articleInsert(article)
    time.Sleep(10000 * time.Millisecond)
}

func articleInsert(article *Article) {
    session, err := mgo.Dial("127.0.0.1")
    if err != nil {
        panic(err)
    }

    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
    c := session.DB("test").C("article")

    article.RegDate = time.Now()
    article.NickName = "Bot"
    err = c.Insert(article)
    if err != nil {
        log.Println(err)
        return
    }
    log.Println("Success")
}

func fetchUrl() string {
    urlItems := []string{
        "https://namu.wiki/random",
        //"https://ko.wikipedia.org/wiki/%ED%8A%B9%EC%88%98:%EC%9E%84%EC%9D%98%EB%AC%B8%EC%84%9C",
    }
    index := rand.Intn(1)
    if( len(urlItems) >= index ) {
        return urlItems[index]
    }
    return ""
}

func crawl(url string) (*Article, error) {
    article := &Article{}
    doc, err := goquery.NewDocument(url)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }

    bodyHtml, _ := doc.Find("body").Html()
    content, err := readability.NewDocument(bodyHtml)
    if err != nil {
      log.Println("Error:", err)
    }

    article.Title = doc.Find("h1.title").Text()
    article.Content = content.Content()
    return article, nil
}
