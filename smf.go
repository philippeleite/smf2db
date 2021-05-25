package main

var parseTable = map[int]func([]byte){
	4: smf004Parse,
}

func smfInit() {
	smf004Init()
}

func smfParse(rec []byte) {
	if f, ok := parseTable[int(rec[1])]; ok {
		f(rec)
	}
}
