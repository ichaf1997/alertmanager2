package media

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/ichaf1997/alertmanager2/utils"
	"github.com/sirupsen/logrus"
)

type FeiShuCustomRobot struct {
	MsgTmplName string `form:"msg_tmpl_name" binding:"required"`
	MsgData     any
	Template    *template.Template
	Api         string
}

type FeiShuCustomRobotResponse struct {
	Code int
	Msg  string
	Data any
}

func (feishu FeiShuCustomRobot) SendMsg() error {
	var buf bytes.Buffer
	if err := feishu.Template.ExecuteTemplate(&buf, feishu.MsgTmplName, feishu.MsgData); err != nil {
		logrus.Errorf("Failed to render template [%s]: %v", feishu.MsgTmplName, err)
		return err
	}
	logrus.Debugf("Success to render template [%s]", feishu.MsgTmplName)
	logrus.Debugf("rendered template: %s", buf.String())

	req, err1 := http.NewRequest(http.MethodPost, feishu.Api, &buf)
	if err1 != nil {
		logrus.Errorf("Error creating request: %v", err1)
		return err1
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0")
	resp, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		logrus.Errorf("Error sending request: %v", err2)
		return err2
	}

	var f FeiShuCustomRobotResponse
	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		logrus.Warnf("Error parse feishu response With struct FeiShuCustomRobotResponse: %v", err)
		return nil
	}
	logrus.Debugf("Success parsing feishu response with struct FeiShuCustomRobotResponse: %v", f)
	if f.Code != 0 {
		errMsg := fmt.Sprintf("Failed to send feishu message with code [%d]: %s", f.Code, f.Msg)
		logrus.Errorf(errMsg)
		return errors.New(errMsg)
	} else {
		logrus.Infof("Success send feishu message with code [%d]: %s", f.Code, f.Msg)
	}
	defer resp.Body.Close()
	return nil
}

func HelloWorld(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}

func (h *Handler) FeiShuCustomRobotHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			feishu FeiShuCustomRobot
			rsp    Response
			data   utils.Notification
			medium Medium
		)

		if err1, err2 := c.BindQuery(&feishu), c.BindJSON(&data); err1 == nil && err2 == nil {
			feishu.MsgData = data
			feishu.Api = fmt.Sprintf("https://open.feishu.cn/open-apis/bot/v2/hook/%s", c.Param("key"))
			feishu.Template, _ = h.template.Clone()
			rsp.Data = data
			logrus.Debug(feishu)
		} else {
			rsp.Code = 1
			if err1 != nil {
				msg := fmt.Sprintf("Failed to parse query: %v", err1)
				rsp.Message = msg
				logrus.Warn(msg)
			} else if err2 != nil {
				msg := fmt.Sprintf("Failed to parse post data: %v", err2)
				rsp.Message = msg
				logrus.Warn(msg)
			}
			c.JSON(http.StatusBadRequest, rsp)
			return
		}

		medium = feishu
		if err := medium.SendMsg(); err != nil {
			rsp.Code = 2
			rsp.Message = err.Error()
			c.JSON(http.StatusInternalServerError, rsp)
			return
		}
		rsp.Message = "success send message via FeiShuCustomRobot medium"
		c.JSON(http.StatusOK, rsp)
	}
}
