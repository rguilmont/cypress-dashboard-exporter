package cypressclient

import (
	"bytes"
	"encoding/json"
	"io"
	"time"
)

type ProjectID string
type Input struct {
	Page      int `json:"page"`
	TimeRange struct {
		StartDate string `json:"startDate"`
		EndDate   string `json:"endDate"`
	} `json:"timeRange"`
	PerPage int `json:"perPage"`
}

type graphqlQuery struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
}

const cypressDateFormat = "2006-01-02"

func createMetricRequest(projectID string, from time.Time, to time.Time, size int) (io.Reader, error) {
	const graphql = `query RunsList($projectId: String!, $input: ProjectRunsConnectionInput) {
		project(id: $projectId) {
		  id
		  name
		  ...FlakyRateEmptyStateProject
	  
		  runs(input: $input) {
			totalCount
			nodes {
			  id
			  ...RunsListItem
			}
		  }
		}
	  }
	  
	  fragment FlakyRateEmptyStateProject on Project {
		id
		isUsingRetries
		shouldUpdateCypressVersion5
	  }
	  
	  fragment RunsListItem on Run {
		id
		status
		buildNumber
	  
		totalPassed
		totalFailed
		totalPending
		totalSkipped
		totalMutedTests
		startTime
		totalDuration
		scheduledToCompleteAt
		parallelizationDisabled
		cancelledAt
		totalFlakyTests
		project {
		  id
		}
		ci {
		  provider
		  ciBuildNumberFormatted
		}
	  
		commit {
		  branch
		  authorEmail
		}

		testResults(input: { perPage: 500 }) {
		  totalCount
		  nodes {
			id
			...TestOverview
		  }
		}
	  }
	  
	  fragment TestOverview on TestResult {
		id
		titleParts
		isFlaky
		isMuted
		state
		duration
		instance {
		  id
		  ...DrawerRunInstance
		  spec {
			id
			shortPath
		  }
		}
	  }
	  
	  fragment DrawerRunInstance on RunInstance {
		id
		status
		duration
		completedAt
		group {
			id    
			name
		}
		os {
		  ...SpecOs
		}
		browser {
		  ...SpecBrowser
		}
	  }
	  
	  fragment SpecOs on OperatingSystem {
		name
		version
	  }
	  
	  fragment SpecBrowser on BrowserInfo {
		name
		version
	  }	  
`

	variables := Input{
		Page: 1,
		TimeRange: struct {
			StartDate string "json:\"startDate\""
			EndDate   string "json:\"endDate\""
		}{
			StartDate: from.Format(cypressDateFormat),
			EndDate:   to.Format(cypressDateFormat),
		},
		PerPage: size,
	}

	query := graphqlQuery{
		OperationName: "RunsList",
		Query:         graphql,
		Variables: map[string]interface{}{
			"projectId": projectID,
			"input":     variables,
		},
	}

	rawJsonQuery, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(rawJsonQuery), nil
}
