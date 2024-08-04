package conf

type RepositoryConfig struct {
	Name   string `json:"name"`
	Owner  string `json:"owner"`
	Branch string `json:"branch"`
}
