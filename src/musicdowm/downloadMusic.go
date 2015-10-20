// downloadMusic
package musicdowm

import (
	"bytes"
	"fmt"
	"getGuid"
	"gopkg.in/mgo.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	url_mg  = "192.168.100.168:27017"
	mongodb = "music"
)

//Database
func getMongodb() *mgo.Database {
	//	var database mgo.Database
	session, err := mgo.Dial(url_mg)
	session.SetMode(mgo.Monotonic, true)
	session.SetBatch(10)
	if err != nil {
		fmt.Println("mongodb seecion Err!", err)
		getMongodb()
	}
	database := session.DB(mongodb)
	defer session.Close()
	return database

}

func Download(musicInfo map[string]interface{}) interface{} {
	fmt.Println("通过开启协程处理下载工作")
	var url string
	if nil != musicInfo {
		db := getMongodb()
		files, err := db.GridFS("fs").Create(string(getGuid.NewObjectId().Hex()))
		if err != nil {
			fmt.Println("create GridFs err")
		} else {
			url = musicInfo["musictopurl"].(string)
			name := musicInfo["name"]
			format := musicInfo["format"]
			fileName := name.(string) + "." + format.(string)
			out, _ := os.Create(fileName)
			defer out.Close()
			resp, _ := http.Get(url)
			defer resp.Body.Close()
			pix, _ := ioutil.ReadAll(resp.Body)
			n, err := io.Copy(out, bytes.NewReader(pix))
			//			messages, err := os.Open(fileName)
			//			defer messages.Close()
			//			io.Copy(files, messages)
			n1, err := files.Write(pix)
			fmt.Println(n, n1)
			err1 := files.Close()
			if err1 != nil {
				fmt.Println("mongodb Gridfs Close Err")
			}
			if err != nil {
				fmt.Println(n, "下载完成！")
			}
		}

	}

	return nil
}
