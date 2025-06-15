package sipapi

import (
	"errors"
	"fmt"
	sip "go-sip/sip/s"
	"net/http"
	"strings"
)

// 云台控制
func DeviceControl(DeviceID string, leftRight, upDown, inOut, moveSpeed, zoomSpeed int) error {

	u, ok := _activeDevices.Load(DeviceID)
	if !ok {
		return errors.New("活跃设备不存在")
	}
	to := u.(Devices)
	channelURI, _ := sip.ParseURI(to.URIStr)
	to.addr = &sip.Address{URI: channelURI}

	_url, _ := sip.ParseSipURI(fmt.Sprintf("sip:%s@%s:%s", to.DeviceID, to.Host, to.Port))
	contact := &sip.Address{
		URI: &_url,
	}

	cmdstr := cmdString(leftRight, upDown, inOut, moveSpeed, zoomSpeed)
	hb := sip.NewHeaderBuilder().SetTo(contact).SetFrom(_serverDevices.addr).AddVia(&sip.ViaHop{
		Params: sip.NewParams().Add("branch", sip.String{Str: sip.GenerateBranch()}),
	}).SetContentType(&sip.ContentTypeXML).SetMethod(sip.MESSAGE).SetContact(contact)
	req := sip.NewRequest("", sip.MESSAGE, to.addr.URI, sip.DefaultSipVersion, hb.Build(), []byte(sip.GenerateDeviceControl(DeviceID, cmdstr)))
	req.SetDestination(to.source)
	tx, err := srv.Request(req)
	if err != nil {
		return err
	}
	response := tx.GetResponse()
	if response.StatusCode() != http.StatusOK {
		return errors.New(response.Reason())
	}

	return nil
}

func cmdString(leftRight, upDown, inOut, moveSpeed, zoomSpeed int) string {
	cmdCode := 0

	// 移动方向位设置
	if leftRight == 2 {
		cmdCode |= 0x01 // 右移
	} else if leftRight == 1 {
		cmdCode |= 0x02 // 左移
	}
	if upDown == 2 {
		cmdCode |= 0x04 // 下移
	} else if upDown == 1 {
		cmdCode |= 0x08 // 上移
	}
	if inOut == 2 {
		cmdCode |= 0x10 // 放大
	} else if inOut == 1 {
		cmdCode |= 0x20 // 缩小
	}

	var builder strings.Builder
	builder.WriteString("A50F01")

	// cmdCode
	builder.WriteString(fmt.Sprintf("%02X", cmdCode))

	// moveSpeed 两次
	moveHex := fmt.Sprintf("%02X", moveSpeed)
	builder.WriteString(moveHex)
	builder.WriteString(moveHex)

	// 优化 zoomSpeed（最低为16）
	if zoomSpeed > 0 && zoomSpeed < 16 {
		zoomSpeed = 16
	}
	builder.WriteString(fmt.Sprintf("%X0", zoomSpeed>>4)) // 仅取高4位作为一位写入 + “0”

	// 校验码计算
	checkCode := (0xA5 + 0x0F + 0x01 + cmdCode + moveSpeed + moveSpeed + (zoomSpeed & 0xF0)) % 0x100
	builder.WriteString(fmt.Sprintf("%02X", checkCode))

	return builder.String()
}
