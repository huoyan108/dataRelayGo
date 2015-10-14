package olderhc

import (
	nsq "github.com/bitly/go-nsq"
	"encoding/json"
	"log"
	"strconv"
	//"github.com/golang/protobuf/proto"
    "github.com/huoyan108/dataRelayGo/olderhc/pbgo"
)

//func Receive(cs map[string]chan byte) {
//	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
//		defer func() {
//			if x := recover(); x != nil {
//				log.Println("sorry, chan has closed")
//			}
//		}()
//
//		log.Printf("Got one message: %v", message.Body)
//
//       //SendMsg(topic,msg)
//		res := BuildResult(message.Body)
//        if res.commandType ==int32( Report.Command_CMT_REPWARNUP){
//           // mqttChan <- res.info
//        }
//		chanKey := res.mac + strconv.Itoa(int(res.serialNum))
//		ci, ok := cs[chanKey]
//		if !ok {
//			log.Println("no chan!")
//		} else {
//			ci <- res.info
//			close(ci)
//			delete(cs, chanKey)
//			log.Println("delete existed chan!")
//		}
//		return nil
//	}))
//
//	nsqlookup_url, _ := config.GetString("message", "nsqlookup_url")
//	errConsumer := q.ConnectToNSQD(nsqlookup_url)
//	if errConsumer != nil {
//		log.Panic("Consumer could not connect nsq")
//	}
//}

func Receive(cs map[string]chan []*Report.Command_Param) {
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		defer func() {
			if x := recover(); x != nil {
				log.Println("sorry, chan has closed")
			}
		}()

		log.Printf("Got one message: %v", message.Body)

		res := BuildResult(message.Body)
        if res.commandType ==int32( Report.Command_CMT_REPWARNUP){
            
	        var jsonRes []byte //返回内容
            var paras []*Report.Command_Param
            paras=[]*Report.Command_Param(res.paras)

            mapValue :=make( map[string]string)
            mapValue["warntime"]=strconv.Itoa(int(paras[0].Npara)) 
            mapValue["devtype"]=strconv.Itoa(int(paras[1].Npara)) 
            mapValue["devid"]=strconv.Itoa(int(paras[2].Npara)) 
            mapValue["warnstyle"]=strconv.Itoa(int(paras[3].Npara)) 
            
			jsonRes, _ = json.Marshal(mapValue)
            SendMsg(res.mac,string(jsonRes))
        }else if res.commandType == int32(Report.Command_CMT_REPADDDELDEVICE ){

	        var jsonRes []byte //返回内容
            var paras []*Report.Command_Param
            paras=[]*Report.Command_Param(res.paras)

            mapValue :=make( map[string]string)
            mapValue["action"]=strconv.Itoa(int(paras[0].Npara)) 
            mapValue["devtype"]=strconv.Itoa(int(paras[1].Npara)) 
            mapValue["devid"]=strconv.Itoa(int(paras[2].Npara)) 
            mapValue["manufacturers"]=strconv.Itoa(int(paras[3].Npara)) 
            mapValue["devstate"]=strconv.Itoa(int(paras[4].Npara))            
			jsonRes, _ = json.Marshal(mapValue)
            SendMsg(res.mac,string(jsonRes))
            
        }else{
		chanKey := res.mac + strconv.Itoa(int(res.serialNum))
        log.Println("nsqrecv"+chanKey)
		ci, ok := cs[chanKey]
        log.Println(cs)
		if !ok {
			log.Println("no chan!")
		} else {
			ci <- res.paras
			close(ci)
			delete(cs, chanKey)
			log.Println("delete existed chan!")
		}}
		return nil
	}))

	nsqlookup_url, _ := config.GetString("message", "nsqlookup_url")
	errConsumer := q.ConnectToNSQD(nsqlookup_url)
	if errConsumer != nil {
		log.Panic("Consumer could not connect nsq")
	}
}

func Send(command []byte) {

	command_topic, _ := config.GetString("message", "command_topic")

	err := w.Publish(command_topic, command)
	if err != nil {
		log.Panic("communicator Send could not connect nsq")
	}

	log.Println("send command ok!", command)
}
