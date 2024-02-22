package users

type EnrolledUser struct {
	Id             string `json:"id" bson:"_id"`
	ExternalId     string `json:"externalId" bson:"externalId"`
	Email          string `json:"email,omitempty" bson:"email,omitempty"`
	OrganizationId string `json:"organizationId,omitempty" bson:"organizationId,omitempty"`
}
