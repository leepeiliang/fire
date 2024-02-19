package firenet

import (
	"fire/pkg/fire"
	"fire/pkg/fireface"
)

// Request 请求
type Request struct {
	conn fireface.IConnection //已经和客户端建立好的 链接
	msg  fire.FireMessage     //客户端请求的数据
}

// GetConnection 获取请求连接信息
func (r *Request) GetConnection() fireface.IConnection {
	return r.conn
}

// GetData 获取请求消息的数据
func (r *Request) GetMsgData() []byte {
	return r.msg.GetMsgData()
}

// GetMsgID 获取请求的消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// BackFireMessage 获取请求消息的数据
func (r *Request) BackFireMessage() fire.FireMessage {
	return r.msg
}
