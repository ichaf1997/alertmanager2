package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type WxRobotsConfig struct {
	Keys         []string
	TemplateName string
	ContentType  string
}

type WxworkTextMsg struct {
	MsgType string `json:"msgtype"`
	Text    Text   `json:"text"`
}

type Text struct {
	Content               string   `json:"content"`
	Mentioned_List        []string `json:"mented_list"`
	Mentioned_Mobile_List []string `json:"mented_mobile_list"`
}

type WxworkMarkdownMsg struct {
	MsgType  string   `json:"msgtype"`
	MarkDown MarkDown `json:"markdown"`
}

type MarkDown struct {
	Content string `json:"content"`
}

type WxResponse struct {
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

func SendWxWorkRobot(c *gin.Context) {
	keys, Iskey := c.GetQueryArray("key")
	tpname, Istpname := c.GetQuery("tpname")
	msgtype, Ismsgtype := c.GetQuery("msgtype")

	if !Iskey {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query string: 'key'"})
		return
	} else if !Istpname {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query string: 'tpname'"})
		return
	}

	wxrobotConfig := WxRobotsConfig{
		Keys:         keys,
		TemplateName: tpname,
		ContentType:  msgtype,
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
	err = tmpl.ExecuteTemplate(&buf, wxrobotConfig.TemplateName, notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Debug(err.Error())
		return
	}

	var data []byte
	if !Ismsgtype {
		j, _ := json.Marshal(
			WxworkMarkdownMsg{
				MsgType: "markdown",
				MarkDown: MarkDown{
					Content: buf.String(),
				},
			},
		)
		data = j
	} else if strings.ToLower(msgtype) != "markdown" || strings.ToLower(msgtype) != "text" {
		j, _ := json.Marshal(
			WxworkTextMsg{
				MsgType: "text",
				Text: Text{
					Content: buf.String(),
				},
			},
		)
		data = j
	}

	var ErrKeys []string
	log.WithFields(
		logrus.Fields{
			"notification":  notification,
			"wxrobotConfig": wxrobotConfig,
		},
	).Info("Display Request query and body")
	for _, key := range keys {
		url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", key)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.WithFields(
				logrus.Fields{
					"key": key,
					"err": err.Error(),
				},
			).Warn("calling wxwork-robot api failed")
			ErrKeys = append(ErrKeys, key)
			continue
		}
		respDataRaw, _ := io.ReadAll(resp.Body)

		var wxrsp WxResponse
		_ = json.Unmarshal(respDataRaw, &wxrsp)

		if wxrsp.ErrCode == "0" {
			log.WithFields(
				logrus.Fields{
					"key": key,
				},
			).Info("send wxworkrobot successful")
		} else {
			log.WithFields(
				logrus.Fields{
					"errcode": wxrsp.ErrCode,
					"errmsg":  wxrsp.ErrMsg,
				},
			).Warn("send alisms failed")
			ErrKeys = append(ErrKeys, key)
		}

		c.JSON(http.StatusOK, gin.H{"errPhones": strings.Join(ErrKeys, ","), "all": len(keys), "err": len(ErrKeys)})

	}
}
