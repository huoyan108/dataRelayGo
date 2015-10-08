package olderhc

import (
	"encoding/json"
	"log"
    "fmt"
	"net/http"
	"strconv"
	"time"
	"github.com/golang/protobuf/proto"
    "github.com/huoyan108/dataRelayGo/olderhc/pbgo"
)

//控制指令处理
func ControlHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	now := time.Now()
	var s int32 = int32(now.Unix())
	var jsonRes []byte //返回内容
	//处理panic
	defer func() {
		if x := recover(); x != nil {
			jsonRes, _ = json.Marshal(map[string]byte{"result": SERVER_FAILED})
			fmt.Fprint(w, string(jsonRes))
			log.Println("sorry, server break down!")
		}
	}()

	//获取请求参数mac、cmdtype
	loginid := r.Form["loginid"][0]
    tid ,errNum:=strconv.Atoi(loginid)
//	password, errNum := strconv.Atoi(r.Form["password"][0])
    password := r.Form["password"][0]
	// 请求参数不正确
    log.Println(loginid)
    log.Println(password)
    log.Println(errNum)
    log.Println(len(loginid))
	if errNum != nil || len(loginid) != 15{
		jsonRes, _ = json.Marshal(map[string]byte{"result": BAD_PARAMETER})
		fmt.Fprint(w, string(jsonRes))
		log.Println("control command type error!")
		return
	}
	//构造指令内容

    req := &Report.ControlReport{
        Tid:uint64(tid),
        SerialNumber:uint32( s),
        Command: &Report.Command{
        Type: Report.Command_CMT_REQLOGIN,
        Paras: []*Report.Command_Param{
                    &Report.Command_Param{
                     Type:  Report.Command_Param_STRING,
                     Strpara: password,
                     },
                },
            },
        }       
    reqdata, _ := proto.Marshal(req)

	log.Printf("control command request: ", reqdata)

	//chanKey := loginid + strconv.Itoa(int(s))
	chanKey := loginid + strconv.Itoa(1)
    log.Println("http"+chanKey)
	ci, ok := cs[chanKey]
	if !ok {
        log.Println("make chan")
		ci = make(chan []*Report.Command_Param)
		cs[chanKey] = ci
        log.Println(chanKey)
        log.Println(cs[chanKey])
	}
	Send(reqdata)

	select {
	case res := <-ci:
        value:=[]*Report.Command_Param(res)[0].Npara
        log.Println(value)
		jsonRes, _ = json.Marshal(map[string]uint64{"loginresult": value})

//	case <-time.After(time.Duration(delay) * time.Second):
//		log.Println("res : 超时")
//		close(ci)
//		delete(cs, chanKey)
//		jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
	}

	fmt.Fprint(w, string(jsonRes))

}

func CheckOnlineHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	now := time.Now()
	var s int32 = int32(now.Unix())
	var jsonRes []byte //返回内容
	//处理panic
	defer func() {
		if x := recover(); x != nil {
			jsonRes, _ = json.Marshal(map[string]byte{"result": SERVER_FAILED})
			fmt.Fprint(w, string(jsonRes))
			log.Println("sorry, server break down!")
		}
	}()

	//获取请求参数mac、cmdtype
	loginid := r.Form["tid"][0]
    tid ,errNum:=strconv.Atoi(loginid)
	// 请求参数不正确 
    if errNum != nil || len(loginid) != 15{
        jsonRes, _ = json.Marshal(map[string]byte{"result": BAD_PARAMETER})
        fmt.Fprint(w, string(jsonRes)) 
        log.Println("control command type error!")
        return
	}
	//构造指令内容

    req := &Report.ControlReport{
        Tid:uint64(tid),
        SerialNumber:uint32( s),
        Command: &Report.Command{
        Type: Report.Command_CMT_REQONLINE,
        Paras: []*Report.Command_Param{
                    &Report.Command_Param{
                     Type:  Report.Command_Param_UINT32,
                     Npara: uint64(tid),
                     },
                },
            },
        }       
    reqdata, _ := proto.Marshal(req)

	log.Printf("control command request: ", reqdata)

	chanKey := loginid + strconv.Itoa(1)
    log.Println("http"+chanKey)
	ci, ok := cs[chanKey]
	if !ok {
        log.Println("make chan")
		ci = make(chan []*Report.Command_Param)
		cs[chanKey] = ci
        log.Println(chanKey)
        log.Println(cs[chanKey])
	}
	Send(reqdata)

	select {
	case res := <-ci:
        value:=[]*Report.Command_Param(res)[0].Npara
        log.Println(value)
		jsonRes, _ = json.Marshal(map[string]uint64{"checkRes": value})

//	case <-time.After(time.Duration(delay) * time.Second):
//		log.Println("res : 超时")
//		close(ci)
//		delete(cs, chanKey)
//		jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
	}

	fmt.Fprint(w, string(jsonRes))

}

func GetZBNodeHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	now := time.Now()
	var s int32 = int32(now.Unix())
	var jsonRes []byte //返回内容
	//处理panic
	defer func() {
		if x := recover(); x != nil {
			jsonRes, _ = json.Marshal(map[string]byte{"result": SERVER_FAILED})
			fmt.Fprint(w, string(jsonRes))
			log.Println("sorry, server break down!")
		}
	}()

	//获取请求参数mac、cmdtype
	loginid := r.Form["tid"][0]
    tid ,errNum:=strconv.Atoi(loginid)
	// 请求参数不正确
	if errNum != nil || len(loginid) != 15{
		jsonRes, _ = json.Marshal(map[string]byte{"result": BAD_PARAMETER})
		fmt.Fprint(w, string(jsonRes))
		log.Println("control command type error!")
		return
	}
	//构造指令内容

    req := &Report.ControlReport{
        Tid:uint64(tid),
        SerialNumber:uint32( s),
        Command: &Report.Command{
        Type: Report.Command_CMT_REQDEVICELIST,
        Paras: []*Report.Command_Param{
                    &Report.Command_Param{
                     Type:  Report.Command_Param_UINT32,
                     Npara: uint64(tid),
                     },
                },
            },
        }       
    reqdata, _ := proto.Marshal(req)

	log.Printf("control command request: ", reqdata)

	chanKey := loginid + strconv.Itoa(1)
    log.Println("http"+chanKey)
	ci, ok := cs[chanKey]
	if !ok {
        log.Println("make chan")
		ci = make(chan []*Report.Command_Param)
		cs[chanKey] = ci
        log.Println(chanKey)
        log.Println(cs[chanKey])
	}
	Send(reqdata)

    var jsonSource string
	select {
	case res := <-ci:
        var paras []*Report.Command_Param
        paras = res
        var parasLen int
        parasLen = len(paras)
        
        jsonSource = "totalcount:" + strconv.Itoa(parasLen) +","

        jsonSource += "response_params:["
        dataCount := parasLen/3
        for i:=0; i<dataCount;i++{
            jsonSource += "{"
            temp := "devid:"+strconv.Itoa(int(paras[i*3].Npara))
            jsonSource += temp

            temp = "devtype:"+strconv.Itoa(int(paras[i*3+1].Npara))
            jsonSource += temp

            temp = "manufacturers:"+strconv.Itoa(int(paras[i*3+2].Npara))
            jsonSource += temp


            jsonSource += "},"
        }
        jsonSource += "]"
        log.Println("jsonSource"+jsonSource)
	    jsonRes, _ = json.Marshal(jsonSource)
        
//	case <-time.After(time.Duration(delay) * time.Second):
//		log.Println("res : 超时")
//		close(ci)
//		delete(cs, chanKey)
//		jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
	}

	fmt.Fprint(w, string(jsonRes))

}
func MainsOperationHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	now := time.Now()
	var s int32 = int32(now.Unix())
	var jsonRes []byte //返回内容
	//处理panic
	defer func() {
		if x := recover(); x != nil {
			jsonRes, _ = json.Marshal(map[string]byte{"result": SERVER_FAILED})
			fmt.Fprint(w, string(jsonRes))
			log.Println("sorry, server break down!")
		}
	}()

	//获取请求参数mac、cmdtype
	loginid := r.Form["tid"][0]
    tid ,errNum:=strconv.Atoi(loginid)
	// 请求参数不正确 
    if errNum != nil || len(loginid) != 15{
        jsonRes, _ = json.Marshal(map[string]byte{"result": BAD_PARAMETER})
        fmt.Fprint(w, string(jsonRes)) 
        log.Println("control command type error!")
        return
	}
    ep := r.Form["ep"][0]
    nEp,errNum := strconv.Atoi(ep)
    operatortype := r.Form["operatortype"][0]
    nOperatortype,errNum := strconv.Atoi(operatortype)
	//构造指令内容

    req := &Report.ControlReport{
        Tid:uint64(tid),
        SerialNumber:uint32( s),
        Command: &Report.Command{
        Type: Report.Command_CMT_REQOP,
        Paras: []*Report.Command_Param{
                    &Report.Command_Param{
                     Type:  Report.Command_Param_UINT64,
                     Npara: uint64(nEp),
                     },
                    &Report.Command_Param{
                     Type:  Report.Command_Param_UINT8,
                     Npara: uint64(nOperatortype),
                     },
                },
            },
        }       
    reqdata, _ := proto.Marshal(req)

	log.Printf("control command request: ", reqdata)

	chanKey := loginid + strconv.Itoa(1)
    log.Println("http"+chanKey)
	ci, ok := cs[chanKey]
	if !ok {
        log.Println("make chan")
		ci = make(chan []*Report.Command_Param)
		cs[chanKey] = ci
        log.Println(chanKey)
        log.Println(cs[chanKey])
	}
	Send(reqdata)

	select {
	case res := <-ci:
        value:=[]*Report.Command_Param(res)[0].Npara
        log.Println(value)
		jsonRes, _ = json.Marshal(map[string]uint64{"controlRes": value})

//	case <-time.After(time.Duration(delay) * time.Second):
//		log.Println("res : 超时")
//		close(ci)
//		delete(cs, chanKey)
//		jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
	}

	fmt.Fprint(w, string(jsonRes))

}
func ChangeDeviceNameHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	now := time.Now()
	var s int32 = int32(now.Unix())
	var jsonRes []byte //返回内容
	//处理panic
	defer func() {
		if x := recover(); x != nil {
			jsonRes, _ = json.Marshal(map[string]byte{"result": SERVER_FAILED})
			fmt.Fprint(w, string(jsonRes))
			log.Println("sorry, server break down!")
		}
	}()

	//获取请求参数mac、cmdtype
	loginid := r.Form["tid"][0]
    tid ,errNum:=strconv.Atoi(loginid)
	// 请求参数不正确 
    if errNum != nil || len(loginid) != 15{
        jsonRes, _ = json.Marshal(map[string]byte{"result": BAD_PARAMETER})
        fmt.Fprint(w, string(jsonRes)) 
        log.Println("control command type error!")
        return
	}
    ep := r.Form["ep"][0]
    nep,_ := strconv.Atoi(ep)
    newname := r.Form["newname"][0]
	//构造指令内容

    req := &Report.ControlReport{
        Tid:uint64(tid),
        SerialNumber:uint32( s),
        Command: &Report.Command{
        Type: Report.Command_CMT_REQOP,
        Paras: []*Report.Command_Param{
                    &Report.Command_Param{
                     Type:  Report.Command_Param_UINT64,
                     Npara: uint64(nep),
                     },
                    &Report.Command_Param{
                     Type:  Report.Command_Param_STRING,
                     Strpara: newname,
                     },
                },
            },
        }       
    reqdata, _ := proto.Marshal(req)

	log.Printf("control command request: ", reqdata)

	chanKey := loginid + strconv.Itoa(1)
    log.Println("http"+chanKey)
	ci, ok := cs[chanKey]
	if !ok {
        log.Println("make chan")
		ci = make(chan []*Report.Command_Param)
		cs[chanKey] = ci
        log.Println(chanKey)
        log.Println(cs[chanKey])
	}
	Send(reqdata)

	select {
	case res := <-ci:
        value:=[]*Report.Command_Param(res)[0].Npara
        log.Println(value)
		jsonRes, _ = json.Marshal(map[string]uint64{"changePwdRes": value})

//	case <-time.After(time.Duration(delay) * time.Second):
//		log.Println("res : 超时")
//		close(ci)
//		delete(cs, chanKey)
//		jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
	}

	fmt.Fprint(w, string(jsonRes))

}
func ChangeUserPasswordHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	now := time.Now()
	var s int32 = int32(now.Unix())
	var jsonRes []byte //返回内容
	//处理panic
	defer func() {
		if x := recover(); x != nil {
			jsonRes, _ = json.Marshal(map[string]byte{"result": SERVER_FAILED})
			fmt.Fprint(w, string(jsonRes))
			log.Println("sorry, server break down!")
		}
	}()

	//获取请求参数mac、cmdtype
	loginid := r.Form["tid"][0]
    tid ,errNum:=strconv.Atoi(loginid)
	// 请求参数不正确 
    if errNum != nil || len(loginid) != 15{
        jsonRes, _ = json.Marshal(map[string]byte{"result": BAD_PARAMETER})
        fmt.Fprint(w, string(jsonRes)) 
        log.Println("control command type error!")
        return
	}
    oldpwd := r.Form["oldpwd"][0]
    newpwd := r.Form["newpwd"][0]
	//构造指令内容

    req := &Report.ControlReport{
        Tid:uint64(tid),
        SerialNumber:uint32( s),
        Command: &Report.Command{
        Type: Report.Command_CMT_REQOP,
        Paras: []*Report.Command_Param{
                    &Report.Command_Param{
                     Type:  Report.Command_Param_STRING,
                     Strpara: oldpwd,
                     },
                    &Report.Command_Param{
                     Type:  Report.Command_Param_STRING,
                     Strpara: newpwd,
                     },
                },
            },
        }       
    reqdata, _ := proto.Marshal(req)

	log.Printf("control command request: ", reqdata)

	chanKey := loginid + strconv.Itoa(1)
    log.Println("http"+chanKey)
	ci, ok := cs[chanKey]
	if !ok {
        log.Println("make chan")
		ci = make(chan []*Report.Command_Param)
		cs[chanKey] = ci
        log.Println(chanKey)
        log.Println(cs[chanKey])
	}
	Send(reqdata)

	select {
	case res := <-ci:
        value:=[]*Report.Command_Param(res)[0].Npara
        log.Println(value)
		jsonRes, _ = json.Marshal(map[string]uint64{"changePwdRes": value})

//	case <-time.After(time.Duration(delay) * time.Second):
//		log.Println("res : 超时")
//		close(ci)
//		delete(cs, chanKey)
//		jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
	}

	fmt.Fprint(w, string(jsonRes))

}

