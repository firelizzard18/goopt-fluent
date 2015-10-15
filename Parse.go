package goopt_fluent

import (
   "errors"
   "strings"
   "unicode/utf8"
)

var Default = NewOptionSet()

func (s *OptionSet) Parse(args []string) ([]string, error) {
   opts := make(map[*Option]bool)
   for _, opt := range s.opts {
      opts[opt] = false
   }
   
   skip := 0
   last := 0
   
main:
   for i, arg := range args[1:] {
      if len(arg) < 1 {
         // wtf, zero length arg?
         return args[last+1:], errors.New("Invalid zero-length argument")
      }

      if arg == "--" {
         // `--` means stop processing options
         return args[last+1:], nil
      }
      
      if skip > 0 {
         skip--
         continue
      }
      
      if arg[0] != '-' {
         return args[last+1:], errors.New("Unexpected non-option")
      }
      
      if len(arg) < 2 {
         return args[last+1:], errors.New("Invalid empty option ('-')")
      }
      
      shortName := rune(0)
      longName := ""
      if arg[1] == '-' {
         longName = arg[2:]
      } else if len(arg) > 2 {
         return args[last+1:], errors.New("Invalid short-form option: '" + arg + "'")
      } else {
         shortName, _ = utf8.DecodeRuneInString(arg[1:])
      }
      
      var optarg string
      bits := strings.SplitN(longName, "=", 2)
      if len(bits) == 2 {
         // --option=optarg
         longName = bits[0]
         optarg = bits[1]
      } else if len(args) > i {
         nextarg := args[i + 1]
         if len(nextarg) > 0 && nextarg[0] != '-' {
            // --option optarg
            skip = 1
            optarg = nextarg
         }
      }
      
      for opt, found := range opts {
         if found {
            continue
         }
         
         success, err := opt.process(shortName, longName, optarg)
         if err != nil {
            return args[last + 1:], err
         }
         
         if !success {
            continue
         }
         
         opts[opt] = true
         last = i
         continue main
      }
      
      return args[last + 1:], errors.New("Unrecognized option: '" + arg + "'")
   }
   
   return args[last + 1:], nil
}