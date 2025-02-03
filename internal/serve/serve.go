package serve

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/jcocozza/jbf/internal/metadata"
	"github.com/jcocozza/jbf/internal/pandoc"
	"github.com/jcocozza/jbf/internal/service"
)

type Handler struct {
	s              *service.Service
	htmlContentDir string
	baseContentDir string
}

func (h *Handler) HandleAllPage(w http.ResponseWriter, r *http.Request) {
	ml, err := h.s.ListContentByDate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(ml) == 0 {
		var data = struct{ Content string}{ Content: "nothing to see here."}
		pandoc.DefaultLayout.Execute(w, data)
		return
	}

	var currT metadata.Date = ml[0].Created
	s := ml[0].Created.String() + " <ul>"
	for _, m := range ml {
		op, err := service.GetOutputPath(m.Filepath, h.baseContentDir, h.htmlContentDir)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if op == "index.html" {
			continue
		}
		if currT.Equal(m.Created) {
			s += fmt.Sprintf("<li><a href=%s>%s</a></li>", op, m.Title)
			continue
		}
		s += " </ul>"
		currT = m.Created
		s += m.Created.String() + " <ul>"
		s += fmt.Sprintf("<li><a href=%s>%s</a></li>", op, m.Title)
	}

	var data = struct {
		Content template.HTML
		Name    string
	}{
		Content: template.HTML(s),
		Name:    "foo bar",
	}
	err = pandoc.DefaultLayout.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) handleFiles(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/all" || r.URL.Path == "/all/" {
		h.HandleAllPage(w, r)
		return
	}
	http.FileServer(http.Dir(h.htmlContentDir)).ServeHTTP(w, r)
}

func router(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.handleFiles)
	staticDir := filepath.Join(h.htmlContentDir, "static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	return mux
}

func Server(s *service.Service, htmlContentDir string, baseContentDir string) {
	h := &Handler{
		s:              s,
		htmlContentDir: htmlContentDir,
		baseContentDir: baseContentDir,
	}
	r := router(h)
	err := http.ListenAndServe(":55000", r) // TODO: allow this port to be specified
	if err != nil {
		panic(err)
	}
}
