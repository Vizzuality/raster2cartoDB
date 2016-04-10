package importer

import (
    "github.com/agonzalezro/cartodb_go"
    log "github.com/Sirupsen/logrus"
    "io/ioutil"
    "encoding/json"
    "errors"
)

type Error struct {
    Error []string `json:"error"`
}

func ExecuteSQL(line string, client *cartodb.APIKeyClient) ([]byte, error) {
    log.Debug("Executing CartoDBImporter")
    response, err := client.Req(client.ResourceURL, "GET", nil, "select * from lukkom")
    if err != nil {
		log.Error(err)
		return nil, err
	}
    body, err := ioutil.ReadAll(response.Body)
    res := Error{}
    json.Unmarshal(body, &res)
    
    
    if len(res.Error) > 0 {
        
        return nil, errors.New(res.Error[0])
    } 
    
    return body, err
}