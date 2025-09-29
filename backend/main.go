package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type AuditRequest struct {
	Content string `json:"content"`
}

type AuditResponse struct {
	Status string   `json:"status"`
	Issues []string `json:"issues"`
}

func recursiveAudit(data interface{}, issues *[]string, h1Count *int) {
	if m, ok := data.(map[string]interface{}); ok {
		if component, exists := m["component"].(string); exists && component == "heading" {
			if level, levelExists := m["heading_level"].(string); levelExists && level == "h1" {
				*h1Count++
			}
		}

		for key, value := range m {
			if key == "text" {
				if s, isString := value.(string); isString {
					if strings.Contains(strings.ToLower(s), "lorem ipsum") {
						*issues = append(*issues, "Placeholder text ('lorem ipsum') found in a text field. Needs real content.")
					}
				}
			}
			recursiveAudit(value, issues, h1Count)
		}
	}

	if arr, ok := data.([]interface{}); ok {
		for _, item := range arr {
			recursiveAudit(item, issues, h1Count)
		}
	}
}

func auditHandler(c *gin.Context) {

	issues := []string{}
	h1Count := 0

	var req AuditRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var contentMap interface{}
	err := json.Unmarshal([]byte(req.Content), &contentMap)
	if err != nil {
		issues = append(issues, fmt.Sprintf("Error parsing Storyblok content JSON: %v", err))
	} else {
		recursiveAudit(contentMap, &issues, &h1Count)
	}

	if h1Count == 0 {
		issues = append(issues, "Page is missing a mandatory H1 main heading.")
	} else if h1Count > 1 {
		issues = append(issues, fmt.Sprintf("Page has %d H1 headings. Only one H1 is allowed per page for proper structure.", h1Count))
	}

	response := AuditResponse{
		Status: "Pass",
		Issues: issues,
	}
	if len(issues) > 0 {
		response.Status = "Fail"
	}

	c.JSON(http.StatusOK, response)
}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/api/audit", auditHandler)

	fmt.Println("Gin server is running on http://localhost:8080")
	log.Fatal(router.Run(":8080"))
}
