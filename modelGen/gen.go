package modelGen

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//OCGenReturn 生成iOS returns 的模型
func OCGenReturn(savePath string, modelName string) {
	ocGen(savePath, modelName, false, true)
}

//OCGenForm 生成iOS forms 的模型
func OCGenForm(savePath string, modelName string) {
	ocGen(savePath, modelName, true, false)
}

//OCGen 生成iOS 普通模型
func OCGen(savePath string, modelName string) {
	ocGen(savePath, modelName, false, false)
}

//OCAutoGen savePath 保存路径(最后不带斜线)，如 /Users/cooerson/Documents/ios/action_ios/Action/Action
//modelName 文件名，生成文件结尾自动加上Model.go
func ocGen(savePath string, modelName string, isForm bool, isReturn bool) {
	var strFinal string
	var strStruct string
	var strFinal2 string //form or return的 model 的 .m 文件的内容
	strFinal2 += "#import \"" + modelName + "Model.h\"\n\n"
	readLine(modelName, func(s string) {
		//按空格分割此行字符串
		arrs := strings.Fields(s)
		if strings.HasPrefix(s, "@") {
			return
		} else if strings.HasPrefix(s, "=") || strings.HasPrefix(s, "~=") {
			//匹配`=xx`或`=Xx`
			//指定为interface的名称Xx
			if strings.HasPrefix(s, "=") {
				s = strings.Replace(s, "=", "", 1)
			} else if strings.HasPrefix(s, "~=") {
				s = strings.Replace(s, "~=", "", 1)
			}
			strStruct = strings.Title(s) //首字母大写
			// @interface User : NSObject
			strFinal += "@interface " + strStruct + "Model : NSObject\n"
			//@implementation PicMeta @end
			strFinal2 += "@implementation " + strStruct + "Model : NSObject\n@end \n\n"
		} else if strings.HasPrefix(s, "-ID") || strings.HasPrefix(s, "-id") || strings.HasPrefix(s, "-Id") {
			//匹配`-ID`或`-id`或`-Id`
			//指定字段名为 `xxID`
			//指定类型为 `NSString`
			//例如 @property (nonatomic, copy) NSString *userID;
			strStruct := lowerFisterLetter(strStruct)
			strFinal += "@property (nonatomic, strong) NSString *" + strStruct + "ID;\n"
		} else if strings.HasPrefix(s, "-") {
			//1.匹配`-xx`或`-Xx`
			//  指定字段名为`Xx，首字母小写
			strField := strings.Replace(arrs[0], "-", "", 1)
			strField = lowerFisterLetter(strField)
			//2.匹配类型：例如0、1、2、3、4、5、6、自定义字符串
			//  0:指定字段类型为`bson.ObjectId`->(nonatomic, strong) NSString
			//  1:指定字段类型为`string`->(nonatomic, strong) NSString
			//  2:指定字段类型为`bool`->(nonatomic) BOOL
			//  3:指定字段类型为`float64`->(nonatomic) double
			//  4:指定字段类型为`int64`->(nonatomic) NSInteger
			//  5:指定字段类型为`int8`->(nonatomic) NSInteger
			//  6:指定字段类型为`time.Time`->(nonatomic, strong) NSString
			var strType string
			sIntValue, err := strconv.Atoi(arrs[1])
			if err == nil {
				if sIntValue >= 0 || sIntValue <= 6 {
					if sIntValue == 0 {
						strType = "(nonatomic, strong) NSString *"
					} else if sIntValue == 1 {
						strType = "(nonatomic, strong) NSString *"
					} else if sIntValue == 2 {
						strType = "(nonatomic, assign) BOOL "
					} else if sIntValue == 3 {
						strType = "(nonatomic, assign) double "
					} else if sIntValue == 4 {
						strType = "(nonatomic, assign) NSInteger "
					} else if sIntValue == 5 {
						strType = "(nonatomic, assign) NSInteger "
					} else if sIntValue == 6 {
						strType = "(nonatomic, strong) NSString *"
					} else if sIntValue == 7 {
						strType = "(nonatomic, strong) NSArray *"
					}
				}
			} else {
				//其他类型用|符号分割，go为typeArr[0]，oc为typeArr[1]
				if strings.Contains(arrs[1], "|") {
					typeArr := strings.Split(arrs[1], "|")
					strType = "(nonatomic, strong) " + typeArr[1] + " "
				} else {
					strType = arrs[1] + " "
				}
			}

			strFinal += "@property " + strType + strField + ";"
		} else if strings.HasPrefix(s, "~") && strings.HasPrefix(s, "~#") == false {
			// 4.匹配`~`，oc独有
			strField := strings.Replace(arrs[0], "~", "", 1)
			if strings.HasPrefix(strField, "//") {
				strFinal += strings.Replace(s, "~", "", 1) + "\n"
			} else {
				strField = lowerFisterLetter(strField)
				var strType string
				sIntValue, err := strconv.Atoi(arrs[1])
				if err == nil {
					if sIntValue >= 0 || sIntValue <= 6 {
						if sIntValue == 0 {
							strType = "(nonatomic, strong) NSString *"
						} else if sIntValue == 1 {
							strType = "(nonatomic, strong) NSString *"
						} else if sIntValue == 2 {
							strType = "(nonatomic, assign) BOOL "
						} else if sIntValue == 3 {
							strType = "(nonatomic, assign) double "
						} else if sIntValue == 4 {
							strType = "(nonatomic, assign) NSInteger "
						} else if sIntValue == 5 {
							strType = "(nonatomic, assign) NSInteger "
						} else if sIntValue == 6 {
							strType = "(nonatomic, strong) NSString *"
						} else if sIntValue == 7 {
							strType = "(nonatomic, strong) NSArray *"
						}
					}
				} else {
					if strings.Contains(s, "|") {
						strArr := strings.Split(s, "|")
						strType = strArr[1]
					}
				}
				strFinal += "@property " + strType + strField + ";\n"
			}
		} else {
			if strings.Contains(s, "|") {
				strArr := strings.Split(s, "|")
				strFinal += strArr[1] + "\n"
			} else if strings.HasPrefix(s, "//==") {
				//匹配`//`
				//指定为`注释`
				strFinal += s + "\n"
			} else if strings.Contains(s, "###") == false && strings.Contains(s, "|") == false {
				strFinal += s + "\n"
			}
		}

		// 4.匹配行内注释
		if len(arrs) > 1 && strings.Contains(s, "###") == false && strings.HasPrefix(s, "//") == false && strings.HasPrefix(s, "package") == false && strings.HasPrefix(s, "~") == false {
			var strComment string
			if strings.Contains(s, "//") {
				strComment = strings.Split(s, "//")[1]
				strFinal += " //" + strComment + "\n"
			} else {
				strFinal += strComment
			}
		}

		//匹配多个###...
		//struct结束
		if strings.HasPrefix(s, "###") || strings.HasPrefix(s, "~###") {
			strFinal += "@end\n\n"
		}
	})
	//生成model文件
	if isForm || isReturn {
		ioutil.WriteFile(savePath+"/"+modelName+"Model.h", []byte(strFinal), 0666)
		if isReturn {
			//不覆盖
			file, _ := os.OpenFile(savePath+"/"+modelName+"Model.m", os.O_CREATE, 0666)
			defer file.Close()
		} else {
			ioutil.WriteFile(savePath+"/"+modelName+"Model.m", []byte(strFinal2), 0666)
		}

	} else {
		ioutil.WriteFile(savePath+"/"+modelName+"/"+modelName+"Model.h", []byte(strFinal), 0666)
		file, _ := os.OpenFile(savePath+"/"+modelName+"/"+modelName+"Model.m", os.O_CREATE, 0666)
		defer file.Close()
	}

}

//GoGen 普通的model
func GoGen(savePath string, modelName string) {
	goGen(savePath, modelName, false, false)
}

//GoGenForm 是提交或者返回model
func GoGenForm(savePath string, modelName string) {
	goGen(savePath, modelName, true, false)
}

//GoGenReturn 是提交或者返回model
func GoGenReturn(savePath string, modelName string) {
	goGen(savePath, modelName, false, true)
}

//goGen savePath保存路径(最后不带斜线)，如 /Users/cooerson/Documents/go/src/action/models
//modelName,文件名,结尾自动加上.go
func goGen(savePath string, modelName string, isForm bool, isReturn bool) {
	var strFinal string
	var strStruct string
	readLine(modelName, func(s string) {
		//按空格分割此行字符串
		arrs := strings.Fields(s)
		if strings.HasPrefix(s, "~") {
			return
		} else if strings.HasPrefix(s, "=") || strings.HasPrefix(s, "@=") {
			//匹配`=xx`或`=Xx`
			//指定为struct的名称Xx
			if strings.HasPrefix(s, "=") {
				s = strings.Replace(s, "=", "", 1)
			} else if strings.HasPrefix(s, "@=") {
				s = strings.Replace(s, "@=", "", 1)
			}
			strStruct = strings.Title(s) //首字母大写
			strFinal += "type " + strStruct + " struct {\n"
		} else if strings.HasPrefix(s, "-ID") || strings.HasPrefix(s, "@ID") {
			//匹配`-ID`或`@id`
			//指定字段名为 `ID`
			//指定类型为 `bson.ObjectId`
			//指定json为 `-Xx`或`-xx`的`xxID``
			//指定bson为 `_id,omitempty`
			//例如 ID bson.ObjectId `json:"userID" bson:"_id"`
			strStruct := lowerFisterLetter(strStruct)
			if isReturn {
				strFinal += "ID bson.ObjectId `json:\"" + strStruct + "ID\"`"
			} else {
				strFinal += "ID bson.ObjectId `json:\"" + strStruct + "ID\" bson:\"_id,omitempty\"`"
			}
		} else if (strings.HasPrefix(s, "-") || strings.HasPrefix(s, "@")) && strings.HasPrefix(s, "@###") == false && strings.HasPrefix(s, "@//") == false {
			//1.匹配`-xx`或`-Xx`
			//  指定字段名为`Xx，首字母大写
			var strField string
			if strings.HasPrefix(s, "-") {
				strField = strings.Replace(arrs[0], "-", "", 1)
			} else if strings.HasPrefix(s, "@") {
				strField = strings.Replace(arrs[0], "@", "", 1)
			}
			strField = strings.Title(strField)
			//2.匹配类型：例如0、1、2、3、4、5、6、自定义字符串
			//  0:指定字段类型为`bson.ObjectId`,同时指定bson为`omitempty`
			//  1:指定字段类型为`string`
			//  2:指定字段类型为`bool` //如果是生成 form，bool 变成 int8
			//  3:指定字段类型为`float64`
			//  4:指定字段类型为`int64`
			//  5:指定字段类型为`int8`
			//  6:指定字段类型为`time.Time`
			var strType string
			sIntValue, err := strconv.Atoi(arrs[1])
			if err == nil {
				if sIntValue >= 0 || sIntValue <= 6 {
					if sIntValue == 0 {
						if isForm {
							strType = "string"
						} else {
							strType = "bson.ObjectId"
						}
					} else if sIntValue == 1 {
						strType = "string"
					} else if sIntValue == 2 {
						strType = "bool"
					} else if sIntValue == 3 {
						strType = "float64"
					} else if sIntValue == 4 {
						strType = "int64"
					} else if sIntValue == 5 {
						strType = "int8"
					} else if sIntValue == 6 {
						strType = "time.Time"
					} else if sIntValue == 7 {
						strType = "[]string"
					}
				}
			} else {
				//其他类型用|符号分割，go为typeArr[0]，oc为typeArr[1]
				if strings.Contains(arrs[1], "|") {
					typeArr := strings.Split(arrs[1], "|")
					strType = typeArr[0]
				} else {
					strType = arrs[1]
				}
			}

			//3.匹配json/bson设定
			strJOrB := lowerFisterLetter(strField)
			if isForm {
				if len(arrs) > 2 && strings.HasPrefix(arrs[2], "-bfo") {
					//匹配-bfo
					strJOrB = "`json:\"" + strJOrB + "\" bson:\"" + strJOrB + ",omitempty\"`"
				} else if len(arrs) > 2 && strings.HasPrefix(arrs[2], "-fj") {
					//匹配-f、-j
					strJOrB = "`form:\"" + strJOrB + "\" json:\"" + strJOrB + "\"`"
				} else if len(arrs) > 2 && strings.HasPrefix(arrs[2], "-fjbo") {
					//匹配-f、-bfo、-j
					strJOrB = "`form:\"" + strJOrB + "\" json:\"" + strJOrB + "\" bson:\"" + strJOrB + ",omitempty\"`"
				} else if len(arrs) > 2 && strings.HasPrefix(arrs[2], "-f") {
					//匹配-f、无-bfo
					strJOrB = "`form:\"" + strJOrB + "\"`"
				} else {
					//匹配无-f、无-bfo、无-fjbo
					if sIntValue == 0 {
						//form中，带ID的字段，必须加omitempty
						strJOrB = "`json:\"" + strJOrB + "\" bson:\"" + strJOrB + ",omitempty\"`"
					} else {
						strJOrB = "`json:\"" + strJOrB + "\"`"
					}
				}
			} else if isReturn {
				strJOrB = "`json:\"" + strJOrB + "\"`"
			} else {
				if len(arrs) > 2 && strings.HasPrefix(arrs[2], "//") {
					//3.匹配 无指定json\bson内容

					// 默认指定json为 首字母小写，忽略空
					// 默认指定bson为 首字母小写，不忽略空(带ID的字段除外)
					if sIntValue == 0 {
						strJOrB = "`json:\"" + strJOrB + ",omitempty\" bson:\"" + strJOrB + ",omitempty\"`"
					} else {
						strJOrB = "`json:\"" + strJOrB + ",omitempty\" bson:\"" + strJOrB + "\"`"
					}
				} else {
					//3.匹配 有指定的json\bson内容
					// 匹配 `j:xx`或`Xx`  指定json为 `xx` 首字母变小写
					// 匹配 `b:xx`或`Xx`  指定bson为 `xx` 首字母变小写
					// 匹配 `jo:xx`或`Xx`  指定json为 `xx,omitempty` 首字母变小写
					// 匹配 `bo:xx`或`Xx`  指定bson为 `xx,omitempty` 首字母变小写
					// 匹配 `-j`或`－b`  指定为 `json:"-"` 或 `bson:"-"`
					// 匹配 `-jb` 指定为 空
					needLastBackQuate := false

					if len(arrs) > 2 {
						if strings.HasPrefix(arrs[2], "j:") {
							strJ := strings.Replace(arrs[2], "j:", "", 1)
							strJ = lowerFisterLetter(strJ)
							strJOrB = "`json:\"" + strJ
							needLastBackQuate = true
						} else if strings.HasPrefix(arrs[2], "jo:") {
							strJ := strings.Replace(arrs[2], "jo:", "", 1)
							strJ = lowerFisterLetter(strJ)
							strJOrB = "`json:\"" + strJOrB + ",omitempty"
							needLastBackQuate = true
						} else if strings.HasPrefix(arrs[2], "-jb") {
							strJOrB = ""
						} else if strings.HasPrefix(arrs[2], "-j") {
							strJOrB = "`json:\"-\" bson:\"" + strJOrB + ",omitempty\""
							needLastBackQuate = true
						} else if strings.HasPrefix(arrs[2], "b:") {
							strB := strings.Replace(arrs[2], "b:", "", 1)
							strB = lowerFisterLetter(strB)
							strJOrB = "`json:\"" + strJOrB + ",omitempty\" bson:\"" + strB + "\"`"
						} else if strings.HasPrefix(arrs[2], "bo:") {
							strB := strings.Replace(arrs[2], "bo:", "", 1)
							strB = lowerFisterLetter(strB)
							strJOrB = "`json:\"" + strJOrB + ",omitempty\" bson:\"" + strB + ",omitempty\"`"
						} else if strings.HasPrefix(arrs[2], "-b") {
							strJOrB = "`json:\"" + strJOrB + ",omitempty\" bson:\"-\"`"
						}

					}
					if len(arrs) > 3 {
						if strings.HasPrefix(arrs[3], "b:") {
							strB := strings.Replace(arrs[3], "b:", "", 1)
							strB = lowerFisterLetter(strB)
							strJOrB += " bson:\"" + strB
						} else if strings.HasPrefix(arrs[3], "bo:") {
							strB := strings.Replace(arrs[3], "bo:", "", 1)
							strB = lowerFisterLetter(strB)
							strJOrB += " bson:\"" + strB + ",omitempty"
						} else if strings.HasPrefix(arrs[3], "-b") {
							strJOrB += "bson:\"-\""
						}
					}

					if needLastBackQuate {
						strJOrB += "`"
					}
				}
			}

			strFinal += strField + " " + strType + " " + strJOrB
		} else {
			if strings.HasPrefix(s, "@//") {
				s = strings.Replace(s, "@", "", 1)
				strFinal += s + "\n"
			} else if strings.Contains(s, "|") {
				strArr := strings.Split(s, "|")
				strFinal += strArr[0] + "\n"
			} else if strings.Contains(s, "###") == false {
				strFinal += s + "\n"
			} else if strings.HasPrefix(s, "//==") {
				//匹配`//`
				//指定为`注释`
				strFinal += s + "\n"
			}
		}

		// 4.匹配行内注释
		if len(arrs) > 1 && strings.Contains(s, "###") == false && strings.HasPrefix(s, "//") == false && strings.HasPrefix(s, "package") == false && strings.HasPrefix(s, "~") == false && strings.HasPrefix(s, "@//") == false {
			var strComment string
			if strings.Contains(s, "//") {
				strComment = strings.Split(s, "//")[1]
				strFinal += " //" + strComment + "\n"
			} else {
				strFinal += strComment + "\n"
			}
		}

		//匹配多个###...
		//struct结束
		if strings.HasPrefix(s, "###") || strings.HasPrefix(s, "@###") {
			strFinal += "}\n\n"
		}
	})
	//生成model文件
	filePath := savePath + "/" + modelName + "Model.go"
	ioutil.WriteFile(filePath, []byte(strFinal), 0666)

	//导入包、格式化代码
	cmd := exec.Command("goreturns", "-w", filePath)
	cmd.Run()
}

func readLine(fileName string, handler func(string)) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func lowerFisterLetter(str string) string {
	arrLetter := strings.Split(str, "")
	firstLetter := strings.ToLower(arrLetter[0])
	arrLetter = append(arrLetter[:0], arrLetter[0+1:]...)                                 //删除第一个大写字母
	arrLetter = append(arrLetter[:0], append([]string{firstLetter}, arrLetter[0:]...)...) //添加第一个小写字母
	return strings.Join(arrLetter, "")
}
