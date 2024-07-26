package sw

import (
	"log"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
)

const snmpGet = "get"
const snmpGetNex = "getNext"
const snmpWalk = "walk"
const snmpBulkWalk = "bulkWalk"

func RunSnmp(ip, community, oid, method string, retry, timeoutMillisecond int) (snmpPDUs []gosnmp.SnmpPDU, err error) {
	if method == snmpGetNex {
		timeoutMillisecond = 5000
	}
	params := &gosnmp.GoSNMP{
		Target:    ip,
		Port:      161,
		Version:   gosnmp.Version2c,
		Community: community,
		Timeout:   time.Duration(timeoutMillisecond) * time.Millisecond,
	}
	params.Retries = retry
	err = params.Connect()
	if err != nil {
		return
	}
	defer params.Conn.Close()
	snmpPDUs, err = ParseSnmpMethod(oid, method, params)
	return
}

func ParseSnmpMethod(oid, method string, cur_gosnmp *gosnmp.GoSNMP) (snmpPDUs []gosnmp.SnmpPDU, err error) {
	var snmpPacket *gosnmp.SnmpPacket
	defer func() {
		if r := recover(); r != nil {
			log.Println(cur_gosnmp.Target+" Recovered in ParseSnmpMethod", r)
		}
	}()
	switch method {
	case snmpGet:
		snmpPacket, err = cur_gosnmp.Get([]string{oid})
		if err != nil {
			return nil, err
		} else {
			snmpPDUs = snmpPacket.Variables
			return snmpPDUs, err
		}
	case snmpGetNex:
		var oidNext = oid
		var pack *gosnmp.SnmpPacket
		for {
			pack, err = cur_gosnmp.GetNext([]string{oidNext})
			if err != nil || len(pack.Variables) <= 0 {
				break
			}
			oidNext = pack.Variables[0].Name
			if strings.Contains(oidNext, oid) {
				snmpPDUs = append(snmpPDUs, pack.Variables[0])
			} else {
				break
			}
		}

	case snmpBulkWalk:
		err = cur_gosnmp.BulkWalk(oid, func(pdu gosnmp.SnmpPDU) error {
			snmpPDUs = append(snmpPDUs, pdu)
			return nil
		})
		if err != nil {
			return nil, err
		}
	default:
		err = cur_gosnmp.Walk(oid, func(pdu gosnmp.SnmpPDU) error {
			snmpPDUs = append(snmpPDUs, pdu)
			return nil
		})
	}

	return
}

func RunSnmpBulkWalk(ip, community, oid string, retry int, timeout int) ([]gosnmp.SnmpPDU, error) {
	snmpPDUs, err := RunSnmp(ip, community, oid, snmpBulkWalk, retry, timeout)
	return snmpPDUs, err
}

func RunSnmpGetNext(ip, community, oid string, retry int, timeout int) ([]gosnmp.SnmpPDU, error) {
	snmpPDUs, err := RunSnmp(ip, community, oid, snmpGetNex, retry, timeout)
	return snmpPDUs, err
}

func RunSnmpGet(ip, community, oid string, retry int, timeout int) ([]gosnmp.SnmpPDU, error) {
	snmpPDUs, err := RunSnmp(ip, community, oid, snmpGet, retry, timeout)
	return snmpPDUs, err
}
