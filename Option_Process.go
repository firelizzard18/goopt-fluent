package goopt_fluent

import (
   "errors"
   "strconv"
   "os"
)

func (o *Option) MissingArg() error {
   return errors.New("Option '" + o.group + "' requires an argument")
}

func (o *Option) ExtraArg() error {
   return errors.New("Option '" + o.group + "' must not have an argument")
}

func (o *Option) ProcessAsFlag(ptr *bool, value bool) *Option {
   return o.Process(func (_ *Option, input string) error {
      if input != "" {
         return o.ExtraArg()
      }
      *ptr = value
      return nil
   })
}

func (o *Option) AsFlag(def bool) *bool {
   alt := o.alternate
   
   if alt != nil && alt.alternate != nil {
      panic("A flag must have zero or one alternates")
   }
   
   flag := &def
   o.ProcessAsFlag(flag, !def)
   if (alt != nil) {
      alt.ProcessAsFlag(flag, def)
   }
   return flag
}

func (o *Option) ProcessAsInt(ptr *int) *Option {
   return o.Process(func (_ *Option, input string) error {
      if input == "" {
         return o.MissingArg()
      }
      
      var err error
      if *ptr, err = strconv.Atoi(input); err != nil {
         return err
      }
      
      return nil
   })
}

func (o *Option) AsInt(def int) *int {
   val := &def
   for _, alt := range o.allAlternates() {
      alt.ProcessAsInt(val)
   }
   return val
}

func (o *Option) ProcessAsString(ptr *string) *Option {
   return o.Process(func (_ *Option, input string) error {
      if input == "" {
         return o.MissingArg()
      }
      
      *ptr = input
      return nil
   })
}

func (o *Option) AsString(def string) *string {
   val := &def
   for _, alt := range o.allAlternates() {
      alt.ProcessAsString(val)
   }
   return val
}

func (o *Option) ProcessAsChoice(ptr *string, choices ...string) *Option {
   return o.Process(func (_ *Option, input string) error {
      if input == "" {
         return o.MissingArg()
      }
      
      for _, choice := range choices {
         if choice == input {
            *ptr = input
            return nil
         }
      }
      
      return errors.New("Option '" + o.group + "'s value must be one of " + quoteJoin(" or", len(choices), func (i int) string { return choices[i] }))
   })
}

func (o *Option) AsChoice(def string, choices ...string) *string {
   val := &def
   for _, alt := range o.allAlternates() {
      alt.ProcessAsChoice(val, choices...)
   }
   return val
}

func (o *Option) ProcessAsFile(ptr **os.File, create bool, errorh func (string, error) error) *Option {
   return o.Process(func (_ *Option, input string) (err error) {
      if input == "" {
         err = o.MissingArg()
         return
      }
      
      if (create) {
         *ptr, err = os.Create(input)
      } else {
         *ptr, err = os.Open(input)
      }
      if (err != nil) {
         err = errorh(input, err)
      }
      return
   })
}

func (o *Option) AsFile(def *os.File, create bool, errorh func (string, error) error) **os.File {
   val := &def
   for _, alt := range o.allAlternates() {
      alt.ProcessAsFile(val, create, errorh)
   }
   return val
}