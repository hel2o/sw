package sw

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

func SysUpTime(ip, community string, retry, timeout int) (string, error) {
	oid := "1.3.6.1.2.1.1.3.0"
	method := snmpGet
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in Uptime", r)
		}
	}()
	snmpPDUs, err := RunSnmp(ip, community, oid, method, retry, timeout)

	if err == nil {
		for _, pdu := range snmpPDUs {
			if durationStr, err := parseTime(int(gosnmp.ToBigInt(pdu.Value).Int64())); err == nil {
				return durationStr, nil
			}

		}
	}

	return "", err
}

func parseTime(d int) (string, error) {
	timestr := strconv.Itoa(d / 100)
	duration, err := time.ParseDuration(timestr + "s")
	if err != nil {
		return "", err
	}
	totalHour := duration.Hours()
	day := int(totalHour / 24)

	modTime := math.Mod(totalHour, 24)
	modTimeStr := strconv.FormatFloat(modTime, 'f', 3, 64)
	modDuration, err := time.ParseDuration(modTimeStr + "h")
	if err != nil {
		return "", err
	}
	modDurationStr := modDuration.String()
	if strings.Contains(modDurationStr, ".") {
		modDurationStr = strings.Split(modDurationStr, ".")[0] + "s"
	}

	return fmt.Sprintf("%dday %s", day, modDurationStr), nil

}
