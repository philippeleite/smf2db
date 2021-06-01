package main

var parseTable = map[int]func([]byte){
	4:  smf004Parse,
	5:  smf005Parse,
	14: smf014Parse,
	15: smf015Parse,
	17: smf017Parse,
}

func smfInit() {
	smf004Init()
	smf005Init()
	smf014Init()
	smf015Init()
	smf017Init()
}

func smfParse(rec []byte) {
	if f, ok := parseTable[int(rec[1])]; ok {
		f(rec)
	}
}
