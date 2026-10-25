package dockerhub

import "time"

type RepositoryResponse struct {
	Count    int          `json:"count"`
	Next     *string      `json:"next"`
	Previous *string      `json:"previous"`
	Results  []Repository `json:"results"`
}

type Repository struct {
	Name              string    `json:"name"`
	Namespace         string    `json:"namespace"`
	RepositoryType    string    `json:"repository_type"`
	Status            int       `json:"status"`
	StatusDescription string    `json:"status_description"`
	Description       string    `json:"description"`
	IsPrivate         bool      `json:"is_private"`
	StarCount         int       `json:"star_count"`
	PullCount         int       `json:"pull_count"`
	LastUpdated       time.Time `json:"last_updated"`
	LastModified      time.Time `json:"last_modified"`
	DateRegistered    time.Time `json:"date_registered"`
	Affiliation       string    `json:"affiliation"`
	MediaTypes        []string  `json:"media_types"`
	ContentTypes      []string  `json:"content_types"`
	Categories        []string  `json:"categories"`
	StorageSize       *int64    `json:"storage_size"`
	Source            any       `json:"source"`
}
