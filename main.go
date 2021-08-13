package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type config struct {
	Exe  string `json:"exe"`
	Args string `json:"args"`
}

type tSaveData struct {
	File   string  `json:"file"`
	Proofs float64 `json:"proofs"`
}

const (
	configFile = "config.json"
	saveFile   = "plotProofs.json"
	startText  = "Testing plot"
	endText    = "Proofs"
)

//var decoder *encoding.Decoder
//greenBg
//whiteBg      =
//yellowBg     =
//redBg        =
//blueBg       =
//magentaBg    =
//cyanBg       =
//green        =
//white        =
//yellow       =
//red          =
//blue         =
//magenta      =
//cyan         =
//reset        =
var colorBg = [15]string{
	string([]byte{27, 91, 57, 55, 59, 52, 50, 109}),
	string([]byte{27, 91, 57, 48, 59, 52, 55, 109}),
	string([]byte{27, 91, 57, 48, 59, 52, 51, 109}),
	string([]byte{27, 91, 57, 55, 59, 52, 49, 109}),
	string([]byte{27, 91, 57, 55, 59, 52, 52, 109}),
	string([]byte{27, 91, 57, 55, 59, 52, 53, 109}),
	string([]byte{27, 91, 57, 55, 59, 52, 54, 109}),
	string([]byte{27, 91, 51, 50, 109}),
	string([]byte{27, 91, 51, 55, 109}),
	string([]byte{27, 91, 51, 51, 109}),
	string([]byte{27, 91, 51, 49, 109}),
	string([]byte{27, 91, 51, 52, 109}),
	string([]byte{27, 91, 51, 53, 109}),
	string([]byte{27, 91, 51, 54, 109}),
	string([]byte{27, 91, 48, 109})}

func getConfig() (con config, ok bool) {
	data, err := ioutil.ReadFile("./" + configFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &con)
	if err != nil {
		panic(err)
	}
	return con, true
}

func main() {

	con, ok := getConfig()
	if !ok {
		fmt.Println("配置文件有错")
		return
	}

	cmd := exec.Command(con.Exe, strings.Split(con.Args, " ")...)
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stderr)

	b := false
	var saveData []tSaveData
	for scanner.Scan() {
		bs := scanner.Text()
		for _, s := range colorBg {
			bs = strings.Replace(bs, s, "", -1)
		}
		m := bs
		position := strings.Index(m, startText)
		if position > 1 {
			contents := strings.Split(m, startText)
			content := contents[len(contents)-1]
			save := tSaveData{
				File: content[1 : len(content)-5],
			}
			saveData = append(saveData, save)
			b = true
		}
		if b {
			position = strings.Index(m, endText)
			if position > 1 {
				contents := strings.Split(m, endText)
				content := contents[len(contents)-1]
				index := len(saveData) - 1
				saveData[index].Proofs, _ = strconv.ParseFloat(content[10:], 32)
				fmt.Printf("第%d个文件：%s，幸运值%.4f\n", index+1, saveData[index].File, saveData[index].Proofs)
				b = false

			}
		}

	}
	if len(saveData) > 0 {
		sort.Slice(saveData, func(i, j int) bool {
			return saveData[i].Proofs < saveData[j].Proofs
		})
		data, err := json.Marshal(saveData)
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile("./"+saveFile, data, 0777)
	}
	cmd.Wait()
}

//func init() {
//	decoder = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
//}
