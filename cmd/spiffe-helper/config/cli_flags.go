package config

import (
	"flag"
	"fmt"
	"strconv"
)

const (
	daemonModeFlagName = "daemon-mode"
)

type CLIFlags interface {
	Parse() error
	ConfigFile() string
	DaemonModeFlag() (bool, bool)
	PositionalCommandFlag() bool
}

type cliFlagsParser struct {
	configFile        string
	daemonMode        *bool
	positionalCommand bool
}

func NewCLIFlags() CLIFlags {
	f := cliFlagsParser{}
	flag.StringVar(&f.configFile, "config-file", "<configFile> Configuration file path", "helper.conf")
	// We need to be able to tell if -daemon-mode was set or not to decide
	// whether to override the config file value
	flag.BoolFunc(
		daemonModeFlagName,
		"<true|false> Toggle running as a daemon to rotate X.509/JWT or just fetch and exit",
		func(s string) error {
			if b, e := strconv.ParseBool(s); e != nil {
				return fmt.Errorf("parsing -daemon-mode=%s as bool: %w", s, e)
			} else {
				f.daemonMode = &b
			}
			return nil
		},
	)
	flag.BoolVar(
		&f.positionalCommand,
		"command",
		false,
		"<true|false> Positional arguments after -- will be run as a command and its arguments, overriding the cmd and cmd_args in the configuration file",
	)
	return &f
}

func (f *cliFlagsParser) Parse() error {
	flag.Parse()
	return nil
}

func (f *cliFlagsParser) ConfigFile() string {
	return f.configFile
}

func (f *cliFlagsParser) DaemonModeFlag() (bool, bool) {
	if f.daemonMode == nil {
		return false, false
	}
	return true, *f.daemonMode
}

func (f *cliFlagsParser) PositionalCommandFlag() bool {
	return f.positionalCommand
}
