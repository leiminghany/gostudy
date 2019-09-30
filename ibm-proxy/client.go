package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panic(err)
	}

	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		go handleClientRequest(client)
	}
}

func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	var b [409600]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("read %d bytes from client\n", n);

	//var address = "192.168.0.207:8088"
	var address = "mhlei.us-south.cf.appdomain.cloud:80"
	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("dial to [%s] ok", address)

	newHeader := []byte(fmt.Sprintf("POST / HTTP/1.1\r\nHost: mhlei.us-south.cf.appdomain.cloud\r\nConnection: keep-alive\r\nUser-Agent: Mozilla/5.0 \r\nAccept: */*\r\nAccept-Encoding: gzip, deflate\r\nAccept-Language: zh-CN,zh;q=0.9\r\nContent-Length:%d\r\n\r\n", n) )
	sendBuf := append(newHeader, b[:n]...)

	nw, err := server.Write(sendBuf[:len(newHeader)+n])
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("write %d bytes to server, [%s]\n", nw, string(sendBuf))

	/*nr, err := server.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("read %d bytes from server, [%s]\n", nr, string(b[:n]));

	nc, err := client.Write(b[:nr])
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("return %d bytes to client\n", nc);
	*/

	//进行转发
	// go io.Copy(server, client)
	// io.Copy(client, server)
	go proxyRequest(client, server)
	proxyRequest(server, client)
}

// Forward all requests from r to w
func proxyRequest(r net.Conn, w net.Conn) {
    defer r.Close()
    defer w.Close()

    var buffer = make([]byte, 4096000)
    for {
	n, err := r.Read(buffer[:])
        if err != nil {
            log.Printf("Unable to read from input, error: %s\n", err.Error())
            break
        }

	log.Printf("forward:[%s]", string(buffer[:n]) )

        n, err = w.Write(buffer[:n])

	if(string(buffer[9:12]) != "200") {
		continue;
	}

        if err != nil {
            log.Printf("Unable to write to output, error: %s\n", err.Error())
            break
        }
    }
}
