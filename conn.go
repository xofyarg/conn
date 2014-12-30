package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"time"
)

// configurable options & default value
var options = struct {
	SshBin      string
	ScpBin      string
	ConfFile    string
	CacheFile   string
	CacheExpire time.Duration
	Sep         string
	Matcher     string
	Lister      string
}{
	SshBin:      "/usr/bin/ssh",
	ScpBin:      "/usr/bin/scp",
	ConfFile:    "~/.conn.conf",
	CacheFile:   "~/.conn.cache",
	CacheExpire: 3600 * 24,
	Sep:         ".",
	Matcher:     "subtoken, token, substring, string",
	// format: name1|have_args, name2|, name3|args
	Lister: "khost|",
}

// global variables
var (
	version = "unknown"
)

// program begins
type Wrapper interface {
	// parse all the args, return wrapper args don't need later
	ParseArgs() []string
	ForceUpdate()
	Expand() []string
	Run()
}

// baseWrapper implements Wrapper interface
type baseWrapper struct {
	cmd    string   // binary to call
	args   []string // cmdline arguments
	index  int      // index for host abbrev
	prefix string   // user part in host(user@host)
	suffix string   // path part in host(host:/home)
	hosts  []string // expanded hosts
	update bool     // force update cache
	skip   bool     // bypass matching logic, use by scp on remote host
}

// ForceUpdate will make cache update unconditionaly before run
func (w *baseWrapper) ForceUpdate() {
	Debug("force cache update")
	w.update = true
}

func (w *baseWrapper) Expand() []string {
	if w.skip {
		return nil
	}

	if w.index == 0 {
		return nil
	}

	var abbrev string
	w.prefix, abbrev, w.suffix = hostSplit(w.args[w.index])

	// call each lister
	hosts := listHosts(options.Lister, w.update)

	// call matchers
	w.hosts = matchHosts(options.Matcher, abbrev, hosts)

	return w.hosts
}

func (w *baseWrapper) Run() {
	if w.hosts != nil {
		w.args[w.index] = w.prefix + w.hosts[0] + w.suffix
	}

	Debug("call: %v", w.args)
	c := &exec.Cmd{}
	c.Path = w.cmd
	c.Args = w.args
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Run()
}

type sshWrapper struct {
	baseWrapper
}

func (w *sshWrapper) ParseArgs() []string {
	var wArgs, realArgs []string
	// options which accept a argument
	chars := "bcDeFiLlmopRS"
	// jump over option arg above
	pass := false

	realArgs = append(realArgs, w.args[0])
	for i, a := range w.args[1:] {
		// All options used by ssh/scp start with single dash. So we use double
		// dashes options for wrapper.
		if strings.HasPrefix(a, "--") {
			wArgs = append(wArgs, a)
			continue
		}

		realArgs = append(realArgs, a)

		// encounter an option
		if strings.HasPrefix(a, "-") {
			for i, c := range a {
				if strings.Contains(chars, string(c)) {
					// option accepts an argument which will be provided
					// by next option
					if len(a) == i+1 {
						pass = true
					}
					break
				}
			}
			continue
		}

		if pass {
			pass = false
			continue
		}

		w.index = len(realArgs) - 1
		realArgs = append(realArgs, w.args[i+2:]...)
		break
	}

	Debug("args after parse: [a:%s] [w:%s]", realArgs, wArgs)
	// exclude wrapper args from real args
	w.args = realArgs

	return wArgs
}

type scpWrapper struct {
	baseWrapper
}

func (w *scpWrapper) ParseArgs() []string {
	var wArgs, realArgs []string

	realArgs = append(realArgs, w.args[0])
	for i, a := range w.args[1:] {
		if strings.HasPrefix(a, "--") {
			wArgs = append(wArgs, a)
			continue
		}

		if a == "-t" || a == "-f" { // skip when called remote
			w.skip = true
		}

		realArgs = append(realArgs, a)

		if strings.Contains(a, ":") {
			w.index = len(realArgs) - 1
			realArgs = append(realArgs, w.args[i+2:]...)
			w.args = realArgs
			break
		}
	}

	return wArgs
}

func main() {
	//LogLevel(LogDebug)

	// load config
	err := LoadConfig(expandPath(options.ConfFile), &options)
	if err != nil {
		Warn("cannot load config file, using default")
	}
	options.CacheExpire *= time.Second

	w := NewWrapper(os.Args)
	Debug("wapper %#v", w)

	if w == nil {
		usage := "" +
			"conn version: %s\n" +
			"Usage: conn <ssh|scp> [program specified args]\n" +
			"       or make symbolic link named <ssh|scp>\n" +
			"\n" +
			"       use `%s' to seprate host parts, e.g.:\n" +
			"         $ conn ssh baidu%swww%scom\n"
		fmt.Printf(usage, version, options.Sep, options.Sep, options.Sep)
		return
	}

	wArgs := w.ParseArgs()

	var list bool
	for _, v := range wArgs {
		switch v {
		case "--debug":
			LogLevel(LogDebug)
		case "--list":
			list = true
		case "--update":
			w.ForceUpdate()
		}
	}

	hosts := w.Expand()

	if list || len(hosts) > 1 {
		if list {
			fmt.Printf("host list:\n")
		} else {
			fmt.Printf("more than one host:\n")
		}
		fmt.Printf("  %s\n",
			strings.Join(hosts, "\n  "))
	} else {
		w.Run()
	}
}

func NewWrapper(args []string) Wrapper {
	for i, v := range args {
		// we don't need to iter over all the args
		// i == 0: ssh(symlink) foo
		// i == 1: conn ssh foo
		if i > 1 {
			break
		}

		switch path.Base(v) {
		case "ssh":
			w := &sshWrapper{baseWrapper{cmd: options.SshBin, args: args[i:]}}
			return w
		case "scp":
			w := &scpWrapper{baseWrapper{cmd: options.ScpBin, args: args[i:]}}
			return w
		}
	}

	return nil
}

// expandPath replace tilde with user home directory
func expandPath(p string) string {
	if usr, err := user.Current(); err == nil {
		dir := usr.HomeDir
		if p[:2] == "~/" {
			return strings.Replace(p, "~", dir, 1)
		}
	}
	return p
}

// hostSplit splits user@host:path into each part
func hostSplit(s string) (string, string, string) {
	var i int
	var prefix, suffix string

	i = strings.Index(s, "@")
	if i != -1 {
		prefix = s[:i+1]
		s = s[i+1:]
	}

	i = strings.Index(s, ":")
	if i != -1 {
		suffix = s[i:]
		s = s[:i]
	}

	return prefix, s, suffix
}

// load hosts list from cache or listers if it
// expires. If update is true, ignore cache.
func listHosts(ls string, update bool) []string {
	cache := expandPath(options.CacheFile)

	var hosts []string
	st, err := os.Stat(cache)
	if err == nil && !update &&
		(options.CacheExpire == 0 ||
			time.Now().Before(st.ModTime().Add(options.CacheExpire))) {
		Debug("using cache: %s", cache)
		data, _ := ioutil.ReadFile(cache)
		hosts = strings.Split(string(data), "\n")
	} else {
		Debug("building cache: %s", cache)
		hosts = list(ls)
		if hosts == nil {
			Warn("cannot get server list")
			return nil
		}
		ioutil.WriteFile(cache, []byte(strings.Join(hosts, "\n")), 0666)
	}

	return hosts
}

func list(ls string) []string {
	var hosts []string
	for _, lister := range strings.Split(ls, ",") {
		lister = strings.TrimSpace(lister)
		parts := strings.Split(lister, "|")
		if len(parts) != 2 {
			Warn("lister syntax error: %s", lister)
			continue
		}

		name := parts[0]
		arg := parts[1]
		l, ok := listers[name]
		if !ok {
			Warn("no such lister: %s", name)
		}
		h, err := l.List(arg)
		if err != nil {
			Warn("list failed: [%s] %s", name, err)
			continue
		}
		if h == nil {
			Debug("lister empty")
			continue
		}

		Debug("lister %s get %d result", name, len(h))

		hosts = append(hosts, h...)
	}
	return hosts
}

func matchHosts(ms string, pat string, hosts []string) []string {
	var matched, result []string
	for _, m := range strings.Split(ms, ",") {
		m = strings.TrimSpace(m)

		matcher, ok := matchers[m]
		if !ok {
			continue
		}

		matched = matcher.Match(pat, hosts)
		n := len(matched)
		Debug("matched %d by matcher %s", n, m)
		if n == 1 {
			result = matched
			break
		} else if n == 0 {
			continue
		} else {
			hosts = matched
			result = matched
		}
	}
	return result
}

//
// Extensions
//
type Lister interface {
	List(arg string) ([]string, error)
}

var listers map[string]Lister

func registerLister(name string, l Lister) {
	if listers == nil {
		listers = make(map[string]Lister)
	}
	listers[name] = l
}

type Matcher interface {
	Match(pat string, list []string) []string
}

var matchers map[string]Matcher

func registerMatcher(name string, m Matcher) {
	if matchers == nil {
		matchers = make(map[string]Matcher)
	}
	matchers[name] = m
}
