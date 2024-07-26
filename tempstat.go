package sw

import (
	"errors"
	"log"
)

func Temperature(ip, community string, timeout, retry int) (uint64, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in Temperature", r)
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

	var oid string
	switch vendor {
	case Huawei_V5:
		oid = "1.3.6.1.4.1.2011.5.25.31.1.1.1.1.11"
	case FutureMatrix:
		oid = "1.3.6.1.4.1.56813.5.25.31.1.1.1.1.11"
	case Ruijie:
		oid = "1.3.6.1.4.1.4881.1.1.10.2.1.1.16.0"
	case H3C_V7:
		oid = "1.3.6.1.4.1.25506.2.6.1.1.1.1.12.212"
	default:
		return 0, errors.New(ip + " Switch Temperature Vendor is not defined")
	}
	return getCpuMemTemp(ip, community, oid, timeout, retry)
}
