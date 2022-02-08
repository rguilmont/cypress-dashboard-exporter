# Cypress.io dashboard Prometheus exporter

See an example of what is possible [here](https://rguilmont.net/grafana-example/d/scjkVna7k/cypress-dashboard)

Prometheus exporter for a project from Cypress.io dashboards, giving the ability to alert, make special operations or correlate with other datasources by using grafana.

# Installation

There's no public docker image available yet, but you can simply build it :

```
docker build -t cypress-exporter .
```

or the binary :

```
cd cmd && go build
```

## Usage

```
  -debug
        activate debug logging
  -email string
        email to connect to the dashboard
  -keepUntil int
        Time ( in days ) to keep in memory the results of a test/run before removing it. (default 14)
  -listen string
        host:port to listen (default "0.0.0.0:8081")
  -password string
        password to connect to the dashboard
  -project string
        host:port to listen (default "7s5okt")
```

# Available metrics

| Metric                               | Description                                                                  |
| ------------------------------------ | ---------------------------------------------------------------------------- |
| cypress_run_passed_total_last        | Total number of passed test per run processed ( latest value )               |
| cypress_run_failed_total_last        | Total number of failed test per run processed ( latest value )               |
| cypress_run_pending_total_last       | Total number of pending test per run processed ( latest value )              |
| cypress_run_skipped_total_last       | Total number of skipped test per run processed ( latest value )              |
| cypress_run_muted_tests_total_last   | Total number of muted tests processed ( latest value )                       |
| cypress_run_flaky_tests_total_last   | Total number of flaky tests processed ( latest value )                       |
| cypress_run_duration_ms_last         | Duration of a processed run ( latest value )                                 |
| cypress_run_start_time_ms_last       | Start time of a processed run ( latest value )                               |
| cypress_run_passed_sum               | Total number of passed test per run processed ( summed value )               |
| cypress_run_failed_sum               | Total number of failed test per run processed ( summed value )               |
| cypress_run_pending_sum              | Total number of pending test per run processed ( summed value )              |
| cypress_run_skipped_sum              | Total number of skipped test per run processed ( summed value )              |
| cypress_run_muted_tests_sum          | Total number of muted runs processed ( summed value )                        |
| cypress_run_flaky_tests_sum          | Total number of flaky tests processed ( summed value )                       |
| cypress_run_duration_ms_sum          | Duration of a processed run ( summed value )                                 |
| cypress_run_start_time_ms_sum        | Start time of a processed run ( summed value )                               |
| cypress_run_processed_sum            | Count of processed runs                                                      |
| cypress_test_state_last              | Last state of a test ( filter with label `state` and check for value 1.0 )   |
| cypress_test_duration_ms_total_last  | Last duration of a test                                                      |
| cypress_test_state_sum               | Summed state of a test ( filter with label `state` and check for value 1.0 ) |
| cypress_test_duration_ms_total_sum   | Summed duration of a test                                                    |
| cypress_test_processed_count         | Total number of processed tests                                              |
| cypress_dashboard_exporter_available | Availability of CypressDashbboardExporter                                    |

## Labels

For `run` related metrics, the following labels are exposed :

- browser_name
- ci_provider
- git_branch
- is_using_retries
- os_name
- project_id
- project_name

For `test` related metrics :

- browser_name
- ci_provider
- git_branch
- is_using_retries
- name
- os_name
- project_id
- project_name
- run_group
- spec_file

For `cypress_test_state_last` and there's one label `state` with each possible value `CANCELED` `FAILED` `PASSED` `SKIPPED` or `OTHER`. Value of the metric will be 1.0 ( or incremented in case of the sum one ) when it's the corresponding state, and 0 ( or not incremented ) if not.

# Grafana dashboard

There's a grafana dashboard JSON file available in the `grafana` dashboard that you can import.
