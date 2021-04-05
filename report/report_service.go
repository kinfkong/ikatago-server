package report

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/kinfkong/ikatago-server/config"
	"github.com/kinfkong/ikatago-server/utils"
)

// Service type
type Service struct {
	queue        *utils.MB
	PlatformName string
}

var serviceInstance *Service
var serviceMu sync.Mutex

// GetService returns the singleton instance of the Service
func GetService() *Service {
	serviceMu.Lock()
	defer serviceMu.Unlock()

	if serviceInstance == nil {
		serviceInstance = &Service{
			queue: utils.NewMB(4096),
		}
	}
	return serviceInstance
}

func (service *Service) StartReport() {
	failedCount := 0
	for {
		batchData := service.queue.WaitTimeoutOrMax(time.Duration(time.Second*5), 100)
		// log.Printf("fowarding: %d messages", len(batchData))
		// do real sent
		if len(batchData) > 0 {
			data, _ := json.Marshal(batchData)
			_, err := service.doSendWithRetry(data)
			if err != nil {
				log.Printf("ERROR failed to send data")
				failedCount = failedCount + 1
			} else {
				failedCount = 0
			}
			if failedCount >= 1 {
				time.Sleep(time.Duration(failedCount) * time.Second)
			}
		}
	}
}

func (server *Service) doSendWithRetry(data []byte) (string, error) {
	reportUrl := config.GetConfig().GetString("report.url")
	token := "772758d3-fda6-4987-9944-b420ba9ebad0"
	i := 0
	response := ""
	for i = 0; i < 3 && reportUrl != ""; i++ {
		var err error = nil
		response, err = utils.DoHTTPRequest("POST", reportUrl+"?token="+token, map[string]string{
			"Content-Type": "application/json",
		}, data)
		if err != nil {
			log.Printf("ERROR failed to send request")
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}
	if i == 3 {
		return response, errors.New("failed to send requests.")
	}
	return response, nil
}
func (service *Service) AddToQueue(data interface{}) error {
	err := service.queue.Add(data)
	if err != nil {
		log.Printf("ERROR failed to add data to queue: data:%v error: %v", data, err)
		return errors.New("failed_add_data_to_queue")
	}
	return nil
}
