package enum

type ResourceCode string

type ResourceType struct {
	User ResourceCode
	Role ResourceCode
}

var Resource = ResourceType{
	User: "USER",
	Role: "ROLE",
}
