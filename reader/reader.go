package reader

import (
	"bufio"
	"io"
	"os"
    "path/filepath"
	log "github.com/Sirupsen/logrus"
)

func ReadFile() ([]string, error) {
	log.Info("Reading file: out.sql")
    path, err := filepath.Abs("out.sql")
	f, err := os.Open(path)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)

	var lines []string

	line, err := r.ReadString('\n')
	for {		
		if err == io.EOF {        
			break
		}
        if line != "" {
            lines = append(lines, line)
        }
		line, err = r.ReadString('\n')
	}
    if line != "" {
        lines = append(lines, line)
    }
	
	if err == io.EOF {
		log.Debug("Returning lines (len) ", len(lines))
		return lines, nil
	}

	return lines, err
}
