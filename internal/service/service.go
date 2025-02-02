package service

import (
	"github.com/jcocozza/jbf/internal/dal"
	"github.com/jcocozza/jbf/internal/metadata"
	"github.com/jcocozza/jbf/internal/pandoc"
	"github.com/jcocozza/jbf/internal/styles"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Name   string
	Layout *template.Template
}

type Service struct {
	dal dal.Repository
}

func NewService(d dal.Repository) *Service {
	return &Service{dal: d}
}

func (s *Service) NewFile(contentDir string, fname string) error {
	f, err := os.Create(filepath.Join(contentDir, fname))
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(metadata.MetadataTemplate()))
	return err
}

func (s *Service) processTag(m metadata.Metadata, tagName string) error {
	exists := s.dal.ReadTagExists(tagName)
	if exists {
		return nil
	}
	err := s.dal.CreateTag(m.ID, tagName)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) processMetadataTags(m metadata.Metadata) error {
	for _, tag := range m.Tags {
		err := s.processTag(m, tag)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) updateMetadata(m metadata.Metadata) error {
	err := s.processMetadataTags(m)
	if err != nil {
		return err
	}
	return s.dal.UpdateMetadata(m)
}

func (s *Service) createMetadata(m *metadata.Metadata) error {
	id, err := s.dal.CreateMetadata(*m)
	if err != nil {
		return err
	}
	m.ID = id
	return s.processMetadataTags(*m)
}

func (s *Service) ListContentByDate() ([]metadata.Metadata, error) {
	return s.dal.ReadAllMetadata()
}

func (s *Service) updateDB(filepath string) error {
	md, err := metadata.ExtractFromFile(filepath)
	if err != nil {
		return err
	}
	fmt.Printf("got metadata for file %s:\n%s\n", filepath, md.String())
	exists := s.dal.ReadMetadataExists(filepath)
	if exists {
		return s.updateMetadata(md)
	}
	return s.createMetadata(&md)
}

func (s *Service) processContentFile(inputPath string, outputPath string, cfg Config) error {
	err := s.updateDB(inputPath)
	if err != nil {
		return err
	}
	base, err := pandoc.PandocToHTML(inputPath)
	if err != nil {
		return err
	}
	var htmlContentBuilder strings.Builder
	var data = struct {
		Content template.HTML
		Name    string
	}{
		Content: template.HTML(base),
		Name:    cfg.Name,
	}
	err = cfg.Layout.Execute(&htmlContentBuilder, data)
	if err != nil {
		return err
	}
	ext := filepath.Ext(inputPath) // this should be .md
	name := outputPath
	cutoff := len(name) - len(ext)
	newName := name[0:cutoff] + ".html"
	filepath.Join()
	f, err := os.Create(newName)
	if err != nil {
		return err
	}
	fmt.Println("writing content", inputPath, inputPath)
	_, err = f.WriteString(htmlContentBuilder.String())
	if err != nil {
		return err
	}
	return f.Chmod(0444)
}

func (s *Service) clearCompilation(dir string) error {
	entires, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, entry := range entires {
		entryPath := filepath.Join(dir, entry.Name())
		err := os.RemoveAll(entryPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetOutputPath(inputPath, inputDir, outputDir string) (string, error) {
	relPath, err := filepath.Rel(inputDir, inputPath)
	if err != nil {
		return "", err
	}
	return pandoc.RenameMdToHtml(relPath), nil
}

func (s *Service) Compilation(contentDir string, outputDir string, staticDir string, cfg Config) error {
	// create content, converts md to html
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(contentDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(outputDir, relPath)
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		return s.processContentFile(path, destPath, cfg)
	}
	err := s.clearCompilation(outputDir)
	if err != nil {
		return err
	}
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			return err
		}
	}
	err = filepath.Walk(contentDir, walkFunc)
	if err != nil {
		return err
	}
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		static := filepath.Join(outputDir, "static")
		err := os.MkdirAll(static, 0755)
		if err != nil {
			return err
		}
		f, err := os.Create(filepath.Join(static, "styles.css"))
		if err != nil {
			return err
		}
		_, err = f.Write(styles.DefaultCSSStyles)
		if err != nil {
			return err
		}
		return f.Chmod(0444)
	}

	// copy static directory to the output dir under /static
	walkStaticFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(staticDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(outputDir, "static", relPath)
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		f, err := os.Create(destPath)
		if err != nil {
			return err
		}
		_, err = f.Write(content)
		if err != nil {
			return err
		}
		return f.Chmod(0444)
	}
	return filepath.Walk(staticDir, walkStaticFunc)
}
