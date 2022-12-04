package generate

import (
	"regexp"
	"strings"

	pgs "github.com/lyft/protoc-gen-star"
)

func schemaName(e pgs.Entity) string {
	return strings.TrimPrefix(e.FullyQualifiedName(), ".")
}

func propertyName(f pgs.Entity) string {
	return f.Name().LowerCamelCase().String()
}

func description(e pgs.Entity) string {
	comments := e.SourceCodeInfo().LeadingComments()
	comments = strings.TrimSpace(comments)
	reg := regexp.MustCompile(`\n[ \t]`)
	for reg.MatchString(comments) {
		comments = reg.ReplaceAllString(comments, "\n")
		comments = strings.TrimSpace(comments)
	}
	return comments
}
