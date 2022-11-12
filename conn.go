package sw

import (
	"github.com/gosnmp/gosnmp"
	"log"
)

func ConnectionStat(ip, community string, timeout, retry int) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in Conntilization", r)
		}
	}()
	var vendor string
	var err error
	if v, ok := VendorMap.Load(ip); !ok {
		vendor, err = SysVendor(ip, community, retry, timeout)
		if err != nil {
			return 0, err
		}
		VendorMap.Store(ip, vendor)
	} else {
		vendor = v.(string)
	}

	method := snmpGet
	var oid string
	switch vendor {
	case "Cisco_ASA", "Cisco_ASA_OLD":
		oid = "1.3.6.1.4.1.9.9.147.1.2.2.2.1.5.40.6"
	default:
		return 0, err
	}

	var snmpPDUs []gosnmp.SnmpPDU
	snmpPDUs, err = RunSnmp(ip, community, oid, method, retry, timeout)
	if err == nil {
		for _, pdu := range snmpPDUs {
			return pdu.Value.(int), err
		}
	}

	return 0, err
}
