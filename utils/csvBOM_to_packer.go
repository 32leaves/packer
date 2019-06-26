// This utility takes the output of CSV BOM for Autodesk Fusion 360 (https://apps.autodesk.com/FUSION/en/Detail/Index?id=2857656002406992603&appLang=en&os=Mac)
// and transforms it into a JSON file that the packer CLI can work with.
//
// We assume that the smallest dimension is the material thickness.

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/32leaves/packer/pkg/twod"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <bom.csv>\n", os.Args[0])
		os.Exit(1)
	}

	fn := os.Args[1]
	in, err := os.OpenFile(fn, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	groups := make(map[string][]twod.Item)
	r := csv.NewReader(in)
	firstLine := true
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if firstLine {
			firstLine = false
			continue
		}

		qty, err := strconv.Atoi(record[1])
		if err != nil {
			log.Fatalf("cannot parse quantity of %v: %v", record, err)
		}
		item, thickness, err := rowToItem(record)
		if err != nil {
			log.Fatal(err)
		}

		grp, ok := groups[thickness]
		if !ok {
			grp = make([]twod.Item, 0)
		}
		for q := 0; q < qty; q++ {
			itm := *item
			itm.Name = fmt.Sprintf("%s.%03d", itm.Name, q)
			grp = append(grp, itm)
		}
		groups[thickness] = grp
	}

	basename := strings.TrimSuffix(fn, filepath.Ext(fn))
	for k, grp := range groups {
		prob := twod.Problem{
			Sheet: twod.Sheet{},
			Items: grp,
		}
		fc, err := json.MarshalIndent(prob, "", "    ")
		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s_sheet%s.json", basename, k), fc, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func rowToItem(record []string) (item *twod.Item, thickness string, err error) {
	dims := make([]int64, 3)
	smallestIdx := 0
	smallestDim := int64(math.MaxInt64)
	for i, v := range record[2:5] {
		v, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, "", err
		}
		dims[i] = int64(v)

		if dims[i] < smallestDim {
			smallestIdx = i
			smallestDim = dims[i]
		}
	}

	thickness = strconv.FormatInt(dims[smallestIdx], 10)
	dims = append(dims[:smallestIdx], dims[smallestIdx+1:]...)

	item = &twod.Item{
		Name:   record[0],
		Width:  dims[0],
		Height: dims[1],
	}
	return
}
