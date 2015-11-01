package goopt_fluent

type Option struct {
	set   *OptionSet
	group string
	help  string

	names      []string
	shortnames string

	processor  Processor
	validators []Validator
	found      bool

	alternate *Option
}

func NewOption(set *OptionSet, group, help string) *Option {
	o := new(Option)
	o.set = set
	o.group = group
	o.help = help

	o.names = make([]string, 0, 1)
	o.shortnames = ""

	o.validators = make([]Validator, 0, 1)
	o.found = false
	return o
}

func (opt *Option) Clone(set *OptionSet) *Option {
	o := new(Option)
	o.set = set
	o.group = opt.group
	o.help = opt.help
	o.shortnames = opt.shortnames
	o.found = opt.found
	o.processor = opt.processor

	o.names = make([]string, len(opt.names))
	copy(o.names, opt.names)

	o.validators = make([]Validator, len(opt.validators))
	copy(o.validators, opt.validators)

	if opt.alternate != nil {
		o.alternate = opt.alternate.Clone(set)
	}
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

func (o *Option) Names(names ...string) *Option {
	o.names = append(o.names, names...)
	return o
}

func (o *Option) ShortNames(names ...rune) *Option {
	o.shortnames += string(names)
	return o
}

func (o *Option) Validation(validators ...Validator) *Option {
	// if o.alternate != nil {
	// 	panic("Validation can only be defined at the end of the chain of alternates")
	// }

	o.validators = append(o.validators, validators...)
	return o
}

func (o *Option) Process(processor Processor) *Option {
	if o.processor != nil {
		panic(o.group + ": processing has already been defined for this option")
	}

	o.processor = processor
	return o
}
