package obj

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ReadFile reads the content of a file given its path and returns it as a string.
func ReadFile(filePath string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	f, err := os.Open(filepath.Join(dir, filePath))

	if err != nil {
		return "", err
	}

	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return "", err
	}

	size := fileInfo.Size()
	buffer := make([]byte, size)

	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}

// ParseOBJ takes the content of an OBJ file as a string and parses it into an Object struct.
// Only read vertex and face data, ignore other lines.
func ParseOBJ(content string) (*Object, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	object := &Object{
		Vertices: []Vertice{},
		Faces:    []Face{},
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "v":
			if len(parts) >= 4 {
				x, _ := strconv.ParseFloat(parts[1], 64)
				y, _ := strconv.ParseFloat(parts[2], 64)
				z, _ := strconv.ParseFloat(parts[3], 64)
				object.Vertices = append(object.Vertices, Vertice{X: x, Y: y, Z: z})
			}
		case "f":
			face := Face{Vertices: []Vertice{}}
			for _, part := range parts[1:] {
				vertexIndex, _ := strconv.Atoi(part)
				if vertexIndex > 0 && vertexIndex <= len(object.Vertices) {
					face.Vertices = append(face.Vertices, object.Vertices[vertexIndex-1])
				}
			}
			object.Faces = append(object.Faces, face)
		}
	}

	return object, scanner.Err()
}
