package main

// Git command result
type gitCmdResult struct {
	RepoPath string
	Command  string
	Output   []string
	Error    error
}

func (r gitCmdResults) Len() int {
	return len(r)
}

func (r gitCmdResults) Less(i, j int) bool {
	return r[i].RepoPath < r[j].RepoPath
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
	URL      string
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
