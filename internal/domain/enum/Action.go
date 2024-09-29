package enum

type ActionCode string

type ActionType struct {
	Create ActionCode
	View   ActionCode
	Update ActionCode
	Delete ActionCode
}

var Action = ActionType{
	Create: "CREATE", // Correct initialization without 'ActionCode ='
	View:   "VIEW",
	Update: "UPDATE",
	Delete: "DELETE",
}
