package nationalgrid

import (
	"fmt"
	"strconv"
)

func ValidateGridRef(ref string) error {
	l := len(ref)

	if l < 2 || l > 6 {
		return fmt.Errorf("a valid grid ref must be 2, 4, or 6 chars in length %v", ref)
	}

	validateSquare := func(ref string) error {
		_, err := strconv.Atoi(ref)
		if err == nil {
			return fmt.Errorf("the first two characters of a gridref cannot be numeric %v", ref)
		}
		return nil
	}

	validateSubSquare := func(ref string) error {
		_, err := strconv.Atoi(ref)
		if err != nil {
			return fmt.Errorf("the third and fourth characters of a gridref must be numeric %v", ref)
		}
		return nil
	}

	validateQuadrant := func(ref string) error {
		if ref != string(NE) && ref != string(NW) && ref != string(SE) && ref != string(SW) {
			return fmt.Errorf("the fifth and sixth characters of a gridref must be a valid quadrant %v", ref)
		}
		return nil
	}

	switch l {
	case 2:
		err := validateSquare(ref)
		if err != nil {
			return err
		}

	case 4:
		square := ref[0:2]
		err := validateSquare(square)
		if err != nil {
			return err
		}
		subsquare := ref[2:4]
		err = validateSubSquare(subsquare)
		if err != nil {
			return err
		}

	case 6:
		square := ref[0:2]
		err := validateSquare(square)
		if err != nil {
			return err
		}
		subsquare := ref[2:4]
		err = validateSubSquare(subsquare)
		if err != nil {
			return err
		}
		quadrant := ref[4:6]
		err = validateQuadrant(quadrant)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateSquare(square string) error {
	if _, ok := NationalGridSquares[square]; !ok {
		return fmt.Errorf("unable to load sector %v", square)
	}

	return nil
}
