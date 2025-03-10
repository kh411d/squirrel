package sqx

import (
	"fmt"
	"log"
	"strings"
)

var debug bool

func Debug(v bool) {
	debug = v
}

func logDebug(sqlStr string, args []interface{}, err error) {
	if debug {
		if err != nil {
			log.Printf(`
==========================================================================================
SQX DEBUG
------------------------------------------------------------------------------------------
ERROR: %s
==========================================================================================
`, err.Error())
			return
		}
		for _, v := range args {
			sqlStr = strings.Replace(sqlStr, "?", fmt.Sprintf("%#v", v), 1)
		}
		log.Printf(`
==========================================================================================
SQX DEBUG
------------------------------------------------------------------------------------------
%s
==========================================================================================
`, sqlStr)
	}
}
