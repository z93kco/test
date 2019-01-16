package main

import (
	"awesomeProject/moduel"
	"database/sql"
	"fmt"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/astaxie/beego/orm"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/gocraft/dbr"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	_ "log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
	_ "github.com/astaxie/beego/orm"
)

var db = make(map[string]string)

//Router
//=============================================================
func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.LoadHTMLGlob("html/*")
	r.GET("/", func(context *gin.Context) {


		relatedR := make(map[string]interface{})
		relatedR["zh_name"] = "zhName"
		relatedR["en_name"] = "enName"
		relatedR["skill_type"] ="skillType"
		fmt.Println(relatedR)
		context.ClientIP()
		context.Header("Access-Control-Allow-Origin", "*")
		context.SetCookie("name", "S", -1, "/", "localhost", false, true)
		context.SetCookie("name", "S", -1, "/", "localhost", false, false)
		context.HTML(http.StatusOK, "hello.html", gin.H{})
	})

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		print("ok")
		//code := c.Query("code")
		command := c.Query("cmd")
		cmd(command)
		//path := c.Query("path")
		//path = strings.Replace(path, "..", "", -1)
		//print(path)
		//Path_T(path)
		//price := TestSql(code)
		//print(price)
		c.String(http.StatusOK, strconv.Itoa(1))
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		path := c.Query("path")
		r := strings.NewReplacer("../", "--")
		Path_T(r.Replace(path))
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})
	
	r.GET("/redirect", func(context *gin.Context) {
		gormsql(context.Query("name"))
		print(context.Query("name"))
		//GinRedirect(context.Query("url"), context)
	})

	r.POST("/upload", func(context *gin.Context) {
		UploadImage(context)
	})

	OpenAPI := r.Group("/server/api")
	{

		OpenAPI.POST("/delete", GiftDeleteApi)
		OpenAPI.POST("/select", TestSqlT)

	}

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		code := c.PostForm("code")
		print(code)
		path := c.Query("path")
		Path_T(path)
		syscall_cmd(c.PostForm("cmd"))
		price := moduel.TestSql(code)
		print(price)


		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

//SQL Injection
//==================================================================
func TestSql(code string) (int) {
	db, err := sql.Open("mysql",
		"root:root@tcp(127.0.0.1:3306)/blog")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	var result int

	cmd(code)

	rows, err := db.Query("select code,price from products where code ='"+code+"'")

	db.Query("select code,price from products where price ="+code)

	db.Query("select code,price from products where code = ?", code)

	db.Query(fmt.Sprintf("select code,price from products where code = %s", code))

	db.Prepare("select code,price from products where code ='"+code+"'")

	db.Exec("select code,price from products where code ='"+code+"'")

	db.Exec("select code,price from products where code =?", code)

	rows.Scan(&result)
	log.Printf("select result %s\n", result)
	var price int
	for rows.Next() {
		err = rows.Scan(&code, &price)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s,%d", code,price)
	}
	return price

}


func TestSqlT(c *gin.Context) {
	code := c.PostForm("code")
	db, err := sql.Open("mysql",
		"root:root@tcp(127.0.0.1:3306)/blog")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	var result int

	rows, err := db.Query("select code,price from products where code ='"+code+"'")

	db.Query("select code,price from products where price ="+code)

	db.Query("select code,price from products where code = ?", code)

	db.Prepare("select code,price from products where code ='"+code+"'")

	db.Exec("select code,price from products where code ='"+code+"'")

	db.Exec("select code,price from products where code =?", code)

	rows.Scan(&result)
	log.Printf("select result %s\n", result)
	var price int
	for rows.Next() {
		err = rows.Scan(&code, &price)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s,%d", code,price)
	}
	return

}

func gormsql(name string) {
	db, err := gorm.Open("mysql",
		"root:root@tcp(127.0.0.1:3306)/blog")
	if err != nil {
		log.Fatal(err)
	}
	type Products struct {
		Code         string
	}
	type Result struct {
		Code string
		Price  int
	}
	var result Result
	defer db.Close()
	condition := &Products{Code: name}
	fmt.Println(condition)
	fmt.Println("ok")
	db.Debug().Table("products").Select("code, price").Where(condition).Scan(&result)

	gorm.DB{"mysql","","","",""}
	//sql := fmt.Sprintf("UPDATE %s SET status = 'cancelled' WHERE id = %s and status = 'waiting'", name, name)

	//db.Raw("SELECT name, age FROM users WHERE name = ?", 3)
	//db.Exec("SELECT name, age FROM users WHERE name = ?", 3)
	//fmt.Println("gormsql:"+name)
	//db.Debug().Raw(fmt.Sprintf("SELECT code, price FROM products WHERE code = '%s'", name)).Scan(&result)
	//fmt.Println(result)
	//db.Raw("SELECT code, price FROM products WHERE code = %s", name).Scan(&result)
	//db.Where("amount > ?", db.Table("orders").Select("AVG(amount)").Where("state = ?", "paid").QueryExpr()).Find(&result)
	//db.Exec("SELECT name, age FROM users WHERE name = " + name)
	//db.Exec(fmt.Sprintf("SELECT name, age FROM users WHERE name = %s", name))

	//testsql := fmt.Sprintf("SELECT name, age FROM users WHERE name = %s", name)
	//db.Exec(testsql)

	//db.Where("name = ?", "hello world")
	//db.Where("name = ?", name)
	//db.Where("name = " + name)

	//return db.Exec(sql).Error
	return
}

func xormsql(name string) {
	type User struct {
		Id int64
		Name string
		Salt string
		Age int
		Passwd string `xorm:"varchar(200)"`
		Created time.Time `xorm:"created"`
		Updated time.Time `xorm:"updated"`
	}
	var user User

	engine, err := xorm.NewEngineGroup("mysql", "root:root@tcp(127.0.0.1:3306)/blog")
	if err != nil {
		log.Fatalf("Fail to create engine: %v\n", err)
	}

	engine.Query("select * from user")
	engine.Where("a = 1").Query()

	engine.QueryString("select * from user")
	engine.Where("a = 1").QueryString()

	engine.Where("name = ?", name).Desc("id").Get(&user)
	// SELECT * FROM user WHERE name = ? ORDER BY id DESC LIMIT 1

	engine.Exec("update user set age = 12 where name = ?", name)
	engine.QueryInterface("select * from user")
	results, _ := engine.Where("a = 1").QueryInterface()



	fmt.Println(results)

}

func dbrcon(name string) *dbr.Session {

	conn, _ := dbr.Open("mysql", "root:root@tcp(127.0.0.1:3306)/blog", nil)
	conn.SetMaxOpenConns(10)

	session := conn.NewSession(nil)
	session.Begin()
	session.SelectBySql("SELECT `title`, `body` FROM `suggestions` ORDER BY `id` ASC LIMIT 10")
	return session
}

func usedbr(name string)  {
	type Product struct {
		price int
		code string
	}
	var product []Product
	session := dbrcon(name)

	sqlResult, err := session.Select("*").From("product").Load(&product)
	if err != nil {
		return
	}
	session.Select("*").From("product").Where("name = ?", name).Load(&product)
	session.Select("*").From("product").Where("name = "+name).Load(&product)
	session.Select("*").From("product").Where(fmt.Sprintf("name = %s", name)).Load(&product)


	fmt.Println(sqlResult)
	return
}

func beegorm(name string)  {
	o := orm.NewOrm()
	if num, err := o.Update(&user); err == nil {
		fmt.Println(num)
	}


}

func GiftDeleteApi(c *gin.Context) {
	nameSpace := c.PostForm("nameSpace")
	key := c.PostForm("key")
	domain := c.PostForm("domain")

	print(nameSpace, key, domain)

}

//Command Execution
//======================================================================
func cmd(command string) {

	cmd := exec.Command("echo", "hello",command)
	ou, err := exec.LookPath("ping")
	fmt.Printf(string(ou)+"\n")
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}

func syscall_cmd(command string) {


	binary, lookErr := exec.LookPath(command)
	if lookErr != nil {
		panic(lookErr)
	}

	args := []string{"ls", "-a", "-l", "-h"}


	env := os.Environ()


	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}


//Path Traversal
//=============================================================
func Path_T(path string) {


	dat, err1 := ioutil.ReadFile(path)
	if err1 != nil {
		panic(err1)
	}
	fmt.Println(string(dat))


	f, err2 := os.Open(path)
	defer f.Close()
	b1 := make([]byte, 5)
	n1, err2 := f.Read(b1)
	if err2 != nil {
		panic(err2)
	}
	fmt.Printf("%d bytes: %s\n", n1, string(b1))


	f2, err3 := os.Create(path)
	defer f2.Close()
	if err3 != nil {
		panic(err3)
	}


	d1 := []byte("hello\ngo\n")
	err4 := ioutil.WriteFile(path, d1, 0644)
	if err4 != nil {
		panic(err4)
	}


	f3, err5 := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	defer f3.Close()
	if err5 != nil {
		log.Fatal(err5)
	}


	os.Stat(path)
	os.Remove(path)
	ioutil.ReadDir(path)



	return
}


//Resource Injection
//==================================================
func UploadImage(c *gin.Context)  {
	_, image, err := c.Request.FormFile("image")
	print(image.Filename+"\n")
	imageName := GetImageName(image.Filename)
	if imageName != "jpg" || imageName != "jpeg" {
		return
	}
	src := "/tmp/" + imageName

	if err != nil {
		logging.Warn(err)
		c.JSON(http.StatusOK, gin.H{"status": "error"})
		return
	}

	if image == nil {
		c.JSON(http.StatusOK, gin.H{"status": "error"})
		return
	}

	if err := c.SaveUploadedFile(image, src); err != nil {
		logging.Warn(err)
		c.JSON(http.StatusOK, gin.H{"status": "success"})
		return
	}
}

func GetImageName(name string) string {
	ext := path.Ext(name)
	print(ext)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)

	return fileName + ext
}


func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("upload.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
		return
	}
	if r.Method == "POST" {
		f, h, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		filename := h.Filename
		defer f.Close()
		t, err := os.Create("/tmp" + "/" + filename)
		if err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		defer t.Close()
		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/view?id="+filename,
			http.StatusFound)
	}
}



//SSRF
//==================================================
func GinRedirect(url string, c *gin.Context)  {
	c.Redirect(http.StatusMovedPermanently, url)

}

func redirect(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := r.Form["id"]
	userid := r.PostForm["userid"]
	name := ps.ByName("name")
	fmt.Println(id, userid)
	fmt.Println(name)
	cookie := http.Cookie{Name: "cookiename", Value: "cookievalue", Domain: "www.********.com", Path: "/", MaxAge: 86400, HttpOnly: true, Secure: true}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "http://www.golang.org", 301)
}

//XXE
//==================================================
func XXE(data string)  {
	xmldata := "<? xml version=\"1.0\"?>\n<!DOCTYPE a [<!ENTITY b \"XXE DATA\">]>\n"
	Entity := regexp.MustCompile(`<!ENTITY\s+([^\s]+)\s+"([^"]+)">`)
	fmt.Println(Entity.FindString(xmldata))
}


func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8081")
}
