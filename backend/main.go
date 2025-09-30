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

type Issue struct {
	Message     string                 `json:"message"`
	Severity    string                 `json:"severity"`
	Component   string                 `json:"component"`
	ProblemData map[string]interface{} `json:"problemData"`
}

type AuditRequest struct {
	Content     string   `json:"content"`
	ActiveRules []string `json:"activeRules"`
}

type AuditResponse struct {
	Status string  `json:"status"`
	Issues []Issue `json:"issues"`
}

func recursiveAudit(data interface{}, issues *[]Issue, activeRulesMap map[string]bool) int {
	var h1Count int

	if m, ok := data.(map[string]interface{}); ok {

		if component, exists := m["component"].(string); exists {
			if activeRulesMap["h1_count"] && component == "heading" {
				if level, levelExists := m["heading_level"].(string); levelExists && level == "h1" {
					h1Count++
				}
			}

			if activeRulesMap["alt_text"] && component == "image" {
				if _, altExists := m["alt"].(string); !altExists || m["alt"].(string) == "" {
					*issues = append(*issues, Issue{
						Message:     "Image is missing mandatory Alt Text for screen readers.",
						Severity:    "Error",
						Component:   component,
						ProblemData: m,
					})
				}
			}
		}

		for key, value := range m {
			if activeRulesMap["lorem_ipsum"] && key == "text" {
				if s, isString := value.(string); isString {
					if strings.Contains(strings.ToLower(s), "lorem ipsum") {
						problemData := map[string]interface{}{}
						for k, v := range m {
							problemData[k] = v
						}
						problemData[key] = "!!! LOREM IPSUM FOUND !!!"

						*issues = append(*issues, Issue{
							Message:     "Placeholder text ('lorem ipsum') found in a text field. Needs real content.",
							Severity:    "Warning",
							Component:   m["component"].(string),
							ProblemData: problemData,
						})
					}
				}
			}

			h1Count += recursiveAudit(value, issues, activeRulesMap)
		}
	}

	if arr, ok := data.([]interface{}); ok {
		for _, item := range arr {
			h1Count += recursiveAudit(item, issues, activeRulesMap)
		}
	}

	return h1Count
}

func auditHandler(c *gin.Context) {
	var req AuditRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	issues := []Issue{}
	rulesMap := make(map[string]bool)
	for _, ruleID := range req.ActiveRules {
		rulesMap[ruleID] = true
	}

	var contentMap interface{}
	err := json.Unmarshal([]byte(req.Content), &contentMap)
	if err != nil {
		issues = append(issues, Issue{Message: fmt.Sprintf("Error parsing Storyblok content JSON: %v", err), Severity: "Error"})
		c.JSON(http.StatusBadRequest, AuditResponse{Status: "Fail", Issues: issues})
		return
	}

	h1Count := recursiveAudit(contentMap, &issues, rulesMap)

	if rulesMap["h1_count"] {
		if h1Count == 0 {
			issues = append(issues, Issue{
				Message:   "Page is missing a mandatory H1 main heading. Only one is allowed for proper structure.",
				Severity:  "Error",
				Component: "Global Structure",
			})
		} else if h1Count > 1 {
			issues = append(issues, Issue{
				Message:   fmt.Sprintf("Page has %d H1 headings. Only one H1 is allowed per page for proper structure.", h1Count),
				Severity:  "Error",
				Component: "Global Structure",
			})
		}
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

	fmt.Println("Gin server is running on https://storyblok-accessibility-auditor.onrender.com/")
	log.Fatal(router.Run("https://storyblok-accessibility-auditor.onrender.com/"))
}
