package main

type connectionAttributes struct {
	Url      string
	Username string
	Password string
}

type repositoryDetail struct {
	ParentName   string
	Name         string
	ProtocolUrls map[string]string
}

type provider interface {
	getRepos(string) ([]repositoryDetail, error)
}
