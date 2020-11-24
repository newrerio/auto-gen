
## 鸟币app使用的能自动生成 golang model 和 ios model代码的代码。

### 快速上手：  

以鸟币的oc和go的model生成为示例    
1. 修改在ibisModel/main.go 中的path  

	  //把这两项替换成你自己的
	  goModelPath := "/Users/cooerson/Documents/go/src/niaoshenhao.com/ibis/models"
	  iosModelPath := "/Users/cooerson/Documents/ios/Ibis-ios/Ibiscoin/Ibiscoin/Model"

  
2. ：执行 go run main.go 即可



## 书写规则
1. 示例

```
-ID 0 //用户ID
-Phone 1 //用户手机号，不可重复
-pwd 1 -j //密码不返回
-avatar utils.Pic //头像 原图/640/320/160
-userActivity 0 bo:activity
-isFounder 2 jo:founder b:founder //是否是发起人
```

## go匹配规则

1. 匹配以`//`开头的一行  
指定为`注释`  

1. 匹配`=xx`或`Xx`  
指定为struct的名称Xx  

1. 匹配`-ID`或`-id`或`-Id`  
指定字段名为 `ID`  
指定类型为 `bson.ObjectId`  
指定json为 -Xx或-xx的XxID  
指定bson为 `_id,omitempty`  

1. 匹配`-xx`或`-Xx`  
指定字段名为`Xx`  

1. 匹配类型：例如0、1、2、3、4、5、6、自定义字符串  
0:指定字段类型为`bson.ObjectId`,同时指定bson为`omitempty`  
1:指定字段类型为`string`  
2:指定字段类型为`bool`  
3:指定字段类型为`float64`  
4:指定字段类型为`int64`  
5:指定字段类型为`int8`  
6:指定字段类型为`time.Time`  

1. 匹配 无指定json\bson内容  
默认指定json为 首字母小写，忽略空  
默认指定bson为 首字母小写  

1. 匹配 `j:xx`或`Xx`  
指定json为 xx 首字母变小写  

1. 匹配 `jo:xx`或`Xx`  
指定json为 xx,omitempty 首字母变小写  

1. 匹配 `b:xx`或`Xx`  
指定bson为 xx 首字母变小写  

1. 匹配 `bo:xx`或`Xx`  
指定bson为 `xx,omitempty` 首字母变小写  

1. 匹配 `-j`或`－b`  
指定为 `json:"-"` 或 `bson:"-"`  

1. 匹配 `@`
指定为go独有

1. 匹配`多个###`  
指定为一个struct结束  



## oc匹配规则

在go的匹配规则基础上，

1. 匹配`~`
指定为oc独有

1. 匹配`=xx`或`=Xx`  
指定为interface的名称Xx  

1. 匹配`-ID`或`-id`或`-Id` 
指定字段名为 `xxID` 
指定类型为 `NSString` 

1. 匹配`-xx`或`-Xx`  
指定字段名为`Xx，首字母小写

1. 匹配类型：例如0、1、2、3、4、5、6、自定义字符串  
0:指定字段类型为`bson.ObjectId`->`(nonatomic, strong) NSString`  
1:指定字段类型为`string`->`(nonatomic, strong) NSString`  
2:指定字段类型为`bool`->`(nonatomic) BOOL`  
3:指定字段类型为`float64`->`(nonatomic) double`  
4:指定字段类型为`int64`->`(nonatomic) NSInteger`  
5:指定字段类型为`int8`->`(nonatomic) NSInteger`  
6:指定字段类型为`time.Time`->`(nonatomic, strong) NSString`  

1. 自定义类型用|符号分割  
不带空格，go的index为[0]，oc的index为[1]  
如 -weixin utils.Pic|Pic*|  

1. 匹配`多个###`  
指定为一个interface结束

## 在以上两个的基础上

1. Form 文件中是所有提交数据所要用到的字段  
Return 文件中是所有返回数据中所要用的到的字段

2. 字段属性后全部变成json
#### 注意，gin 的默认 form 类型是 json，所以 formsmodel 中，提交图片时，form类型是"multipart/form-data"，记得在 model中用 form 标记，而不是 json 标记
改用 form标记加上 -f



		1. f(form)\fj（form+json）\bfo(form+bson,omitempty)  \fjbo(form+json+bson,omitempty)  
		2. '新建的'form 和 '更新的'form 中，xxID的bson需要为omitempty  
		3. 字段不接受为空的更新（默认值可以为空，只是更新不为空），匹配 -bfo (bson form omitempty)，指定 bson 为omitempty  


## 其他
指定 tag 的必须按照 f-->j(包括jo、-j)-->b(包括bo、fbo、-b)的前后顺序出现！