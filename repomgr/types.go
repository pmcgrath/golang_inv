package main

type repositoryDetail struct {
	ParentName   string
	Name         string
	ProtocolUrls map[string]string
}

type connectionAttributes struct {
	Url        string
	Username   string
	Password   string
	ParentName string
}
