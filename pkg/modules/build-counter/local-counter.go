package build_counter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	path_utils "github.krafton.com/sbx/version-maker/pkg/modules/path-utils"
	"go.uber.org/zap"
)

type LocalCounter struct {
	db      map[string]int
	project string

	absPath string
}

func NewLocalCounter(path string, project string) (*LocalCounter, error) {
	absPath, err := path_utils.ResolvePathToAbs(path)
	if err != nil {
		return nil, err
	}

	db, err := readFileOrNew(absPath)
	if err != nil {
		return nil, fmt.Errorf("NewLocalCounterFailed, error: %s", err.Error())
	}

	return &LocalCounter{
		db:      db,
		project: project,
		absPath: absPath,
	}, nil
}

func (c *LocalCounter) String() string {
	return fmt.Sprintf("%#v", c)
}

func (c *LocalCounter) Increase(ctx context.Context) (uint, error) {
	if count, ok := c.db[c.project]; !ok {
		c.db[c.project] = 1

		err := writeFile(c.absPath, c.db)
		if err != nil {
			return 0, err
		}
		return 1, nil
	} else {
		count++
		c.db[c.project] = count

		err := writeFile(c.absPath, c.db)
		if err != nil {
			return 0, err
		}
		return uint(count), nil
	}
}

func (c *LocalCounter) Get(ctx context.Context) (uint, error) {
	if count, ok := c.db[c.project]; !ok {
		return 1, nil
	} else {
		return uint(count), nil
	}
}

func (c *LocalCounter) WriteToFile() error {
	return writeFile(c.absPath, c.db)
}

func readFileOrNew(path string) (map[string]int, error) {
	zap.S().Debugf("Read Local Count DB from %s", path)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		zap.S().Debug("File Not Exists, Initialize New DB")
		return map[string]int{}, nil
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	db := map[string]int{}
	err = json.Unmarshal(buf, &db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func writeFile(path string, db map[string]int) error {
	buf, err := json.Marshal(db)
	if err != nil {
		return err
	}

	zap.S().Debugf("Write Local Count DB to %s", path)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		zap.S().Debug("File Not Exists, Create Directory")
		err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(path, buf, 0644)
	if err != nil {
		return err
	}
	return nil
}
