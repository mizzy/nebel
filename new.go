package nebel

import (
	"fmt"
	"os"
	"time"
)

func CreateNewPost(title string) error {
	now := time.Now()
	date := now.Format("2006-01-02")
	location, _ := time.LoadLocation("Asia/Tokyo")
	dateTime := now.In(location).Format("2006-01-02 15:04:05 +0900")

	template := `---
title: %s
date: %s
---
	`

	template = fmt.Sprintf(template, title, dateTime)

	err := os.WriteFile(fmt.Sprintf("posts/%s-%s.markdown", date, title), []byte(template), 0644)
	if err != nil {
		return err
	}

	return nil
}
