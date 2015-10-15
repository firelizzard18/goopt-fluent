package goopt_fluent

type Option struct {
   set *OptionSet
   group string
   help string
   
   names []string
   shortnames string
   
   processor Processor
   validate []Validator
   found bool
   
   alternate *Option
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