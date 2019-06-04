package twod

import (
	"fmt"
)

// Problem describes the situation we wish to optimize
type Problem struct {
	Sheet Sheet  `json:"sheet"`
	Items []Item `json:"items"`
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
	sheets := ""
	for i, s := range s.Sheets {
		sheets += fmt.Sprintf("<g transform=\"translate(0,%f)\">%s</g>", float32(i)*float32(s.Height)*1.05, s.SVG())
	}
	return fmt.Sprintf("<svg>%s</svg>", sheets)
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
		items += "\n\t" + i.SVG()
	}
	return fmt.Sprintf("<g class=\"packed-sheet\">\n\t%s\n\t%s</g>", sheet, items)
}

// PackedItem is an item placed somewhere in space
type PackedItem struct {
	Item
	OffsetX int64 `json:"x"`
	OffsetY int64 `json:"y"`
}

// SVG produces an SVG representation of this item
func (p *PackedItem) SVG() string {
	return fmt.Sprintf("<g transform=\"translate(%d,%d)\">%s</g>", p.OffsetX, p.OffsetY, p.Item.SVG())
}

// NoSolutionError indicates that a problem has no solution
type NoSolutionError struct {
	Item   Item
	Reason string
}

func (e *NoSolutionError) Error() string {
	return e.Reason
}
