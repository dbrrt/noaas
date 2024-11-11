package readuri

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
)

// Default execCommand function
var execCommand = exec.Command

func readRemoteUriPayload(uri string, script bool) (string, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return "", fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to retrieve content. Status code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	output := ""

	if script {
		// Run the content as a shell script
		cmd := exec.Command("/bin/sh", "-c", string(body))
		combinedOutput, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("error executing script: %v", err)
		}
		output = string(combinedOutput)
	} else {
		output = string(body)
	}

	return output, nil
}
