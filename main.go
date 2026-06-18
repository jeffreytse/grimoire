package main

import "github.com/jeffreytse/grimoire/cmd"

var Version = "dev"

func main() {
	cmd.SetVersion(Version)
	cmd.Execute()
}
