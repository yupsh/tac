package command_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/run"
	command "github.com/yupsh/tac"
)

// ==============================================================================
// Test Basic Functionality
// ==============================================================================

func TestTac_ThreeLines(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("line1", "line2", "line3").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"line3",
		"line2",
		"line1",
	})
}

func TestTac_SingleLine(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("single").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"single"})
}

func TestTac_TwoLines(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("first", "second").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"second",
		"first",
	})
}

func TestTac_EmptyInput(t *testing.T) {
	result := run.Quick(command.Tac())

	assertion.NoError(t, result.Err)
	assertion.Empty(t, result.Stdout)
}

func TestTac_EmptyLine(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{""})
}

func TestTac_MultipleEmptyLines(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("", "", "").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"", "", ""})
}

// ==============================================================================
// Test With Various Content
// ==============================================================================

func TestTac_Numbers(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("1", "2", "3", "4", "5").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"5",
		"4",
		"3",
		"2",
		"1",
	})
}

func TestTac_Alphabet(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("a", "b", "c", "d", "e").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"e",
		"d",
		"c",
		"b",
		"a",
	})
}

func TestTac_MixedContent(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines(
			"header",
			"",
			"body content",
			"more content",
			"",
			"footer",
		).Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"footer",
		"",
		"more content",
		"body content",
		"",
		"header",
	})
}

// ==============================================================================
// Test With Whitespace
// ==============================================================================

func TestTac_LeadingSpaces(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("  line1", "    line2", "line3").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"line3",
		"    line2",
		"  line1",
	})
}

func TestTac_TrailingSpaces(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("line1  ", "line2    ", "line3").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"line3",
		"line2    ",
		"line1  ",
	})
}

func TestTac_Tabs(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("a\tb", "c\td", "e\tf").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"e\tf",
		"c\td",
		"a\tb",
	})
}

func TestTac_OnlyWhitespace(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("   ", "\t\t", "  \t  ").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"  \t  ",
		"\t\t",
		"   ",
	})
}

// ==============================================================================
// Test With Unicode
// ==============================================================================

func TestTac_Unicode_Japanese(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("日本語", "中文", "한국어").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"한국어",
		"中文",
		"日本語",
	})
}

func TestTac_Unicode_Mixed(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("Hello", "世界", "123").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"123",
		"世界",
		"Hello",
	})
}

func TestTac_Unicode_Emoji(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("😀", "👋", "🌍").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"🌍",
		"👋",
		"😀",
	})
}

func TestTac_Unicode_Arabic(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("مرحبا", "سلام", "أهلا").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 3)
	// Lines are reversed
	assertion.Equal(t, result.Stdout[0], "أهلا", "last line first")
}

// ==============================================================================
// Test With Special Characters
// ==============================================================================

func TestTac_Punctuation(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("Hello!", "How are you?", "Goodbye.").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"Goodbye.",
		"How are you?",
		"Hello!",
	})
}

func TestTac_SpecialCharacters(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("!@#$", "%^&*()", "{}[]").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"{}[]",
		"%^&*()",
		"!@#$",
	})
}

func TestTac_Quotes(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines(`"quoted"`, `'single'`, "`backtick`").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"`backtick`",
		`'single'`,
		`"quoted"`,
	})
}

// ==============================================================================
// Test Edge Cases
// ==============================================================================

func TestTac_VeryLongLines(t *testing.T) {
	line1 := strings.Repeat("a", 5000)
	line2 := strings.Repeat("b", 5000)
	line3 := strings.Repeat("c", 5000)

	result := run.Command(command.Tac()).
		WithStdinLines(line1, line2, line3).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 3)
	assertion.Equal(t, result.Stdout[0], line3, "line 3 first")
	assertion.Equal(t, result.Stdout[1], line2, "line 2 second")
	assertion.Equal(t, result.Stdout[2], line1, "line 1 last")
}

func TestTac_ManyLines(t *testing.T) {
	lines := make([]string, 1000)
	expected := make([]string, 1000)
	for i := range lines {
		lines[i] = string(rune('0' + (i % 10)))
		expected[999-i] = lines[i]
	}

	result := run.Command(command.Tac()).
		WithStdinLines(lines...).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 1000)
	assertion.Lines(t, result.Stdout, expected)
}

func TestTac_Palindrome(t *testing.T) {
	// Lines that read the same forward and backward
	result := run.Command(command.Tac()).
		WithStdinLines("racecar", "level", "noon").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"noon",
		"level",
		"racecar",
	})
}

func TestTac_ReversedTwice(t *testing.T) {
	// Reversing twice should give original
	input := []string{"a", "b", "c", "d", "e"}

	result1 := run.Command(command.Tac()).
		WithStdinLines(input...).
		Run()

	assertion.NoError(t, result1.Err)

	// Reverse again
	result2 := run.Command(command.Tac()).
		WithStdinLines(result1.Stdout...).
		Run()

	assertion.NoError(t, result2.Err)
	assertion.Lines(t, result2.Stdout, input)
}

// ==============================================================================
// Test With Empty Lines Interspersed
// ==============================================================================

func TestTac_EmptyLinesAtStart(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("", "", "content").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"content",
		"",
		"",
	})
}

func TestTac_EmptyLinesAtEnd(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("content", "", "").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"",
		"",
		"content",
	})
}

func TestTac_EmptyLinesInMiddle(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("before", "", "", "after").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"after",
		"",
		"",
		"before",
	})
}

// ==============================================================================
// Test Error Handling
// ==============================================================================

func TestTac_InputError(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinError(errors.New("read failed")).
		Run()

	assertion.ErrorContains(t, result.Err, "read failed")
}

func TestTac_OutputError(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines("test").
		WithStdoutError(errors.New("write failed")).
		Run()

	assertion.ErrorContains(t, result.Err, "write failed")
}

// ==============================================================================
// Table-Driven Tests
// ==============================================================================

func TestTac_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "three lines",
			input:    []string{"a", "b", "c"},
			expected: []string{"c", "b", "a"},
		},
		{
			name:     "two lines",
			input:    []string{"first", "second"},
			expected: []string{"second", "first"},
		},
		{
			name:     "single line",
			input:    []string{"alone"},
			expected: []string{"alone"},
		},
		{
			name:     "empty line",
			input:    []string{""},
			expected: []string{""},
		},
		{
			name:     "numbers",
			input:    []string{"1", "2", "3"},
			expected: []string{"3", "2", "1"},
		},
		{
			name:     "with empty lines",
			input:    []string{"a", "", "b"},
			expected: []string{"b", "", "a"},
		},
		{
			name:     "unicode",
			input:    []string{"こんにちは", "你好", "안녕"},
			expected: []string{"안녕", "你好", "こんにちは"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := run.Command(command.Tac()).
				WithStdinLines(tt.input...).
				Run()

			assertion.NoError(t, result.Err)
			assertion.Lines(t, result.Stdout, tt.expected)
		})
	}
}

// ==============================================================================
// Test Real-World Scenarios
// ==============================================================================

func TestTac_LogFile(t *testing.T) {
	// Simulate reversing a log file to see newest first
	result := run.Command(command.Tac()).
		WithStdinLines(
			"[2024-01-01] First entry",
			"[2024-01-02] Second entry",
			"[2024-01-03] Third entry",
			"[2024-01-04] Latest entry",
		).Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"[2024-01-04] Latest entry",
		"[2024-01-03] Third entry",
		"[2024-01-02] Second entry",
		"[2024-01-01] First entry",
	})
}

func TestTac_Code(t *testing.T) {
	// Reversing code lines
	result := run.Command(command.Tac()).
		WithStdinLines(
			"func main() {",
			"    fmt.Println(\"Hello\")",
			"}",
		).Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"}",
		"    fmt.Println(\"Hello\")",
		"func main() {",
	})
}

func TestTac_CSV(t *testing.T) {
	// Reversing CSV rows
	result := run.Command(command.Tac()).
		WithStdinLines(
			"Name,Age,City",
			"Alice,30,NYC",
			"Bob,25,LA",
			"Carol,35,SF",
		).Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"Carol,35,SF",
		"Bob,25,LA",
		"Alice,30,NYC",
		"Name,Age,City",
	})
}

func TestTac_NumberedList(t *testing.T) {
	result := run.Command(command.Tac()).
		WithStdinLines(
			"1. First item",
			"2. Second item",
			"3. Third item",
		).Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"3. Third item",
		"2. Second item",
		"1. First item",
	})
}

// ==============================================================================
// Test Flags (for coverage)
// ==============================================================================

func TestTac_WithSeparator(t *testing.T) {
	// Separator flag is defined but not currently used in implementation
	// This test ensures the flag can be set without errors
	result := run.Command(command.Tac(command.Separator("\n"))).
		WithStdinLines("a", "b", "c").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"c", "b", "a"})
}

func TestTac_WithBeforeFlag(t *testing.T) {
	// Before flag is defined but not currently used in implementation
	result := run.Command(command.Tac(command.Before)).
		WithStdinLines("a", "b").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"b", "a"})
}

func TestTac_WithRegexFlag(t *testing.T) {
	// Regex flag is defined but not currently used in implementation
	result := run.Command(command.Tac(command.Regex)).
		WithStdinLines("a", "b").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"b", "a"})
}

