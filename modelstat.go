package sw

import (
	"log"

	"github.com/gosnmp/gosnmp"
)

func SysModel(ip, community string, retry int, timeout int) (model string, err error) {
	var vendor string
	if v, ok := VendorMap.Load(ip); !ok {
		vendor, err = SysVendor(ip, community, retry, timeout)
		if err != nil {
			return "", err
		}
		VendorMap.Store(ip, vendor)
	} else {
		vendor = v.(string)
	}

	method := snmpGet
	var oid string

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in sw.modelstat.go SysModel", r)
		}
	}()

	switch vendor {
	case "Cisco_NX", "Cisco", "Cisco_old", "Cisco_IOS_XR", "Cisco_IOS_XE", "Ruijie":
		method = "bulkWalk"
		oid = "1.3.6.1.2.1.47.1.1.1.1.13"
	case "Huawei_ME60", "Huawei_V5", "Huawei_V3.10":
		method = "bulkWalk"
		oid = "1.3.6.1.2.1.47.1.1.1.1.2"
	case "H3C_V3.1", "H3C_S9500", "H3C", "H3C_V5", "H3C_V7", "Cisco_ASA":
		oid = "1.3.6.1.2.1.47.1.1.1.1.13"
		return getSwmodle(ip, community, oid, timeout, retry)
	case "Linux":
		return "Linux", nil
	default:
		return "", err
	}

	snmpPDUs, err := RunSnmp(ip, community, oid, method, retry, timeout)

	for _, pdu := range snmpPDUs {
		if len(string(pdu.Value.([]byte))) > 0 {
			model = model + "\n" + string(pdu.Value.([]byte))
		}
	}
	return

}

func getSwmodle(ip, community, oid string, timeout, retry int) (value string, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in getSwmodle", r)
		}
	}()
	method := "bulkWalk"
	var snmpPDUs []gosnmp.SnmpPDU
	snmpPDUs, err = RunSnmp(ip, community, oid, method, retry, timeout)
	for _, pdu := range snmpPDUs {
		if len(string(pdu.Value.([]byte))) > 0 {
			value = pdu.Value.(string)
		}
	}
	return value, err
}
