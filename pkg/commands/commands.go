package commands

import "gopkg.in/alecthomas/kingpin.v2"

type KingpinCommander interface {
	Command(name, help string) *kingpin.CmdClause
}
