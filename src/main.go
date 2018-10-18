package main
import (
	"io/ioutil"
	"fmt"
	"net/http"
)

var cacheMap map[string]string;

func handler(writer http.ResponseWriter, request *http.Request){
	language := "";
	topic := "";

	languages:= request.URL.Query()["language"];
	if(len(languages)>0){
		language = languages[0];
	}
	
	topics:= request.URL.Query()["topic"];
	if(len(topics)>0){
		topic = topics[0];
	}
	

	body := httpGet(language,topic)
	if value, ok := cacheMap[language+topic]; ok {  
		fmt.Fprintf(writer,string(value));
	} else {
		cacheMap[language+topic] = string(body);
		fmt.Fprintf(writer,string(body));
	} 
	
}

func main(){
	cacheMap = make(map[string]string)
	http.HandleFunc("/",handler)
	http.ListenAndServe(":3001",nil)
}

func httpGet(language string,topic string) string{
	url := "https://api.github.com/search/repositories" + "?q=language:" + language + "+topic:" + topic + "&sort=stars&order=desc"
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