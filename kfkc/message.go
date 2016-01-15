package kfkc

type Message struct {
	/* orgin message struct */
	Ips				[]string
	Host			string
	Flowid			string
	Taskruntime		string
	Batinterval		string
	File			string
	Loguri			string
	Batnum			string
	Dir				string
	Webfile			string
	Rights			string
	Sid				string
	Localfile_md5	string
	Postscript		string

	/* dealed message struct */
	right			string
	owner			string
}
