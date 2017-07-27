package service

type Raw struct {
	StackFrames map[string]struct {
		Name   string `json:"name"`
		Parent int `json:"parent"`
	} `json:"stackFrames"`
}

type WorkFlowData struct {
	WorkFlowNumber int
	Data           []Raw
}

type workFlow map[string][]Raw