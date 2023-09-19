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
	case "Huawei_V5":
		oid = "1.3.6.1.4.1.2011.5.25.31.1.1.1.1.11"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	case "FutureMatrix":
		oid = "1.3.6.1.4.1.56813.5.25.31.1.1.1.1.11"
		return getCpuMemTemp(ip, community, oid, timeout, retry)
	default:
		err = errors.New(ip + " Switch Temperature Vendor is not defined")
	}
	return 0, err
}
