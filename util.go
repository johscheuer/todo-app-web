package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

func getAllAddresses(ifaces []net.Interface) ([]string, error) {
	var addresses []string
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			return addresses, err
		}

		for _, addr := range addrs {
			addresses = append(addresses, addr.String())
		}
	}

	return addresses, nil
}

func generateJSONResponse(rw http.ResponseWriter, toMarshal interface{}) {
	responseJSON, err := json.MarshalIndent(toMarshal, "", "  ")
	if err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), 500)
	}
	rw.Write(responseJSON)
}
