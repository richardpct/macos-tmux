// tmux package
package main

import (
	"flag"
	"fmt"
	"github.com/richardpct/pkgsrc"
	"log"
	"os"
	"os/exec"
	"path"
)

var destdir = flag.String("destdir", "", "directory installation")
var pkg pkgsrc.Pkg

const (
	name     = "tmux"
	vers     = "3.0a"
	ext      = "tar.gz"
	url      = "https://github.com/tmux/tmux/releases/download/" + vers
	hashType = "sha256"
	hash     = "34901232f486fd99f3a39e864575e658b5d49f43289ccc6ee57c365f2e2c2980"
)

func checkArgs() error {
	if *destdir == "" {
		return fmt.Errorf("Argument destdir is missing")
	}
	return nil
}

func configure() {
	fmt.Println("Waiting while configuring ...")
	cmd := exec.Command("./configure", "--prefix="+*destdir)
	cmd.Env = append(os.Environ(),
		"LDFLAGS="+"-L"+*destdir+"/lib",
		"CPPFLAGS="+"-I"+*destdir+"/include",
	)
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}

func build() {
	fmt.Println("Waiting while compiling ...")
	cmd := exec.Command("make", "-j"+pkgsrc.Ncpu)
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}

func install() {
	fmt.Println("Waiting while installing ...")
	cmd := exec.Command("make", "install")
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}

func main() {
	flag.Parse()
	if err := checkArgs(); err != nil {
		log.Fatal(err)
	}

	pkg.Init(name, vers, ext, url, hashType, hash)
	pkg.CleanWorkdir()
	if !pkg.CheckSum() {
		pkg.DownloadPkg()
	}
	if !pkg.CheckSum() {
		log.Fatal("Package is corrupted")
	}

	pkg.Unpack()
	wdPkgName := path.Join(pkgsrc.Workdir, pkg.PkgName)
	if err := os.Chdir(wdPkgName); err != nil {
		log.Fatal(err)
	}
	configure()
	build()
	install()
}
