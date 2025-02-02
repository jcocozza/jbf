package pandoc

import (
	_ "embed"
	"html/template"
	"os/exec"
	"path/filepath"
)

//go:embed layout.html
var defaultLayout string

var DefaultLayout = template.Must(template.New("default_layout").Parse(defaultLayout))

func PandocToHTML(filepath string) (string, error) {
	cmd := exec.Command("pandoc", filepath, "--to", "html")
	fbyte, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(fbyte), nil
}

func RenameMdToHtml(mdPath string) string {
	ext := filepath.Ext(mdPath)
	cutoff := len(mdPath) - len(ext)
	newName := mdPath[0:cutoff] + ".html"
	return newName
}
