package cypresscollector

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rguilmont/cypress-dashboard-exporter/pkg/cypressclient"
)

type labelsEvaluatorImpl struct {
	label     func() string
	evaluator func(cypressclient.StatsFromCypressDashboard, interface{}) string // Second argument is context specific
}

var RunsOrderedLabels []labelsEvaluatorImpl = []labelsEvaluatorImpl{
	{
		label: func() string { return "project_id" },
		evaluator: func(s cypressclient.StatsFromCypressDashboard, _ interface{}) string {
			return s.Data.Project.ID
		},
	},
	{
		label: func() string { return "project_name" },
		evaluator: func(s cypressclient.StatsFromCypressDashboard, _ interface{}) string {
			return s.Data.Project.Name
		},
	},
	{
		label: func() string { return "is_using_retries" },
		evaluator: func(s cypressclient.StatsFromCypressDashboard, _ interface{}) string {
			return func() string {
				if s.Data.Project.IsUsingRetries {
					return "1"
				}
				return "0"
			}()
		},
	},
}

func castInterfaceToRunContext(i interface{}) cypressclient.RunResult {
	res, ok := i.(cypressclient.RunResult)
	if ok {
		return res
	}
	panic(fmt.Sprintf("Can't cast into RunResult %v - %v", reflect.TypeOf(i), i))
}

var RunInstanceOrderedLabels []labelsEvaluatorImpl = []labelsEvaluatorImpl{
	{
		func() string { return "project_id" },
		func(s cypressclient.StatsFromCypressDashboard, _ interface{}) string {
			return s.Data.Project.ID
		},
	},
	{
		func() string { return "project_name" },
		func(s cypressclient.StatsFromCypressDashboard, _ interface{}) string {
			return s.Data.Project.Name
		},
	},
	{
		func() string { return "is_using_retries" },
		func(s cypressclient.StatsFromCypressDashboard, _ interface{}) string {
			return func() string {
				if s.Data.Project.IsUsingRetries {
					return "1"
				}
				return "0"
			}()
		},
	},
	{
		func() string { return "ci_provider" },
		func(_ cypressclient.StatsFromCypressDashboard, i interface{}) string {
			ctx := castInterfaceToRunContext(i)
			return ctx.Ci.Provider
		},
	},
	{
		func() string { return "os_name" },
		func(_ cypressclient.StatsFromCypressDashboard, i interface{}) string {
			ctx := castInterfaceToRunContext(i)
			if len(ctx.TestResults.Nodes) > 0 {
				return fmt.Sprintf("%v_%v", ctx.TestResults.Nodes[0].Instance.Os.Name, ctx.TestResults.Nodes[0].Instance.Os.Version)
			}
			return "unknown"
		},
	},
	{
		func() string { return "browser_name" },
		func(_ cypressclient.StatsFromCypressDashboard, i interface{}) string {
			ctx := castInterfaceToRunContext(i)
			if len(ctx.TestResults.Nodes) > 0 {
				return fmt.Sprintf("%v_%v", ctx.TestResults.Nodes[0].Instance.Browser.Name, ctx.TestResults.Nodes[0].Instance.Browser.Version)
			}
			return "unknown"
		},
	},
	{
		func() string { return "git_branch" },
		func(_ cypressclient.StatsFromCypressDashboard, i interface{}) string {
			ctx := castInterfaceToRunContext(i)
			return ctx.Commit.Branch
		},
	},
}

type testContext struct {
	runResult  cypressclient.RunResult
	testResult cypressclient.TestResult
}

func castInterfaceToTestContext(i interface{}) testContext {
	res, ok := i.(testContext)
	if ok {
		return res
	}
	panic(fmt.Sprintf("Can't cast into RunResult %v - %v", reflect.TypeOf(i), i))
}

var TestInstanceOrderedLabels []labelsEvaluatorImpl = []labelsEvaluatorImpl{
	{
		func() string { return "project_id" },
		func(s cypressclient.StatsFromCypressDashboard, _ interface{}) string {
			return s.Data.Project.ID
		},
	},
	{
		func() string { return "project_name" },
		func(s cypressclient.StatsFromCypressDashboard, _ interface{}) string {
			return s.Data.Project.Name
		},
	},
	{
		func() string { return "is_using_retries" },
		func(s cypressclient.StatsFromCypressDashboard, _ interface{}) string {
			return func() string {
				if s.Data.Project.IsUsingRetries {
					return "1"
				}
				return "0"
			}()
		},
	},
	{
		func() string { return "ci_provider" },
		func(s cypressclient.StatsFromCypressDashboard, i interface{}) string {
			ctx := castInterfaceToTestContext(i)
			return ctx.runResult.Ci.Provider
		},
	},
	{
		func() string { return "os_name" },
		func(_ cypressclient.StatsFromCypressDashboard, i interface{}) string {
			ctx := castInterfaceToTestContext(i)
			return fmt.Sprintf("%v_%v", ctx.testResult.Instance.Os.Name, ctx.testResult.Instance.Os.Version)
		},
	},
	{
		func() string { return "browser_name" },
		func(_ cypressclient.StatsFromCypressDashboard, i interface{}) string {
			ctx := castInterfaceToTestContext(i)
			return fmt.Sprintf("%v_%v", ctx.testResult.Instance.Browser.Name, ctx.testResult.Instance.Browser.Version)
		},
	},
	{
		func() string { return "spec_file" },
		func(_ cypressclient.StatsFromCypressDashboard, i interface{}) string {
			ctx := castInterfaceToTestContext(i)
			return ctx.testResult.Instance.Spec.ShortPath
		},
	},
	{
		func() string { return "name" },
		func(_ cypressclient.StatsFromCypressDashboard, i interface{}) string {
			ctx := castInterfaceToTestContext(i)
			return strings.Join(ctx.testResult.TitleParts, " ")
		},
	},
	{
		func() string { return "git_branch" },
		func(_ cypressclient.StatsFromCypressDashboard, i interface{}) string {
			ctx := castInterfaceToTestContext(i)
			return ctx.runResult.Commit.Branch
		},
	},
	{
		func() string { return "run_group" },
		func(_ cypressclient.StatsFromCypressDashboard, i interface{}) string {
			ctx := castInterfaceToTestContext(i)
			return ctx.testResult.Instance.Group.Name
		},
	},
}

func TestResultInstanceOrderedLabels(state string) []labelsEvaluatorImpl {
	return []labelsEvaluatorImpl{
		{
			func() string { return "project_id" },
			func(s cypressclient.StatsFromCypressDashboard, _ interface{}) string {
				return s.Data.Project.ID
			},
		},
		{
			func() string { return "project_name" },
			func(s cypressclient.StatsFromCypressDashboard, _ interface{}) string {
				return s.Data.Project.Name
			},
		},
		{
			func() string { return "is_using_retries" },
			func(s cypressclient.StatsFromCypressDashboard, _ interface{}) string {
				return func() string {
					if s.Data.Project.IsUsingRetries {
						return "1"
					}
					return "0"
				}()
			},
		},
		{
			func() string { return "ci_provider" },
			func(_ cypressclient.StatsFromCypressDashboard, i interface{}) string {
				ctx := castInterfaceToTestContext(i)
				return ctx.runResult.Ci.Provider
			},
		},
		{
			func() string { return "os_name" },
			func(_ cypressclient.StatsFromCypressDashboard, i interface{}) string {
				ctx := castInterfaceToTestContext(i)
				return fmt.Sprintf("%v_%v", ctx.testResult.Instance.Os.Name, ctx.testResult.Instance.Os.Version)
			},
		},
		{
			func() string { return "browser_name" },
			func(_ cypressclient.StatsFromCypressDashboard, i interface{}) string {
				ctx := castInterfaceToTestContext(i)
				return fmt.Sprintf("%v_%v", ctx.testResult.Instance.Browser.Name, ctx.testResult.Instance.Browser.Version)
			},
		},
		{
			func() string { return "spec_file" },
			func(_ cypressclient.StatsFromCypressDashboard, i interface{}) string {
				ctx := castInterfaceToTestContext(i)
				return ctx.testResult.Instance.Spec.ShortPath
			},
		},
		{
			func() string { return "name" },
			func(_ cypressclient.StatsFromCypressDashboard, i interface{}) string {
				ctx := castInterfaceToTestContext(i)
				return strings.Join(ctx.testResult.TitleParts, " ")
			},
		},
		{
			func() string { return "git_branch" },
			func(_ cypressclient.StatsFromCypressDashboard, i interface{}) string {
				ctx := castInterfaceToTestContext(i)
				return ctx.runResult.Commit.Branch
			},
		},
		{
			func() string { return "run_group" },
			func(_ cypressclient.StatsFromCypressDashboard, i interface{}) string {
				ctx := castInterfaceToTestContext(i)
				return ctx.testResult.Instance.Group.Name
			},
		},
		{
			func() string { return "state" },
			func(s cypressclient.StatsFromCypressDashboard, i interface{}) string {
				return state
			},
		},
	}
}

func labelsInOrder(evaluators []labelsEvaluatorImpl) []string {
	res := []string{}
	for _, e := range evaluators {
		res = append(res, e.label())
	}
	return res
}

func evaluateLabels(evaluators []labelsEvaluatorImpl, s cypressclient.StatsFromCypressDashboard, ctx interface{}) []string {
	res := []string{}
	for _, e := range evaluators {
		res = append(res, e.evaluator(s, ctx))
	}
	return res
}
