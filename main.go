package main

import (
	"fmt"
	"go_second/app"

	//"net/http"

	"github.com/spf13/viper"
	//. "go_second/app"
)

func main() {
	viper.AddConfigPath("./env")
	viper.SetConfigName("db")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	var err error
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	var config app.Config
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Cant config by viper")
		//http.StatusBadRequest()
		return
	}
	db, handler := app.NewApp(config)
	fmt.Println(handler)
	//fmt.Println(db)
	var r app.Router
	r.Router(db, handler)

}
