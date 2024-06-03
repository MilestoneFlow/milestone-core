package publicapi

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"milestone_core/apiclient"
	"milestone_core/flow"
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

func startFlow() flow.Flow {
	return flow.Flow{
		ID:          deterministicObjectID(1),
		WorkspaceID: WorkspaceID(),
		Name:        "Starter Flow",
		BaseURL:     "testBaseURL",
		Segments:    []flow.Segment{},
		Steps: []flow.Step{
			{
				StepID: "step_1",
				Data: flow.StepData{
					TargetUrl:   "/dashboard",
					ElementType: flow.StepElementTypeTooltip,
				},
				Opts: flow.StepOpts{
					IsSource: true,
				},
				ParentNodeId: "",
			},
		},
		Opts: flow.Opts{
			ThemeColor:      "blue",
			ElementTemplate: flow.StepElementTemplateDark,
		},
		Live: true,
	}
}

func generalFlowDependentOnStarter() flow.Flow {
	newId := deterministicObjectID(2)
	return flow.Flow{
		ID:          newId,
		WorkspaceID: WorkspaceID(),
		Name:        "Second Flow",
		BaseURL:     "testBaseURL",
		Segments:    []flow.Segment{},
		Steps: []flow.Step{
			{
				StepID: "step_1",
				Data: flow.StepData{
					TargetUrl:   "/dashboard",
					ElementType: flow.StepElementTypeTooltip,
				},
				Opts: flow.StepOpts{
					IsSource: true,
				},
				ParentNodeId: "",
			},
		},
		Opts: flow.Opts{
			ThemeColor:      "blue",
			ElementTemplate: flow.StepElementTemplateDark,
			DependsOn:       []string{startFlow().ID.Hex()},
		},
		Live: true,
	}
}

func secondGeneralFlowDependentOnStarter() flow.Flow {
	newId := deterministicObjectID(3)
	return flow.Flow{
		ID:          newId,
		WorkspaceID: WorkspaceID(),
		Name:        "Third Flow",
		BaseURL:     "testBaseURL",
		Segments:    []flow.Segment{},
		Steps: []flow.Step{
			{
				StepID: "step_1",
				Data: flow.StepData{
					TargetUrl:   "/dashboard",
					ElementType: flow.StepElementTypePopup,
				},
				Opts: flow.StepOpts{
					IsSource: true,
				},
				ParentNodeId: "",
			},
		},
		Opts: flow.Opts{
			ThemeColor:      "blue",
			ElementTemplate: flow.StepElementTemplateDark,
			DependsOn:       []string{startFlow().ID.Hex()},
		},
		Live: true,
	}
}

func thirdGeneralFlowDependentOnStarterAndTargeted() flow.Flow {
	newId := deterministicObjectID(4)
	return flow.Flow{
		ID:          newId,
		WorkspaceID: WorkspaceID(),
		Name:        "Targeted General Flow",
		BaseURL:     "testBaseURL",
		Segments:    []flow.Segment{},
		Steps: []flow.Step{
			{
				StepID: "step_1",
				Data: flow.StepData{
					TargetUrl:   "/dashboard",
					ElementType: flow.StepElementTypePopup,
				},
				Opts: flow.StepOpts{
					IsSource: true,
				},
				ParentNodeId: "",
			},
		},
		Opts: flow.Opts{
			ThemeColor:      "blue",
			ElementTemplate: flow.StepElementTemplateDark,
			DependsOn:       []string{startFlow().ID.Hex()},
			Targeting: flow.Targeting{
				Rules: []flow.TargetingRule{
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

func flowDependentOnGeneralFirstAndSecond() flow.Flow {
	newId := deterministicObjectID(5)
	return flow.Flow{
		ID:          newId,
		WorkspaceID: WorkspaceID(),
		Name:        "Dependent Flow On First and Second General",
		BaseURL:     "testBaseURL",
		Segments:    []flow.Segment{},
		Steps: []flow.Step{
			{
				StepID: "step_1",
				Data: flow.StepData{
					TargetUrl:   "/dashboard",
					ElementType: flow.StepElementTypePopup,
				},
				Opts: flow.StepOpts{
					IsSource: true,
				},
				ParentNodeId: "",
			},
		},
		Opts: flow.Opts{
			ThemeColor:      "blue",
			ElementTemplate: flow.StepElementTemplateDark,
			DependsOn:       []string{generalFlowDependentOnStarter().ID.Hex(), secondGeneralFlowDependentOnStarter().ID.Hex()},
		},
		Live: true,
	}
}
