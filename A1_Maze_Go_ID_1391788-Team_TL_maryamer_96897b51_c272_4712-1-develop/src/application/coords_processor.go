package application

import (
	"errors"
	"maze/domain"
)

type CoordsParser interface {
	Parse(string, string) (domain.Coords, domain.Coords, error)
}

type CoordsProcessor struct {
	parser CoordsParser
}

func NewCoordsProcessor(parser CoordsParser) CoordsProcessor {
	return CoordsProcessor{parser: parser}
}

func (p CoordsProcessor) GetCoords(startString, endString string, rows, cols int) (domain.Coords, domain.Coords, error) {
	start, end, err := p.parser.Parse(startString, endString)
	if err != nil {
		return domain.Coords{}, domain.Coords{}, err
	}

	if err = coordsIsValid(start, cols, rows); err != nil {
		return domain.Coords{}, domain.Coords{}, err
	}
	if err = coordsIsValid(end, cols, rows); err != nil {
		return domain.Coords{}, domain.Coords{}, err
	}

	return start, end, nil
}

// helper
func coordsIsValid(coords domain.Coords, cols, rows int) error {
	if coords.X < 0 || coords.Y < 0 || coords.X > cols-1 || coords.Y > rows-1 {
		return errors.New("Некорректные координаты")
	}
	return nil
}
