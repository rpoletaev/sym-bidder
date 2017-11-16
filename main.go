package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/rpoletaev/sym-bidder/api"
)

//Config describes configuration file

const exampleData = `{"id":"CA1DF4146DE67248","imp":[{"id":"1","banner":{"w":320,"h":480,"pos":1,"btype":[4],"battr":[3,8,10,14],"api":[3,5]},"video":{"mimes":["video/3gpp","video/mp4"],"minduration":16,"maxduration":120,"protocols":[2,5,3,6],"w":320,"h":480,"linearity":1,"sequence":1,"battr":[3,8,10,14],"api":[3,5],"companiontype":[1,2,3]},"displaymanager":"mopub","displaymanagerver":"4.7.1","instl":1,"tagid":"01a16f7530db4b6689d10d7d5cab3183","bidfloor":25.22,"ext":{"brsrclk":1,"dlp":1}}],"app":{"bundle":"com.erenapps.beachatvsim","ver":"29","id":"7271c98bfac340ff83c0bf3008f31d55","name":"Beach ATV Simulator - 16478","cat":["IAB1","IAB9","IAB9-30","entertainment","games"],"publisher":{"id":"57e33763e5804d75802705596150f3c1","name":"Appodeal, Inc."}},"device":{"ua":"Mozilla/5.0 (Linux; Android 4.4.4; SM-T561 Build/KTU84P) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/33.0.0.0 Safari/537.36","ip":"178.206.218.71","geo":{"lat":54.55,"lon":52.8,"accuracy":50,"ipservice":3,"country":"RUS","region":"73","city":"Bugulma","zip":"423230","type":2,"utcoffset":180,"ext":{"old_geo":{"country":"RUS","region":"73","city":"Bugulma","zip":"423230"}}},"carrier":"250-02","language":"ru","make":"samsung","model":"SM-T561","os":"Android","osv":"4.4.4","js":1,"connectiontype":2,"ifa":"951bfcee-3fc4-45c3-82b2-50371681e28a","h":800,"w":1280,"pxratio":1},"at":2,"cur":["USD"],"bcat":["IAB25","IAB26","IAB9-9"],"ext":{"envisionx":{"ssp":9}}}`

func main() {
	var configPath string

	logger := log.New(os.Stdout, "bidder", log.LstdFlags)
	flag.StringVar(&configPath, "c", "config.json", "-c /path/to/config/file")
	flag.Parse()

	bts, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Fatalln("не удалость прочесть файл конфига: ", err)
	}

	var config api.Config
	if err := json.Unmarshal(bts, &config); err != nil {
		logger.Fatalln("ошибка при разборе конфига: ", err)
	}

	api := api.CreateApi(&config, logger)
	api.Run()
}
