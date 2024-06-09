package gobs

import "strings"

func compactName(name string) string {
	names := strings.Split(name, ".")
	return names[len(names)-1]
}
