package main

import (
	"io/ioutil"
	"net"
	"os"
	"strings"
)

type lister struct {
	name string
	list func(string) ([]string, error)
}

func (l *lister) List(arg string) ([]string, error) {
	return l.list(arg)
}

func init() {
	registerLister("khost", &lister{"khost", listKnownHosts})
}

// retrieve host names from ~/.ssh/known_hosts
func listKnownHosts(noarg string) ([]string, error) {
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
