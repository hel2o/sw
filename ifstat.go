package sw

import (
	"fmt"
	"github.com/spf13/cast"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
)

var (
	ifNameOid        = []string{"1.3.6.1.2.1.31.1.1.1.1"}
	ifNameOidPrefix  = ".1.3.6.1.2.1.31.1.1.1.1."
	ifHCInOid        = []string{"1.3.6.1.2.1.31.1.1.1.6"}
	ifHCInOidPrefix  = ".1.3.6.1.2.1.31.1.1.1.6."
	ifHCOutOid       = []string{"1.3.6.1.2.1.31.1.1.1.10"}
	ifHCOutOidPrefix = ".1.3.6.1.2.1.31.1.1.1.10."

	ifHCInPktsOid       = []string{"1.3.6.1.2.1.31.1.1.1.7"}
	ifHCInPktsOidPrefix = ".1.3.6.1.2.1.31.1.1.1.7."
	ifHCOutPktsOid      = []string{"1.3.6.1.2.1.31.1.1.1.11"}
	//up(1),down(2)
	ifOperStatusOid              = []string{"1.3.6.1.2.1.2.2.1.8"}
	ifOperStatusOidPrefix        = ".1.3.6.1.2.1.2.2.1.8."
	ifHCInBroadcastPktsOid       = []string{"1.3.6.1.2.1.31.1.1.1.9"}
	ifHCInBroadcastPktsOidPrefix = ".1.3.6.1.2.1.31.1.1.1.9."
	ifHCOutBroadcastPktsOid      = []string{"1.3.6.1.2.1.31.1.1.1.13"}
	// multicastpkt
	ifHCInMulticastPktsOid       = []string{"1.3.6.1.2.1.31.1.1.1.8"}
	ifHCInMulticastPktsOidPrefix = ".1.3.6.1.2.1.31.1.1.1.8."
	ifHCOutMulticastPktsOid      = []string{"1.3.6.1.2.1.31.1.1.1.12"}
	// speed 配置
	ifSpeedOid       = []string{"1.3.6.1.2.1.31.1.1.1.15"}
	ifSpeedOidPrefix = ".1.3.6.1.2.1.31.1.1.1.15."

	// Discards配置
	ifInDiscardsOid       = []string{"1.3.6.1.2.1.2.2.1.13"}
	ifInDiscardsOidPrefix = ".1.3.6.1.2.1.2.2.1.13."
	ifOutDiscardsOid      = []string{"1.3.6.1.2.1.2.2.1.19"}

	// Errors配置
	ifInErrorsOid        = []string{"1.3.6.1.2.1.2.2.1.14"}
	ifInErrorsOidPrefix  = ".1.3.6.1.2.1.2.2.1.14."
	ifOutErrorsOid       = []string{"1.3.6.1.2.1.2.2.1.20"}
	ifOutErrorsOidPrefix = ".1.3.6.1.2.1.2.2.1.20."

	//ifInUnknownProtos 由于未知或不支持的网络协议而丢弃的输入报文的数量
	ifInUnknownProtosOid    = []string{"1.3.6.1.2.1.2.2.1.15"}
	ifInUnknownProtosPrefix = ".1.3.6.1.2.1.2.2.1.15."

	//ifOutQLen 接口上输出报文队列长度
	ifOutQLenOid    = []string{"1.3.6.1.2.1.2.2.1.21"}
	ifOutQLenPrefix = ".1.3.6.1.2.1.2.2.1.21."

	//2011 华为 早期华三交换机
	//25506 新华三交换机
	//56813 华为智选
	//4881 锐捷交换机
	//45577 信锐安视交换机

	//二层端口类型
	//华为 trunk(1) invalid(0) access(2) hybrid(3) fabric(4) qinq(5) desirable(6) auto(7)
	//华三vLANTrunk(1),access(2),hybrid(3),fabric(4) hh3cifVLANType
	//锐捷 access(1), trunk(2), dot1q-tunnel(3),hybrid(4), other(5), uplink(6),host(7) or promiscuous(8) port.
	l2IfPortTypeOid = []string{"1.3.6.1.4.1.2011.5.25.42.1.1.1.3.1.3", "1.3.6.1.4.1.56813.5.25.42.1.1.1.3.1.3", "1.3.6.1.4.1.4881.1.1.10.2.10.1.1.1.24", "1.3.6.1.4.1.25506.8.35.1.1.1.5", "1.3.6.1.4.1.45577.5.7.9.1.3"}

	//二层端口的VLAN ID
	//取值范围为0～4094。如果设置为0，则hwL2IfPVID恢复为缺省值1
	l2IfPvidOid    = []string{"1.3.6.1.4.1.2011.5.25.42.1.1.1.3.1.4", "1.3.6.1.4.1.56813.5.25.42.1.1.1.3.1.4", "1.3.6.1.4.1.4881.1.1.10.2.9.1.6.1.3", "1.3.6.1.4.1.45577.5.7.9.1.6", "1.3.6.1.2.1.17.7.1.4.5.1.1"}
	l2IfPvidPrefix = []string{".1.3.6.1.4.1.2011.5.25.42.1.1.1.3.1.4.", ".1.3.6.1.4.1.56813.5.25.42.1.1.1.3.1.4.", ".1.3.6.1.4.1.4881.1.1.10.2.9.1.6.1.3.", ".1.3.6.1.4.1.45577.5.7.9.1.6.", ".1.3.6.1.2.1.17.7.1.4.5.1.1."}

	//trunk放行的VLAN  信锐有这个属性
	l2IfVlanUntaggedOid = []string{".1.3.6.1.4.1.45577.5.7.9.1.4"}

	//二层端口的模式
	//INTEGER : copper(1)2: fiber(2)3: other(3)
	ethernetPortModeOid    = []string{"1.3.6.1.4.1.2011.5.25.157.1.1.1.1.39", "1.3.6.1.4.1.56813.5.25.157.1.1.1.1.39", "1.3.6.1.4.1.4881.1.1.10.2.10.1.1.1.15"}
	ethernetPortModePrefix = []string{".1.3.6.1.4.1.2011.5.25.157.1.1.1.1.39.", ".1.3.6.1.4.1.56813.5.25.157.1.1.1.1.39.", ".1.3.6.1.4.1.4881.1.1.10.2.10.1.1.1.15."}

	//以太网接口的双工模式
	//INTEGER : full(1) half(2)
	ethernetDuplexOid    = []string{"1.3.6.1.4.1.2011.5.25.157.1.1.1.1.14", "1.3.6.1.4.1.56813.5.25.157.1.1.1.1.14", "1.3.6.1.4.1.4881.1.1.10.2.10.1.1.1.8", "1.3.6.1.4.1.25506.8.35.5.1.4.1.3"}
	ethernetDuplexPrefix = []string{".1.3.6.1.4.1.2011.5.25.157.1.1.1.1.14.", ".1.3.6.1.4.1.56813.5.25.157.1.1.1.1.14.", ".1.3.6.1.4.1.4881.1.1.10.2.10.1.1.1.8.", ".1.3.6.1.4.1.25506.8.35.5.1.4.1.3."}

	//是否admin down up
	//up(1),down(2),testing(3)
	ifAdminStatusOid    = []string{"1.3.6.1.2.1.2.2.1.7"}
	ifAdminStatusPrefix = ".1.3.6.1.2.1.2.2.1.7."

	//接口描述
	ifDescrOid    = []string{"1.3.6.1.2.1.2.2.1.2"}
	ifDescrPrefix = ".1.3.6.1.2.1.2.2.1.2."

	ifVlanPortListOid = []string{"1.3.6.1.4.1.25506.8.35.2.1.1.1.19"}

	hwL2VlanPortListOid = "1.3.6.1.4.1.2011.5.25.42.3.1.1.1.1.3"  //华为 port vlan list
	fmL2VlanPortListOid = "1.3.6.1.4.1.56813.5.25.42.3.1.1.1.1.3" //华为智选 port vlan list
	l2VlanPortListOid   = "1.3.6.1.2.1.17.7.1.4.3.1.2"            //通用 port vlan list

	hwL2IfPortNameOid = "1.3.6.1.4.1.2011.5.25.42.1.1.1.3.1.19"  //获取二层接口索引对应的接口列表 华为
	fmL2IfPortNameOid = "1.3.6.1.4.1.56813.5.25.42.1.1.1.3.1.19" //获取二层接口索引对应的接口列表 华为智选

	//STP状态
	//1：disabled 2：discarding 4：learning 5：forwarding
	dot1dStpPortStateOid    = []string{"1.3.6.1.2.1.17.2.15.1.3"}
	dot1dStpPortStatePrefix = []string{".1.3.6.1.2.1.17.2.15.1.3."}

	//STP使能状态
	//1:enabled 2:disabled
	dot1dStpPortEnabledOid    = []string{"1.3.6.1.2.1.17.2.15.1.4"}
	dot1dStpPortEnabledPrefix = []string{".1.3.6.1.2.1.17.2.15.1.4."}

	//这个是在获取PortType和PVID时用的，因为只有二层接口才有这个属性
	l2IfName = []string{"1.3.6.1.4.1.2011.5.25.42.1.1.1.3.1.19", "1.3.6.1.4.1.56813.5.25.42.1.1.1.3.1.19", "1.3.6.1.4.1.45577.5.7.9.1.2"}
)

type IfStats struct {
	IfName               string `json:"ifName"`
	IfIndex              int    `json:"ifIndex"`
	IfHCInOctets         uint64 `json:"ifHCInOctets"`
	IfHCOutOctets        uint64 `json:"ifHCOutOctets"`
	IfHCInUcastPkts      uint64 `json:"ifHCInUcastPkts"`
	IfHCOutUcastPkts     uint64 `json:"ifHCOutUcastPkts"`
	IfHCInBroadcastPkts  uint64 `json:"ifHCInBroadcastPkts"`
	IfHCOutBroadcastPkts uint64 `json:"ifHCOutBroadcastPkts"`
	IfHCInMulticastPkts  uint64 `json:"ifHCInMulticastPkts"`
	IfHCOutMulticastPkts uint64 `json:"ifHCOutMulticastPkts"`
	IfSpeed              uint64 `json:"ifSpeed"`
	IfInDiscards         uint64 `json:"ifInDiscards"`
	IfOutDiscards        uint64 `json:"ifOutDiscards"`
	IfInErrors           uint64 `json:"ifInErrors"`
	IfOutErrors          uint64 `json:"ifOutErrors"`
	IfInUnknownProtos    uint64 `json:"ifInUnknownProtos"`
	IfOutQLen            uint64 `json:"ifOutQLen"`
	IfOperStatus         int    `json:"ifOperStatus"`
	L2IfPortType         uint64 `json:"l2IfPortType"`
	L2IfPVID             uint64 `json:"l2IfPVID"`
	L2IfTrunkAllowedVlan string `json:"l2IfTrunkAllowedVlan"`
	EthernetPortMode     uint64 `json:"ethernetPortMode"`
	EthernetDuplex       uint64 `json:"ethernetDuplex"`
	IfAdminStatus        uint64 `json:"ifAdminStatus"`
	IfDescr              string `json:"ifDescr"`
	StpStatus            uint64 `json:"stpStatus"`
	StpEnabled           uint64 `json:"stpEnabled"`
	TS                   int64  `json:"ts"`
}

func (this *IfStats) String() string {
	return fmt.Sprintf("<IfName:%s, IfIndex:%d, IfHCInOctets:%d, IfHCOutOctets:%d>", this.IfName, this.IfIndex, this.IfHCInOctets, this.IfHCOutOctets)
}

type IgnoreIfStats struct {
	IgnoreInOutOctets,
	IgnorePkt,
	IgnoreOperStatus,
	IgnoreBroadcastPkt,
	IgnoreMulticastPkt,
	IgnoreDiscards,
	IgnoreErrors,
	IgnoreUnknownProtos,
	IgnoreOutQLen,
	IgnoreL2IfPortType,
	IgnoreL2IfPVID,
	IgnoreEthernetPortMode,
	IgnoreEthernetDuplex,
	IgnoreIfAdminStatus,
	IgnoreIfDesc,
	IgnoreStpStatus,
	IgnoreIfSpeed bool
}

func ListIfStats(ip, community string, timeout int, interfaces []string, retry int, limitConn int, manufacturer string, useSnmpGetNext bool, ignoreIfStats IgnoreIfStats) ([]IfStats, error) {
	var ifStatsList []IfStats
	var limitCh chan bool
	if limitConn > 0 {
		limitCh = make(chan bool, limitConn)
	} else {
		limitCh = make(chan bool, 1)
	}
	var sleep = 5 * time.Millisecond
	if useSnmpGetNext {
		sleep = 700 * time.Millisecond
	}

	chIfNameList := make(chan []gosnmp.SnmpPDU)
	limitCh <- true
	go ListIfName(ip, community, timeout, chIfNameList, retry, limitCh, useSnmpGetNext)

	var ifInList, ifOutList, ifSpeedList []gosnmp.SnmpPDU
	chIfInList := make(chan []gosnmp.SnmpPDU)
	chIfOutList := make(chan []gosnmp.SnmpPDU)
	chIfSpeedList := make(chan []gosnmp.SnmpPDU)

	if ignoreIfStats.IgnoreInOutOctets == false {
		limitCh <- true
		go ListIfHCInOctets(ip, community, timeout, chIfInList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
		limitCh <- true
		go ListIfHCOutOctets(ip, community, timeout, chIfOutList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	//Speed
	if ignoreIfStats.IgnoreIfSpeed == false {
		limitCh <- true
		go ListIfSpeed(ip, community, timeout, chIfSpeedList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	// OperStatus
	var ifStatusList []gosnmp.SnmpPDU
	chIfStatusList := make(chan []gosnmp.SnmpPDU)
	if ignoreIfStats.IgnoreOperStatus == false {
		limitCh <- true
		go ListIfOperStatus(ip, community, timeout, chIfStatusList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	//L2IfPortType
	var l2IfPortTypeList []gosnmp.SnmpPDU
	chL2IfPortTypeList := make(chan []gosnmp.SnmpPDU)
	if ignoreIfStats.IgnoreL2IfPortType == false {
		limitCh <- true
		go ListL2IfPortType(ip, community, timeout, chL2IfPortTypeList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	//L2IfPVID
	var l2IfPVIDList []gosnmp.SnmpPDU
	chL2IfPVIDList := make(chan []gosnmp.SnmpPDU)
	if ignoreIfStats.IgnoreL2IfPVID == false {
		limitCh <- true
		go ListL2IfPVlanId(ip, community, timeout, chL2IfPVIDList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	//ethernetPortModeOid
	var ethernetPortModeList []gosnmp.SnmpPDU
	chEthernetPortModeList := make(chan []gosnmp.SnmpPDU)
	if ignoreIfStats.IgnoreL2IfPortType == false {
		limitCh <- true
		go ListEthernetPortMode(ip, community, timeout, chEthernetPortModeList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	//EthernetDuplex
	var ethernetDuplexList []gosnmp.SnmpPDU
	chEthernetDuplexList := make(chan []gosnmp.SnmpPDU)
	if ignoreIfStats.IgnoreL2IfPortType == false {
		limitCh <- true
		go ListEthernetDuplex(ip, community, timeout, chEthernetDuplexList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	//AdminStatus
	var ifAdminStatusList []gosnmp.SnmpPDU
	chIfAdminStatusList := make(chan []gosnmp.SnmpPDU)
	if ignoreIfStats.IgnoreIfAdminStatus == false {
		limitCh <- true
		go ListIfAdminStatus(ip, community, timeout, chIfAdminStatusList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	//IfDescr
	var ifDescrList []gosnmp.SnmpPDU
	chIfDescrList := make(chan []gosnmp.SnmpPDU)
	if ignoreIfStats.IgnoreIfDesc == false {
		limitCh <- true
		go ListIfDescr(ip, community, timeout, chIfDescrList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	var ifStpStatusList, ifStpEnabledList []gosnmp.SnmpPDU
	chIfStpStatusList := make(chan []gosnmp.SnmpPDU)
	chIfStpEnabledList := make(chan []gosnmp.SnmpPDU)
	if ignoreIfStats.IgnoreStpStatus == false {
		limitCh <- true
		go ListIfStpStatus(ip, community, timeout, chIfStpStatusList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
		limitCh <- true
		go ListIfStpEnabled(ip, community, timeout, chIfStpEnabledList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	var ifInPktList, ifOutPktList []gosnmp.SnmpPDU
	chIfInPktList := make(chan []gosnmp.SnmpPDU)
	chIfOutPktList := make(chan []gosnmp.SnmpPDU)
	if ignoreIfStats.IgnorePkt == false {
		limitCh <- true
		go ListIfHCInUcastPkts(ip, community, timeout, chIfInPktList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
		limitCh <- true
		go ListIfHCOutUcastPkts(ip, community, timeout, chIfOutPktList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	var ifInBroadcastPktList, ifOutBroadcastPktList []gosnmp.SnmpPDU
	chIfInBroadcastPktList := make(chan []gosnmp.SnmpPDU)
	chIfOutBroadcastPktList := make(chan []gosnmp.SnmpPDU)
	if ignoreIfStats.IgnoreBroadcastPkt == false {
		limitCh <- true
		go ListIfHCInBroadcastPkts(ip, community, timeout, chIfInBroadcastPktList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
		limitCh <- true
		go ListIfHCOutBroadcastPkts(ip, community, timeout, chIfOutBroadcastPktList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	var ifInMulticastPktList, ifOutMulticastPktList []gosnmp.SnmpPDU
	chIfInMulticastPktList := make(chan []gosnmp.SnmpPDU)
	chIfOutMulticastPktList := make(chan []gosnmp.SnmpPDU)
	if ignoreIfStats.IgnoreMulticastPkt == false {
		limitCh <- true
		go ListIfHCInMulticastPkts(ip, community, timeout, chIfInMulticastPktList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
		limitCh <- true
		go ListIfHCOutMulticastPkts(ip, community, timeout, chIfOutMulticastPktList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	//Discards
	var ifInDiscardsList, ifOutDiscardsList []gosnmp.SnmpPDU
	chIfInDiscardsList := make(chan []gosnmp.SnmpPDU)
	chIfOutDiscardsList := make(chan []gosnmp.SnmpPDU)
	if ignoreIfStats.IgnoreDiscards == false {
		limitCh <- true
		go ListIfInDiscards(ip, community, timeout, chIfInDiscardsList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
		limitCh <- true
		go ListIfOutDiscards(ip, community, timeout, chIfOutDiscardsList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	//Errors
	var ifInErrorsList, ifOutErrorsList []gosnmp.SnmpPDU
	chIfInErrorsList := make(chan []gosnmp.SnmpPDU)
	chIfOutErrorsList := make(chan []gosnmp.SnmpPDU)
	if ignoreIfStats.IgnoreErrors == false {
		limitCh <- true
		go ListIfInErrors(ip, community, timeout, chIfInErrorsList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
		limitCh <- true
		go ListIfOutErrors(ip, community, timeout, chIfOutErrorsList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	//UnknownProtos
	var ifInUnknownProtosList []gosnmp.SnmpPDU
	chIfInUnknownProtosList := make(chan []gosnmp.SnmpPDU)
	if ignoreIfStats.IgnoreUnknownProtos == false {
		limitCh <- true
		go ListIfInUnknownProtos(ip, community, timeout, chIfInUnknownProtosList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	//QLen
	var ifOutQLenList []gosnmp.SnmpPDU
	chIfOutQLenList := make(chan []gosnmp.SnmpPDU)
	if ignoreIfStats.IgnoreOutQLen == false {
		limitCh <- true
		go ListIfOutQLen(ip, community, timeout, chIfOutQLenList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	//开始取出数据
	ifNameList := <-chIfNameList

	if ignoreIfStats.IgnoreInOutOctets == false {
		ifInList = <-chIfInList
		ifOutList = <-chIfOutList
	}
	if ignoreIfStats.IgnoreIfSpeed == false {
		ifSpeedList = <-chIfSpeedList
	}
	if ignoreIfStats.IgnoreOperStatus == false {
		ifStatusList = <-chIfStatusList
	}
	if ignoreIfStats.IgnoreL2IfPortType == false {
		l2IfPortTypeList = <-chL2IfPortTypeList
	}
	if ignoreIfStats.IgnoreL2IfPVID == false {
		l2IfPVIDList = <-chL2IfPVIDList
	}
	if ignoreIfStats.IgnoreEthernetPortMode == false {
		ethernetPortModeList = <-chEthernetPortModeList
	}
	if ignoreIfStats.IgnoreEthernetDuplex == false {
		ethernetDuplexList = <-chEthernetDuplexList
	}
	if ignoreIfStats.IgnoreIfAdminStatus == false {
		ifAdminStatusList = <-chIfAdminStatusList
	}
	if ignoreIfStats.IgnorePkt == false {
		ifInPktList = <-chIfInPktList
		ifOutPktList = <-chIfOutPktList
	}
	if ignoreIfStats.IgnoreIfDesc == false {
		ifDescrList = <-chIfDescrList
	}
	if ignoreIfStats.IgnoreBroadcastPkt == false {
		ifInBroadcastPktList = <-chIfInBroadcastPktList
		ifOutBroadcastPktList = <-chIfOutBroadcastPktList
	}
	if ignoreIfStats.IgnoreMulticastPkt == false {
		ifInMulticastPktList = <-chIfInMulticastPktList
		ifOutMulticastPktList = <-chIfOutMulticastPktList
	}
	if ignoreIfStats.IgnoreDiscards == false {
		ifInDiscardsList = <-chIfInDiscardsList
		ifOutDiscardsList = <-chIfOutDiscardsList
	}
	if ignoreIfStats.IgnoreErrors == false {
		ifInErrorsList = <-chIfInErrorsList
		ifOutErrorsList = <-chIfOutErrorsList
	}
	if ignoreIfStats.IgnoreUnknownProtos == false {
		ifInUnknownProtosList = <-chIfInUnknownProtosList
	}
	if ignoreIfStats.IgnoreOutQLen == false {
		ifOutQLenList = <-chIfOutQLenList
	}
	if ignoreIfStats.IgnoreStpStatus == false {
		ifStpStatusList = <-chIfStpStatusList
		ifStpEnabledList = <-chIfStpEnabledList
	}

	if len(ifNameList) > 0 {
		now := time.Now().Unix()
		for _, ifNamePDU := range ifNameList {
			ifName := string(ifNamePDU.Value.([]byte))

			check := false
			if len(interfaces) > 0 {
				for _, inter := range interfaces {
					//以包含的接口名开头
					if strings.HasPrefix(strings.ToLower(ifName), strings.ToLower(inter)) {
						check = true
						break
					}
				}
			}

			if check {
				var ifStats IfStats
				ifStats.TS = now
				ifStats.IfName = ifName
				ifIndexStr := strings.Replace(ifNamePDU.Name, ifNameOidPrefix, "", 1)
				ifStats.IfIndex, _ = strconv.Atoi(ifIndexStr)

				if ignoreIfStats.IgnoreInOutOctets == false && len(ifInList) > 0 && len(ifOutList) > 0 {
					for ti, ifHCInOctetsPDU := range ifInList {
						if strings.Replace(ifHCInOctetsPDU.Name, ifHCInOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfHCInOctets = gosnmp.ToBigInt(ifInList[ti].Value).Uint64()
							break
						}
					}
					for ti, ifHCOutOctetsPDU := range ifOutList {
						if strings.Replace(ifHCOutOctetsPDU.Name, ifHCOutOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfHCOutOctets = gosnmp.ToBigInt(ifOutList[ti].Value).Uint64()
							break
						}
					}
				}
				if ignoreIfStats.IgnorePkt == false {
					for ti, ifHCInPktsPDU := range ifInPktList {
						if strings.Replace(ifHCInPktsPDU.Name, ifHCInPktsOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfHCOutUcastPkts = gosnmp.ToBigInt(ifOutPktList[ti].Value).Uint64()
							ifStats.IfHCInUcastPkts = gosnmp.ToBigInt(ifInPktList[ti].Value).Uint64()
							break
						}
					}
				}
				if ignoreIfStats.IgnoreBroadcastPkt == false {
					for ti, ifHCInBroadcastPktPDU := range ifInBroadcastPktList {
						if strings.Replace(ifHCInBroadcastPktPDU.Name, ifHCInBroadcastPktsOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfHCInBroadcastPkts = gosnmp.ToBigInt(ifInBroadcastPktList[ti].Value).Uint64()
							ifStats.IfHCOutBroadcastPkts = gosnmp.ToBigInt(ifOutBroadcastPktList[ti].Value).Uint64()
							break
						}
					}
				}
				if ignoreIfStats.IgnoreMulticastPkt == false {
					for ti, ifHCInMulticastPktPDU := range ifInMulticastPktList {
						if strings.Replace(ifHCInMulticastPktPDU.Name, ifHCInMulticastPktsOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfHCInMulticastPkts = gosnmp.ToBigInt(ifInMulticastPktList[ti].Value).Uint64()
							ifStats.IfHCOutMulticastPkts = gosnmp.ToBigInt(ifOutMulticastPktList[ti].Value).Uint64()
							break
						}
					}
				}

				if ignoreIfStats.IgnoreDiscards == false {
					for ti, ifInDiscardsPDU := range ifInDiscardsList {
						if strings.Replace(ifInDiscardsPDU.Name, ifInDiscardsOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfInDiscards = gosnmp.ToBigInt(ifInDiscardsList[ti].Value).Uint64()
							ifStats.IfOutDiscards = gosnmp.ToBigInt(ifOutDiscardsList[ti].Value).Uint64()
							break
						}
					}
				}

				if ignoreIfStats.IgnoreErrors == false {
					for ti, ifInErrorsPDU := range ifInErrorsList {
						if strings.Replace(ifInErrorsPDU.Name, ifInErrorsOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfInErrors = gosnmp.ToBigInt(ifInErrorsList[ti].Value).Uint64()
							break
						}
					}
					for ti, ifOutErrorsPDU := range ifOutErrorsList {
						if strings.Replace(ifOutErrorsPDU.Name, ifOutErrorsOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfOutErrors = gosnmp.ToBigInt(ifOutErrorsList[ti].Value).Uint64()
							break
						}
					}
				}

				if ignoreIfStats.IgnoreOperStatus == false {
					for ti, ifOperStatusPDU := range ifStatusList {
						if strings.Replace(ifOperStatusPDU.Name, ifOperStatusOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfOperStatus = int(gosnmp.ToBigInt(ifStatusList[ti].Value).Int64())
							break
						}
					}
				}

				if ignoreIfStats.IgnoreL2IfPVID == false {
					for ti, l2IfPvidPDU := range l2IfPVIDList {
						for _, oidPrefix := range l2IfPvidPrefix {
							if strings.Replace(l2IfPvidPDU.Name, oidPrefix, "", 1) == ifIndexStr {
								ifStats.L2IfPVID = gosnmp.ToBigInt(l2IfPVIDList[ti].Value).Uint64()
								break
							}
						}
					}
				}

				if ignoreIfStats.IgnoreStpStatus == false {
					for ti, ifStpStatusPDU := range ifStpStatusList {
						for _, oidPrefix := range dot1dStpPortStatePrefix {
							if strings.Replace(ifStpStatusPDU.Name, oidPrefix, "", 1) == ifIndexStr {
								ifStats.StpStatus = gosnmp.ToBigInt(ifStpStatusList[ti].Value).Uint64()
								break
							}
						}
					}

					if ignoreIfStats.IgnoreStpStatus == false {
						for ti, ifStpEnabledPDU := range ifStpEnabledList {
							for _, oidPrefix := range dot1dStpPortEnabledPrefix {
								if strings.Replace(ifStpEnabledPDU.Name, oidPrefix, "", 1) == ifIndexStr {
									ifStats.StpEnabled = gosnmp.ToBigInt(ifStpEnabledList[ti].Value).Uint64()
									break
								}
							}
						}
					}
				}

				if ignoreIfStats.IgnoreEthernetPortMode == false {
					for ti, ethernetPortModePDU := range ethernetPortModeList {
						for _, oidPrefix := range ethernetPortModePrefix {
							if strings.Replace(ethernetPortModePDU.Name, oidPrefix, "", 1) == ifIndexStr {
								ifStats.EthernetPortMode = gosnmp.ToBigInt(ethernetPortModeList[ti].Value).Uint64()
								break
							}
						}
					}
				}

				if ignoreIfStats.IgnoreEthernetDuplex == false {
					for ti, ethernetDuplexPDU := range ethernetDuplexList {
						for _, oidPrefix := range ethernetDuplexPrefix {
							if strings.Replace(ethernetDuplexPDU.Name, oidPrefix, "", 1) == ifIndexStr {
								ifStats.EthernetDuplex = gosnmp.ToBigInt(ethernetDuplexList[ti].Value).Uint64()
								break
							}
						}
					}
				}

				if ignoreIfStats.IgnoreIfAdminStatus == false {
					for ti, ifAdminStatusPDU := range ifAdminStatusList {
						if strings.Replace(ifAdminStatusPDU.Name, ifAdminStatusPrefix, "", 1) == ifIndexStr {
							ifStats.IfAdminStatus = gosnmp.ToBigInt(ifAdminStatusList[ti].Value).Uint64()
							break
						}
					}
				}

				if ignoreIfStats.IgnoreIfDesc == false {
					for ti, ifDescrPDU := range ifDescrList {
						if strings.Replace(ifDescrPDU.Name, ifDescrPrefix, "", 1) == ifIndexStr {
							ifStats.IfDescr = string(ifDescrList[ti].Value.([]byte))
							break
						}
					}
				}

				if ignoreIfStats.IgnoreUnknownProtos == false {
					for ti, ifInUnknownProtosPDU := range ifInUnknownProtosList {
						if strings.Replace(ifInUnknownProtosPDU.Name, ifInUnknownProtosPrefix, "", 1) == ifIndexStr {
							ifStats.IfInUnknownProtos = gosnmp.ToBigInt(ifInUnknownProtosList[ti].Value).Uint64()
							break
						}
					}
				}

				if ignoreIfStats.IgnoreOutQLen == false {
					for ti, ifOutQLenPDU := range ifOutQLenList {
						if strings.Replace(ifOutQLenPDU.Name, ifOutQLenPrefix, "", 1) == ifIndexStr {
							ifStats.IfOutQLen = gosnmp.ToBigInt(ifOutQLenList[ti].Value).Uint64()
							break
						}
					}
				}
				if ignoreIfStats.IgnoreIfSpeed == false {
					for ti, ifSpeedPDU := range ifSpeedList {
						if strings.Replace(ifSpeedPDU.Name, ifSpeedOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfSpeed = 1000 * 1000 * gosnmp.ToBigInt(ifSpeedList[ti].Value).Uint64()
							break
						}
					}
				}
				ifStatsList = append(ifStatsList, ifStats)
			}
		}
	}

	//以下是处理采集接口类型和PVID的情况
	if ignoreIfStats.IgnoreL2IfPortType == false || ignoreIfStats.IgnoreL2IfPVID == false {
		var ifPortVlanMap = make(map[string][]int)
		var vlanPortListPDU, portNameListPDU []gosnmp.SnmpPDU
		var vlanPortListOid, portNameListOid []string
		var offset int

		//处理华为和华为智选的情况
		if strings.Contains(manufacturer, Huawei) || strings.Contains(manufacturer, FutureMatrix) {
			vlanPortListOid = []string{hwL2VlanPortListOid, fmL2VlanPortListOid}
			portNameListOid = []string{hwL2IfPortNameOid, fmL2IfPortNameOid}
			offset = 0 //华为和智选不需要偏移
		}

		//处理通用情况如锐捷和H3C的情况
		if strings.Contains(manufacturer, Ruijie) || strings.Contains(manufacturer, H3C) {
			vlanPortListOid = []string{l2VlanPortListOid}
			portNameListOid = nil
			offset = 1 //锐捷和H3C需要偏移1位

			//H3C V3.1的情况就不使用PortListVLAN去判断了 直接使用pvid这个OID去查
			if manufacturer == H3C_V3_1 {
				vlanPortListOid = []string{}
			}
		}

		if len(vlanPortListOid) > 0 {
			chIfVlanPortList := make(chan []gosnmp.SnmpPDU)
			limitCh <- true
			go RunSnmpRetry(ip, community, timeout, chIfVlanPortList, retry, limitCh, useSnmpGetNext, vlanPortListOid)
			vlanPortListPDU = <-chIfVlanPortList
		}

		if portNameListOid != nil {
			chL2IfPortNameList := make(chan []gosnmp.SnmpPDU)
			limitCh <- true
			go RunSnmpRetry(ip, community, timeout, chL2IfPortNameList, retry, limitCh, useSnmpGetNext, portNameListOid)
			portNameListPDU = <-chL2IfPortNameList
		} else {
			portNameListPDU = ifNameList
		}

		if len(vlanPortListPDU) > 0 && len(portNameListPDU) > 0 {
			var ifIdxNameMap = make(map[int]string)
			for _, pdu := range portNameListPDU {
				idx := cast.ToInt(pdu.Name[strings.LastIndex(pdu.Name, ".")+1:])
				ifIdxNameMap[idx] = string(pdu.Value.([]byte))
			}
			var allVlanList []int

			for _, pdu := range vlanPortListPDU {
				vlanId := cast.ToInt(pdu.Name[strings.LastIndex(pdu.Name, ".")+1:])
				allVlanList = append(allVlanList, vlanId)
				hex := fmt.Sprintf("%08b", pdu.Value)
				binaryStr := strings.ReplaceAll(hex[1:len(hex)-1], " ", "")
				binary := binaryStr[:strings.LastIndex(binaryStr, "1")+1]
				//fmt.Println(vlanId, ":", binary)
				for i, s := range binary {
					if fmt.Sprintf("%c", s) == "1" {
						if ifIdxNameMap[i+offset] != "" {
							//通用的处理如：华三/锐捷要偏移1位
							//华为和智选不需要偏移
							ifPortVlanMap[ifIdxNameMap[i+offset]] = append(ifPortVlanMap[ifIdxNameMap[i+offset]], vlanId)
						}
					}
				}
			}
			//fmt.Println(ifIdxNameMap)
			//for i, v := range ifPortVlanMap {
			//	fmt.Println(i, ":", v)
			//}
			allVlan := strings.Join(findRanges(allVlanList), ",")
			for i, stats := range ifStatsList {
				if vlanIds, ok := ifPortVlanMap[stats.IfName]; ok {
					if len(vlanIds) > 0 {
						if len(vlanIds) == 1 && vlanIds[0] == 1 {
							ifStatsList[i].L2IfPortType = 3 //hybrid
							ifStatsList[i].L2IfTrunkAllowedVlan = "1-4094"
							ifStatsList[i].L2IfPVID = 1
						} else if len(vlanIds) == 1 && vlanIds[0] != 1 {
							ifStatsList[i].L2IfPortType = 2 //access
							ifStatsList[i].L2IfPVID = cast.ToUint64(vlanIds[0])
						} else {
							ifStatsList[i].L2IfPortType = 1 //trunk
							ifStatsList[i].L2IfPVID = 1
							allowedVlan := strings.Join(findRanges(vlanIds), ",")
							if allowedVlan == allVlan {
								ifStatsList[i].L2IfTrunkAllowedVlan = "all"
							} else {
								ifStatsList[i].L2IfTrunkAllowedVlan = allowedVlan
							}
						}
					}
				}
			}
		}

		//信锐交换机的处理
		if manufacturer == Sundray {
			chIfNameList_ := make(chan []gosnmp.SnmpPDU)
			chIfVlanUntaggedList := make(chan []gosnmp.SnmpPDU)
			limitCh <- true
			go RunSnmpRetry(ip, community, timeout, chIfNameList_, retry, limitCh, useSnmpGetNext, l2IfName)
			limitCh <- true
			go RunSnmpRetry(ip, community, timeout, chIfVlanUntaggedList, retry, limitCh, useSnmpGetNext, l2IfVlanUntaggedOid)
			ifNameList_ := <-chIfNameList_
			ifVlanUntaggedList := <-chIfVlanUntaggedList
			if len(ifNameList_) > 0 {
				var ifNamePortTypeMap = make(map[string]uint64)
				var ifNamePVIDMap = make(map[string]uint64)
				var ifVlanUntaggedMap = make(map[string]string)
				for i, snmpPDU := range ifNameList_ {
					ifName := string(snmpPDU.Value.([]byte))
					if t, ok := l2IfPortTypeList[i].Value.([]byte); ok {
						switch string(t) {
						case "access":
							ifNamePortTypeMap[ifName] = 2
						case "trunk":
							ifNamePortTypeMap[ifName] = 1
						case "hybrid":
							ifNamePortTypeMap[ifName] = 3
						}
					} else {
						ifNamePortTypeMap[ifName] = gosnmp.ToBigInt(l2IfPortTypeList[i].Value).Uint64()
					}
					ifNamePVIDMap[ifName] = gosnmp.ToBigInt(l2IfPVIDList[i].Value).Uint64()
					ifVlanUntaggedMap[ifName] = string(ifVlanUntaggedList[i].Value.([]byte))
				}

				for i, stats := range ifStatsList {
					ifStatsList[i].L2IfPortType = ifNamePortTypeMap[stats.IfName]
					ifStatsList[i].L2IfPVID = ifNamePVIDMap[stats.IfName]
					ifStatsList[i].L2IfTrunkAllowedVlan = ifVlanUntaggedMap[stats.IfName]
				}
			}
		}
	}
	return ifStatsList, nil
}

func ListIfOperStatus(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifOperStatusOid)
}

func ListIfName(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifNameOid)
}

func ListIfHCInOctets(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifHCInOid)
}

func ListIfHCOutOctets(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifHCOutOid)
}

func ListIfHCInUcastPkts(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifHCInPktsOid)
}

func ListIfHCInBroadcastPkts(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifHCInBroadcastPktsOid)
}

func ListIfHCOutBroadcastPkts(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifHCOutBroadcastPktsOid)
}

func ListIfHCInMulticastPkts(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifHCInMulticastPktsOid)
}

func ListIfHCOutMulticastPkts(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifHCOutMulticastPktsOid)
}

func ListIfInDiscards(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifInDiscardsOid)
}

func ListIfOutDiscards(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifOutDiscardsOid)
}

func ListIfInErrors(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifInErrorsOid)
}

func ListIfOutErrors(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifOutErrorsOid)
}

func ListIfHCOutUcastPkts(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifHCOutPktsOid)
}

func ListIfInUnknownProtos(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifInUnknownProtosOid)
}

func ListIfOutQLen(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifOutQLenOid)
}

func ListIfSpeed(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifSpeedOid)
}

func ListL2IfPortType(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, l2IfPortTypeOid)
}

func ListL2IfPVlanId(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, l2IfPvidOid)
}

func ListL2IfVlanUntagged(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, l2IfVlanUntaggedOid)
}

func ListEthernetPortMode(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ethernetPortModeOid)
}

func ListEthernetDuplex(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ethernetDuplexOid)
}

func ListIfAdminStatus(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifAdminStatusOid)
}

func ListIfDescr(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifDescrOid)
}

func ListIfStpStatus(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, dot1dStpPortStateOid)
}

func ListIfStpEnabled(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, dot1dStpPortEnabledOid)
}

func RunSnmpRetry(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool, oids []string) {
	var snmpPDUs []gosnmp.SnmpPDU
	var err error
	for _, oid := range oids {
		if useSnmpGetNext {
			snmpPDUs, err = RunSnmpGetNext(ip, community, oid, retry, timeout)
		} else {
			snmpPDUs, err = RunSnmpBulkWalk(ip, community, oid, retry, timeout)
		}
		if len(snmpPDUs) > 0 {
			err = nil
			break
		}
	}
	if err != nil {
		log.Println(ip, oids, err)
		close(ch)
		<-limitCh
		return

	}
	<-limitCh
	ch <- snmpPDUs
	return
}

func findRanges(nums []int) []string {
	var ranges []string
	if len(nums) == 0 {
		return ranges
	}

	start := nums[0]
	for i := 0; i < len(nums)-1; i++ {
		if nums[i]+1 != nums[i+1] {
			end := nums[i]
			if start == end {
				ranges = append(ranges, fmt.Sprintf("%d", start))
			} else {
				ranges = append(ranges, fmt.Sprintf("%d-%d", start, end))
			}
			start = nums[i+1]
		}
	}

	// Handle the last number
	if start == nums[len(nums)-1] {
		ranges = append(ranges, fmt.Sprintf("%d", start))
	} else {
		ranges = append(ranges, fmt.Sprintf("%d-%d", start, nums[len(nums)-1]))
	}

	return ranges
}
