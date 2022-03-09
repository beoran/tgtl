package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	//	"sort"
)

import "github.com/beoran/tgtl"
import "github.com/peterh/liner"

func runLine(env *tgtl.Environment, in string) *tgtl.Error {
	parsed, err := tgtl.Parse(in)
	if err != nil {
		return err
	}
	if parsed == nil {
		return tgtl.ErrorFromString("No parse results")
	}
	val, eff := parsed.Eval(env)
	if val != nil {
		env.Printi(">>${1}\n", val)
	} else {
		env.Printi(">>nil\n")
	}
	err, ok := eff.(*tgtl.Error)
	if ok {
		return err
	}
	return nil
}

func runLines(env *tgtl.Environment, line *liner.State) error {
	buf := ""
	for {
		if in, err := line.Prompt("> "); err == nil {
			first := ';'
			if len(in) > 0 {
				first = rune(in[0])
			}
			if first == '\\' {
				buf = buf + "\n" + in[1:len(in)]
			} else {
				if len(buf) > 0 {
					buf = buf + "\n" + in
				} else {
					buf = in + "\n"
				}
				rerr := runLine(env, buf)
				buf = ""
				if rerr != nil {
					env.Printi("Error ${1}: \n", tgtl.String(rerr.Message))
				}
			}
			line.AppendHistory(in)
		} else if err == liner.ErrPromptAborted {
			env.Printi("Aborted\n")
			return nil
		} else if err == io.EOF {
			return nil
		} else {
			env.Printi("Error reading line: ${1}\n", tgtl.ErrorFromError(err))
		}
	}
	return nil
}

func runFile(env *tgtl.Environment, name string) *tgtl.Error {
	fin, err := os.Open(name)
	if err != nil {
		return tgtl.ErrorFromError(err)
	}
	defer fin.Close()
	buf, err := ioutil.ReadAll(fin)
	if err != nil {
		return tgtl.ErrorFromError(err)
	}
	in := string(buf)

	parsed, rerr := tgtl.Parse(in)
	if rerr != nil {
		return rerr
	}
	if parsed == nil {
		return tgtl.ErrorFromString("Parse result is empty.")
	}
	args := tgtl.List{}
	for _, a := range os.Args {
		args = append(args, tgtl.String(a))
	}
	_, reff := parsed.Eval(env, args...)
	rerr, ok := reff.(*tgtl.Error)
	if ok {
		return rerr
	}
	return nil
}

func main() {
	// console := muesli.NewStdConsole()
	env := &tgtl.Environment{}
	env.Out = os.Stdout
	env.Push()

	env.RegisterBuiltins()
	env.RegisterTuringCompleteBuiltins()
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	home, _ := os.UserHomeDir()
	historyName := filepath.Join(home, ".tgtl_history")

	if f, err := os.Open(historyName); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			name := os.Args[i]
			rerr := runFile(env, name)
			if rerr != nil {
				sname := tgtl.String(name)
				env.Printi("error in ${1}: ${2}\n", sname,
					rerr)
			}
		}
		return
	}
	line.SetWordCompleter(func(line string, pos int) (head string, c []string, tail string) {
		return tgtl.WordCompleter(*env, line, pos)
	})
	runLines(env, line)

	if f, err := os.Create(historyName); err != nil {
		env.Printi("Error writing history file: ${1}\n", tgtl.ErrorFromError(err))
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}
