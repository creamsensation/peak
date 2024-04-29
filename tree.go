package peak

import (
	"slices"
)

func createTree(builders ...QueryBuilder) queryBuildersTree {
	filters := make([]*filterBuilder, 0)
	relationships := make([]*relationshipBuilder, 0)
	selectors := make([]*selectorBuilder, 0)
	shapes := make([]*shapeBuilder, 0)
	temporaries := make([]*temporaryBuilder, 0)
	values := make([]*valuesBuilder, 0)
	for _, builder := range builders {
		switch b := builder.(type) {
		case *filterBuilder:
			filters = append(filters, b)
		case *relationshipBuilder:
			relationships = append(relationships, b)
		case *selectorBuilder:
			selectors = append(selectors, b)
		case *shapeBuilder:
			shapes = append(shapes, b)
		case *temporaryBuilder:
			temporaries = append(temporaries, b)
		case *valuesBuilder:
			values = append(values, b)
		}
	}
	if doesExistAggregatePart(selectors) {
		shapes = updateShapesWithNonAggregateFields(shapes, selectors)
	}
	return queryBuildersTree{
		filters:       filters,
		relationships: relationships,
		selectors:     selectors,
		shapes:        shapes,
		temporaries:   temporaries,
		values:        values,
	}
}

func updateShapesWithNonAggregateFields(shapes []*shapeBuilder, selectors []*selectorBuilder) []*shapeBuilder {
	shapesExist := len(shapes) > 0
	nonAggregateFields := make([]QueryBuilder, 0)
	for _, selector := range selectors {
		for _, part := range selector.parts {
			if part.partType == selectorFieldPart && !slices.Contains(nonAggregateFields, part.builder) {
				nonAggregateFields = append(nonAggregateFields, part.builder)
			}
		}
	}
	if !shapesExist {
		shapes = append(
			shapes, &shapeBuilder{
				groupFields: nonAggregateFields,
			},
		)
	}
	if shapesExist {
		for i, shape := range shapes {
			if len(shape.groupFields) == 0 {
				continue
			}
			shapes[i].groupFields = append(
				shapes[i].groupFields,
				nonAggregateFields...,
			)
		}
	}
	return shapes
}

func doesExistAggregatePart(selectors []*selectorBuilder) bool {
	for _, selector := range selectors {
		for _, part := range selector.parts {
			if part.partType == selectorAggregatePart {
				return true
			}
		}
	}
	return false
}
