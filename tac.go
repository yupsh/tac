package tac

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	yup "github.com/yupsh/framework"
	"github.com/yupsh/framework/opt"
	localopt "github.com/yupsh/tac/opt"
)

// Flags represents the configuration options for the tac command
type Flags = localopt.Flags

// Command implementation
type command opt.Inputs[string, Flags]

// Tac creates a new tac command with the given parameters
func Tac(parameters ...any) yup.Command {
	cmd := command(opt.Args[string, Flags](parameters...))
	// Set default separator to newline
	if cmd.Flags.Separator == "" {
		cmd.Flags.Separator = "\n"
	}
	return cmd
}

func (c command) Execute(ctx context.Context, input io.Reader, output, stderr io.Writer) error {
	// Check for cancellation before starting
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	// If no files specified, read from stdin
	if len(c.Positional) == 0 {
		return c.processReader(ctx, input, output)
	}

	// Process each file
	for _, filename := range c.Positional {
		// Check for cancellation before each file
		if err := yup.CheckContextCancellation(ctx); err != nil {
			return err
		}

		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(stderr, "tac: %s: %v\n", filename, err)
			continue
		}

		if err := c.processReader(ctx, file, output); err != nil {
			file.Close()
			fmt.Fprintf(stderr, "tac: %s: %v\n", filename, err)
			continue
		}
		file.Close()
	}

	return nil
}

func (c command) processReader(ctx context.Context, reader io.Reader, output io.Writer) error {
	// Check for cancellation before reading
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	// Read all content with context-aware reading
	content, err := c.readAllWithContext(ctx, reader)
	if err != nil {
		return err
	}

	if len(content) == 0 {
		return nil
	}

	// Check for cancellation after reading large content
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	contentStr := string(content)

	// Split by separator
	var records []string
	sep := string(c.Flags.Separator)

	if bool(c.Flags.Regex) {
		// Use separator as regex
		re, err := regexp.Compile(sep)
		if err != nil {
			return fmt.Errorf("tac: invalid regular expression: %v", err)
		}

		// Check for cancellation before potentially expensive regex operation
		if err := yup.CheckContextCancellation(ctx); err != nil {
			return err
		}

		records = re.Split(contentStr, -1)
	} else {
		// Use separator as literal string
		records = strings.Split(contentStr, sep)
	}

	// Check for cancellation after splitting
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	// Handle the separator placement and output in reverse order
	if bool(c.Flags.Before) {
		// Separator is before the record
		for i := len(records) - 1; i >= 0; i-- {
			// Check for cancellation every 1000 records to avoid excessive overhead
			if (len(records)-1-i)%1000 == 0 {
				if err := yup.CheckContextCancellation(ctx); err != nil {
					return err
				}
			}

			if i < len(records)-1 && records[i] != "" {
				// Don't add separator before the last (first when reversed) record
				fmt.Fprint(output, sep)
			}
			fmt.Fprint(output, records[i])
		}
	} else {
		// Separator is after the record (default)
		for i := len(records) - 1; i >= 0; i-- {
			// Check for cancellation every 1000 records to avoid excessive overhead
			if (len(records)-1-i)%1000 == 0 {
				if err := yup.CheckContextCancellation(ctx); err != nil {
					return err
				}
			}

			fmt.Fprint(output, records[i])
			if i > 0 && records[i] != "" {
				// Don't add separator after the last (first when reversed) record
				fmt.Fprint(output, sep)
			}
		}
	}

	// Final check for cancellation
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	return nil
}

// readAllWithContext reads all content from a reader with context cancellation support
func (c command) readAllWithContext(ctx context.Context, reader io.Reader) ([]byte, error) {
	var result []byte
	buffer := make([]byte, 32*1024) // 32KB buffer
	totalRead := 0

	for {
		// Check for cancellation before each read
		if err := yup.CheckContextCancellation(ctx); err != nil {
			return result, err
		}

		n, err := reader.Read(buffer)
		if n > 0 {
			result = append(result, buffer[:n]...)
			totalRead += n

			// Check for cancellation every 1MB read to avoid excessive overhead
			if totalRead%(1024*1024) == 0 {
				if err := yup.CheckContextCancellation(ctx); err != nil {
					return result, err
				}
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func (c command) String() string {
	return fmt.Sprintf("tac %v", c.Positional)
}
