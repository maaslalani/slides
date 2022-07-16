package code

import (
	"regexp"
	"strings"
)

const comment = "///"

var commentRegexp = regexp.MustCompile("(?m)[\r\n]+^" + comment + ".*$")

// HideComments removes all comments from the given content.
func HideComments(content string) string {
	return commentRegexp.ReplaceAllString(content, "")
}

// RemoveComments strips all the comments from the given content.
// This is useful for when we want to actually use the content of the comments.
func RemoveComments(content string) string {
	return strings.ReplaceAll(content, comment, "")
}
