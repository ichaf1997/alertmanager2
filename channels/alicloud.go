package channels

import (
	"bytes"
	"net/http"
	"os"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dyvmsapi"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SmsConfig struct {
	Phones       []string `json:"phones"`
	SignName     string   `json:"sign_name"`
	TemplateCode string   `json:"template_code"`
	AccessKeyId  string   `json:"access_key_id"`
	AccessSecret string   `json:"access_secret"`
	TemplateName string   `json:"template_name"`
}

type VmsConfig struct {
	Phones         []string `json:"phones"`
	TemplateCode   string   `json:"template_code"`
	AccessKeyId    string   `json:"access_key_id"`
	AccessSecret   string   `json:"access_secret"`
	TemplateName   string   `json:"template_name"`
	CallShowNumber string   `json:"call_show_number"`
}

func SendFastSms(c *gin.Context) {

	phones, Isphone := c.GetQueryArray("phone")
	templatecode, Istemplatecode := c.GetQuery("templatecode")
	ak, Isak := os.LookupEnv("ALICLOUD_AK")
	as, Isas := os.LookupEnv("ALICLOUD_AS")

	if !Isphone {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query string: 'phone'"})
		return
	} else if !Istemplatecode {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query string: 'templatecode'"})
		return
	}
	if !Isak || !Isas {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AccessKeyId or AccessSecret Not set On the environment, Please contact administrator for help"})
		log.Debug("AccessKeyId or AccessSecret Not set On the environment")
		return
	}

	cfg := SmsConfig{
		Phones:       phones,
		SignName:     "太平洋互联网",
		TemplateCode: templatecode,
		AccessKeyId:  ak,
		AccessSecret: as,
	}

	var simpletextnotification SimpleTextNotification
	if err := c.BindJSON(&simpletextnotification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errmsg": err.Error(), "reason": "request body is invalid"})
		return
	}

	message := simpletextnotification.Content

	client, _ := dysmsapi.NewClientWithAccessKey("cn-shenzhen", cfg.AccessKeyId, cfg.AccessSecret)
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = cfg.SignName
	request.TemplateCode = cfg.TemplateCode

	// Must Setting `summary` variables in Ali SMS template
	request.TemplateParam = `{"summary":"` + message + `"}`

	var ErrPhones []string
	log.WithFields(
		logrus.Fields{
			"content":   message,
			"smsconfig": cfg,
		},
	).Info("Display Request query and body")
	for _, phonenumber := range phones {
		request.PhoneNumbers = phonenumber
		response, err := client.SendSms(request)
		if err != nil {
			log.WithFields(
				logrus.Fields{
					"phone": phonenumber,
					"err":   err.Error(),
				},
			).Warn("calling alisms api failed")
			ErrPhones = append(ErrPhones, phonenumber)
			continue
		}
		if strings.ToLower(response.Code) == "ok" {
			log.WithFields(
				logrus.Fields{
					"phone": phonenumber,
				},
			).Info("send alisms successful")
		} else {
			log.WithFields(
				logrus.Fields{
					"phone":     phonenumber,
					"code":      response.Code,
					"reason":    response.Message,
					"bizid":     response.BizId,
					"requestid": response.RequestId,
				},
			).Warn("send alisms failed")
			ErrPhones = append(ErrPhones, phonenumber)
		}
	}

	c.JSON(http.StatusOK, gin.H{"errPhones": strings.Join(ErrPhones, ","), "all": len(phones), "err": len(ErrPhones)})

}

func SendFastVms(c *gin.Context) {

	phones, Isphone := c.GetQueryArray("phone")
	templatecode, Istemplatecode := c.GetQuery("templatecode")
	callnumber, Iscallnumber := c.GetQuery("callnumber")
	ak, Isak := os.LookupEnv("ALICLOUD_AK")
	as, Isas := os.LookupEnv("ALICLOUD_AS")

	if !Isphone {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query string: 'phone'"})
		return
	} else if !Istemplatecode {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query string: 'templatecode'"})
		return
	}
	if !Isak || !Isas {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AccessKeyId or AccessSecret Not set On the environment, Please contact administrator for help"})
		log.Debug("AccessKeyId or AccessSecret Not set On the environment")
		return
	}

	cfg := VmsConfig{
		Phones:         phones,
		TemplateCode:   templatecode,
		AccessKeyId:    ak,
		AccessSecret:   as,
		CallShowNumber: callnumber,
	}

	var simpletextnotification SimpleTextNotification
	if err := c.BindJSON(&simpletextnotification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errmsg": err.Error(), "reason": "request body is invalid"})
		return
	}

	message := simpletextnotification.Content

	client, _ := dyvmsapi.NewClientWithAccessKey("cn-shenzhen", cfg.AccessKeyId, cfg.AccessSecret)
	request := dyvmsapi.CreateSingleCallByTtsRequest()
	request.Scheme = "https"
	request.TtsCode = cfg.TemplateCode
	if Iscallnumber {
		request.CalledShowNumber = cfg.CallShowNumber
	}

	// Must Setting `summary` variables in Ali VMS template
	request.TtsParam = `{"summary":"` + message + `"}`

	var ErrPhones []string
	log.WithFields(
		logrus.Fields{
			"content":   simpletextnotification,
			"vmsconfig": cfg,
		},
	).Info("Display Request query and body")
	for _, phonenumber := range phones {
		request.CalledNumber = phonenumber
		response, err := client.SingleCallByTts(request)
		if err != nil {
			log.WithFields(
				logrus.Fields{
					"phone": phonenumber,
					"err":   err.Error(),
				},
			).Warn("calling alivms api failed")
			ErrPhones = append(ErrPhones, phonenumber)
			continue
		}
		if strings.ToLower(response.Code) == "ok" {
			log.WithFields(
				logrus.Fields{
					"phone": phonenumber,
				},
			).Info("send alivms successful")
		} else {
			log.WithFields(
				logrus.Fields{
					"phone":     phonenumber,
					"code":      response.Code,
					"reason":    response.Message,
					"callid":    response.CallId,
					"requestid": response.RequestId,
				},
			).Warn("send alivms failed")
			ErrPhones = append(ErrPhones, phonenumber)
		}
	}

	c.JSON(http.StatusOK, gin.H{"errPhones": strings.Join(ErrPhones, ","), "all": len(phones), "err": len(ErrPhones)})

}

func SendVms(c *gin.Context) {

	phones, Isphone := c.GetQueryArray("phone")
	templatecode, Istemplatecode := c.GetQuery("templatecode")
	callnumber, Iscallnumber := c.GetQuery("callnumber")
	ak, Isak := c.GetQuery("ak")
	as, Isas := c.GetQuery("as")
	tpname, Istpname := c.GetQuery("tpname")

	if !Isphone {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query string: 'phone'"})
		return
	} else if !Istemplatecode {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query string: 'templatecode'"})
		return
	} else if !Isak {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query string: 'ak'"})
		return
	} else if !Isas {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query string: 'as'"})
		return
	} else if !Istpname {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query string: 'tpname'"})
		return
	}
	vmsconfig := VmsConfig{
		Phones:         phones,
		TemplateCode:   templatecode,
		AccessKeyId:    ak,
		AccessSecret:   as,
		TemplateName:   tpname,
		CallShowNumber: callnumber,
	}

	var notification Notification
	if err := c.BindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errmsg": err.Error(), "reason": "request body is invalid"})
		return
	}
	tmpl, err := ParseTemplates(tmplDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Debug(err.Error())
		return
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, vmsconfig.TemplateName, notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Debug(err.Error())
		return
	}

	client, _ := dyvmsapi.NewClientWithAccessKey("cn-shenzhen", vmsconfig.AccessKeyId, vmsconfig.AccessSecret)
	request := dyvmsapi.CreateSingleCallByTtsRequest()
	request.Scheme = "https"
	request.TtsCode = vmsconfig.TemplateCode
	if Iscallnumber {
		request.CalledShowNumber = vmsconfig.CallShowNumber
	}

	// Must Setting `summary` variables in Ali VMS template
	request.TtsParam = `{"summary":"` + buf.String() + `"}`

	var ErrPhones []string
	log.WithFields(
		logrus.Fields{
			"notification": notification,
			"vmsconfig":    vmsconfig,
		},
	).Info("Display Request query and body")
	for _, phonenumber := range phones {
		request.CalledNumber = phonenumber
		response, err := client.SingleCallByTts(request)
		if err != nil {
			log.WithFields(
				logrus.Fields{
					"phone": phonenumber,
					"err":   err.Error(),
				},
			).Warn("calling alivms api failed")
			ErrPhones = append(ErrPhones, phonenumber)
			continue
		}
		if strings.ToLower(response.Code) == "ok" {
			log.WithFields(
				logrus.Fields{
					"phone": phonenumber,
				},
			).Info("send alivms successful")
		} else {
			log.WithFields(
				logrus.Fields{
					"phone":     phonenumber,
					"code":      response.Code,
					"reason":    response.Message,
					"callid":    response.CallId,
					"requestid": response.RequestId,
				},
			).Warn("send alivms failed")
			ErrPhones = append(ErrPhones, phonenumber)
		}
	}

	c.JSON(http.StatusOK, gin.H{"errPhones": strings.Join(ErrPhones, ","), "all": len(phones), "err": len(ErrPhones)})
}

func SendSms(c *gin.Context) {

	phones, Isphone := c.GetQueryArray("phone")
	sn, Issn := c.GetQuery("signname")
	templatecode, Istemplatecode := c.GetQuery("templatecode")
	ak, Isak := c.GetQuery("ak")
	as, Isas := c.GetQuery("as")
	tpname, Istpname := c.GetQuery("tpname")

	if !Isphone {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query string: 'phone'"})
		return
	} else if !Istemplatecode {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query string: 'templatecode'"})
		return
	} else if !Issn {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query string: 'signname'"})
		return
	} else if !Isak {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query string: 'ak'"})
		return
	} else if !Isas {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query string: 'as'"})
		return
	} else if !Istpname {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query string: 'tpname'"})
		return
	}
	smsconfig := SmsConfig{
		Phones:       phones,
		SignName:     sn,
		TemplateCode: templatecode,
		AccessKeyId:  ak,
		AccessSecret: as,
		TemplateName: tpname,
	}

	var notification Notification
	if err := c.BindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errmsg": err.Error(), "reason": "request body is invalid"})
		return
	}
	tmpl, err := ParseTemplates(tmplDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Debug(err.Error())
		return
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, smsconfig.TemplateName, notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Debug(err.Error())
		return
	}

	client, _ := dysmsapi.NewClientWithAccessKey("cn-shenzhen", smsconfig.AccessKeyId, smsconfig.AccessSecret)
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = smsconfig.SignName
	request.TemplateCode = smsconfig.TemplateCode

	// Must Setting `summary` variables in Ali SMS template
	request.TemplateParam = `{"summary":"` + buf.String() + `"}`

	var ErrPhones []string
	log.WithFields(
		logrus.Fields{
			"notification": notification,
			"smsconfig":    smsconfig,
		},
	).Info("Display Request query and body")
	for _, phonenumber := range phones {
		request.PhoneNumbers = phonenumber
		response, err := client.SendSms(request)
		if err != nil {
			log.WithFields(
				logrus.Fields{
					"phone": phonenumber,
					"err":   err.Error(),
				},
			).Warn("calling alisms api failed")
			ErrPhones = append(ErrPhones, phonenumber)
			continue
		}
		if strings.ToLower(response.Code) == "ok" {
			log.WithFields(
				logrus.Fields{
					"phone": phonenumber,
				},
			).Info("send alisms successful")
		} else {
			log.WithFields(
				logrus.Fields{
					"phone":     phonenumber,
					"code":      response.Code,
					"reason":    response.Message,
					"bizid":     response.BizId,
					"requestid": response.RequestId,
				},
			).Warn("send alisms failed")
			ErrPhones = append(ErrPhones, phonenumber)
		}
	}

	c.JSON(http.StatusOK, gin.H{"errPhones": strings.Join(ErrPhones, ","), "all": len(phones), "err": len(ErrPhones)})

}
