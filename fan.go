package sw

import (
	"errors"
	"github.com/gosnmp/gosnmp"
	"log"
)

type SlotValue struct {
	Value uint64
	Slot  uint64
}

func FanStatus(ip, community string, timeout, retry int) ([]SlotValue, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in Temperature", r)
		}
	}()
	var vendor string
	var err error
	if v, ok := VendorMap.Load(ip); !ok {
		vendor, _, err = SysVendor(ip, community, retry, timeout)
		if err != nil {
			return nil, err
		}
		VendorMap.Store(ip, vendor)
	} else {
		vendor = v.(string)
	}

	var oid string
	switch vendor {
	case Huawei_V5_130, Huawei_V5_150, Huawei_V5_170:
		oid = "1.3.6.1.4.1.2011.5.25.31.1.1.10.1.7"
	case FutureMatrix:
		oid = "1.3.6.1.4.1.56813.5.25.31.1.1.10.1.7"
	default:
		return nil, errors.New(ip + " Switch Temperature Vendor is not defined")
	}
	return generalGetValue(ip, community, oid, timeout, retry)
}

// generalGetValue 返回一个数组，数组的每个元素是一个 SlotValue 结构体，包含了每个槽的状态值和槽位号
func generalGetValue(ip, community, oid string, timeout, retry int) (sv []SlotValue, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in commGetValue", r, oid)
		}
	}()
	method := snmpBulkWalk
	var snmpPDUs []gosnmp.SnmpPDU
	snmpPDUs, err = RunSnmp(ip, community, oid, method, retry, timeout)
	if len(snmpPDUs) < 1 {
		if err == nil {
			err = errors.New("snmpPDUs is nil")
		}
	}
	for i, pdu := range snmpPDUs {
		v := gosnmp.ToBigInt(pdu.Value).Uint64()
		sv = append(sv, SlotValue{Value: v, Slot: uint64(i)})
	}
	return
}
