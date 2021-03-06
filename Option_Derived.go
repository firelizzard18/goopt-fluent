package goopt_fluent

import (
   "errors"
)

func (o *Option) Required() *Option {
   return o.Validation(func(o *Option, get GetOption) error {
      if o.Found() {
         return nil
      }
      return errors.New("Option group '" + o.group + "' is required and was not specified")
   })
}

func (o *Option) RequiredIfSpecified(otherOptGrp string) *Option {
   var other *Option
   var ok bool
   
   if other, ok = o.set.opts[otherOptGrp]; !ok {
      panic("Option '" + otherOptGrp + "' has not been defined yet")
   }
   
   return o.Validation(func(o *Option, get GetOption) error {
      if o.Found() {
         return nil
      }
      
      if other.Found() {
         return nil
      }
      
      return errors.New("Option '" + o.group + "' is required and was not specified")
   })
}

func (o *Option) MutuallyExclusiveWithAlternates(required bool) *Option {
   return o.Validation(func(o *Option, get GetOption) error {
      count := 0
      for _, alt := range o.allAlternates() {
         if alt.found {
            count++
         }
      }
      if count == 0 && required {
         return errors.New("One option of group '" + o.group + "' must be specified")
      }
      if count > 1 {
         return errors.New("Only one option of group '" + o.group + "' may be specified")
      }
      return nil
   })
}

func (o *Option) MutuallyExclusiveWith(required bool, otherOptGrps ...string) *Option {
   opts := make([]*Option, len(otherOptGrps) + 1)
   opts[0] = o
   
   for i, id := range otherOptGrps {
      var ok bool
      var other *Option
      
      if other, ok = o.set.opts[id]; !ok {
         panic("Option '" + id + "' has not been defined yet")
      }
      
      opts[i + 1] = other
   }
   
   return o.Validation(func(o *Option, get GetOption) error {
      count := 0
      
      // count the number of options that were specicied
      for _, opt := range opts {
         if opt.Found() {
            count++
         }
      }
      
      groups := func(word string) (str string) {
         return quoteJoin(word, len(opts), func (i int) string { return opts[i].group })
      }
      
      switch count {
         case 0:
            // if none were specified, it depends on whether or not one is required to be specified
            if required {
               return errors.New("Neither " + groups(" nor ") +" were specified")
            }
            return nil
            
         case 1:
            // if exactly 1 was specified, we're good
            return nil
            
         default:
            return errors.New("More than one of options " + groups(" and ") + " were specified")
      }
   })
}