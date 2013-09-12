// Copyright 2013 Xing Xing <mikespook@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a commercial
// license that can be found in the LICENSE file.

package main

import (
	"time"
)

type Author struct {
	Name  string
	Email string
}

type GitLabRequest struct {
	Request
	UserId            int        `json:"user_id"`
	UserName          string     `json:"user_name"`
	Repo              Repository `json:"repository"`
	Commits           []Commit
	TotalCommitsCount int `json:"total_commits_count"`
}

type GitHubRequest struct {
	Request
	Repo    Repository `json:"repository"`
	Commits []GitHubCommit
}

type GitHubCommit struct {
	Commit
	Added    []string
	Removed  []string
	Modified []string
}

type GitHubRepo struct {
	Repository
	Pledgie  string
	Watchers int
	Forks    int
	Private  bool
	Owner    Author
}

type Request struct {
	Before string
	After  string
	Ref    string
}

type Commit struct {
	Id        string
	Message   string
	Timestamp time.Time
	Url       string
	Author    Author `json:"author"`
}

type Repository struct {
	Name        string
	Url         string
	Description string
	Homepage    string
}

type ghokoErr struct {
	msg    string
	status int
}

func NewError(msg string, status int) (err *ghokoErr) {
	return &ghokoErr{
		msg:    msg,
		status: status,
	}
}

func (err *ghokoErr) Error() string {
	return err.msg
}

func (err *ghokoErr) Errno() int {
	return err.status
}
