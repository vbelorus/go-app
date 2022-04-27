package clickhouse

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/vbelorus/go-app/v2/src/models"
	"os"
	"strconv"
	"time"
)

type AppDB struct {
	Conn          driver.Conn
	ErrorsChannel chan error
}

const defaultBatchSize = 1000

func (db *AppDB) BatchSize() int {
	batchSizeEnv, ok := os.LookupEnv("DEVICE_EVENT_BATCH_SIZE")
	if !ok {
		return defaultBatchSize
	} else {
		batchSize, err := strconv.Atoi(batchSizeEnv)
		if err != nil {
			return defaultBatchSize
		}
		return batchSize
	}
}

func (db *AppDB) ListenDeviceEvents(ch chan models.DeviceEvent) {
	var bufferedDeviceEvents []models.DeviceEvent

	for {
		select {
		case v := <-ch:
			bufferedDeviceEvents = append(bufferedDeviceEvents, v)
			if len(bufferedDeviceEvents) > db.BatchSize() {
				err := db.SaveBatchDeviceEvents(bufferedDeviceEvents)
				if err != nil {
					db.ErrorsChannel <- err
				}
				bufferedDeviceEvents = bufferedDeviceEvents[:0]
			}
		//if no new data for 5 sec - send existing
		case <-time.After(5 * time.Second):
			if len(bufferedDeviceEvents) > 0 {
				err := db.SaveBatchDeviceEvents(bufferedDeviceEvents)
				if err != nil {
					db.ErrorsChannel <- err
				}
				bufferedDeviceEvents = bufferedDeviceEvents[:0]
			}
			fmt.Println("Time out: 5 second, save if something is in bufferedDeviceEvents")
		}
	}
}

func (db *AppDB) SaveBatchDeviceEvents(bufferedDeviceEvents []models.DeviceEvent) error {
	//todo move to main?
	ctx := context.Background()

	batch, err := db.Conn.PrepareBatch(ctx, "INSERT INTO events")
	if err != nil {
		return err
	}

	for i := 0; i < len(bufferedDeviceEvents); i++ {
		e := bufferedDeviceEvents[i]
		err := batch.Append(
			e.ClientTime.GetTime(), e.DeviceId, e.DeviceOs, e.Session, e.Sequence, e.Event, e.ParamInt, e.ParamString, e.ServerTime, e.Ip,
		)
		if err != nil {
			return err
		}
	}

	fmt.Println("Send batch to clickhouse")
	return batch.Send()
}
