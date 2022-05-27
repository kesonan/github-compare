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
	"fmt"

	"github.com/anqiansong/github-compare/pkg/stat"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func render(printStyle style, list ...stat.Data) error {
	var prettyText string
	switch printStyle {
	case styleJSON:
		data, _ := json.MarshalIndent(list, "", "  ")
		prettyText = string(data)
	case styleYAML:
		data, _ := yaml.Marshal(list)
		prettyText = string(data)
	default:
		data, err := convert2ViperList(list)
		if err != nil {
			return err
		}

		t := createTable(data, true)
		t.SetStyle(table.StyleLight)
		prettyText = t.Render()
	}
	fmt.Println(prettyText)
	return nil
}

func convert2ViperList(list []stat.Data) ([]*viper.Viper, error) {
	var data []*viper.Viper
	for _, e := range list {
		v, err := convert2Viper(e)
		if err != nil {
			return nil, err
		}
		data = append(data, v)
	}
	return data, nil
}

func createTable(data []*viper.Viper, emoji bool) table.Writer {
	t := table.NewWriter()
	t.AppendHeader(createRow("metrics", "fullName", false, data...))
	t.AppendRows([]table.Row{
		createRow("homepage", "homepage", emoji, data...),
		createRow("language", "language", emoji, data...),
		createRow("license", "license", emoji, data...),
		createRow("age", "age", emoji, data...),
		createRow("stars", "starCount", emoji, data...),
		createRow("latestDayStarCount", "latestDayStarCount", emoji, data...),
		createRow("latestWeekStarCount", "latestWeekStarCount", emoji, data...),
		createRow("latestMonthStarCount", "latestMonthStarCount", emoji, data...),
		createRow("forks", "forkCount", emoji, data...),
		createRow("watchers", "watcherCount", emoji, data...),
		createRow("issues", "issue", emoji, data...),
		createRow("pull requests", "pull", emoji, data...),
		createRow("contributors", "contributorCount", emoji, data...),
		createRow("releases", "releaseCount", emoji, data...),
		createRow("release circle(avg)", "avgReleasePeriod", emoji, data...),
		createRow("lastRelease", "latestReleaseAt", emoji, data...),
		createRow("lastCommit", "lastPushedAt", emoji, data...),
		createRow("lastUpdate", "lastUpdatedAt", emoji, data...),
	})
	return t
}

func convert2Viper(e stat.Data) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("json")
	d, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	err = v.ReadConfig(bytes.NewBuffer(d))
	if err != nil {
		return nil, err
	}
	return v, nil
}

var emojiMap = map[string]string{
	"homepage":             "🏠 ",
	"language":             "🌎 ",
	"license":              "📌 ",
	"age":                  "⏰ ",
	"stars":                "🌟 ",
	"latestDayStarCount":   "📊 ",
	"latestWeekStarCount":  "📉 ",
	"latestMonthStarCount": "📈 ",
	"forks":                "👏 ",
	"watchers":             "👀 ",
	"issues":               "💪 ",
	"pull requests":        "💯 ",
	"contributors":         "👥 ",
	"releases":             "🚀 ",
	"release circle(avg)":  "🔭 ",
	"lastRelease":          "🎯 ",
	"lastCommit":           "🕦 ",
	"lastUpdate":           "📝 ",
}

func createRow(title string, field string, emoji bool, data ...*viper.Viper) table.Row {
	ret := table.Row{title}
	for _, e := range data {
		title := fmt.Sprintf("%v", e.Get(field))
		if emoji {
			title += emojiMap[field]
		}
		ret = append(ret, e.Get(field))
	}
	return ret
}
