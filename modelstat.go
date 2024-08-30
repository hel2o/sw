package sw

import (
	"log"

	"github.com/gosnmp/gosnmp"
)

func SysModel(ip, community string, retry int, timeout int) (model string, err error) {
	var vendor string
	if v, ok := VendorMap.Load(ip); !ok {
		vendor, _, err = SysVendor(ip, community, retry, timeout)
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
			log.Println("Recovered in SysModel", r)
		}
	}()

	switch vendor {
	case Cisco_NX, Cisco, Cisco_old, Cisco_IOS_XE, Cisco_IOS_XR, Ruijie:
		method = snmpBulkWalk
		oid = "1.3.6.1.2.1.47.1.1.1.1.13"
	case Huawei_ME60, Huawei_V5, Huawei_V3_10, Huawei_V5_150, Huawei_V5_130, Huawei_V5_70:
		method = "bulkWalk"
		oid = "1.3.6.1.2.1.47.1.1.1.1.2"
	case H3C_V3_1, H3C_S9500, H3C, H3C_V5, H3C_V7, Cisco_ASA:
		oid = "1.3.6.1.2.1.47.1.1.1.1.13"
		return getModule(ip, community, oid, timeout, retry)
	case Linux:
		return Linux, nil
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

func getModule(ip, community, oid string, timeout, retry int) (value string, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in getModule", r)
		}
	}()
	method := snmpBulkWalk
	var snmpPDUs []gosnmp.SnmpPDU
	snmpPDUs, err = RunSnmp(ip, community, oid, method, retry, timeout)
	for _, pdu := range snmpPDUs {
		if len(string(pdu.Value.([]byte))) > 0 {
			value = pdu.Value.(string)
		}
	}
	return value, err
}
