package main
import (
	"io/ioutil"
	"fmt"
	"net/http"
	"github.com/robfig/cron"
)

var cacheMap map[string]string;

func handler(writer http.ResponseWriter, request *http.Request){
	languages:= request.URL.Query()["language"];
	language := "";
	if(len(languages)>0){
        	language = languages[0];
        }

	topics:= request.URL.Query()["topic"];
	topic := "";
	if(len(topics)>0){
		topic = topics[0];
	}
        
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
        writer.Header().Set("Access-Control-Allow-Origin", "*")
        
	if value, ok := cacheMap[language+topic]; ok {  
		fmt.Fprintf(writer,string(value));
	} else {
		body := httpGet(language,topic)
		cacheMap[language+topic] = string(body);
		fmt.Fprintf(writer,string(body));
	} 
	
}

func main(){
	cacheMap = make(map[string]string)
	
	cron := cron.New();
	spec :="0 0 0 1/1 * ?"
    cron.AddFunc(spec, func() {
		cacheMap = make(map[string]string);
    })
	cron.Start()
	
	http.HandleFunc("/",handler)
	http.ListenAndServe(":3001",nil)
}

func httpGet(language string,topic string) string{
	url := ""
	if(topic=="" && language==""){
		url = "https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"
	}else{
		url = "https://api.github.com/search/repositories" + "?q=language:" + language;
		if(topic!=""){
			url += "+topic:" + topic;
		}
		url += "&sort=stars&order=desc";
	}
    resp, err :=   http.Get(url)
    if err != nil {
        // handle error
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        // handle error
    }

    return string(body);
}