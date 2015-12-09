package main

import (
	"fmt"

	"github.com/mostlygeek/dns-ec2-addr/awsapi"
)

func main() {

	ipdb := awsapi.NewIpDB()

	fmt.Println("Fetching public ip addresses...")
	ipdb.Update("autopush", "autopush")

	for i := 0; i < 25; i++ {
		if ip := ipdb.Get(); ip != nil {
			fmt.Println(*ip)
		} else {
			fmt.Println("Err: nil")
			break
		}
	}
}
