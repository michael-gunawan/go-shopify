package goshopify

import (
	"errors"
	"fmt"
)

type MetafieldDefinitionService interface {
	List(ownerType OwnerType, namespace string) ([]MetafieldDefinition, error)
	Create(MetafieldDefinition) (*MetafieldDefinition, error)
}

type MetafieldDefinitionServiceOp struct {
	client     *Client
	resource   string
	resourceID int64
}

type MetafieldDefinitionsQueryResponse struct {
	Data struct {
		MetafieldDefinitions struct {
			Edges []struct {
				Node MetafieldDefinition `json:"node"`
			}
		}
	}
}

type MetafieldDefinitionsMutationResponse struct {
	Data struct {
		MetafieldDefinitionCreate struct {
			CreatedDefinition MetafieldDefinition
		}
	}
}

type MetafieldDefinition struct {
	Id          string                `json:"id,omitempty"`
	Name        string                `json:"name,omitempty"`
	Description string                `json:"description,omitempty"`
	Namespace   string                `json:"namespace,omitempty"`
	Key         string                `json:"key,omitempty"`
	OwnerType   string                `json:"ownerType,omitempty"`
	Type        string                `json:"type,omitempty"` // Note: type of [Type] is different on query vs mutation.
	Validations []MetafieldValidation `json:"validations,omitempty"`

	// Mutation only
	Pin                    bool `json:"pin,omitempty"`
	VisibleToStorefrontApi bool `json:"visibleToStorefrontApi,omitempty"`
}

type MetafieldType struct {
	Name     string `json:"name,omitempty"`
	Category string `json:"category,omitempty"`
}

type MetafieldValidation struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

type OwnerType string

const (
	OwnerTypeAny     OwnerType = ""
	OwnerTypeProduct OwnerType = "PRODUCT"
)

// List metafields
func (s *MetafieldDefinitionServiceOp) List(ownerType OwnerType, namespace string) ([]MetafieldDefinition, error) {
	ownerTypeStr := ""
	if len(ownerType) > 0 {
		ownerTypeStr = fmt.Sprintf(", ownerType: %s", ownerType)
	}
	namespaceStr := ""
	if len(namespace) > 0 {
		namespaceStr = fmt.Sprintf(", namespace: \"%s\"", namespace)
	}
	query := &GraphQLQuery{
		Query: fmt.Sprintf("{ metafieldDefinitions(first: 100%s%s) { edges { node { id name namespace key validations { name type value } } } } }", ownerTypeStr, namespaceStr),
	}
	path := "/graphql.json"
	response := new(MetafieldDefinitionsQueryResponse)
	err := s.client.Post(path, query, response)

	definitions := make([]MetafieldDefinition, len(response.Data.MetafieldDefinitions.Edges))
	for index, f := range response.Data.MetafieldDefinitions.Edges {
		definitions[index] = f.Node
	}

	return definitions, err
}

// Create a new metafield
func (s *MetafieldDefinitionServiceOp) Create(metafield MetafieldDefinition) (*MetafieldDefinition, error) {
	query := &GraphQLQuery{
		Query:     "mutation CreateMetafieldDefinition($definition: MetafieldDefinitionInput!) { metafieldDefinitionCreate(definition: $definition) { createdDefinition { id name } userErrors { field message code } } }",
		Variables: GraphQLQueryVariables{Definition: metafield},
	}

	path := "/graphql.json"
	response := new(MetafieldDefinitionsMutationResponse)
	err := s.client.Post(path, query, response)

	createdDefinition := &response.Data.MetafieldDefinitionCreate.CreatedDefinition
	if len(createdDefinition.Id) == 0 {
		return nil, errors.New(fmt.Sprintf("Unable to create metafield definition %v", metafield.Name))
	}

	return createdDefinition, err
}
