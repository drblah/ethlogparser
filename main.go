package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/drblah/ethlogparser/parser"
)

func main() {
	var str = `DEBUG[10-11|08:21:00.588] Inserting downloaded chain               items=1 firstnum=1 firsthash=9c2008…0bbea8 lastnum=1 lasthash=9c2008…0bbea8
DEBUG[10-11|08:21:00.588] Trie cache stats after commit            misses=0 unloads=0
DEBUG[10-11|08:21:00.588] Chain split detected                     number=0 hash=351c48…6c9ea9 drop=1 dropfrom=ba0794…5bfc7a add=1 addfrom=9c2008…0bbea8
DEBUG[10-11|08:21:00.588] Inserted new block                       number=1 hash=9c2008…0bbea8 uncles=0 txs=0 gas=0 elapsed=400.009µs
DEBUG[10-11|08:21:00.950] Synchronisation terminated               elapsed=1.192096689s
INFO [10-11|08:21:00.950] Fast sync complete, auto disabling 
TRACE[10-11|08:21:00.950] Announced block                          hash=04ee97…2d7ed7 recipients=1 duration=2562047h47m16.854s`

	scanner := bufio.NewScanner(strings.NewReader(str))

	for scanner.Scan() {

		tmpStr := parser.SplitByCol(scanner.Text())

		for _, s := range tmpStr {
			fmt.Println(s)
		}

		status, timestamp := parser.ParseStatusTimestamp(tmpStr[0])

		fmt.Println(status, timestamp)
		os.Exit(0)
	}
}
