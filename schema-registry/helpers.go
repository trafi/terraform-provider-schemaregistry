package schemaregistry

import (
	"fmt"
	"strconv"
	"strings"
)

const separator = ":"

func parseID(serializedID string) (string, int, error) {
	parts := strings.Split(serializedID, separator)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("Unable to parse ID %v", serializedID)
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("Cannot parse schema id")
	}

	return parts[0], id, nil
}

func serializeID(subject string, id int) string {
	s := []string{subject, strconv.Itoa(id)}

	return strings.Join(s, separator)
}
