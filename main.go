package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	repoUser, repoFull := extractDirname(os.Args[1:])

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("unable to find home directory")
	}

	repoDir := fmt.Sprintf("%s/dev/src/github.com/%s", home, repoFull)
	if exists(repoDir) {
		fmt.Printf("directory exists: %q\n", repoDir)
		return
	}

	// mkdir
	repoUserDir := fmt.Sprintf("%s/dev/src/github.com/%s", home, repoUser)
	if !exists(repoUserDir) {
		if err := os.Mkdir(repoUserDir, 0755); err != nil {
			log.Fatalf("cannot create directory: %q", repoUserDir)
		}
	}

	// cmd: gh repo clone
	cmd := exec.Command("gh", "repo", "clone", repoFull, "--", "--filter=blob:none")
	cmd.Dir = repoUserDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	fmt.Printf("running command: %s\n", cmd)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func extractDirname(args []string) (string, string) {
	var repo string
	// gh repo clone hashicorp/terraform
	if args[0] == "gh" && len(args) == 4 {
		repo = args[3]
	} else {
		repo = args[0]
	}

	s := strings.Split(repo, "/")
	return s[0], repo
}

func exists(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	return info != nil
}
