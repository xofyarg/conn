package main

import (
	"strings"
	"testing"
)

func TestHostSplit(t *testing.T) {
	var u, h, p string

	u, h, p = hostSplit("")
	if h != "" {
		t.Error(h)
	}

	u, h, p = hostSplit("foo")
	if h != "foo" {
		t.Error(h)
	}

	u, h, p = hostSplit("u@foo:/bar")
	if h != "foo" {
		t.Error(h)
	}
	if u != "u@" {
		t.Error(u)
	}
	if p != ":/bar" {
		t.Error(p)
	}
}

func TestSSHParseArgs(t *testing.T) {
	var w Wrapper
	var s *sshWrapper

	w = NewWrapper(strings.Split("", " "))

	if w != nil {
		t.Error(w)
	}

	type sshw struct {
		index int
		args  string
		wargs string
	}
	cases := map[string]sshw{
		"ssh":                  {0, "ssh", ""},
		"ssh --update":         {0, "ssh", "--update"},       // ignore wrapper options
		"ssh --list foo":       {1, "ssh|foo", "--list"},     // ignore wrapper options
		"ssh -a host":          {2, "ssh|-a|host", ""},       // option a accepts no parameter
		"ssh -b param host":    {3, "ssh|-b|param|host", ""}, // option b accepts single parameter
		"ssh -bparam host":     {2, "ssh|-bparam|host", ""},
		"ssh -abparam host":    {2, "ssh|-abparam|host", ""},
		"ssh -bparam host pwd": {2, "ssh|-bparam|host|pwd", ""}, // parse remote command
	}

	for cmd, c := range cases {
		w = NewWrapper(strings.Split(cmd, " "))
		wargs := w.ParseArgs()
		s = w.(*sshWrapper)
		if s.index != c.index {
			t.Errorf("index testing [%s] => [%d != %d]", cmd, s.index, c.index)
		}

		a := strings.Join(s.args, "|")
		if a != c.args {
			t.Errorf("args testing [%s] => [%s != %s]", cmd, a, c.args)
		}

		a = strings.Join(wargs, "|")
		if a != c.wargs {
			t.Errorf("wrapper args testing [%s] => [%s != %s]", cmd, a, c.wargs)
		}
	}
}

func TestSCPParseArgs(t *testing.T) {
	var w Wrapper
	var s *scpWrapper

	w = NewWrapper(strings.Split("", " "))

	if w != nil {
		t.Error(w)
	}

	cases := map[string]int{
		"scp":                  0,
		"scp foo test:/bar":    2,
		"scp -r foo test:/bar": 3,
		"scp foo test:/bar -r": 2,
		"scp fo:o test:/bar":   1,
	}

	for cmd, index := range cases {
		w = NewWrapper(strings.Split(cmd, " "))
		w.ParseArgs()
		s = w.(*scpWrapper)
		if s.index != index {
			t.Errorf("testing [%s] => [%d != %d]", cmd, index, s.index)
		}
	}
}
