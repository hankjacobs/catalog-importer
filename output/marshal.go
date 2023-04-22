package output

import (
	"context"
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/incident-io/catalog-importer/client"
	"github.com/incident-io/catalog-importer/expr"
	"github.com/incident-io/catalog-importer/source"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type CatalogTypeModel struct {
	Name        string
	Description string
	TypeName    string
	Attributes  []client.CatalogTypeAttributePayloadV2
}

type CatalogEntryModel struct {
	ExternalID      string
	Name            string
	Aliases         []string
	AttributeValues map[string]client.CatalogAttributeBindingPayloadV2
}

// Marshal builds payloads to configure both catalog type and the entries for that type
// from the output configuration and entries that have already been filtered.
//
// The majority of the work comes from compiling and evaluating the CEL expressions that
// marshal the catalog entries from source.
func MarshalType(output *Output) *CatalogTypeModel {
	catalogTypeModel := &CatalogTypeModel{
		Name:        output.Name,
		Description: output.Description,
		TypeName:    output.TypeName,
		Attributes:  []client.CatalogTypeAttributePayloadV2{},
	}
	for _, attr := range output.Attributes {
		catalogTypeModel.Attributes = append(
			catalogTypeModel.Attributes, client.CatalogTypeAttributePayloadV2{
				Id:    lo.ToPtr(attr.ID),
				Name:  attr.Name,
				Type:  attr.Type,
				Array: attr.Array,
			})
	}

	return catalogTypeModel
}

// MarshalEntries builds payloads to for the entries of the given output, assuming those
// entries have already been filtered.
//
// The majority of the work comes from compiling and evaluating the CEL expressions that
// marshal the catalog entries from source.
func MarshalEntries(ctx context.Context, output *Output, entries []source.Entry) ([]*CatalogEntryModel, error) {
	nameProgram, err := expr.Compile(output.Source.Name)
	if err != nil {
		return nil, errors.Wrap(err, "source.name")
	}

	externalIDProgram, err := expr.Compile(output.Source.ExternalID)
	if err != nil {
		return nil, errors.Wrap(err, "source.external_id")
	}

	aliasPrograms := []cel.Program{}
	for idx, alias := range output.Source.Aliases {
		prg, err := expr.Compile(alias)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("alias.%d: compiling alias", idx))
		}

		aliasPrograms = append(aliasPrograms, prg)
	}

	var (
		attributeByID     = map[string]*Attribute{}
		attributePrograms = map[string]cel.Program{}
	)
	for idx, attr := range output.Attributes {
		prg, err := expr.Compile(attr.Source)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("attributes.%d (id = %s): compiling source", idx, attr.ID))
		}

		attributeByID[attr.ID] = attr
		attributePrograms[attr.ID] = prg
	}

	catalogEntryModels := []*CatalogEntryModel{}
	for _, entry := range entries {
		name, err := expr.Eval[string](ctx, nameProgram, entry)
		if err != nil {
			return nil, errors.Wrap(err, "evaluating entry name")
		}

		externalID, err := expr.Eval[string](ctx, externalIDProgram, entry)
		if err != nil {
			return nil, errors.Wrap(err, "evaluating entry external ID")
		}

		aliases := []string{}
		for idx, aliasProgram := range aliasPrograms {
			alias, err := expr.Eval[string](ctx, aliasProgram, entry)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("aliases.%d: evaluating entry external ID", idx))
			}

			aliases = append(aliases, alias)
		}

		// Attribute values are built best effort, as it might not be the case that upstream
		// source entries have these fields, or have fields of the correct type.
		attributeValues := map[string]client.CatalogAttributeBindingPayloadV2{}
	eachAttribute:
		for attributeID, prg := range attributePrograms {
			binding := client.CatalogAttributeBindingPayloadV2{}

			if attributeByID[attributeID].Array {
				valueLiterals, err := expr.Eval[[]any](ctx, prg, entry)
				if err != nil {
					continue eachAttribute
				}

				arrayValue := []client.CatalogAttributeValuePayloadV2{}
				for _, literalAny := range valueLiterals {
					literal, ok := literalAny.(string)
					if !ok {
						continue
					}

					arrayValue = append(arrayValue, client.CatalogAttributeValuePayloadV2{
						Literal: lo.ToPtr(literal),
					})
				}

				binding.ArrayValue = &arrayValue
			} else {
				literal, err := expr.Eval[string](ctx, prg, entry)
				if err != nil {
					continue eachAttribute
				}

				binding.Value = &client.CatalogAttributeValuePayloadV2{
					Literal: lo.ToPtr(literal),
				}
			}

			attributeValues[attributeID] = binding
		}

		catalogEntryModels = append(catalogEntryModels, &CatalogEntryModel{
			Name:            name,
			ExternalID:      externalID,
			Aliases:         aliases,
			AttributeValues: attributeValues,
		})
	}

	return catalogEntryModels, nil
}
