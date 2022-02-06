package cypresscollector

import (
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rguilmont/cypress-dashboard-exporter/pkg/cypressclient"
	"github.com/rguilmont/cypress-dashboard-exporter/pkg/cypresscollector/converter"
	"github.com/rguilmont/cypress-dashboard-exporter/pkg/cypresscollector/metricsmap"
	"github.com/rguilmont/cypress-dashboard-exporter/pkg/cypresscollector/set"
	"github.com/rguilmont/cypress-dashboard-exporter/pkg/optional"
	"github.com/sirupsen/logrus"
)

type CypressDashboardCollector struct {
	// Project level
	CypressRunsCount *prometheus.Desc
	// Runs metrics
	CypressRunCount      *prometheus.Desc
	CypressRunPassed     *prometheus.Desc
	CypressRunFailed     *prometheus.Desc
	CypressRunPending    *prometheus.Desc
	CypressRunSkipped    *prometheus.Desc
	CypressRunMutedTests *prometheus.Desc
	CypressRunDuration   *prometheus.Desc
	CypressRunFlakyTests *prometheus.Desc
	CypressRunStartTime  *prometheus.Desc

	CypressRunPassedSum     *prometheus.Desc
	CypressRunFailedSum     *prometheus.Desc
	CypressRunPendingSum    *prometheus.Desc
	CypressRunSkippedSum    *prometheus.Desc
	CypressRunMutedTestsSum *prometheus.Desc
	CypressRunDurationSum   *prometheus.Desc
	CypressRunFlakyTestsSum *prometheus.Desc
	CypressRunStartTimeSum  *prometheus.Desc

	// Tests metrics
	CypressTestCount        *prometheus.Desc
	CypressTestStateSum     *prometheus.Desc
	CypressTestStateLast    *prometheus.Desc
	CypressTestDurationSum  *prometheus.Desc
	CypressTestDurationLast *prometheus.Desc

	// Other metrics for DD availability
	CypressDashboardExporterAvailable *prometheus.Desc

	cli *cypressclient.CypressDashboardMetricsClient

	// This is for keeping state
	LastDateTest        time.Time
	LastBuild           int
	TotalAnalysedBuilds int
	TotalAnalysedTests  int

	AlreadyProcessedBuilds set.IntSet
	runSummary             metricsmap.MetricMapSumValues
	testSummary            metricsmap.MetricMapSumValues
	runLatest              metricsmap.MetricMapKeepFirst
	testLatest             metricsmap.MetricMapKeepFirst
	firstRequest           bool
	project                string
}

func NewCypressDashboardCollector(endpoint url.URL, project, email, password string, keepUntil int64) (*CypressDashboardCollector, error) {

	client := cypressclient.NewCypressDashboardMetricsClient(endpoint, email, password)
	return &CypressDashboardCollector{
		CypressRunsCount: prometheus.NewDesc("cypress_runs_total", "Total number of runs", labelsInOrder(RunsOrderedLabels), prometheus.Labels{}),

		CypressRunPassed:     prometheus.NewDesc("cypress_run_passed_total_last", "Total number of passed test per run processed ( latest value )", labelsInOrder(RunInstanceOrderedLabels), prometheus.Labels{}),
		CypressRunFailed:     prometheus.NewDesc("cypress_run_failed_total_last", "Total number of failed test per run processed ( latest value )", labelsInOrder(RunInstanceOrderedLabels), prometheus.Labels{}),
		CypressRunPending:    prometheus.NewDesc("cypress_run_pending_total_last", "Total number of pending test per run processed ( latest value )", labelsInOrder(RunInstanceOrderedLabels), prometheus.Labels{}),
		CypressRunSkipped:    prometheus.NewDesc("cypress_run_skipped_total_last", "Total number of skipped test per run processed ( latest value )", labelsInOrder(RunInstanceOrderedLabels), prometheus.Labels{}),
		CypressRunMutedTests: prometheus.NewDesc("cypress_run_muted_tests_total_last", "Total number of muted tests processed ( latest value )", labelsInOrder(RunInstanceOrderedLabels), prometheus.Labels{}),
		CypressRunFlakyTests: prometheus.NewDesc("cypress_run_flaky_tests_total_last", "Total number of flaky tests processed ( latest value )", labelsInOrder(RunInstanceOrderedLabels), prometheus.Labels{}),
		CypressRunDuration:   prometheus.NewDesc("cypress_run_duration_ms_last", " Duration of a processed run ( latest value )", labelsInOrder(RunInstanceOrderedLabels), prometheus.Labels{}),
		CypressRunStartTime:  prometheus.NewDesc("cypress_run_start_time_ms_last", "Start time of a processed run ( latest value )", labelsInOrder(RunInstanceOrderedLabels), prometheus.Labels{}),

		CypressRunPassedSum:     prometheus.NewDesc("cypress_run_passed_sum", "Total number of passed test per run processed ( summed value )", labelsInOrder(RunInstanceOrderedLabels), prometheus.Labels{}),
		CypressRunFailedSum:     prometheus.NewDesc("cypress_run_failed_sum", "Total number of failed test per run processed ( summed value )", labelsInOrder(RunInstanceOrderedLabels), prometheus.Labels{}),
		CypressRunPendingSum:    prometheus.NewDesc("cypress_run_pending_sum", "Total number of pending test per run processed ( summed value )", labelsInOrder(RunInstanceOrderedLabels), prometheus.Labels{}),
		CypressRunSkippedSum:    prometheus.NewDesc("cypress_run_skipped_sum", "Total number of skipped test per run processed ( summed value )", labelsInOrder(RunInstanceOrderedLabels), prometheus.Labels{}),
		CypressRunMutedTestsSum: prometheus.NewDesc("cypress_run_muted_tests_sum", "Total number of muted runs processed ( summed value )", labelsInOrder(RunInstanceOrderedLabels), prometheus.Labels{}),
		CypressRunFlakyTestsSum: prometheus.NewDesc("cypress_run_flaky_tests_sum", "Total number of flaky tests processed ( summed value )", labelsInOrder(RunInstanceOrderedLabels), prometheus.Labels{}),
		CypressRunDurationSum:   prometheus.NewDesc("cypress_run_duration_ms_sum", "Duration of a processed run ( summed value )", labelsInOrder(RunInstanceOrderedLabels), prometheus.Labels{}),
		CypressRunStartTimeSum:  prometheus.NewDesc("cypress_run_start_time_ms_sum", "Start time of a processed run ( summed value )", labelsInOrder(RunInstanceOrderedLabels), prometheus.Labels{}),

		CypressRunCount: prometheus.NewDesc("cypress_run_processed_sum", "Count of processed runs", labelsInOrder(RunInstanceOrderedLabels), prometheus.Labels{}),

		CypressTestStateLast:    prometheus.NewDesc("cypress_test_state_last", "Last state of a test ( filter with label `state` and check for value 1.0 )", labelsInOrder(TestResultInstanceOrderedLabels("")), prometheus.Labels{}),
		CypressTestDurationLast: prometheus.NewDesc("cypress_test_duration_ms_total_last", "Last duration of a test", labelsInOrder(TestInstanceOrderedLabels), prometheus.Labels{}),
		CypressTestStateSum:     prometheus.NewDesc("cypress_test_state_sum", "Summed state of a test ( filter with label `state` and check for value 1.0 )", labelsInOrder(TestResultInstanceOrderedLabels("")), prometheus.Labels{}),
		CypressTestDurationSum:  prometheus.NewDesc("cypress_test_duration_ms_total_sum", "Summed duration of a test", labelsInOrder(TestInstanceOrderedLabels), prometheus.Labels{}),
		CypressTestCount:        prometheus.NewDesc("cypress_test_processed_count", "Total number of processed tests", labelsInOrder(TestInstanceOrderedLabels), prometheus.Labels{}),

		CypressDashboardExporterAvailable: prometheus.NewDesc("cypress_dashboard_exporter_available", "Availability of CypressDashbboardExporter", labelsInOrder(RunsOrderedLabels), prometheus.Labels{}),

		cli:          &client,
		project:      project,
		LastDateTest: time.Date(2006, 1, 1, 1, 1, 1, 1, time.UTC),
		LastBuild:    0,

		AlreadyProcessedBuilds: set.NewIntSet(),
		runSummary: metricsmap.MetricMapSumValues{
			KeepUntil: time.Duration(time.Duration(keepUntil)),
		},
		testSummary: metricsmap.MetricMapSumValues{
			KeepUntil: time.Duration(time.Duration(keepUntil)),
		},

		runLatest: metricsmap.MetricMapKeepFirst{
			KeepUntil: time.Duration(time.Duration(keepUntil)),
		},
		testLatest: metricsmap.MetricMapKeepFirst{
			KeepUntil: time.Duration(time.Duration(keepUntil)),
		},

		firstRequest: true,
	}, nil

}

func (c *CypressDashboardCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.CypressRunsCount
	ch <- c.CypressRunPassed
	ch <- c.CypressRunFailed
	ch <- c.CypressRunPending
	ch <- c.CypressRunSkipped
	ch <- c.CypressRunMutedTests
	ch <- c.CypressRunDuration
	ch <- c.CypressRunFlakyTests
	ch <- c.CypressRunStartTime
	ch <- c.CypressRunPassedSum
	ch <- c.CypressRunFailedSum
	ch <- c.CypressRunPendingSum
	ch <- c.CypressRunSkippedSum
	ch <- c.CypressRunMutedTestsSum
	ch <- c.CypressRunDurationSum
	ch <- c.CypressRunFlakyTestsSum
	ch <- c.CypressRunStartTimeSum
	ch <- c.CypressRunCount
	ch <- c.CypressTestCount
	ch <- c.CypressTestStateSum
	ch <- c.CypressTestStateLast
	ch <- c.CypressTestDurationSum
	ch <- c.CypressTestDurationLast
	ch <- c.CypressDashboardExporterAvailable
}

// maybeMetric Send metric, if exist, to chanel. If value of metric is nil, or uncastable to float64, then print a warning or an error.
// Transform func will transform
func maybeMetric(ch chan<- prometheus.Metric, p *prometheus.Desc, valueType prometheus.ValueType, value interface{}, transformFunc metricTransformer, labels []string) {
	prometheusValue, err := converter.ConvertValueForPrometheus(value)

	if err != nil {
		logrus.Errorf("Skipping metric %v : %v", p.String(), err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		p,
		valueType,
		transformFunc(prometheusValue),
		labels...,
	)
}

func (c *CypressDashboardCollector) Collect(ch chan<- prometheus.Metric) {
	opts := cypressclient.EmptyMetricOptions()
	if c.firstRequest {
		backlog := 40
		logrus.Info("Processing the backlog of multiple requests")
		opts.Size = optional.NewOptionalInt(&backlog)
		c.firstRequest = false
	}
	// Set the project in the request
	opts.Project = c.project
	metrics, err := c.cli.GetMetrics(opts)

	if err != nil {
		logrus.Errorln("Error while scraping CypressDashboard metrics:", err)
		ch <- prometheus.NewInvalidMetric(c.CypressDashboardExporterAvailable, err)
		return
	}

	// Project level metrics
	maybeMetric(ch, c.CypressRunsCount, prometheus.GaugeValue, metrics.Data.Project.Runs.TotalCount, noopTransformer, evaluateLabels(RunsOrderedLabels, *metrics, nil))

	for _, runInstance := range metrics.Data.Project.Runs.Nodes.Reverse() {

		// First result => Latest build
		logrus.Infof("Processing build %v started at %v in state %v", runInstance.BuildNumber, runInstance.StartTime, runInstance.Status)
		if c.AlreadyProcessedBuilds.Has(runInstance.BuildNumber) {
			logrus.Infoln("Already processed build id", runInstance.BuildNumber)
		} else if runInstance.Status == "PASSED" || runInstance.Status == "FAILED" {
			logrus.Infoln("Processing build id", runInstance.BuildNumber)

			c.runLatest.Add(c.CypressRunPassed, runInstance.TotalPassed, evaluateLabels(RunInstanceOrderedLabels, *metrics, runInstance)...)
			c.runLatest.Add(c.CypressRunPending, runInstance.TotalPending, evaluateLabels(RunInstanceOrderedLabels, *metrics, runInstance)...)
			c.runLatest.Add(c.CypressRunFailed, runInstance.TotalFailed, evaluateLabels(RunInstanceOrderedLabels, *metrics, runInstance)...)
			c.runLatest.Add(c.CypressRunMutedTests, runInstance.TotalMutedTests, evaluateLabels(RunInstanceOrderedLabels, *metrics, runInstance)...)
			c.runLatest.Add(c.CypressRunSkipped, runInstance.TotalSkipped, evaluateLabels(RunInstanceOrderedLabels, *metrics, runInstance)...)
			c.runLatest.Add(c.CypressRunFlakyTests, runInstance.TotalFlakyTests, evaluateLabels(RunInstanceOrderedLabels, *metrics, runInstance)...)
			c.runLatest.Add(c.CypressRunDuration, runInstance.TotalDuration, evaluateLabels(RunInstanceOrderedLabels, *metrics, runInstance)...)
			c.runLatest.Add(c.CypressRunStartTime, runInstance.StartTime, evaluateLabels(RunInstanceOrderedLabels, *metrics, runInstance)...)
			// c.runLatest.Lock() // As soon as we processed the last build, we lock the map ( since latest build appears first in results )

			c.runSummary.Add(c.CypressRunPassedSum, runInstance.TotalPassed, evaluateLabels(RunInstanceOrderedLabels, *metrics, runInstance)...)
			c.runSummary.Add(c.CypressRunPendingSum, runInstance.TotalPending, evaluateLabels(RunInstanceOrderedLabels, *metrics, runInstance)...)
			c.runSummary.Add(c.CypressRunFailedSum, runInstance.TotalFailed, evaluateLabels(RunInstanceOrderedLabels, *metrics, runInstance)...)
			c.runSummary.Add(c.CypressRunMutedTestsSum, runInstance.TotalMutedTests, evaluateLabels(RunInstanceOrderedLabels, *metrics, runInstance)...)
			c.runSummary.Add(c.CypressRunSkippedSum, runInstance.TotalSkipped, evaluateLabels(RunInstanceOrderedLabels, *metrics, runInstance)...)
			c.runSummary.Add(c.CypressRunFlakyTestsSum, runInstance.TotalFlakyTests, evaluateLabels(RunInstanceOrderedLabels, *metrics, runInstance)...)
			c.runSummary.Add(c.CypressRunDurationSum, runInstance.TotalDuration, evaluateLabels(RunInstanceOrderedLabels, *metrics, runInstance)...)
			c.runSummary.Add(c.CypressRunStartTimeSum, runInstance.StartTime, evaluateLabels(RunInstanceOrderedLabels, *metrics, runInstance)...)

			//Count number of scraped runs
			c.runSummary.Add(c.CypressRunCount, 1.0, evaluateLabels(RunInstanceOrderedLabels, *metrics, runInstance)...)

			c.AlreadyProcessedBuilds.Add(runInstance.BuildNumber)

			for _, testInstance := range runInstance.TestResults.Nodes {
				state := testInstance.State

				c.testLatest.Add(c.CypressTestDurationLast, testInstance.Duration, evaluateLabels(TestInstanceOrderedLabels, *metrics, testContext{runInstance, testInstance})...)

				matched := false
				for _, value := range cypressclient.AllValidState() {
					s := promValueFromState(state, value.String())
					if s == 1.0 {
						matched = true
					}
					c.testSummary.Add(c.CypressTestStateSum, s, evaluateLabels(TestResultInstanceOrderedLabels(value.String()), *metrics, testContext{runInstance, testInstance})...)
					c.testLatest.Add(c.CypressTestStateLast, s, evaluateLabels(TestResultInstanceOrderedLabels(value.String()), *metrics, testContext{runInstance, testInstance})...)
				}
				if !matched {
					logrus.Warnln("Unknown state", state, " while processing test", testInstance.TitleParts)
					c.testSummary.Add(c.CypressTestStateSum, 1.0, evaluateLabels(TestResultInstanceOrderedLabels(cypressclient.Other.String()), *metrics, testContext{runInstance, testInstance})...)
					c.testLatest.Add(c.CypressTestStateLast, 1.0, evaluateLabels(TestResultInstanceOrderedLabels(cypressclient.Other.String()), *metrics, testContext{runInstance, testInstance})...)
				} else {
					c.testSummary.Add(c.CypressTestStateSum, 0.0, evaluateLabels(TestResultInstanceOrderedLabels(cypressclient.Other.String()), *metrics, testContext{runInstance, testInstance})...)
					c.testLatest.Add(c.CypressTestStateLast, 0.0, evaluateLabels(TestResultInstanceOrderedLabels(cypressclient.Other.String()), *metrics, testContext{runInstance, testInstance})...)
				}

				c.testSummary.Add(c.CypressTestDurationSum, testInstance.Duration, evaluateLabels(TestInstanceOrderedLabels, *metrics, testContext{runInstance, testInstance})...)
				c.testSummary.Add(c.CypressTestCount, 1.0, evaluateLabels(TestInstanceOrderedLabels, *metrics, testContext{runInstance, testInstance})...)
			}
			logrus.Debugln("Map of tests and runs : %+v\n%+v\n%+v\n%+v", c.runLatest, c.runSummary, c.testLatest, c.testSummary)

		} else {
			logrus.Infof("Run %v is in state %v, skipping for now...", runInstance.BuildNumber)
		}

	}
	for key, value := range c.runSummary.Map() {
		logrus.Debugln("Processing summary ( counters )", key.Prom.String())
		maybeMetric(ch, key.Prom, prometheus.CounterValue, value.Value, noopTransformer, value.Labels)
	}
	for key, value := range c.runLatest.Map() {
		logrus.Debugln("Processing latests ( gauge )", key.Prom.String())
		maybeMetric(ch, key.Prom, prometheus.GaugeValue, value.Value, noopTransformer, value.Labels)
	}

	for key, value := range c.testSummary.Map() {
		logrus.Debugln("Processing tests summary ( counters )", key.Prom.String())
		maybeMetric(ch, key.Prom, prometheus.CounterValue, value.Value, noopTransformer, value.Labels)
	}
	for key, value := range c.testLatest.Map() {
		logrus.Debugln("Processing test latests ( gauge )", key.Prom.String())
		maybeMetric(ch, key.Prom, prometheus.GaugeValue, value.Value, noopTransformer, value.Labels)
	}
}
