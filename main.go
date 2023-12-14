package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
)

const keyServerAddr = "serverAddr"

func getRoot(w http.ResponseWriter, r *http.Request){
    ctx := r.Context()

    hasFirst := r.URL.Query().Has("first")
    First := r.URL.Query().Get("first")
    hasSecond := r.URL.Query().Has("second")
    second := r.URL.Query().Get("second")

    body, err := io.ReadAll(r.Body)
    if err != nil {
        fmt.Printf("%s : err" , err)
    }

    fmt.Printf("%s: got / , first(%t)= %s, second(%t)= %s \n , body : \n%s\n",
    ctx.Value(keyServerAddr),
    hasFirst, First,
    hasSecond, second,
    body)

    io.WriteString(w, "website\n")
}


func getHello(w http.ResponseWriter, r *http.Request){
    ctx := r.Context()

    fmt.Printf("%s got /hello\n", ctx.Value(keyServerAddr))
    io.WriteString(w,"Hello\n")
}

func main(){
    mux := http.NewServeMux()
    mux.HandleFunc("/",getRoot)
    mux.HandleFunc("/hello", getHello)

    ctx, cancelCtx := context.WithCancel(context.Background())
    serverOne := &http.Server{
        Addr: ":8080",
        Handler: mux,
        BaseContext: func(l net.Listener) context.Context{
            ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
            return ctx
        },
    }
    servertow := &http.Server{
        Addr: ":8000",
        Handler: mux,
        BaseContext: func(l net.Listener) context.Context{
            ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
            return ctx
        },
    }

    go func(){

        err := serverOne.ListenAndServe()
        if err != nil {
            fmt.Println("Could not start server with error: ",err)
        }
        cancelCtx()
    }()

    go func(){

        err := servertow.ListenAndServe()
        if err != nil {
            fmt.Println("Could not start server with error: ",err)
        }
        cancelCtx()
    }()

    <-ctx.Done()
}

