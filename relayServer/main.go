package  main

import(
    "fmt"
     "github.com/huoyan108/dataRelayGo/olderhc" 
    //"os"
    //"io"
    "log"
    "net/http"
)


func main(){

    fmt.Println("Server Start")
	http.HandleFunc("/netgate/network/login", olderhc.ControlHandler)
	http.HandleFunc("/netgate/network/checkOnline", olderhc.CheckOnlineHandler)
	http.HandleFunc("/netgate/network/getzbnodes", olderhc.GetZBNodeHandler)
    http.HandleFunc("/netgate/network/mainsoperation",olderhc.MainsOperationHandler) 
    http.HandleFunc("/netgate/network/changedevicename",olderhc.ChangeDeviceNameHandler)
    http.HandleFunc("/netgate/network/changeuserpassword",olderhc.ChangeUserPasswordHandler) 
	err := http.ListenAndServe(":"+olderhc.Port, nil)
	
	if err != nil {
        log.Fatal("ListenAndServer:",err)
    }
    

}
