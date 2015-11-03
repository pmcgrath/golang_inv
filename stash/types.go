package main

// Provider interface and connection attributes
type provider interface {
	getRepos(string) (repositoryDetails, error)
}

type providerConnectionAttributes struct {
	URL      string
	APIKey   string
	Username string
	Password string
}

// Repository detail
type repositoryDetail struct {
	ParentName   string
	Name         string
	ProtocolURLs map[string]string
}

type repositoryDetails []repositoryDetail
