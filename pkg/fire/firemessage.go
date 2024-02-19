package fire

import "k8s.io/klog/v2"

const (
	DefaultFireMsgID = iota
)

// FireMessage 数据包结构
type FireMessage struct { //启动符‘@@’(2 字节) 数据包的第 1、2 字节，为固定值 64，64。十进制  十六进制是0x40,0x40
	Start [2]byte
	Control
	//数据负载 具体数据具体解析
	Data []byte
	//CRC校验
	CRC byte
	//结束符‘##’ (2 字节)为固定值 35，35。
	End [2]byte
}

// NewFireMessage 创建一个NewFireMessage实例
func NewFireMessage(data []byte) *FireMessage {
	var fireMessage = &FireMessage{}
	fireMessage.UnmarshalRes(data)
	return fireMessage
}

// GetMsgID 获取请求的消息的ID
func (fm *FireMessage) GetMsgID() uint32 {
	return DefaultFireMsgID
}

// 获取消息数据段长度
func (fm *FireMessage) GetFireMessageDataLen() uint16 {
	return StartSignLen + FireControlLen + fm.Control.AppDataLen.AppDataLen + CRCLen + EndSignLen
}

// GetData 获取返回给前端的消息内容
func (fm *FireMessage) GetMsgData() []byte {
	return fm.MarshalResp()
}

// GetControl 获取消息控制单元
func (fm *FireMessage) GetControl() Control {
	return fm.Control
}

// SetData 设置消息数据内容
func (fm *FireMessage) SetData(data []byte) {
	fm.SetDataLen(uint16(len(data)))
	fm.Data = make([]byte, 0)
	//	klog.V(0).Infof("FireMessage.SetData:len[%s][%v:%x]", fm.GetDataLen(), data)
	fm.Data = append(fm.Data, data...)
	klog.V(3).Infof("FireMessage.SetData:len[%d][%x]", fm.GetDataLen(), fm.Data)
	return
}
func (fm *FireMessage) GetData() []byte {
	return fm.Data
}

// SetControl 设置消息控制单元
func (fm *FireMessage) SetControl(header Control) {
	fm.Control = header
	return
}

// SetDataLen 设置消息数据段长度
func (fm *FireMessage) SetDataLen(datalen uint16) {
	fm.Control.SetAppDataLen(datalen)
}

// GetDataLen 获取消息数据段长度
func (fm *FireMessage) GetDataLen() uint16 {
	return fm.Control.GetAppDataLen()
}
