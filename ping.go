package sw

import "time"

func PingRtt(ip string, timeout int, retry int, fastPingMode bool) (float64, error) {
	var rtt float64
	var err error
	for i := 0; i < retry; i++ {
		if fastPingMode == true {
			rtt, err = fastPingRtt(ip, timeout)
		} else {
			rtt, err = goPingRtt(ip, timeout)
		}

		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	return rtt, err
}

func Ping(ip string, timeout int, pingRetry int, fastPingMode bool) bool {
	rtt, _ := PingRtt(ip, timeout, pingRetry, fastPingMode)
	if rtt == -1 {
		return false
	}
	return true
}
