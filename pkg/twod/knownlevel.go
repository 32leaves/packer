package twod

// implements the known-level algorithms of:
// N. Ntene, J.H. van Vuuren,
// A survey and comparison of guillotine heuristics for the 2D oriented offline strip packing problem,
// Discrete Optimization,
// Volume 6, Issue 2,
// 2009,
// Pages 174-188,
// ISSN 1572-5286,
// https://doi.org/10.1016/j.disopt.2008.11.002.
// (http://www.sciencedirect.com/science/article/pii/S1572528608000844)
// Abstract: An overview and comparison is provided of a number of heuristics from the literature for the two-dimensional strip packing problem in which rectangles have to be packed without rotation. Heuristics producing only guillotine packings are considered. A new heuristic is also introduced and a number of modifications are suggested to the existing heuristics. The resulting heuristics (known and new) are then compared statistically with respect to a large set of known benchmarks at a 5% level of significance.
// Keywords: Strip packing problem; Level algorithms

import (
	"sort"
)

// NextFitDecreasingHeight implements the next-fit-decreasing-height algorithm of
// Performance Bounds for Level-Oriented Two-Dimensional Packing Algorithms
// E. G. Coffman, Jr., M. R. Garey, D. S. Johnson, and R. E. Tarjan
// SIAM Journal on Computing 1980 9:4, 808-826
func NextFitDecreasingHeight(problem Problem) (*Solution, error) {
	err := problem.Validate()
	if err != nil {
		return nil, err
	}

	sort.Slice(problem.Items, func(i, j int) bool {
		return problem.Items[j].Height < problem.Items[i].Height
	})

	perSheetCost := problem.Sheet.Width * problem.Sheet.Height
	solution := Solution{Sheets: []PackedSheet{}, Cost: perSheetCost}

	currentSheet := PackedSheet{Sheet: problem.Sheet, Items: make([]PackedItem, 0)}
	thisLevel := int64(0)
	level := int64(0)
	fill := int64(0)
	for _, item := range problem.Items {
		// the first item needs to initialize thisLevel, as it also started a level
		if thisLevel == 0 {
			thisLevel = item.Height
		}

		// need a new level
		if fill+item.Width > currentSheet.Width {
			level += thisLevel
			fill = 0

			// this item marks the start of a new level, thus its height is the offset for the next level
			thisLevel = item.Height

			// need a new sheet
			if 0 < currentSheet.Height && currentSheet.Height < level {
				// append old sheet
				solution.Sheets = append(solution.Sheets, currentSheet)

				// create new one
				currentSheet = PackedSheet{Sheet: problem.Sheet, Items: make([]PackedItem, 0)}
				level = 0

				solution.Cost += perSheetCost
			}
		}

		currentSheet.Items = append(currentSheet.Items, PackedItem{
			Item:    item,
			OffsetX: fill,
			OffsetY: level,
		})
		fill += item.Width
	}
	solution.Sheets = append(solution.Sheets, currentSheet)

	return &solution, nil
}

// FirstFitDecreasingHeight implements the first-fit-decreasing-height algorithm of
// Performance Bounds for Level-Oriented Two-Dimensional Packing Algorithms
// E. G. Coffman, Jr., M. R. Garey, D. S. Johnson, and R. E. Tarjan
// SIAM Journal on Computing 1980 9:4, 808-826
func FirstFitDecreasingHeight(problem Problem) (*Solution, error) {
	err := problem.Validate()
	if err != nil {
		return nil, err
	}

	sort.Slice(problem.Items, func(i, j int) bool {
		return problem.Items[j].Height < problem.Items[i].Height
	})

	type level struct {
		Height    int64
		SpaceLeft int64
		Items     []PackedItem
	}

	// map items to levels
	var levels []*level
	for _, item := range problem.Items {
		var l *level
		for _, lvl := range levels {
			if lvl.Height <= item.Height && lvl.SpaceLeft >= item.Width {
				l = lvl
				break
			}
		}

		if l == nil {
			l = &level{item.Height, problem.Sheet.Width, []PackedItem{}}
			levels = append(levels, l)
		}

		l.Items = append(l.Items, PackedItem{
			Item:    item,
			OffsetX: problem.Sheet.Width - l.SpaceLeft,
			OffsetY: 0,
		})
		l.SpaceLeft -= item.Width
	}

	// map levels to sheets
	sol := Solution{[]PackedSheet{}, 0}
	var sheet *PackedSheet
	var fill int64
	for _, lvl := range levels {
		if sheet == nil {
			sheet = &PackedSheet{
				Sheet: problem.Sheet,
				Items: []PackedItem{},
			}
		}
		if sheet != nil && fill+lvl.Height > sheet.Height {
			sol.Sheets = append(sol.Sheets, *sheet)
			fill = 0
		}

		for _, itm := range lvl.Items {
			sheet.Items = append(sheet.Items, PackedItem{
				Item:    itm.Item,
				OffsetX: itm.OffsetX,
				OffsetY: fill,
			})
		}
		fill += lvl.Height
	}
	sol.Sheets = append(sol.Sheets, *sheet)
	sol.Cost = int64(len(sol.Sheets)) * problem.Sheet.Width * problem.Sheet.Height

	return &sol, nil
}
