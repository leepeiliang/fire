package main

import (
	"fire/pkg/fire"
	"flag"
	"fmt"
	"io"
	"k8s.io/klog/v2"
	"net"
	"time"

	"fire/pkg/firenet"
)

var host = "127.0.0.1:30119"

//var host = "192.168.198.15:30119" //B28
/*
	模拟客户端
*/
func main() {
	var msgType int

	flag.IntVar(&msgType, "msgType", 1, "年龄")

	//解析命令行参数
	flag.Parse()

	klog.Infof("msgType：%d", msgType)

	switch msgType {

	case 1:
		go uploadAirSamplingFireFirstAlarm()
	default:
		fmt.Println("请输入支持的数据类型")
	}

	//等待子进程执行完毕，也可以用wait
	time.Sleep(10 * time.Second)

}

// 空气采样-故障 首次火警
func uploadAirSamplingFireFirstAlarm() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	//上传建筑消防设施系统运行状态1
	var sendMsg = []byte{
		0x40, 0x40,
		0x02, 0x00,
		0x02, 0x03,
		0x2A, 0x29, 0x09, 0x10, 0x01, 0x18,
		0x10, 0x25, 0x51, 0x80, 0x12, 0xFF,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x45, 0x00,
		0x02,
		0x8D,       //类型标识
		0x01,       //信息对象数目
		0x88,       //系统类型
		0x40, 0x00, //数据长度
		0x1B, 0x36, 0x32, 0x34, 0x2F, 0x30, 0x31, 0x2F,
		0x31, 0x32, 0x20, 0x20, 0x31, 0x35, 0x3A, 0x32,
		0x34, 0x0A, 0x20, 0x20, 0x30, 0x30, 0x34, 0x2D,
		0x30, 0x30, 0x32, 0x0A, 0x1B, 0x39, 0xC9, 0xE8,
		0xB1, 0xB8, 0xB9, 0xCA, 0xD5, 0xCF, 0x0A, 0xBC,
		0xD0, 0xB2, 0xE3, 0xC5, 0xE4, 0xB5, 0xE7, 0xCA,
		0xD2, 0xBB, 0xF0, 0xBE, 0xAF, 0x0A, 0xC6, 0xF8,
		0xCC, 0xE5, 0xC3, 0xF0, 0xBB, 0xF0, 0x0A, 0x0A,
		0xAA, //校验位
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(sendMsg)
	if err != nil {
		klog.Errorf("client write err: ", err)
		return
	}
	sendMsgLen, err := conn.Write(msg)
	if err != nil {
		klog.Errorf("client write err: ", err)
		return
	}
	klog.Infof("Client send -------------------success[%d]--------------------", sendMsgLen)
	klog.Infof("Client recv --------------------start-------------------")
	//msg 是有data数据的，需要再次读取data数据
	fm := fire.FireMessage{}
	//1 先读出流中的head部分
	startData := make([]byte, fire.StartSignLen)
	_, err = io.ReadFull(conn, startData) //ReadFull 会把msg填充满为止
	if err != nil {
		klog.Errorf("read head error")
		return
	}
	klog.Infof("==>Client Recv msg:startData=%x", startData)
	copy(fm.Start[:], startData[:fire.StartSignLen])
	//2 先读出流中的控制部分
	controData := make([]byte, fire.FireControlLen)
	_, err = io.ReadFull(conn, controData) //ReadFull 会把msg填充满为止
	if err != nil {
		klog.Errorf("read head error")
		return
	}
	//将Control字节流 拆包到msg中
	msgControl, err := dp.UnpackControl(controData)
	if err != nil {
		klog.Errorf("server unpack err:", err)
		return
	}
	fm.Control = *msgControl
	klog.Infof("==>Client Recv msg:Control Data=%x", controData)
	klog.Infof("==>Client Recv msg:struct:Control=%v", fm.Control)
	fm.Data = make([]byte, msgControl.GetDataLen())
	if fm.Control.GetDataLen() > 0 {
		//根据dataLen从io中读取字节流
		backlen, err := io.ReadFull(conn, fm.Data)
		if err != nil {
			klog.Errorf("server unpack data err:", err)
			return
		}

		klog.Infof("==> Client Recv DataMsg:len=%d, data=%x backlen=%d", msgControl.GetDataLen(), fm.Data, backlen)
	}
	var crc []byte
	crc = make([]byte, fire.CRCLen)
	_, err = io.ReadFull(conn, crc)
	if err != nil {
		klog.Errorf("read msg data error ", err)
		return
	}
	fm.CRC = crc[0]
	klog.Infof("==> Client Recv msg:CRC Data=%x", crc)

	crcData := make([]byte, 25+fm.Control.GetDataLen())
	crcData = append(crcData, controData...)
	crcData = append(crcData, fm.Data...)
	// 校验数据CRC
	klog.Infof("fireEnder.CRC[%x:%x]", fm.CRC, fire.CRC(crcData))
	if fm.CRC != fire.CRC(crcData) {
		return
	}
	var end []byte
	end = make([]byte, fire.StartSignLen)
	if _, err := io.ReadFull(conn, end); err != nil {
		klog.Errorf("Client read msg data error ", err)
		return
	}
	copy(fm.End[:], end[:fire.StartSignLen])

	klog.Infof("==>Client Recv EndData:Data=%x", fm.End)
	klog.Infof("Client recv --------------------END-success------------------")
	time.Sleep(time.Second)
	//	}
}
