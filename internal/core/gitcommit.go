package core

import (
	"time"
)

// GitCommit represents a Git commit object
type GitCommit struct {
	CommitHash string    // The commit's unique identifier (SHA-1 hash)
	Author     string    // Author name and email
	Committer  string    // Committer name and email
	Timestamp  time.Time // Timestamp of the commit
	Message    string    // The commit message
	ParentHash string    // Parent commit hash
	TreeHash   string    // The hash of the tree object associated with the commit
}

// Init initializes a new GitCommit with default values
func (commit *GitCommit) Init(author, committer, message, parentHash, treeHash string) {
	commit.Author = author
	commit.Committer = committer
	commit.Message = message
	commit.ParentHash = parentHash
	commit.TreeHash = treeHash
	commit.Timestamp = time.Now()
}

func (commit *GitCommit) BuildTreeFromIndex(gitIndex *GitIndex) {

}
