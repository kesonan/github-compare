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

package cmd

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/anqiansong/github-compare/pkg/stat"
	"github.com/briandowns/spinner"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/viper"
)

func render(spinner *spinner.Spinner, list ...stat.Data) error {
	spinner.Stop()
	var data []*viper.Viper
	for _, e := range list {
		v := viper.New()
		v.SetConfigType("json")
		d, err := json.Marshal(e)
		if err != nil {
			return err
		}
		err = v.ReadConfig(bytes.NewBuffer(d))
		if err != nil {
			return err
		}
		data = append(data, v)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(createRow("metrics", "fullName", data...))
	t.AppendRows([]table.Row{
		createRow("🏠 homepage", "homepage", data...),
		createRow("🌎 language", "language", data...),
		createRow("📌 license", "license", data...),
		createRow("⏰ age", "age", data...),
		createRow("🌟 stars", "starCount", data...),
		createRow("📊 latestDayStarCount", "latestDayStarCount", data...),
		createRow("📉 latestWeekStarCount", "latestWeekStarCount", data...),
		createRow("📈 latestMonthStarCount", "latestMonthStarCount", data...),
		createRow("👏 forks", "forkCount", data...),
		createRow("👀 watchers", "watcherCount", data...),
		createRow("💪 issues", "issue", data...),
		createRow("💯 pull requests", "pull", data...),
		createRow("👥 contributors", "contributorCount", data...),
		createRow("🚀 releases", "releaseCount", data...),
		createRow("🔭 release circle(avg)", "avgReleasePeriod", data...),
		createRow("🎯 lastRelease", "latestReleaseAt", data...),
		createRow("🕦 lastCommit", "lastPushedAt", data...),
		createRow("📝 lastUpdate", "lastUpdatedAt", data...),
	})
	t.SetStyle(table.StyleLight)
	t.Render()
	return nil
}

func createRow(title string, field string, data ...*viper.Viper) table.Row {
	ret := table.Row{title}
	for _, e := range data {
		ret = append(ret, e.Get(field))
	}
	return ret
}
