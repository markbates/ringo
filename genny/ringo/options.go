package ringo

type Options struct {
	Name string
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if len(opts.Name) == 0 {
		opts.Name = "Starr"
	}
	return nil
}
