package enum

type ActionCode string

type ActionType struct {
	Add    ActionCode
	View   ActionCode
	Update ActionCode
	Delete ActionCode
}

var Action = ActionType{
	Add:    "ADD", // Correct initialization without 'ActionCode ='
	View:   "VIEW",
	Update: "UPDATE",
	Delete: "DELETE",
}
