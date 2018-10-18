package main
import (
	"io/ioutil"
	"log"
	"fmt"
	"net/http"
)

var cacheMap map[string]string;

func handler(writer http.ResponseWriter, request *http.Request){
	languages:= request.URL.Query()["language"];
	language := languages[0];
	log.Println("Url Param 'language' is: " + string(language));

	topics:= request.URL.Query()["topic"];
	topic := topics[0];
	log.Println("Url Param 'topic' is: " + string(topic));

	body := httpGet(language,topic)
	if value, ok := cacheMap[language+topic]; ok {  
		log.Println(999);
		fmt.Fprintf(writer,string(value));
	} else {
		log.Println(888);
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