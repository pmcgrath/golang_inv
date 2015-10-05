package main

// Git command result
type gitCmdResult struct {
	Repo    string
	Command string
	Output  []string
	Error   error
}

func (r gitCmdResults) Len() int {
	return len(r)
}

func (r gitCmdResults) Less(i, j int) bool {
	return r[i].Repo < r[j].Repo
}

func (r gitCmdResults) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type gitCmdResults []gitCmdResult

// Provider interface and connection attributes
type provider interface {
	getRepos(string) (repositoryDetails, error)
}

type providerConnectionAttributes struct {
	Url      string
	Username string
	Password string
}

// Repository detail
type repositoryDetail struct {
	ParentName   string
	Name         string
	ProtocolUrls map[string]string
}

type repositoryDetails []repositoryDetail
