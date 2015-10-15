package goopt_fluent

import (
   "strings"
   "errors"
)

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

func (o *Option) checkShort(short rune) bool {
   if short == rune(0) {
      return false
   }
   
   return strings.ContainsAny(string(short), o.shortnames)
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

func (o *Option) process(short rune, long, arg string) (success bool, err error) {
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