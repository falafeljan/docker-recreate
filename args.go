package main

import (
  "errors"
)

// Args foobar
type Args struct {
  containerID string
  tagName string
  pullImage bool
  deleteContainer bool
}

func parseFlags(args []string, defaultArgs *Args) (Args, int) {
  i := -1

  parsedFlags := *defaultArgs

  for j, arg := range args {
    if arg[0] == '-' {
      switch arg[1] {
      case 'p':
        parsedFlags.pullImage = true

      case 'd':
        parsedFlags.deleteContainer = true
      }

      i = j
    } else {
      break
    }
  }

  return parsedFlags, i+1
}

func parseArgs(args []string) (Args, error) {
  args = args[1:]

  defaultArgs := Args{
    pullImage: false,
    deleteContainer: false }

  parsedArgs, i := parseFlags(args, &defaultArgs)
  remainingArgs := len(args) - i

  for j := 0; j < remainingArgs; j++ {
    arg := args[i+j]

    if arg[0] == '-' {
      return defaultArgs, errors.New("Bad format")
    }

    switch j {
    case 0:
      parsedArgs.containerID = arg

    case 1:
      parsedArgs.tagName = arg
    }
  }

  return parsedArgs, nil
}
