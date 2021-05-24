package main

var parseTable = map[int]func([]byte){
	4: smf004Parse,
}

func smfInit() {
	smf004Init()
}

func smfParse(rec []byte) {
	parseTable[int(rec[3])](rec)
}
