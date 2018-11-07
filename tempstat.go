package sw

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/hel2o/gosnmp"
)

func Temperature(ip, community string, timeout, retry int) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in CPUtilization", r)
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

	method := "get"
	var oid string
	switch vendor {
	case "Huawei_V5":
		oid = "1.3.6.1.4.1.2011.5.25.31.1.1.1.1.11"
		return getH3CHWtemp(ip, community, oid, timeout, retry)
	default:
		err = errors.New(ip + " Switch Temperature Vendor is not defined")
		return 0, err
	}

	var snmpPDUs []gosnmp.SnmpPDU
	for i := 0; i < retry; i++ {
		snmpPDUs, err = RunSnmp(ip, community, oid, method, timeout)
		if len(snmpPDUs) > 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if err == nil {
		for _, pdu := range snmpPDUs {
			return pdu.Value.(int), err
		}
	}

	return 0, err
}

func getH3CHWtemp(ip, community, oid string, timeout, retry int) (value int, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in CPUtilization", r)
		}
	}()
	method := "getnext"
	oidnext := oid
	var snmpPDUs []gosnmp.SnmpPDU

	for {
		for i := 0; i < retry; i++ {
			snmpPDUs, err = RunSnmp(ip, community, oidnext, method, timeout)
			if len(snmpPDUs) > 0 {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		if len(snmpPDUs) < 1 {
			if err == nil {
				err = errors.New("snmpPDUs is nil")
			}
			break
		}
		oidnext = snmpPDUs[0].Name
		if strings.Contains(oidnext, oid) {
			if snmpPDUs[0].Value.(int) != 0 {
				value = snmpPDUs[0].Value.(int)
				break
			}
		} else {
			break
		}

	}
	return value, err
}
