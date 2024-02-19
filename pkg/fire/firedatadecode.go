package fire

import (
	"fire/pkg/data"
	"k8s.io/klog/v2"
)

func (m *FireData) Decode(dAtA []byte) error {

	buf := data.NewBuffer(dAtA)

	m.DataBaseType = DataBaseType(buf.ReadByte())
	m.ObjectNum = buf.ReadByte()
	m.Data = make([]byte, buf.Len())
	m.Data = buf.ReadN(buf.Len())
	return nil
}

func (m *FireData) DecodeSwitchDataTypeToData(dAtA []byte) []Data {
	var back []Data
	m.Decode(dAtA)
	klog.Infof("DecodeSwitchDataTypeToData= %v", m)
	klog.Infof("DecodeSwitchDataType= 0x%x", m.DataBaseType)
	switch m.DataBaseType {
	case UploadFireSystemState: //上传建筑消防设施系统状态-01
		return m.FireBuildFacilitiesSysStatDecodeToData()
	case UploadFireSystemRunState, UploadFireSystemUnitFireStat, UploadFireSystemUnitOtherStat: //上传建筑消防设施部件运行状态-02/204/206
		return m.FireBuildFacilitiesPartRunStatDecodeToData()
	case UploadFireSystemSimulate: //上传建筑消防设施部件模拟量值-03
		return m.FireBuildAnalogDecodeToData()
	case UploadFireSystemOperate: //上传建筑消防设施操作信息-04
		return m.FireBuildOperaterDecodeToData()
	case UploadFireSystemVersion: //上传建筑消防设施软件版本-05
		return m.FireBuildVersionDecodeToData()
	case UploadFireSystemConfig: //建筑消防设施系统配置情况-06
		return m.FireBuildSystemConfigStatDecodeToData()
	case UploadFireComponentConfig: //建筑消防设施系统部件配置情况-07
		return m.FireBuildSystemComponentConfigStatDecodeToData()
	case UploadFireSystemTime, UploadBuildSystemOpen, UploadBuildSystemClose: //上传建筑消防设施系统时间-08/系统开机-132/系统关机-133
		return m.FireBuildFacilitiesTimeStatDecodeToData()
	case UploadFirePareCardPrintMsg: //上传消控主机解析卡打印机信息-09
		return m.FireParseCardAlarmStatDecodeToData()
	case UploadFirePareCardCRTMsg: //上传消控主机解析卡 CRT 信息-10
		return m.FireParseCardCRTAlarmStatDecodeToData()
	case UploadUserSendDeviceRunStat: // 用户信息传输装置运行状态数据-21
		return m.FireBuildUserConfigStatDecodeToData()
	case UploadUserSendDeviceTime, UploadUserSystemProduceTime: //上传用户信息传输装置系统时间-28/128
		return m.FireBuildFacilitiesUsersTimeStatDecodeToData()
	case UploadUserSendDeviceOperate: //上传用户传输装置操作信息记录-24
		return m.FireUserOperaterDecodeToData()
	case UploadFireSystemUSStat: //上传用户信息传输装置与监控中心线路运行状态-200
		return m.FireUserToUSStatDecodeToData()
	case UploadFireSystemUSBStat: //上传用户信息传输装置与监控中心线路恢复状态-201
		return m.FireUserToUSStatDecodeToData()
	case UploadFireSystemLineStat: //上报建筑消防设施系联动状态-205
		return m.FireBuildSystemLineStatDecodeToData()
	case UploadSystemOpenStat, UploadComponentOpenStat: //上传用户信息传输装置开机时间信息-130/上报用户信息传输装置关机时间-131
		return m.FireBuildFacilitiesUsersTimeStatDecodeToData()
	case UploadSystemStatRecover, UploadSystemUintStatRecover: //上传建筑消防设施系统状态恢复-134/上传建筑消防设施系统部件状态恢复-135
		return m.FireBuildFacilitiesPartRunStatDecodeToDataRecover()
	case UploadUserSystemRecover: //上传建筑消防设施用户传输装置运行状态恢复-136
		return m.FireBuildFacilitiesUserRunStatDecodeToDataRecover()
	case UploadEnergy: //上传电能
		return m.FireBuildElectricityDecodeToData()
	case MultiChannel: //多通道状态模拟量
		return m.FireBuildMultiChannelAnalogDecodeToData()
	case UploadDeviceStat: //状态模拟量
		return m.FireBuildDeviceAnalogDecodeToData()
	case ReadFireSystemState: //读建筑消防设施系统状态
		return m.FireBuildFacilitiesSysStatDecodeToData()
	case ReadFireSystemRunState: //读建筑消防设施部件运行状态
		return m.FireBuildFacilitiesPartRunStatDecodeToData()
	case ReadFireSystemSimulate: //读建筑消防设施部件模拟量值
		return m.FireBuildAnalogDecodeToData()
	case ReadFireSystemOperate: //读建筑消防设施操作信息
		return m.FireBuildOperaterDecodeToData()
	case ReadFireSystemVersion: //读建筑消防设施软件版本
		return m.FireBuildVersionDecodeToData()
	case ReadFireSystemConfig: //读建筑消防设施系统配置情况
		return m.FireBuildSystemComponentConfigStatDecodeToData()
	case ReadFireComponentConfig: //读建筑消防设施部件配置情况
		return m.FireBuildSystemComponentConfigStatDecodeToData()
	case ReadFireSystemTime: //读建筑消防设施系统时间
		return m.FireBuildFacilitiesTimeStatDecodeToData()
	case RemoteControl: //读建筑消防设施软件版本
		return m.FireBuildVersionDecodeToData()
	case PassThrough: //读建筑消防设施软件版本
		m.FireBuildTransparentTransmissionDecodeToData()
	default:
		klog.Errorf("%T parse type not currently supported", m.DataBaseType)
	}
	return back
}
