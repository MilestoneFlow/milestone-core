package apigateway

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"milestone_core/identity/apiclient"
	"milestone_core/tours/flows"
	"time"
)

func SetupMockData(mongoConnection *mongo.Database) {
	mongoConnection.Collection("flows").Drop(context.Background())
	mongoConnection.Collection("api_clients").Drop(context.Background())

	mongoConnection.Collection("api_clients").InsertOne(context.Background(), apiclient.ApiClient{
		Token:       ApiToken(),
		WorkspaceID: WorkspaceID(),
	})

	mongoConnection.Collection("flows").InsertMany(context.Background(), []interface{}{
		startFlow(),
		generalFlowDependentOnStarter(),
		secondGeneralFlowDependentOnStarter(),
		thirdGeneralFlowDependentOnStarterAndTargeted(),
		flowDependentOnGeneralFirstAndSecond(),
	})
}

func CleanupMockData(mongoConnection *mongo.Database) {
	mongoConnection.Collection("flows").Drop(context.Background())
	mongoConnection.Collection("api_clients").Drop(context.Background())
}

func WorkspaceID() string {
	return "testWorkspaceId"
}

func ApiToken() string {
	return "token"
}

func deterministicObjectID(unixTimestamp int64) primitive.ObjectID {
	t := time.Unix(unixTimestamp, 0)
	return primitive.NewObjectIDFromTimestamp(t)
}

func startFlow() flows.Flow {
	return flows.Flow{
		ID:          deterministicObjectID(1),
		WorkspaceID: WorkspaceID(),
		Name:        "Starter Flow",
		BaseURL:     "testBaseURL",
		Segments:    []flows.Segment{},
		Steps: []flows.Step{
			{
				StepID: "step_1",
				Data: flows.StepData{
					TargetUrl:   "/dashboard",
					ElementType: flows.StepElementTypeTooltip,
				},
				Opts: flows.StepOpts{
					IsSource: true,
				},
				ParentNodeId: "",
			},
		},
		Opts: flows.Opts{
			ThemeColor:      "blue",
			ElementTemplate: flows.StepElementTemplateDark,
		},
		Live: true,
	}
}

func generalFlowDependentOnStarter() flows.Flow {
	newId := deterministicObjectID(2)
	return flows.Flow{
		ID:          newId,
		WorkspaceID: WorkspaceID(),
		Name:        "Second Flow",
		BaseURL:     "testBaseURL",
		Segments:    []flows.Segment{},
		Steps: []flows.Step{
			{
				StepID: "step_1",
				Data: flows.StepData{
					TargetUrl:   "/dashboard",
					ElementType: flows.StepElementTypeTooltip,
				},
				Opts: flows.StepOpts{
					IsSource: true,
				},
				ParentNodeId: "",
			},
		},
		Opts: flows.Opts{
			ThemeColor:      "blue",
			ElementTemplate: flows.StepElementTemplateDark,
			DependsOn:       []string{startFlow().ID.Hex()},
		},
		Live: true,
	}
}

func secondGeneralFlowDependentOnStarter() flows.Flow {
	newId := deterministicObjectID(3)
	return flows.Flow{
		ID:          newId,
		WorkspaceID: WorkspaceID(),
		Name:        "Third Flow",
		BaseURL:     "testBaseURL",
		Segments:    []flows.Segment{},
		Steps: []flows.Step{
			{
				StepID: "step_1",
				Data: flows.StepData{
					TargetUrl:   "/dashboard",
					ElementType: flows.StepElementTypePopup,
				},
				Opts: flows.StepOpts{
					IsSource: true,
				},
				ParentNodeId: "",
			},
		},
		Opts: flows.Opts{
			ThemeColor:      "blue",
			ElementTemplate: flows.StepElementTemplateDark,
			DependsOn:       []string{startFlow().ID.Hex()},
		},
		Live: true,
	}
}

func thirdGeneralFlowDependentOnStarterAndTargeted() flows.Flow {
	newId := deterministicObjectID(4)
	return flows.Flow{
		ID:          newId,
		WorkspaceID: WorkspaceID(),
		Name:        "Targeted General Flow",
		BaseURL:     "testBaseURL",
		Segments:    []flows.Segment{},
		Steps: []flows.Step{
			{
				StepID: "step_1",
				Data: flows.StepData{
					TargetUrl:   "/dashboard",
					ElementType: flows.StepElementTypePopup,
				},
				Opts: flows.StepOpts{
					IsSource: true,
				},
				ParentNodeId: "",
			},
		},
		Opts: flows.Opts{
			ThemeColor:      "blue",
			ElementTemplate: flows.StepElementTemplateDark,
			DependsOn:       []string{startFlow().ID.Hex()},
			Targeting: flows.Targeting{
				Rules: []flows.TargetingRule{
					{
						Condition: "elapsed_time",
						Value:     "7d",
					},
				},
			},
		},
		Live: true,
	}
}

func flowDependentOnGeneralFirstAndSecond() flows.Flow {
	newId := deterministicObjectID(5)
	return flows.Flow{
		ID:          newId,
		WorkspaceID: WorkspaceID(),
		Name:        "Dependent Flow On First and Second General",
		BaseURL:     "testBaseURL",
		Segments:    []flows.Segment{},
		Steps: []flows.Step{
			{
				StepID: "step_1",
				Data: flows.StepData{
					TargetUrl:   "/dashboard",
					ElementType: flows.StepElementTypePopup,
				},
				Opts: flows.StepOpts{
					IsSource: true,
				},
				ParentNodeId: "",
			},
		},
		Opts: flows.Opts{
			ThemeColor:      "blue",
			ElementTemplate: flows.StepElementTemplateDark,
			DependsOn:       []string{generalFlowDependentOnStarter().ID.Hex(), secondGeneralFlowDependentOnStarter().ID.Hex()},
		},
		Live: true,
	}
}
