package users

type EnrolledUser struct {
	Id             string `json:"id" bson:"_id"`
	ExternalId     string `json:"externalId" bson:"externalId"`
	Email          string `json:"email,omitempty" bson:"email,omitempty"`
	OrganizationId string `json:"organizationId,omitempty" bson:"organizationId,omitempty"`
}

type UserState struct {
	UserId                string          `json:"userId" bson:"userId"`
	CurrentEnrolledFlowId string          `json:"currentEnrolledFlowId" bson:"currentEnrolledFlowId"`
	CurrentStepId         string          `json:"currentStepId" bson:"currentStepId"`
	Times                 []UserStateTime `json:"times" bson:"times"`
	Created               int64           `json:"created" bson:"created"`
}

type UserStateTime struct {
	StepId         string `json:"stepId" bson:"stepId"`
	StartTimestamp int64  `json:"startTimestamp" bson:"startTimestamp"`
	EndTimestamp   int64  `json:"endTimestamp" bson:"endTimestamp"`
}
