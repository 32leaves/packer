package twod

import (
	"fmt"
)

// Problem describes the situation we wish to optimize
type Problem struct {
	Sheet Sheet  `json:"sheet"`
	Items []Item `json:"items"`
}

// Validate checks if the problem can be solved.
// If this returns an error the problem definitely cannot be solved.
// If this does not return an error it's still possible that the problem has no solution.
func (p *Problem) Validate() error {
	for _, i := range p.Items {
		if i.Width > p.Sheet.Width {
			return &NoSolutionError{i, "item is wider than the sheet"}
		}
		if i.Height > p.Sheet.Height {
			return &NoSolutionError{i, "item is taller than the sheet"}
		}
	}
	return nil
}

// Sheet describes the container we want to pack the items on
type Sheet struct {
	Width int64 `json:"width"`

	// Height is the height of the sheet in units. If zero we'll pack a single infinite sheet
	Height int64 `json:"height,omitempty"`
}

// SVG produces an SVG representation of this sheet
func (s *Sheet) SVG() string {
	return fmt.Sprintf("<rect width=\"%d\" height=\"%d\" class=\"sheet\" />", s.Width, s.Height)
}

// Item describes a single item we want to pack on a sheet
type Item struct {
	Name   string `json:"name"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
}

// SVG produces an SVG representation of this item
func (i *Item) SVG() string {
	return fmt.Sprintf("<rect width=\"%d\" height=\"%d\" class=\"item\"><text>%s</text></rect>", i.Width, i.Height, i.Name)
}

// Solution has the items packed on sheets
type Solution struct {
	Sheets []PackedSheet `json:"sheet"`
	Cost   int64         `json:"cost"`
}

// SVG produces an SVG representation of this item
func (s *Solution) SVG() string {
	style := `<style>
svg {
    background-color: white;
}
.packed-sheet {
	stroke: black;
	fill: none;
}
.packed-item {
	stroke: red;
	fill: none;
}
</style>`

	sheets := ""
	for i, s := range s.Sheets {
		sheets += fmt.Sprintf("<g transform=\"translate(0,%f)\">%s</g>", float32(i)*float32(s.Height)*1.05, s.SVG())
	}

	var w, h float32
	if len(s.Sheets) > 0 {
		w = float32(len(s.Sheets)) * float32(s.Sheets[0].Height) * 1.05
		h = float32(s.Sheets[0].Width)
	}

	return fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"%f\" height=\"%f\">%s%s</svg>", w, h, style, sheets)
}

// PackedSheet contains a set of items aranged on a sheet
type PackedSheet struct {
	Sheet
	Items []PackedItem `json:"items"`
}

// SVG produces an SVG representation of this sheet
func (s *PackedSheet) SVG() string {
	sheet := s.Sheet.SVG()
	items := ""
	for _, i := range s.Items {
		items += i.SVG()
	}
	return fmt.Sprintf("<g class=\"packed-sheet\">%s%s</g>", sheet, items)
}

// PackedItem is an item placed somewhere in space
type PackedItem struct {
	Item
	OffsetX int64 `json:"x"`
	OffsetY int64 `json:"y"`
}

// SVG produces an SVG representation of this item
func (p *PackedItem) SVG() string {
	return fmt.Sprintf("<g transform=\"translate(%d,%d)\" class=\"packed-item\">%s</g>", p.OffsetX, p.OffsetY, p.Item.SVG())
}

// NoSolutionError indicates that a problem has no solution
type NoSolutionError struct {
	Item   Item
	Reason string
}

func (e *NoSolutionError) Error() string {
	return e.Reason
}
