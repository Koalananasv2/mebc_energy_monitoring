package main

import (
	"log"
	"context"
	"fmt"
	"os"
	"io/ioutil"
	"strings"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"regexp"
	"time"
	"monitoring/api/team_data_struct"
)

var influxDBClient influxdb2.Client
var influxDBWriteAPI api.WriteAPIBlocking
var ValidationRegex = regexp.MustCompile("^[a-zA-Z0-9]+$")
var TeamInfoMap map[string]TeamInfo
var token = os.Getenv("INFLUXDB_TOKEN")

type TeamInfo struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
}

func initInfluxDB() {
	influxport := os.Getenv("INFLUXDB_PORT")
	influxorg := os.Getenv("INFLUXDB_ORG")
	influxbck := os.Getenv("INFLUXDB_BCK")
	influxDBClient = influxdb2.NewClientWithOptions(fmt.Sprintf("http://localhost:%s", influxport), token, influxdb2.DefaultOptions().SetBatchSize(50))
	influxDBWriteAPI = influxDBClient.WriteAPIBlocking(influxorg, influxbck)
}

func loadTeamInfoFromFile() error {
	data, err := ioutil.ReadFile("Team_info.txt")
	if err != nil {
		return err
	}

	TeamInfoMap = make(map[string]TeamInfo)

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) == 2 {
			token := fields[0]
			name := fields[1]

			TeamInfoMap[token] = TeamInfo{
				Name:   name,
			}
		}
	}

	return nil
}

func main() {
	err := loadTeamInfoFromFile()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de TeamInfoMap depuis le fichier : %v", err)
	}

	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	initInfluxDB()

	// Route for creating monitoring data
	app.Post("/monitoringdata", createMonitoringData)

	// Start the server
	err = app.Listen(":3001")
	if err != nil {
		panic("Failed to start the server!")
	}
}

// Create monitoring data
func createMonitoringData(c *fiber.Ctx) error {

	//extract team info
	var request struct {
		Team string `json:"team"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request Team",
		})
	}
	// Vérifier que le champ Team ne contient que des lettres ou des chiffres
	if !ValidationRegex.MatchString(request.Team) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Team format. Only letters and numbers are allowed.",
		})
	}

	var sensorData team_data_struct.SensorDataInterface

	if teamInfo, ok := TeamInfoMap[request.Team]; ok {
		switch teamInfo.Name {
		case "Generic 1":
			sensorData = &team_data_struct.GenericSensorData{}
		case "Generic 2":
			sensorData = &team_data_struct.GenericSensorDataDualPower{}
		case "Generic 3":
			sensorData = &team_data_struct.GenericSensorDataQuadPower{}
		default:
			sensorData = &team_data_struct.GenericSensorData{}
		}

		if err := c.BodyParser(&sensorData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Bad Request Data",
			})
		}


		if sensorDataQuad, ok := sensorData.(team_data_struct.SensorDataInterfaceQuadPower); ok {
			var result float32
			var dataReceived bool
			if sensorDataQuad.GetVoltage1() != nil && sensorDataQuad.GetCurrent1() != nil {
				result += *sensorDataQuad.GetVoltage1() * *sensorDataQuad.GetCurrent1()
				dataReceived = true
			}
			if sensorDataQuad.GetVoltage2() != nil && sensorDataQuad.GetCurrent2() != nil {
				result += *sensorDataQuad.GetVoltage2() * *sensorDataQuad.GetCurrent2()
				dataReceived = true
			}
			if sensorDataQuad.GetVoltage3() != nil && sensorDataQuad.GetCurrent3() != nil {
				result += *sensorDataQuad.GetVoltage3() * *sensorDataQuad.GetCurrent3()
				dataReceived = true
			}
			if sensorDataQuad.GetVoltage4() != nil && sensorDataQuad.GetCurrent4() != nil {
				result += *sensorDataQuad.GetVoltage4() * *sensorDataQuad.GetCurrent4()
				dataReceived = true
			}
			if dataReceived {
				sensorDataQuad.SetPower(result)
			}
		} else if sensorDatadual, ok := sensorData.(team_data_struct.SensorDataInterfaceDualPower); ok {
			var result float32
			var datareceived bool
			if sensorDatadual.GetVoltage1() != nil && sensorDatadual.GetCurrent1() != nil {
				result += *sensorDatadual.GetVoltage1() * *sensorDatadual.GetCurrent1()
				datareceived = true
			}
			if sensorDatadual.GetVoltage2() != nil && sensorDatadual.GetCurrent2() != nil {
				result += *sensorDatadual.GetVoltage2() * *sensorDatadual.GetCurrent2()
				datareceived = true
			}
			if datareceived {
				sensorDatadual.SetPower(result)
			}
		} else {
			if voltage := sensorData.GetVoltage(); voltage != nil && sensorData.GetCurrent() != nil {
				result := *voltage * *sensorData.GetCurrent()
				sensorData.SetPower(result)
			}
		}

		point := write.NewPointWithMeasurement("monitoring_data")
		sensorDataValue := reflect.ValueOf(sensorData).Elem() // Assurez-vous que sensorData est déjà un pointeur
		sensorDataType := sensorDataValue.Type()

		for i := 0; i < sensorDataType.NumField(); i++ {
			field := sensorDataType.Field(i)
			fieldValue := sensorDataValue.Field(i)

			// Vérifiez si le champ est non-nil et non-Team avant de l'ajouter
			if !fieldValue.IsNil() {
				// Utilisez le nom du champ JSON pour le nom du champ InfluxDB
				jsonTag := field.Tag.Get("json")
				point.AddField(jsonTag, fieldValue.Elem().Interface()) // Utilisez Elem() pour obtenir la valeur pour les types pointeur
			}
		}

		point = point.SetTime(time.Now())
		// Écrire le point dans InfluxDB
		err := influxDBWriteAPI.WritePoint(context.Background(), point)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}
	} else {
		// Si non, renvoyer une erreur
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Team. Team information not found.",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}