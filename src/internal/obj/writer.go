package obj

import (
	"os"
	"strconv"
	"strings"
	"tucil/src/internal/model"
)

func WriteOBJToFile(filename string, object *model.Object) error {
	var sb strings.Builder

	// Write vertexes
	for _, v := range object.Vertexes {
		sb.WriteString("v ")
		sb.WriteString(strconv.FormatFloat(v.X, 'f', 6, 64))
		sb.WriteString(" ")
		sb.WriteString(strconv.FormatFloat(v.Y, 'f', 6, 64))
		sb.WriteString(" ")
		sb.WriteString(strconv.FormatFloat(v.Z, 'f', 6, 64))
		sb.WriteString("\n")
	}

	// Write faces
	for _, f := range object.Faces {
		sb.WriteString("f ")
		for _, idx := range f.Vertexes {
			sb.WriteString(strconv.Itoa(idx + 1)) // OBJ is 1-indexed
			sb.WriteString(" ")
		}
		sb.WriteString("\n")
	}

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}
