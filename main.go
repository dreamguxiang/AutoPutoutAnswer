package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	putoutTitleList()
	//ReadTitle("timu.txt")
}

func ReadTitle(filepath string) {
	data, err := ioutil.ReadFile("title.json")
	m := make(map[string]string)

	json.Unmarshal(data, &m)

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	log.Println("> Reading timu.txt")
	scanner := bufio.NewScanner(file)
	num := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "*") {
			//out := regexp.MustCompile("[0-9]+([^)]*)([^(]*)[*]").FindStringSubmatch(line)
			reg := regexp.MustCompile("[^\u4e00-\u9fa5]")
			line = reg.ReplaceAllString(line, "")
			if _, isOk := m[line]; isOk {
				num += 1
				fmt.Println(strconv.Itoa(num) + "." + line + "\n答案：" + m[line] + "\n")
			} else {
				num += 1
				fmt.Println(strconv.Itoa(num) + "." + line + "\n答案：不存在\n")
			}
		}
	}
}

func createTitleList(filepath string) map[string]string {
	map1 := make(map[string]string)
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Println("> Reading ALL titlelist.txt")
	str := string(data)
	out := regexp.MustCompile("([*])[^*]+(?:;[^*][^*]*)*").FindAllString(str, -1)
	for _, v := range out {
		if strings.Contains(v, "正确答案") {
			outs1 := regexp.MustCompile("：[A-D]+").FindStringSubmatch(v)
			if len(outs1) == 0 {
				outs1 = regexp.MustCompile("：[\u4e00-\u9fa5]+").FindStringSubmatch(v)
			}
			outs2 := regexp.MustCompile("[A-D]+").FindStringSubmatch(outs1[0])
			if len(outs2) == 0 {
				outs2 = regexp.MustCompile("[\u4e00-\u9fa5]+").FindStringSubmatch(outs1[0])
			}
			right := outs2[0]
			title := ""
			scanner := bufio.NewScanner(strings.NewReader(v))

			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, "*") {
					reg := regexp.MustCompile("[^\u4e00-\u9fa5]")
					line = reg.ReplaceAllString(line, "")
					title = line
				}
				if strings.Contains(line, right) {
					right = line
				}
			}
			map1[title] = clearABCD(right)
		}
	}
	return map1
}

func clearABCD(right string) string {
	right = replaceStr(right, "回答错误，正确答案为：", "")
	right = replaceStr(right, "A.", "")
	right = replaceStr(right, "B.", "")
	right = replaceStr(right, "C.", "")
	right = replaceStr(right, "D.", "")
	right = replaceStr(right, "D．", "") //fix bug
	return right
}

func putoutTitleList() {
	maps := make(map[string]string)
	map1 := createTitleList("titlelist.txt")
	for k, v := range map1 {
		maps[k] = v
	}
	mjson, _ := json.Marshal(maps)
	mString := string(mjson)
	jsonFile, err := os.Create("title.json") // 创建 json 文件
	if err != nil {
		log.Printf("create json file %v error [ %v ]", "title.json", err)
		return
	}
	defer jsonFile.Close()

	encode := json.NewEncoder(jsonFile) // 创建编码器
	err = encode.Encode(mString)
}

func replaceStr(str1 string, str2 string, str3 string) string {
	return strings.Replace(str1, str2, str3, -1)
}
