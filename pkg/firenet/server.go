package firenet

import (
	"fmt"
	"net"

	"fire/config"
	"fire/pkg/fireface"
)

var zinxLogo = `                                        
██████████     ██     ██    ▄█▀  ██████████
██             ▀▀     ██   ▄█▀   ██
██           ████     ██  ▄█▀    ██
████████       ██     ██▀▀▀      ████████
██             ██     ██         ██
██          ▄▄▄██▄▄▄  ██         ██
▀▀          ▀▀▀▀▀▀▀▀  ██         ██████████
                                        `
var topLine = `┌──────────────────────────────────────────────────────┐`
var borderLine = `│`
var bottomLine = `└──────────────────────────────────────────────────────┘`

// Server 接口实现，定义一个Server服务类
type Server struct {
	Server config.Server
	//当前Server的消息管理模块，用来绑定MsgID和对应的处理方法
	msgHandler fireface.IMsgHandle
	//当前Server的链接管理器
	ConnMgr fireface.IConnManager
	//该Server的连接创建时Hook函数
	OnConnStart func(conn fireface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn fireface.IConnection)

	packet fireface.Packet
}

// NewServer 创建一个服务器句柄
func NewServer(opts config.Server) fireface.IServer {
	printLogo(opts)

	s := &Server{
		Server:     opts,
		msgHandler: NewMsgHandle(opts.WorkerPoolSize, opts.MaxWorkerTaskLen),
		ConnMgr:    NewConnManager(),
		packet:     NewDataPack(opts.MaxPacketSize),
	}

	return s
}

// ============== 实现 fireface.IServer 里的全部接口方法 ========
func (s *Server) Config() config.Server {
	return s.Server
}

// Start 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Server.Name, s.Server.Host, s.Server.TCPPort)

	//开启一个go去做服务端Linster业务
	go func() {
		//0 启动worker工作池机制
		s.msgHandler.StartWorkerPool()

		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.Server.IPVersion, fmt.Sprintf("%s:%d", s.Server.Host, s.Server.TCPPort))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		//2 监听服务器地址
		listener, err := net.ListenTCP(s.Server.IPVersion, addr)
		if err != nil {
			panic(err)
		}

		//已经监听成功
		fmt.Println("start fire-up server  ", s.Server.Name, " succ, now listenning...")

		//TODO server.go 应该有一个自动生成ID的方法
		var cID uint16
		cID = 0

		//3 启动server网络连接业务
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			fmt.Println("Get conn remote addr = ", conn.RemoteAddr().String())

			//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= s.Server.MaxConn {
				conn.Close()
				continue
			}

			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConnection(s, conn, cID, s.msgHandler)
			cID++

			//3.4 启动当前链接的处理业务
			go dealConn.Start()
		}
	}()
}

// Stop 停止服务
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Server.Name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

// Serve 运行服务
func (s *Server) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	select {}
}

// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgID uint32, router fireface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

// GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() fireface.IConnManager {
	return s.ConnMgr
}

// SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(fireface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(fireface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn fireface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn fireface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}

func (s *Server) Packet() fireface.Packet {
	return s.packet
}

func printLogo(opts config.Server) {
	fmt.Println(zinxLogo)
	fmt.Println(topLine)
	fmt.Println(fmt.Sprintf("%s [Github] https://21vianet.com                 %s", borderLine, borderLine))
	fmt.Println(fmt.Sprintf("%s [tutorial] http://meta42.indc.vnet.com %s", borderLine, borderLine))
	fmt.Println(bottomLine)
	fmt.Printf("[EdgeFire] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		opts.IPVersion,
		opts.MaxConn,
		opts.MaxPacketSize)
}
