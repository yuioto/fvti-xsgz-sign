package notify

import (
	"net/http"
	"strings"
)

func SendNtfyMessage(topic string, level string, title string, message string) error {
	atriSleep := "https://i0.hdslb.com/bfs/article/b11ad7419cd98dfb661f23505a996288ef694932.jpg"

	url := "https://ntfy.sh/" + topic

	req, err := http.NewRequest("POST", url, strings.NewReader(message))
	req.Header.Set("Markdown", "yes")
	req.Header.Set("Title", title)
	req.Header.Set("Priority", level)
	req.Header.Set("Tags", "fvti,xsgz,sign")
	req.Header.Set("Attach", atriSleep)
	http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	//if req.StatusCode != http.StatusOK {
	//	return fmt.Errorf("failed to send message: %s", req.Body)
	//}

	return nil
}
