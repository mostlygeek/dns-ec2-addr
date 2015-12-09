package awsapi

import (
	"log"
	"math/rand"
	"sync"
)

type IpDB struct {
	sync.RWMutex
	list []string
}

// Update refreshes the internal list
func (i *IpDB) Update(tagApp, tagType string) {
	regions := []string{"us-west-2", "us-west-1", "us-east-1", "eu-west-1"}

	ipList := make([]string, 0)

	for _, r := range regions {
		ips, err := GetIPs(r, tagApp, tagType)
		if err != nil {
			log.Fatalf("Err pulling IPs %s, %s", r, err.Error())
		}

		ipList = append(ipList, ips...)
	}

	i.Lock()
	i.list = ipList
	i.Unlock()
}

// Get provides a random ip address or nil if nothing exists
func (i *IpDB) Get() *string {
	if len(i.list) == 0 {
		return nil
	}

	// make a copy and return the ref
	ip := i.list[rand.Intn(len(i.list)-1)]
	return &ip
}

func NewIpDB() *IpDB {
	return &IpDB{list: make([]string, 0)}
}
