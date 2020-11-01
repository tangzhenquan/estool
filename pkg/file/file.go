package file

import "os"

func ReadOpen(path string) (*os.File, error) {
	flag := os.O_RDONLY
	perm := os.FileMode(0)
	return os.OpenFile(path, flag, perm)
}
/*

func (h *Harvester) openFile() error {
	fi, err := os.Stat(h.state.Source)
	if err != nil {
		return fmt.Errorf("failed to stat source file %s: %v", h.state.Source, err)
	}
	if fi.Mode()&os.ModeNamedPipe != 0 {
		return fmt.Errorf("failed to open file %s, named pipes are not supported", h.state.Source)
	}

	f, err := file_helper.ReadOpen(h.state.Source)
	if err != nil {
		return fmt.Errorf("Failed opening %s: %s", h.state.Source, err)
	}

	harvesterOpenFiles.Add(1)

	// Makes sure file handler is also closed on errors
	err = h.validateFile(f)
	if err != nil {
		f.Close()
		harvesterOpenFiles.Add(-1)
		return err
	}

	h.source = File{File: f}
	return nil
}



 */