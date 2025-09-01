package gitignore

// Provides common aliases for gitignore languages.
// They should map commonly used names to the actual names used in the gitignore repository:
// https://github.com/github/gitignore/tree/main
var aliasesMap = map[string]string{
	"golang": "go",
	"cpp":    "c++",
	"py":     "python",
	"c#":     "dotnet",
	"rb":     "ruby",
	"rs":     "rust",
}
