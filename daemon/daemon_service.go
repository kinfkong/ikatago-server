package daemon

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
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
func (service *Service) IsDaemonAvailable() bool {
	port := config.GetConfig().Get("daemon.port")
	if port == nil {
		return false
	}
	v, err := utils.GetJSONIntNumber(port)
	return err == nil && v > 0
}

func (service *Service) StartDaemonReport() {
	log.Printf("Daemon report started")
	go func() {
		for {
			// run the gather infos
			workerData := WorkerData{
				WorkerType:      "ikatago-server",
				Timestamp:       time.Now(),
				RunningCommands: utils.GetCmdManager().GetAllCmdInfo(),
				ExtraInfo:       make(map[string]interface{}),
			}
			service.AddToQueue(workerData)
			time.Sleep(time.Second)
		}

	}()
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
	daemonPort := config.GetConfig().GetInt("daemon.port")
	daemonReportUrl := "http://localhost:" + strconv.Itoa(daemonPort) + "/api/worker/data"
	i := 0
	response := ""
	for i = 0; i < 3; i++ {
		var err error = nil
		response, err = utils.DoHTTPRequest("POST", daemonReportUrl, map[string]string{
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
		return response, errors.New("failed to send requests")
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
