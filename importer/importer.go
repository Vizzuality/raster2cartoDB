package importer

import (
	"sync"
    "gopkg.in/cheggaaa/pb.v1"
	log "github.com/Sirupsen/logrus"
    "github.com/agonzalezro/cartodb_go"
    "time"
)

func Run(lines []string, cartoDBApi string, user string, proyection int, numThreads int) {
	log.Info("Import lines with API: ", cartoDBApi, " Num lines:", len(lines))

	var waitGroup sync.WaitGroup
    bar := pb.StartNew(len(lines) + 3)
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
                    log.Error("Error in executor(", thread ,") with sql: ", line, ". Error: ", err)
                }      
			}
            time.Sleep(1 * time.Second)
			waitGroup.Done()
		}(i)
	}
	for _, line := range lines {
		tasks <- line
	}
    close(tasks)
    log.Debug("Waiting to executors")
	waitGroup.Wait()
    
    
    

}
