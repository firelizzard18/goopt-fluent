package goopt_fluent

func DefineOption(group, help string) *Option {
   if _, ok := options[group]; ok {
      panic("The option group '" + group + "' has already been definied")
   }
   
   o := makeOption(group, help)
   options[group] = o
   return o
}

func (o *Option) DefineAlternate(help string) *Option {
   if len(o.validate) > 0 {
      panic("Validation can only be defined at the end of the chain of alternates")
   }
   
   if o.alternate != nil {
      panic("Each option can only have a single alternate defined")
   }
   
   alt := makeOption(o.group, help)
   options[o.group] = alt
   alt.alternate = o
   return alt
}

func (o *Option) Names(names ...string) *Option {
   o.names = append(o.names, names...)
   return o
}

func (o *Option) ShortNames(names ...rune) *Option {
   o.shortnames += string(names)
   return o
}

func (o *Option) Validation(validators ...Validator) *Option {
   if o.alternate != nil {
      panic("Validation can only be defined at the end of the chain of alternates")
   }
   
   o.validate = append(o.validate, validators...)
   return o
}

func (o *Option) Process(processor Processor) *Option {
   if o.process != nil {
      panic("Processing has already been defined for this option")
   }
   
   o.processor = processor
   return o
}