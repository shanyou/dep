// Package mirrors handles managing mirrors in the running application
package mirrors

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

// Cache the location of the homedirectory.
var homeDir = ""

var mirrors map[string]*mirror

func init() {
	mirrors = make(map[string]*mirror)
}

type mirror struct {
	Repo, Vcs string
}

// Get retrieves information about an mirror. It returns.
// - bool if found
// - new repo location
// - vcs type
func Get(k string) (bool, string, string) {
	o, f := mirrors[k]
	if !f {
		return false, "", ""
	}

	return true, o.Repo, o.Vcs
}

// Load pulls the mirrors into memory
func Load() error {
	home := Home()

	op := filepath.Join(home, "mirrors.yaml")

	var ov *Mirrors
	if _, err := os.Stat(op); os.IsNotExist(err) {
		// log.Println("No mirrors.yaml file exists")
		ov = &Mirrors{
			Repos: make(MirrorRepos, 0),
		}
		return nil
	} else if err != nil {
		ov = &Mirrors{
			Repos: make(MirrorRepos, 0),
		}
		return err
	}

	var err error
	ov, err = ReadMirrorsFile(op)
	if err != nil {
		return fmt.Errorf("Error reading existing mirrors.yaml file: %s", err)
	}

	// log.Println("Loading mirrors from mirrors.yaml file")
	for _, o := range ov.Repos {
		// log.Println(fmt.Sprintf("Found mirror: %s to %s (%s)", o.Original, o.Repo, o.Vcs))
		no := &mirror{
			Repo: o.Repo,
			Vcs:  o.Vcs,
		}
		mirrors[o.Prefix] = no
	}

	return nil
}

// Home returns the Glide home directory ($GLIDE_HOME or ~/.glide, typically).
//
// This normalizes to an absolute path, and passes through os.ExpandEnv.
func Home() string {
	if homeDir != "" {
		return homeDir
	}

	if h, err := homedir.Dir(); err == nil {
		homeDir = filepath.Join(h, ".dep")
	} else {
		cwd, err := os.Getwd()
		if err == nil {
			homeDir = filepath.Join(cwd, ".dep")
		} else {
			homeDir = ".dep"
		}
	}

	return homeDir
}
