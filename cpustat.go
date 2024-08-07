package sw

import (
	"errors"
	"github.com/gosnmp/gosnmp"
	"log"
	"sync"
)

var (
	VendorMap sync.Map
)

func CpuUtilization(ip, community string, timeout, retry int) (uint64, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in CPUtilization", r)
		}
	}()
	var vendor string
	var err error
	if v, ok := VendorMap.Load(ip); !ok {
		vendor, _, err = SysVendor(ip, community, retry, timeout)
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
	case Cisco_NX:
		oid = "1.3.6.1.4.1.9.9.305.1.1.1.0"
	case Cisco, Cisco_old:
		oid = "1.3.6.1.4.1.9.9.109.1.1.1.1.7.1"
	case Cisco_IOS_XE, Cisco_IOS_XR:
		oid = "1.3.6.1.4.1.9.9.109.1.1.1.1.7"
		method = "bulkWalk"
	case Cisco_ASA:
		oid = "1.3.6.1.4.1.9.9.109.1.1.1.1.7"
		return getCiscoASAcpu(ip, community, oid, timeout, retry)
	case Cisco_ASA_OLD:
		oid = "1.3.6.1.4.1.9.9.109.1.1.1.1.4"
		return getCiscoASAcpu(ip, community, oid, timeout, retry)
	case FutureMatrix:
		oid = "1.3.6.1.4.1.56813.6.3.4.1.3"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	case Huawei, Huawei_V5:
		oid = "1.3.6.1.4.1.2011.5.25.31.1.1.1.1.5"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	case Huawei_V3_10, H3C_V3_1:
		oid = "1.3.6.1.4.1.2011.6.1.1.1.3"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	case Huawei_ME60:
		oid = "1.3.6.1.4.1.2011.6.3.4.1.2"
		return getHuawei_ME60cpu(ip, community, oid, timeout, retry)
	case H3C, H3C_V5, H3C_V7:
		oid = "1.3.6.1.4.1.25506.2.6.1.1.1.1.6"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	case H3C_ER:
		oid = "1.3.6.1.2.1.25.3.3.1.2"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	case H3C_S9500, H3C_S5500:
		oid = "1.3.6.1.4.1.2011.10.2.6.1.1.1.1.6"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	case Juniper:
		oid = "1.3.6.1.4.1.2636.3.1.13.1.8"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	case Ruijie:
		oid = "1.3.6.1.4.1.4881.1.1.10.2.36.1.1.2"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	case Dell:
		oid = "1.3.6.1.4.1.674.10895.5000.2.6132.1.1.1.1.4.11"
		return getDellCpu(ip, community, oid, timeout, retry)
	case FortiGate:
		oid = "1.3.6.1.4.1.12356.101.4.1.3"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	case Sundray:
		oid = "1.3.6.1.4.1.2021.11.11.0"
		return getSundrayCpu(ip, community, oid, timeout, retry)
	default:
		err = errors.New(ip + " Switch Cpu Vendor is not defined")
		return 0, err
	}

	var snmpPDUs []gosnmp.SnmpPDU

	snmpPDUs, err = RunSnmp(ip, community, oid, method, retry, timeout)

	if err == nil {
		for _, pdu := range snmpPDUs {
			return pdu.Value.(uint64), err
		}
	}

	return 0, err
}

func getCiscoASAcpu(ip, community, oid string, timeout, retry int) (value uint64, err error) {
	CPU_Value_SUM, CPU_Count, err := snmp_walk_sum(ip, community, oid, timeout, retry)
	if err == nil {
		if CPU_Count > 0 {
			return uint64(CPU_Value_SUM / CPU_Count), err
		}
	}
	return 0, err
}

func getCpuMemTemp(ip, community, oid string, timeout, retry int) (value uint64, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in CpuMemTemp", r, oid)
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
	var valid uint64
	for _, pdu := range snmpPDUs {
		v := gosnmp.ToBigInt(pdu.Value).Uint64()
		if v > 0 && v <= 120 {
			valid++
			value = value + v
		}
	}
	if valid > 0 {
		value = value / valid
	}
	return
}

func getSundrayCpu(ip, community, oid string, timeout, retry int) (uint64, error) {
	ssCpuIdle, err := getCpuMemTemp(ip, community, oid, timeout, retry)
	if err != nil {
		return 0, err
	}
	return 100 - ssCpuIdle, err
}

func getFortiGatecpumem(ip, community, oid string, timeout, retry int) (value uint64, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in CPUtilization", r, oid)
		}
	}()
	method := snmpBulkWalk

	var snmpPDUs []gosnmp.SnmpPDU

	snmpPDUs, err = RunSnmp(ip, community, oid, method, retry, timeout)

	return snmpPDUs[0].Value.(uint64), err
}

func getHuawei_ME60cpu(ip, community, oid string, timeout, retry int) (value uint64, err error) {
	CPU_Value_SUM, CPU_Count, err := snmp_walk_sum(ip, community, oid, timeout, retry)
	if err == nil {
		if CPU_Count > 0 {
			return uint64(CPU_Value_SUM / CPU_Count), err
		}
	}

	return 0, err
}

func getDellCpu(ip, community, oid string, timeout, retry int) (value uint64, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in CPUtilization", r)
		}
	}()
	method := snmpBulkWalk

	var snmpPDUs []gosnmp.SnmpPDU

	snmpPDUs, err = RunSnmp(ip, community, oid, method, retry, timeout)

	return snmpPDUs[0].Value.(uint64), err
}

func snmp_walk_sum(ip, community, oid string, timeout, retry int) (value_sum int, value_count int, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in CPUtilization", r)
		}
	}()
	var snmpPDUs []gosnmp.SnmpPDU
	method := snmpWalk

	snmpPDUs, err = RunSnmp(ip, community, oid, method, retry, timeout)

	var Values []int
	if err == nil {
		for _, pdu := range snmpPDUs {
			Values = append(Values, pdu.Value.(int))
		}
	}
	var Value_SUM int
	Value_SUM = 0
	for _, value := range Values {
		Value_SUM = Value_SUM + value
	}
	return Value_SUM, len(Values), err
}
