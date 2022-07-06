package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	//fetchSMS()
	//fetchMMS()
	//fetchVoices()
	//fetchEmails()
	//fetchBillings()
	//fetchSupport()
	//fetchIncedents()
}

type IncidentData struct {
	Topic  string `json:"topic"`
	Status string `json:"status"` // возможные статусы: active и closed
}

func fetchIncedents() []IncidentData {
	result := []IncidentData{}

	resp, err := http.Get("http://localhost:8383/accendent")
	if err != nil {
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode == 500 {
		fmt.Println("ошибка на удаленном сервере doSupport")
	}
	if resp.StatusCode == 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return result
		}
		json.Unmarshal(b, &result)
	}

	return result
}

type SupportData struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}

func fetchSupport() ([]SupportData, error) {
	data := []SupportData{}
	resp, err := http.Get("http://localhost:8383/support")
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 500 {
		fmt.Println("ошибка на удаленном сервере doSupport")
	}
	if resp.StatusCode == 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return data, err
		}
		json.Unmarshal(b, &data)
	}

	return data, nil
}

type BillingData struct {
	CreateCustomer bool
	Purchase       bool
	Payout         bool
	Recurring      bool
	FraudControl   bool
	CheckoutPage   bool
}

func fetchBillings() (BillingData, error) {
	billingData := BillingData{}
	f, err := os.Open("billing.data")
	if err != nil {
		return billingData, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return billingData, err
	}

	billingData = billingParse(b)
	return billingData, nil
}

func billingParse(b []byte) BillingData {
	return BillingData{
		b[0] == 49,
		b[1] == 49,
		b[2] == 49,
		b[3] == 49,
		b[4] == 49,
		b[5] == 49,
	}
}

type EmailData struct {
	Country      string
	Provider     string
	DeliveryTime int
}

func fetchEmails() ([]EmailData, error) {
	result := []EmailData{}
	f, err := os.Open("email.data")
	if err != nil {
		return result, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return result, err
	}

	result = emailValidate(b)
	return result, nil
}

func emailValidate(data []byte) []EmailData {
	rows := strings.Split(string(data), "\n")

	var res []EmailData
	for _, row := range rows {
		cols := strings.Split(row, ";")
		if len(cols) != 3 {
			continue
		}

		country, ok := countryCodes[cols[0]]
		if !ok {
			continue
		}

		providers := []string{"Gmail", "Yahoo", "Hotmail", "MSN", "Orange", "Comcast", "AOL", "Live", "RediffMail",
			"GMX", "Protonmail", "Yandex", "Mail.ru"}

		flag := false
		for _, p := range providers {
			if p == cols[1] {
				flag = true
				break
			}
		}
		if !flag {
			continue
		}

		deliveryTime, err := strconv.ParseInt(cols[2], 10, 64)
		if err != nil {
			continue
		}

		d := EmailData{
			country,
			cols[1],
			int(deliveryTime),
		}
		res = append(res, d)
	}
	fmt.Println("emailValidate fails count:", len(rows)-len(res)-1)
	return res
}

type VoiceCallData struct {
	Country             string
	Bandwidth           string
	ResponseTime        string
	Provider            string
	ConnectionStability float32
	TTFB                int
	VoicePurity         int
	MedianOfCallsTime   int
}

func fetchVoice() ([]VoiceCallData, error) {
	result := []VoiceCallData{}
	f, err := os.Open("voice.data")
	if err != nil {
		return result, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return result, err
	}

	result = voiceValidate(b)
	return result, nil
}

func voiceValidate(data []byte) []VoiceCallData {
	rows := strings.Split(string(data), "\n")

	var res []VoiceCallData
	for _, row := range rows {
		cols := strings.Split(row, ";")
		if len(cols) != 8 {
			continue
		}

		country, ok := countryCodes[cols[0]]
		if !ok {
			continue
		}

		if !(cols[3] == "TransparentCalls" || cols[3] == "E-Voice" || cols[3] == "JustPhone") {
			continue
		}

		cs, err := strconv.ParseFloat(cols[4], 32)
		if err != nil {
			continue
		}
		ttfb, err := strconv.ParseInt(cols[5], 10, 64)
		if err != nil {
			continue
		}
		voicePurity, err := strconv.ParseInt(cols[6], 10, 64)
		if err != nil {
			continue
		}
		medianOfCallsTime, err := strconv.ParseInt(cols[7], 10, 64)
		if err != nil {
			continue
		}

		d := VoiceCallData{
			country,
			cols[1],
			cols[2],
			cols[3],
			float32(cs),
			int(ttfb),
			int(voicePurity),
			int(medianOfCallsTime),
		}
		res = append(res, d)
	}
	fmt.Println("voiceValidate fails count:", len(rows)-len(res)-1)
	return res
}

type MMSData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

func fetchMMS() ([]MMSData, error) {
	data := []MMSData{}
	resp, err := http.Get("http://localhost:8383/mms")
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	var validData []MMSData
	if resp.StatusCode == 500 {
		fmt.Println("ошибка на удаленном сервере doMMS")
	}
	if resp.StatusCode == 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return data, err
		}

		json.Unmarshal(b, &data)

		for _, mms := range data {
			c, ok := countryCodes[mms.Country]
			if !ok {
				continue
			}
			mms.Country = c

			if !(mms.Provider == "Topolo" || mms.Provider == "Rond" || mms.Provider == "Kildy") {
				continue
			}
			validData = append(validData, mms)
		}
	}

	fmt.Println(validData)
	return validData, nil
}

type SMSData struct {
	Country      string
	Bandwidth    string
	ResponseTime string
	Provider     string
}

func fetchSMS() ([]SMSData, error) {
	result := []SMSData{}
	f, err := os.Open("sms.data")
	if err != nil {
		return result, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return result, err
	}

	result = smsValidate(b)
	return result, nil
}

func smsValidate(data []byte) []SMSData {
	rows := strings.Split(string(data), "\n")

	var res []SMSData
	for _, row := range rows {
		elems := strings.Split(row, ";")
		if len(elems) != 4 {
			continue
		}

		country, ok := countryCodes[elems[0]]
		if !ok {
			continue
		}

		if !(elems[3] == "Topolo" || elems[3] == "Rond" || elems[3] == "Kildy") {
			continue
		}

		d := SMSData{
			country,
			elems[1],
			elems[2],
			elems[3],
		}
		res = append(res, d)
	}
	fmt.Println("smsValidate fails count:", len(rows)-len(res)-1)
	return res
}
