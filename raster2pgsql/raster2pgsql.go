package raster2pgsql

import (
    "os"
    "os/exec"
    "fmt"
    "path/filepath"
    log "github.com/Sirupsen/logrus"
)
var (
    outputPath = filepath.Join("./output")
    bashScript = filepath.Join( "_script.sh" )
)

func ExeRaster(fileName string, tableColumn string, tableName string, constraint bool) ([]byte, error) {
    log.Debug("Preparing folders to execute raster2pgsql")
    os.RemoveAll(outputPath)
    log.Debug("Removed output folder")
    err := os.MkdirAll( outputPath, os.ModePerm|os.ModeDir )
    log.Info("Created folder")
    if err != nil {
        return nil, err
    }
    file, err := os.Create( filepath.Join(outputPath, bashScript))
    if err != nil {
        return nil, err
    }
    defer file.Close()
    log.Debug("Creating bash script")
    file.WriteString("#!/bin/sh\n")
    command := fmt.Sprintf("raster2pgsql -t 128x128 -c %s -x -f %s %s ",fileName, tableColumn, tableName)
    if constraint == true {
        command += " -C "
    }
    command += "> out.sql\n"
    file.WriteString(command)
    log.Debug("command: ", command)
    //set working directory to outputPath
    // err = os.Chdir(outputPath)
    
    if err != nil {
        return nil, err
    }
    log.Debug("Executing bashScript")
    out, err := exec.Command("sh", filepath.Join(outputPath, bashScript)).Output()
    log.Debug("Output ", string(out))
    
    error := os.Remove(filepath.Join(outputPath, bashScript))
    error = os.Remove(filepath.Join(outputPath))
    if error != nil {
        log.Error("Error removing file. ", error)
    }
    return out, err
}
