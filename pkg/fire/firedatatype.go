package fire

import (
	"bytes"
	"encoding/binary"
	"fire/pkg/common"
	"fire/pkg/data"
	"fire/pkg/globals"
	"fmt"
	"github.com/axgle/mahonia"
	"k8s.io/klog/v2"
	"regexp"
	"strings"
)

type FireData struct {
	DataBaseType        //类型标志
	ObjectNum    byte   //信息对象数目
	Data         []byte //信息体
}
type DataBitKey uint8

const (
	DataBaseTypeKey               DataBitKey = iota //0-数据单元标识符
	SystemTypeKey                                   //1-系统类型代码
	SystemAddressKey                                //2-系统地址
	ChannelTypeKey                                  //3-通道类型
	ChannelIDKey                                    //4-通道号
	FireHostIDKey                                   //5-消控主机编号
	ComponentTypeKey                                //6-部件类型 设备名称
	ComponentAddressKey                             //7-部件地址 设备位置
	ComponentChannelTypeKey                         //8-部件通道类型
	ComponentChannelIDKey                           //9-部件通道号 回路
	DeviceIDKey                                     //10-设备编号
	ElectricityCycleKey                             //11-电能周期
	ElectricityIDKey                                //12-电能序号
	ElectricityCurrentTypeTypeKey                   //13-电能电流类型
	PropertyIDKey                                   //14-状态位  //主版本
	AnalogTypeKey                                   //15-模拟量类型 // 次版本
)

// 特殊状态位   目前定义字符串类型告警
const (
	PropertyIDMsg = "s"
)

type DataKey string

func NewDataKey() DataKey {
	// 数据单元标识符-系统类型代码-系统地址-消控主机编号-部件类型-部件地址-通道类型-通道号-部件通道类型-部件通道号-状态位-模拟量类型
	// 特殊结构 上传消控主机解析卡PRT报警字符串  状态位的0，1，2，3代表测点告警状态，4代表告警设备告警信息
	return "#-#-#-#-#-#-#-#-#-#-#-#-#-#-#-#"
}

const sep = "-"

func (d *DataKey) SetDataBitKey(n string, dbk DataBitKey) {
	tmpSlice := strings.Split(string(*d), sep)
	tmpSlice[dbk] = fmt.Sprintf("%s", n)
	newString := strings.Join(tmpSlice, sep)
	*d = DataKey(newString)
	return
}

type Data struct {
	Data map[DataKey]*common.DataValue `json:"data"`
}

// FireBuildFacilitiesSysStat 建筑消防设施系统状态
type FireBuildFacilitiesSysStat struct {
	// 系统类型标志符为 1 字节二进制数，取值范围 0--255，系统类型定义如表 4 所示。
	SystemType
	// 系统地址为 1 字节二进制数，取值范围 0--255，由建筑消防设施设定。
	SystemAddress byte
	// 系统状态数据为 2 字节，低字节传输在前。
	Property2Stat
	//通道类型(1 字节) 选用，配合 8.1.1 类型标志 130 使用
	//01 标识剩余电流，02 标识温度，03标识A相电流，04标识B相电流，05标识C相电流，06标识A相电压,07标识B相电压,08 标识 C 相电压
	ChannelType byte
	// 通道号(1 字节) 配合 8.1.1 类型标志 131 使用
	//选用，需要上报通道号时使用，低字节传输在前。
	ChannelID byte
	//为数据包发出的时间，具体定义见 10.2.2。
	TimeLabels
}

func (m FireData) FireBuildFacilitiesSysStatDecode() []FireBuildFacilitiesSysStat {
	var back = make([]FireBuildFacilitiesSysStat, 0)

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {
		var f FireBuildFacilitiesSysStat
		f.SystemType = SystemType(buf.ReadByte())
		f.SystemAddress = buf.ReadByte()
		copy(f.Property2Stat.Status[:], buf.ReadN(2))
		if m.DataBaseType == UploadSystemStatRecover {
			f.ChannelType = buf.ReadByte()
			f.ChannelID = buf.ReadByte()
		}
		f.FireReadTime(buf)
		klog.Infof("FireBuildFacilitiesSysStatDecode Pos:%d--%+v", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}
func (m FireData) FireBuildFacilitiesSysStatDecodeToData() []Data {
	klog.Infof("FireBuildElectricityDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	back = make([]Data, 0)
	fs := m.FireBuildFacilitiesSysStatDecode()
	//	klog.Infof("FireBuildElectricityDecode :%+v", fs)
	for _, f := range fs {
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}

		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemType), SystemTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemAddress), SystemAddressKey)
		//if m.DataBaseType == UploadSystemStat {
		//	tmp.SetDataBitKey(fmt.Sprintf("%d", f.ChannelType), ChannelTypeKey)
		//	tmp.SetDataBitKey(fmt.Sprintf("%d", f.ChannelID), ChannelIDKey)
		//}
		for i := 0; i < 16; i++ {
			if i == 14 || i == 15 {
				continue
			}
			if f.Property2Stat.GetBit(i) {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     1,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			} else {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     0,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			}
			//tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
			//temp.Data[tmp] = &common.DataValue{
			//	Value:     f.Property2Stat.GetBit(i),
			//	Timestamp: common.GetTimestamp(),
			//	Metadata: common.DataMetadata{
			//		Type:      "boolean",
			//		Timestamp: common.GetTimestamp(),
			//	},
			//}
		}
		back = append(back, temp)
	}

	return back
}

// 建筑消防设施部件地址为 4 字节二进制数，建筑消防设施部件状态数据为 2 字节，低字节先传输。
type ComponentAddress struct {
	Address [4]byte
}

func (c *ComponentAddress) String() string {
	var params string
	switch globals.GetHostNameInfo() {
	case globals.TEST:
		params = fmt.Sprintf("%09d", binary.LittleEndian.Uint32(c.Address[:]))
	default:
		params = fmt.Sprintf("%09d", binary.LittleEndian.Uint32(c.Address[:]))
	}
	fmt.Println("ComponentAddress: " + params)
	return params
}
func (c *ComponentAddress) FireReadComponentAddress(buffer *data.Buffer) {
	copy(c.Address[:], buffer.ReadN(4))
	return
}

// FireBuildFacilitiesPartRunStat 建筑消防设施部件运行状态
type FireBuildFacilitiesPartRunStat struct {
	// 建筑消防设施系统类型标志、系统地址分别为 1 字节二进制数，其定义见 8.2.1.1。
	SystemType
	//建筑消防设施部件类型标志符为 1 字节二进制数，定义如 'const''部件类型代码表' 所示。
	SystemAddress byte
	// 部件类型(1 字节)
	ComponentType
	// 建筑消防设施部件地址为 4 字节二进制数，建筑消防设施部件状态数据为 2 字节，低字节先传输。
	ComponentAddress
	// 部件状态(2 字节)
	Property2Stat
	// 建筑消防设施部件说明为 31 字节的字符串，采用 GB 18030--2005 规定的编码。
	ComponentMsg [31]byte
	// 部件通道类型(1 字节) 选用，配合 8.1.1 类型标志 131 使用
	ComponentChannelType byte
	// 部件通道号(1 字节) 选用，配合 8.1.1 类型标志 131 使用
	ComponentChannelID byte
	TimeLabels
}

func (m FireData) FireBuildFacilitiesPartRunStatDecode() []FireBuildFacilitiesPartRunStat {
	var back = make([]FireBuildFacilitiesPartRunStat, 0)

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {
		var f FireBuildFacilitiesPartRunStat
		f.SystemType = SystemType(buf.ReadByte())
		f.SystemAddress = buf.ReadByte()
		f.ComponentType = ComponentType(buf.ReadInt8())
		f.ComponentAddress.FireReadComponentAddress(buf)
		copy(f.Property2Stat.Status[:], buf.ReadN(2))
		copy(f.ComponentMsg[:], buf.ReadN(31))
		f.FireReadTime(buf)
		klog.Infof("FireBuildFacilitiesPartRunStatDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}
func (m FireData) FireBuildFacilitiesPartRunStatDecodeToData() []Data {
	klog.Infof("FireBuildFacilitiesPartRunStatDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	back = make([]Data, 0)
	fs := m.FireBuildFacilitiesPartRunStatDecode()
	klog.Infof("FireBuildElectricityDecode :%+v", fs)
	for _, f := range fs {
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemType), SystemTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemAddress), SystemAddressKey)
		tmp.SetDataBitKey(fmt.Sprintf("%x", f.ComponentType), ComponentTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%s", f.ComponentAddress.String()), ComponentAddressKey)
		klog.Infof("FireBuildComponentMsg :%s", f.ComponentMsg)

		if f.ComponentType == SFMD {
			for i := 0; i < 4; i++ {
				//tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				//temp.Data[tmp] = &common.DataValue{
				//	Value:     f.Property2Stat.GetBit(1 << i),
				//	Timestamp: common.GetTimestamp(),
				//	Metadata: common.DataMetadata{
				//		Type:      "boolean",
				//		Timestamp: f.TimeLabels.FireToTimeUnixNano(),
				//	},
				//}

				if f.Property2Stat.GetBit(1 << i) {
					tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
					temp.Data[tmp] = &common.DataValue{
						Value:     1,
						Timestamp: common.GetTimestamp(),
						Metadata: common.DataMetadata{
							Type:      "boolean",
							Timestamp: f.TimeLabels.FireToTimeUnixNano(),
						},
					}
				} else {
					tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
					temp.Data[tmp] = &common.DataValue{
						Value:     0,
						Timestamp: common.GetTimestamp(),
						Metadata: common.DataMetadata{
							Type:      "boolean",
							Timestamp: f.TimeLabels.FireToTimeUnixNano(),
						},
					}
				}
			}
			//klog.Infof("FireBuildElectricityDecode Len:%+v", len(temp.Data))
			goto NEXT
		}
		if f.ComponentType == FireDoor {
			bits := []int{0, 1, 7, 15}
			for _, num := range bits {
				//tmp.SetDataBitKey(fmt.Sprintf("%d", num), PropertyIDKey)
				//temp.Data[tmp] = &common.DataValue{
				//	Value:     f.Property2Stat.GetBit(num),
				//	Timestamp: common.GetTimestamp(),
				//	Metadata: common.DataMetadata{
				//		Type:      "boolean",
				//		Timestamp: common.GetTimestamp(),
				//	},
				//}

				if f.Property2Stat.GetBit(num) {
					tmp.SetDataBitKey(fmt.Sprintf("%d", num), PropertyIDKey)
					temp.Data[tmp] = &common.DataValue{
						Value:     1,
						Timestamp: common.GetTimestamp(),
						Metadata: common.DataMetadata{
							Type:      "boolean",
							Timestamp: f.TimeLabels.FireToTimeUnixNano(),
						},
					}
				} else {
					tmp.SetDataBitKey(fmt.Sprintf("%d", num), PropertyIDKey)
					temp.Data[tmp] = &common.DataValue{
						Value:     0,
						Timestamp: common.GetTimestamp(),
						Metadata: common.DataMetadata{
							Type:      "boolean",
							Timestamp: f.TimeLabels.FireToTimeUnixNano(),
						},
					}
				}
			}
			goto NEXT
		}
		if f.ComponentType == FireWaterPressure {
			bits := []int{0, 1, 7}
			for _, num := range bits {
				//tmp.SetDataBitKey(fmt.Sprintf("%d", num), PropertyIDKey)
				//temp.Data[tmp] = &common.DataValue{
				//	Value:     f.Property2Stat.GetBit(num),
				//	Timestamp: common.GetTimestamp(),
				//	Metadata: common.DataMetadata{
				//		Type:      "boolean",
				//		Timestamp: common.GetTimestamp(),
				//	},
				//}
				if f.Property2Stat.GetBit(num) {
					tmp.SetDataBitKey(fmt.Sprintf("%d", num), PropertyIDKey)
					temp.Data[tmp] = &common.DataValue{
						Value:     1,
						Timestamp: common.GetTimestamp(),
						Metadata: common.DataMetadata{
							Type:      "boolean",
							Timestamp: f.TimeLabels.FireToTimeUnixNano(),
						},
					}
				} else {
					tmp.SetDataBitKey(fmt.Sprintf("%d", num), PropertyIDKey)
					temp.Data[tmp] = &common.DataValue{
						Value:     0,
						Timestamp: common.GetTimestamp(),
						Metadata: common.DataMetadata{
							Type:      "boolean",
							Timestamp: f.TimeLabels.FireToTimeUnixNano(),
						},
					}
				}
			}
			goto NEXT
		}
		if f.ComponentType == FireProtectionLevel {
			bits := []int{0, 1, 7}
			for _, num := range bits {
				//tmp.SetDataBitKey(fmt.Sprintf("%d", num), PropertyIDKey)
				//temp.Data[tmp] = &common.DataValue{
				//	Value:     f.Property2Stat.GetBit(num),
				//	Timestamp: common.GetTimestamp(),
				//	Metadata: common.DataMetadata{
				//		Type:      "boolean",
				//		Timestamp: common.GetTimestamp(),
				//	},
				//}

				if f.Property2Stat.GetBit(num) {
					tmp.SetDataBitKey(fmt.Sprintf("%d", num), PropertyIDKey)
					temp.Data[tmp] = &common.DataValue{
						Value:     1,
						Timestamp: common.GetTimestamp(),
						Metadata: common.DataMetadata{
							Type:      "boolean",
							Timestamp: f.TimeLabels.FireToTimeUnixNano(),
						},
					}
				} else {
					tmp.SetDataBitKey(fmt.Sprintf("%d", num), PropertyIDKey)
					temp.Data[tmp] = &common.DataValue{
						Value:     0,
						Timestamp: common.GetTimestamp(),
						Metadata: common.DataMetadata{
							Type:      "boolean",
							Timestamp: f.TimeLabels.FireToTimeUnixNano(),
						},
					}
				}
			}
			goto NEXT
		}
		if f.ComponentType == IntelligentHoodOfFireHydrant {
			bits := []int{0, 1, 2, 3, 7}
			for _, num := range bits {
				//tmp.SetDataBitKey(fmt.Sprintf("%d", num), PropertyIDKey)
				//temp.Data[tmp] = &common.DataValue{
				//	Value:     f.Property2Stat.GetBit(num),
				//	Timestamp: common.GetTimestamp(),
				//	Metadata: common.DataMetadata{
				//		Type:      "boolean",
				//		Timestamp: common.GetTimestamp(),
				//	},
				//}

				if f.Property2Stat.GetBit(num) {
					tmp.SetDataBitKey(fmt.Sprintf("%d", num), PropertyIDKey)
					temp.Data[tmp] = &common.DataValue{
						Value:     1,
						Timestamp: common.GetTimestamp(),
						Metadata: common.DataMetadata{
							Type:      "boolean",
							Timestamp: f.TimeLabels.FireToTimeUnixNano(),
						},
					}
				} else {
					tmp.SetDataBitKey(fmt.Sprintf("%d", num), PropertyIDKey)
					temp.Data[tmp] = &common.DataValue{
						Value:     0,
						Timestamp: common.GetTimestamp(),
						Metadata: common.DataMetadata{
							Type:      "boolean",
							Timestamp: f.TimeLabels.FireToTimeUnixNano(),
						},
					}
				}
			}
			goto NEXT
		}
		for i := 0; i < 10; i++ {
			//tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
			//temp.Data[tmp] = &common.DataValue{
			//	Value:     f.Property2Stat.GetBit(i),
			//	Timestamp: common.GetTimestamp(),
			//	Metadata: common.DataMetadata{
			//		Type:      "boolean",
			//		Timestamp: common.GetTimestamp(),
			//	},
			//}

			if f.Property2Stat.GetBit(i) {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     1,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			} else {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     0,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			}
		}
	NEXT:
		back = append(back, temp)
	}

	return back
}

// FireBuildFacilitiesPartRunStatRecover 建筑消防设施部件运行状态恢复
type FireBuildFacilitiesPartRunStatRecover struct {
	// 建筑消防设施系统类型标志、系统地址分别为 1 字节二进制数，其定义见 8.2.1.1。
	SystemType
	//建筑消防设施部件类型标志符为 1 字节二进制数，定义如 'const''部件类型代码表' 所示。
	SystemAddress byte
	// 部件类型(1 字节)
	ComponentType
	// 建筑消防设施部件地址为 4 字节二进制数，建筑消防设施部件状态数据为 2 字节，低字节先传输。
	ComponentAddress
	// 部件状态(2 字节)
	Property2Stat
	TimeLabels
}

func (m FireData) FireBuildFacilitiesRunStatRecoverDecode() []FireBuildFacilitiesPartRunStatRecover {
	var back = make([]FireBuildFacilitiesPartRunStatRecover, 0)

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {
		var f FireBuildFacilitiesPartRunStatRecover
		f.SystemType = SystemType(buf.ReadByte())
		f.SystemAddress = buf.ReadByte()
		if m.DataBaseType == UploadSystemUintStatRecover {
			f.ComponentType = ComponentType(buf.ReadByte())
			f.ComponentAddress.FireReadComponentAddress(buf)
		}
		copy(f.Property2Stat.Status[:], buf.ReadN(2))
		f.FireReadTime(buf)
		klog.Infof("FireBuildFacilitiesRunStatRecoverDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}
func (m FireData) FireBuildFacilitiesPartRunStatDecodeToDataRecover() []Data {
	klog.Infof("FireBuildFacilitiesRunStatRecoverDecode:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	back = make([]Data, 0)
	fs := m.FireBuildFacilitiesRunStatRecoverDecode()
	klog.Infof("FireBuildElectricityDecode :%+v", fs)
	for _, f := range fs {
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemType), SystemTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemAddress), SystemAddressKey)
		if m.DataBaseType == UploadSystemUintStatRecover {
			//			tmp.SetDataBitKey(fmt.Sprintf("%d", f.ComponentType), ComponentTypeKey)
			tmp.SetDataBitKey(fmt.Sprintf("%s", f.ComponentAddress.String()), ComponentAddressKey)
		}

		for i := 0; i < 10; i++ {
			//tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
			//temp.Data[tmp] = &common.DataValue{
			//	Value:     f.Property2Stat.GetBit(i),
			//	Timestamp: common.GetTimestamp(),
			//	Metadata: common.DataMetadata{
			//		Type:      "boolean",
			//		Timestamp: common.GetTimestamp(),
			//	},
			//}
			if f.Property2Stat.GetBit(i) {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     1,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			} else {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     0,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			}
		}
		back = append(back, temp)
	}

	return back
}

// FireBuildFacilitiesUserRunStatRecover 上传建筑消防设施用户传输装置运行状态恢复
type FireBuildFacilitiesUserRunStatRecover struct {
	PropertyStat //状态类型 0/正常 1/火警 2/故障 3/主电故障 4/备电故障 5/监控中心通信信道故障 6/接线故障 7/预留
	TimeLabels
}

func (m FireData) FireBuildFacilitiesUserRunStatRecoverDecode() []FireBuildFacilitiesUserRunStatRecover {
	var back = make([]FireBuildFacilitiesUserRunStatRecover, 0)

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {
		var f FireBuildFacilitiesUserRunStatRecover

		f.PropertyStat.Status = buf.ReadByte()
		f.FireReadTime(buf)
		klog.Infof("FireBuildFacilitiesRunStatRecoverDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}
func (m FireData) FireBuildFacilitiesUserRunStatDecodeToDataRecover() []Data {
	klog.Infof("FireBuildFacilitiesRunStatRecoverDecode:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	back = make([]Data, 0)
	fs := m.FireBuildFacilitiesUserRunStatRecoverDecode()
	klog.Infof("FireBuildElectricityDecode :%+v", fs)
	for _, f := range fs {
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)

		for i := 0; i < 8; i++ {
			//tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
			//temp.Data[tmp] = &common.DataValue{
			//	Value:     f.PropertyStat.GetBit(i),
			//	Timestamp: common.GetTimestamp(),
			//	Metadata: common.DataMetadata{
			//		Type:      "boolean",
			//		Timestamp: common.GetTimestamp(),
			//	},
			//}
			if f.PropertyStat.GetBit(i) {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     1,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			} else {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     0,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			}
		}
		back = append(back, temp)
	}

	return back
}

// FireBuildAnalog 建筑消防部件模拟量值
type FireBuildAnalog struct {
	// 系统类型标志(1 字节)系统类型标志、系统地址、部件类型、部件地址的定义同 8.2.1.2。
	SystemType
	// 系统地址(1 字节)
	SystemAddress byte
	// 部件类型(1 字节)
	ComponentType
	// 部件地址(4 字节)
	ComponentAddress
	// 模拟量类型为 1 字节二进制数，取值范围 0~255
	AnalogType byte
	// 模拟量值为 2 字节有符号整型数，取值范围为-32768~+32767，低字节传输在前。
	// 模拟量类型和模拟量值的具体定义见表
	Analog
	TimeLabels
}
type Analog struct {
	Analog [2]byte
}

func (a Analog) WriteUint16() uint16 {
	return binary.LittleEndian.Uint16(a.Analog[:])
}
func (a Analog) WriteInt16() int16 {
	return int16(binary.LittleEndian.Uint16(a.Analog[:]))
}

func (m FireData) FireBuildAnalogDecode() []FireBuildAnalog {
	var back = make([]FireBuildAnalog, 0)

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {
		var f FireBuildAnalog
		f.SystemType = SystemType(buf.ReadByte())
		f.SystemAddress = buf.ReadByte()
		f.ComponentType = ComponentType(buf.ReadByte())
		f.ComponentAddress.FireReadComponentAddress(buf)
		f.AnalogType = buf.ReadByte()
		copy(f.Analog.Analog[:], buf.ReadN(2))
		f.FireReadTime(buf)
		klog.Infof("FireBuildFacilitiesPartRunStatDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}
func (m FireData) FireBuildAnalogDecodeToData() []Data {
	klog.Infof("FireBuildFacilitiesPartRunStatDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	back = make([]Data, 0)
	fs := m.FireBuildAnalogDecode()
	//	klog.Infof("FireBuildElectricityDecode :%+v", fs)
	for _, f := range fs {
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemType), SystemTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemAddress), SystemAddressKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.ComponentType), ComponentTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%s", f.ComponentAddress.String()), ComponentAddressKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.AnalogType), AnalogTypeKey)
		temp.Data = map[DataKey]*common.DataValue{}
		if f.AnalogType == ResidualCurrent {
			var p2 Property2Stat
			p2.Status = f.Analog.Analog
			for i := 0; i < 16; i++ {
				//tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				//temp.Data[tmp] = &common.DataValue{
				//	Value:     p2.GetBit(i),
				//	Timestamp: common.GetTimestamp(),
				//	Metadata: common.DataMetadata{
				//		Type:      "boolean",
				//		Timestamp: f.TimeLabels.FireToTimeUnixNano(),
				//	},
				//}
				if p2.GetBit(i) {
					tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
					temp.Data[tmp] = &common.DataValue{
						Value:     1,
						Timestamp: common.GetTimestamp(),
						Metadata: common.DataMetadata{
							Type:      "boolean",
							Timestamp: f.TimeLabels.FireToTimeUnixNano(),
						},
					}
				} else {
					tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
					temp.Data[tmp] = &common.DataValue{
						Value:     0,
						Timestamp: common.GetTimestamp(),
						Metadata: common.DataMetadata{
							Type:      "boolean",
							Timestamp: f.TimeLabels.FireToTimeUnixNano(),
						},
					}
				}
			}
		}
		if f.AnalogType == TemperatureSensorStatus {
			var p2 Property2Stat
			p2.Status = f.Analog.Analog
			for i := 0; i < 16; i++ {
				//tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				//temp.Data[tmp] = &common.DataValue{
				//	Value:     p2.GetBit(i),
				//	Timestamp: common.GetTimestamp(),
				//	Metadata: common.DataMetadata{
				//		Type:      "boolean",
				//		Timestamp: f.TimeLabels.FireToTimeUnixNano(),
				//	},
				//}
				if p2.GetBit(i) {
					tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
					temp.Data[tmp] = &common.DataValue{
						Value:     1,
						Timestamp: common.GetTimestamp(),
						Metadata: common.DataMetadata{
							Type:      "boolean",
							Timestamp: f.TimeLabels.FireToTimeUnixNano(),
						},
					}
				} else {
					tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
					temp.Data[tmp] = &common.DataValue{
						Value:     0,
						Timestamp: common.GetTimestamp(),
						Metadata: common.DataMetadata{
							Type:      "boolean",
							Timestamp: f.TimeLabels.FireToTimeUnixNano(),
						},
					}
				}
			}
		}
		if f.AnalogType == ThreePhaseVoltageCurrentState {
			var p2 Property2Stat
			p2.Status = f.Analog.Analog
			for i := 0; i < 6; i++ {
				//tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				//temp.Data[tmp] = &common.DataValue{
				//	Value:     p2.GetBit(i),
				//	Timestamp: common.GetTimestamp(),
				//	Metadata: common.DataMetadata{
				//		Type:      "boolean",
				//		Timestamp: f.TimeLabels.FireToTimeUnixNano(),
				//	},
				//}
				if p2.GetBit(i) {
					tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
					temp.Data[tmp] = &common.DataValue{
						Value:     1,
						Timestamp: common.GetTimestamp(),
						Metadata: common.DataMetadata{
							Type:      "boolean",
							Timestamp: f.TimeLabels.FireToTimeUnixNano(),
						},
					}
				} else {
					tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
					temp.Data[tmp] = &common.DataValue{
						Value:     0,
						Timestamp: common.GetTimestamp(),
						Metadata: common.DataMetadata{
							Type:      "boolean",
							Timestamp: f.TimeLabels.FireToTimeUnixNano(),
						},
					}
				}
			}
		}
		if f.AnalogType == FireAlarmControlSystem {
			temp.Data[tmp] = &common.DataValue{
				Value:     f.Analog,
				Timestamp: common.GetTimestamp(),
				Metadata: common.DataMetadata{
					Type:      "boolean",
					Timestamp: f.TimeLabels.FireToTimeUnixNano(),
				},
			}

		}
		if f.AnalogType == DSFFuncControl {
			temp.Data[tmp] = &common.DataValue{
				Value:     f.Analog,
				Timestamp: common.GetTimestamp(),
				Metadata: common.DataMetadata{
					Type:      "boolean",
					Timestamp: f.TimeLabels.FireToTimeUnixNano(),
				},
			}
		}
		if f.AnalogType == ReportToControl {
			temp.Data[tmp] = &common.DataValue{
				Value:     f.Analog,
				Timestamp: common.GetTimestamp(),
				Metadata: common.DataMetadata{
					Type:      "boolean",
					Timestamp: f.TimeLabels.FireToTimeUnixNano(),
				},
			}
		}
		if f.AnalogType == EventCount || f.AnalogType == Height ||
			f.AnalogType == Temperature || f.AnalogType == PressureA ||
			f.AnalogType == PressureB || f.AnalogType == GasConcentration ||
			f.AnalogType == Time || f.AnalogType == Voltage ||
			f.AnalogType == Current || f.AnalogType == Flow ||
			f.AnalogType == AirVolume || f.AnalogType == WindSpeed ||
			f.AnalogType == DBm {
			temp.Data[tmp] = &common.DataValue{
				Value:     f.Analog.WriteInt16(),
				Timestamp: common.GetTimestamp(),
				Metadata: common.DataMetadata{
					Type:      "int",
					Timestamp: f.TimeLabels.FireToTimeUnixNano(),
				},
			}
		}
		if f.AnalogType >= Power && f.AnalogType <= ResidualCurrent {
			temp.Data[tmp] = &common.DataValue{
				Value:     f.Analog.WriteUint16(),
				Timestamp: common.GetTimestamp(),
				Metadata: common.DataMetadata{
					Type:      "int",
					Timestamp: f.TimeLabels.FireToTimeUnixNano(),
				},
			}
		}
		back = append(back, temp)
	}

	return back
}

// FireBuildFacilitiesTimeStat  上传建筑消防设施系统时间
type FireBuildFacilitiesTimeStat struct {
	// 系统类型标志(1 字节)系统类型标志、系统地址、部件类型、部件地址的定义同 8.2.1.2。
	SystemType
	//建筑消防设施部件类型标志符为 1 字节二进制数，定义如 'const''部件类型代码表' 所示。
	SystemAddress byte
	TimeLabels
}

func (m FireData) FireBuildFacilitiesTimeStatDecode() []FireBuildFacilitiesTimeStat {
	var back = make([]FireBuildFacilitiesTimeStat, 0)

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {
		var f FireBuildFacilitiesTimeStat
		f.SystemType = SystemType(buf.ReadByte())
		f.SystemAddress = buf.ReadByte()
		f.FireReadTime(buf)
		klog.Infof("FireBuildFacilitiesPartRunStatDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}
func (m FireData) FireBuildFacilitiesTimeStatDecodeToData() []Data {
	klog.Infof("FireBuildFacilitiesPartRunStatDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	back = make([]Data, 0)
	fs := m.FireBuildFacilitiesTimeStatDecode()
	//	klog.Infof("FireBuildElectricityDecode :%+v", fs)
	for _, f := range fs {
		var temp Data
		klog.Infof("FireBuildFacilitiesTimeStatDecodeToData TimeLabels :%+v", f.TimeLabels)
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemType), SystemTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemAddress), SystemAddressKey)
		temp.Data[tmp] = &common.DataValue{
			Value:     f.TimeLabels.FireToTime(),
			Timestamp: common.GetTimestamp(),
			Metadata: common.DataMetadata{
				Type:      "string",
				Timestamp: common.GetTimestamp(),
			},
		}
		back = append(back, temp)
	}

	return back
}

// FireBuildFacilitiesUsersTimeStat 上传用户传输装置系统时间
type FireBuildFacilitiesUsersTimeStat struct {
	// 系统类型标志(1 字节)系统类型标志、系统地址、部件类型、部件地址的定义同 8.2.1.2。
	//操作员ID
	UserID byte
	TimeLabels
}

func (m FireData) FireBuildFacilitiesUsersTimeStatDecode() []FireBuildFacilitiesUsersTimeStat {
	var back = make([]FireBuildFacilitiesUsersTimeStat, 0)

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {
		var f FireBuildFacilitiesUsersTimeStat
		if m.DataBaseType == UploadSystemOpenStat || m.DataBaseType == UploadComponentOpenStat {
			f.UserID = buf.ReadByte()
		}
		f.FireReadTime(buf)

		klog.Infof("FireBuildFacilitiesPartRunStatDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}
func (m FireData) FireBuildFacilitiesUsersTimeStatDecodeToData() []Data {
	klog.Infof("FireBuildFacilitiesUsersTimeStatDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	back = make([]Data, 0)
	fs := m.FireBuildFacilitiesUsersTimeStatDecode()
	klog.Infof("FireBuildFacilitiesUsersTimeStatDecode :%+v", fs)
	for _, f := range fs {
		var temp Data
		klog.Infof("FireBuildFacilitiesUsersTimeStatDecodeToData TimeLabels :%+v", f.TimeLabels)
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		temp.Data[tmp] = &common.DataValue{
			Value:     f.TimeLabels.FireToTime(),
			Timestamp: common.GetTimestamp(),
			Metadata: common.DataMetadata{
				Type:      "string",
				Timestamp: common.GetTimestamp(),
			},
		}
		back = append(back, temp)
	}

	return back
}

// FireBuildOperater 建筑消防设施操作信息
type FireBuildOperater struct {
	// 系统类型标志(1 字节)系统类型标志、系统地址、部件类型、部件地址的定义同 8.2.1.2。
	SystemType
	// 系统地址(1 字节)
	SystemAddress byte
	//操作标志
	PropertyStat
	//操作员ID
	UserID byte
	TimeLabels
}

func (m FireData) FireBuildOperaterDecode() []FireBuildOperater {
	var back = make([]FireBuildOperater, 0)

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {
		var f FireBuildOperater
		f.SystemType = SystemType(buf.ReadByte())
		f.SystemAddress = buf.ReadByte()
		f.PropertyStat.Status = buf.ReadByte()
		f.UserID = buf.ReadByte()

		f.FireReadTime(buf)
		klog.Infof("FireBuildFacilitiesPartRunStatDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}

func (m FireData) FireBuildOperaterDecodeToData() []Data {
	klog.Infof("FireBuildOperaterDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	back = make([]Data, 0)
	fs := m.FireBuildOperaterDecode()
	for _, f := range fs {
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemType), SystemTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemAddress), SystemAddressKey)
		for i := 0; i < 7; i++ {

			//tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
			//temp.Data[tmp] = &common.DataValue{
			//	Value:     f.PropertyStat.GetBit(i),
			//	Timestamp: common.GetTimestamp(),
			//	Metadata: common.DataMetadata{
			//		Type:      "boolean",
			//		Timestamp: f.TimeLabels.FireToTimeUnixNano(),
			//	},
			//}
			if f.PropertyStat.GetBit(i) {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     1,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			} else {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     0,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			}
		}
		back = append(back, temp)
	}

	return back
}

// FireUserOperater 建筑消防设施操作信息
type FireUserOperater struct {

	//操作标志
	PropertyStat
	//操作员ID
	UserID byte
	TimeLabels
}

func (m FireData) FireUserOperaterDecode() []FireUserOperater {
	var back = make([]FireUserOperater, 0)

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {
		var f FireUserOperater

		f.PropertyStat.Status = buf.ReadByte()
		f.UserID = buf.ReadByte()

		f.FireReadTime(buf)
		klog.Infof("FireBuildFacilitiesPartRunStatDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}

// FireUserOperaterDecodeToData 上传用户传输装置操作信息记录-24
func (m FireData) FireUserOperaterDecodeToData() []Data {
	klog.Infof("FireBuildOperaterDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	back = make([]Data, 0)
	fs := m.FireUserOperaterDecode()
	//	klog.Infof("FireUserOperaterDecode :%+v", fs)
	for _, f := range fs {
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		for i := 0; i < 7; i++ {
			//tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
			//temp.Data[tmp] = &common.DataValue{
			//	Value:     f.PropertyStat.GetBit(i),
			//	Timestamp: common.GetTimestamp(),
			//	Metadata: common.DataMetadata{
			//		Type:      "boolean",
			//		Timestamp: f.TimeLabels.FireToTimeUnixNano(),
			//	},
			//}
			if f.PropertyStat.GetBit(i) {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     1,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			} else {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     0,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			}
		}
		back = append(back, temp)
	}

	return back
}

// FireBuildVersion 建筑消防设施的软件版本数据结构
type FireBuildVersion struct {
	// 系统类型标志(1 字节)系统类型标志、系统地址、部件类型、部件地址的定义同 8.2.1.2。
	SystemType
	// 系统地址(1 字节)
	SystemAddress byte
	//主版本
	MainVersion byte
	//副版本
	DeputyVersion byte
	TimeLabels
}

func (m FireData) FireBuildVersionDecode() []FireBuildVersion {
	var back = make([]FireBuildVersion, 0)

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {
		var f FireBuildVersion
		f.SystemType = SystemType(buf.ReadByte())
		f.SystemAddress = buf.ReadByte()
		f.MainVersion = buf.ReadByte()
		f.DeputyVersion = buf.ReadByte()

		//f.FireReadTime(buf)
		klog.Infof("FireBuildFacilitiesPartRunStatDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}
func (m FireData) FireBuildVersionDecodeToData() []Data {
	klog.Infof("FireBuildVersionDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	back = make([]Data, 0)
	fs := m.FireBuildVersionDecode()
	klog.Infof("FireBuildElectricityDecode len:%+v", len(fs))
	for _, f := range fs {
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemType), SystemTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemAddress), SystemAddressKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", 1), PropertyIDKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", 1), AnalogTypeKey)
		temp.Data[tmp] = &common.DataValue{
			Value:     fmt.Sprintf("%02d%02d", f.MainVersion, f.DeputyVersion),
			Timestamp: common.GetTimestamp(),
			Metadata: common.DataMetadata{
				Type:      "string",
				Timestamp: f.TimeLabels.FireToTimeUnixNano(),
			},
		}
		back = append(back, temp)
	}

	return back
}

// FireBuildSystemConfigStat 建筑消防设施系统配置情况
type FireBuildSystemConfigStat struct {
	// 系统类型标志(1 字节)系统类型标志、系统地址、部件类型、部件地址的定义同 8.2.1.2。
	SystemType
	// 系统地址(1 字节)
	SystemAddress byte
	DescribeLen   byte
	Describe      []byte
	TimeLabels
}

func (m FireData) FireBuildSystemConfigStatDecode() []FireBuildSystemConfigStat {
	var back = make([]FireBuildSystemConfigStat, 0)

	buf := data.NewBuffer(m.Data)
	klog.Infof("FireBuildElectricityDecode :%+v", buf.Pos())
	for i := 0; i < int(m.ObjectNum); i++ {
		var f FireBuildSystemConfigStat
		f.SystemType = SystemType(buf.ReadByte())
		f.SystemAddress = buf.ReadByte()
		f.DescribeLen = buf.ReadByte()
		copy(f.Describe[:], buf.ReadN(int(f.DescribeLen)))
		f.FireReadTime(buf)
		klog.Infof("FireBuildSystemConfigStatDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}
func (m FireData) FireBuildSystemConfigStatDecodeToData() []Data {
	klog.Infof("FireBuildSystemConfigStatDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	back = make([]Data, 0)
	fs := m.FireBuildSystemConfigStatDecode()
	//	klog.Infof("FireBuildElectricityDecode :%+v", fs)
	for _, f := range fs {
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemType), SystemTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemAddress), SystemAddressKey)

		tmp.SetDataBitKey(fmt.Sprintf("%d", 0), AnalogTypeKey)
		temp.Data[tmp] = &common.DataValue{
			Value:     string(f.Describe),
			Timestamp: common.GetTimestamp(),
			Metadata: common.DataMetadata{
				Type:      "string",
				Timestamp: f.TimeLabels.FireToTimeUnixNano(),
			},
		}
		back = append(back, temp)
	}

	return back
}

// FireBuildSystemComponentConfigStat  建筑消防设施系统部件配置情况
type FireBuildSystemComponentConfigStat struct {
	// 系统类型标志(1 字节)系统类型标志、系统地址、部件类型、部件地址的定义同 8.2.1.2。
	SystemType
	// 系统地址(1 字节)
	SystemAddress byte
	// 部件类型(1 字节)
	ComponentType
	// 部件地址(4 字节)
	ComponentAddress
	Describe [31]byte
	TimeLabels
}

func (m FireData) FireBuildSystemComponentConfigStatDecode() []FireBuildSystemComponentConfigStat {
	var back = make([]FireBuildSystemComponentConfigStat, 0)

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {
		var f FireBuildSystemComponentConfigStat
		f.SystemType = SystemType(buf.ReadByte())
		f.SystemAddress = buf.ReadByte()
		f.ComponentType = ComponentType(buf.ReadByte())
		f.ComponentAddress.FireReadComponentAddress(buf)
		copy(f.Describe[:], buf.ReadN(31))
		f.FireReadTime(buf)
		klog.Infof("FireBuildUserConfigStatDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}
func (m FireData) FireBuildSystemComponentConfigStatDecodeToData() []Data {
	klog.Infof("FireBuildSystemComponentConfigStatDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	back = make([]Data, 0)
	fs := m.FireBuildSystemComponentConfigStatDecode()
	//	klog.Infof("FireBuildElectricityDecode :%+v", fs)
	for _, f := range fs {
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemType), SystemTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemAddress), SystemAddressKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.ComponentType), ComponentTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%x", f.ComponentAddress.String()), ComponentAddressKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", 0), AnalogTypeKey)
		temp.Data[tmp] = &common.DataValue{
			Value:     string(f.Describe[:]),
			Timestamp: common.GetTimestamp(),
			Metadata: common.DataMetadata{
				Type:      "string",
				Timestamp: f.TimeLabels.FireToTimeUnixNano(),
			},
		}

		back = append(back, temp)

	}
	return back
}

// FireBuildUserConfigStat 	// 用户信息传输装置运行状态数据-21
type FireBuildUserConfigStat struct {
	PropertyStat
	TimeLabels
}

func (m FireData) FireBuildUserConfigStatDecode() []FireBuildUserConfigStat {
	var back = make([]FireBuildUserConfigStat, 0)

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {
		var f FireBuildUserConfigStat
		f.PropertyStat.Status = buf.ReadByte()
		f.FireReadTime(buf)
		klog.Infof("FireBuildUserConfigStatDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}

// 用户信息传输装置运行状态数据-21
func (m FireData) FireBuildUserConfigStatDecodeToData() []Data {
	klog.Infof("FireBuildSystemConfigStatDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var out []Data
	out = make([]Data, 0)
	var back Data
	back.Data = map[DataKey]*common.DataValue{}
	fs := m.FireBuildUserConfigStatDecode()
	//	klog.Infof("FireBuildElectricityDecode :%+v", fs)
	for _, f := range fs {
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)

		for i := 0; i < 7; i++ {
			//tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
			//temp.Data[tmp] = &common.DataValue{
			//	Value:     f.PropertyStat.GetBit(i),
			//	Timestamp: common.GetTimestamp(),
			//	Metadata: common.DataMetadata{
			//		Type:      "boolean",
			//		Timestamp: f.TimeLabels.FireToTimeUnixNano(),
			//	},
			//}
			if f.PropertyStat.GetBit(i) {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     1,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			} else {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     0,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			}
		}
		for key, value := range temp.Data {
			back.Data[key] = value
		}
	}
	out = append(out, back)
	return out
}

// FireUserToUSStat 	// 上传用户信息传输装置与监控中心线路运行状态-200
type FireUserToUSStat struct {

	//// 系统类型标志(1 字节)系统类型标志、系统地址、部件类型、部件地址的定义同 8.2.1.2。
	//SystemType
	//// 系统地址(1 字节)
	//SystemAddress byte
	PropertyStat
	TimeLabels
}

func (m FireData) FireUserToUSDecode() []FireUserToUSStat {
	var back = make([]FireUserToUSStat, 0)

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {
		var f FireUserToUSStat
		//f.SystemType = SystemType(buf.ReadByte())
		//f.SystemAddress = buf.ReadByte()
		//copy(f.Property2Stat.Status[:], buf.ReadN(2))
		f.PropertyStat.Status = buf.ReadByte()
		f.FireReadTime(buf)
		klog.Infof("FireBuildUserConfigStatDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}

// 上报建筑消防设施系联动状态-200/201
func (m FireData) FireUserToUSStatDecodeToData() []Data {
	klog.Infof("FireBuildSystemLineStatDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var out []Data
	out = make([]Data, 0)
	var back Data
	back.Data = map[DataKey]*common.DataValue{}
	fs := m.FireUserToUSDecode()
	//	klog.Infof("FireBuildSystemLineStatDecode :%+v", fs)
	for _, f := range fs {
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		//tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemType), SystemTypeKey)
		//tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemAddress), SystemAddressKey)

		for i := 0; i < 8; i++ {
			//tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
			//temp.Data[tmp] = &common.DataValue{
			//	Value:     f.PropertyStat.GetBit(i),
			//	Timestamp: common.GetTimestamp(),
			//	Metadata: common.DataMetadata{
			//		Type:      "boolean",
			//		Timestamp: f.TimeLabels.FireToTimeUnixNano(),
			//	},
			//}
			if f.PropertyStat.GetBit(i) {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     1,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			} else {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     0,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			}
		}
		for key, value := range temp.Data {
			back.Data[key] = value
		}
	}
	out = append(out, back)
	return out
}

// FireBuildSystemLineStat 	// 上报建筑消防设施系联动状态-205
type FireBuildSystemLineStat struct {

	// 系统类型标志(1 字节)系统类型标志、系统地址、部件类型、部件地址的定义同 8.2.1.2。
	SystemType
	// 系统地址(1 字节)
	SystemAddress byte
	Property2Stat
	TimeLabels
}

func (m FireData) FireBuildSystemLineStatDecode() []FireBuildSystemLineStat {
	var back = make([]FireBuildSystemLineStat, 0)

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {
		var f FireBuildSystemLineStat
		f.SystemType = SystemType(buf.ReadByte())
		f.SystemAddress = buf.ReadByte()
		copy(f.Property2Stat.Status[:], buf.ReadN(2))
		f.FireReadTime(buf)
		klog.Infof("FireBuildUserConfigStatDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}

// 上报建筑消防设施系联动状态-205
func (m FireData) FireBuildSystemLineStatDecodeToData() []Data {
	klog.Infof("FireBuildSystemLineStatDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var out []Data
	out = make([]Data, 0)
	var back Data
	back.Data = map[DataKey]*common.DataValue{}
	fs := m.FireBuildSystemLineStatDecode()
	//	klog.Infof("FireBuildSystemLineStatDecode :%+v", fs)
	for _, f := range fs {
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemType), SystemTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemAddress), SystemAddressKey)

		for i := 0; i < 16; i++ {
			//tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
			//temp.Data[tmp] = &common.DataValue{
			//	Value:     f.Property2Stat.GetBit(i),
			//	Timestamp: common.GetTimestamp(),
			//	Metadata: common.DataMetadata{
			//		Type:      "boolean",
			//		Timestamp: f.TimeLabels.FireToTimeUnixNano(),
			//	},
			//}
			if f.Property2Stat.GetBit(i) {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     1,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			} else {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     0,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			}
		}
		for key, value := range temp.Data {
			back.Data[key] = value
		}
	}
	out = append(out, back)
	return out
}

// FireParseCardPRTAlarmStat 上传消控主机解析卡PRT信息 上传消控主机解析卡报警字符串
type FireParseCardPRTAlarmStat struct {
	// 通道号(1 字节)
	ChannelID   byte
	FireHostNum byte
	Reserved    [2]byte
	PropertyStat
	Name  CardAlarmMsg
	Place CardAlarmMsg
	Msg   CardAlarmMsg
	TimeLabels
}
type CardAlarmMsg struct {
	MsgLen byte
	Msg    string
}

func (m FireData) FireParseCardAlarmStatDecode() []FireParseCardPRTAlarmStat {
	var back = make([]FireParseCardPRTAlarmStat, 0)

	buf := data.NewBuffer(m.Data)
	klog.Infof("FireBuildElectricityDecode-ObjectNum:[%d]", m.ObjectNum)
	klog.Infof("FireBuildElectricityDecode-datalen:[%d]", len(m.Data))
	enc := mahonia.NewDecoder("gbk")
	for i := 0; i < int(m.ObjectNum); i++ {
		var f FireParseCardPRTAlarmStat
		f.ChannelID = buf.ReadByte()
		f.FireHostNum = buf.ReadByte()
		copy(f.Reserved[:], buf.ReadN(2))

		f.PropertyStat.Status = buf.ReadByte()
		f.Name.MsgLen = buf.ReadByte()
		f.Name.Msg = strings.Replace(enc.ConvertString(string(buf.ReadN(int(f.Name.MsgLen)))), " ", "", -1)
		f.Place.MsgLen = buf.ReadByte()
		f.Place.Msg = strings.Replace(enc.ConvertString(string(buf.ReadN(int(f.Place.MsgLen)))), " ", "", -1)
		f.Msg.MsgLen = buf.ReadByte()
		f.Msg.Msg = strings.Replace(enc.ConvertString(string(buf.ReadN(int(f.Msg.MsgLen)))), " ", "", -1)

		// ConvertString converts a string from d's encoding to UTF-8.
		//fmt.Println(enc.ConvertString("hello,世界"))

		//		klog.Infof("中文乱码:%s", "火警:13191")
		klog.Infof("FireBuildUserConfigStatDecode status:[%x]", f.PropertyStat.Status)
		klog.Infof("FireBuildUserConfigStatDecode Name:%d--%s", f.Name.MsgLen, f.Name.Msg)
		klog.Infof("FireBuildUserConfigStatDecode Place:%d--%s", f.Place.MsgLen, f.Place.Msg)
		klog.Infof("FireBuildUserConfigStatDecode Msg:%d--%s", f.Msg.MsgLen, f.Msg.Msg)
		f.FireReadTime(buf)
		klog.Infof("FireBuildUserConfigStatDecode Pos:%d--%+v", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}
func (m FireData) FireParseCardAlarmStatDecodeToData() []Data {
	klog.Infof("FireParseCardAlarmStatDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	back = make([]Data, 0)
	fs := m.FireParseCardAlarmStatDecode()
	klog.Infof("FireParseCardAlarmStatDecode :%+v", fs)
	for _, f := range fs {
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.ChannelID), ChannelIDKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.FireHostNum), FireHostIDKey)
		//取出设备信息里的回路和设备编码
		var params string
		switch globals.GetHostNameInfo() {
		case globals.M6v3floors12HostName:
			params = StringStrip(f.Name.Msg)
		case globals.SoftwareparkHostName:
			params = StringStrip(f.Name.Msg)
		case globals.JiyunHostName:
			params = StringStripJiyun(f.Place.Msg, f.Msg.Msg)
		case globals.TongJiHostName:
			params = StringStrip(f.Msg.Msg)
		case globals.BoXingHostName:
			params = StringStripBoxxing(f.Msg.Msg)
		case globals.XIASHAHostName:
			tag := StringStrip(f.Place.Msg)
			params = StringStrip(f.Name.Msg)
			if strings.HasPrefix(tag, "7") {
				params = "7" + StringStrip(f.Name.Msg)
			}
			if strings.HasPrefix(tag, "9") {
				params = "9" + StringStrip(f.Name.Msg)
			}
			f.PropertyStat.StringStripDefaultPropertyStat(f.Name.Msg, f.Place.Msg, f.Msg.Msg)
		case globals.BJM3HostName:
			params = StringStrip(f.Msg.Msg)
		case globals.TengrenHostName:
			if strings.Contains(f.Msg.Msg, "消防主机通讯故障") {
				params = "99999"
			}
			if strings.Contains(f.Msg.Msg, "消防主机通讯恢复") {
				params = "99998"
			}
			if params == "" {
				params = StringStripTengren(f.Place.Msg)
			}
			if strings.Contains(f.Msg.Msg, "控制器复原") {
				params = "0000000"
			}
			klog.Infof("腾仁清洗结果: %s", params)
			f.PropertyStat.StringStripDefaultPropertyStat(f.Name.Msg, f.Place.Msg, f.Msg.Msg)
		case globals.SHWAIGAOQIAOHostName:
			params = StringStripWaigaoqiao(f.Place.Msg)
			klog.Infof("外高桥清洗结果: %s", params)
			f.PropertyStat.StringStripDefaultPropertyStat(f.Name.Msg, f.Place.Msg, f.Msg.Msg)
		case globals.BJM6:
			params = StringStripDefault(f.Name.Msg, f.Place.Msg, f.Msg.Msg)
			klog.Infof("m6a清洗结果: %s", params)
			f.PropertyStat.StringStripDefaultPropertyStat(f.Name.Msg, f.Place.Msg, f.Msg.Msg)
		case globals.BEIJINGB28:
			params = StringStripbtobaone(f.Place.Msg, f.Msg.Msg)
			klog.Infof("b28清洗结果: %s", params)
			f.PropertyStat.StringStripDefaultPropertyStat(f.Name.Msg, f.Place.Msg, f.Msg.Msg)
		case globals.HEDAN:
			params = StringStriphedan(f.Place.Msg, f.Msg.Msg)
			klog.Infof("hedan清洗结果: %s", params)
			f.PropertyStat.StringStripDefaultPropertyStat(f.Name.Msg, f.Place.Msg, f.Msg.Msg)
		case globals.TongZhouB, globals.SongJiang, globals.WULAN4B1, globals.WULAN3A1:
			params = StringStripBluebird(f.Name.Msg, f.Place.Msg, f.Msg.Msg)
			klog.Infof("青鸟通用清洗结果: %s", params)
			f.PropertyStat.StringStripDefaultPropertyStat(f.Name.Msg, f.Place.Msg, f.Msg.Msg)
		case globals.SANHE4, globals.SANHE6, globals.SANHE7:
			params = StringStripSanhe(f.Name.Msg)
			klog.Infof("三河海湾打印机清洗结果: %s", params)
			f.PropertyStat.StringStripDefaultPropertyStat(f.Name.Msg, f.Place.Msg, f.Msg.Msg)
		default:
			params = StringStripDefault(f.Name.Msg, f.Place.Msg, f.Msg.Msg)
			klog.Infof("测试清洗结果: %s", params)
			f.PropertyStat.StringStripDefaultPropertyStat(f.Name.Msg, f.Place.Msg, f.Msg.Msg)
		}
		if params != "" {
			tmp.SetDataBitKey(fmt.Sprintf("%s", params), DeviceIDKey)
		}
		//		tmp.SetDataBitKey(fmt.Sprintf("%s", f.Place.Msg), ComponentAddressKey)
		for i := 0; i < 8; i++ {
			if f.PropertyStat.GetBit(i) == true {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     1,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			} else {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     0,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			}
		}
		if f.Msg.MsgLen > 0 && string(f.Msg.Msg) != "" {
			tmp.SetDataBitKey(fmt.Sprintf("%s", "s"), PropertyIDKey)
			//tmp.SetDataBitKey(fmt.Sprintf("%s", PropertyIDMsg), AnalogTypeKey)
			temp.Data[tmp] = &common.DataValue{
				Value:     f.Msg.Msg,
				Timestamp: common.GetTimestamp(),
				Metadata: common.DataMetadata{
					Type:      "string",
					Timestamp: f.TimeLabels.FireToTimeUnixNano(),
				},
			}
		}
		back = append(back, temp)
	}
	return back
}

type DeviceEventType uint8

const (
	DisEnabled DeviceEventType = iota
	Strings
	Numerical
)

// 事件字段 1 表示类型(设备名称)
// 事件字段 1 相应应用数据
// 事件字段 2 表示类型(设备描述)
// 事件字段 2 相应应用数据
// 事件字段 3 表示类型(设备信息)
// 事件字段 3 相应应用数据
// 事件字段 4 表示类型(设备编号)
// 事件字段 4 相应应用数据
// 事件字段 5 表示类型(其他)
// 事件字段 5 相应应用数据
// ...
// 事件字段 n 表示类型
// 事件字段 n 相应应用数据
//
//	如果DeviceType =1 存在一个字节的数据长度 根据长度去匹配数据
type EventData struct {
	DeviceEventType
	DataInfo string //事件字段表示字符串
	DataLen  byte   //事件数据长度
	Data     string //事件字段表示字符串
	DataNum  uint32 //事件字段表示数值
}

// 上传消控主机解析卡CRT信息
type FireParseCardCRTAlarmStat struct {
	ChannelID    byte //通道号
	FireHostNum  byte //消控主机编号
	Reserved     byte //预留
	EventNum     byte //2-10
	PropertyStat      //报警类型
	DeviceInfos  []EventData
	TimeLabels
}

func (m FireData) FireParseCardCRTAlarmStatDecode() []FireParseCardCRTAlarmStat {
	var back []FireParseCardCRTAlarmStat

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	klog.Infof("FireBuildElectricityDecode-ObjectNum:[%d]", m.ObjectNum)
	for i := 0; i < int(m.ObjectNum); i++ {
		var f FireParseCardCRTAlarmStat
		f.ChannelID = buf.ReadByte()
		f.FireHostNum = buf.ReadByte()
		f.Reserved = buf.ReadByte()
		f.EventNum = buf.ReadByte()
		f.PropertyStat.Status = buf.ReadByte()
		klog.Infof("类型个数：[%d]报警类型:%x", f.EventNum, f.PropertyStat.Status)
		f.DeviceInfos = make([]EventData, 0)
		for j := 0; j < int(f.EventNum); j++ {
			var tmp EventData
			tmp.DeviceEventType = DeviceEventType(buf.ReadByte())
			klog.Infof("表示类型:%x[%d]", tmp.DeviceEventType, j)
			//			if tmp.DeviceEventType == DisEnabled {
			//				continue
			//			}
			switch j {
			case 0:
				tmp.DataInfo = DeviceNameInfo
			case 1:
				tmp.DataInfo = DeviceDescribeInfo
			case 2:
				tmp.DataInfo = DeviceMsgInfo
			case 3:
				tmp.DataInfo = DeviceNumInfo
			case 4:
				tmp.DataInfo = "其他"
			default:
				tmp.DataInfo = "其他"
			}
			if tmp.DeviceEventType == Strings {
				tmp.DataLen = buf.ReadByte()
				klog.Infof("表示类型{字符串}:%x", tmp.DataLen)
				if tmp.DataLen > 0 {
					tmp.Data = string(buf.ReadN(int(tmp.DataLen)))
					klog.Infof("%s{字符串}:%s", tmp.DataInfo, tmp.Data)
				}
			}
			if tmp.DeviceEventType == Numerical {
				tmp.DataLen = 4
				tmp.DataNum = buf.ReadUint32()
				klog.Infof("%s{数值}:%d", tmp.DataInfo, tmp.DataNum)
			}
			f.DeviceInfos = append(f.DeviceInfos, tmp)
		}

		f.FireReadTime(buf)
		klog.Infof("FireParseCardCRTAlarmStatDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}
	klog.Infof("FireParseCardCRTAlarmStatDecode-back:[%d]", len(back))
	return back
}
func (m FireData) FireParseCardCRTAlarmStatDecodeToData() []Data {
	klog.Infof("FireParseCardCRTAlarmStatDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	fs := m.FireParseCardCRTAlarmStatDecode()
	klog.Infof("FireParseCardCRTAlarmStatDecode-len[%d]", len(fs))
	for _, f := range fs {
		klog.Infof("FireParseCardCRTAlarmStatDecodeToData[%+v]", f)
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.FireHostNum), FireHostIDKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.ChannelID), ChannelIDKey)

		for _, one := range f.DeviceInfos {
			klog.Infof("FireParseCardCRTAlarmStatDecodeToData-DeviceInfos[%d][%+v]", len(f.DeviceInfos), f.DeviceInfos)
			if globals.TengrenHostName == globals.GetHostNameInfo() {
				if one.DeviceEventType == Strings && one.DataLen > 0 {
					fmt.Sprintf("%x", f.DeviceInfos[DeviceName].Data)

					// 设备名称

					// 设备描述

					// 设备信息
					if one.DataInfo == DeviceMsgInfo {
						space := strings.TrimSpace(one.Data)
						infos := strings.SplitN(space, " ", 2)
						// todo 腾仁打印机接crt  部件编号，需要确认都是这样的格式
						if len(infos) > 0 && isNumeric(infos[0]) {
							tmp.SetDataBitKey(fmt.Sprintf("%s", infos[0]), DeviceIDKey)
						}

					}

				}

				if one.DeviceEventType == Numerical {
					if one.DataInfo == DeviceNumInfo {

						// 腾仁打印机接crt  部件回路号
						tmp.SetDataBitKey(fmt.Sprintf("%d", one.DataNum), ComponentChannelIDKey)
					}

				}
			}
			if globals.BoXingHostName == globals.GetHostNameInfo() {
				if one.DeviceEventType == Strings && one.DataLen > 0 {
					fmt.Sprintf("%x", f.DeviceInfos[DeviceName].Data)

					// 设备名称

					// 设备描述

					// 设备信息
					if one.DataInfo == DeviceDescribeInfo {
						space := strings.TrimSpace(one.Data)
						infos := strings.SplitN(space, ",", 4)
						// todo 博兴消防接crt  回路设备编号，&&QN01,机器01,回路16,部件101
						if len(infos) > 0 {
							var oneTmp, twoTmp, threeTmp string
							for k, v := range infos {
								klog.Infof("设备信息解析[%d]:%+v", k, v)
								jiqi := []byte{0xbb, 0xfa, 0xc6, 0xf7} //GB2312 机器
								if bytes.HasPrefix([]byte(v), jiqi) {
									oneTmp = fmt.Sprintf("%s", StringStrip(v))
									//tmp.SetDataBitKey(fmt.Sprintf("%s", StringStrip(v)), ComponentAddressKey)
								}
								huilu := []byte{0xbb, 0xd8, 0xc2, 0xb7} //GB2312 回路
								if bytes.HasPrefix([]byte(v), huilu) {
									twoTmp = fmt.Sprintf("%s", StringStrip(v))
									//tmp.SetDataBitKey(fmt.Sprintf("%s", StringStrip(v)), ComponentChannelIDKey)
								}
								bujian := []byte{0xb2, 0xbf, 0xbc, 0xfe} //GB2312 部件
								if bytes.HasPrefix([]byte(v), bujian) {
									threeTmp = fmt.Sprintf("%s", StringStrip(v))
									//tmp.SetDataBitKey(fmt.Sprintf("%s", StringStrip(v)), DeviceIDKey)
								}
							}
							tmp.SetDataBitKey(fmt.Sprintf("%s%s%s", oneTmp, twoTmp, threeTmp), DeviceIDKey)
						}

					}

				}

			}

		}
		for i := 0; i < 4; i++ {
			if f.PropertyStat.GetBit(i) {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     1,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: common.GetTimestamp(),
					},
				}

			} else {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     0,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: common.GetTimestamp(),
					},
				}
			}

		}
		back = append(back, temp)

	}
	klog.Infof("FireParseCardCRTAlarmStatDecodeToData-back[%d]", len(back))
	return back
}
func isNumeric(input string) bool {
	match, _ := regexp.MatchString("^[0-9]+$", input)
	return match
}

type ElectricityType uint8

const (
	ReservedElectricityType ElectricityType = iota //预留
	AElectricity                                   //A项电能
	BElectricity                                   //B项电能
	CElectricity                                   //C项电能
	AllElectricity                                 //All项电能
)

type ElectricityCurrentType uint8

const (
	Reserved              ElectricityCurrentType = iota
	AElectricityCurrent                          //A项电流
	BElectricityCurrent                          //A项电流
	CElectricityCurrent                          //A项电流
	AllElectricityCurrent                        //all项电流
)
const (
	SendElectricityType = iota
	SendElectricityCurrentType
)

// FireBuildElectricity 建筑消防上报电能量
type FireBuildElectricity struct {
	// 系统类型标志(1 字节)系统类型标志、系统地址、部件类型、部件地址的定义同 8.2.1.2。
	SystemType
	// 部件类型(1 字节)
	ComponentType
	// 部件地址(4 字节)
	ComponentAddress
	//电能周期
	ElectricityCycle byte
	// 电能类型
	ElectricityType
	// 电能ID
	ElectricityID byte
	//电能值
	ElectricityAnalog uint16
	// 电流类型
	ElectricityCurrentType
	//电流值
	ElectricityCurrenAnalog uint16
	TimeLabels
}

func (m FireData) FireBuildElectricityDecode() []FireBuildElectricity {
	var back = make([]FireBuildElectricity, 0)

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {

		var f FireBuildElectricity
		f.SystemType = SystemType(buf.ReadByte())
		f.ComponentType = ComponentType(buf.ReadByte())
		copy(f.ComponentAddress.Address[:], buf.ReadN(4))
		f.ElectricityCycle = buf.ReadByte()

		f.ElectricityType = ElectricityType(buf.ReadByte())
		f.ElectricityID = buf.ReadByte()
		f.ElectricityAnalog = uint16(buf.ReadInt16())
		f.ElectricityCurrentType = ElectricityCurrentType(buf.ReadByte())
		f.ElectricityCurrenAnalog = uint16(buf.ReadInt16())
		f.FireReadTime(buf)
		klog.Infof("FireBuildElectricityDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}
func (m FireData) FireBuildElectricityDecodeToData() []Data {
	klog.Infof("FireBuildElectricityDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	back = make([]Data, 0)
	fs := m.FireBuildElectricityDecode()
	//	klog.Infof("FireBuildElectricityDecode :%+v", fs)
	for _, f := range fs {
		var (
			temp    Data
			tempkey string
		)
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemType), SystemTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.ComponentType), ComponentTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%s", f.ComponentAddress.String()), ComponentAddressKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.ElectricityCycle), ElectricityCycleKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.ElectricityID), ElectricityIDKey)

		tmp.SetDataBitKey(fmt.Sprintf("%d", f.ElectricityType), ElectricityCurrentTypeTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", SendElectricityType), AnalogTypeKey) //"0"标识电能值

		tempkey = string(tmp[:36])
		temp.Data[tmp] = &common.DataValue{
			Value:     f.ElectricityAnalog,
			Timestamp: common.GetTimestamp(),
			Metadata: common.DataMetadata{
				Type: "int",

				Timestamp: f.TimeLabels.FireToTimeUnixNano(),
			},
		}
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.ElectricityCurrentType), ElectricityCurrentTypeTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", SendElectricityCurrentType), AnalogTypeKey) //"1"标识电流值

		temp.Data[tmp] = &common.DataValue{
			Value:     f.ElectricityCurrenAnalog,
			Timestamp: common.GetTimestamp(),
			Metadata: common.DataMetadata{
				Type:      "int",
				Timestamp: f.TimeLabels.FireToTimeUnixNano(),
			},
		}

		for i, v := range back {

			for key, _ := range v.Data {
				//				klog.Infof("FireBuildElectricityDecodeToData %s===%s", tempkey, string(key))
				if strings.HasPrefix(string(key), tempkey) {

					for addkey, addval := range temp.Data {
						klog.Infof("FireBuildElectricityDecodeToData ADD dev[%d] temkey[%s]===Property key[%s]", i, tempkey, string(addkey))
						back[i].Data[addkey] = addval
					}
					goto NEXT
				}
			}

		}
		back = append(back, temp)
	NEXT:
	}
	klog.Infof("FireBuildElectricityDecodeToData Len:%d map[] Len[%d]", len(back), len(back[0].Data))
	return back
}

type MultiAnalog struct {
	// 模拟量类型为 1 字节二进制数，取值范围 0~255
	AnalogType byte
	// 模拟量值为 2 字节有符号整型数，取值范围为-32768~+32767，低字节传输在前。
	// 模拟量类型和模拟量值的具体定义见表
	Analog int16
}

// FireBuildMultiChannelAnalog 上报多通道设备状态、模拟量
type FireBuildMultiChannelAnalog struct {
	// 系统类型标志(1 字节)系统类型标志、系统地址、部件类型、部件地址的定义同 8.2.1.2。
	SystemType
	// 系统地址(1 字节)
	SystemAddress byte
	// 部件类型(1 字节)
	ComponentType
	// 部件地址(4 字节)
	ComponentAddress
	// 通道号(1 字节) 配合 8.1.1 类型标志 131 使用
	//选用，需要上报通道号时使用，低字节传输在前。
	ChannelID byte
	// 系统状态数据为 2 字节，低字节传输在前。
	Property2Stat
	// 模拟量类型为 1 字节二进制数，取值范围 0~255
	AnalogNum   uint8
	MultiAnalog []MultiAnalog
	TimeLabels
}

func (m FireData) FireBuildMultiChannelAnalogDecode() []FireBuildMultiChannelAnalog {
	var back = make([]FireBuildMultiChannelAnalog, 0)

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {

		var f FireBuildMultiChannelAnalog

		f.SystemType = SystemType(buf.ReadByte())
		f.SystemAddress = buf.ReadByte()
		f.ComponentType = ComponentType(buf.ReadByte())
		f.ComponentAddress.FireReadComponentAddress(buf)
		f.ChannelID = buf.ReadByte()
		copy(f.Property2Stat.Status[:], buf.ReadN(2))
		f.AnalogNum = buf.ReadByte()
		f.MultiAnalog = make([]MultiAnalog, 0)
		for i := 0; i < int(f.AnalogNum); i++ {
			var tmp MultiAnalog
			tmp.AnalogType = buf.ReadByte()
			tmp.Analog = buf.ReadInt16()
		}

		f.FireReadTime(buf)
		klog.Infof("FireBuildElectricityDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}

func (m FireData) FireBuildMultiChannelAnalogDecodeToData() []Data {
	klog.Infof("FireBuildElectricityDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	back = make([]Data, 0)
	fs := m.FireBuildMultiChannelAnalogDecode()
	//	klog.Infof("FireBuildElectricityDecode :%+v", fs)
	for _, f := range fs {
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemType), SystemTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemAddress), SystemAddressKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.ComponentType), ComponentTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%s", f.ComponentAddress.String()), ComponentAddressKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.ChannelID), ChannelIDKey)

		if f.ComponentType == SFMD {
			for i := 0; i < 4; i++ {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value: f.Property2Stat.GetBit(1 << i),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			}
		}
		if f.ComponentType == FireDoor {
			bits := []int{0, 1, 7, 15}
			for i, num := range bits {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     f.Property2Stat.GetBit(num),
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: common.GetTimestamp(),
					},
				}
			}
		}
		if f.ComponentType == FireWaterPressure {
			bits := []int{0, 1, 7}
			for i, num := range bits {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     f.Property2Stat.GetBit(num),
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: common.GetTimestamp(),
					},
				}
			}
		}
		if f.ComponentType == FireProtectionLevel {
			bits := []int{0, 1, 7}
			for i, num := range bits {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     f.Property2Stat.GetBit(num),
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: common.GetTimestamp(),
					},
				}
			}
		}
		if f.ComponentType == IntelligentHoodOfFireHydrant {
			bits := []int{0, 1, 2, 3, 7}
			for i, num := range bits {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     f.Property2Stat.GetBit(num),
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: common.GetTimestamp(),
					},
				}
			}

		}
		for i := 0; i < 10; i++ {
			tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
			temp.Data[tmp] = &common.DataValue{
				Value:     f.Property2Stat.GetBit(i),
				Timestamp: common.GetTimestamp(),
				Metadata: common.DataMetadata{
					Type:      "boolean",
					Timestamp: common.GetTimestamp(),
				},
			}
		}

		for i := 0; i < int(f.AnalogNum); i++ {

			tmp.SetDataBitKey(fmt.Sprintf("%d", f.MultiAnalog[i].AnalogType), AnalogTypeKey)
			temp.Data[tmp] = &common.DataValue{
				Value:     f.MultiAnalog[i].Analog,
				Timestamp: common.GetTimestamp(),
				Metadata: common.DataMetadata{
					Type:      "int",
					Timestamp: f.TimeLabels.FireToTimeUnixNano(),
				},
			}
		}
		back = append(back, temp)
	}

	return back
}

// FireBuildDeviceAnalog 上报设备状态、模拟量
type FireBuildDeviceAnalog struct {
	// 系统类型标志(1 字节)系统类型标志、系统地址、部件类型、部件地址的定义同 8.2.1.2。
	SystemType
	// 系统地址(1 字节)
	SystemAddress byte
	// 部件类型(1 字节)
	ComponentType
	// 部件地址(4 字节)
	ComponentAddress

	// 系统状态数据为 2 字节，低字节传输在前。
	Property2Stat
	// 模拟量类型为 1 字节二进制数，取值范围 0~255
	AnalogNum   uint8
	MultiAnalog []MultiAnalog
	TimeLabels
}

func (m FireData) FireBuildDeviceAnalogDecode() []FireBuildDeviceAnalog {
	var back = make([]FireBuildDeviceAnalog, 0)

	buf := data.NewBuffer(m.Data)
	//	klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {

		var f FireBuildDeviceAnalog
		f.SystemType = SystemType(buf.ReadByte())
		f.SystemAddress = buf.ReadByte()
		f.ComponentType = ComponentType(buf.ReadByte())
		f.ComponentAddress.FireReadComponentAddress(buf)
		copy(f.Property2Stat.Status[:], buf.ReadN(2))
		f.AnalogNum = buf.ReadByte()
		f.MultiAnalog = make([]MultiAnalog, 0)
		for i := 0; i < int(f.AnalogNum); i++ {
			var tmp MultiAnalog
			tmp.AnalogType = buf.ReadByte()
			tmp.Analog = buf.ReadInt16()
		}

		f.FireReadTime(buf)
		klog.Infof("FireBuildElectricityDecode Pos:%d--%x", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}
func (m FireData) FireBuildDeviceAnalogDecodeToData() []Data {
	klog.Infof("FireBuildElectricityDecodeToData:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	back = make([]Data, 0)
	fs := m.FireBuildDeviceAnalogDecode()
	//	klog.Infof("FireBuildElectricityDecode :%+v", fs)
	for _, f := range fs {
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemType), SystemTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemAddress), SystemAddressKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.ComponentType), ComponentTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%s", f.ComponentAddress.String()), ComponentAddressKey)

		if f.ComponentType == SFMD {
			for i := 0; i < 4; i++ {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     f.Property2Stat.GetBit(1 << i),
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: f.TimeLabels.FireToTimeUnixNano(),
					},
				}
			}
		}
		if f.ComponentType == FireDoor {
			bits := []int{0, 1, 7, 15}
			for i, num := range bits {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     f.Property2Stat.GetBit(num),
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: common.GetTimestamp(),
					},
				}
			}
		}
		if f.ComponentType == FireWaterPressure {
			bits := []int{0, 1, 7}
			for i, num := range bits {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     f.Property2Stat.GetBit(num),
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: common.GetTimestamp(),
					},
				}
			}
		}
		if f.ComponentType == FireProtectionLevel {
			bits := []int{0, 1, 7}
			for i, num := range bits {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     f.Property2Stat.GetBit(num),
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: common.GetTimestamp(),
					},
				}
			}
		}
		if f.ComponentType == IntelligentHoodOfFireHydrant {
			bits := []int{0, 1, 2, 3, 7}
			for i, num := range bits {
				tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
				temp.Data[tmp] = &common.DataValue{
					Value:     f.Property2Stat.GetBit(num),
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: common.GetTimestamp(),
					},
				}
			}

		}
		for i := 0; i < 10; i++ {
			tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
			temp.Data[tmp] = &common.DataValue{
				Value:     f.Property2Stat.GetBit(i),
				Timestamp: common.GetTimestamp(),
				Metadata: common.DataMetadata{
					Type:      "boolean",
					Timestamp: common.GetTimestamp(),
				},
			}
		}
		for i := 0; i < int(f.AnalogNum); i++ {

			tmp.SetDataBitKey(fmt.Sprintf("%d", f.MultiAnalog[i].AnalogType), AnalogTypeKey)
			temp.Data[tmp] = &common.DataValue{
				Value:     f.MultiAnalog[i].Analog,
				Timestamp: common.GetTimestamp(),
				Metadata: common.DataMetadata{
					Type:      "int",
					Timestamp: common.GetTimestamp(),
				},
			}
		}
		back = append(back, temp)
	}

	return back
}

// FireBuildTransparentTransmission 透传模式
type FireBuildTransparentTransmission struct {
	// 系统类型标志(1 字节)系统类型标志、系统地址、部件类型、部件地址的定义同 8.2.1.2。
	SystemType
	// 数据长度
	DataLen uint16
	// 数据
	Data []byte
}

func (m FireData) FireBuildTransparentTransmissionDecode() []FireBuildTransparentTransmission {
	var back = make([]FireBuildTransparentTransmission, 0)

	buf := data.NewBuffer(m.Data)
	//klog.Infof("FireBuildElectricityDecode :%+v", buf)
	for i := 0; i < int(m.ObjectNum); i++ {

		var f FireBuildTransparentTransmission
		f.SystemType = SystemType(buf.ReadByte())
		f.DataLen = buf.ReadUint16()

		if f.DataLen > 0 {
			klog.Infof("FireBuildElectricityDecode-DataLen :%d", f.DataLen)
			f.Data = buf.ReadN(int(f.DataLen))
			klog.Infof("FireBuildElectricityDecode-Data :%x", f.Data)
		}
		//klog.Infof("FireBuildTransparentTransmissionDecode Pos:%d--%v", buf.Pos(), f)
		back = append(back, f)
	}

	return back
}

func (m FireData) FireBuildTransparentTransmissionDecodeToData() []Data {
	klog.Infof("FireBuildTransparentTransmissionDecode:Type:%x,%d", m.DataBaseType, len(m.Data))
	var back []Data
	back = make([]Data, 0)
	fs := m.FireBuildTransparentTransmissionDecode()
	klog.Infof("FireBuildTransparentTransmissionDecode :%+v", fs)
	for _, f := range fs {
		var temp Data
		temp.Data = map[DataKey]*common.DataValue{}
		tmp := NewDataKey()
		tmp.SetDataBitKey(fmt.Sprintf("%d", m.DataBaseType), DataBaseTypeKey)
		tmp.SetDataBitKey(fmt.Sprintf("%d", f.SystemType), SystemTypeKey)

		if f.DataLen > 0 {

			//取出设备信息里的回路和设备编码
			var params string
			switch globals.GetHostNameInfo() {
			case globals.M6v3floors12HostName:

			case globals.SoftwareparkHostName:

			case globals.JiyunHostName:

			case globals.TongJiHostName:

			case globals.BoXingHostName:

			case globals.XIASHAHostName:

			case globals.BJM3HostName:
			case globals.TengrenHostName:
			case globals.SHWAIGAOQIAOHostName:
			case globals.BJM6:
			case globals.BEIJINGB28:
			case globals.TEST:
			default:
				f.ParseDatabtobaone()
				klog.Infof("测试清洗结果: %s", params)

			}
			//for i := 0; i < 4; i++ {
			//	tmp.SetDataBitKey(fmt.Sprintf("%d", i), PropertyIDKey)
			//	temp.Data[tmp] = &common.DataValue{
			//		Value:     f.Property2Stat.GetBit(1 << i),
			//		Timestamp: common.GetTimestamp(),
			//		Metadata: common.DataMetadata{
			//			Type:      "boolean",
			//			Timestamp: f.TimeLabels.FireToTimeUnixNano(),
			//		},
			//	}
			//}
		}

		back = append(back, temp)
	}

	return back
}

type FirePrintMessage struct {
	// 时间
	PrintTime []byte // 0x1b0x36 时间
	// 回路信息
	PrintLoopAddress []byte // 0x1b0x39 时间
	// 配置信息
	PrintConfigMessage []byte
}

type FireAlarmInformation struct {
	// 时间
	LoopAddress string //回路地址
	// 系统状态数据为 2 字节，低字节传输在前。
	Property2Stat
	// 配置信息
	Message []string
	TimeLabels
}

func (ftt *FireBuildTransparentTransmission) ParseDatabtobaone() []FireAlarmInformation {
	//	var back = make([]FireAlarmInformation, 0)

	if ftt.DataLen > 0 {
		var fm FirePrintMessage
		one := bytes.Split(ftt.Data, []byte{0x0a})
		if len(one) < 2 {
			return nil
		}
		fm.PrintTime = bytes.TrimPrefix(one[0], []byte{0x1b, 0x36})
		klog.Infof("B28清洗: PrintTime[%s]", fm.PrintTime)
		klog.Infof("B28清洗: PrintLoopAddress[%x]", one[1])
		fm.PrintLoopAddress = bytes.TrimSpace(one[1])
		klog.Infof("B28清洗: PrintLoopAddress[%s]", fm.PrintLoopAddress)
	}
	return nil
}
func StringStripJiyun(place, msg string) string {
	klog.Infof("纪蕴清洗: place[%s]", place)
	klog.Infof("纪蕴清洗: msg[%s]", msg)
	if place == "" && msg == "" {
		return ""
	}
	place = strings.TrimSpace(place)
	msg = strings.TrimSpace(msg)

	if strings.Contains(place, "QN01") {
		infos := strings.SplitN(place, ",", 4)

		reg := regexp.MustCompile(`[\W|_]{1,}`)
		return reg.ReplaceAllString(strings.Join(infos[2:], ""), "")
	}
	if strings.Contains(place, "回路") {
		reg := regexp.MustCompile(`[\W|_]{1,}`)
		return reg.ReplaceAllString(place, "")
	}
	if strings.Contains(msg, "QN01") {
		infos := strings.SplitN(place, ",", 4)
		reg := regexp.MustCompile(`[\W|_]{1,}`)
		return reg.ReplaceAllString(strings.Join(infos[2:], ""), "")
	}
	if strings.Contains(msg, "回路") {
		reg := regexp.MustCompile(`[\W|_]{1,}`)
		return reg.ReplaceAllString(msg, "")
	}
	return ""
}

func StringStripDefault(name, place, msg string) string {
	klog.Infof("通用清洗: name[%s]", name)
	klog.Infof("通用清洗: place[%s]", place)
	klog.Infof("通用清洗: msg[%s]", msg)
	if name == "" && place == "" && msg == "" {
		return ""
	}
	name = strings.TrimSpace(name)
	place = strings.TrimSpace(place)
	msg = strings.TrimSpace(msg)

	if strings.ContainsAny(name, "AL") {
		reg := regexp.MustCompile(`[\W|_]{1,}`)
		start := strings.IndexAny(name, "AL")
		return reg.ReplaceAllString(name[start:], "")
	}
	if strings.Contains(name, "回路") {
		reg := regexp.MustCompile(`[\W|_]{1,}`)
		return reg.ReplaceAllString(name, "")
	}
	if strings.Contains(place, "QN01") {
		infos := strings.SplitN(place, ",", 4)
		reg := regexp.MustCompile(`[\W|_]{1,}`)
		return strings.TrimPrefix(reg.ReplaceAllString(strings.Join(infos[2:], ""), ""), "QN01")
	}
	if strings.Contains(place, "回路") {
		reg := regexp.MustCompile(`[\W|_]{1,}`)
		return reg.ReplaceAllString(place, "")
	}
	if strings.Contains(msg, "QN01") {
		infos := strings.SplitN(msg, ",", 4)
		reg := regexp.MustCompile(`[\W|_]{1,}`)
		return reg.ReplaceAllString(strings.Join(infos[2:], ""), "")
	}
	if strings.Contains(msg, "回路") {
		reg := regexp.MustCompile(`[\W|_]{1,}`)
		return reg.ReplaceAllString(msg, "")
	}
	return ""
}

func StringStripBluebird(name, place, msg string) string {
	klog.Infof("青鸟通用清洗: name[%s]", name)
	klog.Infof("青鸟通用清洗: place[%s]", place)
	klog.Infof("青鸟通用清洗: msg[%s]", msg)
	if name == "" && place == "" && msg == "" {
		return ""
	}
	name = strings.TrimSpace(name)
	place = strings.TrimSpace(place)
	msg = strings.TrimSpace(msg)

	if strings.ContainsAny(name, "AL") {
		reg := regexp.MustCompile(`[\W|_]{1,}`)
		start := strings.IndexAny(name, "AL")
		return reg.ReplaceAllString(name[start:], "")
	}
	if strings.Contains(name, "回路") {
		reg := regexp.MustCompile(`[\W|_]{1,}`)
		return reg.ReplaceAllString(name, "")
	}
	if strings.Contains(place, "QN01") {
		infos := strings.SplitN(place, ",", 4)
		reg := regexp.MustCompile(`[\W|_]{1,}`)
		return strings.TrimPrefix(reg.ReplaceAllString(strings.Join(infos[1:], ""), ""), "QN01")
	}
	if strings.Contains(place, "回路") {
		reg := regexp.MustCompile(`[\W|_]{1,}`)
		return reg.ReplaceAllString(place, "")
	}
	if strings.Contains(msg, "QN01") {
		infos := strings.SplitN(msg, ",", 4)
		reg := regexp.MustCompile(`[\W|_]{1,}`)
		return reg.ReplaceAllString(strings.Join(infos[1:], ""), "")
	}
	if strings.Contains(msg, "回路") {
		reg := regexp.MustCompile(`[\W|_]{1,}`)
		return reg.ReplaceAllString(msg, "")
	}
	return ""
}

func StringStripSanhe(name string) string {
	klog.Infof("通用清洗: name[%s]", name)
	if name == "" {
		return ""
	}
	name = strings.TrimSpace(name)
	return StringStrip(name)
}
func StringStripSanhe9(name string) string {
	klog.Infof("通用清洗: name[%s]", name)

	if name == "" {
		return ""
	}
	name = strings.TrimSpace(name)
	if strings.ContainsAny(name, "AL") {
		reg := regexp.MustCompile(`[\W|_]{1,}`)
		start := strings.IndexAny(name, "AL")
		return reg.ReplaceAllString(name[start:], "")
	}
	return ""
}

func StringStripbtobaone(place, msg string) string {
	klog.Infof("B28清洗: place[%s]", place)
	klog.Infof("B28清洗: msg[%s]", msg)
	if place == "" {
		return ""
	}
	place = strings.TrimSpace(place)
	reg := regexp.MustCompile(`[\W|_]{1,}`)
	return reg.ReplaceAllString(place, "")
}
func StringStriphedan(place, msg string) string {
	klog.Infof("荷担清洗: place[%s]", place)
	klog.Infof("荷担清洗: msg[%s]", msg)
	if place == "" {
		return ""
	}
	place = strings.TrimSpace(place)
	reg := regexp.MustCompile(`[\W|_]{1,}`)
	infos := strings.SplitN(place, ":", 2)
	return reg.ReplaceAllString(infos[0], "")
}
func StringStrip(str string) string {
	if str == "" {
		return ""
	}
	str = strings.TrimSpace(str)
	reg := regexp.MustCompile(`[\W|_]{1,}`)
	return reg.ReplaceAllString(str, "")
}
func StringStripTengren(str string) string {
	if str == "" {
		return ""
	}
	klog.Infof("腾仁清洗: %s", str)
	space := strings.TrimSpace(str)
	infos := strings.SplitN(space, ",", 4)

	reg := regexp.MustCompile(`[\W|_]{1,}`)
	return reg.ReplaceAllString(strings.Join(infos[1:], ""), "")
}
func StringStripBoxxing(str string) string {
	if str == "" {
		return ""
	}
	klog.Infof("博兴清洗: %s", str)
	space := strings.TrimSpace(str)
	infos := strings.SplitN(space, ",", 4)

	reg := regexp.MustCompile(`[\W|_]{1,}`)
	return reg.ReplaceAllString(strings.Join(infos[1:], ""), "")
}

func StringStripWaigaoqiao(str string) string {
	if str == "" {
		return ""
	}
	klog.Infof("外高桥: %s", str)
	space := strings.TrimSpace(str)
	infos := strings.SplitN(space, ",", 4)

	reg := regexp.MustCompile(`[\W|_]{1,}`)
	return reg.ReplaceAllString(strings.Join(infos[0:], ""), "")
}
