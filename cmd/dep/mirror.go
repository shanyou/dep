package main

import (
	"flag"
	"strings"

	"github.com/golang/dep"
	"github.com/golang/dep/gps/mirrors"
	"github.com/pkg/errors"
)

const mirrorShortHelp = `mirror config for vendor project`
const mirrorLongHelp = `
mirror command try to help developer to manage vendor with mirror reporoot. it will init an cache
file in "$HOME/.dep/mirror.yaml". such as
repos:
- prefix: golang.org/x/sys
  repo: https://github.com/golang/sys.git
  vcs: git
- prefix: cloud.google.com/go
  repo: https://github.com/googleapis/google-cloud-go.git
  vcs: git

when use ensure command, find package with prefix, dep will use given mirror repo to retrive code.
`
const mirrorExamples = `
dep mirror

	try to mange mirror for vendor package.

dep mirror -add -p k8s.io/apimachinery -r https://github.com/kubernetes/apimachinery.git -s git

	it will add k8s.io/apimachinery mirrors
`

type mirrorCommand struct {
	examples bool
	add      bool
	remove   bool
	list     bool
	prefix   string
	repo     string
	vcs      string
}

func (cmd *mirrorCommand) Name() string { return "mirror" }
func (cmd *mirrorCommand) Args() string {
	return "[-add | -remove | -list ] [-p <prefix> -r <repo> -v <vcs>]"
}
func (cmd *mirrorCommand) ShortHelp() string { return mirrorShortHelp }
func (cmd *mirrorCommand) LongHelp() string  { return mirrorLongHelp }
func (cmd *mirrorCommand) Hidden() bool      { return false }

func (cmd *mirrorCommand) Register(fs *flag.FlagSet) {
	fs.BoolVar(&cmd.examples, "examples", false, "print detailed usage examples")
	fs.BoolVar(&cmd.add, "add", false, "add mirror to cache")
	fs.BoolVar(&cmd.remove, "remove", false, "remove mirror from cache")
	fs.BoolVar(&cmd.list, "list", false, "list all mirrors in cache")
	fs.StringVar(&cmd.prefix, "p", "", "package name to mirror")
	fs.StringVar(&cmd.repo, "r", "", "repo mirror for package")
	fs.StringVar(&cmd.vcs, "s", "git", "vcs for package default is git")
}

func (cmd *mirrorCommand) Run(ctx *dep.Ctx, args []string) error {
	if cmd.examples {
		ctx.Err.Println(strings.TrimSpace(mirrorExamples))
		return nil
	}

	if err := cmd.validateFlags(); err != nil {
		return err
	}

	mirrors.Load()

	if cmd.add {

	} else if cmd.remove {

	} else if cmd.list {

	}

	return nil
}

func (cmd *mirrorCommand) validateFlags() error {
	if cmd.add && cmd.remove && cmd.list {
		return errors.New("cannot pass both -add and -remove and -list")
	}

	return nil
}
