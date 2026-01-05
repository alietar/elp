package findfiles

import (
	lgo "github.com/yageek/lambertgo"

	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func FromGpsWgs84ToLambert93(longStr string, latStr string) (float64, float64, error) {
	longitude, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		return 0, 0, err
	}

	latitude, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return 0, 0, err
	}

	// Point WGS84 en degrés
	point := &lgo.Point{
		X:    longitude,
		Y:    latitude,
		Z:    0,
		Unit: lgo.Degree,
	}

	// degrés → radians
	point.ToRadian()

	// WGS84 → Lambert-93 (EPSG:2154)
	point.ToLambert(lgo.Lambert93)

	return point.X, point.Y, nil
}

func GetFilesNameFolder(folderName string) ([]string, error) {
	entries, err := os.ReadDir(folderName)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

func ReadCoordinateLambert93File(fileName string) (float64, float64, float64, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, 0, 0, err
	}
	defer file.Close()

	var xll, yll, cellsize float64
	foundX, foundY, foundCell := false, false, false

	r := bufio.NewReader(file)

	for {
		line, err := r.ReadBytes('\n')
		if err != nil && len(line) == 0 {
			break
		}

		fields := strings.Fields(string(line))
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "xllcorner":
			xll, err = strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return 0, 0, 0, err
			}
			foundX = true

		case "yllcorner":
			yll, err = strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return 0, 0, 0, err
			}
			foundY = true

		case "cellsize":
			cellsize, err = strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return 0, 0, 0, err
			}
			foundCell = true
		}

		if foundX && foundY && foundCell {
			return xll, yll, cellsize, nil
		}
	}

	return 0, 0, 0, fmt.Errorf("header incomplete in %s", fileName)
}

func GetFileForMyCoordinate(longStr string, latStr string, folderPath string) (string, error) {

	// 1. Conversion WGS84 → Lambert-93
	x, y, err := FromGpsWgs84ToLambert93(longStr, latStr)
	if err != nil {
		return "", err
	}

	// 2. Liste des fichiers
	files, err := GetFilesNameFolder(folderPath)
	if err != nil {
		return "", err
	}

	// 3. Parcours des fichiers
	for _, file := range files {
		if !strings.HasSuffix(file, ".asc") {
			continue
		}

		fullPath := folderPath + "/" + file

		xll, yll, cellsize, err := ReadCoordinateLambert93File(fullPath)
		if err != nil {
			continue
		}

		// 4. Calcul des bornes
		xmin := xll
		xmax := xll + 1000*cellsize
		ymin := yll
		ymax := yll + 1000*cellsize

		// 5. Test d’appartenance
		if x >= xmin && x <= xmax && y >= ymin && y <= ymax {
			return file, nil
		}
	}

	return "", fmt.Errorf("coordinate not found in any file")
}
