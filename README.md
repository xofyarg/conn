# conn

A wrapper of ssh/scp, use a alternative way to find a hostname which
is hard to remember.

# Overview

Conn is useful to people who has lots of hosts that is not suitable
for put them into ssh config. These hosts have long names, which is
hard to remember and type. This tool build a list of hostnames by
using some kind of _listers_ which can retrieve hostnames from
various place first. Then this list is filtered according to the user
input by several _matchers_ to find the final destination. The default
lister grabs hostname from `~/.ssh/known_hosts`.

# Installation

1. `go get go.papla.net/conn`
2. put the _conn_ binary in your $PATH
3. (optional) create a symlink point to it with the name of ssh/scp or
   make aliases:

```shell
alias ssh='conn ssh'
alias scp='conn scp'
```
	
# Usage

Insert _conn_ before calling ssh/scp, or run ssh/scp directly if you do
the step 3 of installation. Then you can use a kind of abbreviation to
specify hostname.

The abbreviation is a list of tokens separated by `.`(this can be
changed to other characters). These tokens are used to match one or
several hosts without to type the full name of the host. The order of
the tokens is not important. It should do what you think. The more
characters you type, the less results you get.

The default Lister is khost, which requires the target host should be
recorded in the _known_hosts_ file.

## Examples

- login to remote machine:

```shell
conn ssh foo.bar
```

this may match any hosts with name of foo-bar.com, bar.foo.org, or
barr.fooo.toooooo.looooong.

- copy files:

```shell
conn scp foo bar.baz:/tmp
```

will copy foo to a host matching _bar.baz_.

- list matched names:

```
conn ssh --list log
```

can list all log related hosts.

- update cache manaually before it expire:

```shell
conn ssh --update foo
```

will update the cache before match.


## Configuration

Conn reads its config file sitting at `~/.conn.conf`, you can assign
the following options in it. These are default value also.

```shell
# path to ssh binary
ssh_bin      = /usr/bin/ssh
# path to scp binary 
scp_bin      = /usr/bin/scp
# lister to use
lister       = khost|
# matchers and their order
matcher      = subtoken, token, substring, string
sep          = .
cache_file   = ~/.conn.cache
cache_expire = 86400
```

## More

The khost lister is not the first lister used by conn. If your
facility have an API call to retrieve host names, you may want to add
such a lister.

# Notes

1. conn do not parse the Host keyword in `~/.ssh/config` now.
2. use the full path to run original binary, like /usr/bin/{ssh,scp}.

