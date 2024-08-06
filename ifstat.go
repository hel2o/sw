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
	ifNameOid           = []string{"1.3.6.1.2.1.31.1.1.1.1"}
	ifNameOidPrefix     = ".1.3.6.1.2.1.31.1.1.1.1."
	ifHCInOid           = []string{"1.3.6.1.2.1.31.1.1.1.6"}
	ifHCInOidPrefix     = ".1.3.6.1.2.1.31.1.1.1.6."
	ifHCOutOid          = []string{"1.3.6.1.2.1.31.1.1.1.10"}
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
	l2IfPvidOid = []string{"1.3.6.1.4.1.2011.5.25.42.1.1.1.3.1.4", "1.3.6.1.4.1.56813.5.25.42.1.1.1.3.1.4", "1.3.6.1.4.1.4881.1.1.10.2.9.1.6.1.3", "1.3.6.1.4.1.45577.5.7.9.1.6"}

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
	TS                   int64  `json:"ts"`
}

func (this *IfStats) String() string {
	return fmt.Sprintf("<IfName:%s, IfIndex:%d, IfHCInOctets:%d, IfHCOutOctets:%d>", this.IfName, this.IfIndex, this.IfHCInOctets, this.IfHCOutOctets)
}

func ListIfStats(ip, community string, timeout int, interfaces []string, retry int, limitConn int, ignorePkt bool, ignoreOperStatus bool, ignoreBroadcastPkt bool, ignoreMulticastPkt bool, ignoreDiscards bool, ignoreErrors bool, ignoreUnknownProtos bool, ignoreOutQLen bool, ignoreL2IfPortType bool, ignoreL2IfPVID bool, ignoreEthernetPortMode bool, ignoreEthernetDuplex bool, ignoreIfAdminStatus bool, ignoreIfDescr bool, useSnmpGetNext bool) ([]IfStats, error) {
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
	chIfInList := make(chan []gosnmp.SnmpPDU)
	chIfOutList := make(chan []gosnmp.SnmpPDU)

	chIfNameList := make(chan []gosnmp.SnmpPDU)
	chIfSpeedList := make(chan []gosnmp.SnmpPDU)

	limitCh <- true
	go ListIfHCInOctets(ip, community, timeout, chIfInList, retry, limitCh, useSnmpGetNext)
	time.Sleep(sleep)
	limitCh <- true
	go ListIfHCOutOctets(ip, community, timeout, chIfOutList, retry, limitCh, useSnmpGetNext)
	time.Sleep(sleep)
	limitCh <- true
	go ListIfName(ip, community, timeout, chIfNameList, retry, limitCh, useSnmpGetNext)
	time.Sleep(sleep)
	limitCh <- true
	go ListIfSpeed(ip, community, timeout, chIfSpeedList, retry, limitCh, useSnmpGetNext)
	time.Sleep(sleep)

	// OperStatus
	var ifStatusList []gosnmp.SnmpPDU
	chIfStatusList := make(chan []gosnmp.SnmpPDU)
	if ignoreOperStatus == false {
		limitCh <- true
		go ListIfOperStatus(ip, community, timeout, chIfStatusList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	//L2IfPortType
	var l2IfPortTypeList []gosnmp.SnmpPDU
	chL2IfPortTypeList := make(chan []gosnmp.SnmpPDU)
	if ignoreL2IfPortType == false {
		limitCh <- true
		go ListL2IfPortType(ip, community, timeout, chL2IfPortTypeList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	//L2IfPVID
	var l2IfPVIDList []gosnmp.SnmpPDU
	chL2IfPVIDList := make(chan []gosnmp.SnmpPDU)
	if ignoreL2IfPVID == false {
		limitCh <- true
		go ListL2IfPVlanId(ip, community, timeout, chL2IfPVIDList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}
	//ethernetPortModeOid
	var ethernetPortModeList []gosnmp.SnmpPDU
	chEthernetPortModeList := make(chan []gosnmp.SnmpPDU)
	if ignoreL2IfPortType == false {
		limitCh <- true
		go ListEthernetPortMode(ip, community, timeout, chEthernetPortModeList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}
	//EthernetDuplex
	var ethernetDuplexList []gosnmp.SnmpPDU
	chEthernetDuplexList := make(chan []gosnmp.SnmpPDU)
	if ignoreL2IfPortType == false {
		limitCh <- true
		go ListEthernetDuplex(ip, community, timeout, chEthernetDuplexList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	//AdminStatus
	var ifAdminStatusList []gosnmp.SnmpPDU
	chIfAdminStatusList := make(chan []gosnmp.SnmpPDU)
	if ignoreIfAdminStatus == false {
		limitCh <- true
		go ListIfAdminStatus(ip, community, timeout, chIfAdminStatusList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}
	//IfDescr
	var ifDescrList []gosnmp.SnmpPDU
	chIfDescrList := make(chan []gosnmp.SnmpPDU)
	if ignoreIfDescr == false {
		limitCh <- true
		go ListIfDescr(ip, community, timeout, chIfDescrList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}
	chIfInPktList := make(chan []gosnmp.SnmpPDU)
	chIfOutPktList := make(chan []gosnmp.SnmpPDU)

	var ifInPktList, ifOutPktList []gosnmp.SnmpPDU

	if ignorePkt == false {
		limitCh <- true
		go ListIfHCInUcastPkts(ip, community, timeout, chIfInPktList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
		limitCh <- true
		go ListIfHCOutUcastPkts(ip, community, timeout, chIfOutPktList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)

	}

	chIfInBroadcastPktList := make(chan []gosnmp.SnmpPDU)
	chIfOutBroadcastPktList := make(chan []gosnmp.SnmpPDU)

	var ifInBroadcastPktList, ifOutBroadcastPktList []gosnmp.SnmpPDU

	if ignoreBroadcastPkt == false {
		limitCh <- true
		go ListIfHCInBroadcastPkts(ip, community, timeout, chIfInBroadcastPktList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
		limitCh <- true
		go ListIfHCOutBroadcastPkts(ip, community, timeout, chIfOutBroadcastPktList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)

	}

	chIfInMulticastPktList := make(chan []gosnmp.SnmpPDU)
	chIfOutMulticastPktList := make(chan []gosnmp.SnmpPDU)

	var ifInMulticastPktList, ifOutMulticastPktList []gosnmp.SnmpPDU

	if ignoreMulticastPkt == false {
		limitCh <- true
		go ListIfHCInMulticastPkts(ip, community, timeout, chIfInMulticastPktList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
		limitCh <- true
		go ListIfHCOutMulticastPkts(ip, community, timeout, chIfOutMulticastPktList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)

	}

	//Discards
	chIfInDiscardsList := make(chan []gosnmp.SnmpPDU)
	chIfOutDiscardsList := make(chan []gosnmp.SnmpPDU)

	var ifInDiscardsList, ifOutDiscardsList []gosnmp.SnmpPDU

	if ignoreDiscards == false {
		limitCh <- true
		go ListIfInDiscards(ip, community, timeout, chIfInDiscardsList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
		limitCh <- true
		go ListIfOutDiscards(ip, community, timeout, chIfOutDiscardsList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)

	}

	//Errors
	chIfInErrorsList := make(chan []gosnmp.SnmpPDU)
	chIfOutErrorsList := make(chan []gosnmp.SnmpPDU)

	var ifInErrorsList, ifOutErrorsList []gosnmp.SnmpPDU

	if ignoreErrors == false {
		limitCh <- true
		go ListIfInErrors(ip, community, timeout, chIfInErrorsList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
		limitCh <- true
		go ListIfOutErrors(ip, community, timeout, chIfOutErrorsList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)

	}

	//UnknownProtos
	chIfInUnknownProtosList := make(chan []gosnmp.SnmpPDU)

	var ifInUnknownProtosList []gosnmp.SnmpPDU

	if ignoreUnknownProtos == false {
		limitCh <- true
		go ListIfInUnknownProtos(ip, community, timeout, chIfInUnknownProtosList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)

	}
	//QLen
	chIfOutQLenList := make(chan []gosnmp.SnmpPDU)

	var ifOutQLenList []gosnmp.SnmpPDU

	if ignoreOutQLen == false {
		limitCh <- true
		go ListIfOutQLen(ip, community, timeout, chIfOutQLenList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)

	}
	ifInList := <-chIfInList
	ifOutList := <-chIfOutList
	ifNameList := <-chIfNameList
	ifSpeedList := <-chIfSpeedList
	if ignoreOperStatus == false {
		ifStatusList = <-chIfStatusList
	}
	if ignoreL2IfPortType == false {
		l2IfPortTypeList = <-chL2IfPortTypeList
	}
	if ignoreL2IfPVID == false {
		l2IfPVIDList = <-chL2IfPVIDList
	}
	if ignoreEthernetPortMode == false {
		ethernetPortModeList = <-chEthernetPortModeList
	}
	if ignoreEthernetDuplex == false {
		ethernetDuplexList = <-chEthernetDuplexList
	}
	if ignoreIfAdminStatus == false {
		ifAdminStatusList = <-chIfAdminStatusList
	}
	if ignorePkt == false {
		ifInPktList = <-chIfInPktList
		ifOutPktList = <-chIfOutPktList
	}
	if ignoreIfDescr == false {
		ifDescrList = <-chIfDescrList
	}
	if ignoreBroadcastPkt == false {
		ifInBroadcastPktList = <-chIfInBroadcastPktList
		ifOutBroadcastPktList = <-chIfOutBroadcastPktList
	}
	if ignoreMulticastPkt == false {
		ifInMulticastPktList = <-chIfInMulticastPktList
		ifOutMulticastPktList = <-chIfOutMulticastPktList
	}
	if ignoreDiscards == false {
		ifInDiscardsList = <-chIfInDiscardsList
		ifOutDiscardsList = <-chIfOutDiscardsList
	}
	if ignoreErrors == false {
		ifInErrorsList = <-chIfInErrorsList
		ifOutErrorsList = <-chIfOutErrorsList
	}
	if ignoreUnknownProtos == false {
		ifInUnknownProtosList = <-chIfInUnknownProtosList
	}
	if ignoreOutQLen == false {
		ifOutQLenList = <-chIfOutQLenList
	}
	var ifCount int

	if len(ifNameList) > 0 && len(ifInList) > 0 && len(ifOutList) > 0 {
		now := time.Now().Unix()
		for _, ifNamePDU := range ifNameList {
			ifName := string(ifNamePDU.Value.([]byte))

			check := false
			if len(interfaces) > 0 {
				for _, inter := range interfaces {
					//以包含的接口名开头
					if strings.HasPrefix(strings.ToLower(ifName), strings.ToLower(inter)) {
						check = true
						ifCount++
						break
					}
				}
			}

			if check {
				var ifStats IfStats

				ifIndexStr := strings.Replace(ifNamePDU.Name, ifNameOidPrefix, "", 1)
				ifStats.IfIndex, _ = strconv.Atoi(ifIndexStr)

				for ti, ifHCInOctetsPDU := range ifInList {
					if strings.Replace(ifHCInOctetsPDU.Name, ifHCInOidPrefix, "", 1) == ifIndexStr {
						ifStats.IfHCInOctets = gosnmp.ToBigInt(ifInList[ti].Value).Uint64()
						ifStats.IfHCOutOctets = gosnmp.ToBigInt(ifOutList[ti].Value).Uint64()
						break
					}
				}
				if ignorePkt == false {
					for ti, ifHCInPktsPDU := range ifInPktList {
						if strings.Replace(ifHCInPktsPDU.Name, ifHCInPktsOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfHCOutUcastPkts = gosnmp.ToBigInt(ifOutPktList[ti].Value).Uint64()
							ifStats.IfHCInUcastPkts = gosnmp.ToBigInt(ifInPktList[ti].Value).Uint64()
							break
						}
					}
				}
				if ignoreBroadcastPkt == false {
					for ti, ifHCInBroadcastPktPDU := range ifInBroadcastPktList {
						if strings.Replace(ifHCInBroadcastPktPDU.Name, ifHCInBroadcastPktsOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfHCInBroadcastPkts = gosnmp.ToBigInt(ifInBroadcastPktList[ti].Value).Uint64()
							ifStats.IfHCOutBroadcastPkts = gosnmp.ToBigInt(ifOutBroadcastPktList[ti].Value).Uint64()
							break
						}
					}
				}
				if ignoreMulticastPkt == false {
					for ti, ifHCInMulticastPktPDU := range ifInMulticastPktList {
						if strings.Replace(ifHCInMulticastPktPDU.Name, ifHCInMulticastPktsOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfHCInMulticastPkts = gosnmp.ToBigInt(ifInMulticastPktList[ti].Value).Uint64()
							ifStats.IfHCOutMulticastPkts = gosnmp.ToBigInt(ifOutMulticastPktList[ti].Value).Uint64()
							break
						}
					}
				}

				if ignoreDiscards == false {
					for ti, ifInDiscardsPDU := range ifInDiscardsList {
						if strings.Replace(ifInDiscardsPDU.Name, ifInDiscardsOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfInDiscards = gosnmp.ToBigInt(ifInDiscardsList[ti].Value).Uint64()
							ifStats.IfOutDiscards = gosnmp.ToBigInt(ifOutDiscardsList[ti].Value).Uint64()
							break
						}
					}
				}

				if ignoreErrors == false {
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

				if ignoreOperStatus == false {
					for ti, ifOperStatusPDU := range ifStatusList {
						if strings.Replace(ifOperStatusPDU.Name, ifOperStatusOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfOperStatus = int(gosnmp.ToBigInt(ifStatusList[ti].Value).Int64())
							break
						}
					}
				}

				//if ignoreL2IfPortType == false {
				//	for ti, l2IfPortTypePDU := range l2IfPortTypeList {
				//		for _, oidPrefix := range l2IfPortTypeOidPrefix {
				//			if strings.Replace(l2IfPortTypePDU.Name, oidPrefix, "", 1) == ifIndexStr {
				//				ifStats.L2IfPortType = gosnmp.ToBigInt(l2IfPortTypeList[ti].Value).Uint64()
				//				break
				//			}
				//
				//		}
				//	}
				//}

				//if ignoreL2IfPVID == false {
				//	for ti, l2IfPVIDPDU := range l2IfPVIDList {
				//		for _, oidPrefix := range l2IfPVIDOidPrefix {
				//			if strings.Replace(l2IfPVIDPDU.Name, oidPrefix, "", 1) == ifIndexStr {
				//				ifStats.L2IfPVID = gosnmp.ToBigInt(l2IfPVIDList[ti].Value).Uint64()
				//				break
				//			}
				//		}
				//	}
				//}

				if ignoreEthernetPortMode == false {
					for ti, ethernetPortModePDU := range ethernetPortModeList {
						for _, oidPrefix := range ethernetPortModePrefix {
							if strings.Replace(ethernetPortModePDU.Name, oidPrefix, "", 1) == ifIndexStr {
								ifStats.EthernetPortMode = gosnmp.ToBigInt(ethernetPortModeList[ti].Value).Uint64()
								break
							}
						}
					}
				}

				if ignoreEthernetDuplex == false {
					for ti, ethernetDuplexPDU := range ethernetDuplexList {
						for _, oidPrefix := range ethernetDuplexPrefix {
							if strings.Replace(ethernetDuplexPDU.Name, oidPrefix, "", 1) == ifIndexStr {
								ifStats.EthernetDuplex = gosnmp.ToBigInt(ethernetDuplexList[ti].Value).Uint64()
								break
							}
						}
					}
				}

				if ignoreIfAdminStatus == false {
					for ti, ifAdminStatusPDU := range ifAdminStatusList {
						if strings.Replace(ifAdminStatusPDU.Name, ifAdminStatusPrefix, "", 1) == ifIndexStr {
							ifStats.IfAdminStatus = gosnmp.ToBigInt(ifAdminStatusList[ti].Value).Uint64()
							break
						}
					}
				}

				if ignoreIfDescr == false {
					for ti, ifDescrPDU := range ifDescrList {
						if strings.Replace(ifDescrPDU.Name, ifDescrPrefix, "", 1) == ifIndexStr {
							ifStats.IfDescr = string(ifDescrList[ti].Value.([]byte))
							break
						}
					}
				}

				if ignoreUnknownProtos == false {
					for ti, ifInUnknownProtosPDU := range ifInUnknownProtosList {
						if strings.Replace(ifInUnknownProtosPDU.Name, ifInUnknownProtosPrefix, "", 1) == ifIndexStr {
							ifStats.IfInUnknownProtos = gosnmp.ToBigInt(ifInUnknownProtosList[ti].Value).Uint64()
							break
						}
					}
				}

				if ignoreOutQLen == false {
					for ti, ifOutQLenPDU := range ifOutQLenList {
						if strings.Replace(ifOutQLenPDU.Name, ifOutQLenPrefix, "", 1) == ifIndexStr {
							ifStats.IfOutQLen = gosnmp.ToBigInt(ifOutQLenList[ti].Value).Uint64()
							break
						}
					}
				}

				for ti, ifSpeedPDU := range ifSpeedList {
					if strings.Replace(ifSpeedPDU.Name, ifSpeedOidPrefix, "", 1) == ifIndexStr {
						ifStats.IfSpeed = 1000 * 1000 * gosnmp.ToBigInt(ifSpeedList[ti].Value).Uint64()
						break
					}
				}

				ifStats.TS = now
				ifStats.IfName = ifName
				ifStatsList = append(ifStatsList, ifStats)
			}
		}
	}
	//以下是处理了锐捷 信锐采集接口类型和PVID的情况
	if ignoreL2IfPortType == false || ignoreL2IfPVID == false {
		chIfNameList_ := make(chan []gosnmp.SnmpPDU)
		limitCh <- true
		go ListIfName_(ip, community, timeout, chIfNameList_, retry, limitCh, useSnmpGetNext)
		ifNameList_ := <-chIfNameList_

		var ifNamePortTypeMap = make(map[string]uint64)
		var ifNamePVIDMap = make(map[string]uint64)
		if len(ifNameList_) == 0 {
			if len(l2IfPVIDList) > 0 && len(l2IfPortTypeList) > 0 {
				ifNameList_ = ifNameList
				for _, pdu := range l2IfPortTypeList {
					start := strings.LastIndex(pdu.Name, ".")
					ifNamePortTypeMap[pdu.Name[start+1:]] = gosnmp.ToBigInt(pdu.Value).Uint64()
				}
				for _, pdu := range l2IfPVIDList {
					start := strings.LastIndex(pdu.Name, ".")
					ifNamePVIDMap[pdu.Name[start+1:]] = gosnmp.ToBigInt(pdu.Value).Uint64()
				}
				for i, stats := range ifStatsList {
					pt, ok1 := ifNamePortTypeMap[fmt.Sprintf("%d", stats.IfIndex)]
					if ok1 {
						ifStatsList[i].L2IfPortType = pt
					}
					pvid, ok2 := ifNamePVIDMap[fmt.Sprintf("%d", stats.IfIndex)]
					if ok2 {
						ifStatsList[i].L2IfPVID = pvid
					}
				}
				//	新华三的处理
			} else if len(l2IfPVIDList) == 0 {
				limitCh <- true
				chIfVlanPortList := make(chan []gosnmp.SnmpPDU)
				go ListIfVlanPortList(ip, community, timeout, chIfVlanPortList, retry, limitCh, useSnmpGetNext)
				ifVlanPortList := <-chIfVlanPortList
				if len(ifVlanPortList) > 0 {
					var ifPortVlanMap = make(map[string][]string)
					for _, pdu := range ifVlanPortList {
						start := strings.LastIndex(pdu.Name, ".")
						vlanId := pdu.Name[start+1:]
						for _, port := range strings.Split(string(pdu.Value.([]byte)), ",") {
							if port != "" {
								ifPortVlanMap[port] = append(ifPortVlanMap[port], vlanId)
							}
						}
					}
					for i, stats := range ifStatsList {
						if vlanIds, ok := ifPortVlanMap[fmt.Sprintf("%d", stats.IfIndex)]; ok {
							if len(vlanIds) > 0 {
								if len(vlanIds) == 1 && vlanIds[0] == "1" {
									ifStatsList[i].L2IfPortType = 3 //hybrid
									ifStatsList[i].L2IfTrunkAllowedVlan = "1-4094"
									ifStatsList[i].L2IfPVID = 1
								} else if len(vlanIds) == 1 && vlanIds[0] != "1" {
									ifStatsList[i].L2IfPortType = 2 //access
									ifStatsList[i].L2IfPVID = cast.ToUint64(vlanIds[0])
								} else {
									ifStatsList[i].L2IfPortType = 1 //trunk
									ifStatsList[i].L2IfPVID = 1
									ifStatsList[i].L2IfTrunkAllowedVlan = strings.Join(vlanIds, ",")
								}

							}
						}
					}
				}
			}
		} else {
			for i, snmpPDU := range ifNameList_ {
				//信锐的处理
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
			}
			for i, stats := range ifStatsList {
				ifStatsList[i].L2IfPortType = ifNamePortTypeMap[stats.IfName]
				ifStatsList[i].L2IfPVID = ifNamePVIDMap[stats.IfName]
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

func ListIfName_(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, l2IfName)
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

func ListIfVlanPortList(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifVlanPortListOid)
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
