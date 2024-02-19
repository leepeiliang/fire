package main

import (
	"bufio"
	parsing2 "fire/pkg/fire"
	"fmt"
	"net"
	"time"
)

// tcp/server/main.go

// TCP server端

// 处理函数
func process(conn net.Conn) {
	var (
		readerbuf []byte
		readlen   int
	)
	defer conn.Close() // 关闭连接
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:]) // 读取数据
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}

		fmt.Printf("收到client端发来的数据：[%d]%x\n", n, buf[:n])
		readerbuf = append(readerbuf, buf[:n]...)
		readlen += n
		if parsing2.ValidationEndSign(buf[:n]) {
			break
		}
		if parsing2.ValidationStartSign(readerbuf) {
			break
		}
	}

	fmt.Printf("收到client端发来的数据all：[%s][%d]%x\n", time.Now().Format("2006-01-02 15:04:05.000 Mon Jan"), readlen, readerbuf[:readlen])
	conn.Write(readerbuf[:readlen]) // 发送数据
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:9139")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn) // 启动一个goroutine处理连接
	}
}
