package server

import (
	"encoding/json"
	"net/http"
	"sort"
)

type ResultT struct {
	Status bool       `json:"status"` // true, если все этапы сбора
	Data   ResultSetT `json:"data"`   // заполнен, если все этапы сбора
	Error  string     `json:"error"`  // пустая строка если все этапы
}

type ResultSetT struct {
	SMS       [][]SMSData     `json:"sms"`
	MMS       [][]MMSData     `json:"mms"`
	VoiceCall []VoiceCallData `json:"voice_call"`
	Email     [][]EmailData   `json:"email"`
	Billing   BillingData     `json:"billing"`
	Support   []int           `json:"support"`
	Incidents []IncidentData  `json:"incident"`
}

func getResultHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		res := ResultT{}
		d, err := getResultData()
		if err != nil {
			res.Error = "Error on collect data" + err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		if d.SMS != nil && d.MMS != nil && d.VoiceCall != nil && d.Email != nil && d.Support != nil && d.Incidents != nil {
			res.Status = true
			res.Data = d
		} else {
			res.Error = "Error on collect data"
		}

		json.NewEncoder(w).Encode(res)
	}
}

func getResultData() (ResultSetT, error) {
	res := ResultSetT{}
	var err error

	res.SMS, err = getSMS()
	if err != nil {
		return res, err
	}

	res.MMS, err = getMMS()
	if err != nil {
		return res, err
	}

	res.VoiceCall, err = fetchVoice()
	if err != nil {
		return res, err
	}

	emails := map[string][][]EmailData{}
	emails, err = getEmails()
	if err != nil {
		return res, err
	}
	res.Email = emails["Russian Federation"]

	res.Billing, err = fetchBillings()
	if err != nil {
		return res, err
	}

	res.Support, err = getSupports()
	if err != nil {
		return res, err
	}

	res.Incidents, err = fetchIncedents()
	if err != nil {
		return res, err
	}
	sort.Slice(res.Incidents, func(i, j int) bool {
		return res.Incidents[i].Status < res.Incidents[j].Status
	})

	return res, nil
}

func getSupports() ([]int, error) {
	res := []int{}
	supports, err := fetchSupport()
	if err != nil {
		return res, err
	}
	var sum int
	for _, s := range supports {
		sum += s.ActiveTickets
	}
	if sum < 9 {
		res = []int{1}
	} else if sum >= 9 && sum <= 16 {
		res = []int{2}
	} else if sum > 16 {
		res = []int{3}
	}
	res = append(res, sum*60/18)

	return res, nil
}

func getEmails() (map[string][][]EmailData, error) {
	emailsRes := map[string][][]EmailData{}

	emails, err := fetchEmails()
	if err != nil {
		return emailsRes, err
	}

	for _, email := range emails {
		if email.DeliveryTime == 0 {
			continue
		}

		country := email.Country
		if _, ok := emailsRes[country]; !ok {
			emailsRes[country] = make([][]EmailData, 2)
		}

		if len(emailsRes[country][0]) < 3 {
			emailsRes[country][0] = append(emailsRes[country][0], email)
			if len(emailsRes[country][0]) == 3 {
				sort.Slice(emailsRes[country][0], func(i, j int) bool {
					return emailsRes[country][0][i].DeliveryTime < emailsRes[country][0][j].DeliveryTime
				})
			}
		} else {
			placeToInsert := len(emailsRes[country][0])
			for i := range emailsRes[country][0] {
				if emailsRes[country][0][i].DeliveryTime > email.DeliveryTime {
					placeToInsert = i
					break
				}
			}

			var e1, e2 EmailData
			e2 = email
			for i := placeToInsert; i < len(emailsRes[country][0]); i++ {
				e1 = emailsRes[country][0][i]
				emailsRes[country][0][i] = e2
				e2 = e1
			}
		}

		if len(emailsRes[country][1]) < 3 {
			emailsRes[country][1] = append(emailsRes[country][1], email)
			if len(emailsRes[country][1]) == 3 {
				sort.Slice(emailsRes[country][1], func(i, j int) bool {
					return emailsRes[country][1][i].DeliveryTime > emailsRes[country][1][j].DeliveryTime
				})
			}
		} else {
			placeToInsert := len(emailsRes[country][1])
			for i := range emailsRes[country][1] {
				if emailsRes[country][1][i].DeliveryTime < email.DeliveryTime {
					placeToInsert = i
					break
				}
			}

			var e EmailData
			for i := placeToInsert; i < len(emailsRes[country][1]); i++ {
				e = emailsRes[country][1][i]
				emailsRes[country][1][i] = email
				email = e
			}
		}
	}

	return emailsRes, nil
}

func getSMS() ([][]SMSData, error) {
	res := make([][]SMSData, 0)
	sms, err := fetchSMS()
	if err != nil {
		return res, err
	}

	sort.Slice(sms, func(i, j int) bool {
		return sms[i].Provider < sms[j].Provider
	})
	res = append(res, sms)

	sms1 := make([]SMSData, len(sms))
	copy(sms1, sms)
	sort.Slice(sms1, func(i, j int) bool {
		return sms1[i].Country < sms1[j].Country
	})
	res = append(res, sms1)

	return res, nil
}

func getMMS() ([][]MMSData, error) {
	res := make([][]MMSData, 0)
	mms, err := fetchMMS()
	if err != nil {
		return res, err
	}

	sort.Slice(mms, func(i, j int) bool {
		return mms[i].Provider < mms[j].Provider
	})
	res = append(res, mms)

	mms1 := make([]MMSData, len(mms))
	copy(mms1, mms)
	sort.Slice(mms1, func(i, j int) bool {
		return mms1[i].Country < mms1[j].Country
	})
	res = append(res, mms1)

	return res, nil
}
