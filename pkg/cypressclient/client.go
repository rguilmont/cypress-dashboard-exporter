package cypressclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rguilmont/cypress-dashboard-exporter/pkg/optional"
	"github.com/sirupsen/logrus"
)

const (
	defaultTimeout = 20 * time.Second
	defaultPaging  = 3 // I doubt there's a lot of run going all the time on Cypress :)
)

type StatsFromCypressDashboard struct {
	Data struct {
		Project struct {
			ID                          string `json:"id"`
			Name                        string `json:"name"`
			IsUsingRetries              bool   `json:"isUsingRetries"`
			ShouldUpdateCypressVersion5 bool   `json:"shouldUpdateCypressVersion5"`
			Runs                        Runs   `json:"runs"`
		} `json:"project"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type Runs struct {
	TotalCount int        `json:"totalCount"`
	Nodes      RunResults `json:"nodes"`
}

// Reverse the order of nodes.
//  Reason is that cypress returns the queries in descending order, and we would like to start by
//  processing the results the ascending order.
type RunResults []RunResult

func (res RunResults) Reverse() []RunResult {
	res2 := make(RunResults, len(res))
	for i := 0; i < len(res); i++ {
		res2[len(res)-i-1] = res[i]
	}
	return res2
}

type RunResult struct {
	ID                      string      `json:"id"`
	Status                  string      `json:"status"`
	BuildNumber             int         `json:"buildNumber"`
	TotalPassed             int         `json:"totalPassed"`
	TotalFailed             int         `json:"totalFailed"`
	TotalPending            int         `json:"totalPending"`
	TotalSkipped            int         `json:"totalSkipped"`
	TotalMutedTests         int         `json:"totalMutedTests"`
	StartTime               time.Time   `json:"startTime"`
	TotalDuration           int         `json:"totalDuration"`
	ScheduledToCompleteAt   time.Time   `json:"scheduledToCompleteAt"`
	ParallelizationDisabled bool        `json:"parallelizationDisabled"`
	CancelledAt             interface{} `json:"cancelledAt"`
	TotalFlakyTests         int         `json:"totalFlakyTests"`
	Project                 struct {
		ID string `json:"id"`
	} `json:"project"`
	Ci struct {
		Provider               string `json:"provider"`
		CiBuildNumberFormatted string `json:"ciBuildNumberFormatted"`
	} `json:"ci"`
	Commit struct {
		Branch      string `json:"branch"`
		AuthorEmail string `json:"authorEmail"`
		Typename    string `json:"__typename"`
	} `json:"commit"`
	TestResults struct {
		TotalCount int          `json:"totalCount"`
		Nodes      []TestResult `json:"nodes"`
	} `json:"testResults"`
}

type TestResult struct {
	ID         string   `json:"id"`
	TitleParts []string `json:"titleParts"`
	IsFlaky    bool     `json:"isFlaky"`
	IsMuted    bool     `json:"isMuted"`
	State      string   `json:"state"`
	Duration   int      `json:"duration"`
	Instance   struct {
		ID          string    `json:"id"`
		Status      string    `json:"status"`
		Duration    int       `json:"duration"`
		CompletedAt time.Time `json:"completedAt"`
		Os          struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"os"`
		Browser struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"browser"`
		Group struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"group"`
		Spec struct {
			ID        string `json:"id"`
			ShortPath string `json:"shortPath"`
		} `json:"spec"`
	} `json:"instance"`
}

type CypressDashboardMetricsClient struct {
	httpClient *http.Client
	endpoint   url.URL

	email               string
	password            string
	authenticationToken string
}

func (cli *CypressDashboardMetricsClient) Authenticate() error {

	body := map[string]string{
		"email":    cli.email,
		"password": cli.password,
	}

	content, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := http.Post("https://authenticate.cypress.io/login/local?source=dashboard", "application/json", bytes.NewBuffer(content))
	if err != nil {
		return err
	}

	logrus.Debugln("Header results of the request to authentication", resp.Header)

	cookie := resp.Header.Get("set-cookie")
	token := strings.ReplaceAll(strings.Split(cookie, ";")[0], "cy_dashboard=", "")

	logrus.Debugln("Extracted token : ", token)
	cli.authenticationToken = token
	return nil
}

func NewCypressDashboardMetricsClient(endpoint url.URL, email, password string) CypressDashboardMetricsClient {

	client := CypressDashboardMetricsClient{
		httpClient:          http.DefaultClient,
		endpoint:            endpoint,
		email:               email,
		password:            password,
		authenticationToken: ``,
	}

	client.httpClient.Timeout = defaultTimeout
	// TEMP CONFIG
	return client
}

// GetMetrics returns a list of StatsFromCypressDashboard, or an error

type GetMetricOptions struct {
	From    optional.OptionalTime
	To      optional.OptionalTime
	Size    optional.OptionalInt
	Project string
}

func EmptyMetricOptions() GetMetricOptions {
	from := time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Now()
	defaultPaging := 10
	return GetMetricOptions{
		From:    optional.NewOptionalTime(&from),
		To:      optional.NewOptionalTime(&to),
		Size:    optional.NewOptionalInt(&defaultPaging),
		Project: "7s5okt", // This is the realworld example from Cypress
	}
}

func (client *CypressDashboardMetricsClient) GetMetrics(opts GetMetricOptions) (*StatsFromCypressDashboard, error) {

	statsURL := client.endpoint

	createReq := func() (*http.Request, error) {
		body, err := createMetricRequest(opts.Project, optional.OrElseTime(opts.From, time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC)),
			optional.OrElseTime(opts.To, time.Now()),
			optional.OrElseInt(opts.Size, defaultPaging))
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest(http.MethodPost, statsURL.String(), body)
		if err != nil {
			return nil, err
		}
		req.Header.Add("cookie", fmt.Sprintf("cy_dashboard=%v", client.authenticationToken))
		req.Header.Add("content-type", "application/json")
		return req, nil
	}

	getAnswer := func(req *http.Request) (*StatsFromCypressDashboard, error) {
		resp, err := client.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		stats := StatsFromCypressDashboard{}
		err = json.NewDecoder(resp.Body).Decode(&stats)
		if err != nil {
			return nil, err
		}
		return &stats, nil
	}

	req, err := createReq()
	if err != nil {
		return nil, err
	}

	resp, err := getAnswer(req)
	if err != nil {
		return nil, err
	}
	// Check if we had the authorization to get the dashboard. Otherwise, log in and retry.
	//  For now we'll consider that the only error we could get is unauthorized. If after the 2nd attempt it fails again
	//  then we'll return a proper error.
	if len(resp.Errors) > 0 {
		logrus.Warnf("error on first request, trying to authenticate. Error was : %v", resp.Errors)
		err := client.Authenticate()

		req2, err := createReq()
		if err != nil {
			return nil, err
		}

		resp2, err2 := getAnswer(req2)
		if err2 != nil {
			return nil, err
		}
		if len(resp2.Errors) > 0 {
			return nil, fmt.Errorf("Unrecoverable error occured : %v . Check your credentials and project ID.", resp2.Errors)
		}
		return resp2, nil
	}

	return resp, nil
}
