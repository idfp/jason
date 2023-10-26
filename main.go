package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/BishopFox/jsluice"
)


func main() {
	stat, _ := os.Stdin.Stat()
  var links []string
  var files [][]byte
  var urlResults []*jsluice.URL
  var secretResults []*jsluice.Secret

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		var buf []byte
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			buf = append(buf, scanner.Bytes()...)
      buf = append(buf, 10)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
    buffer := string(buf[:])
    temp := strings.Split(buffer, "\n")

    links = temp[:len(temp) - 1]
	}else{
    fmt.Println("Please pass the links through stdin piping\nExample: cat urls.txt | jason [OPTIONS]")
    return
  }
  for _, link := range links{
    if !(strings.Contains(link, "http://") || strings.Contains(link, "https://")){
      continue
    }
    resp, err := http.Get(link)
    if err != nil{
      log.Fatal(err)
    }
    body, err := io.ReadAll(resp.Body)
    files = append(files, body)
  }

  for _, file := range files{
    analyzer := jsluice.NewAnalyzer(file)
    for _, url := range analyzer.GetURLs() {
      urlResults = append(urlResults, url)
    }
    for _, secret := range analyzer.GetSecrets() {
      secretResults = append(secretResults, secret)
    }
  }

  res, err := json.MarshalIndent(urlResults, "", " ")
  if err != nil{
    log.Fatal(err)
  }
  if err := os.WriteFile("results-url.txt", []byte(res), 0666); err != nil{
    log.Fatal(err)
  }

  res2, err := json.MarshalIndent(secretResults, "", " ")
  if err != nil{
    log.Fatal(err)
  }
  if err := os.WriteFile("results-secret.txt", []byte(res2), 0666); err != nil{
    log.Fatal(err)
  }


}
