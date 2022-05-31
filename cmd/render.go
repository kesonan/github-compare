// MIT License
//
// Copyright (c) 2022 anqiansong
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anqiansong/github-compare/pkg/stat"
	ui "github.com/dcorbe/termui-dpc"
	"github.com/dcorbe/termui-dpc/widgets"
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
		if len(list) == 1 {
			return renderDetail(list[0])
		}

		t, err := createTable(list, true, false)
		if err != nil {
			return err
		}

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

func createTable(list []stat.Data, emoji bool, exportCSV bool) (table.Writer, error) {
	data, err := convert2ViperList(list)
	if err != nil {
		return nil, err
	}

	t := table.NewWriter()
	t.AppendHeader(createRow("name", "fullName", false, data...))
	if exportCSV {
		t.AppendRows([]table.Row{
			createRow("description", "description", false, data...),
			createRow("tags", "tags", false, data...),
			createRow("latestMonthStargazers", "latestMonthStargazers.data", false, data...),
			createRow("latestWeekForks", "latestWeekForks.data", false, data...),
			createRow("latestWeekCommits", "latestWeekCommits.data", false, data...),
			createRow("latestWeekIssues", "latestWeekIssues.data", false, data...),
		})
	}
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

	return t, nil
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
	"starCount":            "🌟 ",
	"latestDayStarCount":   "📊 ",
	"latestWeekStarCount":  "📉 ",
	"latestMonthStarCount": "📈 ",
	"forkCount":            "👏 ",
	"watcherCount":         "👀 ",
	"issue":                "💪 ",
	"pull":                 "💯 ",
	"contributorCount":     "👥 ",
	"releaseCount":         "🚀 ",
	"avgReleasePeriod":     "🔭 ",
	"latestReleaseAt":      "🎯 ",
	"lastPushedAt":         "🕦 ",
	"lastUpdatedAt":        "📝 ",
}

func createRow(title string, field string, emoji bool, data ...*viper.Viper) table.Row {
	if emoji {
		title = emojiMap[field] + title
	}

	ret := table.Row{title}
	for _, e := range data {
		ret = append(ret, e.Get(field))
	}

	return ret
}

func renderDetail(st stat.Data) error {
	data, err := convert2Viper(st)
	if err != nil {
		return err
	}

	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	starBar := createBarChart(st.LatestMonthStargazers, "Star(Latest Week)", ui.ColorRed,
		func() []ui.Color {
			var colorList []ui.Color
			for i := 1; i < 18; i++ {
				colorList = append(colorList, ui.Color(i))
			}
			return colorList
		}()...)

	forkBar := createBarChart(st.LatestWeekForks, "Forks(Latest Week)", ui.ColorGreen)
	commitBar := createBarChart(st.LatestWeekCommits, "Commits(Latest Week)", ui.ColorYellow)
	pullBar := createBarChart(st.LatestWeekPulls, "Pulls(Latest Week)", ui.ColorWhite)
	issueBar := createBarChart(st.LatestWeekIssues, "Issues(Latest Week)", ui.ColorCyan)

	desc := creatParagraph("About", ui.ColorYellow, func() []string {
		return []string{
			fmt.Sprintf("[◉ Homepage: %s](fg:blue)", data.GetString("homepage")),
			fmt.Sprintf("[◉ Description: %s](fg:white)", data.GetString("description")),
			fmt.Sprintf("◉ Tags: %s", formatTags(data.GetStringSlice("tags"))),
		}
	}()...)
	desc.TextStyle = ui.NewStyle(ui.ColorGreen)

	metrics1 := creatParagraph("Metrics1", ui.ColorRed, func() []string {
		return []string{
			fmt.Sprintf("[◉ TotalStars: %s](fg:red)", data.GetString("starCount")),
			fmt.Sprintf("[◉ TotalForks: %s](fg:green)", data.GetString("forkCount")),
			fmt.Sprintf("[◉ TotalWatcers: %s](fg:yellow)", data.GetString("watcherCount")),
			fmt.Sprintf("[◉ TotalContributors: %s](fg:cyan)", data.GetString("contributorCount")),
		}
	}()...)

	metrics2 := creatParagraph("Metrics2", ui.ColorGreen, func() []string {
		return []string{
			fmt.Sprintf("[◉ LatestDayStars: %s](fg:red)", data.GetString("latestDayStarCount")),
			fmt.Sprintf("[◉ LatestWeekStars: %s](fg:green)", data.GetString("latestWeekStarCount")),
			fmt.Sprintf("[◉ LatestMonthStars: %s](fg:yellow)",
				data.GetString("latestMonthStarCount")),
			fmt.Sprintf("[◉ ReleaseCount: %s](fg:cyan)", data.GetString("releaseCount")),
		}
	}()...)

	metrics3 := creatParagraph("Metrics3", ui.ColorYellow, func() []string {
		return []string{
			fmt.Sprintf("[◉ Issue: %s](fg:red)", data.GetString("issue")),
			fmt.Sprintf("[◉ Pull: %s](fg:green)", data.GetString("pull")),
			fmt.Sprintf("[◉ License: %s](fg:yellow)", data.GetString("license")),
			fmt.Sprintf("[◉ Language: %s](fg:cyan)", data.GetString("language")),
		}
	}()...)

	metrics4 := creatParagraph("Metrics4", ui.ColorCyan, func() []string {
		return []string{
			fmt.Sprintf("[◉ Age: %s](fg:red)", data.GetString("age")),
			fmt.Sprintf("[◉ LastRelease: %s](fg:green)", data.GetString("latestReleaseAt")),
			fmt.Sprintf("[◉ LastPushed: %s](fg:yellow)", data.GetString("lastPushedAt")),
			fmt.Sprintf("[◉ LastUpdated: %s](fg:cyan)", data.GetString("lastUpdatedAt")),
		}
	}()...)

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(
		ui.NewRow(1.0/4, ui.NewCol(1.0, starBar)),
		ui.NewRow(1.0/4,
			ui.NewCol(1.0/4, forkBar),
			ui.NewCol(1.0/4, commitBar),
			ui.NewCol(1.0/4, pullBar),
			ui.NewCol(1.0/4, issueBar),
		),
		ui.NewRow(1.0/4, ui.NewCol(1.0, desc)),
		ui.NewRow(1.0/4,
			ui.NewCol(1.0/4, metrics1),
			ui.NewCol(1.0/4, metrics2),
			ui.NewCol(1.0/4, metrics3),
			ui.NewCol(1.0/4, metrics4),
		),
	)
	ui.Render(grid)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>", "<Escape>":
			ui.Clear()
			return nil
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			grid.SetRect(0, 0, payload.Width, payload.Height)
			ui.Clear()
			ui.Render(grid)
		}
	}
}

var colorString = []string{"black", "red", "green", "blue", "magenta", "cyan"}

func formatTags(tags []string) string {
	var ret []string
	for idx, e := range tags {
		colorIndex := idx % len(colorString)
		ret = append(ret, fmt.Sprintf("[%s](fg:white,bg:%s)", e, colorString[colorIndex]))
	}
	return strings.Join(ret, " ")
}

func creatParagraph(title string, titleColor ui.Color, list ...string) *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Text = strings.Join(list, "\n")
	p.Title = title
	p.TitleStyle = ui.NewStyle(titleColor)
	return p
}

func createBarChart(data stat.Chart, title string, titleColor ui.Color,
	barColors ...ui.Color) *widgets.BarChart {
	maxVal := func() float64 {
		var maxVal float64
		for _, e := range data.Data {
			if e > maxVal {
				maxVal = e
			}
		}
		return maxVal + 2
	}

	bar := widgets.NewBarChart()
	bar.Title = title
	bar.Data = data.Data
	bar.Labels = data.Labels
	bar.MaxVal = maxVal()
	bar.TitleStyle = ui.NewStyle(titleColor)
	if len(barColors) > 0 {
		bar.BarColors = barColors
		bar.LabelStyles = func() []ui.Style {
			var colorList []ui.Style
			for _, c := range barColors {
				colorList = append(colorList, ui.NewStyle(c))
			}
			return colorList
		}()
	}

	bar.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
	return bar
}
