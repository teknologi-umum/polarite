package paste

import (
	"strings"
)

func (c *Dependency) UpdateIDListFromDB(pastes []Item) ([]string, error) {
	var temp []string
	for i := 0; i < len(pastes); i++ {
		temp = append(temp, pastes[i].ID)
	}

	s := strings.Join(temp, ",")
	err := c.Memory.Set("ids", []byte(s))
	if err != nil {
		return []string{}, err
	}

	return temp, nil
}

func (c *Dependency) UpdateIDListFromCache(pastes []string, new string) (int, error) {
	pastes = append(pastes, new)
	s := strings.Join(pastes, ",")
	err := c.Memory.Set("ids", []byte(s))
	if err != nil {
		return 0, err
	}

	return len(pastes), nil
}
