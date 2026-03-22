package obj

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"tucil/src/internal/model"
)

// ReadFile reads the content of a file given its path and returns it as a string.
func ReadFile(filePath string) (string, error) {
	f, err := os.Open(filePath)

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
func ParseOBJ(content string) (*model.Object, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	object := &model.Object{
		Faces:    []model.Face{},
		Vertexes: []model.Vertex{},
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "#") {
			line = line[:strings.Index(line, "#")]
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		prefix, args := parts[0], parts[1:]

		switch prefix {
		case "v":
			if len(args) >= 3 {
				x, err := strconv.ParseFloat(args[0], 64)
				if err != nil {
					return nil, err
				}
				y, err := strconv.ParseFloat(args[1], 64)
				if err != nil {
					return nil, err
				}
				z, err := strconv.ParseFloat(args[2], 64)
				if err != nil {
					return nil, err
				}
				object.Vertexes = append(object.Vertexes, model.Vertex{X: x, Y: y, Z: z})
			}
		case "f":
			face := model.Face{Vertexes: [3]int{}}
			for _, part := range args {
				vertexIndex, err := strconv.Atoi(part)
				if err != nil {
					return nil, err
				}
				if vertexIndex > 0 && vertexIndex <= len(object.Vertexes) {
					for i := range 3 {
						if face.Vertexes[i] == 0 {
							face.Vertexes[i] = vertexIndex - 1
							break
						}
					}
				}
			}
			object.Faces = append(object.Faces, face)
		}
	}

	if len(object.Faces) == 0 || len(object.Vertexes) == 0 {
		return nil, errors.New("no valid faces or vertexes found in OBJ content")
	}

	for i := range object.Faces {
		for j := range object.Faces[i].Vertexes {
			if object.Faces[i].Vertexes[j] < 0 || object.Faces[i].Vertexes[j] >= len(object.Vertexes) {
				return nil, errors.New("face vertex index out of bounds")
			}
		}
	}

	return object, scanner.Err()
}
