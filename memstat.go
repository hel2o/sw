package sw

import (
	"errors"
	"github.com/gosnmp/gosnmp"
	"log"
	"strconv"
)

func MemUtilization(ip, community string, timeout, retry int) (uint64, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in MemUtilization", r)
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
	case Cisco_NX:
		oid = "1.3.6.1.4.1.9.9.305.1.1.2.0"
	case Cisco, Cisco_IOS_XE, Cisco_old:
		memUsedOid := "1.3.6.1.4.1.9.9.48.1.1.1.5.1"
		snmpMemUsed, _ := RunSnmp(ip, community, memUsedOid, method, retry, timeout)

		memFreeOid := "1.3.6.1.4.1.9.9.48.1.1.1.6.1"
		snmpMemFree, _ := RunSnmp(ip, community, memFreeOid, method, retry, timeout)

		if len(snmpMemFree) == 0 || len(snmpMemUsed) == 0 {
			err := errors.New(ip + " No Such Object available on this agent at this OID")
			return 0, err
		} else {
			if snmpMemUsed[0].Value == nil || snmpMemFree[0].Value == nil {
				err := errors.New(ip + " mem value return nil")
				return 0, err
			}
			memUsed := snmpMemUsed[0].Value.(int)
			memFree := snmpMemFree[0].Value.(int)
			if memUsed+memFree != 0 {
				memUtili := float64(memUsed) / float64(memUsed+memFree)
				return uint64(memUtili * 100), nil
			}
		}
	case Cisco_IOS_XR:
		return getCisco_IOS_XR_Mem(ip, community, timeout, retry)
	case Cisco_ASA, Cisco_ASA_OLD:
		return getCisco_ASA_Mem(ip, community, timeout, retry)
	case FutureMatrix:
		oid = "1.3.6.1.4.1.56813.5.25.31.1.1.1.1.7"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	case Huawei, Huawei_V5:
		oid = "1.3.6.1.4.1.2011.5.25.31.1.1.1.1.7"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	case Huawei_V3_10:
		return getOldHuawei_Mem(ip, community, timeout, retry)
	case Huawei_ME60:
		return getHuawei_Me60_Mem(ip, community, timeout, retry)
	case H3C_V3_1:
		oid = "1.3.6.1.4.1.2011.10.2.6.1.1.1.1.8"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	case H3C, H3C_V5, H3C_V7:
		oid = "1.3.6.1.4.1.25506.2.6.1.1.1.1.8"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	case H3C_S9500:
		oid = "1.3.6.1.4.1.2011.10.2.6.1.1.1.1.8"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	case Juniper:
		oid = "1.3.6.1.4.1.2636.3.1.13.1.11"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	case Ruijie:
		oid = "1.3.6.1.4.1.4881.1.1.10.2.35.1.1.1.3"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	case Dell:
		return GetDellMem(ip, community, timeout, retry)
	case FortiGate:
		oid = "1.3.6.1.4.1.12356.101.4.1.4"
		return getFortiGatecpumem(ip, community, oid, timeout, retry)
	default:
		err = errors.New(ip + " Switch Mem Vendor is not defined")
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

func getCisco_IOS_XR_Mem(ip, community string, timeout, retry int) (uint64, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in MemUtilization", r)
		}
	}()
	cpuindex := "1.3.6.1.4.1.9.9.109.1.1.1.1.2"
	method := snmpBulkWalk
	var snmpPDUs []gosnmp.SnmpPDU
	var err error
	var index string
	snmpPDUs, err = RunSnmp(ip, community, cpuindex, method, retry, timeout)
	index = strconv.Itoa(snmpPDUs[0].Value.(int))
	method = snmpGet
	memUsedOid := "1.3.6.1.4.1.9.9.221.1.1.1.1.18." + index + ".1"
	snmpMemUsed, _ := RunSnmp(ip, community, memUsedOid, method, retry, timeout)
	memFreeOid := "1.3.6.1.4.1.9.9.221.1.1.1.1.20." + index + ".1"
	snmpMemFree, _ := RunSnmp(ip, community, memFreeOid, method, retry, timeout)
	if len(snmpMemFree) == 0 || len(snmpMemUsed) == 0 {
		err := errors.New(ip + " No Such Object available on this agent at this OID")
		return 0, err
	} else {
		if snmpMemUsed[0].Value == nil || snmpMemFree[0].Value == nil {
			err := errors.New(ip + " mem value return nil")
			return 0, err
		}
		memUsed := snmpMemUsed[0].Value.(uint64)
		memFree := snmpMemFree[0].Value.(uint64)
		if memUsed+memFree != 0 {
			memUtili := float64(memUsed) / float64(memUsed+memFree)
			return uint64(memUtili * 100), nil
		}
	}
	return 0, err
}

func getOldHuawei_Mem(ip, community string, timeout, retry int) (uint64, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in MemUtilization", r)
		}
	}()
	method := snmpWalk
	memTotalOid := "1.3.6.1.4.1.2011.6.1.2.1.1.2"
	snmpMemTotal, err := RunSnmp(ip, community, memTotalOid, method, retry, timeout)
	memFreeOid := "1.3.6.1.4.1.2011.6.1.2.1.1.3"
	snmpMemFree, err := RunSnmp(ip, community, memFreeOid, method, retry, timeout)
	if len(snmpMemFree) == 0 || len(snmpMemTotal) == 0 {
		err := errors.New(ip + " No Such Object available on this agent at this OID")
		return 0, err
	} else {
		memTotal := snmpMemTotal[0].Value.(uint)
		memFree := snmpMemFree[0].Value.(uint)
		if memTotal != 0 {
			memUtili := float64(memTotal-memFree) / float64(memTotal)
			return uint64(memUtili * 100), err
		}
	}
	return 0, err
}

func getCisco_ASA_Mem(ip, community string, timeout, retry int) (uint64, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in MemUtilization", r)
		}
	}()
	method := snmpBulkWalk
	memUsedOid := "1.3.6.1.4.1.9.9.221.1.1.1.1.18"
	snmpMemUsed, err := RunSnmp(ip, community, memUsedOid, method, retry, timeout)
	memFreeOid := "1.3.6.1.4.1.9.9.221.1.1.1.1.20"
	snmpMemFree, err := RunSnmp(ip, community, memFreeOid, method, retry, timeout)
	if len(snmpMemFree) == 0 || len(snmpMemUsed) == 0 {
		err := errors.New(ip + " No Such Object available on this agent at this OID")
		return 0, err
	} else {
		if snmpMemUsed[0].Value == nil || snmpMemFree[0].Value == nil {
			err := errors.New(ip + " mem value return nil")
			return 0, err
		}
		memUsed := snmpMemUsed[0].Value.(uint64)
		memFree := snmpMemFree[0].Value.(uint64)
		if memUsed+memFree != 0 {
			memUtili := float64(memUsed) / float64(memUsed+memFree)
			return uint64(memUtili * 100), nil
		}
	}
	return 0, err
}

func getHuawei_Me60_Mem(ip, community string, timeout, retry int) (uint64, error) {
	memTotalOid := "1.3.6.1.4.1.2011.6.3.5.1.1.2"

	memTotal, _, err := snmp_walk_sum(ip, community, memTotalOid, timeout, retry)

	memFreeOid := "1.3.6.1.4.1.2011.6.3.5.1.1.3"
	memFree, _, err := snmp_walk_sum(ip, community, memFreeOid, timeout, retry)
	if memTotal != 0 && memFree != 0 {
		memUtili := float64(memTotal-memFree) / float64(memTotal)
		return uint64(memUtili * 100), nil
	}
	return 0, err
}

func GetDellMem(ip, community string, timeout, retry int) (uint64, error) {
	method := snmpBulkWalk
	memTotalOid := "1.3.6.1.4.1.674.10895.5000.2.6132.1.1.1.1.4.2"
	memTotal, err := RunSnmp(ip, community, memTotalOid, method, retry, timeout)
	memFreeOid := "1.3.6.1.4.1.674.10895.5000.2.6132.1.1.1.1.4.1"
	memFree, err := RunSnmp(ip, community, memFreeOid, method, retry, timeout)
	if &memTotal[0] != nil && &memFree[0] != nil {
		memfree := memFree[0].Value.(int)
		memtotal := memTotal[0].Value.(int)
		memUtili := float64(memtotal-memfree) / float64(memtotal)
		return uint64(memUtili * 100), nil
	}
	return 0, err
}
func GetLinuxMem(ip, community string, timeout, retry int) (uint64, error) {
	method := snmpGet
	memTotalOid := "1.3.6.1.4.1.2021.4.5.0"
	memTotal, err := RunSnmp(ip, community, memTotalOid, method, retry, timeout)
	memFreeOid := "1.3.6.1.4.1.2021.4.11.0"
	memFree, err := RunSnmp(ip, community, memFreeOid, method, retry, timeout)
	if &memTotal[0] != nil && &memFree[0] != nil {
		memfree := memFree[0].Value.(int)
		memtotal := memTotal[0].Value.(int)
		memUtili := float64(memtotal-memfree) / float64(memtotal)
		return uint64(memUtili * 100), nil
	}
	return 0, err
}
