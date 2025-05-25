package ai

import (
	"fmt"
)

func CreateLogger(logString string) {
	fmt.Println(logString)
}

func GeneratePrompt(prompt string, outputFormat string) string {
	// This function is a placeholder for generating a prompt.
	// In a real application, this would interface with an AI model.
	context := fmt.Sprintf(`
		Context: You are a fitness app that have features that create personalised workout routine and counts calorie intake.

		Action: %s

		Output: JSON format with these suggested keys:
		%s
	`, prompt, outputFormat)
	return context
}
