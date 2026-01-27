package gpsfiles

import (
	"math"

	lgo "github.com/yageek/lambertgo"

	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ConvertWgs84ToLambert93(longitude, latitude float64) (float64, float64, error) {
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

func ConvertLambert93ToWgs84(x, y float64) (float64, float64, error) {
	// Crée le point avec les coordonnées Lambert93
	point := &lgo.Point{
		X:    x,
		Y:    y,
		Z:    0,
		Unit: lgo.Meter, // Préciser que c'est en mètres
	}

	point.ToWGS84(lgo.Lambert93) // Conversion Lambert-93 → WGS84

	// Convertion de résultat en Radians vers des Degrés pour Leaflet de JS
	point.ToDegree()

	// point.Y est la Latitude, point.X est la Longitude
	return point.Y, point.X, nil
}

func GetFilesNameFolder(folderName string) ([]string, error) {
	entries, err := os.ReadDir(folderName)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, folderName+entry.Name())
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

func GetFileForMyCoordinate(x, y float64, folderPath string, accuracy string) (string, error, float64, float64) {
	path := folderPath

	var cellSize float64
	switch accuracy {
	case "1M":
		cellSize = 1
		path += "RGEALTI_FXX_"
	case "5M":
		path += "RGEALTI_FXX_"
		cellSize = 5
	case "25M":
		path += "BDALTIV2_25M_FXX_"
		cellSize = 25
	}

	fileX := math.Floor((x+cellSize/2)/(cellSize*1000)) * cellSize
	fileY := math.Floor((y-cellSize/2)/(cellSize*1000))*cellSize + cellSize

	path = fmt.Sprintf("%s%04.0f_%04.0f_MNT_LAMB93_IGN69.asc", path, fileX, fileY)

	xll := fileX*1000 - cellSize/2
	yll := (fileY-cellSize)*1000 + cellSize/2

	// 4. Calcul des bornes
	xmin := xll
	xmax := xll + 1000*cellSize
	ymin := yll
	ymax := yll + 1000*cellSize

	// 5. Test d’appartenance
	if x >= xmin && x <= xmax && y >= ymin && y <= ymax {
		return path, nil, xll, yll
	}

	return "", fmt.Errorf("coordinate not found in any file"), -1, -1
}

func BuildBDIndex(baseDir string) ([]BDFileInfo, error) {

	var results []BDFileInfo

	err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(d.Name()) != ".asc" {
			return nil
		}

		xll, yll, cell, err := ReadCoordinateLambert93File(path)
		if err != nil {
			return fmt.Errorf("erreur lecture %s : %w", path, err)
		}

		info := BDFileInfo{
			Path:      filepath.Base(path),
			XllCorner: xll,
			YllCorner: yll,
			CellSize:  cell,
		}

		results = append(results, info)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}
