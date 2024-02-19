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

// var host = "192.168.212.10:30119" //外高桥
// var host = "192.168.195.14:30119" //下沙
// var host = "192.168.248.10:30119" //腾仁
// var host = "192.168.198.15:30119" //m6a
// var host = "127.0.0.1:30119"
// var host = "192.168.226.11:30119"
var host = "192.168.194.17:30119"

/*
模拟客户端
*/
func main() {
	var msgType int

	flag.IntVar(&msgType, "msgType", 36, "年龄")

	//解析命令行参数
	flag.Parse()

	klog.Infof("msgType：%d", msgType)

	switch msgType {
	case 1:

		go loginFun()
	case 2:
		go uploadTimeFun()
	case 3:
		go uploadAnalogFunc1()
	case 4:
		go uploadAnalogFunc2()
	case 5:
		go uploadAnalogFunc3()
	case 6:
		go uploadEnergyFunc()
	case 7:
		go uploadFireSystemctlFunc1()
	case 8:
		go uploadFireSystemctlFunc2()
	case 9:
		go uploadFireBuildFacilitiesPartFunc()
	case 10:
		go uploadFireBuildFacilitiesPartFunc3()
	case 11:
		go uploadFireBuildFacilitiesPartFunc4()
	case 12:
		go uploadFireBuildFacilitiesPartFunc5()
	case 13:
		go uploadFireBuildFacilitiesPartFunc6()
	case 14:
		go uploadFireBuildFacilitiesPartFunc6()
	case 15:
		go uploadFireBuildFacilitiesPartFunc7()
	case 16:
		go uploadSystemctlVersionFunc()
	case 17:
		go uploadSystemctlPRTFunc()
	case 18:
		go uploadSystemctlCRTFunc()
	case 19:
		go uploadSystemctlPRTFunc2()
	case 20:
		go uploadSystemctlPRTFunc3()
	case 21:
		go uploadSystemctlPRTFunc4()
	case 22:
		go uploadSystemctlPRTFunc5()
	case 23:
		go uploadSystemctlPRTFuncResetM6V()
	case 24:
		go uploadSystemctlPRTFuncResetsoft()
	case 25:
		go uploadFireSystemctlFunc0()
	case 28:
		// 9--消防广播P 25--多线广播盘0区域0楼层0房间 22--011回路002地址联动请求
		go uploadSystemctlPRTFuncJiyun()
	case 29:
		go uploadSystemctlPRTFuncJiyun2()
	case 30:
		go uploadSystemctlPRTTongji()
	case 31:
		go uploadSystemctlCRTFuncTongjiReSet()
	case 32:
		go uploadSystemctlCRTFuncBoxing()
	case 33:
		go uploadSystemctlPRTFuncResetM3()
	case 34:
		go uploadSystemctlPRTFuncM3()
	case 35:
		go uploadSystemctlPRTFuncM32()
	case 36:
		//上传消控主机解析卡PRT信息  经开
		go uploadSystemctlPRTFuncXIANJK()
	case 37:
		//上传消控主机解析卡PRT信息  经开
		go uploadSystemctlPRTFuncXIANJKReset()
	case 38:
		//上传消控主机解析卡PRT信息  经开
		go uploadSystemctlPRTFuncXIANJK2()

	case 42:
		go uploadSystemctlCRTFunctengren()
	case 43:
		go uploadSystemctlCRTFunctengren2()
	case 44:
		go uploadSystemctlCRTFunctengrenReSet()
	case 46:
		//上传消控主机解析卡PRT信息  外高桥
		go uploadSystemctlPRTFuncWGQReset()
	case 49:
		go uploadSystemctlPRTFuncWaiGaoqiao3()
	case 50:
		go uploadSystemctlPRTFuncTongji()
	case 51:
		go uploadSystemctlPRTFunM3qimei()
	case 52:
		go uploadSystemctlPRTFunBoXingKongcai()

	case 53:
		go uploadSystemctlPRTFunjIyun333()
	case 54:
		go uploadSystemctlPRTFuntengrenhuida()
	case 55:
		go uploadSystemctlPRTFuntengrenfuwei()

	case 60:
		//上传消控主机解析卡PRT信息  B28
		go uploadSystemctlPRTFuncB28Reset()

	case 61:
		//上传消控主机解析卡PRT信息  B28空调
		go uploadSystemctlPRTFuncB28air()
	case 70:
		go uploadSystemctlPRTFuncXIASHAReset()
	case 71:
		go uploadSystemctlPRTFuncXIASHA()
	case 72:
		go uploadSystemctlPRTFuncXIASHA2()
	case 73:
		go uploadSystemctlPRTFuncXIASHA3()
	case 74:
		go uploadSystemctlPRTFuncXIASHA4()
	default:
		fmt.Println("请输入支持的数据类型")
	}

	//等待子进程执行完毕，也可以用wait
	time.Sleep(10 * time.Second)

}

// 上传建筑消防设施系统心跳
func uploadFireSystemctlFunc0() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	//上传建筑消防设施系统运行状态1
	var sendMsg = []byte{
		0x40, 0x40,
		0x07, 0x25, 0x51, 0x80, 0x12, 0xff, 0x04, 0x88,
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

// 上传建筑消防设施系统运行状态1
func uploadFireSystemctlFunc1() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	//上传建筑消防设施系统运行状态1
	var sendMsg = []byte{
		0x40, 0x40,
		0x01, 0x00,
		0x02, 0x03,
		0x35, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x15, 0x27, 0x08, 0x20, 0x20, 0xFE,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x0c, 0x00,
		0x02,
		0x01,
		0x01,
		0x01,
		0x87,
		0x57, 0x05,
		0x35, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x7d,
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

// 上传建筑消防设施系统运行状态2
func uploadFireSystemctlFunc2() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	//上传建筑消防设施系统运行状态2
	var sendMsg = []byte{
		0x40, 0x40,
		0x01, 0x00,
		0x02, 0x03,
		0x35, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x15, 0x27, 0x08, 0x20, 0x20, 0xFE,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x0e, 0x00,
		0x02,
		0x82,
		0x01,
		0x81,
		0x87,
		0xd7, 0x15,
		0x01,
		0x08,
		0x35, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x19,
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

// 上传建筑消防设施部件模拟量值-消防管道
func uploadAnalogFunc1() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	//03消防管道压力系统 数量3 128消防水压
	var uploadAnalog = []byte{
		0x40, 0x40,
		0x02, 0x00,
		0x02, 0x03,
		0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x15, 0x27, 0x08, 0x20, 0x20, 0xFE,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x32, 0x00,
		0x02,
		0x03,
		0x03,
		0x83, 0x03, 0x80, 0x15, 0x27, 0x08, 0x20, 0x04, 0x80, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x83, 0x03, 0x80, 0x15, 0x27, 0x08, 0x20, 0x05, 0x80, 0x01, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x83, 0x03, 0x80, 0x15, 0x27, 0x08, 0x20, 0x9a, 0x60, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x12,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传建筑消防设施部件模拟量值2-电气火灾监控主机系统
func uploadAnalogFunc2() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	//129电气火灾系统 数量14 17剩余电流式电气火灾监控探测器
	var uploadAnalog = []byte{
		0x40, 0x40,
		0x02, 0x00,
		0x02, 0x03,
		0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x15, 0x27, 0x08, 0x20, 0x20, 0xFE,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0xE2, 0x00,
		0x02,
		0x03,
		0x0E,
		0x81, 0x03, 0x11, 0x15, 0x27, 0x08, 0x20, 0x82, 0x00, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x81, 0x03, 0x11, 0x15, 0x27, 0x08, 0x20, 0x86, 0x0C, 0xFE, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x81, 0x03, 0x11, 0x15, 0x27, 0x08, 0x20, 0x8D, 0x0C, 0xFE, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x81, 0x03, 0x11, 0x15, 0x27, 0x08, 0x20, 0x8E, 0x0C, 0xFE, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x81, 0x03, 0x11, 0x15, 0x27, 0x08, 0x20, 0x8F, 0x0C, 0xFE, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x81, 0x03, 0x11, 0x15, 0x27, 0x08, 0x20, 0x94, 0x00, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x81, 0x03, 0x11, 0x15, 0x27, 0x08, 0x20, 0x95, 0x00, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x81, 0x03, 0x11, 0x15, 0x27, 0x08, 0x20, 0x96, 0x00, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x81, 0x03, 0x11, 0x15, 0x27, 0x08, 0x20, 0x97, 0x00, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x81, 0x03, 0x11, 0x15, 0x27, 0x08, 0x20, 0x98, 0x00, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x81, 0x03, 0x11, 0x15, 0x27, 0x08, 0x20, 0x99, 0x00, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x81, 0x03, 0x11, 0x15, 0x27, 0x08, 0x20, 0xB1, 0x00, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x81, 0x03, 0x11, 0x15, 0x27, 0x08, 0x20, 0xB2, 0x00, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x81, 0x03, 0x11, 0x15, 0x27, 0x08, 0x20, 0xB3, 0x38, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x79,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传建筑消防设施部件模拟量值3-LoRa 烟感网关系统、
func uploadAnalogFunc3() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	//03LoRa 烟感网关系统 数量3 40感烟火灾探测器
	var uploadAnalog = []byte{
		0x40, 0x40,
		0x02, 0x00,
		0x02, 0x03,
		0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x15, 0x27, 0x08, 0x20, 0x20, 0xFE,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x82, 0x00,
		0x02,
		0x03,
		0x08,
		0x82, 0x04, 0x28, 0x15, 0x27, 0x08, 0x20, 0x03, 0x10, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x82, 0x04, 0x28, 0x15, 0x27, 0x08, 0x20, 0x06, 0x10, 0x01, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x82, 0x04, 0x28, 0x15, 0x27, 0x08, 0x20, 0x80, 0x20, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x82, 0x04, 0x28, 0x15, 0x27, 0x08, 0x20, 0xfa, 0x01, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x82, 0x04, 0x28, 0x15, 0x27, 0x08, 0x20, 0xfb, 0x20, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x82, 0x04, 0x28, 0x15, 0x27, 0x08, 0x20, 0xfc, 0x30, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x82, 0x04, 0x28, 0x15, 0x27, 0x08, 0x20, 0xfd, 0x80, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x82, 0x04, 0x28, 0x15, 0x27, 0x08, 0x20, 0xfe, 0x00, 0x00, 0x38, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0xcb,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传建筑消防设施部件运行状态
func uploadFireBuildFacilitiesPartFunc() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	//上传建筑消防设施部件运行状态
	var sendMsg = []byte{
		0x40, 0x40,
		0x01, 0x00,
		0x02, 0x03,
		0x35, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x15, 0x27, 0x08, 0x20, 0x20, 0xFE,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x30, 0x00,
		0x02,
		0x02,
		0x01,
		0x81, 0x03, 0x11, 0x15, 0x27, 0x08, 0x20,
		0x02, 0x00, 0xB0, 0xA1, 0xA9, 0x96, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x35, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x49,
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

// 上传建筑消防设施部件运行状态2-火灾报警系统-烟感
func uploadFireBuildFacilitiesPartFunc3() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	//上传建筑消防设施部件运行状态-火灾报警系统-烟感
	var sendMsg = []byte{
		0x40, 0x40,
		0x01, 0x00,
		0x02, 0x03,
		0x35, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x15, 0x27, 0x08, 0x20, 0x20, 0xFE,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x32, 0x00,
		0x02,
		0x02,                   //标识
		0x01,                   //数量
		0x87,                   //系统
		0x28,                   //系统地址
		0x28,                   //部件
		0x03, 0x03, 0x03, 0x03, //部件地址
		0x57, 0x05,
		0x02, 0x00, 0xB0, 0xA1, 0xA9, 0x96, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00,
		0x00,
		0x35, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x91,
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

// 上传建筑消防设施部件运行状态2-火灾报警系统-防火门
func uploadFireBuildFacilitiesPartFunc4() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	//上传建筑消防设施部件运行状态-火灾报警系统-防火门
	var sendMsg = []byte{
		0x40, 0x40,
		0x01, 0x00,
		0x02, 0x03,
		0x35, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x15, 0x27, 0x08, 0x20, 0x20, 0xFE,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x32, 0x00,
		0x02,
		0x02,                   //标识
		0x01,                   //数量
		0x87,                   //系统
		0x04,                   //系统地址
		0x66,                   //部件
		0x04, 0x04, 0x04, 0x04, //部件地址
		0x57, 0x05,
		0x02, 0x00, 0xB0, 0xA1, 0xA9, 0x96, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00,
		0x00,
		0x35, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0xaf,
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

// 上传建筑消防设施部件运行状态2-火灾报警系统-消防水压
func uploadFireBuildFacilitiesPartFunc5() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	//上传建筑消防设施部件运行状态-火灾报警系统-消防水压
	var sendMsg = []byte{
		0x40, 0x40,
		0x01, 0x00,
		0x02, 0x03,
		0x35, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x15, 0x27, 0x08, 0x20, 0x20, 0xFE,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x32, 0x00,
		0x02,
		0x02,                   //标识
		0x01,                   //数量
		0x87,                   //系统
		0x05,                   //系统地址
		0x80,                   //部件
		0x05, 0x05, 0x05, 0x05, //部件地址
		0x57, 0x05,
		0x02, 0x00, 0xB0, 0xA1, 0xA9, 0x96, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00,
		0x00,
		0x35, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0xce,
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

// 上传建筑消防设施部件运行状态2-火灾报警系统-消防液位
func uploadFireBuildFacilitiesPartFunc6() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	//上传建筑消防设施部件运行状态-火灾报警系统-消防液位
	var sendMsg = []byte{
		0x40, 0x40,
		0x01, 0x00,
		0x02, 0x03,
		0x35, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x15, 0x27, 0x08, 0x20, 0x20, 0xFE,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x32, 0x00,
		0x02,
		0x02,                   //标识
		0x01,                   //数量
		0x87,                   //系统
		0x07,                   //系统地址
		0x82,                   //部件
		0x07, 0x07, 0x07, 0x07, //部件地址
		0x57, 0x05,
		0x02, 0x00, 0xB0, 0xA1, 0xA9, 0x96, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00,
		0x00,
		0x35, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0xda,
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

// 上传建筑消防设施部件运行状态2-火灾报警系统-消防栓智能闷盖
func uploadFireBuildFacilitiesPartFunc7() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	//上传建筑消防设施部件运行状态-火灾报警系统-消防液位
	var sendMsg = []byte{
		0x40, 0x40,
		0x01, 0x00,
		0x02, 0x03,
		0x35, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0x15, 0x27, 0x08, 0x20, 0x20, 0xFE,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x32, 0x00,
		0x02,
		0x02,                   //标识
		0x01,                   //数量
		0x87,                   //系统
		0x06,                   //系统地址
		0x81,                   //部件
		0x06, 0x06, 0x06, 0x06, //部件地址
		0x57, 0x05,
		0x02, 0x00, 0xB0, 0xA1, 0xA9, 0x96, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00,
		0x00,
		0x35, 0x12, 0x0F, 0x0F, 0x06, 0x15,
		0xd4,
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

// 上报电能
func uploadEnergyFunc() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	//电能
	var uploadEnergy = []byte{
		0x40, 0x40,
		0xEB, 0x01,
		0x02, 0x03,
		0x15, 0x06, 0x0A, 0x07, 0x04, 0x15,
		0x3E, 0x00, 0x00, 0x00, 0x06, 0xFF,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x66, 0x00,
		0x02,
		0x8a,
		0x05,
		0x81, 0x11, 0x3E, 0x00, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x11, 0x32, 0x09, 0x07, 0x04, 0x15,
		0x81, 0x11, 0x3E, 0x00, 0x00, 0x00, 0x01, 0x02, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x11, 0x32, 0x09, 0x07, 0x04, 0x15,
		0x81, 0x11, 0x3E, 0x00, 0x00, 0x00, 0x01, 0x03, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x11, 0x32, 0x09, 0x07, 0x04, 0x15,
		0x81, 0x11, 0x3E, 0x00, 0x00, 0x00, 0x01, 0x04, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x11, 0x32, 0x09, 0x07, 0x04, 0x15,
		0x81, 0x11, 0x3E, 0x00, 0x00, 0x00, 0x01, 0x04, 0x01, 0x00, 0x00, 0x04, 0x00, 0x00, 0x11, 0x32, 0x09, 0x07, 0x04, 0x15,
		0xBF,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadEnergy)
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

// 上报时间
func uploadTimeFun() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	var sendMsg = []byte{
		0x40, 0x40,
		0x00, 0x00,
		0x02, 0x03,
		0x3a, 0x01, 0x0f, 0x08, 0x06, 0x16,
		0x58, 0x13, 0x97, 0x10, 0xa1, 0xff,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x0a, 0x00,
		0x02,
		0x08, 0x01, 0x00, 0x00, 0x3a, 0x01, 0x0f, 0x08, 0x06, 0x16,
		0xa8,
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

// 心跳
func loginFun() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	var login = []byte{
		0x40, 0x40,
		//		0x00, 0x00,
		//		0x02, 0x03,
		0x3a, 0x01, 0x0f, 0x08, 0x06, 0x16,
		//		0x58, 0x13, 0x97, 0x10, 0xa1, 0xff,
		//		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		//		0x0a, 0x00,
		0x02, 0x00,
		//		0x08, 0x01, 0x00, 0x00, 0x3a, 0x01, 0x0f, 0x08, 0x06, 0x16,
		//		0xa8,
		0x23, 0x23,
	}
	for {

		dp := firenet.NewDataPack(2048)
		msg, err := dp.Pack(login)
		if err != nil {
			klog.Errorf("client write err: ", err)
			return
		}
		_, err = conn.Write(msg)
		if err != nil {
			klog.Errorf("client write err: ", err)
			return
		}
		klog.Infof("Client send ---------------Login----success--------------------")
		klog.Infof("Client recv -----------------Login---start-------------------")
		//1 先读出流中的head部分
		startData := make([]byte, fire.StartSignLen)
		_, err = io.ReadFull(conn, startData) //ReadFull 会把msg填充满为止
		if err != nil {
			klog.Errorf("read head error")
			return
		}
		klog.Infof("==>Client Recv Login msg:startData=%x", startData)
		//2 在读原地址等信息
		controData := make([]byte, 8)
		_, err = io.ReadFull(conn, controData) //ReadFull 会把msg填充满为止
		if err != nil {
			klog.Errorf("read head error")
			return
		}
		klog.Infof("==>Client Recv Login msg:controData=%x", controData)
		var end []byte
		end = make([]byte, fire.EndSignLen)
		if _, err := io.ReadFull(conn, end); err != nil {
			klog.Errorf("Client read msg data error ", err)
			return
		}

		klog.Infof("==>Client Recv Login EndData:%x", end)
		klog.Infof("Client recv --------------------END--Login-----------------")
		time.Sleep(time.Second)
	}
	conn.Close()
}

// 平台查询终端软件版本号
func uploadSystemctlVersionFunc() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	//电能
	var uploadAnalog = []byte{
		0x40, 0x40,
		0x96, 0x00,
		0x02, 0x03,
		0x10, 0x2C, 0x0E, 0x0F, 0x06, 0x15,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x15, 0x27, 0x08, 0x20, 0x20, 0xFE,
		0x17, 0x00,
		0x01,
		0x81, 0x01, 0x81, 0x01, 0x11, 0x15, 0x27, 0x08,
		0x20, 0x06, 0x00, 0x41, 0x54, 0x2B, 0x56, 0x45,
		0x52, 0x10, 0x2C, 0x0E, 0x0F, 0x06, 0x15,
		0x49,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息
func uploadSystemctlPRTFunc() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	var paly = []byte{0x39, 0x30, 0x30, 0xC2, 0xA5, 0x30, 0xB2, 0xE3, 0x30, 0xB7, 0xBF, 0xBC, 0xE4}
	fmt.Println(string(paly))

	var uploadAnalog = []byte{
		0x40, 0x40,
		0x01, 0x00,
		0x02, 0x03,
		0x22, 0x29, 0x0B, 0x10, 0x06, 0x15,
		0x10, 0x46, 0x94, 0x50, 0x61, 0xFF,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x30, 0x00,
		0x02,
		0x09,
		0x01,
		0x00,
		0x00,
		0x33, 0x33,
		0x05,
		0x0F,
		0xB8, 0xD0, 0xD1, 0xCC, 0xCC, 0xBD, 0xB2, 0xE2, 0xC6, 0xF7, 0x30, 0x30, 0x30, 0x30, 0x31,
		0x0D,
		0x39, 0x30, 0x30, 0xC2, 0xA5, 0x30, 0xB2, 0xE3, 0x30, 0xB7, 0xBF, 0xBC, 0xE4,
		0x04, 0xB9, 0xCA, 0xD5, 0xCF,
		0x00, 0x11, 0x0E, 0x11, 0x09, 0x12,
		0x55,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡CRT信息
func uploadSystemctlCRTFunc() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	var paly = []byte{0x39, 0x30, 0x30, 0xC2, 0xA5, 0x30, 0xB2, 0xE3, 0x30, 0xB7, 0xBF, 0xBC, 0xE4}
	fmt.Println(string(paly))
	var uploadAnalog = []byte{
		0x40, 0x40,
		0x01, 0x00,
		0x02, 0x03,
		0x1D, 0x34, 0x0B, 0x10, 0x06, 0x15,
		0x10, 0x46, 0x94, 0x50, 0x61, 0xFF,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x39, 0x00,
		0x02,
		0x0A,
		0x01,
		0x00,
		0x00,
		0x33,
		0x04,
		0x03,
		0x01, 0x00,
		0x01,
		0x1C,
		0x26, 0x26, 0x51, 0x4E, 0x30, 0x31, 0x2C, 0xBB,
		0xFA, 0xC6, 0xF7, 0x31, 0x31, 0x2C, 0xBB, 0xD8,
		0xC2, 0xB7, 0x33, 0x41, 0x2C, 0xB2, 0xBF, 0xBC,
		0xFE, 0x32, 0x43, 0x44,
		0x01, 0x08,
		0xB1, 0xBE, 0xBB, 0xFA, 0xBB, 0xF0, 0xBE, 0xAF,
		0x01, 0x00,
		0x1D, 0x34, 0x0B, 0x10, 0x06, 0x15,
		0x9A,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息
func uploadSystemctlPRTFunc2() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40,
		0x01, 0x00,
		0x02, 0x03,
		0x22, 0x29, 0x0B, 0x10, 0x06, 0x15,
		0x10, 0x46, 0x94, 0x50, 0x61, 0xFF,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x30, 0x00,
		0x02,
		0x09,
		0x01,
		0x00,
		0x00,
		0x33, 0x33,
		0x05,
		0x0F,
		0xB8, 0xD0, 0xD1, 0xCC, 0xCC, 0xBD, 0xB2, 0xE2, 0xC6, 0xF7, 0x30, 0x30, 0x30, 0x30, 0x31,
		0x0D,
		0x39, 0x30, 0x30, 0xC2, 0xA5, 0x30, 0xB2, 0xE3, 0x30, 0xB7, 0xBF, 0xBC, 0xE4,
		0x04, 0xB9, 0xCA, 0xD5, 0xCF,
		0x00, 0x11, 0x0E, 0x11, 0x09, 0x12,
		0x55,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息
func uploadSystemctlPRTFunc3() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40,
		0x05, 0x00,
		0x02, 0x03,
		0x07, 0x07, 0x0A, 0x17, 0x09, 0x16,
		0x58, 0x13, 0x97, 0x10, 0xA1, 0xFF,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x40, 0x00,
		0x02,
		0x09,
		0x01,
		0x00,
		0x05,
		0x33, 0x33,
		0x09,
		0x0A,
		0xB5, 0xE7, 0xD4, 0xB4, 0x3A, 0x31, 0x33, 0x31, 0x38, 0x38,
		0x22,
		0xA3, 0xD3, 0xA3, 0xD4, 0xA3, 0xCF, 0xA3, 0xD2, 0xA3, 0xCF,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x30, 0x30, 0xC2, 0xA5, 0x30, 0x34, 0xB2, 0xE3, 0x30, 0x30,
		0xB7, 0xBF, 0xBC, 0xE4,
		0x04,
		0xB7, 0xB4, 0xC0, 0xA1,
		0x08, 0x22, 0x0B, 0x08, 0x06, 0x16,
		0x3E,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息
func uploadSystemctlPRTFunc4() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40,
		0x07, 0x00,
		0x02, 0x03,
		0x10, 0x27, 0x12, 0x0B, 0x04, 0x17,
		0x58, 0x13, 0x97, 0x10, 0xA1, 0xFF,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x2D, 0x00,
		0x02,
		0x09,
		0x01,
		0x00,
		0x00,
		0x33, 0x33,
		0x01,
		0x13,
		0xC3, 0xC5, 0xBD, 0xFB, 0x20, 0x20, 0x20, 0x20, 0xCA, 0xCD, 0xB7, 0xC5, 0x5F,
		0x39, 0x39, 0x30, 0x31, 0x34, 0x31,
		0x00,
		0x0A,
		0xCE, 0xDE, 0xC6, 0xA5, 0xC5, 0xE4, 0xC0, 0xE0, 0xD0, 0xCD,
		0x0F, 0x27, 0x12, 0x0B, 0x04, 0x17,
		0xBF,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息
func uploadSystemctlPRTFunc5() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40,
		0x16, 0x00,
		0x02, 0x03,
		0x02, 0x18, 0x10, 0x17, 0x09, 0x16,
		0x58, 0x13, 0x97, 0x10, 0xA1, 0xFF,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x3B, 0x00,
		0x02,
		0x09,
		0x01,
		0x00,
		0x00,
		0x33, 0x33,
		0x01,
		0x09, 0xB5, 0xE7, 0xD4, 0xB4, 0x31, 0x30, 0x31, 0x33, 0x31,
		0x1E, 0xB5, 0xCD, 0xD1, 0xB9, 0xA3, 0xB2, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x30, 0x30, 0xC2, 0xA5, 0x30, 0x34, 0xB2,
		0xE3, 0x30, 0x30, 0xB7, 0xBF, 0xBC, 0xE4,
		0x04, 0xB7, 0xB4, 0xC0, 0xA1,
		0x01, 0x1F, 0x0B, 0x08, 0x06, 0x16, 0x12,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息-M6V复位
func uploadSystemctlPRTFuncResetM6V() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40, 0xe5, 0x00, 0x02, 0x03, 0x19, 0x25, 0x03, 0x0e, 0x06, 0x17, 0x24, 0x25, 0x51, 0x80, 0x12,
		0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x02,
		0x09,
		0x01,
		0x00,
		0x00,
		0x33, 0x33,
		0x01,
		0x0b, 0xb8, 0xb4, 0xce, 0xbb, 0x5f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x04, 0xb8, 0xb4, 0xce, 0xbb,
		0x19, 0x25, 0x03, 0x0e, 0x06, 0x17,
		0xd7,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息-复位
func uploadSystemctlPRTFuncResetsoft() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40, 0xe5, 0x00, 0x02, 0x03, 0x19, 0x25, 0x03, 0x0e, 0x06, 0x17, 0x24, 0x25, 0x51, 0x80, 0x12,
		0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x02,
		0x09,
		0x01,
		0x00,
		0x00,
		0x33, 0x33,
		0x01,
		0x0b, 0xb8, 0xb4, 0xce, 0xbb, 0x5f, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x00,
		0x04, 0xb8, 0xb4, 0xce, 0xbb,
		0x19, 0x25, 0x03, 0x0e, 0x06, 0x17,
		0xf7,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡CRT信息tengren
func uploadSystemctlCRTFunctengren() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{

		0x40, 0x40, 0x37, 0x02, 0x02, 0x03, 0x11, 0x25, 0x02, 0x19, 0x08, 0x17, 0x58, 0x13, 0x97, 0x10,

		0xA1, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33,

		0x33, 0x01, 0x00, 0x1C, 0x26, 0x26, 0x51, 0x4E, 0x30, 0x31, 0x2C, 0xBB, 0xFA, 0xC6, 0xF7, 0x30,

		0x31, 0x2C, 0xBB, 0xD8, 0xC2, 0xB7, 0x30, 0x35, 0x2C, 0xB2, 0xBF, 0xBC, 0xFE, 0x30, 0x35, 0x31,

		0x08, 0xC9, 0xE8, 0xB1, 0xB8, 0xBB, 0xD8, 0xB4, 0xF0, 0x32, 0x24, 0x0A, 0x19, 0x08, 0x17, 0xE9,

		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡CRT信息tengren
func uploadSystemctlCRTFunctengren2() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{

		0x40, 0x40, 0x36, 0x02, 0x02, 0x03, 0x11, 0x25, 0x02, 0x19, 0x08, 0x17, 0x58, 0x13, 0x97, 0x10,

		0xA1, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x30, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33,

		0x33, 0x03, 0x00, 0x1C, 0x26, 0x26, 0x51, 0x4E, 0x30, 0x31, 0x2C, 0xBB, 0xFA, 0xC6, 0xF7, 0x32,

		0x30, 0x2C, 0xBB, 0xD8, 0xC2, 0xB7, 0x30, 0x33, 0x2C, 0xB2, 0xBF, 0xBC, 0xFE, 0x30, 0x30, 0x33,

		0x04, 0xBB, 0xF0, 0xBE, 0xAF, 0x32, 0x24, 0x0A, 0x19, 0x08, 0x17, 0xA5, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡CRT信息tengrenReSet
func uploadSystemctlCRTFunctengrenReSet() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{

		0x40, 0x40,
		0x08, 0x00,
		0x02, 0x03,
		0x37, 0x33, 0x07, 0x04, 0x08, 0x17,
		0x58, 0x13, 0x97, 0x10, 0xA1, 0xFF,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x17, 0x00,
		0x02,
		0x0A,
		0x01,
		0x00,
		0x00,
		0x33,
		0x03,
		0x01,
		0x00,
		0x00,
		0x01,
		0x06,
		0xB8, 0xB4, 0xCE, 0xBB, 0x0D, 0x0A,
		0x19, 0x33, 0x0F, 0x04, 0x08, 0x17,
		0x3F,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息Jiyun
func uploadSystemctlPRTFuncJiyun() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{

		0x40, 0x40,
		0x04, 0x00,
		0x02, 0x03,
		0x1D, 0x0D, 0x0B, 0x08, 0x08, 0x17,
		0x10, 0x25, 0x51, 0x80, 0x12, 0xFF,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x48, 0x00,
		0x02,
		0x09,
		0x01,
		0x00,
		0x00,
		0x33, 0x33,
		0x09,
		0x09,
		0xCF, 0xFB, 0xB7, 0xC0, 0xB9, 0xE3, 0xB2, 0xA5, 0x50,
		0x19,
		0xB6, 0xE0, 0xCF, 0xDF, 0xB9, 0xE3, 0xB2, 0xA5,
		0xC5, 0xCC, 0x30, 0xC7, 0xF8, 0xD3, 0xF2, 0x30,
		0xC2, 0xA5, 0xB2, 0xE3, 0x30, 0xB7, 0xBF, 0xBC,
		0xE4,
		0x16,
		0x30, 0x31, 0x31, 0xBB, 0xD8, 0xC2, 0xB7, 0x30,
		0x30, 0x32, 0xB5, 0xD8, 0xD6, 0xB7, 0xC1, 0xAA,
		0xB6, 0xAF, 0xC7, 0xEB, 0xC7, 0xF3,
		0x00, 0x00, 0x0F, 0x07, 0x08, 0x17,
		0xA4,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息Jiyun
func uploadSystemctlPRTFuncJiyun2() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{

		0x40, 0x40,
		0x05, 0x00,
		0x02, 0x03,
		0x24, 0x0D, 0x0B, 0x08, 0x08, 0x17,
		0x10, 0x25, 0x51, 0x80, 0x12, 0xFF,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x48, 0x00,
		0x02,
		0x09,
		0x01,
		0x00,
		0x00,
		0x33, 0x33,
		0x09,
		0x09,
		0xCF, 0xFB, 0xB7, 0xC0, 0xB9, 0xE3, 0xB2, 0xA5, 0x50,
		0x19,
		0xB6, 0xE0, 0xCF, 0xDF, 0xB9, 0xE3, 0xB2, 0xA5,
		0xC5, 0xCC, 0x30, 0xC7, 0xF8, 0xD3, 0xF2, 0x30,
		0xC2, 0xA5, 0xB2, 0xE3, 0x30, 0xB7, 0xBF, 0xBC,
		0xE4,
		0x16,
		0x30, 0x31, 0x31, 0xBB, 0xD8, 0xC2, 0xB7, 0x30,
		0x30, 0x33, 0xB5, 0xD8, 0xD6, 0xB7, 0xC1, 0xAA,
		0xB6, 0xAF, 0xC7, 0xEB, 0xC7, 0xF3,
		0x00, 0x00, 0x0F, 0x07, 0x08, 0x17,
		0xAD,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息Jiyun
func uploadSystemctlPRTTongji() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{

		0x40, 0x40,
		0x04, 0x00,
		0x02, 0x03,
		0x21, 0x3A, 0x11, 0x08, 0x08, 0x17,
		0x10, 0x25, 0x51, 0x80, 0x12, 0xFF,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x3A, 0x00,
		0x02,
		0x09,
		0x01,
		0x00,
		0x00,
		0x33, 0x33,
		0x03,
		0x08,
		0xCA, 0xD6, 0xB6, 0xAF, 0xB0, 0xB4, 0xC5, 0xA5,
		0x0C,
		0xD2, 0xBB, 0xB2, 0xE3, 0xD4, 0xCB, 0xCE, 0xAC,
		0xCA, 0xD2, 0xCD, 0xE2,
		0x16,
		0x30, 0x30, 0x35, 0xBB, 0xD8, 0xC2, 0xB7, 0x30,
		0x30, 0x31, 0xB5, 0xD8, 0xD6, 0xB7, 0xCA, 0xD7,
		0xB4, 0xCE, 0xBB, 0xF0, 0xBE, 0xAF,
		0x00, 0x1B, 0x11, 0x08, 0x08, 0x17,
		0xBF,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡CRT信息tengrenReSet
func uploadSystemctlCRTFuncTongjiReSet() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{

		0x40, 0x40,
		0x03, 0x00, 0x02, 0x03, 0x00, 0x39, 0x11, 0x08, 0x08, 0x17, 0x10, 0x25, 0x51, 0x80,
		0x12, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33,
		0x33, 0x00, 0x00, 0x00, 0x08, 0xCF, 0xB5, 0xCD, 0xB3, 0xB8, 0xB4, 0xCE, 0xBB, 0x00, 0x19, 0x11,
		0x08, 0x08, 0x17, 0x6C,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡CRT信息Boxing
func uploadSystemctlCRTFuncBoxing() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40,
		0x98, 0x01,
		0x02, 0x03,
		0x29, 0x20, 0x06, 0x11, 0x08, 0x17,
		0x30, 0x25, 0x51, 0x80, 0x12, 0xff,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x33, 0x00,
		0x02, //发送
		0x0a, //类型标识
		0x01, //信息对象数目
		0x00, //通道号
		0x00, //消空主机编号
		0x33, //预
		0x04, //类型个数
		0x05, //报警类型
		0x00, //空    1个
		0x01, //字符串 2个
		0x1c, //长度28
		0x26, 0x26, 0x51, 0x4e, 0x30, 0x31,
		0x2c,
		0xbb, 0xfa, 0xc6, 0xf7, 0x30, 0x31,
		0x2c,
		0xbb, 0xd8, 0xc2, 0xb7, 0x31, 0x36, //回路16
		0x2c,
		0xb2, 0xbf, 0xbc, 0xfe, 0x31, 0x30, 0x31, //部件101
		0x01, // 3个
		0x04,
		0xb9, 0xca, 0xd5, 0xcf,
		0x00, // 4个
		0x00, 0x21, 0x0e, 0x11, 0x08, 0x17,
		0x4b,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息-复位
func uploadSystemctlPRTFuncResetM3() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40,
		0x02, 0x00, 0x02, 0x03, 0x07, 0x0A, 0x03, 0x16, 0x08, 0x17, 0x08, 0x25, 0x51, 0x80,
		0x12, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33,
		0x33, 0x00, 0x00, 0x00, 0x08, 0xCF, 0xB5, 0xCD, 0xB3, 0xB8, 0xB4, 0xCE, 0xBB, 0x00, 0x0E, 0x0B,
		0x16, 0x08, 0x17, 0x38,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息M3
func uploadSystemctlPRTFuncM3() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{

		0x40, 0x40,
		0x03, 0x00,
		0x02, 0x03,
		0x0B, 0x0B, 0x03, 0x16, 0x08, 0x17,
		0x08, 0x25, 0x51, 0x80, 0x12, 0xFF,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x3E, 0x00,
		0x02,
		0x09,
		0x01,
		0x00,
		0x00,
		0x33, 0x33,
		0x03,
		0x0A, 0xB8, 0xD0, 0xD1, 0xCC, 0xCC, 0xBD, 0xB2, 0xE2, 0xC6, 0xF7, 0x0E, 0xA3, 0xB1,
		0xB2, 0xE3, 0xA3, 0xB2, 0xBA, 0xC5, 0xBC, 0xAF, 0xD7, 0xB0, 0xCF, 0xE4, 0x16, 0x30, 0x30, 0x36,
		0xBB, 0xD8, 0xC2, 0xB7, 0x30, 0x30, 0x36, 0xB5, 0xD8, 0xD6, 0xB7, 0xCA, 0xD7, 0xB4, 0xCE, 0xBB,
		0xF0, 0xBE, 0xAF, 0x00, 0x16, 0x0B, 0x0C, 0x08, 0x17, 0x80,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息M3-2
func uploadSystemctlPRTFuncM32() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{

		0x40, 0x40, 0x0B, 0x00, 0x02, 0x03, 0x1D, 0x0B, 0x03, 0x16, 0x08, 0x17, 0x08, 0x25, 0x51, 0x80,
		0x12, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3E, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33,
		0x33, 0x03, 0x0A, 0xB8, 0xD0, 0xD1, 0xCC, 0xCC, 0xBD, 0xB2, 0xE2, 0xC6, 0xF7, 0x0E, 0xC2, 0xA5,
		0xB6, 0xA5, 0xA3, 0xB2, 0xBA, 0xC5, 0xBC, 0xAF, 0xD7, 0xB0, 0xCF, 0xE4, 0x16, 0x30, 0x30, 0x36,
		0xBB, 0xD8, 0xC2, 0xB7, 0x30, 0x32, 0x36, 0xB5, 0xD8, 0xD6, 0xB7, 0xCA, 0xD7, 0xB4, 0xCE, 0xBB,
		0xF0, 0xBE, 0xAF, 0x00, 0x15, 0x14, 0x14, 0x08, 0x17, 0x85, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息  经开
func uploadSystemctlPRTFuncXIANJK() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{

		0x40, 0x40, 0x05, 0x00, 0x02, 0x03, 0x07, 0x29, 0x02, 0x18, 0x08, 0x17, 0x21, 0x25, 0x51, 0x80,
		0x12, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xB2, 0x01, 0x02, 0x09, 0x06, 0x00, 0x00, 0x33,
		0x33, 0x00, 0x08, 0xC6, 0xF8, 0xCC, 0xE5, 0xC3, 0xF0, 0xBB, 0xF0, 0x1E, 0xD2, 0xBB, 0xB2, 0xE3,
		0xB6, 0xAB, 0xD3, 0xCD, 0xCF, 0xE4, 0xD4, 0xA4, 0xBE, 0xAF, 0x30, 0x30, 0x30, 0xC7, 0xF8, 0x30,
		0x30, 0x31, 0xB2, 0xE3, 0x30, 0x30, 0x30, 0x30, 0xCA, 0xD2, 0x16, 0x30, 0x30, 0x31, 0xBB, 0xD8,
		0xC2, 0xB7, 0x31, 0x32, 0x30, 0xB5, 0xD8, 0xD6, 0xB7, 0xC5, 0xE7, 0xC8, 0xF7, 0xBD, 0xFB, 0xD6,
		0xB9, 0x00, 0x2C, 0x0A, 0x18, 0x08, 0x17, 0x00, 0x00, 0x33, 0x33, 0x00, 0x08, 0xC6, 0xF8, 0xCC,
		0xE5, 0xC3, 0xF0, 0xBB, 0xF0, 0x20, 0xD2, 0xBB, 0xB2, 0xE3, 0xA3, 0xD5, 0xA3, 0xD0, 0xA3, 0xD3,
		0xA3, 0xB2, 0xD4, 0xA4, 0xBE, 0xAF, 0x30, 0x30, 0x30, 0xC7, 0xF8, 0x30, 0x30, 0x31, 0xB2, 0xE3,
		0x30, 0x30, 0x30, 0x30, 0xCA, 0xD2, 0x16, 0x30, 0x30, 0x31, 0xBB, 0xD8, 0xC2, 0xB7, 0x31, 0x33,
		0x34, 0xB5, 0xD8, 0xD6, 0xB7, 0xC5, 0xE7, 0xC8, 0xF7, 0xBD, 0xFB, 0xD6, 0xB9, 0x00, 0x2C, 0x0A,
		0x18, 0x08, 0x17, 0x00, 0x00, 0x33, 0x33, 0x00, 0x08, 0xC6, 0xF8, 0xCC, 0xE5, 0xC3, 0xF0, 0xBB,
		0xF0, 0x1E, 0xA3, 0xD5, 0xA3, 0xD0, 0xA3, 0xD3, 0xC5, 0xE4, 0xB5, 0xE7, 0xCA, 0xD2, 0xCE, 0xE5,
		0x30, 0x30, 0x30, 0xC7, 0xF8, 0x30, 0x30, 0x30, 0xB2, 0xE3, 0x30, 0x30, 0x30, 0x30, 0xCA, 0xD2,
		0x16, 0x30, 0x30, 0x37, 0xBB, 0xD8, 0xC2, 0xB7, 0x31, 0x32, 0x31, 0xB5, 0xD8, 0xD6, 0xB7, 0xC5,
		0xE7, 0xC8, 0xF7, 0xBD, 0xFB, 0xD6, 0xB9, 0x00, 0x2C, 0x0A, 0x18, 0x08, 0x17, 0x00, 0x00, 0x33,
		0x33, 0x00, 0x08, 0xC6, 0xF8, 0xCC, 0xE5, 0xC3, 0xF0, 0xBB, 0xF0, 0x16, 0xBB, 0xFA, 0xB7, 0xBF,
		0xBE, 0xC5, 0x30, 0x30, 0x30, 0xC7, 0xF8, 0x30, 0x30, 0x30, 0xB2, 0xE3, 0x30, 0x30, 0x30, 0x30,
		0xCA, 0xD2, 0x16, 0x30, 0x30, 0x37, 0xBB, 0xD8, 0xC2, 0xB7, 0x31, 0x32, 0x33, 0xB5, 0xD8, 0xD6,
		0xB7, 0xC5, 0xE7, 0xC8, 0xF7, 0xBD, 0xFB, 0xD6, 0xB9, 0x00, 0x2C, 0x0A, 0x18, 0x08, 0x17, 0x00,
		0x00, 0x33, 0x33, 0x00, 0x08, 0xC6, 0xF8, 0xCC, 0xE5, 0xC3, 0xF0, 0xBB, 0xF0, 0x20, 0xA3, 0xD5,
		0xA3, 0xD0, 0xA3, 0xD3, 0xC5, 0xE4, 0xB5, 0xE7, 0xC8, 0xFD, 0xBA, 0xCD, 0xCB, 0xC4, 0x30, 0x30,
		0x30, 0xC7, 0xF8, 0x30, 0x30, 0x30, 0xB2, 0xE3, 0x30, 0x30, 0x30, 0x30, 0xCA, 0xD2, 0x16, 0x30,
		0x30, 0x37, 0xBB, 0xD8, 0xC2, 0xB7, 0x31, 0x32, 0x37, 0xB5, 0xD8, 0xD6, 0xB7, 0xC5, 0xE7, 0xC8,
		0xF7, 0xBD, 0xFB, 0xD6, 0xB9, 0x00, 0x2C, 0x0A, 0x18, 0x08, 0x17, 0x00, 0x00, 0x33, 0x33, 0x00,
		0x08, 0xC6, 0xF8, 0xCC, 0xE5, 0xC3, 0xF0, 0xBB, 0xF0, 0x16, 0xBB, 0xFA, 0xB7, 0xBF, 0xCA, 0xAE,
		0x30, 0x30, 0x30, 0xC7, 0xF8, 0x30, 0x30, 0x30, 0xB2, 0xE3, 0x30, 0x30, 0x30, 0x30, 0xCA, 0xD2,
		0x16, 0x30, 0x30, 0x37, 0xBB, 0xD8, 0xC2, 0xB7, 0x31, 0x33, 0x33, 0xB5, 0xD8, 0xD6, 0xB7, 0xC5,
		0xE7, 0xC8, 0xF7, 0xBD, 0xFB, 0xD6, 0xB9, 0x00, 0x2C, 0x0A, 0x18, 0x08, 0x17, 0xB0, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息  经开
func uploadSystemctlPRTFuncXIANJK2() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{

		0x40, 0x40, 0x02, 0x00, 0x02, 0x03, 0x0F, 0x24,
		0x02, 0x18, 0x08, 0x17, 0x21, 0x25, 0x51, 0x80,
		0x12, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xAA, 0x01, 0x02, 0x09, 0x06, 0x00, 0x00, 0x33,
		0x33, 0x00, 0x08, 0xC6, 0xF8, 0xCC, 0xE5, 0xC3,
		0xF0, 0xBB, 0xF0, 0x20, 0xD2, 0xBB, 0xB2, 0xE3,
		0xA3, 0xD5, 0xA3, 0xD0, 0xA3, 0xD3, 0xA3, 0xB2,
		0xD4, 0xA4, 0xBE, 0xAF, 0x30, 0x30, 0x30, 0xC7,
		0xF8, 0x30, 0x30, 0x31, 0xB2, 0xE3, 0x30, 0x30,
		0x30, 0x30, 0xCA, 0xD2, 0x16, 0x30, 0x30, 0x31,
		0xBB, 0xD8, 0xC2, 0xB7, 0x31, 0x33, 0x34, 0xB5,
		0xD8, 0xD6, 0xB7, 0xC5, 0xE7, 0xC8, 0xF7, 0xBD,
		0xFB, 0xD6, 0xB9, 0x00, 0x27, 0x0A, 0x18, 0x08,
		0x17, 0x00, 0x00, 0x33, 0x33, 0x00, 0x08, 0xC6,
		0xF8, 0xCC, 0xE5, 0xC3, 0xF0, 0xBB, 0xF0, 0x1E,
		0xA3, 0xD5, 0xA3, 0xD0, 0xA3, 0xD3, 0xC5, 0xE4,
		0xB5, 0xE7, 0xCA, 0xD2, 0xCE, 0xE5, 0x30, 0x30,
		0x30, 0xC7, 0xF8, 0x30, 0x30, 0x30, 0xB2, 0xE3,
		0x30, 0x30, 0x30, 0x30, 0xCA, 0xD2, 0x16, 0x30,
		0x30, 0x37, 0xBB, 0xD8, 0xC2, 0xB7, 0x31, 0x32,
		0x31, 0xB5, 0xD8, 0xD6, 0xB7, 0xC5, 0xE7, 0xC8,
		0xF7, 0xBD, 0xFB, 0xD6, 0xB9, 0x00, 0x27, 0x0A,
		0x18, 0x08, 0x17, 0x00, 0x00, 0x33, 0x33, 0x00,
		0x08, 0xC6, 0xF8, 0xCC, 0xE5, 0xC3, 0xF0, 0xBB,
		0xF0, 0x16, 0xBB, 0xFA, 0xB7, 0xBF, 0xBE, 0xC5,
		0x30, 0x30, 0x30, 0xC7, 0xF8, 0x30, 0x30, 0x30,
		0xB2, 0xE3, 0x30, 0x30, 0x30, 0x30, 0xCA, 0xD2,
		0x16, 0x30, 0x30, 0x37, 0xBB, 0xD8, 0xC2, 0xB7,
		0x31, 0x32, 0x33, 0xB5, 0xD8, 0xD6, 0xB7, 0xC5,
		0xE7, 0xC8, 0xF7, 0xBD, 0xFB, 0xD6, 0xB9, 0x00,
		0x27, 0x0A, 0x18, 0x08, 0x17, 0x00, 0x00, 0x33,
		0x33, 0x00, 0x08, 0xC6, 0xF8, 0xCC, 0xE5, 0xC3,
		0xF0, 0xBB, 0xF0, 0x20, 0xA3, 0xD5, 0xA3, 0xD0,
		0xA3, 0xD3, 0xC5, 0xE4, 0xB5, 0xE7, 0xC8, 0xFD,
		0xBA, 0xCD, 0xCB, 0xC4, 0x30, 0x30, 0x30, 0xC7,
		0xF8, 0x30, 0x30, 0x30, 0xB2, 0xE3, 0x30, 0x30,
		0x30, 0x30, 0xCA, 0xD2, 0x16, 0x30, 0x30, 0x37,
		0xBB, 0xD8, 0xC2, 0xB7, 0x31, 0x32, 0x37, 0xB5,
		0xD8, 0xD6, 0xB7, 0xC5, 0xE7, 0xC8, 0xF7, 0xBD,
		0xFB, 0xD6, 0xB9, 0x00, 0x27, 0x0A, 0x18, 0x08,
		0x17, 0x00, 0x00, 0x33, 0x33, 0x00, 0x08, 0xC6,
		0xF8, 0xCC, 0xE5, 0xC3, 0xF0, 0xBB, 0xF0, 0x16,
		0xBB, 0xFA, 0xB7, 0xBF, 0xCA, 0xAE, 0x30, 0x30,
		0x30, 0xC7, 0xF8, 0x30, 0x30, 0x30, 0xB2, 0xE3,
		0x30, 0x30, 0x30, 0x30, 0xCA, 0xD2, 0x16, 0x30,
		0x30, 0x37, 0xBB, 0xD8, 0xC2, 0xB7, 0x31, 0x33,
		0x33, 0xB5, 0xD8, 0xD6, 0xB7, 0xC5, 0xE7, 0xC8,
		0xF7, 0xBD, 0xFB, 0xD6, 0xB9, 0x00, 0x27, 0x0A,
		0x18, 0x08, 0x17, 0x00, 0x00, 0x33, 0x33, 0x00,
		0x08, 0xC6, 0xF8, 0xCC, 0xE5, 0xC3, 0xF0, 0xBB,
		0xF0, 0x16, 0xBB, 0xFA, 0xB7, 0xBF, 0xB0, 0xCB,
		0x30, 0x30, 0x30, 0xC7, 0xF8, 0x30, 0x30, 0x30,
		0xB2, 0xE3, 0x30, 0x30, 0x30, 0x30, 0xCA, 0xD2,
		0x16, 0x30, 0x30, 0x37, 0xBB, 0xD8, 0xC2, 0xB7,
		0x31, 0x33, 0x35, 0xB5, 0xD8, 0xD6, 0xB7, 0xC5,
		0xE7, 0xC8, 0xF7, 0xBD, 0xFB, 0xD6, 0xB9, 0x00,
		0x28, 0x0A, 0x18, 0x08, 0x17, 0x79, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息  经开复位
func uploadSystemctlPRTFuncXIANJKReset() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{

		0x40, 0x40, 0x01, 0x00, 0x02, 0x03, 0x1D, 0x23,
		0x02, 0x18, 0x08, 0x17, 0x21, 0x25, 0x51, 0x80,
		0x12, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x18, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33,
		0x33, 0x00, 0x00, 0x00, 0x08, 0xCF, 0xB5, 0xCD,
		0xB3, 0xB8, 0xB4, 0xCE, 0xBB, 0x00, 0x27, 0x0A,
		0x18, 0x08, 0x17, 0x9A, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息  外高桥复位
func uploadSystemctlPRTFuncWGQReset() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{

		0x40, 0x40, 0x01, 0x00, 0x02, 0x03, 0x1D, 0x23,
		0x02, 0x18, 0x08, 0x17, 0x21, 0x25, 0x51, 0x80,
		0x12, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x18, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33,
		0x33, 0x00, 0x00, 0x00, 0x08, 0xCF, 0xB5, 0xCD,
		0xB3, 0xB8, 0xB4, 0xCE, 0xBB, 0x00, 0x27, 0x0A,
		0x18, 0x08, 0x17, 0x9A, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息XIASHA 下沙
func uploadSystemctlPRTFuncXIASHA() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{

		0x40, 0x40, 0x09, 0x00, 0x02, 0x03, 0x2F, 0x03, 0x08, 0x19, 0x08, 0x17, 0x22, 0x25, 0x51, 0x80,
		0x12, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2D, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33,
		0x33, 0x03, 0x0F, 0xCA, 0xD6, 0xB6, 0xAF, 0xB0, 0xB4, 0xC5, 0xA5, 0x5F, 0x30, 0x32, 0x31, 0x31,
		0x35, 0x39, 0x0A, 0x39, 0x20, 0xB6, 0xFE, 0xB5, 0xE7, 0xCC, 0xDD, 0xCC, 0xFC, 0x04, 0xBB, 0xF0,
		0xBE, 0xAF, 0x2F, 0x03, 0x08, 0x19, 0x08, 0x17, 0x70, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息XIASHA 下沙
func uploadSystemctlPRTFuncXIASHA2() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40, 0x02, 0x00, 0x02, 0x03, 0x2E, 0x30, 0x07, 0x19, 0x08, 0x17, 0x22, 0x25, 0x51, 0x80,
		0x12, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2D, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33,
		0x33, 0x21, 0x0F, 0xBE, 0xED, 0xC1, 0xB1, 0xC3, 0xC5, 0xCF, 0xC2, 0x5F, 0x30, 0x31, 0x32, 0x33,
		0x31, 0x30, 0x0A, 0x37, 0x20, 0x20, 0x20, 0xB5, 0xE7, 0xCC, 0xDD, 0xCC, 0xFC, 0x04, 0xB7, 0xB4,
		0xC0, 0xA1, 0x2E, 0x30, 0x07, 0x19, 0x08, 0x17, 0x73, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息XIASHA 下沙
func uploadSystemctlPRTFuncXIASHA3() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40, 0x02, 0x00, 0x02, 0x03, 0x2E, 0x30, 0x07, 0x19, 0x08, 0x17, 0x22, 0x25, 0x51, 0x80,
		0x12, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2D, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33,
		0x33, 0x21, 0x0F, 0xBE, 0xED, 0xC1, 0xB1, 0xC3, 0xC5, 0xCF, 0xC2, 0x5F, 0x30, 0x31, 0x32, 0x33,
		0x31, 0x30, 0x0A, 0x37, 0x20, 0x20, 0x20, 0xB5, 0xE7, 0xCC, 0xDD, 0xCC, 0xFC, 0x04, 0xB7, 0xB4,
		0xC0, 0xA1, 0x2E, 0x30, 0x07, 0x19, 0x08, 0x17, 0x73, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息XIASHA 下沙
func uploadSystemctlPRTFuncXIASHA4() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40, 0x0A, 0x00, 0x02, 0x03, 0x30, 0x3B, 0x01, 0x08, 0x0C, 0x17, 0x22, 0x25, 0x51, 0x80,

		0x12, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2C, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33,

		0x33, 0x03, 0x0F, 0xCE, 0xFC, 0xC6, 0xF8, 0xD4, 0xA4, 0xBE, 0xAF, 0x5F, 0x39, 0x32, 0x33, 0x37,

		0x37, 0x35, 0x09, 0x39, 0x20, 0xB6, 0xFE, 0x39, 0x20, 0x32, 0x20, 0x33, 0x04, 0xBB, 0xF0, 0xBE,

		0xAF, 0x2F, 0x3B, 0x01, 0x08, 0x0C, 0x17, 0x32, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息XIASHA 下沙
func uploadSystemctlPRTFuncXIASHAReset() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40, 0x01, 0x00, 0x02, 0x03, 0x26, 0x30,
		0x07, 0x19, 0x08, 0x17, 0x22, 0x25, 0x51, 0x80,
		0x12, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x1F, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33,
		0x33, 0x01, 0x0B, 0xB8, 0xB4, 0xCE, 0xBB, 0x5F,
		0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x00, 0x04,
		0xB8, 0xB4, 0xCE, 0xBB, 0x25, 0x30, 0x07, 0x19,
		0x08, 0x17, 0x62, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息-waigaoqiao
func uploadSystemctlPRTFuncWaiGaoqiao() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40, 0x0A, 0x00, 0x02, 0x03, 0x08, 0x00, 0x0A, 0x06, 0x0B, 0x17, 0x06, 0x25, 0x51, 0x80,
		0x11, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xE6, 0x00, 0x02, 0x09, 0x02, 0x00, 0x00, 0x33,
		0x33, 0x03, 0x26, 0xCA, 0xD6, 0xB6, 0xAF, 0xB0, 0xB4, 0xC5, 0xA5, 0x5F, 0xBB, 0xFA, 0xC6, 0xF7,
		0x30, 0x31, 0x2C, 0xBB, 0xD8, 0xC2, 0xB7, 0x30, 0x34, 0x2C, 0xB2, 0xBF, 0xBC, 0xFE, 0x31, 0x36,
		0x35, 0x2C, 0xB7, 0xD6, 0xC7, 0xF8, 0x31, 0x31, 0x30, 0x36, 0xB2, 0xE3, 0xCE, 0xF7, 0xB2, 0xE0,
		0xB1, 0xB1, 0xD7, 0xDF, 0xB5, 0xC0, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x26, 0x26, 0x51, 0x4E,
		0x30, 0x31, 0x2C, 0xBB, 0xFA, 0xC6, 0xF7, 0x30, 0x31, 0x2C, 0xBB, 0xD8, 0xC2, 0xB7, 0x30, 0x34,
		0x2C, 0xB2, 0xBF, 0xBC, 0xFE, 0x31, 0x36, 0x35, 0x2C, 0xB7, 0xD6, 0xC7, 0xF8, 0x31, 0x31, 0x30,
		0x08, 0xBB, 0xF0, 0xD4, 0xD6, 0xB1, 0xA8, 0xBE, 0xAF, 0x00, 0x0B, 0x0F, 0x1B, 0x06, 0x17, 0x00,
		0x00, 0x33, 0x33, 0x03, 0x26, 0xB5, 0xE3, 0xD0, 0xCD, 0xB8, 0xD0, 0xD1, 0xCC, 0x5F, 0xBB, 0xFA,
		0xC6, 0xF7, 0x30, 0x31, 0x2C, 0xBB, 0xD8, 0xC2, 0xB7, 0x30, 0x34, 0x2C, 0xB2, 0xBF, 0xBC, 0xFE,
		0x31, 0x34, 0x33, 0x2C, 0xB7, 0xD6, 0xC7, 0xF8, 0x31, 0x31, 0x30, 0x36, 0xB2, 0xE3, 0xCE, 0xF7,
		0xB2, 0xE0, 0xB1, 0xB1, 0xD7, 0xDF, 0xB5, 0xC0, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x26, 0x26,
		0x51, 0x4E, 0x30, 0x31, 0x2C, 0xBB, 0xFA, 0xC6, 0xF7, 0x30, 0x31, 0x2C, 0xBB, 0xD8, 0xC2, 0xB7,
		0x30, 0x34, 0x2C, 0xB2, 0xBF, 0xBC, 0xFE, 0x31, 0x34, 0x33, 0x2C, 0xB7, 0xD6, 0xC7, 0xF8, 0x31,
		0x31, 0x30, 0x08, 0xBB, 0xF0, 0xD4, 0xD6, 0xB1, 0xA8, 0xBE, 0xAF, 0x00, 0x0B, 0x0F, 0x1B, 0x06,
		0x17, 0xB1, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息-waigaoqiao
func uploadSystemctlPRTFuncWaiGaoqiao2() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40, 0x09, 0x00, 0x02, 0x03, 0x07, 0x00, 0x0A, 0x06, 0x0B, 0x17, 0x06, 0x25, 0x51, 0x80,
		0x11, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xA2, 0x02, 0x02, 0x09, 0x06, 0x00, 0x00, 0x33,
		0x33, 0x03, 0x25, 0xB5, 0xE3, 0xD0, 0xCD, 0xB8, 0xD0, 0xD1, 0xCC, 0x5F, 0xBB, 0xFA, 0xC6, 0xF7,
		0x30, 0x31, 0x2C, 0xBB, 0xD8, 0xC2, 0xB7, 0x31, 0x32, 0x2C, 0xB2, 0xBF, 0xBC, 0xFE, 0x31, 0x30,
		0x2C, 0xB7, 0xD6, 0xC7, 0xF8, 0x31, 0x31, 0x34, 0x35, 0xB2, 0xE3, 0xCE, 0xF7, 0xB2, 0xE0, 0xB1,
		0xB1, 0xD7, 0xDF, 0xB5, 0xC0, 0xCE, 0xF7, 0x20, 0x20, 0x20, 0x20, 0x26, 0x26, 0x51, 0x4E, 0x30,
		0x31, 0x2C, 0xBB, 0xFA, 0xC6, 0xF7, 0x30, 0x31, 0x2C, 0xBB, 0xD8, 0xC2, 0xB7, 0x31, 0x32, 0x2C,
		0xB2, 0xBF, 0xBC, 0xFE, 0x31, 0x30, 0x2C, 0xB7, 0xD6, 0xC7, 0xF8, 0x31, 0x31, 0x34, 0x08, 0xBB,
		0xF0, 0xD4, 0xD6, 0xB1, 0xA8, 0xBE, 0xAF, 0x00, 0x0F, 0x0F, 0x1B, 0x06, 0x17, 0x00, 0x00, 0x33,
		0x33, 0x03, 0x25, 0xB5, 0xE3, 0xD0, 0xCD, 0xB8, 0xD0, 0xD1, 0xCC, 0x5F, 0xBB, 0xFA, 0xC6, 0xF7,
		0x30, 0x31, 0x2C, 0xBB, 0xD8, 0xC2, 0xB7, 0x31, 0x32, 0x2C, 0xB2, 0xBF, 0xBC, 0xFE, 0x30, 0x33,
		0x2C, 0xB7, 0xD6, 0xC7, 0xF8, 0x31, 0x31, 0x34, 0x35, 0xB2, 0xE3, 0xCE, 0xF7, 0xB2, 0xE0, 0xB1,
		0xB1, 0xD7, 0xDF, 0xB5, 0xC0, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x26, 0x26, 0x51, 0x4E, 0x30,
		0x31, 0x2C, 0xBB, 0xFA, 0xC6, 0xF7, 0x30, 0x31, 0x2C, 0xBB, 0xD8, 0xC2, 0xB7, 0x31, 0x32, 0x2C,
		0xB2, 0xBF, 0xBC, 0xFE, 0x30, 0x33, 0x2C, 0xB7, 0xD6, 0xC7, 0xF8, 0x31, 0x31, 0x34, 0x08, 0xBB,
		0xF0, 0xD4, 0xD6, 0xB1, 0xA8, 0xBE, 0xAF, 0x00, 0x0F, 0x0F, 0x1B, 0x06, 0x17, 0x00, 0x00, 0x33,
		0x33, 0x03, 0x24, 0xCA, 0xD6, 0xB6, 0xAF, 0xB0, 0xB4, 0xC5, 0xA5, 0x5F, 0xBB, 0xFA, 0xC6, 0xF7,
		0x30, 0x39, 0x2C, 0xBB, 0xD8, 0xC2, 0xB7, 0x30, 0x31, 0x2C, 0xB2, 0xBF, 0xBC, 0xFE, 0x30, 0x31,
		0x2C, 0xB7, 0xD6, 0xC7, 0xF8, 0x30, 0x30, 0x34, 0xB2, 0xE3, 0xCE, 0xF7, 0xB2, 0xE0, 0xB1, 0xB1,
		0xD7, 0xDF, 0xB5, 0xC0, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x26, 0x26, 0x51, 0x4E, 0x30, 0x31,
		0x2C, 0xBB, 0xFA, 0xC6, 0xF7, 0x30, 0x39, 0x2C, 0xBB, 0xD8, 0xC2, 0xB7, 0x30, 0x31, 0x2C, 0xB2,
		0xBF, 0xBC, 0xFE, 0x30, 0x31, 0x2C, 0xB7, 0xD6, 0xC7, 0xF8, 0x30, 0x30, 0x08, 0xBB, 0xF0, 0xD4,
		0xD6, 0xB1, 0xA8, 0xBE, 0xAF, 0x00, 0x0F, 0x0F, 0x1B, 0x06, 0x17, 0x00, 0x00, 0x33, 0x33, 0x03,
		0x25, 0xCA, 0xD6, 0xB6, 0xAF, 0xB0, 0xB4, 0xC5, 0xA5, 0x5F, 0xBB, 0xFA, 0xC6, 0xF7, 0x30, 0x31,
		0x2C, 0xBB, 0xD8, 0xC2, 0xB7, 0x31, 0x32, 0x2C, 0xB2, 0xBF, 0xBC, 0xFE, 0x32, 0x34, 0x2C, 0xB7,
		0xD6, 0xC7, 0xF8, 0x31, 0x31, 0x34, 0x35, 0xB2, 0xE3, 0xB1, 0xB1, 0xC2, 0xA5, 0xCC, 0xDD, 0xBC,
		0xE4, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x26, 0x26, 0x51, 0x4E, 0x30, 0x31, 0x2C,
		0xBB, 0xFA, 0xC6, 0xF7, 0x30, 0x31, 0x2C, 0xBB, 0xD8, 0xC2, 0xB7, 0x31, 0x32, 0x2C, 0xB2, 0xBF,
		0xBC, 0xFE, 0x32, 0x34, 0x2C, 0xB7, 0xD6, 0xC7, 0xF8, 0x31, 0x31, 0x34, 0x08, 0xBB, 0xF0, 0xD4,
		0xD6, 0xB1, 0xA8, 0xBE, 0xAF, 0x00, 0x0E, 0x0F, 0x1B, 0x06, 0x17, 0x00, 0x00, 0x33, 0x33, 0x03,
		0x25, 0xCA, 0xD6, 0xB6, 0xAF, 0xB0, 0xB4, 0xC5, 0xA5, 0x5F, 0xBB, 0xFA, 0xC6, 0xF7, 0x30, 0x31,
		0x2C, 0xBB, 0xD8, 0xC2, 0xB7, 0x31, 0x32, 0x2C, 0xB2, 0xBF, 0xBC, 0xFE, 0x32, 0x33, 0x2C, 0xB7,
		0xD6, 0xC7, 0xF8, 0x31, 0x31, 0x34, 0x35, 0xB2, 0xE3, 0xB1, 0xB1, 0xC2, 0xA5, 0xCC, 0xDD, 0xBC,
		0xE4, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x26, 0x26, 0x51, 0x4E, 0x30, 0x31, 0x2C,
		0xBB, 0xFA, 0xC6, 0xF7, 0x30, 0x31, 0x2C, 0xBB, 0xD8, 0xC2, 0xB7, 0x31, 0x32, 0x2C, 0xB2, 0xBF,
		0xBC, 0xFE, 0x32, 0x33, 0x2C, 0xB7, 0xD6, 0xC7, 0xF8, 0x31, 0x31, 0x34, 0x08, 0xBB, 0xF0, 0xD4,
		0xD6, 0xB1, 0xA8, 0xBE, 0xAF, 0x00, 0x0E, 0x0F, 0x1B, 0x06, 0x17, 0x00, 0x00, 0x33, 0x33, 0x03,
		0x26, 0xCA, 0xD6, 0xB6, 0xAF, 0xB0, 0xB4, 0xC5, 0xA5, 0x5F, 0xBB, 0xFA, 0xC6, 0xF7, 0x30, 0x31,
		0x2C, 0xBB, 0xD8, 0xC2, 0xB7, 0x30, 0x34, 0x2C, 0xB2, 0xBF, 0xBC, 0xFE, 0x31, 0x36, 0x33, 0x2C,
		0xB7, 0xD6, 0xC7, 0xF8, 0x31, 0x31, 0x30, 0x36, 0xB2, 0xE3, 0xCE, 0xF7, 0xB2, 0xE0, 0xB1, 0xB1,
		0xD7, 0xDF, 0xB5, 0xC0, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x26, 0x26, 0x51, 0x4E, 0x30, 0x31,
		0x2C, 0xBB, 0xFA, 0xC6, 0xF7, 0x30, 0x31, 0x2C, 0xBB, 0xD8, 0xC2, 0xB7, 0x30, 0x34, 0x2C, 0xB2,
		0xBF, 0xBC, 0xFE, 0x31, 0x36, 0x33, 0x2C, 0xB7, 0xD6, 0xC7, 0xF8, 0x31, 0x31, 0x30, 0x08, 0xBB,
		0xF0, 0xD4, 0xD6, 0xB1, 0xA8, 0xBE, 0xAF, 0x00, 0x0C, 0x0F, 0x1B, 0x06, 0x17, 0xBE, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息-waigaoqiao
func uploadSystemctlPRTFuncWaiGaoqiao3() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40, 0x67, 0x00, 0x02, 0x03, 0x3b, 0x35, 0x06, 0x1d, 0x0b, 0x17,
		0x06, 0x25, 0x51, 0x80, 0x11, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x5b, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33, 0x33, 0x09, 0x0c, 0xd4,
		0xa4, 0xd7, 0xf7, 0xd3, 0xc3, 0xb7, 0xa7, 0xd7, 0xe9, 0xbc, 0xe4, 0x21,
		0xbb, 0xfa, 0xba, 0xc5, 0x32, 0x2c, 0xbb, 0xd8, 0xc2, 0xb7, 0x31, 0x2c,
		0xb5, 0xe3, 0xba, 0xc5, 0x31, 0x34, 0x33, 0x2c, 0xc2, 0xa5, 0xba, 0xc5,
		0x31, 0x2c, 0xb7, 0xd6, 0xc7, 0xf8, 0x31, 0x38, 0x30, 0x1e, 0xd0, 0xc5,
		0xba, 0xc5, 0xb5, 0xfb, 0xb7, 0xa7, 0x20, 0xc5, 0xb5, 0xfb, 0xb7, 0xa7,
		0x2c, 0x20, 0xd7, 0xb4, 0xcc, 0xac, 0xa3, 0xba, 0xb7, 0xb4, 0xc0, 0xa1,
		0xbb, 0xd6, 0xb8, 0xb4, 0x3b, 0x35, 0x06, 0x1d, 0x0b, 0x17, 0x25, 0x23,
		0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机消息-同济
func uploadSystemctlPRTFuncTongjione() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40,
		0x4a, 0x04,
		0x02, 0x03,
		0x2b, 0x03, 0x09, 0x07, 0x0b, 0x17,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x18, 0x25, 0x51, 0x80, 0x12, 0xff,
		0x00, 0x00,
		0x03,
		0xd5,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机消息-同济
func uploadSystemctlPRTFuncTongji() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40,
		0x18, 0x25, 0x51, 0x80, 0x12, 0xff, 0x04, 0x88,
		0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机消息-m3
func uploadSystemctlPRTFunM3qimei() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40, 0x0F, 0x00, 0x02, 0x03, 0x3B, 0x15, 0x02, 0x0E, 0x0B, 0x17, 0x08, 0x25, 0x51, 0x80,
		0x12, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x43, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33,
		0x33, 0x09, 0x09, 0xC6, 0xF8, 0xCC, 0xE5, 0xC3, 0xF0, 0xBB, 0xF0, 0x50, 0x14, 0xA3, 0xC6, 0xA3,
		0xB3, 0xBF, 0xD5, 0xB5, 0xF7, 0xBB, 0xFA, 0xB7, 0xBF, 0xA3, 0xB1, 0xA3, 0xB0, 0xD4, 0xA4, 0xBE,
		0xAF, 0x16, 0x30, 0x30, 0x37, 0xBB, 0xD8, 0xC2, 0xB7, 0x30, 0x35, 0x31, 0xB5, 0xD8, 0xD6, 0xB7,
		0xC1, 0xAA, 0xB6, 0xAF, 0xC6, 0xF4, 0xB6, 0xAF, 0x00, 0x18, 0x0A, 0x0E, 0x0B, 0x17, 0xFD, 0x23,
		0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机消息-博兴空菜
func uploadSystemctlPRTFunBoXingKongcai() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40, 0xee, 0x00, 0x02, 0x03, 0x03, 0x2b, 0x02, 0x11, 0x0b, 0x17, 0x30, 0x25, 0x51, 0x80,
		0x12, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x37, 0x00, 0x02, 0x0a, 0x01, 0x00, 0x00, 0x33,
		0x04, 0x09, 0x00, 0x01, 0x1c, 0x26, 0x26, 0x51, 0x4e, 0x30, 0x31, 0x2c, 0xbb, 0xfa, 0xc6, 0xf7,
		0x30, 0x31, 0x2c, 0xbb, 0xd8, 0xc2, 0xb7, 0x31, 0x31, 0x2c, 0xb2, 0xbf, 0xbc, 0xfe, 0x30, 0x39,
		0x36, 0x01, 0x08, 0xd7, 0xdc, 0xcf, 0xdf, 0xbb, 0xd8, 0xb4, 0xf0, 0x00, 0x00, 0x2b, 0x0a, 0x11,
		0x0b, 0x17, 0x12, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机消息-纪云
func uploadSystemctlPRTFunjIyun333() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40,
		0x0b, 0x00,
		0x02, 0x03,
		0x09, 0x08, 0x02, 0x14, 0x0b, 0x17,
		0x11, 0x25, 0x51, 0x80, 0x12, 0xff,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x01,
		0x02,
		0x09,
		0x06,
		0x00,
		0x00, 0x33, 0x33, 0x00, 0x00, 0x1c,
		0x26, 0x26, 0x51, 0x4e, 0x30, 0x31,
		0x2c, 0xbb, 0xfa, 0xc6, 0xf7, 0x30,
		0x30, 0x2c, 0xbb, 0xd8, 0xc2, 0xb7,
		0x35, 0x35, 0x2c, 0xb2, 0xbf, 0xbc,
		0xfe, 0x31, 0x30, 0x32, 0x00, 0x08,
		0x08, 0x02, 0x14, 0x0b, 0x17, 0x00,
		0x00, 0x33, 0x33, 0x00, 0x00, 0x1d,
		0x26, 0x26, 0x51, 0x4e, 0x30, 0x31,
		0x2c, 0xbb, 0xfa, 0xc6, 0xf7, 0x30,
		0x30, 0x2c, 0xbb, 0xd8, 0xc2, 0xb7,
		0x31, 0x31, 0x39, 0x2c, 0xb2, 0xbf,
		0xbc, 0xfe, 0x30, 0x30, 0x33, 0x00,
		0x08, 0x08, 0x02, 0x14, 0x0b, 0x17,
		0x00, 0x00, 0x33, 0x33, 0x00, 0x00,
		0x1d, 0x26, 0x26, 0x51, 0x4e, 0x30,
		0x31, 0x2c, 0xbb, 0xfa, 0xc6, 0xf7,
		0x30, 0x30, 0x2c, 0xbb, 0xd8, 0xc2,
		0xb7, 0x31, 0x30, 0x32, 0x2c, 0xb2,
		0xbf, 0xbc, 0xfe, 0x30, 0x32, 0x33,
		0x00, 0x08, 0x08, 0x02, 0x14, 0x0b,
		0x17, 0x00, 0x00, 0x33, 0x33, 0x00,
		0x00, 0x1c, 0x26, 0x26, 0x51, 0x4e,
		0x30, 0x31, 0x2c, 0xbb, 0xfa, 0xc6,
		0xf7, 0x30, 0x30, 0x2c, 0xbb, 0xd8,
		0xc2, 0xb7, 0x30, 0x37, 0x2c, 0xb2,
		0xbf, 0xbc, 0xfe, 0x31, 0x32, 0x36,
		0x00, 0x08, 0x08, 0x02, 0x14, 0x0b,
		0x17, 0x00, 0x00, 0x33, 0x33, 0x00,
		0x00, 0x1c, 0x26, 0x26, 0x51, 0x4e,
		0x30, 0x31, 0x2c, 0xbb, 0xfa, 0xc6,
		0xf7, 0x30, 0x30, 0x2c, 0xbb, 0xd8,
		0xc2, 0xb7, 0x30, 0x30, 0x2c, 0xb2,
		0xbf, 0xbc, 0xfe, 0x31, 0x31, 0x38,
		0x00, 0x08, 0x08, 0x02, 0x14, 0x0b,
		0x17, 0x00, 0x00, 0x33, 0x33, 0x00,
		0x00, 0x1d, 0x26, 0x26, 0x51, 0x4e,
		0x30, 0x31, 0x2c, 0xbb, 0xfa, 0xc6,
		0xf7, 0x30, 0x30, 0x2c, 0xbb, 0xd8,
		0xc2, 0xb7, 0x31, 0x31, 0x39, 0x2c,
		0xb2, 0xbf, 0xbc, 0xfe, 0x30, 0x30,
		0x37, 0x00, 0x08, 0x08, 0x02, 0x14,
		0x0b, 0x17, 0xd8, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机消息-tengren回答
func uploadSystemctlPRTFuntengrenhuida() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40, 0x10, 0x01, 0x02, 0x03, 0x3a, 0x35, 0x05, 0x15, 0x0b, 0x17,
		0x58, 0x13, 0x97, 0x10, 0xa1, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x34, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33, 0x33, 0x01, 0x00, 0x1c,
		0x26, 0x26, 0x51, 0x4e, 0x30, 0x31, 0x2c, 0xbb, 0xfa, 0xc6, 0xf7, 0x30,
		0x31, 0x2c, 0xbb, 0xd8, 0xc2, 0xb7, 0x30, 0x35, 0x2c, 0xb2, 0xbf, 0xbc,
		0xfe, 0x30, 0x36, 0x34, 0x08, 0xc9, 0xe8, 0xb1, 0xb8, 0xbb, 0xd0, 0x18,
		0xb4, 0xf0, 0x36, 0x37, 0x0d, 0x15, 0x0b, 0x17, 0x19, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机消息-tengrenf复位
func uploadSystemctlPRTFuntengrenfuwei() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{
		0x40, 0x40, 0xc0, 0x01, 0x02, 0x03, 0x07, 0x34, 0x07, 0x15, 0x0b, 0x17,
		0x58, 0x13, 0x97, 0x10, 0xa1, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x36, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33, 0x33, 0x09, 0x00, 0x1c,
		0x26, 0x26, 0x51, 0x4e, 0x30, 0x31, 0x2c, 0xbb, 0xfa, 0xc6, 0xf7, 0x30,
		0x31, 0x2c, 0xbb, 0xd8, 0xc2, 0xb7, 0x30, 0x30, 0x2c, 0xb2, 0xbf, 0xbc,
		0xfe, 0x30, 0x30, 0x30, 0x0a, 0xbf, 0xd8, 0xd6, 0xc6, 0xc6, 0xf7, 0xb8,
		0xb4, 0xd4, 0xad, 0x15, 0x35, 0x0f, 0x15, 0x0b, 0x17, 0xff, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息  B28复位
func uploadSystemctlPRTFuncB28Reset() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{

		0x40, 0x40, 0x1d, 0x00, 0x02, 0x03, 0x26, 0x3b, 0x08, 0x18, 0x0b, 0x17,
		0x09, 0x25, 0x51, 0x80, 0x12, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x18, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33, 0x33, 0x01, 0x00, 0x00,
		0x08, 0xcf, 0xb5, 0xcd, 0xb3, 0xb8, 0xb4, 0xce, 0xbb, 0x00, 0x02, 0x11,
		0x18, 0x0b, 0x17, 0xae, 0x23, 0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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

// 上传消控主机解析卡PRT信息  B28复位
func uploadSystemctlPRTFuncB28air() {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("clientx start err, exit!", err)
		return
	}

	var uploadAnalog = []byte{

		0x40, 0x40, 0x1f, 0x00, 0x02, 0x03, 0x1d, 0x11, 0x09, 0x18, 0x0b, 0x17,
		0x09, 0x25, 0x51, 0x80, 0x12, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x2b, 0x00, 0x02, 0x09, 0x01, 0x00, 0x00, 0x33, 0x33, 0x05, 0x08, 0xbf,
		0xd5, 0xb5, 0xf7, 0xd0, 0xc2, 0xb7, 0xe7, 0x0b, 0x36, 0x31, 0x36, 0x32,
		0x20, 0xc7, 0xd0, 0xbf, 0xd5, 0xb5, 0xf7, 0x08, 0xc9, 0xe8, 0xb1, 0xb8,
		0xb9, 0xca, 0xd5, 0xcf, 0x00, 0x05, 0x10, 0x10, 0x09, 0x17, 0x1e, 0x23,
		0x23,
	}

	//	for {

	dp := firenet.NewDataPack(2048)
	msg, err := dp.Pack(uploadAnalog)
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
