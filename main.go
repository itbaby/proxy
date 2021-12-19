// main.go
package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gonutz/wui/v2"
	"github.com/robfig/cron/v3"
)

func getRemote() {

}
func updateRemote() (string, error) {
	resp, err := http.Get("https://gitee.com/ineo6/hosts/raw/master/hosts")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		wui.MessageBoxError("错误", "Github链接失败")
		return "", errors.New("Github链接失败")
	}
	ioutil.WriteFile("C:/Windows/System32/drivers/etc/hosts", body, 0644)
	return string(body), nil
}
func main() {
	c := cron.New()

	windowFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	window := wui.NewWindow()
	window.SetFont(windowFont)
	window.SetInnerSize(575, 284)
	window.SetTitle("自动刷新GitHub Host工具")
	window.SetHasMaxButton(false)
	window.SetResizable(false)

	slider1 := wui.NewSlider()
	slider1.SetBounds(11, 7, 552, 45)
	slider1.SetArrowIncrement(10)
	slider1.SetMouseIncrement(10)
	slider1.SetCursorPosition(10)
	slider1.SetMinMax(10, 60)
	slider1.SetTickFrequency(10)
	window.Add(slider1)

	textEdit1 := wui.NewTextEdit()
	textEdit1.SetBounds(21, 58, 530, 199)
	textEdit1.SetText("当前为:10 分钟间隔 \r\n正在准备更新... ")
	window.Add(textEdit1)
	label1 := wui.NewLabel()
	label1.SetBounds(442, 263, 119, 16)
	label1.SetText("82500583@github.com")
	window.Add(label1)

	slider1.SetOnChange(func(cursor int) {
		c.Stop()
		for _, e := range c.Entries() {
			c.Remove(e.ID)
		}

		c.AddFunc("@every "+strconv.Itoa(cursor)+"m", func() {
			textEdit1.SetText("刷新时间: " + strconv.Itoa(cursor) + " 分钟")
			hosts, err := updateRemote()
			if err != nil {
				wui.MessageBoxError("错误", "Github链接失败")
			}
			textEdit1.SetText(strings.ReplaceAll(hosts, "\n", "\r\n"))
		})
		c.Start()
	})

	c.AddFunc("@every 10m", func() {
		hosts, err := updateRemote()
		if err != nil {
			wui.MessageBoxError("错误", "Github链接失败")
		}
		textEdit1.SetText(strings.ReplaceAll(hosts, "\n", "\r\n"))
	})
	c.Start()
	window.Show()
}
