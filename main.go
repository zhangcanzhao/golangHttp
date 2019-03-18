package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// 数据库字段
type Student struct {
	Name   string `gorm:"column:name"`
	Gender string `gorm:"column:gender"`
	Age    string `gorm:"column:age"`
	Id     string `gorm:"column:id"`
}

// 接收参数结构体
type SQueryCondition struct {
	name   string
	gender string
	age    string
	id     string
}

// 操作数据库的表名
func (Student) TableName() string {
	return "student"
}

// 返回数据结构体
type BaseJsonBean struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// -----
func NewBaseJsonBean() *BaseJsonBean {
	return &BaseJsonBean{}
}

// 指针
var db *gorm.DB

func main() {
	// 定义路由
	http.HandleFunc("/test", postData)
	// 开启服务
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		fmt.Println(err)
	}

}
func postData(w http.ResponseWriter, r *http.Request) {
	//链接数据库
	var err error
	result := NewBaseJsonBean()
	db, err := gorm.Open("mysql", "rtK8xf:rtK8xfWjNeqyIz0VE0BO@tcp(116.62.192.194:3306)/webFront?charset=utf8")
	if err != nil {
		fmt.Println(err)
		result.Code = 101
		result.Message = "数据库链接失败"
		return
	} else {
		result.Code = 100
		fmt.Println("connection succedssed")
	}
	// 定义头信息
	fmt.Println("loginTask is running...")
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
	w.Header().Set("Connection", "close")

	if r.Method == "POST" {
		// 接收参数
		r.ParseForm()
		var qryCondition SQueryCondition
		qryCondition.name = r.Form.Get("name")
		qryCondition.id = r.Form.Get("id")
		qryCondition.age = r.Form.Get("age")
		qryCondition.gender = r.Form.Get("gender")

		student := &Student{
			Id:     qryCondition.id,
			Gender: qryCondition.gender,
			Age:    qryCondition.age,
			Name:   qryCondition.name,
		}

		if err := db.Create(student).Error; err != nil {
			result.Code = 101
			fmt.Println("数据添加失败")
			result.Message = "数据库添加失败"
		} else {
			var users []Student
			result.Code = 100
			result.Message = "数据库添加成功"
			result.Data = db.Find(&users)
			fmt.Println("数据添加成功")
			// 修改成功后返回全部数据
			bytes, _ := json.Marshal(result)
			fmt.Fprint(w, string(bytes))
		}
	}
	return

}
