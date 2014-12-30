# conn: convenient way to use ssh/scp in command line
conn is a wrapper of ssh/scp, use a alternative way to find a hostname without looking up the records in `~/.ssh/config`.

conn use some kind of listers which can retrieve hostnames from various place. The default lister grabs hostname from `~/.ssh/known_hosts`.

# Usage:
Insert conn before ssh/scp :)

Then you can use a kind of abbreviation to specify hostname. After add the following lines to `~/.bashrc`, you may forget existence of conn.

```Shell
    alias ssh='conn ssh'
    alias scp='conn scp'
```

If you need to use original ssh/scp, please use full path, like /usr/bin/{ssh,scp}.


# Configuration:
conn reads its config file sitting at `~/.conn.conf`, you can assign the following options in it. These are default value also.

```Shell
# ssh original path
ssh_bin:      = "/usr/bin/ssh"
# scp original path
scp_bin:      = "/usr/bin/scp"
# lister to use
lister = khost|
```

If you need more, look at the source code, there are more options such as matcher, cache that you can tune.

# Examples:
0. basic matching rule
foo.bar will be split to two keywords, foo and bar, which will be used to match the hostnames in the list. If it cannot do what you think, please contact the author.

1. login to remote machine
`conn ssh foo.bar` may match any hosts with name of foo-bar.com, bar.foo.org, or barr.fooo.toooooo.looooong

2. copy files
`conn scp foo bar.baz:/tmp` copy foo to a host matching *bar.baz*.

3. list matched names
`conn ssh --list syslog`: list all syslog related hosts.

# Notes:
1. Do not parse the Host keyword in `~/.ssh/config` now.


# Bug & Feature request:
  Please create an issue.
