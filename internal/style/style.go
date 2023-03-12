package style

import (
	"time"

	"github.com/fatih/color"
	"github.com/theckman/yacspin"
)

var (
	Spinner, _ = yacspin.New(
		yacspin.Config{
			Frequency:       100 * time.Millisecond,
			CharSet:         yacspin.CharSets[43],
			Suffix:          " retrieving posts and comments",
			SuffixAutoColon: true,
			Message:         "", // Set this to the page "after" setting from struct
			StopCharacter:   "✓",
			StopColors:      []string{"fgGreen"},
		},
	)

	Warn        = color.New(color.FgRed)
	Result      = color.New(color.FgGreen)
	Information = color.New(color.FgHiMagenta)
	Start       = color.New(color.FgHiYellow)
)
