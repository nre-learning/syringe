package main

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	influx "github.com/influxdata/influxdb/client/v2"
)

func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func initAntidoteStats() AntidoteStats {
	mockSyringeConfig := GetmockSyringeConfig()
	return AntidoteStats{
		InfluxURL: mockSyringeConfig.InfluxURL,
		InfluxUsername: mockSyringeConfig.InfluxUsername,
		InfluxPassword: mockSyringeConfig.InfluxPassword,
		Curriculum: GetCurriculum(GetmockSyringeConfig()),
		LiveLessonState: GetMockLiveLessonState()
	}
}

func createInfluxClient() influx.Client {
	mockSyringeConfig := GetmockSyringeConfig()
	client, _ := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:               mockSyringeConfig.InfluxURL,
		Username:           mockSyringeConfig.InfluxUsername,
		Password:           mockSyringeConfig.InfluxPassword,
		InsecureSkipVerify: true,
	})

	if err != nil {
		log.Error("Error creating InfluxDB Client: ", err.Error())
		return err
	}

	return client;
}

func getRowCount(client influx.Client, table string) int {
	query := influx.NewQuery(fmt.Sprintf("SELECT * FROM %s", table), "syringe_metrics", "")
	res, err := client.Query(query)

	if err != nil {
		fmt.Println(err.Error())
	}

	if len(res.Results) > 0 {
		if len(res.Results[0].Series) > 0 {
			series := res.Results[0].Series[0]
			return len(series.Values)
		}
	}

	return 0
}

func dropTable(client influx.Client, table string) {
	query := influx.NewQuery(fmt.Sprintf("DELETE FROM %s", table), "syringe_metrics", "")
	_, err := client.Query(query)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestStartTSDBExport(t testing.T) {
	stats := initAntidoteStats()
	err := stats.StartTSDBExport()
	ok(t, err)

	influxClient := createInfluxClient()
	defer influxClient.Close()

	dropTable(influxClient, "sessionStatus")
	time.Sleep(2 * time.Minute)
	rowCount := getRowCount(influxClient, "sessionStatus")

	assert(t, rowCount > 0, "")
}
