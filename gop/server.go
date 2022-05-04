package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"log"
	"net"
	"net/url"
	"strings"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
        outfile, err := os.OpenFile("goproxy.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666);
        if err != nil {
                fmt.Println(err);
                os.Exit(55);
        }
	log.SetOutput(outfile);
	l, err := net.Listen("tcp", ":45670")
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
	//log.Printf("recieve:[%s]\n", string(b[:n]));

	var i, offset int
	for i=0; i<n; i++ {
		if (string(b[i:i+4])  ==  "\r\n\r\n") {
			break;
		}
	}
	offset = i+4
	offset = 0
	//log.Printf("offset:%d", offset)
	if offset >= n {
		return
	}

	var method, host, address string
	firstLine := string(b[offset:offset+bytes.IndexByte(b[offset:], '\n')]);
	//log.Printf("firstLine:%s", firstLine)
	fmt.Sscanf(firstLine, "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return
	}

	if hostPortURL.Opaque == "443" { //https访问
		address = hostPortURL.Scheme + ":443"
	} else { //http访问
		if strings.Index(hostPortURL.Host, ":") == -1 { //host不带端口， 默认80
			address = hostPortURL.Host + ":80"
		} else {
			address = hostPortURL.Host
		}
	}

	//获得了请求的host和port，就开始拨号吧
	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("dial to [%s] ok", address)

	if method == "CONNECT" {
		_, err := fmt.Fprint(client, "HTTP/1.1 200 OK\r\n\r\n")
		//_, err := client.Write([]byte("HTTP/1.0 200 OK\r\nContent-Length: 9\r\n\r\nGOODBYE\r\n"))
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		_, err := server.Write(b[offset:n])
		if err != nil {
			log.Println(err)
			return
		}
	}

	//进行转发
	go io.Copy(server, client)
	io.Copy(client, server)
	// go proxyRequest(client, server)
	// proxyRequest(server, client)
}

func proxyRequest(r net.Conn, w net.Conn) {
    defer r.Close()
    defer w.Close()

    var buffer = make([]byte, 4096000)
    for {
        n, err := r.Read(buffer)
        if err != nil {
            log.Printf("Unable to read from input, error: %s\n", err.Error())
            break
        }

	log.Printf("forward:[%s]", string(buffer[:n]) )

        n, err = w.Write(buffer[:n])
        if err != nil {
            log.Printf("Unable to write to output, error: %s\n", err.Error())
            break
        }
    }
}
