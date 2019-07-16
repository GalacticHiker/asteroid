package feeder

import (
	"log"
	"net"
)

// GetLocalIP returns a non loopback local IP of the host
func getLocalIP() string {

	localIP := getPreferredLocalIP()
	if localIP != nil {
		return localIP.String()
	}

	// nothing preferred - pick one of these
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	log.Fatalf("No local ip address found.")

	return ""
}

func getPreferredLocalIP() net.IP {

	ifnames := []string{"en3", "eth0"}

	for _, ifname := range ifnames {

		nif, err := net.InterfaceByName(ifname)
		if err != nil {
			return nil
		}

		addrs, _ := nif.Addrs()
		if err != nil {
			return nil
		}

		for _, addr := range addrs {
			ip, _, _ := net.ParseCIDR(addr.String())
			if ip.To4() != nil { // we want only the ipv4
				return ip
			}
		}

	}

	return nil
}
