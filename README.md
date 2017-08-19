# github-webhooks

[![GoDoc](https://godoc.org/github.com/XescuGC/github-webhooks?status.svg)](https://godoc.org/github.com/XescuGC/github-webhooks)
[![Build Status](https://travis-ci.org/XescuGC/github-webhooks.svg?branch=master)](https://travis-ci.org/XescuGC/github-webhooks)

This is a simple web server that handles Github Webhooks and returns structured objects for each event.

## Usage

```bash
$> go get github.com/XescuGC/github-webhooks
```


```go
 wh := webhooks.New(3000, []string{"project_card"})
 go wh.Start()
 for {
  select {
  case e := <- webhooks.ProjectCards
  // deal with the event e
  }
 }
```
