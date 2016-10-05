package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/src-d/iba"
)

func main() {
	parser := flags.NewNamedParser("iba", flags.Default)
	parser.AddCommand("pack", "Create a new archive containing the specified items.", "", &CmdPack{})
	parser.AddCommand("unpack", "Extract to disk from the archive.", "", &CmdUnpack{})
	parser.AddCommand("list", "List the items contained on a file.", "", &CmdList{})
	parser.AddCommand("version", "Show the version information.", "", &CmdVersion{})

	_, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrCommandRequired {
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}
}

type cmd struct {
	Args struct {
		File string `positional-arg-name:"iba-file" required:"true" description:"iba file."`
	} `positional-args:"yes"`

	f *os.File
	r *iba.Reader
	w *iba.Writer
}

func (c *cmd) validate() error {
	if c.Args.File == "" {
		return fmt.Errorf("Missing .iba file, please provide a valid one.")
	}
	return nil
}

func (c *cmd) buildReader() (err error) {
	c.f, err = os.Open(c.Args.File)
	if err != nil {
		return fmt.Errorf("error opening file: %s", err)
	}

	c.r = iba.NewReader(c.f)
	return nil
}

func (c *cmd) buildWriter(append bool) (err error) {
	flags := os.O_WRONLY
	if append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_CREATE | os.O_TRUNC
	}

	c.f, err = os.OpenFile(c.Args.File, flags, 0666)
	if err != nil {
		return fmt.Errorf("error creating file: %s", err)
	}

	c.w = iba.NewWriter(c.f)
	return nil
}

func (c *cmd) close() error {
	if c.w != nil {
		if err := c.w.Close(); err != nil {
			return err
		}
	}

	return c.f.Close()
}
