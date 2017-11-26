/*
Copyright (c) 2017 Jean-François PHILIPPE
Package goconfig implements a config object initialised from files and env variables.


Usage

     builder := NewBuilder("CTX_", nil)
	 // Main config values
	 _, err := builder.LoadJsonFile("/etc/myapp/config.json")

	 // complete config with values from another file
	 // existing values will NOT be updated.
	 _,err = builder.LoadTxtFile("/etc/default/mayapp.txt")

	 // OR
	 // _,err := builder.LoadFiles("/etc/myapp/config.json","/etc/default/mayapp.txt")
	 // Use it
	 config := builder.Config()

	 name := config.GetString("name", "default name")
	 num_proc := config.GetInt("proc.limit", 5)

	 dbconfig := config.GetConfig("database")

	 db_url := dbconfig.GetString("url", "")
	 db_port := dbconfig.GetInt("port", 1234)

	 // Or
	 // db_url := config.GetString("database.url","")
	 ...

Config store values as nested map[string]interface{}.

Txt File Format

	# comment.
	// other type of comment
	# comments must be alone on a line
	name = app name
	database.url = user:${db.passwd}@/dbname
	database.port = 3456

*/

package goconfig

// vi:set fileencoding=utf-8 tabstop=4 ai
