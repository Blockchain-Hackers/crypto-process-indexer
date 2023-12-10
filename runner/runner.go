package runner

import (
	"fmt"
	"time"

	"github.com/blockchain-hackers/indexer/database"
	"github.com/blockchain-hackers/indexer/functions"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Run(flow database.Workflow) {
	var steps []database.StepRun
	for _, step := range flow.Steps {
		fmt.Printf("Running step: %+v\n", step.Name)
		// run the step
		fmt.Println("Function: ", step.Function)

		resp, err := functions.CallFunc(step.Function, functions.ConvertDBParamsToFunctionParams(step.Parameters, step.Name))
		if err.Exists() {
			fmt.Println("Error: ", err)
			steps = append(steps, functions.ConvertFunctionErrorToDBStep(err))
			break
		} else {
			fmt.Println("Response: ", resp)
			// save result to flow runs
			steps = append(steps, functions.ConvertFunctionResponseToDBStep(resp))
		}
	}
	// save the run to flow runs
	run := database.FlowRun{
		ID:		primitive.NewObjectID(),
		FlowID:    flow.ID,
		Trigger:   flow.Trigger,
		Steps:     steps,
		V:         0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	database.WriteRunToFlow(flow.ID, run)
}
