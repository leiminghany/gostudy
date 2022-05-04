package main

import  (
	"fmt"
	"time"
)

func main() {
    //a := []byte{1, 2, 3}
    a := []byte("POST / HTTP/1.1\r\nHost: mhlei.us-south.cf.appdomain.cloud\r\nConnection: keep-alive\r\nUser-Agent: Mozilla/5.0 \r\nAccept: */*\r\nAccept-Encoding: gzip, deflate\r\nAccept-Language: zh-CN,zh;q=0.9\r\nContent-Length:40\r\n\r\n")

    fmt.Printf("len:%d\n", len(a))
    b := []byte{2, 3, 4, 5, 6}
    a = append(a, b...)
    time.Sleep(time.Duration(5)*time.Second)
    fmt.Println(a[:20])
}
