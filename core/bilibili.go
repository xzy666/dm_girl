package core

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"goim/libs/bytes"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strconv"
)

const (
	// 订阅所有cmd事件时使用
	BILIBILI_CMD_All string = ""
	// 直播开始
	BILIBILI_CMD_Live string = "LIVE"
	// 直播准备中
	BILIBILI_CMD_Preparing string = "PREPARING"
	// 弹幕消息
	BILIBILI_CMD_DanmuMsg string = "DANMU_MSG"
	// 管理进房
	BILIBILI_CMD_WelcomeGuard string = "WELCOME_GUARD"
	// 群众进房
	BILIBILI_CMD_Welcome string = "WELCOME"
	// 赠送礼物
	BILIBILI_CMD_SendGift string = "SEND_GIFT"
	// 在线人数变动,这不是一个标准cmd类型,仅为了统一handler接口而加入
	BILIBILI_CMD_OnlineChange string = "ONLINE_CHANGE"
)

const (
	BILIBILI_ROOM_INIT_API   = "https://api.live.bilibili.com/room/v1/Room/room_init?id="
	BILIBILI_ROOM_SERVER_API = "https://api.live.bilibili.com/api/player?id=cid:"
	BILIBILI_ROOM_PLAY_API   = "https://api.live.bilibili.com/room/v1/Room/playUrl?cid="
)
const (
	min = 1000000000
	max = 2000000000
)
var roomId *int

type roomInitResult struct {
	Code int `json:"code"`
	Data struct {
		Encrypted   bool `json:"encrypted"`
		HiddenTill  int  `json:"hidden_till"`
		IsHidden    bool `json:"is_hidden"`
		IsLocked    bool `json:"is_locked"`
		LockTill    int  `json:"lock_till"`
		NeedP2p     int  `json:"need_p2p"`
		PwdVerified bool `json:"pwd_verified"`
		RoomID      int  `json:"room_id"`
		ShortID     int  `json:"short_id"`
		UID         int  `json:"uid"`
	} `json:"data"`
	Message string `json:"message"`
	Msg     string `json:"msg"`
}

type dmServerXml struct {
	Server string `xml:"server"`
}

func init() {
	roomId = flag.Int("r", 0, "room id")
}

// 弹幕抓取
func fetch() {
	flag.Parse()
	fmt.Println("输入房间号为：" + strconv.Itoa(*roomId))

	// 获取初始化房间信息
	response, err := http.Get(BILIBILI_ROOM_INIT_API + strconv.Itoa(*roomId))
	if err != nil {
		fmt.Println("获取初始化房间信息失败")
		return
	}
	var res roomInitResult

	jbytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error read init json all")
		return
	}

	if err = json.Unmarshal(jbytes, &res); err != nil {
		fmt.Println("error json Unmarshal")
		return
	}
	if res.Code == 0 {
		fmt.Println("获取到真实的房间号：" + strconv.Itoa(res.Data.RoomID))
	}
	response.Body.Close()

	// 获取弹幕地址
	//dmResponse, err := http.Get(ROOM_SERVER_API + strconv.Itoa(res.Data.RoomID))
	//if err != nil {
	//	fmt.Println("获取弹幕地址失败")
	//	return
	//}
	//xBytes, err := ioutil.ReadAll(dmResponse.Body)
	//fmt.Println(string(xBytes))
	//if err != nil {
	//	fmt.Println("error read dm server all")
	//	return
	//}
	//var xRes dmServerXml
	//if err = xml.Unmarshal(xBytes, &xRes); err != nil {
	//	fmt.Println("error xml Unmarshal")
	//	return
	//}
	//
	//if xRes.Server != "" {
	//	fmt.Println("弹幕服务器的地址：" + xRes.Server)
	//}
	//dmResponse.Body.Close()

	dmServer := "livecmt-1.bilibili.com"
	dmPort := 788

	dstAddr := fmt.Sprintf("%s:%d", dmServer, dmPort)
	dstConn, err := net.Dial("tcp", dstAddr)
	if err != nil {
		fmt.Println("创建弹幕服务器 连接失败")
		return
	}
	fmt.Println("弹幕连接中....")
	uid := rand.Intn(max) + min
	body := fmt.Sprintf("{\"roomid\":%d,\"uid\":%d}", res.Data.RoomID,uid)
	sendSocketData(dstConn,0, 16, 1, 7, 1, body)

	for {
		buf := make([]byte, 4)
		io.ReadAtLeast(dstConn, buf, 4)
		expr := binary.BigEndian.Uint32(buf)
		io.ReadAtLeast(dstConn, buf, 4)
		io.ReadAtLeast(dstConn, buf, 4)
		num := binary.BigEndian.Uint32(buf)
		io.ReadAtLeast(dstConn, buf, 4)

		bLen := int(expr - 16)
		if bLen <= 0 {
			continue
		}
		num = num - 1
		switch num {
		case 3, 4:
			buf = make([]byte, bLen)
			io.ReadAtLeast(dstConn, buf, bLen)
			messages := string(buf)
			fmt.Println(messages)
		}
	}

}

// sendSocketData define
func sendSocketData(dstConn net.Conn,packetlength uint32, magic uint16, ver uint16, action uint32, param uint32, body string) error {
	bodyBytes := []byte(body)
	if packetlength == 0 {
		packetlength = uint32(len(bodyBytes) + 16)
	}
	headerBytes := new(bytes.Buffer)
	var data = []interface{}{
		packetlength,
		magic,
		ver,
		action,
		param,
	}
	for _, v := range data {
		err := binary.Write(headerBytes, binary.BigEndian, v)
		if err != nil {
			return err
		}
	}
	socketData := append(headerBytes.Bytes(), bodyBytes...)
	_, err := dstConn.Write(socketData)
	return err
}
func BlHandler(params ...interface{}){
}
