package models

import "github.com/google/uuid"

type Request struct {
	ID            uuid.UUID `json:"id"`
	GitHubUsername string    `json:"github_username"`
	CommitHash    string    `json:"commit_hash"`
}