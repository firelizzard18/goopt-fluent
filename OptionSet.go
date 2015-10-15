package goopt_fluent

type OptionSet struct {
   opts map[string]*Option
}

func NewOptionSet() *OptionSet {
   s := new(OptionSet)
   s.opts = make(map[string]*Option)
   return s
}

func MergeSets(a, b *OptionSet) *OptionSet {
   s := new(OptionSet)
   s.opts = make(map[string]*Option)
   
   for name, opt := range a.opts {
      s.opts[name] = opt
   }
   
   for name, opt := range a.opts {
      if _, ok := s.opts[name]; ok {
         panic("Cannot merge options sets with the same option groups")
      }
      s.opts[name] = opt
   }
   
   return s
}

func (s *OptionSet) makeOption(group, help string) *Option {
   o := new(Option)
   o.set = s
   o.group = group
   o.help = help
   
   o.names = make([]string, 0, 1)
   o.shortnames = ""
   
   o.validate = make([]Validator, 0, 1)
   o.found = false
   return o
}

func (s *OptionSet) DefineOption(group, help string) *Option {
   if _, ok := s.opts[group]; ok {
      panic("The option group '" + group + "' has already been definied")
   }
   
   o := s.makeOption(group, help)
   s.opts[group] = o
   return o
}

func (o *Option) DefineAlternate(help string) *Option {
   if len(o.validate) > 0 {
      panic("Validation can only be defined at the end of the chain of alternates")
   }
   
   if o.alternate != nil {
      panic("Each option can only have a single alternate defined")
   }
   
   alt := o.set.makeOption(o.group, help)
   o.set.opts[o.group] = alt
   alt.alternate = o
   return alt
}