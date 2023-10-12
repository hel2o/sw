package sw

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"strconv"
	"strings"
)

func SysDescription(ip, community string, retry int, timeout int) (desc string, err error) {
	oid := "1.3.6.1.2.1.1.1.0"
	method := snmpGet
	var snmpPDUs []gosnmp.SnmpPDU
	snmpPDUs, err = RunSnmp(ip, community, oid, method, retry, timeout)
	for _, pdu := range snmpPDUs {
		if len(string(pdu.Value.([]byte))) > 0 {
			desc = desc + "\n" + string(pdu.Value.([]byte))
		}
	}
	return
}

func SysVendor(ip, community string, retry int, timeout int) (string, error) {
	sysDescr, err := SysDescription(ip, community, retry, timeout)
	sysDescrLower := strings.ToLower(sysDescr)

	if strings.Contains(sysDescrLower, "cisco nx-os") {
		return "Cisco_NX", err
	}

	if strings.Contains(sysDescr, "Cisco Internetwork Operating System Software") {
		return "Cisco_old", err
	}

	if strings.Contains(sysDescrLower, "cisco ios") {
		if strings.Contains(sysDescr, "IOS-XE Software") {
			return "Cisco_IOS_XE", err
		} else if strings.Contains(sysDescr, "Cisco IOS XR") {
			return "Cisco_IOS_XR", err
		} else {
			return "Cisco", err
		}
	}

	if strings.Contains(sysDescrLower, "cisco adaptive security appliance") {
		version_number, err := strconv.ParseFloat(getVersionNumber(sysDescr), 32)
		if err == nil && version_number < 9.2 {
			return "Cisco_ASA_OLD", err
		}
		return "Cisco_ASA", err
	}
	if strings.Contains(sysDescrLower, "h3c") {
		if strings.Contains(sysDescr, "Software Version 5") {
			return "H3C_V5", err
		}

		if strings.Contains(sysDescr, "Software Version 7") {
			return "H3C_V7", err
		}

		if strings.Contains(sysDescr, "Version S9500") {
			return "H3C_S9500", err
		}

		if strings.Contains(sysDescr, "Version 3.1") {
			return "H3C_V3.1", err
		}

		if strings.Contains(sysDescr, "Version ER") {
			return "H3C_ER", err
		}

		if strings.Contains(sysDescr, "S5024P") {
			return "H3C_S5024P", err
		}

		if strings.Contains(sysDescr, "S2126T") {
			return "H3C_S2126T", err
		}

		return "H3C", err
	}
	if strings.Contains(sysDescrLower, "futurematrix") {
		return "FutureMatrix", err
	}
	if strings.Contains(sysDescrLower, "huawei") {
		if strings.Contains(sysDescr, "MultiserviceEngine 60") {
			return "Huawei_ME60", err
		}
		if strings.Contains(sysDescr, "Version 5.") {
			return "Huawei_V5", err
		}
		if strings.Contains(sysDescr, "Version 3.10") {
			return "Huawei_V3.10", err
		}
		return "Huawei", err
	}

	if strings.Contains(sysDescrLower, "ruijie") {
		return "Ruijie", err
	}

	if strings.Contains(sysDescrLower, "juniper networks") {
		return "Juniper", err
	}

	if strings.Contains(sysDescrLower, "dell networking") {
		return "Dell", err
	}
	if strings.Contains(sysDescrLower, "draytek") {
		return "Draytek", err
	}
	if strings.Contains(sysDescrLower, "fortigate") {
		return "FortiGate", err
	}
	if strings.Contains(sysDescrLower, "linux") {
		if strings.Contains(sysDescrLower, "armv7l") {
			return "Sundray", err
		}
		return "Linux", err
	}

	return "", err
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
