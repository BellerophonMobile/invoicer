package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/BellerophonMobile/goflagbuilder"
)

const (
	togglBase       = "https://www.toggl.com/api/v8"
	togglWorkspaces = togglBase + "/workspaces"
	togglClients    = togglBase + "/clients"

	togglReportBase    = "https://toggl.com/reports/api/v2"
	togglReportSummary = togglReportBase + "/summary"
)

const (
	userAgent = "contact@bellerophonmobile.com"

	timeInputFormat  = "2006-01"
	timeOutputFormat = "2006-01-02"
	timeTextFormat   = "2006/01/02"
)

type reportData struct {
	Since time.Time
	Until time.Time
	Today time.Time

	TotalTime time.Duration
	Rate      float64
	TotalDue  float64

	Entries []timeEntry
}

type timeEntry struct {
	Title string
	Time  time.Duration
}

var config = struct {
	TokenFile string `help:"File to store your Toggl API token."`

	Template string `help:"Name of input template file."`
	Output   string `help:"Name of output file."`

	WorkspaceID string `help:"ID of workspace to report on."`
	ClientID    string `help:"ID of client to report on."`

	Month string  `help:"The month to report (format YYYY-MM)."`
	Rate  float64 `help:"Hourly rate."`
}{
	TokenFile: defaultTokenPath(),
	Template:  "invoice.tex.tpl",
	Output:    "invoice.tex",
	Month:     time.Now().AddDate(0, -1, 0).Format(timeInputFormat),
	Rate:      120,
}

var tplFuncs = template.FuncMap{
	"time": func(t time.Time) string {
		return t.Format(timeTextFormat)
	},
	"texDuration": func(d time.Duration) string {
		suffix := "hrs"
		hrs := d.Hours()
		if hrs == 1 {
			suffix = "hr"
		}
		return fmt.Sprintf("%.2f %s", hrs, suffix)
	},
	"texEscape": func(s string) string {
		return strings.Replace(s, "&", "\\&", -1)
	},
	"texCash": func(f float64) string {
		return fmt.Sprintf("\\$%.2f", f)
	},
}

func main() {
	var err error

	if _, err = goflagbuilder.From(&config); err != nil {
		log.Fatalln("Failed to parse config:", err)
		return
	}

	flag.Parse()

	var token string
	if token, err = loadToken(); err != nil || token == "" {
		log.Fatalln("Failed to load token:", err)
		return
	}

	if config.WorkspaceID == "" {
		if err := listWorkspaces(token); err != nil {
			log.Fatalln("Failed to fetch workspaces:", err)
			return
		}
	}

	if config.ClientID == "" {
		if err := listClients(token); err != nil {
			log.Fatalln("Failed to fetch clients:", err)
			return
		}
	}

	if config.WorkspaceID == "" || config.ClientID == "" {
		log.Fatalln("Supply WorkspaceID and ClientID flags.")
		return
	}

	report := getReport(token)
	processTemplate(report)
}

func defaultTokenPath() string {
	config := os.Getenv("XDG_CONFIG_HOME")
	if config == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatalln("Failed to get user:", err)
			return ""
		}
		config = filepath.Join(usr.HomeDir, ".config")
	}

	return filepath.Join(config, "invoicer", "token.txt")
}

func loadToken() (string, error) {
	token, err := ioutil.ReadFile(config.TokenFile)
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(token)), nil
}

func doRequest(url, token string, data interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(token, "api_token")

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return fmt.Errorf("status %d: %s",
			rsp.StatusCode, http.StatusText(rsp.StatusCode))
	}

	dec := json.NewDecoder(rsp.Body)
	return dec.Decode(data)
}

type workspacesResponse []struct {
	ID   uint64
	Name string
}

func listWorkspaces(token string) error {
	var rsp workspacesResponse
	err := doRequest(togglWorkspaces, token, &rsp)
	if err != nil {
		return err
	}

	fmt.Println("=== Workspaces ===\nID\tName\n------------------")
	for _, ws := range rsp {
		fmt.Printf("%d\t%s\n", ws.ID, ws.Name)
	}
	fmt.Println()

	return nil
}

type clientsResponse []struct {
	ID   uint64
	Name string
}

func listClients(token string) error {
	var rsp clientsResponse
	err := doRequest(togglClients, token, &rsp)
	if err != nil {
		return err
	}

	fmt.Println("=== Clients ===\nID\tName\n---------------")
	for _, client := range rsp {
		fmt.Printf("%d\t%s\n", client.ID, client.Name)
	}
	fmt.Println()

	return nil
}

type summaryResponse struct {
	Data []struct {
		Time int64

		Items []struct {
			Title struct {
				TimeEntry string `json:"time_entry"`
			}
			Time int64
		}
	}
}

func getReport(token string) *reportData {
	since, err := time.Parse(timeInputFormat, config.Month)
	if err != nil {
		log.Fatalln("Failed to parse month format:", err)
		return nil
	}
	until := since.AddDate(0, 1, -1)

	params := url.Values{}
	params.Set("user_agent", userAgent)
	params.Set("workspace_id", config.WorkspaceID)
	params.Set("client_ids", config.ClientID)
	params.Set("grouping", "clients")
	params.Set("subgrouping", "time_entries")
	params.Set("since", since.Format(timeOutputFormat))
	params.Set("until", until.Format(timeOutputFormat))

	url := fmt.Sprintf("%s?%s", togglReportSummary, params.Encode())

	var rsp summaryResponse
	err = doRequest(url, token, &rsp)
	if err != nil {
		log.Fatalln("failed to make request:", err)
		return nil
	}

	//fmt.Printf("Data Since: %s, Until: %s\n",
	//	 params.Get("since"), params.Get("until"))
	//for _, sum := range rsp.Data {
	//	total := time.Duration(sum.Time * int64(time.Millisecond))
	//	fmt.Printf("Total Time: %s\n", total)

	//	for _, item := range sum.Items {
	//		total = time.Duration(item.Time * int64(time.Millisecond))
	//		fmt.Printf(" - %s\t%s\n", total, item.Title.TimeEntry)
	//	}
	//}

	rep := reportData{
		Since: since,
		Until: until,
		Today: time.Now(),
		Rate:  config.Rate,
	}
	for _, client := range rsp.Data {
		rep.TotalTime += time.Duration(client.Time * int64(time.Millisecond))

		for _, item := range client.Items {
			rep.Entries = append(rep.Entries, timeEntry{
				Title: item.Title.TimeEntry,
				Time:  time.Duration(item.Time * int64(time.Millisecond)),
			})
		}
	}

	rep.TotalDue = rep.Rate * float64(rep.TotalTime/time.Hour)

	return &rep
}

func getOutputFile() *os.File {
	if config.Output == "" || config.Output == "-" {
		return os.Stdout
	}
	f, err := os.Create(config.Output)
	if err != nil {
		log.Fatalln("Failed to open output file:", err)
		return nil
	}
	return f
}

func processTemplate(data *reportData) {
	tpl, err := template.
		New(filepath.Base(config.Template)).
		Funcs(tplFuncs).
		ParseFiles(config.Template)
	if err != nil {
		log.Fatalln("Failed to read template:", err)
		return
	}

	f := getOutputFile()

	err = tpl.Execute(f, data)
	if err != nil {
		log.Fatalln("failed to execute template:", err)
	}
}
