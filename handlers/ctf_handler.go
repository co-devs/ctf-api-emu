package handlers

import (
	"database/sql"
	"web-service-gin-tut/database"
	"web-service-gin-tut/models"

	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

// GetHeartbeat response with the current status
func GetHeartbeat(c *gin.Context) {
	var heartbeat = models.Heartbeat{Status: "ok"}
	c.JSON(http.StatusOK, heartbeat)
}

func GetEndpoints(c *gin.Context) {
	rows, err := database.DB.Query("SELECT endpoints.team_id, endpoints.service_id, services.service_name, endpoints.hostname FROM endpoints JOIN services ON endpoints.service_id = services.id;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var endpoints []models.Endpoint
	for rows.Next() {
		var endpoint models.Endpoint
		if err:= rows.Scan(&endpoint.TeamID, &endpoint.ServiceID, &endpoint.ServiceName, &endpoint.Hostname); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		endpoints = append(endpoints, endpoint)
	}
	c.JSON(http.StatusOK, endpoints)
}

func GetLiveFlags(c *gin.Context) {
	rows, err := database.DB.Query("SELECT flags.flag_identifier, endpoints.team_id, endpoints.service_id, flags.tick, strftime('%FT%TZ', flags.expiration) AS expiration, endpoints.hostname FROM flags JOIN endpoints ON flags.endpoint_id = endpoints.id;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var flags []models.Flag
	for rows.Next() {
		var flag models.Flag
		if err := rows.Scan(&flag.Identifier, &flag.TeamID, &flag.ServiceID, &flag.Tick, &flag.Expiration, &flag.Hostname); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		flags = append(flags, flag)
	}
	c.JSON(http.StatusOK, flags)
}

// GetSubmittedFlags returns all flags submitted by the current user's team
func GetSubmittedFlags(c *gin.Context) {
    apiKeyStr := c.GetHeader("team-token")
    rows, err := database.DB.Query("SELECT flags.flag, endpoints.team_id, endpoints.service_id, submitted_flags.timestamp FROM submitted_flags JOIN flags on submitted_flags.flag_id = flags.id JOIN endpoints on flags.endpoint_id = endpoints.id JOIN teams on endpoints.team_id = teams.id WHERE teams.key = ?;", apiKeyStr)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    type Result struct {
        Flag      string `json:"flag"`
        TeamID    int    `json:"team_id"`
        ServiceID int    `json:"service_id"`
        Timestamp string `json:"timestamp"`
    }

    var results []Result
    for rows.Next() {
        var result Result
        if err := rows.Scan(&result.Flag, &result.TeamID, &result.ServiceID, &result.Timestamp); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        results = append(results, result)
    }
    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, results)
}

// PostFlag submits a flag and stores it if it is still a live flag.
func PostFlag(c *gin.Context) {
	var requestBody struct {
		FlagIn string `json:"flag_in"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	flag := requestBody.FlagIn
	// log.Printf("Received flag: %s", flag)
	// validate flag format
	flagPattern := `^mctf\{[a-zA-Z0-9+/=]+\}$`
	matched, err := regexp.MatchString(flagPattern, flag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking flag format"})
        return
	}
	if !matched {
		// log.Printf("Flag did not match pattern: %s", flag)
		c.JSON(http.StatusBadRequest, gin.H{"response": "invalid_format"})
		return
	}

	// get team ID from API token
    apiKeyStr := c.GetHeader("team-token")
	var teamID int
	err = database.DB.QueryRow("SELECT id FROM teams WHERE key = ?;", apiKeyStr).Scan(&teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if flag is still active
	var live bool
	err = database.DB.QueryRow(`SELECT EXISTS (SELECT 1 FROM flags WHERE flag = ? AND expiration > datetime('now'));`, flag).Scan(&live)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Flag has expired or does not exist"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Check if flag has already been submitted by team
	var exists bool
	err = database.DB.QueryRow(`SELECT EXISTS (SELECT 1 FROM submitted_flags JOIN flags ON submitted_flags.flag_id = flags.id WHERE flags.flag = ? AND submitted_flags.team_id = ?)`, flag, teamID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Print("Error checking for prior submission")
		// log.Printf("flag: %v", flag)
		// log.Printf("teamID: %v", teamID)
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"response": "duplicate"})
		return
	}

	// If conditions have all been met, then insert into submitted_flags
	_, err = database.DB.Exec(`INSERT INTO submitted_flags (flag_id, team_id, timestamp) SELECT id, ?, datetime('now') FROM flags WHERE flag = ?`, teamID, flag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": "success"})

}

// GetStatus returns status information for the current team
func GetStatus(c *gin.Context) {
    // Define the response structure
    type Service struct {
        Hostname     string   `json:"hostname"`
        ServiceID    int      `json:"service_id"`
        ServiceName  string   `json:"service_name"`
        Availability []string `json:"availability"`
    }

    type Team struct {
        TeamName string    `json:"team_name"`
        TeamID   int       `json:"team_id"`
        Services []Service `json:"services"`
    }

    type Response struct {
        Tick  int    `json:"tick"`
        Teams []Team `json:"teams"`
    }

    var response Response

    // Query the current tick
    err := database.DB.QueryRow("SELECT MAX(id) FROM ticks").Scan(&response.Tick)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Query the list of teams
    rows, err := database.DB.Query("SELECT id, name FROM teams")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var teams []Team
    for rows.Next() {
        var team Team
        if err := rows.Scan(&team.TeamID, &team.TeamName); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Query the list of services for the team
        serviceRows, err := database.DB.Query(`
            SELECT endpoints.hostname, services.id, services.service_name
            FROM endpoints
            JOIN services ON endpoints.service_id = services.id
            WHERE endpoints.team_id = ?
        `, team.TeamID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer serviceRows.Close()

        var services []Service
        for serviceRows.Next() {
            var service Service
            if err := serviceRows.Scan(&service.Hostname, &service.ServiceID, &service.ServiceName); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }

            // Query the availability for the last 3 ticks
            availabilityRows, err := database.DB.Query(`
                SELECT status
                FROM status_checks
                WHERE endpoint_id = (SELECT id FROM endpoints WHERE hostname = ?)
                ORDER BY tick DESC
                LIMIT 3
            `, service.Hostname)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            defer availabilityRows.Close()

            var availability []string
            for availabilityRows.Next() {
                var status string
                if err := availabilityRows.Scan(&status); err != nil {
                    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                    return
                }
                availability = append(availability, status)
            }
            service.Availability = availability

            services = append(services, service)
        }
        team.Services = services
        teams = append(teams, team)
    }
    response.Teams = teams

    c.JSON(http.StatusOK, response)
}

// GetAllFlagSubmissions returns all flag submissions
func GetAllFlagSubmissions(c *gin.Context) {
    rows, err := database.DB.Query(`
        SELECT
            submitted_flags.id,
            submitted_flags.timestamp,
            flags.flag_identifier,
            flags.flag,
            teams.name AS team_name,
            endpoints.hostname,
            services.service_name
        FROM submitted_flags
        JOIN flags ON submitted_flags.flag_id = flags.id
        JOIN teams ON submitted_flags.team_id = teams.id
        JOIN endpoints ON flags.endpoint_id = endpoints.id
        JOIN services ON endpoints.service_id = services.id
    `)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    type FlagSubmission struct {
        ID            int    `json:"id"`
        Timestamp     string `json:"timestamp"`
        FlagIdentifier string `json:"flag_identifier"`
        Flag          string `json:"flag"`
        TeamName      string `json:"team_name"`
        Hostname      string `json:"hostname"`
        ServiceName   string `json:"service_name"`
    }

    var submissions []FlagSubmission
    for rows.Next() {
        var submission FlagSubmission
        if err := rows.Scan(&submission.ID, &submission.Timestamp, &submission.FlagIdentifier, &submission.Flag, &submission.TeamName, &submission.Hostname, &submission.ServiceName); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        submissions = append(submissions, submission)
    }

    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, submissions)
}
