package cli

import (
	"flag"
	"fmt"
	"github.com/jcocozza/jbf/internal/dal/sqlite"
	"github.com/jcocozza/jbf/internal/pandoc"
	"github.com/jcocozza/jbf/internal/serve"
	"github.com/jcocozza/jbf/internal/service"
	"html/template"
	"os"
	"path/filepath"
)

const (
	defaultContentDir string = "content"
	defaultOutputDir  string = "served_content"
	defaultStaticDir  string = ""
	defaultLayoutPath string = ""
)

func initService() (*service.Service, error) {
	db, err := sqlite.Connect()
	if err != nil {
		return nil, err
	}
	dal := sqlite.NewSQLiteRepository(db)
	return service.NewService(dal), nil
}

func initServiceWithClean() (*service.Service, error) {
	db, err := sqlite.ConnectAndClean()
	if err != nil {
		return nil, err
	}
	dal := sqlite.NewSQLiteRepository(db)
	return service.NewService(dal), nil
}

func help() {
	fmt.Fprintln(os.Stdout, "Usage")
	fmt.Fprintf(os.Stdout, "  %s <command> [options]\n", os.Args[0])
	fmt.Fprintln(os.Stdout, "Commands:")
	fmt.Fprintln(os.Stdout, "  init     set up the project")
	fmt.Fprintln(os.Stdout, "  compile  compile input to html")
	fmt.Fprintln(os.Stdout, "  serve    serve content")
	fmt.Fprintln(os.Stdout, "  new      create a new file in the content directory")
	fmt.Fprintf(os.Stdout, "use %s <command> --help for more details", os.Args[0])
}

func checkHelp(cmd *flag.FlagSet) bool {
	if cmd == nil {
		if flag.Lookup("help") != nil {
			flag.Usage()
			return true
		}
	}
	if cmd.Lookup("help") != nil {
		cmd.Usage()
		return true
	}
	return false
}

func initCmd() {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	h := checkHelp(initCmd)
	if h {
		return
	}
	initCmd.Parse(os.Args[2:])
	fmt.Fprintln(os.Stdout, "setting up...")
	err := sqlite.CreateDB()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	fmt.Fprintln(os.Stdout, "init complete")
}

func compileCmd() {
	var inputDir string
	var outputDir string
	var templateLayoutPath string
	var staticDir string
	compileCmd := flag.NewFlagSet("compile", flag.ExitOnError)
	compileCmd.StringVar(&inputDir, "content-dir", defaultContentDir, "the root directory of your content")
	compileCmd.StringVar(&outputDir, "output-dir", defaultOutputDir, "the root directory of where you want output to be written to")
	compileCmd.StringVar(&templateLayoutPath, "template-path", defaultLayoutPath, "point to a template file which will wrap each created file during compilation (empty uses default)")
	compileCmd.StringVar(&staticDir, "static-dir", defaultStaticDir, "the root directory of where static files are located (empty uses default styles)")
	h := checkHelp(compileCmd)
	if h {
		return
	}
	compileCmd.Parse(os.Args[2:])
	var layout *template.Template = pandoc.DefaultLayout
	if templateLayoutPath != "" {
		var err error
		layout, err = template.ParseFiles(templateLayoutPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
	}
	fmt.Fprintf(os.Stderr, "content dir: %s\n", inputDir)
	fmt.Fprintf(os.Stderr, "output dir: %s\n", outputDir)
	s, err := initServiceWithClean()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	cfg := service.Config{
		Layout: layout,
		Name:   "foo bar",
	}
	err = s.Compilation(inputDir, outputDir, staticDir, cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	fmt.Fprintf(os.Stderr, "content compliled to %s\n", outputDir)
}

func serveCmd() {
	var serveDir string
	var contentDir string
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	serveCmd.StringVar(&serveDir, "serve-dir", defaultOutputDir, "the root directory of where you files to be served from")
	serveCmd.StringVar(&contentDir, "content-dir", defaultContentDir, "the root directory of your content")
	h := checkHelp(serveCmd)
	if h {
		return
	}
	serveCmd.Parse(os.Args[2:])
	fmt.Fprintln(os.Stdout, "running serve")
	s, err := initService()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	serve.Server(s, serveDir, contentDir)
}

func newContentCmd() {
	newCmd := flag.NewFlagSet("new", flag.ExitOnError)
	var name string
	var contentDir string
	newCmd.StringVar(&name, "name", "", "name of new content to be created")
	newCmd.StringVar(&contentDir, "content-dir", defaultContentDir, "the root directory of your content")

	h := checkHelp(newCmd)
	if h {
		return
	}
	newCmd.Parse(os.Args[2:])
	if name == "" {
		fmt.Fprintf(os.Stderr, "name required")
		return
	}
	fmt.Fprintf(os.Stdout, "creating file: %s\n", filepath.Join(contentDir, name))
	s, err := initService()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	err = s.NewFile(contentDir, name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
}

func root() {
	flag.Usage = help
	var h bool
	flag.BoolVar(&h, "help", false, "show help")
	flag.Parse()
	if h || len(os.Args) < 2 {
		help()
		return
	}
	switch os.Args[1] {
	case "init":
		initCmd()
	case "compile":
		compileCmd()
	case "serve":
		serveCmd()
	case "new":
		newContentCmd()
	default:
		help()
	}
}

func CLI() {
	root()
}
