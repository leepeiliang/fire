package fire

import (
	"fire/pkg/data"
	"k8s.io/klog/v2"
	"strings"
)

// 建筑消防设施部件运行状态数据结构
// 每一个位代表什么
const (
	GoodRunComponent    = iota //系统运行正常
	FireAlarmComponent         //火警
	FailComponent              //故障
	ScreenComponent            //屏蔽
	RegulatoryComponent        //监管
	StartComponent             //启动
	FeedbackComponent          //反馈
	RockStateComponent         //岩石状态
	PowerFailComponent         //电源故障
	PreWarningComponent        //预警

)

// 烟感的部件运行状态
const (
	SmokeAlarm        = iota //烟雾告警指示位(1:告警;0:正常)
	FailAlarm                //故障告警指示位(1:告警;0:正常)
	TamperAlarm              //防拆告警指示位(1:告警;0:正常)
	UnderVoltageAlarm        //欠压报警指示位(1:告警;0 正常)
)

// 防火门的部件运行状态
const (
	NowStatDoor    = 0  //当前状态指示位(1:开启;0:关闭)
	FailAlarmDoor  = 1  //异常告警指示位(1:告警;0:正常)
	ElectAlarmDoor = 7  //电量告警指示位(1:告警;0:正常)
	TypeDoor       = 15 //防火门类型指示位(1:常开;0:常闭)
)

// 消防水压&液位的部件运行状态
const (
	WaterLowLimit   = 0 //低于下限值告警指示位(1:告警;0:正常)
	WaterHighLimit  = 1 //高于上限值告警指示位(1:告警;0:正常
	ElectAlarmWater = 7 //电量告警指示位(1:告警;0:正常)
)

// 消防栓智能闷盖的部件运行状态
const (
	BlockStat             = iota //封盖状态(1:封盖开启;0:封盖关闭)
	WaterButton                  //水浸开关(60mm 口有无取水) (1:有水流;0:无水流)
	BlockChange                  //100mm 口有无变动(1:有变动;0:无变动)
	BlockTwisted                 //100mm 口有无被拧动(即 100mm 口有无取水) (1:正在被拧动;0:无变动)
	ElectricityAlarmBlock = 7    //电量告警指示位(1:告警;0:正常)
)

// 模拟量定义
const (
	EventCount       = 0x1  //1事件计数（件）电气火灾监控主机系统
	Height           = 0x2  //2高度（m）
	Temperature      = 0x3  //3LoRa 烟感网关系统，电器火灾 监控主机系统
	PressureA        = 0x4  //4消防管道压力系统
	PressureB        = 0x5  //5消防管道压力系统
	GasConcentration = 0x6  //6消防管道压力系统
	Time             = 0x7  //7时间
	Voltage          = 0x8  //8电压 电气火灾监控主机系统
	Current          = 0x9  //9电流 电气火灾监控主机系统
	Flow             = 0x0a //10流量
	AirVolume        = 0x0b //11风量
	WindSpeed        = 0x0c //12风速

	DBm               = 0x80 //128信号强度 dBm LoRa 烟感网关系统、 市政消防栓监测系统
	Power             = 0x81 //129电量 mAh  LoRa 烟感网关系统
	ResidualCurrentI1 = 0x82 //130剩余电流I1 mA 电气火灾监控主机系统
	ResidualCurrentI2 = 0x83 //131剩余电流I2 mA 电气火灾监控主机系统
	ResidualCurrentI3 = 0x84 //132剩余电流I3 mA 电气火灾监控主机系统
	ResidualCurrentI4 = 0x85 //133剩余电流I4 mA 电气火灾监控主机系统
	TemperatureT1     = 0x86 //134温度T1 mA 电气火灾监控主机系统
	FireHydraulic     = 0x87 //135消防液压，液位ADC原始值 智慧消防水系统
	BatteryVoltage    = 0x88 //136电池电压 mV 市政消防栓监测系统
	ResidualCurrentI5 = 0x89 //137剩余电流I5 mA 电气火灾监控主机系统
	ResidualCurrentI6 = 0x8a //138剩余电流I6 mA 电气火灾监控主机系统
	ResidualCurrentI7 = 0x8b //139剩余电流I7 mA 电气火灾监控主机系统
	ResidualCurrentI8 = 0x8c //140剩余电流I8 mA 电气火灾监控主机系统
	TemperatureT2     = 0x8d //141温度T2 0.1C 电气火灾监控主机系统
	TemperatureT3     = 0x8e //142温度T3 0.1C 电气火灾监控主机系统
	TemperatureT4     = 0x8f //143温度T4 0.1C 电气火灾监控主机系统
	TemperatureT5     = 0x90 //144温度T5 0.1C 电气火灾监控主机系统
	TemperatureT6     = 0x91 //145温度T6 0.1C 电气火灾监控主机系统
	TemperatureT7     = 0x92 //146温度T7 0.1C 电气火灾监控主机系统
	TemperatureT8     = 0x93 //147温度T8 0.1C 电气火灾监控主机系统
	ElectricalFireUa  = 0x94 //148Ua V 电气火灾监控主机系统。电气 火灾报警器(2 相)为 V1
	ElectricalFireUb  = 0x95 //149Ub V 电气火灾监控主机系统。电气 火灾报警器(2 相)为 V2
	ElectricalFireUc  = 0x96 //150Uc V 电气火灾监控主机系统
	ElectricalFireIa  = 0x97 //151Ia V 电气火灾监控主机系统。电气 火灾报警器(2 相)为 I1
	ElectricalFireIb  = 0x98 //152Ib V 电气火灾监控主机系统。电气 火灾报警器(2 相)为 I2
	ElectricalFireIc  = 0x99 //153Ic V 电气火灾监控主机系统
	LiquidLevel       = 0x9a //154液位 cm 消防水池液位系统、 消防水箱液位系统

	AlarmLimitIa                  = 0x9b //155Ia 告警限值       A  电气火灾监控主机系统
	WarningLimitIa                = 0x9c //156Ia 预警限值       A  电气火灾监控主机系统
	AlarmLimitIb                  = 0x9d //157Ib 告警限值       A  电气火灾监控主机系统
	WarningLimitIb                = 0x9e //158Ib 预警限值       A  电气火灾监控主机系统
	AlarmLimitIc                  = 0x9f //159Ic 告警限值       A  电气火灾监控主机系统
	WarningLimitIc                = 0xa0 //160Ic 预警限值       A  电气火灾监控主机系统
	AlarmLimitIn                  = 0xa1 //161In 告警限值      mA  电气火灾监控主机系统
	WarningLimitIn                = 0xa2 //162In 预警限值      mA  电气火灾监控主机系统
	AlarmLimitUa                  = 0xa3 //163Ua 告警限值       A  电气火灾监控主机系统
	WarningLimitUa                = 0xa4 //164Ua 预警限值       A  电气火灾监控主机系统
	AlarmLimitUb                  = 0xa5 //165Ub 告警限值       A  电气火灾监控主机系统
	WarningLimitUb                = 0xa6 //166Ub 预警限值       A  电气火灾监控主机系统
	AlarmLimitIaUc                = 0xa7 //167Uc 告警限值       A  电气火灾监控主机系统
	WarningLimitUc                = 0xa8 //168Uc 预警限值       A  电气火灾监控主机系统
	AlarmLimitT1                  = 0xa9 //169T1 告警限值    0.1C  电气火灾监控主机系统
	WarningLimitT1                = 0xaa //170T1 预警限值    0.1C  电气火灾监控主机系统
	AlarmLimitT2                  = 0xab //171T2 告警限值    0.1C  电气火灾监控主机系统
	WarningLimitT2                = 0xac //172T2 预警限值    0.1C  电气火灾监控主机系统
	AlarmLimitT3                  = 0xad //173T3 告警限值    0.1C  电气火灾监控主机系统
	WarningLimitT3                = 0xae //174T3 预警限值    0.1C  电气火灾监控主机系统
	AlarmLimitT4                  = 0xaf //175T4 告警限值    0.1C  电气火灾监控主机系统
	WarningLimitT4                = 0xb0 //176T4 预警限值    0.1C  电气火灾监控主机系统
	ResidualCurrent               = 0xb1 //177剩余电流 传感器状态    电气火灾监控主机系统
	TemperatureSensorStatus       = 0xb2 //178温度传感器状态 电气火灾监控主机系统
	ThreePhaseVoltageCurrentState = 0xb3 //179三相电压/电流状态  电气火灾监控主机系统
	FireAlarmControlSystem        = 0xf9 //249火灾报警控制  火灾报警系统
	DSFFuncControl                = 0xfa //250设备疑似误报功能 控制 LoRa 烟感网关系统
	DSFJumpCycle                  = 0xfb //251设备疑似误报跳变 周期设置 LoRa 烟感网关系统
	DSFJumpNumber                 = 0xfc //252设备疑似误报跳变 次数设置 LoRa 烟感网关系统
	DSFDisabledCycle              = 0xfd //253设备疑似误报禁用 周期设置 LoRa 烟感网关系统
	ReportToControl               = 0xfe //254设备烟感以及防拆 上报控制 LoRa 烟感网关系统
	ReportInterval                = 0xff //255上报周期 可写的模拟量

)

// 剩余电流传感器状态
// 子部件模拟量类型为剩余电流传感器状态(用 155 表示)，模拟量值为状态的实际值，从低到高位分别 表示剩余电流传感器报警状态(8 位)+剩余电流传感器故障状态，0 都表示正常，如下
const (
	AlarmRCI1 = iota // 0~7 剩余电流传感器报警状态
	AlarmRCI2
	AlarmRCI3
	AlarmRCI4
	AlarmRCI5
	AlarmRCI6
	AlarmRCI8
	FailRCI1 // 8~15 剩余电流传感器故障状态
	FailRCI2
	FailRCI3
	FailRCI4
	FailRCI5
	FailRCI6
	FailRCI7
	FailRCI8
)

// 温度传感器状态
// 子部件模拟量类型为温度传感器状态(用 156 表示)，模拟量值为状态的实际值，从低到高位分别表示
// 温度传感器报警状态(8 位)+温度传感器故障状态，0 都表示正常，如下
const (
	AlarmTST1 = iota //0~7 温度传感器报警状态(8 位
	AlarmTST2
	AlarmTST3
	AlarmTST4
	AlarmTST5
	AlarmTST6
	AlarmTST8
	FailTST1 //8~15温度传感器故障状态
	FailTST2
	FailTST3
	FailTST4
	FailTST5
	FailTST6
	FailTST7
	FailTST8
)

// 三相电压/电流状态
// 子部件模拟量类型为三相电压/电流状态(用 157 表示)，模拟量值为状态的实际值，从低到高位分别
// 表示 Ua Ub Uc Ia Ib Ic，0 都表示正常，如下
const (
	ThreePhaseUa = iota
	ThreePhaseUb
	ThreePhaseUc
	ThreePhaseIa
	ThreePhaseIb
	ThreePhaseIc
)

// 建筑消防设施系统/部件状态
// 报警类型状态数据结构
// 告警部件: 烟感 1 烟雾告警 2 故障告警 4 拆卸告警 8 欠压告警
//
//	防火门 1烟雾告警 2 故障告警 4 拆卸告警 8 欠压告警
const (
	RunStateSYS        = iota //0 系统运行正常
	FireAlarmSYS              //1 火警
	FireMalfunctionSYS        //2 故障
	OperateSYS                //3 屏蔽/操作/联动
	RegulatorySYS             //4 监管
	StartSYS                  //5 启动
	FeedbackSYS               //6 反馈/停止/恢复 //延时 允许
	DeferStateSYS             //7 延时状态
	MainPowerFailSYS          //8 主电故障
	BackupPowerFailSYS        //9 备电故障
	BusFailSYS                //10 总线故障
	ManualStateSYS            //11 手动状态
	ConfigChangesSYS          //12 配置改变
	ResetSYS                  //13 复位

)

// 用户信息传输装置运行状态
const (
	GoodRunUserComponent         = iota //系统运行正常
	FireAlarmUserComponent              //火警
	FailUserComponent                   //故障
	MainPowerFailUserComponent          //主电源故障
	DeputyPowerFailUserComponent        //备电源故障
	FailChannelComponent                //与监控中心信道故障
	FailLineComponent                   //监控连接线故障
)

// 用户信息传输装置操作信息
const (
	UserResetSYS      = iota //复位
	ErasureSYS               //消音
	UserManualState          //手动状态
	AlarmRemove              //火警消除
	SelfInspection           //自检
	OnInspectionReply        //查岗应答
	TestSYS                  //Test
)

// 建筑消防设施操作信息
const (
	OperateReset = iota
	OperateErasure
	OperateManual
	OperateEliminate
	OperateSelfInspection
	OperateConfirm
	OperateTest
)

type Property2Stat struct {
	Status [2]byte
}

func (p *Property2Stat) FireReadProperty2Stat(buffer *data.Buffer) {
	copy(p.Status[:], buffer.ReadN(2))
	return
}

// GetBit 系统状态数据为 2 字节，低字节传输在前。
func (p Property2Stat) GetBit(bit int) bool {
	if bit < 8 {
		return (p.Status[0]&(1<<bit))>>bit == 1
	}

	return (p.Status[1]&(1<<(bit-8)))>>(bit-8) == 1
}

// SetBit 系统状态数据为 2 字节，低字节传输在前。
func (p *Property2Stat) SetBit(bit int) {
	if bit > 16 || bit < 0 {
		return
	}
	if bit < 8 {
		p.Status[0] |= 1 << bit
		return
	}
	p.Status[1] |= 1 << (bit - 8)
	return
}

type PropertyStat struct {
	Status byte
}

// GetBit 系统状态数据为 2 字节，低字节传输在前。
func (ss PropertyStat) GetBit(bit int) bool {
	if bit < 8 {
		return (ss.Status&(1<<bit))>>bit == 1
	}

	return false
}

// SetBit 系统状态数据为 2 字节，低字节传输在前。
func (ss *PropertyStat) SetBit(bit int) {
	if bit > 8 || bit < 0 {
		return
	}
	if bit < 8 {
		ss.Status |= 1 << bit
		return
	}
	return
}

func (ss *PropertyStat) StringStripDefaultPropertyStat(name, place, msg string) {
	klog.Infof("通用状态转换name: name[%s]", place)
	klog.Infof("通用状态转换place: place[%s]", place)
	klog.Infof("通用状态转换msg: msg[%s]", msg)
	if name == "" && place == "" && msg == "" {
		return
	}
	name = strings.TrimSpace(place)
	place = strings.TrimSpace(place)
	msg = strings.TrimSpace(msg)

	// 如果出现四信判断出火警 什么也不考虑 直接返回火警
	if ss.GetBit(1) {
		return
	}
	// 如果出现四信判断出故障 什么也不考虑 直接返回故障
	if ss.GetBit(2) {
		return
	}
	ss.Status = 0x00
	ss.SetBit(RunStateSYS) //状态位-正常
	// 如果出现四信判断出故障 什么也不考虑 直接返回故障
	if ss.GetBit(3) {
		ss.SetBit(OperateSYS) //状态位-正常
	}

	// 如果出现火警 什么也不考虑 直接返回火警
	if strings.Contains(name, "火警") ||
		strings.Contains(place, "火警") ||
		strings.Contains(msg, "火警") ||
		strings.Contains(name, "火灾报警") ||
		strings.Contains(place, "火灾报警") ||
		strings.Contains(msg, "火灾报警") {
		if strings.Contains(name, "恢复") ||
			strings.Contains(place, "恢复") ||
			strings.Contains(msg, "恢复") ||
			strings.Contains(name, "解除") ||
			strings.Contains(place, "解除") ||
			strings.Contains(msg, "解除") ||
			strings.Contains(name, "取消") ||
			strings.Contains(place, "取消") ||
			strings.Contains(msg, "取消") {
			return
		}
		ss.SetBit(FireAlarmSYS)
		return
	}

	//返回故障
	if strings.Contains(name, "故障") ||
		strings.Contains(place, "故障") ||
		strings.Contains(msg, "故障") {
		if strings.Contains(name, "恢复") ||
			strings.Contains(place, "恢复") ||
			strings.Contains(msg, "恢复") ||
			strings.Contains(name, "解除") ||
			strings.Contains(place, "解除") ||
			strings.Contains(msg, "解除") ||
			strings.Contains(name, "取消") ||
			strings.Contains(place, "取消") ||
			strings.Contains(msg, "取消") {
			return
		}

		ss.SetBit(FireMalfunctionSYS)
		return
	}
	//
	if strings.Contains(place, "回答") ||
		strings.Contains(place, "屏蔽") ||
		strings.Contains(msg, "回答") ||
		strings.Contains(msg, "屏蔽") {
		if strings.Contains(name, "恢复") ||
			strings.Contains(place, "恢复") ||
			strings.Contains(msg, "恢复") ||
			strings.Contains(name, "解除") ||
			strings.Contains(place, "解除") ||
			strings.Contains(msg, "解除") ||
			strings.Contains(name, "取消") ||
			strings.Contains(place, "取消") ||
			strings.Contains(msg, "取消") {
			return
		}
		klog.Infof("回答状态转换msg: status[%x]", ss.Status)
		ss.SetBit(OperateSYS)
		return

	}

	if strings.Contains(name, "监管") ||
		strings.Contains(place, "监管") ||
		strings.Contains(msg, "监管") {
		if strings.Contains(name, "恢复") ||
			strings.Contains(place, "恢复") ||
			strings.Contains(msg, "恢复") ||
			strings.Contains(name, "解除") ||
			strings.Contains(place, "解除") ||
			strings.Contains(msg, "解除") ||
			strings.Contains(name, "取消") ||
			strings.Contains(place, "取消") ||
			strings.Contains(msg, "取消") {
			return
		}

		klog.Infof("监管状态转换msg: status[%x]", ss.Status)
		ss.SetBit(OperateSYS)
		ss.SetBit(RegulatorySYS)
		return

	}
	//  启动的是操作启动
	if strings.Contains(name, "启动") ||
		strings.Contains(place, "启动") ||
		strings.Contains(msg, "启动") {
		if strings.Contains(name, "恢复") ||
			strings.Contains(place, "恢复") ||
			strings.Contains(msg, "恢复") ||
			strings.Contains(name, "解除") ||
			strings.Contains(place, "解除") ||
			strings.Contains(msg, "解除") ||
			strings.Contains(name, "取消") ||
			strings.Contains(place, "取消") ||
			strings.Contains(msg, "取消") {
			return
		}

		ss.SetBit(OperateSYS)
		ss.SetBit(StartSYS)
		klog.Infof("启动状态转换msg: status[%x]", ss.Status)
	}

	// 反馈的是操作停止
	if strings.Contains(name, "恢复") ||
		strings.Contains(name, "停止") ||
		strings.Contains(name, "反馈") ||
		strings.Contains(name, "允许") ||
		strings.Contains(place, "恢复") ||
		strings.Contains(place, "停止") ||
		strings.Contains(place, "反馈") ||
		strings.Contains(place, "允许") ||
		strings.Contains(msg, "恢复") ||
		strings.Contains(msg, "停止") ||
		strings.Contains(msg, "反馈") ||
		strings.Contains(msg, "允许") {
		if strings.Contains(name, "解除") ||
			strings.Contains(place, "解除") ||
			strings.Contains(msg, "解除") ||
			strings.Contains(name, "取消") ||
			strings.Contains(place, "取消") ||
			strings.Contains(msg, "取消") {
			return
		}

		ss.Status = 0x00
		ss.SetBit(OperateSYS)
		ss.SetBit(FeedbackSYS) //状态位-正常
		klog.Infof("恢复状态转换msg: status[%x]", ss.Status)
		return

	}
	if strings.Contains(name, "自动允许") ||
		strings.Contains(place, "自动允许") ||
		strings.Contains(msg, "自动允许") {
		ss.Status = 0x00
		ss.SetBit(OperateSYS)
		ss.SetBit(FeedbackSYS) //状态位-正常
		klog.Infof("自动允许状态转换msg: status[%x]", ss.Status)
		return

	}
	if strings.Contains(name, "联动请求") ||
		strings.Contains(place, "联动请求") ||
		strings.Contains(msg, "联动请求") {
		klog.Infof("联动请求msg: status[%x]", ss.Status)
		ss.Status = 0x00
		ss.SetBit(OperateSYS)
		ss.SetBit(RunStateSYS) //
		return

	}
	klog.Infof("特殊处理状态转换-status:[%x]", ss.Status)
	return
}

// 上传消控主机解析卡 CRT 数据
const (
	DeviceName     = iota //设备名称
	DeviceDescribe        //设备描述
	DeviceMsg             //设备信息
	DeviceNum             //设备编号
)

// 上传消控主机解析卡 CRT 数据
const (
	DeviceNameInfo     = "11111111" //设备名称
	DeviceDescribeInfo = "22222222" //设备描述
	DeviceMsgInfo      = "33333333" //设备信息
	DeviceNumInfo      = "44444444" //设备编号 回路号
)
