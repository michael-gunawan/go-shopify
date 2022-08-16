package goshopify

type GraphQLQuery struct {
	Query     string                `json:"query"`
	Variables GraphQLQueryVariables `json:"variables,omitempty"`
}

type GraphQLQueryVariables struct {
	Definition interface{} `json:"definition"`
}
