package cypressclient

type state string

func (s state) String() string {
	return string(s)
}

const Passed state = "PASSED"
const Failed state = "FAILED"
const Canceled state = "CANCELLED"
const Skipped state = "SKIPPED"
const Other state = "OTHER"

var allValidState []state = []state{
	Passed,
	Failed,
	Canceled,
	Skipped,
}

func AllValidState() []state {
	return allValidState
}
