package visitor

import (
	memory_retention "github.com/yuto51942/memory-retention"
)

// Get the topic.
//
// Arguments:
// 	token {string} - Session token.
//
// Returns:
// 	{string} - Topic.
func GetTopic(token string) (string, error) {
	return memory_retention.GetTopic(token)
}

// Set ansert.
//
// Arguments:
// 	token {string} - Session token.
// 	answer {string} - Answer.
func AddAnswer(token string, answer string) error {
	return memory_retention.AddAnswer(token, answer)
}
