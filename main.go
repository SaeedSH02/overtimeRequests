package main

import (
	"shiftdony/cmd"
	log "shiftdony/logs"

	"github.com/spf13/cobra"
)

func main() {
	log.Initialize()

	var root = &cobra.Command{
		Use:   "shiftdoni",
		Short: "Change shift very simple",
	}

	root.AddCommand(cmd.Start())

	if err := root.Execute(); err != nil {
		log.Gl.Fatal(err.Error())
	}
}
