package main

import (
	"flag"
	"os"
	"path/filepath"
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
	home := mirrors.Home()
	op := filepath.Join(home, "mirrors.yaml")
	_, err := os.Stat(op)
	fileNotExists := os.IsNotExist(err)
	if cmd.add {
		var ov *mirrors.Mirrors
		if fileNotExists {
			ctx.Out.Println("No mirrors.yaml file exists. Creating new one")
			ov = &mirrors.Mirrors{
				Repos: make(mirrors.MirrorRepos, 0),
			}
		} else {
			ov, err = mirrors.ReadMirrorsFile(op)
			if err != nil {
				ctx.Err.Printf("Error reading existing mirrors.yaml file: %s\n", err)
			}
		}

		found := false
		for i, re := range ov.Repos {
			if re.Prefix == cmd.prefix {
				found = true
				ctx.Out.Printf("%s found in mirrors. Replacing with new settings\n", cmd.prefix)
				ov.Repos[i].Repo = cmd.repo
				ov.Repos[i].Vcs = cmd.vcs
				break
			}
		}

		if !found {
			nr := &mirrors.MirrorRepo{
				Prefix: cmd.prefix,
				Repo:   cmd.repo,
				Vcs:    cmd.vcs,
			}
			ov.Repos = append(ov.Repos, nr)
		}

		ctx.Out.Printf("%s being set to %s\n", cmd.prefix, cmd.repo)

		err := ov.WriteFile(op)
		if err != nil {
			ctx.Err.Printf("Error writing mirrors.yaml file: %s\n", err)
		} else {
			ctx.Out.Printf("mirrors.yaml written with changes")
		}
	} else if cmd.remove {
		if fileNotExists {
			ctx.Out.Println("mirrors.yaml file not found")
			return nil
		}

		ov, err := mirrors.ReadMirrorsFile(op)
		if err != nil {
			ctx.Err.Printf("Unable to read mirrors.yaml file: %s\n", err)
			return nil
		}

		var nre mirrors.MirrorRepos
		var found bool
		for _, re := range ov.Repos {
			if re.Prefix != cmd.prefix {
				nre = append(nre, re)
			} else {
				found = true
			}
		}

		if !found {
			ctx.Out.Printf("%s was not found in mirrors\n", cmd.prefix)
		} else {
			ctx.Out.Printf("%s was removed from mirrors\n", cmd.prefix)
			ov.Repos = nre

			err = ov.WriteFile(op)
			if err != nil {
				ctx.Err.Printf("Error writing mirrors.yaml file: %s\n", err)
			} else {
				ctx.Out.Println("mirrors.yaml written with changes")
			}
		}
	} else if cmd.list {
		if fileNotExists {
			ctx.Out.Println("mirrors.yaml file not found")
			return nil
		}
		ov, err := mirrors.ReadMirrorsFile(op)
		if err != nil {
			ctx.Err.Printf("Unable to read mirrors.yaml file: %s\n", err)
			return nil
		}

		if len(ov.Repos) == 0 {
			ctx.Out.Println("No mirrors found")
			return nil
		}

		ctx.Out.Println("Mirrors...")
		for _, r := range ov.Repos {
			if r.Vcs == "" {
				ctx.Out.Printf("--> %s replaced by %s\n", r.Prefix, r.Repo)
			} else {
				ctx.Out.Printf("--> %s replaced by %s (%s)\n", r.Prefix, r.Repo, r.Vcs)
			}
		}
	}

	return nil
}

func (cmd *mirrorCommand) validateFlags() error {
	if cmd.add && cmd.remove && cmd.list {
		return errors.New("cannot pass both -add and -remove and -list")
	}

	return nil
}
