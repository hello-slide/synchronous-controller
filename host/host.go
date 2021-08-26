package host

import (
	"github.com/hello-slide/synchronous-controller/util"
	memory_retention "github.com/yuto51942/memory-retention"
)

// Create a session.
//
// Arguments:\n
// ip {string} - Ip address in user.
//
// Returns:\n
// {string} - token.
func CreateSession(ip string) string {
	token := util.NewDateSeed()
	token = token.AddSeed(ip)

	hash := token.CreateSpecifyLength(5)

	memory_retention.CreateKey(hash)

	return hash
}

// Set topic and delete results.
//
// Arguments:\n
// token {string} - Session token.\n
// data {string} - Topic data.
func SetTopic(token string, data string) error {
	if err := memory_retention.DeleteAnswer(token); err != nil {
		return err
	}
	return memory_retention.SetTopic(token, data)
}

// Get results.
//
// Arguments;\n
// token {string} - Session token.\n
//
// Returns:
// {[]string} - results.
func GetResult(token string) ([]string, error) {
	return memory_retention.GetAnswer(token)
}

// Close session.
//
// Arguments:\n
// token {string} - Session token.
func Close(token string) error {
	return memory_retention.DeleteKey(token)
}
