package opt

// Custom types for parameters
type Separator string

// Boolean flag types with constants
type BeforeFlag bool
const (
	Before   BeforeFlag = true
	NoBefore BeforeFlag = false
)

type RegexFlag bool
const (
	Regex   RegexFlag = true
	NoRegex RegexFlag = false
)

// Flags represents the configuration options for the tac command
type Flags struct {
	Separator Separator  // Record separator (-s)
	Before    BeforeFlag // Separator is before instead of after (-b)
	Regex     RegexFlag  // Interpret separator as regular expression (-r)
}

// Configure methods for the opt system
func (s Separator) Configure(flags *Flags) { flags.Separator = s }
func (b BeforeFlag) Configure(flags *Flags) { flags.Before = b }
func (r RegexFlag) Configure(flags *Flags)  { flags.Regex = r }
