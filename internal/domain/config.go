package domain

// Config is the project-generation request, independent of CLI and filesystem concerns.
type Config struct {
	Name     string     `json:"name" yaml:"name"`
	Profile  string     `json:"profile" yaml:"profile"`
	Features []string   `json:"features" yaml:"features"`
	Next     NextConfig `json:"next" yaml:"next"`
}

type NextConfig struct {
	TypeScript bool `json:"typescript" yaml:"typescript"`
	ESLint     bool `json:"eslint" yaml:"eslint"`
	Tailwind   bool `json:"tailwind" yaml:"tailwind"`
	SrcDir     bool `json:"src_dir" yaml:"src_dir"`
	AppRouter  bool `json:"app_router" yaml:"app_router"`
	Configured bool `json:"-" yaml:"-"`
}
