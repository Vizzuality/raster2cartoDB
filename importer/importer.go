package importer

import (
	"sync"
    "gopkg.in/cheggaaa/pb.v1"
	log "github.com/Sirupsen/logrus"
    "github.com/agonzalezro/cartodb_go"
    "time"
    "fmt"
)

func Run(lines []string, cartoDBApi string, user string, proyection int, numThreads int, tableName string) {
	log.Info("Import lines with API: ", cartoDBApi, " Num lines:", len(lines), "; Num Threads: ", numThreads)

	var waitGroup sync.WaitGroup
    bar := pb.StartNew(len(lines) + 2)
    //refresh 50ms
    bar.SetRefreshRate(time.Millisecond * 50)

	tasks := make(chan string)

	for i := 0; i < numThreads; i++ {
		waitGroup.Add(1)
		go func(thread int) {     
            client := cartodb.NewAPIKeyClient(cartoDBApi, user, "", "", "")       
			for line := range tasks {  
                _, err := ExecuteSQL(line, client)
                bar.Increment()
                if err != nil {
                    log.Error("Error in executor(", thread ,") with sql: ", ". Error: ", err)
                }      
			}
            time.Sleep(1 * time.Second)
			waitGroup.Done()
		}(i)
	}
    client := cartodb.NewAPIKeyClient(cartoDBApi, user, "", "", "")
    log.Info("Creating table")    
    ExecuteSQL(lines[1], client)
    bar.Increment()
    
    // remove BEGIN command and create table
	for _, line := range lines[2:] {
		tasks <- line
	}
    close(tasks)
    log.Debug("Waiting to executors")
	waitGroup.Wait()
    
    log.Debug("Updating tables");
    ExecuteSQL(fmt.Sprintf("update %s set the_raster_webmercator =st_transform(st_setSrid(the_raster_webmercator, %d), 3857)", tableName, proyection), client)
    bar.Increment()
    ExecuteSQL(fmt.Sprintf("GRANT SELECT ON %s TO tileuser;", tableName), client)
    bar.Increment()
    ExecuteSQL(fmt.Sprintf("GRANT SELECT ON %s TO publicuser;", tableName), client)
    bar.Increment()
}
