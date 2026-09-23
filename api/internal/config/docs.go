package config

type Docs struct {
	DocsDir string `env:"DOCS_DIR" envDefault:"./docs"`
}
