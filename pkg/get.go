package pkg

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-rod/rod"
)

func GetCombinedText(item *rod.Element, sel ...string) string {
	var combinedText strings.Builder
	for _, s := range sel {
		elems, err := item.Elements(s)
		if err != nil || len(elems) == 0 {
			continue
		}
		for _, elem := range elems {
			t, _ := elem.Text()
			combinedText.WriteString(t)
		}
		return combinedText.String()
	}
	return ""
}

func GetAttr(item *rod.Element, attr string, sel ...string) string {
	var elem *rod.Element
	var err error
	if len(sel) > 0 {
		for _, s := range sel {
			elem, err = item.Element(s)
			if err == nil {
				break
			}
		}
	} else {
		elem = item
	}
	if elem == nil {
		return "未找到"
	}
	a, err := elem.Attribute(attr)
	if err != nil || a == nil {
		return "未找到"
	}
	return *a
}
func ExtractLink(item *rod.Element) string {
	gokey, err := item.Attribute("data-gokey")
	if err != nil || gokey == nil || *gokey == "" {
		return "未找到"
	}
	// 正则提取itemId
	re := regexp.MustCompile(`itemId=(\d+)`)
	match := re.FindStringSubmatch(*gokey)
	if len(match) < 2 {
		return "解析失败"
	}
	// 构造1688商品详情链接
	return fmt.Sprintf("https://detail.1688.com/offer/%s.html", match[1])
}
func FilterPurePrice(rawPrice string) string {
	priceRegex := regexp.MustCompile(`[¥￥]\d+(\.\d{2})?`)
	priceMatches := priceRegex.FindAllString(rawPrice, -1)
	if len(priceMatches) > 0 {
		return priceMatches[0]
	}
	return "未找到"
}
