/*
 * MIT License
 *
 * Copyright (c) 2022 anqiansong
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package stat

import (
	"github.com/shurcooL/githubv4"
)

type (
	IssueConnection struct {
		TotalCount githubv4.Int
	}

	Issue struct {
		Issues IssueConnection `graphql:"issues(first: 1, states: $issueStates)"`
	}

	IssueQuery struct {
		Issue Issue `graphql:"repository(owner: $owner, name: $name)"`
	}
)

func (s Stat) OpenIssueCount() githubv4.Int {
	var issueQuery IssueQuery
	_ = s.graphqlClient.Query(s.ctx, &issueQuery, map[string]interface{}{
		"owner":       githubv4.String(s.owner),
		"name":        githubv4.String(s.repo),
		"issueStates": []githubv4.IssueState{githubv4.IssueStateOpen},
	})
	return issueQuery.Issue.Issues.TotalCount
}
