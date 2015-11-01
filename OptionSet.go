package goopt_fluent

type OptionSet struct {
	opts map[string]*Option
}

func NewOptionSet() *OptionSet {
	s := new(OptionSet)
	s.opts = make(map[string]*Option)
	return s
}

func NewMergedSet(sets ...*OptionSet) *OptionSet {
	s := NewOptionSet()

	for _, set := range sets {
		for name, opt := range set.opts {
			if _, ok := s.opts[name]; ok {
				panic("Cannot merge options sets with the same option groups (" + name + ")")
			}
			s.opts[name] = opt.Clone(s)
		}
	}

	return s
}

func (s *OptionSet) DefineOption(group, help string) *Option {
	if _, ok := s.opts[group]; ok {
		panic("The option group '" + group + "' has already been definied")
	}

	o := NewOption(s, group, help)
	s.opts[group] = o
	return o
}

func (o *Option) DefineAlternate(help string) *Option {
	if len(o.validators) > 0 {
		panic("Validation can only be defined at the end of the chain of alternates")
	}

	if o.alternate != nil {
		panic("Each option can only have a single alternate defined")
	}

	alt := NewOption(o.set, o.group, help)
	o.set.opts[o.group] = alt
	alt.alternate = o
	return alt
}
