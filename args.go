package cli

// Args represents command-line arguments.
type Args []string

// Len returns the number of the specified index.
func (a Args) Len(i int) int {
	if i >= 0 && i < len(a) {
		return len(a[i])
	}
	return 0
}

// Get returns the argument at the specified index.
func (a Args) Get(index int) string {
	if index >= 0 && index < len(a) {
		return a[index]
	}
	return ""
}

// Slice returns a slice of all arguments.
func (a Args) Slice() []string {
	return a
}

// Num returns the number of all arguments.
func (a Args) Num() int {
	return len(a)
}
