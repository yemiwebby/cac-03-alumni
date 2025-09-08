package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	colName = "FULL NAME (SCHOOL SURNAME FIRST)"
	colDOB  = "DATE OF BIRTH"
)

type Row struct {
	Name string
	MM   int
	DD   int
}

func main() {
	// Load environment variables from .env file if it exists (for local development)
	// In production (GitHub Actions), environment variables are set directly
	if err := godotenv.Load(); err != nil {
		// Silently ignore if .env file doesn't exist (normal in production)
		if !strings.Contains(err.Error(), "no such file or directory") {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
	}

	csvPath := flag.String("csv", envOr("CSV_PATH", "file/CAC-03-SET-BIODATA-FORM.csv"), "path to the CSV file containing biodata")
	lookahead := flag.Int("lookahead", 1, "days ahead to remind")
	tzName := flag.String("tz", envOr("TIMEZONE", "Europe/London"), "IANA timezone")
	dry := flag.Bool("dry", false, "print instead of sending")
	monthly := flag.Bool("monthly", false, "send monthly birthday summary instead of daily reminders")
	targetMonth := flag.Int("target-month", 0, "override target month for monthly reports (1-12, 0=auto)")
	flag.Parse()

	// WhatsApp config
	phoneID := mustEnv("WA_PHONE_ID")
	token := mustEnv("WA_TOKEN")
	template := envOr("WA_TEMPLATE", "hello_world")
	lang := envOr("WA_LANG", "en")

	// Recipients (collectors)
	toList := strings.Split(strings.ReplaceAll(mustEnv("WA_TO_LIST"), " ", ""), ",")
	if len(toList) == 0 || toList[0] == "" {
		panic("WA_TO_LIST is empty; set comma-separated E.164 numbers without '+'")
	}

	loc, err := time.LoadLocation(*tzName)
	check(err)

	rows, bad := readRows(*csvPath)
	if len(bad) > 0 {
		fmt.Printf("Skipped %d row(s) with invalid DOB format (expect YYYY-MM-DD)\n", len(bad))
	}

	now := time.Now().In(loc)

	if *monthly {
		sendMonthlyReport(rows, now, phoneID, token, template, lang, toList, *dry, *targetMonth)
		return
	}

	target := now.AddDate(0, 0, *lookahead)
	tm, td := int(target.Month()), target.Day()

	// Find birthdays N days ahead
	var hits []Row
	for _, r := range rows {
		if r.MM == tm && r.DD == td {
			hits = append(hits, r)
		}
	}

	if len(hits) == 0 {
		fmt.Println("No birthdays in window.")
		return
	}

	for _, h := range hits {
		birthdayDate := fmt.Sprintf("%04d-%02d-%02d", now.Year(), h.MM, h.DD) // message param
		// Send one message to each collector
		for _, to := range toList {
			if *dry {
				fmt.Printf("[DRY] to=%s | %s | %s\n", to, h.Name, birthdayDate)
				continue
			}
			if err := sendTemplate(phoneID, token, template, lang, to, h.Name, birthdayDate); err != nil {
				fmt.Printf("send error to %s: %v\n", to, err)
			} else {
				fmt.Printf("sent to %s: %s (%s)\n", to, h.Name, birthdayDate)
			}
			time.Sleep(250 * time.Millisecond)
		}
	}

	fmt.Printf("Done. Birthdays matched: %d, messages sent per-collector: %d each.\n", len(hits), len(toList))
}

func readRows(path string) ([]Row, []string) {
	f, err := os.Open(path)
	check(err)
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	header, err := r.Read()
	check(err)

	find := func(needle string) int {
		for i, h := range header {
			if strings.EqualFold(strings.TrimSpace(h), needle) {
				return i
			}
		}
		return -1
	}
	iName, iDOB := find(colName), find(colDOB)
	if iName < 0 || iDOB < 0 {
		check(errors.New("CSV missing required headers: " + colName + " and/or " + colDOB))
	}

	var out []Row
	var bad []string
	seen := make(map[string]bool) // Track duplicates by name+birthday

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		check(err)

		name := strings.TrimSpace(rec[iName])
		dob := strings.TrimSpace(rec[iDOB])
		if name == "" || dob == "" {
			continue
		}
		// Expect YYYY-MM-DD; ignore the year; just get month/day
		if len(dob) != 10 || dob[4] != '-' || dob[7] != '-' {
			bad = append(bad, dob)
			continue
		}
		mm, dd := dob[5:7], dob[8:10]
		m, d := atoi(mm), atoi(dd)
		if m < 1 || m > 12 || d < 1 || d > 31 {
			bad = append(bad, dob)
			continue
		}

		// Create a unique key for deduplication
		key := strings.ToLower(name) + "-" + mm + "-" + dd
		if seen[key] {
			continue // Skip duplicate
		}
		seen[key] = true

		out = append(out, Row{Name: name, MM: m, DD: d})
	}
	return out, bad
}

func sendTemplate(phoneID, token, template, lang, to, param1, param2 string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v20.0/%s/messages", phoneID)
	body := map[string]any{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "template",
		"template": map[string]any{
			"name":     template,
			"language": map[string]string{"code": lang},
			"components": []map[string]any{{
				"type": "body",
				"parameters": []map[string]any{
					{"type": "text", "text": param1},
					{"type": "text", "text": param2},
				},
			}},
		},
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("meta status %s", resp.Status)
	}
	return nil
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}

func envOr(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic("missing env: " + k)
	}
	return v
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func sendMonthlyReport(rows []Row, now time.Time, phoneID, token, template, lang string, toList []string, dry bool, targetMonth int) {
	// Get next month's birthdays - fix for month-end edge case
	// Calculate next month properly to avoid August 31 -> October issue
	var nextMonth time.Time
	var reportMonth int

	if targetMonth > 0 && targetMonth <= 12 {
		// Use specified target month
		nextMonth = time.Date(now.Year(), time.Month(targetMonth), 1, 0, 0, 0, 0, now.Location())
		reportMonth = targetMonth
	} else {
		// Use next month (default behavior)
		nextMonth = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
		reportMonth = int(nextMonth.Month())
	}

	var monthlyBirthdays []Row
	for _, r := range rows {
		if r.MM == reportMonth {
			monthlyBirthdays = append(monthlyBirthdays, r)
		}
	}

	if len(monthlyBirthdays) == 0 {
		fmt.Printf("No birthdays in %s %d.\n", nextMonth.Month().String(), nextMonth.Year())
		return
	}

	// Sort birthdays by day
	for i := 0; i < len(monthlyBirthdays)-1; i++ {
		for j := i + 1; j < len(monthlyBirthdays); j++ {
			if monthlyBirthdays[i].DD > monthlyBirthdays[j].DD {
				monthlyBirthdays[i], monthlyBirthdays[j] = monthlyBirthdays[j], monthlyBirthdays[i]
			}
		}
	}

	// Create birthday list string
	var namesList []string
	for _, b := range monthlyBirthdays {
		namesList = append(namesList, fmt.Sprintf("%s (%d)", b.Name, b.DD))
	}
	birthdayListText := strings.Join(namesList, ", ")

	// Use existing template format
	monthYear := fmt.Sprintf("%s %d", nextMonth.Month().String(), nextMonth.Year())

	for _, to := range toList {
		if dry {
			fmt.Printf("[DRY MONTHLY] to=%s | Monthly Report for %s | %s\n", to, monthYear, birthdayListText)
			continue
		}

		// Send monthly report using template with 2 parameters
		if err := sendTemplate(phoneID, token, template, lang, to, birthdayListText, monthYear); err != nil {
			fmt.Printf("send monthly report error to %s: %v\n", to, err)
		} else {
			fmt.Printf("sent monthly report to %s for %s\n", to, monthYear)
		}
		time.Sleep(250 * time.Millisecond)
	}

	fmt.Printf("Done. Monthly report sent to %d recipients for %d birthdays in %s.\n", len(toList), len(monthlyBirthdays), nextMonth.Month().String())
}
