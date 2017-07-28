package service

// Raw struct used to unmarshal json from go trace
type Raw struct {
	StackFrames map[string]struct {
		Name   string `json:"name"`
		Parent int    `json:"parent"`
	} `json:"stackFrames"`
}

// WorkFlowData used to group work flow functions
// type WorkFlowData struct {
// 	WorkFlowNumber int
// 	Data           []Raw
// }

//type workFlow map[string][]Raw
