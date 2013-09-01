package models

import (
    "fmt"
    "strconv"
    "time"
    "github.com/coopernurse/gorp"
)

type Link struct {
    Id   int
    Slug string
    Url  string
    Created string
    Updated string
}

type LinkView struct {
    Id  int64
    Url string
}

func (this *Link) Save() {
    dbmap.Insert(this)
}

func (this *Link) PreInsert(s gorp.SqlExecutor) error {
    now := time.Now()
    this.Created = fmt.Sprintf("%d-%d-%d %02d:%02d:%02d\n", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
    this.Updated = fmt.Sprintf("%d-%d-%d %02d:%02d:%02d\n", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
    if this.Slug == "" {
        this.Slug = generateSlug()
    }
    return nil
}

func Find(obj interface{}, query string, args ...interface{}) error {
    if _, err := dbmap.Select(obj, query, args...); err != nil {
        return err
    }
    return nil
}

func generateSlug() string {
    now := time.Now()
    count, _  := dbmap.SelectStr("select count(*) from links where DATE(Created) = DATE(NOW())")
    dayOfYear := strconv.Itoa(now.YearDay())
    year      := strconv.Itoa(now.Year())

    dateCount, _ := strconv.Atoi(year + dayOfYear + count)
    sgx := ""
    for dateCount > 0 {
        mod := dateCount % 60
        sgx = NewBase60Vocab[mod] + sgx
        dateCount = (dateCount - mod)/60
    }
    return sgx
}
