// Copyright 2013 Xing Xing <mikespook@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a commercial
// license that can be found in the LICENSE file.

package main

import (
	"time"
)

type Repository struct {
	Name        string
	Url         string
	Description string
	Homepage    string
}

type Author struct {
	Name  string
	Email string
}

type Commit struct {
	Id        string
	Message   string
	Timestamp time.Time // format "2011-12-12T14:27:31+02:00"
	Url       string
	Author    Author `json:"author"`
}

type Request struct {
	Before            string
	After             string
	Ref               string
	UserId            int        `json:"user_id"`
	UserName          string     `json:"user_name"`
	Repo              Repository `json:"repository"`
	Commits           []Commit
	TotalCommitsCount int `json:"total_commits_count"`
}
