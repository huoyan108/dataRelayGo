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
func char2byte(c string) byte {
    switch c {
        case "0":
            return 0
            case "1":
            return 1
            case "2":
            return 2
            case "3":
            return 3
            case "4":
            return 4
            case "5":
            return 5
            case "6":
            return 6
            case "7":
            return 7
            case "8":
            return 8
            case "9":
            return 9
            case "a":
            return 10
            case "b":
            return 11
            case "c":
            return 12
            case "d":
            return 13
            case "e":
            return 14
            case "f":
            return 15
    }
    return 0
}

func Macaddr2uint64(mac string) uint64 {
    var buffer []byte
        buffer = append(buffer, 0)
        buffer = append(buffer, 0)
        value := char2byte(string(mac[0]))*16 + char2byte(string(mac[1]))
        buffer = append(buffer, value)
        value = char2byte(string(mac[2]))*16 + char2byte(string(mac[3]))
        buffer = append(buffer, value)
        value = char2byte(string(mac[4]))*16 + char2byte(string(mac[5]))
        buffer = append(buffer, value)
        value = char2byte(string(mac[6]))*16 + char2byte(string(mac[7]))
        buffer = append(buffer, value)
        value = char2byte(string(mac[8]))*16 + char2byte(string(mac[9]))
        buffer = append(buffer, value)
        value = char2byte(string(mac[10]))*16 + char2byte(string(mac[11]))
        buffer = append(buffer, value)

        return binary.BigEndian.Uint64(buffer)
}






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

         chanKey := loginid + strconv.Itoa(int(s))
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

             case <-time.After(time.Duration(delay) * time.Second):
          log.Println("res : 超时")
          close(ci)
          delete(cs, chanKey)
          jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
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
        SerialNumber:uint32(s),
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

         chanKey := loginid + strconv.Itoa(int(s))
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

             case <-time.After(time.Duration(delay) * time.Second):
          log.Println("res : 超时")
          close(ci)
          delete(cs, chanKey)
          jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
         }

     fmt.Fprint(w, string(jsonRes))

}
type DeviceInfo struct{
    Devid uint64
        Devtype int
        Manufacturers int
        Devstate int
        Devname  string
}
type ZBNodeGroup struct{
    TidName string
        TotalCount int
        DeviceGroup []DeviceInfo
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
        SerialNumber:uint32(s),
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

         chanKey := loginid + strconv.Itoa(int(s))
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
                       var paras []*Report.Command_Param
                           paras = res

                           var group ZBNodeGroup
                           group.TidName  = paras[0].Strpara
                           group.TotalCount = int(paras[1].Npara)
                           for i:=0 ; i<group.TotalCount;i++{

                               var deviceInfo DeviceInfo
                                   deviceInfo.Devid = paras[i*5+2].Npara
                                   deviceInfo.Devtype = int(paras[i*5+1+2].Npara)
                                   deviceInfo.Manufacturers = int(paras[i*5+2+2].Npara)
                                   deviceInfo.Devstate = int(paras[i*5+3+2].Npara)
                                   deviceInfo.Devname = paras[i*5+4+2].Strpara

                                   group.DeviceGroup = append(group.DeviceGroup,deviceInfo)
                           }
                       log.Println(group)
                           jsonRes, _ = json.Marshal(group)

             case <-time.After(time.Duration(delay) * time.Second):
                           log.Println("res : 超时")
                           close(ci)
                           delete(cs, chanKey)
                           jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
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
ep := r.Form["devid"][0]
        nEp,errNum := strconv.Atoi(ep)
        operatortype := r.Form["operatortype"][0]
        nOperatortype,errNum := strconv.Atoi(operatortype)
        ser := r.Form["ser"][0]
        nSer,errNum := strconv.Atoi(ser)
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
       Npara: uint64(nSer),
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

        chanKey := loginid + strconv.Itoa(int(s))
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

          m := make(map[string]uint64)
          m["controlRes"] = value
          m["ser"] = uint64(nSer)
          jsonRes,_ = json.Marshal(m)
          //jsonRes, _ = json.Marshal(map[string]uint64{"controlRes": value})

            case <-time.After(time.Duration(delay) * time.Second):
          log.Println("res : 超时")
          close(ci)
          delete(cs, chanKey)
          jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
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
Type: Report.Command_CMT_REQSETDEVICENAME ,
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

        chanKey := loginid + strconv.Itoa(int(s))
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
          jsonRes, _ = json.Marshal(map[string]uint64{"changeNameRes": value})

            case <-time.After(time.Duration(delay) * time.Second):
          log.Println("res : 超时")
          close(ci)
          delete(cs, chanKey)
          jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
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
Type: Report.Command_CMT_REQCHANGEPASSWD,
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

            chanKey := loginid + strconv.Itoa(int(s))
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

                case <-time.After(time.Duration(delay) * time.Second):
          log.Println("res : 超时")
          close(ci)
          delete(cs, chanKey)
          jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
            }

        fmt.Fprint(w, string(jsonRes))

}
func DelDeviceHandler(w http.ResponseWriter, r *http.Request) {

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
devid := r.Form["devid"][0]
           nDevid,errNum := strconv.Atoi(devid)

           //构造指令内容
           req := &Report.ControlReport{
Tid:uint64(tid),
        SerialNumber:uint32(s),
        Command: &Report.Command{
Type: Report.Command_CMT_REQONLINE,
      Paras: []*Report.Command_Param{
          &Report.Command_Param{
Type:  Report.Command_Param_UINT32,
       Npara: uint64(nDevid),
          },
      },
        },
           }       
       reqdata, _ := proto.Marshal(req)

           log.Printf("control command request: ", reqdata)

           chanKey := loginid + strconv.Itoa(int(s))
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
          jsonRes, _ = json.Marshal(map[string]uint64{"delRes": value})

               case <-time.After(time.Duration(delay) * time.Second):
          log.Println("res : 超时")
          close(ci)
          delete(cs, chanKey)
          jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
           }

       fmt.Fprint(w, string(jsonRes))

}

