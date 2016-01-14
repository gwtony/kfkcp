package kfkc

type Message struct {
	/* orgin message struct */
	Ips				[]string
	Host			string
	Flowid			int64
	Taskruntime		string
	Batinterval		string
	File			string
	Loguri			string
	Batnum			string
	Dir				string
	Webfile			string
	Rights			string

	/* dealed message struct */
	right			string
	owner			string
}
