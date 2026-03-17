//go:build darwin

package robotgo

import (
	"fmt"
	"os/exec"
	"strings"
)

func typeByOSAScript(text string) error {
	cmd := exec.Command("osascript", "-e", buildTypeAppleScript(text))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("osascript type failed: %w, output=%s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func buildTypeAppleScript(text string) string {
	var b strings.Builder
	b.WriteString(`tell application "System Events"` + "\n")

	var chunk strings.Builder
	flush := func() {
		if chunk.Len() == 0 {
			return
		}
		b.WriteString(`  keystroke "`)
		b.WriteString(escapeAppleScriptText(chunk.String()))
		b.WriteString(`"` + "\n")
		chunk.Reset()
	}

	for _, r := range text {
		switch r {
		case '\n':
			flush()
			b.WriteString("  key code 36\n")
		case '\r':
			// Ignore CR; LF handles enter behavior.
		case '\t':
			flush()
			b.WriteString("  key code 48\n")
		default:
			chunk.WriteRune(r)
		}
	}
	flush()

	b.WriteString("end tell")
	return b.String()
}

func escapeAppleScriptText(s string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"\"", "\\\"",
	)
	return replacer.Replace(s)
}
