package fire

import (
	"fire/pkg/data"
	"time"
)

// Control 控制单元
type Control struct {
	//业务流水号 (2 字节)应答方按请求包的业务流水 号返回。低字节传输在前。业务流水号是一个 2 字节的正整数，
	SerialNumber
	//协议版本号 (2 字节)协议版本号包含主版本号(第 5 字节)和用户版本号(第 6 字节)。主版本号为 固定值 2，用户版本号由用户自行定义。
	ProtocolVersion
	//数据包的第 7~12 字节，为数据包发出的时间，具体定义见 10.2.2。
	TimeLabels
	//数据包的第 13~18 字节，为数据包的源地址(上位机、消防控制显示装置或火 灾自动报警设备地址)。低字节传输在前。(参考建议:若发送方为平台，默认 地址为 0x00 0x00 0x00 0x00 0x00 0x00，若发送方非平台，可将源地址设置为 ID 的前六位。ID 组成为:四信 ID(1 字节)+设备类型(1 字节)+设备编号(4 字节)+预留(2 字节))
	SourceAddress
	//数据包的第 19~24 字节，为数据包的目的地址(上位机、消防控制显示装置或 火灾自动报警设备地址)。低字节传输在前。
	TargetAddress
	//数据包的第 25、26 字节，为应用数据单元的长度，长度不应大于 1024;低字节 传输在前。
	AppDataLen
	//数据包的第 27 字节，控制单元的命令字节，具体定义见表 2。
	ControlCommand
}

// 获取消息数据段长度
func (f *Control) GetDataLen() uint16 {
	return f.AppDataLen.AppDataLen
}

// 设置消息数据段长度
func (f *Control) SetDataLen(datalen uint16) {
	f.AppDataLen.AppDataLen = datalen
	return
}

type SerialNumber struct {
	SerialNumber uint16
}

type ProtocolVersion struct {
	ProtocolVersion [2]byte
}
type TimeLabels struct {
	TimeLabels [6]byte
}

func (t *TimeLabels) FireReadTime(buffer *data.Buffer) {
	tmp := buffer.ReadN(6)
	copy(t.TimeLabels[:], tmp)
	return
}
func (t TimeLabels) FireToTime() string {
	if len(t.TimeLabels) != 6 {
		return time.Now().Format("2006-01-02 15:04:05.000 Mon Jan")
	}
	var (
		year, day, hour, min, sec int
		month                     time.Month
	)
	// TODO 默认年是20XX  其中字节数据范围是0-99 所以年范围就是2000～2099
	year = int(t.TimeLabels[5]) + 2000
	month = time.Month(int(t.TimeLabels[4]))
	day = int(t.TimeLabels[3])
	hour = int(t.TimeLabels[2])
	min = int(t.TimeLabels[1])
	sec = int(t.TimeLabels[0])
	return time.Date(year, month, day, hour, min, sec, 12345600, time.Local).Format("2006-01-02 15:04:05.000 Mon Jan")
}
func (t TimeLabels) FireToTimeUnixNano() int64 {
	if len(t.TimeLabels) != 6 {
		return time.Now().UnixNano() / 1e9
	}
	var (
		year, day, hour, min, sec int
		month                     time.Month
	)
	// TODO 默认年是20XX  其中字节数据范围是0-99 所以年范围就是2000～2099
	year = int(t.TimeLabels[5]) + 2000
	month = time.Month(int(t.TimeLabels[4]))
	day = int(t.TimeLabels[3])
	hour = int(t.TimeLabels[2])
	min = int(t.TimeLabels[1])
	sec = int(t.TimeLabels[0])

	back := time.Date(year, month, day, hour, min, sec, 0, time.Local)
	return back.UnixNano() / 1e9
}

func FireToSetTime(now time.Time) *TimeLabels {
	var (
		fireTime    = &TimeLabels{}
		fireTimebuf [6]byte
	)
	// TODO 默认年是20XX  其中字节数据范围是0-99 所以年范围就是2000～2099
	fireTimebuf[5] = byte(now.Year() % 100)
	fireTimebuf[4] = byte(now.Month())
	fireTimebuf[3] = byte(now.Day())
	fireTimebuf[2] = byte(now.Hour())
	fireTimebuf[1] = byte(now.Minute())
	fireTimebuf[0] = byte(now.Second())
	fireTime.TimeLabels = fireTimebuf
	return fireTime
}

type SourceAddress struct {
	Address
}

type Address struct {
	Address [6]byte
}
type TargetAddress struct {
	Address
}

type AppDataLen struct {
	AppDataLen uint16
}

func (t *AppDataLen) SetAppDataLen(datalen uint16) {
	t.AppDataLen = datalen
}
func (t AppDataLen) GetAppDataLen() uint16 {
	return t.AppDataLen
}
