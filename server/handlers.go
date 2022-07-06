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
	SMS       [][]SMSData              `json:"sms"`
	MMS       [][]MMSData              `json:"mms"`
	VoiceCall []VoiceCallData          `json:"voice_call"`
	Email     map[string][][]EmailData `json:"email"`
	Billing   BillingData              `json:"billing"`
	Support   []int                    `json:"support"`
	Incidents []IncidentData           `json:"incident"`
}

func getResultHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res := ResultT{}
		d, err := getResultData()
		if err != nil {
			res.Error = "Error on collect data" + err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		if d.SMS != nil && d.MMS != nil && d.VoiceCall != nil && d.Email != nil && d.Support != nil && d.Incidents != nil {
			res.Status = true
		} else {
			res.Error = "Error on collect data"
		}
		res.Data = d
		json.NewEncoder(w).Encode(res)
	}
}

func getResultData() (ResultSetT, error) {
	res := ResultSetT{}
	sms, err := fetchSMS()
	if err != nil {
		return res, err
	}
	sort.Slice(sms, func(i, j int) bool {
		return sms[i].Provider < sms[j].Provider
	})
	res.SMS = append(res.SMS, sms)

	sms1 := make([]SMSData, len(sms))
	copy(sms1, sms)
	sort.Slice(sms1, func(i, j int) bool {
		return sms1[i].Country < sms1[j].Country
	})
	res.SMS = append(res.SMS, sms1)

	mms, err := fetchMMS()
	if err != nil {
		return res, err
	}
	sort.Slice(mms, func(i, j int) bool {
		return mms[i].Provider < mms[j].Provider
	})
	res.MMS = append(res.MMS, mms)

	mms1 := make([]MMSData, len(mms))
	copy(mms1, mms)
	sort.Slice(mms1, func(i, j int) bool {
		return mms1[i].Country < mms1[j].Country
	})
	res.MMS = append(res.MMS, mms1)

	voiceCall, err := fetchVoice()
	if err != nil {
		return res, err
	}
	res.VoiceCall = voiceCall

	//emails, err := fetchEmails()
	//if err != nil {
	//	return res, err
	//}
	//var emailsRes map[string][]EmailData
	//for _, v := range emails {
	//	if emailsRes[v.Country]
	//}
	//fmt.Println(emails)

	return res, nil
}
