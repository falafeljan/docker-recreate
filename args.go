package main

import (
  "errors"
)

type Args struct {
  containerId string
  tagName string
  pullImage bool
}

func parseFlags(args []string, defaultArgs *Args) (Args, int) {
  var i int

  parsedFlags := *defaultArgs

  for j, arg := range args {
    if arg[0] == '-' {
      switch arg[1] {
      case 'p':
        parsedFlags.pullImage = true
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
    pullImage: false }

  parsedArgs, i := parseFlags(args, &defaultArgs)
  remainingArgs := len(args) - i

  for j := 0; j < remainingArgs; j++ {
    arg := args[i+j]

    if arg[0] == '-' {
      return defaultArgs, errors.New("Bad format")
    }

    switch j {
    case 0:
      parsedArgs.containerId = arg

    case 1:
      parsedArgs.tagName = arg
    }
  }

  return parsedArgs, nil
}
