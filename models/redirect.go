package models

type Redirect struct {
	Path string `yaml:"path" json:"path"`
	Url  string `yaml:"url"  json:"url"`
}
