package coords

import (
	"errors"
	"maze/domain"
	"strconv"
	"strings"
)

type Parser struct{}

func NewParser() Parser {
	return Parser{}
}

func (p Parser) Parse(exit1Coords, exit2Coords string) (exit1 domain.Coords, exit2 domain.Coords, err error) {
	exit1, err = parseCoord(exit1Coords)
	if err != nil {
		return domain.Coords{}, domain.Coords{}, err
	}

	exit2, err = parseCoord(exit2Coords)
	if err != nil {
		return domain.Coords{}, domain.Coords{}, err
	}

	return exit1, exit2, nil
}

func parseCoord(coordstring string) (domain.Coords, error) {
	var err error
	coordsArray := strings.Split(coordstring, ":")

	if len(coordsArray) != 2 {
		return domain.Coords{}, errors.New("Неверный формат координат")
	}

	exit := domain.Coords{}

	exit.X, err = strconv.Atoi(coordsArray[0])
	if err != nil {
		return domain.Coords{}, errors.New("Неверный формат координат")
	}

	exit.Y, err = strconv.Atoi(coordsArray[1])
	if err != nil {
		return domain.Coords{}, errors.New("Неверный формат координат")
	}

	return exit, nil
}
