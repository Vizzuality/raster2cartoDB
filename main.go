package main

import (
	"os"
    log "github.com/Sirupsen/logrus"
    "github.com/rrequero/importraster/importer"
    "github.com/rrequero/importraster/raster2pgsql"
    "github.com/rrequero/importraster/reader"
    "github.com/jessevdk/go-flags"
)

// init is called prior to main.
func init() {
	log.SetLevel(log.InfoLevel)
}

var opts struct {
    File string `short:"f" long:"file" description:"Path to tiff file" required:"true"` 
    TableName string `short:"t" long:"table" description:"Schema an Table where save data" required:"true"` 
    TableColumn string `short:"c" long:"col" description:"Name of column." default:"the_raster_webmercator" ` 
    Constraints bool `short:"C" description:"Add constraints"`
    Proyection int `short:"p" long:"proj" description:"Source proyection" required:"true"`
    Threads int `long:"threads" description:"Number of threads that execute INSERT commands to CartoDB" default:"2"`
    CartoAPI string `long:"api" description:"CartoDB API" required:"true"`
    CartoUser string `long:"user" description:"CartoDB User" required:"true"`
    Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information"`
}

// main is the entry point for the program.
func main() {
    
   
    _, err := flags.ParseArgs(&opts, os.Args[1:])
    
    if err != nil {
        os.Exit(1)
    }
    
    if opts.Verbose {
        log.SetLevel(log.DebugLevel)
    }
    
    log.Info("Start import tiff")
    log.Info("First: Create sql file from tiff image")
    _, err = raster2pgsql.ExeRaster(opts.File, opts.TableColumn, opts.TableName, opts.Constraints)
    if err != nil {
        log.Error(err)
        panic(err)        
    }
    log.Info("Second: Reading file")
    lines, err := reader.ReadFile()
    if err != nil {
        panic(err)
    } 
    log.Info("Third: Importing in CartoDB")
    importer.Run(lines, opts.CartoAPI, opts.CartoUser, opts.Proyection, opts.Threads)
    
    if err != nil {
        panic(err)
    }
    log.Info("Successful import!!")
}
