# Blog

A cli tool and framework for running a self hosted blog.

## Overview

1. Content
   - write your content (can't help you here)
2. Compilation
   - markdown is translated into html (via pandoc)
     - all content is wrapped by a `layout.html` file
   - markdown metadata is written to a SQLite database
   - an output directory is created that mirrors the content directory
3. Serve
   - html content is served based on the file structure of the output directory.

Here is a sample. Notice how directory structure is preserved from content -> served_content

```
content
├── blog
│   ├── foo.md
│   ├── bar.md
│   └── baz.md
└── index.md

---

served_content
├── blog
│   ├── foo.html
│   ├── bar.html
│   └── baz.html
└── index.html
```

## Defaults

Everything works out of the box with no customization.

### Layout

Each file is wrapped in with the [layout file](internal/pandoc/layout.html).
This takes advantage of Go's templating system. You can modify this file, or create your own. Simply just include `{{ .Content }}` where you'd like your content to go.

### Styling

The default styles are found [styles.css](internal/styles/styles.css).
You can modify this file, or create your own.
These will be written to `/static/styles.css` in the compilation directory.

## Peculiarities

- The compilation step is idempotent. Moreover, it never modifies anything. Everything (including static files) is copied and written to the compiled target directory.
  This may be bad practice, but I like being able to delete the entire thing without having to worry about losing anything.
- `index.md` in the root of your content directory will be mapped to root(`/`)
  - note that `index.md` in subdirectories just behave as regular files
- `/all` is a reserved route - it will show a date ordered list of all your content
  - note: `/all/other/path` is not affected by this rule.
- The `/static` path is a reserved set of routes(e.g. `/static/*`). Use this to store css and images if you like
- At least for now, the database tables are truncated on each recompile. This is because I am lazy. More specifically, it is much easier to just reindex/recompute all the metadata each time.

## Dependencies

- sqlite
- pandoc
