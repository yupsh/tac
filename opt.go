package command

type Separator string

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

type flags struct {
	Separator Separator
	Before    BeforeFlag
	Regex     RegexFlag
}

func (s Separator) Configure(flags *flags)  { flags.Separator = s }
func (b BeforeFlag) Configure(flags *flags) { flags.Before = b }
func (r RegexFlag) Configure(flags *flags)  { flags.Regex = r }
