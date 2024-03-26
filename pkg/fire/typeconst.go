package fire

var (
	StartSign = [...]byte{0x40, 0x40}
	EndSign   = [...]byte{0x23, 0x23}
	Source    = [...]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

const (
	StartSignLen        = 2
	SerialNumberLen     = 2
	ProtocolVersionLen  = 2
	TimeDateFlagLen     = 6
	SourceAddressLen    = 6
	TargetAddressLen    = 6
	AppDataLenLen       = 2
	CommandLen          = 1
	CRCLen              = 1
	EndSignLen          = 2
	FireControlTwoLen   = 15
	FireControlFirstLen = 10 //first+two len = 25

	FireControlLen   = 25 //first+two len = 25
	RemoveDataAllLen = 30
)

const (
	BigEndianness    = 1
	LittleEndianness = 0
)

// 控制单元命令字节定义表

type CommunicationProtocol uint8

const (
	Ethernet      CommunicationProtocol = 0x00 //以太网
	NBModule      CommunicationProtocol = 0x01 //NB网络
	LoraModule    CommunicationProtocol = 0x02 //Lora 模块
	LorawanModule CommunicationProtocol = 0x03 //Lorawan 模块
	GPRSModule    CommunicationProtocol = 0x04 //GPRS 模块(解析卡)
//	ControlResponse      CommunicationProtocol = 0x05～255 //用户自定义

)

// 控制单元命令字节定义表

type ControlCommand uint8

const (
	ControlReserved      ControlCommand = 0x00 //预留
	RemoteControlCommand ControlCommand = 0x01 //控制命令-时间同步、下发远程控制
	ControlSendData      ControlCommand = 0x02 //发送数据-发送火灾自动报警系统火灾报警、运行状态等信息
	ControlConfirm       ControlCommand = 0x03 //确认-对控制命令和发送信息的确认回答
	ControlRequest       ControlCommand = 0x04 //请求-查询火灾自动报警系统的火灾报警、运行状态等信息
	ControlResponse      ControlCommand = 0x05 //应答-返回查询的信息、上报远程应答
	ControlRepudiate     ControlCommand = 0x06 //拒绝，否定-对控制命令和发送信息的否认回答
	Heartbeat            ControlCommand = 0xfe //心跳
)

// 应用数据单元基本格式
//应用数据单元基本格式如图 2 所示
//-------------------------------------------
//|数据单元标识符  |信息对象    | 1字节
//|             |信息对象数目 | 1字节
//|信息对象 1     |信息体     | 根据类型不同长度不同
//|             |时间标签 1  | 6字节
//|...
//|信息对象 n     |信息体     |根据类型不同长度不同
//|              |时间标签1   |6字节
//--------------------------------------------
// 数据定义
// 数据单元标识符
// 类型标志
/*
   11～20预留(建筑消防设施信息)
   23~23 预留
   27 预留
   29～40 预留(用户传输装置信息)
   41～60 预留(用户传输装置信息)
   69～80 预留(用户传输装置信息)

   82~83 预留
   87 预留
   92～127 预留
*/

// DataBaseType 数据单元标识
type DataBaseType uint8

const (
	DataReserved               DataBaseType = 0x00 //预留
	UploadFireSystemState      DataBaseType = 0x01 //上传建筑消防设施系统状态
	UploadFireSystemRunState   DataBaseType = 0x02 //上传建筑消防设施部件运行状态
	UploadFireSystemSimulate   DataBaseType = 0x03 //上传建筑消防设施部件模拟量值
	UploadFireSystemOperate    DataBaseType = 0x04 //上传建筑消防设施操作信息
	UploadFireSystemVersion    DataBaseType = 0x05 //上传建筑消防设施软件版本
	UploadFireSystemConfig     DataBaseType = 0x06 //上传建筑消防设施系统配置情况
	UploadFireComponentConfig  DataBaseType = 0x07 //上传建筑消防设施部件配置情况
	UploadFireSystemTime       DataBaseType = 0x08 //上传建筑消防设施系统时间
	UploadFirePareCardPrintMsg DataBaseType = 0x09 //上传消控主机解析卡打印机信息
	UploadFirePareCardCRTMsg   DataBaseType = 0x0a //上传消控主机解析卡 CRT 信息
	//

	UploadUserSendDeviceRunStat DataBaseType = 0x15 //上传用户传输装置运行状态

	UploadUserSendDeviceOperate DataBaseType = 0x18 //上行 上传用户传输装置操作信息
	UploadUserSendDeviceVersion DataBaseType = 0x19 //上传用户传输装置软件版本
	UploadUserSendDeviceConfig  DataBaseType = 0x1a //上传用户传输装置配置情况

	UploadUserSendDeviceTime DataBaseType = 0x1c //上传用户传输装置系统时间

	ReadFireSystemState     DataBaseType = 0x3d //读建筑消防设施系统状态
	ReadFireSystemRunState  DataBaseType = 0x3e //读建筑消防设施部件运行状态
	ReadFireSystemSimulate  DataBaseType = 0x3f //读建筑消防设施部件模拟量值
	ReadFireSystemOperate   DataBaseType = 0x40 //读建筑消防设施操作信息
	ReadFireSystemVersion   DataBaseType = 0x41 //读建筑消防设施软件版本
	ReadFireSystemConfig    DataBaseType = 0x42 //读建筑消防设施系统配置情况
	ReadFireComponentConfig DataBaseType = 0x43 //读建筑消防设施部件配置情况
	ReadFireSystemTime      DataBaseType = 0x44 //读建筑消防设施系统时间

	ReadUserSendDeviceRunStat DataBaseType = 0x51 //读用户传输装置运行状态

	ReadUserSendDeviceOperate DataBaseType = 0x54 //上传用户传输装置操作信息记录
	ReadUserSendDeviceVersion DataBaseType = 0x55 //上传用户传输装置软件版本
	ReadUserSendDeviceConfig  DataBaseType = 0x56 //上传用户传输装置配置情况

	ReadUserSendDeviceTime    DataBaseType = 0x58 //上传用户传输装置系统时间
	InitUserSendDeviceRunStat DataBaseType = 0x59 //初始化用户传输装置
	SyncUserSendDeviceRunStat DataBaseType = 0x5a //同步用户传输装置时钟
	InspectionCommand         DataBaseType = 0x5b //查岗命令

	UploadUserSystemProduceTime DataBaseType = 0x80 //上传用户传输装置生产时间
	RemoteControl               DataBaseType = 0x81 //下发远程控制
	UploadSystemOpenStat        DataBaseType = 0x82 //上报用户信息传输装置开机时间
	UploadComponentOpenStat     DataBaseType = 0x83 //上报用户信息传输装置关机时间
	UploadBuildSystemOpen       DataBaseType = 0x84 //上报建筑消防设施系统开机信息
	UploadBuildSystemClose      DataBaseType = 0x85 //上报建筑消防设施系统关机信息
	UploadSystemStatRecover     DataBaseType = 0x86 //上传建筑消防设施系统运行状态恢复
	UploadSystemUintStatRecover DataBaseType = 0x87 //上传建筑消防设施部件运行状态恢复
	UploadUserSystemRecover     DataBaseType = 0x88 //上报用户信息传输装置运行状态恢复
	StartUpgrade                DataBaseType = 0x89 //启动升级
	UploadEnergy                DataBaseType = 0x8a //上报电能量
	Breaker                     DataBaseType = 0x8b //断路器
	MultiChannel                DataBaseType = 0x8c //上报多通道设备状态、模拟量
	PassThrough                 DataBaseType = 0x8d //解析卡透传数据
	UploadDeviceStat            DataBaseType = 0x8e //上报设备状态、模拟量
	UploadFireSystemUSStat      DataBaseType = 0xc8 //上传用户信息传输装置与监控中心线路运行状态

	UploadFireSystemUSBStat       DataBaseType = 0xc9 //上传用户信息传输装置与监控中心线路恢复状态
	UploadFireSystemUnitFireStat  DataBaseType = 0xcc //上报建筑消防设施部件火警状态-204
	UploadFireSystemLineStat      DataBaseType = 0xcd //上报建筑消防设施系联动状态-205
	UploadFireSystemUnitOtherStat DataBaseType = 0xce //上报建筑消防设施部件其他状态-206

)

// SystemType 系统类型定义表
// 2~9 预留
// 25～127预留

type SystemType uint8

const (
	General   SystemType = 0x00 //通用
	FireAlarm SystemType = 0x01 //火灾报警系统

	FireLinkController                SystemType = 0x0a //消防联动控制器
	FireTiedSystem                    SystemType = 0x0b //消防栓系统
	FireExtinguisherAutoWater         SystemType = 0x0c //自动喷水灭火器
	FireExtinguisherGas               SystemType = 0x0d //气体灭火器系统
	FireExtinguisherGasWaterPump      SystemType = 0x0e //水喷雾灭火系统（泵启动方式）
	FireExtinguisherGasWaterContainer SystemType = 0x0f //水喷雾灭火系统（压力容器启动方式）
	FireExtinguisherBubble            SystemType = 0x10 //泡沫灭火器
	FireExtinguisherPowder            SystemType = 0x11 //干粉灭火器

	PreventSmokeSystem SystemType = 0x12 //防烟排烟系统
	FireDoorSystem     SystemType = 0x13 //防火门及卷帘系统
	FireElevator       SystemType = 0x14 //消防电梯
	FireEmergencyRadio SystemType = 0x15 //消防应急广播
	FireEmergencyLight SystemType = 0x16 //消防应急照明和疏散指示系统
	FireEmergencyPower SystemType = 0x17 //消防电源
	FireEmergencyPhone SystemType = 0x18 //消防电话

	SeparateComponentSystem SystemType = 0x80 //独立部件系统（无系统，直接接部件）
	FireSystemPower         SystemType = 0x81 //电气火灾系统
	LoRaGateWay             SystemType = 0x82 //LoRa网关系统
	FireSystemPipe          SystemType = 0x83 //消防管道压力系统
	FireSystemPoolLevel     SystemType = 0x84 //消防水池液位系统
	FireSystemMunicipal     SystemType = 0x85 //市政消防栓检测系统
	FireSystemWaterTank     SystemType = 0x86 //消防水箱液位系统
	FireSystemSmoke         SystemType = 0x87 //烟感系统
	FireSystemProtocol      SystemType = 0x88 //协议系统
	DTU                     SystemType = 0x89 //DTU
	//用户自定义

)

// ComponentType 部件类型代码表
type ComponentType uint8

func (c ComponentType) ComponentType10() uint8 {
	return uint8(c)/16*10 + uint8(c)%16
}

// 系统类型定义表
// 2~9   预留
// 14～15预留
// 20预留
// 26～29预留
// 38～39预留
// 45～49预留
// FMD FireMonitoringDetector
// FMA FireMonitoringAlarm
// TFMD TemperatureFMD
// SFMD 	SmokeFMD
// CFMD   CompositeFMD
// FD    Flame Detector
//

const (
	ComponentGeneral   ComponentType = 0x00 //00通用
	ComponentFireAlarm ComponentType = 0x01 //01火灾报警系统

	ControlCombustibleGas      ComponentType = 0x0a //10可燃气体报警控制器
	CombustibleGasDetector     ComponentType = 0x0b //11点型可燃气体探测器
	OwnCombustibleGasDetector  ComponentType = 0x0c //12独立式可燃气体探测器
	LineCombustibleGasDetector ComponentType = 0x0d //13线型可燃气体探测器

	ElectricalFMA                   ComponentType = 0x10 //16电气火灾监控报警器
	ResidualCurrentElectricalFMA    ComponentType = 0x11 //17剩余电流式电气火灾监控探测器
	TemperatureMeasuringElectricFMD ComponentType = 0x12 //18测温式电气火灾监控探测器
	ElectricalFMA2                  ComponentType = 0x13 //19电气火灾监控报警器(2 相)

	DetectionCircuit      ComponentType = 0x15 //21探测回路
	FireDisplayPanel      ComponentType = 0x16 //22火灾显示盘
	ManualFireAlarmButton ComponentType = 0x17 //23手动火灾报警按钮
	HydrantButton         ComponentType = 0x18 //24消火栓按钮
	FMD                   ComponentType = 0x19 //25火灾探测器

	TFMD             ComponentType = 0x1e //30感温火灾探测器
	SpotTFMD         ComponentType = 0x1f //31点型感温火灾探测器
	SSpotTFMD        ComponentType = 0x20 //32点型感温火灾探测器(S 型)
	RSpotTFMD        ComponentType = 0x21 //33点型感温火灾探测器(R 型)
	LineTFMD         ComponentType = 0x22 //34线型感温火灾探测器
	SLineTFMD        ComponentType = 0x23 //35线型感温火灾探测器(S 型)
	RLineTFMD        ComponentType = 0x24 //36线型感温火灾探测器(R 型)
	OpticalFiberTFMD ComponentType = 0x25 //37光纤感温火灾探测器

	SFMD                  ComponentType = 0x28 //40感烟火灾探测器
	SpotSFMD              ComponentType = 0x29 //41点型离子感烟火灾探测器
	SpotPhotoelectricSFMD ComponentType = 0x2a //42点型光电感烟火灾探测器
	LineBeamSFMD          ComponentType = 0x2b //43线型光束感烟火灾探测器
	AspiratedSFMD         ComponentType = 0x2c //44吸气式感烟火灾探测器

	CFMD                    ComponentType = 0x32 //50复合式火灾探测器
	Temperature2ndSmokeCFMD ComponentType = 0x33 //51复合式感烟感温火灾探测器
	Beam2ndTemperatureCFMD  ComponentType = 0x34 //52复合式感光感温火灾探测器
	Beam2ndSmokeCFMD        ComponentType = 0x35 //53复合式感光感烟火灾探测器

	UVFD       ComponentType = 0x3d //61紫外火焰探测器
	InfraredFD ComponentType = 0x3e //62红外火焰探测器

	LightFMD ComponentType = 0x45 //69感光火灾探测器

	GasDetector         ComponentType = 0x4a //74气体探测器
	CameraFMD           ComponentType = 0x4e //78图像摄像方式火灾探测器
	VoiceFMD            ComponentType = 0x4f //79感声火灾探测器
	AcoustoOpticAlarm   ComponentType = 0x50 //80声光报警器
	GasController       ComponentType = 0x51 //81气体灭火控制器
	ElectricFireControl ComponentType = 0x52 //82消防电气控制装置
	GraphicsFireControl ComponentType = 0x53 //83消防控制室图形显示装置
	Module              ComponentType = 0x54 //84模块
	InputModule         ComponentType = 0x55 //85输入模块
	Output              ComponentType = 0x56 //86输出模块
	IOModule            ComponentType = 0x57 //87输入/输出模块
	RelayModule         ComponentType = 0x58 //88中继模块

	FirePump      ComponentType = 0x5b //91消防水泵
	FireWaterTank ComponentType = 0x5c //92消防水箱

	SprayPump      ComponentType = 0x5f //95喷淋泵
	FlowIndicator  ComponentType = 0x60 //96水流指示器
	SignalValve    ComponentType = 0x61 //97信号阀
	AlarmValve     ComponentType = 0x62 //98报警阀
	PressureSwitch ComponentType = 0x63 //99压力开关

	ValveDriver               ComponentType = 0x65 //101阀驱动装置
	FireDoor                  ComponentType = 0x66 //102防火门
	FireValve                 ComponentType = 0x67 //103防火阀
	VentilationAirConditioner ComponentType = 0x68 //104通风空调
	FoamPump                  ComponentType = 0x69 //105泡沫液泵
	PipeSolenoidValve         ComponentType = 0x6a //106管网电磁阀

	SmokeControl2ndExhaustFan ComponentType = 0x6f //111防烟排烟风机

	SmokeExhaustFireValve             ComponentType = 0x71 //113排烟防火阀
	AlwaysCloseTheAirSupply           ComponentType = 0x71 //114常闭送风口
	ExhaustPort                       ComponentType = 0x73 //115排烟口
	ElectricControlSmokeRetainingWall ComponentType = 0x74 //116电控挡烟垂壁
	FireShutterController             ComponentType = 0x75 //117防火卷帘控制器
	FireDoorMonitor                   ComponentType = 0x76 //118防火门监控器

	AlarmDeviceComponentType = 0x79 //121警报装置

	FireWaterPressure            ComponentType = 0x80 //128消防水压
	FireProtectionLevel          ComponentType = 0x81 //129消防液位
	IntelligentHoodOfFireHydrant ComponentType = 0x82 //130消防栓智能闷盖
	ParseTheCard                 ComponentType = 0x83 //131解析卡
	LORARouter                   ComponentType = 0x84 //132LORA+路由
	UserInfoTransDevice          ComponentType = 0x85 //133用户信息传输装置

)

type ChannelType uint8

const (
	AfterCurrentItem ChannelType = 0x01 //标识剩余电流，
	TemperatureItem  ChannelType = 0x02 //标识温度，
	ACurrentItem     ChannelType = 0x03 //标识A相电流，
	BCurrentItem     ChannelType = 0x04 //标识B相电流，
	CCurrentItem     ChannelType = 0x05 // 标识C相电流，
	AVoltageItem     ChannelType = 0x06 //标识A相电压,
	BVoltageItem     ChannelType = 0x07 //标识B相电压,
	CVoltageItem     ChannelType = 0x08 //标识 C 相电压。
)
