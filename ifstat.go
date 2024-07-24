package sw

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
)

const (
	ifNameOid           = "1.3.6.1.2.1.31.1.1.1.1"
	ifNameOidPrefix     = ".1.3.6.1.2.1.31.1.1.1.1."
	ifHCInOid           = "1.3.6.1.2.1.31.1.1.1.6"
	ifHCInOidPrefix     = ".1.3.6.1.2.1.31.1.1.1.6."
	ifHCOutOid          = "1.3.6.1.2.1.31.1.1.1.10"
	ifHCInPktsOid       = "1.3.6.1.2.1.31.1.1.1.7"
	ifHCInPktsOidPrefix = ".1.3.6.1.2.1.31.1.1.1.7."
	ifHCOutPktsOid      = "1.3.6.1.2.1.31.1.1.1.11"
	//up(1),down(2)
	ifOperStatusOid              = "1.3.6.1.2.1.2.2.1.8"
	ifOperStatusOidPrefix        = ".1.3.6.1.2.1.2.2.1.8."
	ifHCInBroadcastPktsOid       = "1.3.6.1.2.1.31.1.1.1.9"
	ifHCInBroadcastPktsOidPrefix = ".1.3.6.1.2.1.31.1.1.1.9."
	ifHCOutBroadcastPktsOid      = "1.3.6.1.2.1.31.1.1.1.13"
	// multicastpkt
	ifHCInMulticastPktsOid       = "1.3.6.1.2.1.31.1.1.1.8"
	ifHCInMulticastPktsOidPrefix = ".1.3.6.1.2.1.31.1.1.1.8."
	ifHCOutMulticastPktsOid      = "1.3.6.1.2.1.31.1.1.1.12"
	// speed 配置
	ifSpeedOid       = "1.3.6.1.2.1.31.1.1.1.15"
	ifSpeedOidPrefix = ".1.3.6.1.2.1.31.1.1.1.15."

	// Discards配置
	ifInDiscardsOid       = "1.3.6.1.2.1.2.2.1.13"
	ifInDiscardsOidPrefix = ".1.3.6.1.2.1.2.2.1.13."
	ifOutDiscardsOid      = "1.3.6.1.2.1.2.2.1.19"

	// Errors配置
	ifInErrorsOid        = "1.3.6.1.2.1.2.2.1.14"
	ifInErrorsOidPrefix  = ".1.3.6.1.2.1.2.2.1.14."
	ifOutErrorsOid       = "1.3.6.1.2.1.2.2.1.20"
	ifOutErrorsOidPrefix = ".1.3.6.1.2.1.2.2.1.20."

	//ifInUnknownProtos 由于未知或不支持的网络协议而丢弃的输入报文的数量
	ifInUnknownProtosOid    = "1.3.6.1.2.1.2.2.1.15"
	ifInUnknownProtosPrefix = ".1.3.6.1.2.1.2.2.1.15."

	//ifOutQLen 接口上输出报文队列长度
	ifOutQLenOid    = "1.3.6.1.2.1.2.2.1.21"
	ifOutQLenPrefix = ".1.3.6.1.2.1.2.2.1.21."

	//二层端口类型
	//trunk(1) invalid(0) access(2) hybrid(3) fabric(4) qinq(5) desirable(6) auto(7)
	hwL2IfPortTypeOid       = "1.3.6.1.4.1.2011.5.25.42.1.1.1.3.1.3"
	hwL2IfPortTypeOidPrefix = ".1.3.6.1.4.1.2011.5.25.42.1.1.1.3.1.3."

	//二层端口的VLAN ID
	//取值范围为0～4094。如果设置为0，则hwL2IfPVID恢复为缺省值1
	hwL2IfPVIDOid       = "1.3.6.1.4.1.2011.5.25.42.1.1.1.3.1.4"
	hwL2IfPVIDOidPrefix = ".1.3.6.1.4.1.2011.5.25.42.1.1.1.3.1.4."

	//二层端口的模式
	//INTEGER : copper(1)2: fiber(2)3: other(3)
	hwEthernetPortModeOid    = "1.3.6.1.4.1.2011.5.25.157.1.1.1.1.39"
	hwEthernetPortModePrefix = ".1.3.6.1.4.1.2011.5.25.157.1.1.1.1.39."

	//以太网接口的双工模式
	//INTEGER : full(1) half(2)
	hwEthernetDuplexOid    = "1.3.6.1.4.1.2011.5.25.157.1.1.1.1.14"
	hwEthernetDuplexPrefix = ".1.3.6.1.4.1.2011.5.25.157.1.1.1.1.14."

	//是否admin down up
	//up(1),down(2),testing(3)
	ifAdminStatusOid    = "1.3.6.1.2.1.2.2.1.7"
	ifAdminStatusPrefix = ".1.3.6.1.2.1.2.2.1.7."

	//接口描述
	ifDescrOid    = "1.3.6.1.2.1.2.2.1.2"
	ifDescrPrefix = ".1.3.6.1.2.1.2.2.1.2."
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
	HwL2IfPortType       uint64 `json:"hwL2IfPortType"`
	HwL2IfPVID           uint64 `json:"hwL2IfPVID"`
	HwEthernetPortMode   uint64 `json:"hwEthernetPortMode"`
	HwEthernetDuplex     uint64 `json:"hwEthernetDuplex"`
	IfAdminStatus        uint64 `json:"ifAdminStatus"`
	IfDescr              string `json:"ifDescr"`
	TS                   int64  `json:"ts"`
}

func (this *IfStats) String() string {
	return fmt.Sprintf("<IfName:%s, IfIndex:%d, IfHCInOctets:%d, IfHCOutOctets:%d>", this.IfName, this.IfIndex, this.IfHCInOctets, this.IfHCOutOctets)
}

func ListIfStats(ip, community string, timeout int, ignoreIface []string, retry int, limitConn int, ignorePkt bool, ignoreOperStatus bool, ignoreBroadcastPkt bool, ignoreMulticastPkt bool, ignoreDiscards bool, ignoreErrors bool, ignoreUnknownProtos bool, ignoreOutQLen bool, ignoreHwL2IfPortType bool, ignoreHwL2IfPVID bool, ignoreHwEthernetPortMode bool, ignoreHwEthernetDuplex bool, ignoreIfAdminStatus bool, ignoreIfDescr bool, useSnmpGetNext bool) ([]IfStats, error) {
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

	//hwL2IfPortType
	var hwL2IfPortTypeList []gosnmp.SnmpPDU
	chHwL2IfPortTypeList := make(chan []gosnmp.SnmpPDU)
	if ignoreHwL2IfPortType == false {
		limitCh <- true
		go ListHwL2IfPortType(ip, community, timeout, chHwL2IfPortTypeList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}

	//hwL2IfPVID
	var hwL2IfPVIDList []gosnmp.SnmpPDU
	chHwL2IfPVIDList := make(chan []gosnmp.SnmpPDU)
	if ignoreHwL2IfPVID == false {
		limitCh <- true
		go ListHwL2IfPVlanId(ip, community, timeout, chHwL2IfPVIDList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}
	//hwEthernetPortModeOid
	var hwEthernetPortModeList []gosnmp.SnmpPDU
	chHwEthernetPortModeList := make(chan []gosnmp.SnmpPDU)
	if ignoreHwL2IfPortType == false {
		limitCh <- true
		go ListHwEthernetPortMode(ip, community, timeout, chHwEthernetPortModeList, retry, limitCh, useSnmpGetNext)
		time.Sleep(sleep)
	}
	//hwEthernetDuplex
	var hwEthernetDuplexList []gosnmp.SnmpPDU
	chHwEthernetDuplexList := make(chan []gosnmp.SnmpPDU)
	if ignoreHwL2IfPortType == false {
		limitCh <- true
		go ListHwEthernetDuplex(ip, community, timeout, chHwEthernetDuplexList, retry, limitCh, useSnmpGetNext)
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
	if ignoreHwL2IfPortType == false {
		hwL2IfPortTypeList = <-chHwL2IfPortTypeList
	}
	if ignoreHwL2IfPVID == false {
		hwL2IfPVIDList = <-chHwL2IfPVIDList
	}
	if ignoreHwEthernetPortMode == false {
		hwEthernetPortModeList = <-chHwEthernetPortModeList
	}
	if ignoreHwEthernetDuplex == false {
		hwEthernetDuplexList = <-chHwEthernetDuplexList
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

	if len(ifNameList) > 0 && len(ifInList) > 0 && len(ifOutList) > 0 {
		now := time.Now().Unix()

		for _, ifNamePDU := range ifNameList {
			//fmt.Printf("%+v\n", ifNamePDU.Type)
			ifName := string(ifNamePDU.Value.([]byte))

			check := true
			if len(ignoreIface) > 0 {
				for _, ignore := range ignoreIface {
					if strings.Contains(ifName, ignore) {
						check = false
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

				if ignoreHwL2IfPortType == false {
					for ti, hwL2IfPortTypePDU := range hwL2IfPortTypeList {
						if strings.Replace(hwL2IfPortTypePDU.Name, hwL2IfPortTypeOidPrefix, "", 1) == ifIndexStr {
							ifStats.HwL2IfPortType = gosnmp.ToBigInt(hwL2IfPortTypeList[ti].Value).Uint64()
							break
						}
					}
				}

				if ignoreHwL2IfPVID == false {
					for ti, hwL2IfPVIDPDU := range hwL2IfPVIDList {
						if strings.Replace(hwL2IfPVIDPDU.Name, hwL2IfPVIDOidPrefix, "", 1) == ifIndexStr {
							ifStats.HwL2IfPVID = gosnmp.ToBigInt(hwL2IfPVIDList[ti].Value).Uint64()
							break
						}
					}
				}

				if ignoreHwEthernetPortMode == false {
					for ti, hwEthernetPortModePDU := range hwEthernetPortModeList {
						if strings.Replace(hwEthernetPortModePDU.Name, hwEthernetPortModePrefix, "", 1) == ifIndexStr {
							ifStats.HwEthernetPortMode = gosnmp.ToBigInt(hwEthernetPortModeList[ti].Value).Uint64()
							break
						}
					}
				}

				if ignoreHwEthernetDuplex == false {
					for ti, hwEthernetDuplexPDU := range hwEthernetDuplexList {
						if strings.Replace(hwEthernetDuplexPDU.Name, hwEthernetDuplexPrefix, "", 1) == ifIndexStr {
							ifStats.HwEthernetDuplex = gosnmp.ToBigInt(hwEthernetDuplexList[ti].Value).Uint64()
							break
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

func ListHwL2IfPortType(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, hwL2IfPortTypeOid)
}

func ListHwL2IfPVlanId(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, hwL2IfPVIDOid)
}

func ListHwEthernetPortMode(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, hwEthernetPortModeOid)
}

func ListHwEthernetDuplex(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, hwEthernetDuplexOid)
}

func ListIfAdminStatus(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifAdminStatusOid)
}

func ListIfDescr(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool) {
	RunSnmpRetry(ip, community, timeout, ch, retry, limitCh, useSnmpGetNext, ifDescrOid)
}
func RunSnmpRetry(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, limitCh chan bool, useSnmpGetNext bool, oid string) {
	var snmpPDUs []gosnmp.SnmpPDU
	var err error
	if useSnmpGetNext {
		snmpPDUs, err = RunSnmpGetNext(ip, community, oid, retry, timeout)
	} else {
		snmpPDUs, err = RunSnmpBulkWalk(ip, community, oid, retry, timeout)
	}
	if err != nil {
		log.Println(ip, oid, err)
		close(ch)
		<-limitCh
		return

	}
	<-limitCh
	ch <- snmpPDUs
	return
}
