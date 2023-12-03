package main

import (
	"log"
	"context"
	"fmt"
	"os"
	"io/ioutil"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"regexp"
	"time"
)

// Model structure
type SensorData struct {
	Temp1   float32 `json:"temp1"`
	Temp2   float32 `json:"temp2"`
	Temp3   float32 `json:"temp3"`
	Temp4   float32 `json:"temp4"`
	Temp5   float32 `json:"temp5"`
	Voltage float32 `json:"voltage"`
	Current float32 `json:"current"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Team    string  `json:"team"`
}

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
	var sensorData SensorData
	if err := c.BodyParser(&sensorData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	// Vérifier que le champ Team ne contient que des lettres ou des chiffres
	if !ValidationRegex.MatchString(sensorData.Team) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Team format. Only letters and numbers are allowed.",
		})
	}

	// Vérifier si le token de l'équipe existe dans la map
	if teamInfo, ok := TeamInfoMap[sensorData.Team]; ok {
		// Si oui, mettre à jour les champs Team avec les informations correspondantes
		sensorData.Team = teamInfo.Name

		tp1 := fmt.Sprintf("%s%s", "temp1_", sensorData.Team)
		tp2 := fmt.Sprintf("%s%s", "temp2_", sensorData.Team)
		tp3 := fmt.Sprintf("%s%s", "temp3_", sensorData.Team)
		tp4 := fmt.Sprintf("%s%s", "temp4_", sensorData.Team)
		tp5 := fmt.Sprintf("%s%s", "temp5_", sensorData.Team)
		volt := fmt.Sprintf("%s%s", "voltage_", sensorData.Team)
		curr := fmt.Sprintf("%s%s", "current_", sensorData.Team)
		lat := fmt.Sprintf("%s%s", "lat_", sensorData.Team)
		lon := fmt.Sprintf("%s%s", "lon_", sensorData.Team)
		pow := fmt.Sprintf("%s%s", "power_", sensorData.Team)

		// Ajouter le numéro de l'équipe à la base de données
		point := write.NewPointWithMeasurement("monitoring_data").
				AddField(tp1, sensorData.Temp1).
				AddField(tp2, sensorData.Temp2).
				AddField(tp3, sensorData.Temp3).
				AddField(tp4, sensorData.Temp4).
				AddField(tp5, sensorData.Temp5).
				AddField(volt, sensorData.Voltage).
				AddField(curr, sensorData.Current).
				AddField(lat, sensorData.Lat).
				AddField(lon, sensorData.Lon).
				AddField(pow, sensorData.Voltage * sensorData.Current).
				SetTime(time.Now())

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

	return c.JSON(sensorData)
}
