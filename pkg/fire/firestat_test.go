package fire

import (
	"github.com/imroc/biu"
	"k8s.io/klog/v2"
	"testing"
)

func TestFireStat(t *testing.T) {
	var stat = Property2Stat{
		Status: [2]byte{0x00, 0x00},
	}
	//GoodRunSYS         = iota //系统运行正常
	//FireAlarmSYS              //火警
	//FireMalfunctionSYS        //故障
	//ScreenSYS                 //屏蔽
	//RegulatorySYS             //监管
	//StartSYS                  //启动
	//FeedbackSYS               //反馈
	//RockStateSYS              //岩石状态
	//MainPowerFailSYS          //主电故障
	//BackupPowerFailSYS        //备电故障
	//BusFailSYS                //总线故障
	//ManualStateSYS            //手动状态
	//ConfigChangesSYS          //配置改变
	//ResetSYS                  //复位
	stat.SetBit(RunStateSYS)
	klog.Infof("%s%s", biu.ToBinaryString(stat.Status[0]), biu.ToBinaryString(stat.Status[1]))
	stat.SetBit(FireAlarmSYS)
	klog.Infof("%s%s", biu.ToBinaryString(stat.Status[0]), biu.ToBinaryString(stat.Status[1]))
	stat.SetBit(FireMalfunctionSYS)
	klog.Infof("%s%s", biu.ToBinaryString(stat.Status[0]), biu.ToBinaryString(stat.Status[1]))
	stat.SetBit(OperateSYS)
	klog.Infof("%s%s", biu.ToBinaryString(stat.Status[0]), biu.ToBinaryString(stat.Status[1]))
	stat.SetBit(RegulatorySYS)
	klog.Infof("%s%s", biu.ToBinaryString(stat.Status[0]), biu.ToBinaryString(stat.Status[1]))
	stat.SetBit(StartSYS)
	klog.Infof("%s%s", biu.ToBinaryString(stat.Status[0]), biu.ToBinaryString(stat.Status[1]))
	stat.SetBit(FeedbackSYS)
	klog.Infof("%s%s", biu.ToBinaryString(stat.Status[0]), biu.ToBinaryString(stat.Status[1]))
	stat.SetBit(DeferStateSYS)
	klog.Infof("%s%s", biu.ToBinaryString(stat.Status[0]), biu.ToBinaryString(stat.Status[1]))
	stat.SetBit(MainPowerFailSYS)
	klog.Infof("%s%s", biu.ToBinaryString(stat.Status[0]), biu.ToBinaryString(stat.Status[1]))
	stat.SetBit(BackupPowerFailSYS)
	klog.Infof("%s%s", biu.ToBinaryString(stat.Status[0]), biu.ToBinaryString(stat.Status[1]))
	stat.SetBit(BusFailSYS)
	klog.Infof("%s%s", biu.ToBinaryString(stat.Status[0]), biu.ToBinaryString(stat.Status[1]))
	stat.SetBit(ManualStateSYS)
	klog.Infof("%s%s", biu.ToBinaryString(stat.Status[0]), biu.ToBinaryString(stat.Status[1]))
	stat.SetBit(ConfigChangesSYS)
	klog.Infof("%s%s", biu.ToBinaryString(stat.Status[0]), biu.ToBinaryString(stat.Status[1]))
	stat.SetBit(ResetSYS)
	klog.Infof("%s%s", biu.ToBinaryString(stat.Status[0]), biu.ToBinaryString(stat.Status[1]))
	for i := 0; i < 16; i++ {
		klog.Infof("第%d位状态%t", i, stat.GetBit(i))
	}

}
func TestFireSFMDStat(t *testing.T) {
	klog.Infof("1<<0 = %d", 1<<0)
	klog.Infof("1<<1 = %d", 1<<1)
	klog.Infof("1<<2 = %d", 1<<2)
	klog.Infof("1<<3 = %d", 1<<3)
}
func TestFireFireDoorStat(t *testing.T) {

	bits := []int{0, 1, 7, 15}
	for i, num := range bits {
		klog.Infof("%d= %d", i, num)
	}
}

func TestFireSetDataBitKey(t *testing.T) {
	tmp := NewDataKey()

	tmp.SetDataBitKey("2", 2)
	klog.Infof("%v", tmp)
}
