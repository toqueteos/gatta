package gatta

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	TokenURL     = "https://api.cognitive.microsoft.com/sts/v1.0/issueToken"
	TranslateURL = "http://api.microsofttranslator.com/v2/Http.svc/Translate"
)

type Translator struct {
	client *http.Client
	key    string
	token  string
}

func New(key string) (*Translator, error) {
	tr := &Translator{
		client: &http.Client{Timeout: 5 * time.Second},
		key:    key,
	}

	if err := tr.getToken(); err != nil {
		return nil, err
	}

	return tr, nil
}

// TOKEN=$(curl -XPOST -H "Ocp-Apim-Subscription-Key: $SUBSCRIPTION_KEY" https://api.cognitive.microsoft.com/sts/v1.0/issueToken --data "")
func (t *Translator) getToken() error {
	req, err := http.NewRequest(http.MethodPost, TokenURL, nil)
	req.Header.Add("Ocp-Apim-Subscription-Key", t.key)
	if err != nil {
		return err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	t.token = string(body)
	log.Println("token", t.token)
	return nil
}

// Translate is a direct translation of the following cURL command:
//  curl -XGET -H "Authorization: Bearer $TOKEN" -H "Accept: application/xml" 'http://api.microsofttranslator.com/v2/Http.svc/Translate?text=$TEXT&to=$TO'
func (t *Translator) Translate(text, to string) (*http.Response, error) {
	u, err := url.Parse(TranslateURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("text", text)
	q.Add("to", to)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	req.Header.Add("Authorization", "Bearer "+t.token)
	req.Header.Add("Accept", "application/xml")
	if err != nil {
		return nil, err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
