package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type IPv4 struct {
	DefaultGateway string `json:"DefaultGateway"`
	DhcpEnable     bool   `json:"DhcpEnable"`
	IPAddress      string `json:"IPAddress"`
	SubnetMask     string `json:"SubnetMask"`
}

type IPv6 struct {
	DefaultGateway   string `json:"DefaultGateway"`
	DhcpEnable       bool   `json:"DhcpEnable"`
	IPAddress        string `json:"IPAddress"`
	LinkLocalAddress string `json:"LinkLocalAddress"`
}

type DeviceInfo struct {
	AlarmInputChannels       int    `json:"AlarmInputChannels"`
	AlarmOutputChannels      int    `json:"AlarmOutputChannels"`
	DeviceClass              string `json:"DeviceClass"`
	DeviceType               string `json:"DeviceType"`
	HttpPort                 int    `json:"HttpPort"`
	IPv4Address              IPv4   `json:"IPv4Address"`
	IPv6Address              IPv6   `json:"IPv6Address"`
	MachineName              string `json:"MachineName"`
	Manufacturer             string `json:"Manufacturer"`
	Port                     int    `json:"Port"`
	RemoteVideoInputChannels int    `json:"RemoteVideoInputChannels"`
	SerialNo                 string `json:"SerialNo"`
	Vendor                   string `json:"Vendor"`
	Version                  string `json:"Version"`
	VideoInputChannels       int    `json:"VideoInputChannels"`
	VideoOutputChannels      int    `json:"VideoOutputChannels"`
}

type Response struct {
	Mac    string `json:"mac"`
	Method string `json:"method"`
	Params struct {
		DeviceInfo DeviceInfo `json:"deviceInfo"`
	} `json:"params"`
}

func main() {
	conn, err := net.ListenPacket("udp", "0.0.0.0:37810")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("[DVR] Listening for connections ...")
	for {
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFrom((buf))
		if err != nil {
			log.Fatal(err)
			continue
		}

		fmt.Printf("%s 0x44, 0x48, 0x49, 0x50\n", addr.String())
		if n == 4 && buf[0] == 0x44 && buf[1] == 0x48 && buf[2] == 0x49 && buf[3] == 0x50 {

			// create response payload
			response := Response{
				Mac:    "80:91:33:93:28:5f",
				Method: "client.notifyDevInfo",
			}
			response.Params.DeviceInfo = DeviceInfo{
				AlarmInputChannels:  0,
				AlarmOutputChannels: 0,
				DeviceClass:         "IPC",
				DeviceType:          "IPC-HDW5421S",
				HttpPort:            443,
				IPv4Address: IPv4{
					DefaultGateway: "93.95.227.1",
					DhcpEnable:     false,
					IPAddress:      "93.95.228.103",
					SubnetMask:     "255.255.255.0",
				},
				IPv6Address: IPv6{
					DefaultGateway:   "2001:db8:1234:ffff:ffff:ffff:ffff:ffff",
					DhcpEnable:       false,
					IPAddress:        "2001:db8:1234::/48",
					LinkLocalAddress: "fe80::4e11:bfff:fec2:0b9d/64",
				},
				MachineName:              "1F006E4PAX00075",
				Manufacturer:             "Private",
				Port:                     37777,
				RemoteVideoInputChannels: 0,
				SerialNo:                 "1F006E4PAX00075",
				Vendor:                   "Dahua",
				Version:                  "2.400.0.8",
				VideoInputChannels:       1,
				VideoOutputChannels:      16,
			}

			responseData, err := json.Marshal(response)
			if err != nil {
				log.Println("JSON error: ", err)
				continue
			}

			// send response as 800 bytes
			pad := make([]byte, 800)
			copy(pad, responseData)

			_, err = conn.WriteTo(pad, addr)
			if err != nil {
				log.Println("WriteTo error: err")
			}
		}
	}

}
