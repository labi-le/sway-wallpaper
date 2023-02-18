package browser

import (
	"database/sql"
	"os"
)

type History struct {
	DB      *sql.DB
	TmpFile *os.File
}

func (h *History) Cleanup() error {
	if err := h.DB.Close(); err != nil {
		return err
	}

	if err := h.TmpFile.Close(); err != nil {
		return err
	}

	if err := os.Remove(h.TmpFile.Name()); err != nil {
		return err
	}

	return nil
}
