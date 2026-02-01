package main

import (
	"context"
	"demo/goods/service/basic/logger"
	"demo/pkg"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	InitMysql()
	browser := rod.New().ControlURL(
		launcher.New().Headless(false).
			Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36").
			Set("disable-blink-features", "AutomationControlled").
			MustLaunch(),
	).MustConnect()
	defer browser.MustClose()
	page := browser.MustPage("https://sale.1688.com/factory/u0vjcc4j.html?spm=a260k.home2025.centralDoor.ddoor.66333597BBbHgE&topOfferIds=1005591171200")
	page.MustWaitLoad()
	time.Sleep(2 * time.Second)
	err := page.Timeout(20 * time.Second).MustElement(".offerItem").WaitVisible()
	page.Mouse.Scroll(0, 1000, 3)
	if err != nil {
		logger.Error("加载超时", zap.Error(err))
		return
	}
	offerItems := page.Timeout(20 * time.Second).MustElements(".offerItem")
	if len(offerItems) == 0 {
		logger.Error("未获取到商品")
	}
	fmt.Printf("获取到%d个商品\n", len(offerItems))

	for i, item := range offerItems {
		fmt.Printf("===商品%d===\n", i+1)
		fmt.Printf("商品图片:%s\n", pkg.GetAttr(item, `src`, ".offerImg"))
		fmt.Printf("商品标题:%s\n", pkg.GetCombinedText(item, `.offerTitle`))
		fmt.Printf("商品价格:%s\n", pkg.FilterPurePrice(pkg.GetCombinedText(item, `.offerItem`, `span.text`)))
		fmt.Printf("商品连接:%s\n", pkg.ExtractLink(item))
		commodity := Commodity{
			Title:    pkg.GetCombinedText(item, `.offerTitle`),
			Price:    pkg.FilterPurePrice(pkg.GetCombinedText(item, `.offerItem`, `span.text`)),
			GoodsImg: pkg.GetAttr(item, `src`, ".offerImg"),
			Link:     pkg.ExtractLink(item),
		}
		err = commodity.CommodityCreate(DB)
		if err != nil {
			return
		}
		Upload(pkg.GetAttr(item, `src`, ".offerImg"))
	}
	fmt.Println("商品爬取完成")
	fmt.Println("爬取的数据已经存储到数据库")
	fmt.Println("图片已上传")
}

var (
	DB  *gorm.DB
	err error
)

func InitMysql() {
	once := sync.Once{}
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:4ay1nkal3u8ed77y@tcp(115.190.54.31:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"
	once.Do(func() {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Println("数据库连接失败")
			return
		}
		fmt.Println("数据库连接成功")
	})

	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, _ := DB.DB()
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	err = DB.AutoMigrate(&Commodity{})
	if err != nil {
		fmt.Println("迁移失败")
	}
	fmt.Println("迁移成功")
}

type Commodity struct {
	Title    string `gorm:"type:varchar(255);comment:商品标题"`
	Price    string `gorm:"type:varchar(50);comment:商品价格"`
	GoodsImg string `gorm:"type:varchar(255);comment:商品图片"`
	Link     string `gorm:"type:varchar(255);comment:商品链接"`
}

func (c *Commodity) CommodityCreate(db *gorm.DB) error {
	return db.Create(&c).Error
}

func Upload(img string) string {
	objectName := img
	region := "us-east-1"
	access_key_id := "Ctxwr2HEaMLQjo0Buh5z"
	secret_access_key := "fHlLmXCFBJ82iSrjxea0KndQvIWRsq4tYMb5wV76"
	endpoint := "http://115.190.54.31:9501"

	// build aws.Config
	cfg := aws.Config{
		Region: region,
		EndpointResolver: aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: endpoint,
			}, nil
		}),
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(access_key_id, secret_access_key, "")),
	}

	// build S3 client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	ext := strings.Split(objectName, ".")[1]
	objectName = fmt.Sprintf("%d.%s", time.Now().UnixMilli(), ext)
	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("rui"),
		Key:    aws.String(objectName),
		Body:   strings.NewReader("hello rustfs"),
	})
	if err != nil {
		log.Fatalf("上传文件失败: %v", err)
	}
	return fmt.Sprintf("文件上传成功：%s", objectName)
}
