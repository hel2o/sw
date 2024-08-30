package sw

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"strconv"
	"strings"
	"sync"
)

// VendorMap 全局记录设备型号的map
var (
	VendorMap sync.Map
)

const (
	H3C_V5        = "H3C_V5"
	H3C_V7        = "H3C_V7"
	H3C_S9500     = "H3C_S9500"
	H3C_S5500     = "H3C_S5500"
	H3C_V3_1      = "H3C_V3.1"
	H3C_ER        = "H3C_ER"
	H3C_S5024P    = "H3C_S5024P"
	H3C_S2126T    = "H3C_S2126T"
	H3C           = "H3C"
	Cisco_NX      = "Cisco_NX"
	Cisco_ASA_OLD = "Cisco_ASA_OLD"
	Cisco_ASA     = "Cisco_ASA"
	Cisco_IOS_XE  = "Cisco_IOS_XE"
	Cisco_IOS_XR  = "Cisco_IOS_XR"
	Cisco_old     = "Cisco_old"
	Cisco         = "Cisco"
	Huawei_ME60   = "Huawei_ME60"
	Huawei_V5     = "Huawei_V5"
	Huawei_V5_170 = "Huawei_V5.170" //S5720 S7700 S5730 S5735
	Huawei_V5_150 = "Huawei_V5.150" //S5700 V200R005C10SPC50
	Huawei_V5_130 = "Huawei_V5.130" //S5700 V200R003C10SPC00
	Huawei_V5_70  = "Huawei_V5.70"  //S2326 S2700
	Huawei_V3_10  = "Huawei_V3.10"
	Huawei        = "Huawei"
	Ruijie        = "Ruijie"
	Ruijie_S1900  = "Ruijie_S1900"
	Ruijie_S5700  = "Ruijie_S5700"
	Juniper       = "Juniper"
	Dell          = "Dell"
	Draytek       = "Draytek"
	FortiGate     = "FortiGate"
	Sundray       = "Sundray"
	Linux         = "Linux"
	FutureMatrix  = "FutureMatrix"
)

func SystemName(ip, community string, retry int, timeout int) (name string, err error) {
	oid := "1.3.6.1.2.1.1.5.0"
	method := snmpGet
	var snmpPDUs []gosnmp.SnmpPDU
	snmpPDUs, err = RunSnmp(ip, community, oid, method, retry, timeout)
	for _, pdu := range snmpPDUs {
		if len(string(pdu.Value.([]byte))) > 0 {
			name = string(pdu.Value.([]byte))
			break
		}
	}
	return

}

func SysDescription(ip, community string, retry int, timeout int) (desc string, err error) {
	//sysDescr
	oid := "1.3.6.1.2.1.1.1.0"
	method := snmpGet
	var snmpPDUs []gosnmp.SnmpPDU
	snmpPDUs, err = RunSnmp(ip, community, oid, method, retry, timeout)
	if err != nil {
		return
	}
	if len(snmpPDUs) == 0 || string(snmpPDUs[0].Value.([]byte)) == "" {
		//sysContact
		oid = "1.3.6.1.2.1.1.4.0"
		snmpPDUs, err = RunSnmp(ip, community, oid, method, retry, timeout)
		if err != nil {
			return
		}
	}
	for _, pdu := range snmpPDUs {
		if pdu.Value != nil && len(string(pdu.Value.([]byte))) > 0 {
			desc = desc + "\n" + string(pdu.Value.([]byte))
		}
	}
	return
}

func SerialNumber(ip, community string, retry int, timeout int) (sn string, err error) {
	oid := "1.3.6.1.2.1.47.1.1.1.1.11"
	method := snmpBulkWalk
	var snmpPDUs []gosnmp.SnmpPDU
	snmpPDUs, err = RunSnmp(ip, community, oid, method, retry, timeout)
	if err != nil {
		return
	}
	if len(snmpPDUs) == 0 {
		//信锐
		oid = ".1.3.6.1.4.1.45577.5.7.6.0"
		snmpPDUs, err = RunSnmp(ip, community, oid, method, retry, timeout)
		if err != nil {
			return
		}
	}
	for _, pdu := range snmpPDUs {
		if pdu.Value != nil && len(string(pdu.Value.([]byte))) > 0 {
			sn = string(pdu.Value.([]byte))
			break
		}
	}
	return
}

func SoftwareVersion(ip, community string, retry int, timeout int) (software string) {
	var chSnmpPDUs = make(chan []gosnmp.SnmpPDU)
	limitCh := make(chan bool, 1)
	limitCh <- true
	go RunSnmpRetry(ip, community, timeout, chSnmpPDUs, retry, limitCh, false, []string{"1.3.6.1.2.1.47.1.1.1.1.10", "1.3.6.1.4.1.45577.5.7.7.0"})
	var s []string
	snmpPDUs := <-chSnmpPDUs
	for _, pdu := range snmpPDUs {
		if pdu.Value != nil && len(string(pdu.Value.([]byte))) > 0 {
			s = append(s, string(pdu.Value.([]byte)))
		}
	}
	software = strings.Join(RemoveDuplicateElement(s), "  ")
	return
}

func ProductClass(ip, community string, retry int, timeout int) (productClass string) {
	var chSnmpPDUs = make(chan []gosnmp.SnmpPDU)
	limitCh := make(chan bool, 1)
	limitCh <- true
	go RunSnmpRetry(ip, community, timeout, chSnmpPDUs, retry, limitCh, true, []string{"1.3.6.1.2.1.47.1.1.1.1.2", "1.3.6.1.4.1.45577.5.7.8.0"})
	snmpPDUs := <-chSnmpPDUs
	for _, pdu := range snmpPDUs {
		if pdu.Value != nil && len(string(pdu.Value.([]byte))) > 3 {
			productClass = string(pdu.Value.([]byte))
			break
		}
	}
	return
}

func SysVendor(ip, community string, retry int, timeout int) (version string, sysDesc string, err error) {
	sysDesc, err = SysDescription(ip, community, retry, timeout)
	if err != nil {
		return
	}
	version, err = VersionDetect(sysDesc)
	return
}

func VersionDetect(sysDesc string) (version string, err error) {
	sysDescLower := strings.ToLower(sysDesc)
	if strings.Contains(sysDescLower, "cisco nx-os") {
		version = Cisco_NX
		return
	}
	if strings.Contains(sysDesc, "Cisco Internetwork Operating System Software") {
		version = Cisco_old
		return
	}
	if strings.Contains(sysDescLower, "cisco ios") {
		if strings.Contains(sysDesc, "IOS-XE Software") {
			version = Cisco_IOS_XE
			return
		} else if strings.Contains(sysDesc, "Cisco IOS XR") {
			version = Cisco_IOS_XR
			return
		} else {
			version = Cisco
			return
		}
	}
	if strings.Contains(sysDescLower, "cisco adaptive security appliance") {
		var versionNumber float64
		versionNumber, err = strconv.ParseFloat(getVersionNumber(sysDesc), 32)
		if err == nil && versionNumber < 9.2 {
			version = Cisco_ASA_OLD
			return
		}
		version = Cisco_ASA
		return
	}
	if strings.Contains(sysDescLower, "h3c") {
		if strings.Contains(sysDesc, "Software Version 5") {
			version = H3C_V5
			return
		}
		if strings.Contains(sysDesc, "Software Version 7") {
			version = H3C_V7
			return
		}
		if strings.Contains(sysDesc, "S5500-SI") {
			version = H3C_S5500
			return
		}
		if strings.Contains(sysDesc, "Version S9500") {
			version = H3C_S9500
			return
		}
		if strings.Contains(sysDesc, "Version 3.1") {
			version = H3C_V3_1
			return
		}
		if strings.Contains(sysDesc, "Version ER") {
			version = H3C_ER
			return
		}
		if strings.Contains(sysDesc, "S5024P") {
			version = H3C_S5024P
			return
		}
		if strings.Contains(sysDesc, "S2126T") {
			version = H3C_S2126T
			return
		}
		version = H3C
		return
	}
	if strings.Contains(sysDescLower, "futurematrix") {
		version = FutureMatrix
		return
	}
	if strings.Contains(sysDescLower, "huawei") {
		if strings.Contains(sysDesc, "Version 5.") {
			if strings.Contains(sysDesc, "Version 5.170") {
				version = Huawei_V5_170
				return
			}
			if strings.Contains(sysDesc, "Version 5.150") {
				version = Huawei_V5_150
				return
			}
			if strings.Contains(sysDesc, "Version 5.130") {
				version = Huawei_V5_130
				return
			}
			if strings.Contains(sysDesc, "Version 5.70") {
				version = Huawei_V5_70
				return
			}
			version = Huawei_V5
			return
		}
		if strings.Contains(sysDesc, "Version 3.10") {
			version = Huawei_V3_10
			return
		}
		version = Huawei
		return
	}
	if strings.Contains(sysDescLower, "ruijie") {
		if strings.Contains(sysDesc, "S57") {
			version = Ruijie_S5700
			return
		}
		if strings.Contains(sysDesc, "S19") {
			version = Ruijie_S1900
			return
		}
		version = Ruijie
		return
	}
	if strings.Contains(sysDescLower, "juniper networks") {
		version = Juniper
		return
	}
	if strings.Contains(sysDescLower, "dell networking") {
		version = Dell
		return
	}
	if strings.Contains(sysDescLower, "draytek") {
		version = Draytek
		return
	}
	if strings.Contains(sysDescLower, "fortigate") {
		version = FortiGate
		return
	}
	if strings.Contains(sysDescLower, "linux") {
		if strings.Contains(sysDescLower, "armv7l") {
			version = Sundray
			return
		}
		version = Linux
		return
	}

	err = fmt.Errorf("unknown vendor")
	return
}

func getVersionNumber(sysdescr string) string {
	version_number := ""
	s := strings.Fields(sysdescr)
	for index, value := range s {
		if strings.ToLower(value) == "version" {
			version_number = s[index+1]
		}
	}
	version_number = strings.Replace(version_number, "(", "", -1)
	version_number = strings.Replace(version_number, ")", "", -1)
	return version_number
}

func SysPatchInfo(ip, community string, retry int, timeout int) (patch string, err error) {
	oid := "1.3.6.1.4.1.2011.5.25.19.1.8.5.1.1.4"
	method := snmpBulkWalk
	var snmpPDUs []gosnmp.SnmpPDU
	snmpPDUs, err = RunSnmp(ip, community, oid, method, retry, timeout)
	for i, pdu := range snmpPDUs {
		if len(string(pdu.Value.([]byte))) > 0 {
			patch = fmt.Sprintf("%s<br>slot%d patch: %s", patch, i, string(pdu.Value.([]byte)))
		}
	}
	return
}

func RemoveDuplicateElement(s []string) []string {
	result := make([]string, 0, len(s))
	temp := map[string]struct{}{}
	for _, item := range s {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
