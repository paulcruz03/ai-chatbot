package ai

import (
	"fmt"
)

func CreateLogger(logString string) {
	fmt.Println(logString)
}

func GeneratePrompt(prompt string) string {
	// This function is a placeholder for generating a prompt.
	// In a real application, this would interface with an AI model.
	context := fmt.Sprintf(`
		Context: You are a fitness app that have features that create personalised workout routine and counts calorie intake.

		Action: %s

		Output: JSON format with the following keys:
		{
			"response": "The response to the action",
			"status": "success or error",
			"error": "Error message if any"
			"data": [
				{
					"id": "<value>",
					"question": "<value>",
					"options": [
						{
							text: "<description>",
							value: "<value>"
						}
				}
			]
		}
	`, prompt)
	return context
}
