package main

import (
	"github.com/yungsem/db-desc/cnf"
	_ "github.com/yungsem/db-desc/cnf"
	"github.com/yungsem/db-desc/db"
	"log"
	"os"
	"text/template"
)

func main() {

	conf, err := cnf.NewConf()
	if err != nil {
		log.Fatal(err.Error())
	}

	describer, err := db.NewTableDescriber(conf.DB.Type,
		conf.DB.Host,
		conf.DB.Port,
		conf.DB.Username,
		conf.DB.Password,
		conf.DB.Schema)
	if err != nil {
		log.Fatal(err.Error())
	}

	tableInfos, err := describer.DescribeTable()
	if err != nil {
		log.Fatal(err.Error())
	}

	tmpl, err := template.ParseFiles("./templates/ER.sql")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("out.sql")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(f, map[string]any{
		"tableInfos": tableInfos,
	})

	if err != nil {
		panic(err)
	}
}
