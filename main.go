/*
=================================================================
Proxy : go env -w GOPROXY=https://goproxy.cn/,direct go build

Init  : go init gm_csv
Mod   : go get -u github.com/gin-gonic/gin

Run   : go run main.go
        GIN_MODE="release" && go run main.go
        $env:GIN_MODE="release"; go run main.go

Build : go build -o gm_csv.exe main.go
		GIN_MODE="release" && go build -o gm_csv.exe main.go     # Linux
		$env:GIN_MODE="release"; go build -o gm_csv.exe main.go  # Windows Powershell

URl:
	https://github.com/lmzxtek/gmcsv_go/releases/latest/download/gmcsv-linux-amd64
	https://github.com/lmzxtek/gmcsv_go/releases/latest/download/gmcsv-windows-amd64.exe

=================================================================
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Debug     bool   `json:"debug"`
	FldData   string `json:"fld_data"`
	ServerTag string `json:"servertag"`
}

// var cfg = &Config{
// 	Debug:     false,
// 	Port:      5002,
// 	FldData:   "data_year",
// 	ServerTag: "gmcsv",
// }

var cfg Config
var baseDir string

func main() {
	// 设置为 Release 模式
	gin.SetMode(gin.ReleaseMode)
	fmt.Printf("\n Set gin mode to Release\n")

	// 解析命令行参数
	var configFile string
	flag.StringVar(&configFile, "c", "gmcsv.json", "配置文件路径")
	flag.Parse()

	// 加载或创建配置
	cfg, err := LoadOrCreateConfig(configFile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nLoaded config: %+v\n", cfg)

	baseDir, _ = filepath.Abs(cfg.FldData)
	os.MkdirAll(baseDir, os.ModePerm)

	r := gin.Default()

	r.GET("/test", routeTest)
	r.GET("/test2", routeTest2)
	r.GET("/test3", routeTest3)
	r.GET("/usage", routeUsage)
	r.GET("/download/*filepath", routeDownload)
	r.POST("/upload/*filepath", routeUpload)

	// addr := fmt.Sprintf("%s:%d", cfg.Host, config.Port)
	addr := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("\nServer running at http://*%s\n\n", addr)
	r.Run(addr)
}

// 修改后的配置加载函数，接收文件名参数
func LoadOrCreateConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		defaultCfg := &Config{
			Debug:     false,
			Port:      5002,
			FldData:   "data_year",
			ServerTag: "gmcsv",
		}
		if os.IsNotExist(err) {
			// 文件不存在时创建默认配置

			data, err := json.MarshalIndent(defaultCfg, "", "    ")
			if err != nil {
				return defaultCfg, fmt.Errorf("JSON编码失败: %v", err)
			}

			if err := os.WriteFile(filename, data, 0644); err != nil {
				return defaultCfg, fmt.Errorf("文件创建失败: %v", err)
			}

			fmt.Printf(" %s配置文件不存在，已创建默认配置\n", filename)
			return defaultCfg, nil
		}
		return defaultCfg, fmt.Errorf("文件打开错误: %v, \n 使用默认参数 ", err)
	}
	defer file.Close()

	// 文件存在时解析配置
	// var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("配置解析失败: %v", err)
	}

	return &cfg, nil
}

func getFilePathYear(symbol string, tag string, year int) string {
	// 构造行情数据文件路径
	key := fmt.Sprintf("%s-%s", symbol[:2], symbol[5:7])
	var subfld string
	if tag == "vv" || tag == "pe" {
		subfld = fmt.Sprintf("kbars-%s/%s-%d/%s-%d--%s/", tag, tag, year, tag, year, key)
	} else {
		subfld = fmt.Sprintf("kbars-year/year-%d/year-%d--%s/", year, year, key)
	}
	fname := fmt.Sprintf("kbars-%s--%s--%d-.csv.xz", tag, symbol, year)
	fpath := fmt.Sprintf("%s%s", subfld, fname)
	// fpath := filepath.Join(subfld, fname)
	return fpath
}

func getFilePathMonth(symbol string, year int, month int) string {
	// 构造分时行情文件路径
	tag := "1m"
	key := fmt.Sprintf("%s-%s", symbol[:2], symbol[5:7])
	subfld := fmt.Sprintf("kbars-month/month-%d/month-%d-%02d--%s/", year, year, month, key)
	fname := fmt.Sprintf("kbars-%s--%s--%d-%02d-.csv.xz", tag, symbol, year, month)
	fpath := fmt.Sprintf("%s%s", subfld, fname)
	// fpath := filepath.Join(subfld, fname)
	return fpath
}

func createTestFile() {
	// 如果 test.txt 不存在，创建它
	testFile := filepath.Join(baseDir, "test.txt")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		os.WriteFile(testFile, []byte("This is a file for url request test."), 0644)
	}
}

// RandomStock 从预定义的股票代码数组中随机返回一个元素
func randomStock() string {
	// var localRand *rand.Rand
	// stocks := []string{"AAPL", "GOOG", "AMZN", "MSFT", "TSLA"}
	// // 使用本地的随机数生成器
	// return stocks[localRand.Intn(len(stocks))]

	stocks := []string{"AAPL", "GOOG", "AMZN", "MSFT", "TSLA"}

	// 使用当前时间作为随机种子
	// rand.Seed(time.Now().UnixNano())

	// 生成随机索引
	randomIndex := rand.Intn(len(stocks))

	return stocks[randomIndex]
}

// randomIntWithRange 返回一个在 [a-b, a+b] 区间内的随机整数
func randomIntWithRange(a, b int) int {
	// 使用本地的随机数生成器生成一个在 [a-b, a+b] 区间内的随机整数
	// var localRand *rand.Rand
	// return localRand.Intn(2*b+1) + a - b

	// 计算最小值和最大值
	min := a - b
	max := a + b

	// 生成范围内的随机整数
	return rand.Intn(max-min+1) + min
}

// randomDateInLastYear 返回过去一年之内的随机日期
func randomDateInLastYear() string {
	// var localRand *rand.Rand
	// 获取当前时间
	now := time.Now()
	// 计算一年前的时间
	oneYearAgo := now.AddDate(-1, 0, 0)
	// 计算一年前和现在的时间差（以秒为单位）
	diff := now.Unix() - oneYearAgo.Unix()
	// 生成一个在过去一年内的随机时间点（以秒为单位）
	randomUnix := oneYearAgo.Unix() + int64(rand.Intn(int(diff)))
	// 将随机时间点转换为时间对象
	randomTime := time.Unix(randomUnix, 0)
	// 格式化随机时间为 "YYYY-MM-DD" 格式
	return randomTime.Format("2006-01-02")
}

func routeTest(c *gin.Context) {
	createTestFile()

	pa := randomIntWithRange(100, 20)
	pp := int(float64(pa) * 0.1)
	va := randomIntWithRange(100000, 1000)
	vv := int(float64(pa) * 0.1)

	// symbol := c.DefaultQuery("symbol", "AAPL")
	// time := c.DefaultQuery("time", "2024-01-01")
	data := gin.H{
		"Symbol": []string{randomStock(), randomStock(), randomStock(), randomStock(), randomStock()},
		"Time":   []string{randomDateInLastYear(), randomDateInLastYear(), randomDateInLastYear(), randomDateInLastYear(), randomDateInLastYear()},
		"Price":  []int{randomIntWithRange(pa, pp), randomIntWithRange(pa, pp), randomIntWithRange(pa, pp), randomIntWithRange(pa, pp), randomIntWithRange(pa, pp)},
		"Volume": []int{randomIntWithRange(va, vv), randomIntWithRange(va, vv), randomIntWithRange(va, vv), randomIntWithRange(va, vv), randomIntWithRange(va, vv)},
	}
	c.JSON(http.StatusOK, data)
}

func routeTest2(c *gin.Context) {
	createTestFile()

	pa := randomIntWithRange(100, 20)
	pp := int(float64(pa) * 0.1)
	va := randomIntWithRange(100000, 1000)
	vv := int(float64(pa) * 0.1)

	data := gin.H{
		"columns": []string{"Symbol", "Time", "Price", "Volume"},
		"data": [][]any{
			{randomStock(), randomDateInLastYear(), randomIntWithRange(pa, pp), randomIntWithRange(va, vv)},
			{randomStock(), randomDateInLastYear(), randomIntWithRange(pa, pp), randomIntWithRange(va, vv)},
			{randomStock(), randomDateInLastYear(), randomIntWithRange(pa, pp), randomIntWithRange(va, vv)},
			{randomStock(), randomDateInLastYear(), randomIntWithRange(pa, pp), randomIntWithRange(va, vv)},
			{randomStock(), randomDateInLastYear(), randomIntWithRange(pa, pp), randomIntWithRange(va, vv)},
			// {"AAPL", "2025-05-01", 100, 200},
			// {"AAPL", "2025-05-01", 100, 200},
			// {"AAPL", "2025-05-01", 100, 200},
			// {"AAPL", "2025-05-01", 100, 200},
		},
	}
	c.JSON(http.StatusOK, data)
}

// 定义一个record结构体，包含每个字段
type record struct {
	Symbol string `json:"Symbol"`
	Time   string `json:"Time"`
	Price  int    `json:"Price"`
	Volume int    `json:"Volume"`
}

func routeTest3(c *gin.Context) {
	createTestFile()

	pa := randomIntWithRange(100, 20)
	pp := int(float64(pa) * 0.1)
	va := randomIntWithRange(100000, 1000)
	vv := int(float64(pa) * 0.1)

	// 创建一个包含5个record的切片
	records := make([]record, 5)
	for i := range records {
		records[i] = record{
			Symbol: randomStock(),
			Time:   randomDateInLastYear(),
			Price:  randomIntWithRange(pa, pp),
			Volume: randomIntWithRange(va, vv),
		}
	}

	c.JSON(http.StatusOK, records)
}

func routeDownload(c *gin.Context) {
	fileRelPath := c.Param("filepath")
	filePath := filepath.Join(baseDir, fileRelPath)

	// 防止目录穿越
	if !filepathHasPrefix(filePath, baseDir) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.FileAttachment(filePath, filepath.Base(filePath))
}

func routeUpload(c *gin.Context) {
	fileRelPath := c.Param("filepath")
	filePath := filepath.Join(baseDir, fileRelPath)

	// 防止目录穿越
	if !filepathHasPrefix(filePath, baseDir) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// 创建父目录
	os.MkdirAll(filepath.Dir(filePath), os.ModePerm)

	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "No file uploaded")
		return
	}

	c.SaveUploadedFile(file, filePath)
	c.String(http.StatusOK, "File uploaded successfully")
}

// 使用Go模板引擎生成
const htmlTemplate = `
<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><title>gm-csv (%s)</title></head>
<body>
    <h1>gm-csv usage manual</h1>
    <h2>Test: (%s)</h2>
    <ul>
        <li>说明: <a href="http://%s/usage" target="_blank">http://%s/usage</a></li>
        <li>测试: <a href="http://%s/test" target="_blank">http://%s/test</a></li>
    </ul>
    <h2>Download(csv.xz)</h2>
    <ul>
        <li>链接: <a href="http://%s/download" target="_blank">http://%s/download</a></li>
        <li>测试: <a href="http://%s/download/test.txt" target="_blank">http://%s/download/test.txt</a></li>
    </ul>
    <h2>Deploy</h2>
    <ul>
        <li>Init  :  go init gm_csv </li>
        <li>Mod   :  go get -u github.com/gin-gonic/gin</li>
        <li>Run   :  go run main.go</li>
        <li>Deploy:  go build -o gm_csv.exe main.go </li>
        <li>Proxy :  go env -w GOPROXY=https://goproxy.cn,direct </li>
    </ul>
</body>
</html>`

// 定义配置结构体
type HTMLConfig struct {
	ServerTag string
	HostURL   string
	Symbol    string
	Sididx    string
}

func RenderTemplate(cfg HTMLConfig) string {
	var buf strings.Builder
	tmpl := template.Must(template.New("html").Parse(htmlTemplate))

	// 通过切片传递参数
	data := []any{
		cfg.ServerTag, cfg.ServerTag,
		cfg.HostURL, cfg.HostURL,
		cfg.HostURL, cfg.HostURL,
		cfg.HostURL, cfg.HostURL,
		cfg.HostURL, cfg.HostURL,
	}

	tmpl.Execute(&buf, data)
	return buf.String()
}

// 生成HTML的构造函数
func BuildHTML(cfg HTMLConfig) string {
	url := cfg.HostURL
	// 获取当前时间
	now := time.Now()
	// 格式化当前日期为 "YYYY-MM-DD" 格式
	today := now.Format("2006-01-02")
	year := now.Year()
	month := int(now.Month())
	// prday := now.AddDate(0, -1, 0).Format("2006-01-02")
	// ydate := now.AddDate(-1, 0, 0).Format("2006-01-02")

	sym := cfg.Symbol
	idx := cfg.Sididx

	fpathMonth1 := getFilePathMonth(sym, year, month)
	fpathMonth2 := getFilePathMonth(idx, year, month)
	fpathMonth3 := getFilePathMonth(sym, year-1, month)
	fpathMonth4 := getFilePathMonth(idx, year-1, month)

	fpathYear1 := getFilePathYear(sym, "1m", year)
	fpathYear2 := getFilePathYear(sym, "pe", year)
	fpathYear3 := getFilePathYear(sym, "vv", year)

	// fmt.Println(fpathYear1, fpathYear2)
	// fmt.Println(fpathMonth1, fpathMonth2)

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><title>GM-csv -> [ %s ] </title></head>
<body>
    <h1>GM-CSV Instructions</h1>
    <h2>服务器  : %s </h2>
    <h3>当前日期: %s </h3>
    <ul>
        <li>说明: <a href="http://%s/usage" target="_blank">http://%s/usage</a></li>
        <li>测试: <a href="http://%s/test" target="_blank">http://%s/test</a></li>
        <li>测试: <a href="http://%s/test2" target="_blank">http://%s/test2</a></li>
        <li>测试: <a href="http://%s/test3" target="_blank">http://%s/test3</a></li>
    </ul>

    <h3>下载接口</h3>
    <ul>
        <li>链接: <a href="http://%s/download" target="_blank">http://%s/download</a></li>
        <li>测试: <a href="http://%s/download/test.txt" target="_blank">http://%s/download/test.txt</a></li>
    </ul>

    <h3>分时数据</h3>
    <ul>
	<li>个股: <a href="http://%s/download/%s" target="_blank">http://%s/download/%s</a></li>
	<li>大盘: <a href="http://%s/download/%s" target="_blank">http://%s/download/%s</a></li>
	<li>个股: <a href="http://%s/download/%s" target="_blank">http://%s/download/%s</a></li>
	<li>大盘: <a href="http://%s/download/%s" target="_blank">http://%s/download/%s</a></li>
    </ul>

    <h3>年度数据</h3>
    <ul>
        <li>1m: <a href="http://%s/download/%s" target="_blank">http://%s/download/%s</a></li>
        <li>pe: <a href="http://%s/download/%s" target="_blank">http://%s/download/%s</a></li>
        <li>vv: <a href="http://%s/download/%s" target="_blank">http://%s/download/%s</a></li>
    </ul>

    <h3>GO语言</h3>
    <ul>
        <li>Init  :  go init gm_csv </li>
        <li>Mod   :  go get -u github.com/gin-gonic/gin</li>
        <li>Run   :  go run main.go</li>
        <li>Deploy:  go build -o gm_csv.exe main.go </li>
        <li>Proxy :  go env -w GOPROXY=https://goproxy.cn,direct </li>
    </ul>
</body>
</html>`,
		cfg.ServerTag, cfg.ServerTag, today,

		url, url,
		url, url,
		url, url,
		url, url,

		url, url,
		url, url,

		url, fpathMonth1, url, fpathMonth1,
		url, fpathMonth2, url, fpathMonth2,
		url, fpathMonth3, url, fpathMonth3,
		url, fpathMonth4, url, fpathMonth4,

		url, fpathYear1, url, fpathYear1,
		url, fpathYear2, url, fpathYear2,
		url, fpathYear3, url, fpathYear3,
		// url, fpathMonth1, url, fpathMonth1,
		// url, fpathYear2, url, fpathYear2,
	)
}

func routeUsage(c *gin.Context) {
	hostURL := c.Request.Host

	// html := RenderTemplate(HTMLConfig{
	// 	ServerTag: config.ServerTag,
	// 	HostURL:   hostURL,
	// })
	html := BuildHTML(HTMLConfig{
		ServerTag: cfg.ServerTag,
		HostURL:   hostURL,
		Symbol:    "SHSE.601088",
		Sididx:    "SHSE.000001",
	})

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// 防止目录穿越攻击
func filepathHasPrefix(path, prefix string) bool {
	rel, err := filepath.Rel(prefix, path)
	return err == nil && !startsWithDotDot(rel)
}

func startsWithDotDot(path string) bool {
	return path == ".." || filepath.HasPrefix(path, ".."+string(filepath.Separator))
}
