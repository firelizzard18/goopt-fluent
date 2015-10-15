package goopt_fluent

import (
   "strings"
   "errors"
)

type Option struct {
   group string
   help string
   
   names []string
   shortnames string
   
   processor Processor
   validate []Validator
   found bool
   
   alternate *Option
}

func makeOption(group, help string) *Option {
   o := new(Option)
   o.group = group
   o.help = help
   
   o.names = make([]string, 0, 1)
   o.shortnames = ""
   
   o.validate = make([]Validator, 0, 1)
   o.found = false
   return o
}



type GetOption func(group string) (opt *Option, ok bool)
type Validator func(opt *Option, get GetOption) (err error)
type Processor func(opt *Option, input string) (err error)



func (o *Option) Found() bool {
   for _, alt := range o.allAlternates() {
      if alt.found {
         return true
      }
   }
   return false
}

func (o *Option) allAlternates() []*Option {
   return o.addAlternates(make([]*Option, 0, 2))
}

func (o *Option) addAlternates(alts []*Option) []*Option {
   for _, alt := range alts {
      if o == alt {
         panic("There is a self-referencing loop of alternates")
      }
   }
   
   alts = append(alts, o)
   if o.alternate == nil {
      return alts
   }
   
   return o.alternate.addAlternates(alts)
}

func (o *Option) checkShort(short string) bool {
   if short == "" {
      return false
   }
   
   return strings.ContainsAny(short, o.shortnames)
}

func (o *Option) checkLong(long string) bool {
   if long == "" {
      return false
   }
   for _, name := range o.names {
      if name == long {
         return true
      }
   }
   return false
}

func (o *Option) process(short, long, arg string) (success bool, err error) {
   success = false
   for _, alt := range o.allAlternates() {
      if alt.processor == nil {
         panic("Options must have processors")
      }
      
      if (!o.checkShort(short) && !o.checkLong(long)) {
         continue
      }
      
      if success == true {
         err = errors.New("Only one of the alternates of '" + o.group + "' can be specified")
         return
      }
      
      alt.found = true
      success = true
      if err = alt.processor(alt, arg); err != nil {
         return
      }
   }
   return
}