### Promise In Go

## Install
`go get github.com/fabsolute/promise-in-go`

## Usage 

The basic usage is to just do 

```go
// create file main.go
package main

import (
       promise "github.com/fabsolute/promise-in-go"
       "fmt"
       "time"
       "strconv"
       )
func main(){
  response := promise.New(func(resolve, reject func(interface{})) {
    time.Sleep(2 * time.Second)
    resolve(2)
  }).Then(func(value interface{}) interface{} {
    return value.(int) + 4
  }).Then(func(value interface{}) interface{} {
    return "This message has " + strconv.Itoa(value.(int)) + " words."
  }).Await()

  fmt.Println(response)
}
```
Execute
```go run main.go```

You will see

```This message has 6 words.```
