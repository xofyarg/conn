package main

import (
	"io/ioutil"
	"net"
	"os"
	"strings"
)

type listerFunc func([]string) ([]string, error)

func (lf listerFunc) List(args []string) ([]string, error) {
	return lf(args)
}

func init() {
	registerLister("khost", listerFunc(listKnownHosts))
}

// retrieve host names from ~/.ssh/known_hosts
func listKnownHosts(noargs []string) ([]string, error) {
	f, err := os.Open(expandPath("~/.ssh/known_hosts"))
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var hosts []string
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		// ignore comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// ignore empty line
		if strings.Trim(line, " ") == "" {
			continue
		}

		// ignore hashed hostname
		if strings.HasPrefix(line, "|") {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			continue
		}

		hs := strings.Split(parts[0], ",")
		for _, h := range hs {
			// ignore ip address
			if net.ParseIP(h) != nil {
				continue
			}

			// ignore hostname with port
			if strings.Contains(h, ":") {
				continue
			}
			hosts = append(hosts, h)
		}
	}
	return hosts, nil
}
