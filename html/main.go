package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/sirupsen/logrus"
)

func main() {
	//exampleSimple()
	agentExample()
}

func exampleSimple() {
	page := rod.New().MustConnect().MustPage("https://www.facebook.com/ads/library/?active_status=active&ad_type=all&country=ALL&is_targeted_country=false&media_type=all&search_type=page&source=nav-header&view_all_page_id=104905298804759")
	page.MustWindowFullscreen()
	//page.MustWaitStable().MustScreenshot("a.png")
	time.Sleep(time.Hour)
}

func agentExample() {
	spider, err := NewSpyWebsiteSp3HeadlessSpider(false, time.Duration(60)*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	url := "https://www.tokyo-tiger.com/"
	productCount, err := spider.CrawlWebsite(context.Background(), url, "shopline", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(productCount)
}

type SpyProduct struct {
	ID              string  `json:"id"`
	Handle          string  `json:"handle"`
	Title           string  `json:"title"`
	BodyHTML        string  `json:"body_html"`
	URL             string  `json:"url"`
	Currency        string  `json:"currency"`
	Price           float64 `json:"price"`
	CompareAtPrice  float64 `json:"compare_at_price"`
	Weight          float64 `json:"weight"`
	WeightUnit      string  `json:"weight_unit"`
	Options         string  `json:"options"`  // JSON格式字符串
	Variants        string  `json:"variants"` // JSON格式字符串
	Tags            string  `json:"tags"`     // JSON格式字符串
	Images          string  `json:"images"`   // JSON格式字符串
	WebsiteUrl      string  `json:"website_url"`
	WebsitePlatform string  `json:"website_platform"`
	Referrer        string  `json:"referrer"`
}

type CrawlingProgress struct {
	ProductCount int
	LastUpdate   time.Time
}

// Shopify Shopline Shoplazza 的headless浏览器爬虫方案
type SpyWebsiteSp3HeadlessSpider struct {
	headlessMode       bool
	logger             *logrus.Entry
	stuckCheckInterval time.Duration
}

// NewSpyWebsiteSp3HeadlessSpider 创建一个新的爬虫实例
func NewSpyWebsiteSp3HeadlessSpider(headlessMode bool, stuckCheckInterval time.Duration) (*SpyWebsiteSp3HeadlessSpider, error) {
	logger := logrus.WithFields(logrus.Fields{
		"component": "SpyWebsiteSp3HeadlessSpider",
	})

	if stuckCheckInterval == 0 {
		stuckCheckInterval = 60 * time.Second
	}

	return &SpyWebsiteSp3HeadlessSpider{
		headlessMode:       headlessMode,
		logger:             logger,
		stuckCheckInterval: stuckCheckInterval,
	}, nil
}

// CrawlWebsite 爬取指定网站的商品信息
func (s *SpyWebsiteSp3HeadlessSpider) CrawlWebsite(ctx context.Context, websiteUrl string, websitePlatform string, productChan chan<- *SpyProduct) (productCount int, gerr error) {
	logger := s.logger.WithFields(logrus.Fields{
		"websitePlatform": websitePlatform,
		"websiteUrl":      websiteUrl,
	})
	fmt.Print("开始爬取网站\n")

	timeoutCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 创建一个完成信号channel
	done := make(chan struct{})

	// 用于跟踪爬取进度的变量
	crawlingProgress := CrawlingProgress{ProductCount: 0, LastUpdate: time.Now()}

	// 启动一个goroutine来监控爬虫是否停滞
	go func() {
		stuckCheckInterval := 60 * time.Second
		ticker := time.NewTicker(stuckCheckInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// 如果超过60秒没有发现新产品，认为爬虫可能停滞
				fmt.Printf("爬虫进度: %d, 最后更新时间: %s\n", crawlingProgress.ProductCount, crawlingProgress.LastUpdate.Format(time.RFC3339))
				fmt.Printf("爬虫停滞时间: %s\n", time.Since(crawlingProgress.LastUpdate).String())
				if time.Since(crawlingProgress.LastUpdate) > stuckCheckInterval {
					logger.Warnf("爬虫可能已停滞，超过%v秒未发现新产品\n", stuckCheckInterval)
					// 不直接取消，给它多一次机会
					if time.Since(crawlingProgress.LastUpdate) > 5*stuckCheckInterval {
						logger.Errorf("爬虫已停滞，超过%v秒未发现新产品，强制终止\n", 5*stuckCheckInterval)
						cancel() // 取消上下文，触发超时
						return
					}
				}
			case <-done:
				return
			case <-timeoutCtx.Done():
				return
			}
		}
	}()

	// 创建一个随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 设置常见的现代浏览器 User-Agent 列表
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36 Edg/92.0.902.78",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0",
	}

	// 随机选择一个User-Agent
	randomUserAgent := userAgents[r.Intn(len(userAgents))]

	// 启动goroutine进行爬取，以便能够响应上下文取消
	go func() {
		defer close(done)

		// 创建一个启动器并设置参数
		l := launcher.New().
			Headless(s.headlessMode).
			Set("user-agent", randomUserAgent).
			Set("disable-blink-features", "AutomationControlled"). // 防止被检测为自动化工具
			Set("window-size", "1920,1080").                       // 设置常见的窗口大小
			Set("start-maximized", "true").                        // 最大化窗口
			Set("disable-infobars", "true").                       // 禁用信息栏
			Set("disable-notifications", "true")                   // 禁用通知

		url := l.MustLaunch()
		browser := rod.New().ControlURL(url).MustConnect()
		defer browser.MustClose()

		fmt.Printf("浏览器启动成功: %s\n", url)

		// 创建一个上下文来保存 cookies 等状态
		browser = browser.NoDefaultDevice()

		// 模拟设备
		page := browser.MustPage()

		fmt.Printf("浏览器页面创建成功: %s\n", url)

		// 绕过机器人检测
		page.MustEvalOnNewDocument(`
			Object.defineProperty(navigator, 'webdriver', {
				get: () => false,
			});
			
			// 添加缺失的浏览器特性
			window.chrome = {
				runtime: {},
			};
			
			// 重新定义navigator.languages
			Object.defineProperty(navigator, 'languages', {
				get: () => ['zh-CN', 'zh', 'en-US', 'en'],
			});
		`)

		// 模拟人类行为的操作函数
		humanLikeDelay := func() {
			// 200ms 到 2000ms 的随机延迟
			delay := 200 + r.Intn(180)
			time.Sleep(time.Duration(delay) * time.Millisecond)
			fmt.Printf("模拟人类行为延迟: %dms\n", delay)

			// 检查上下文是否已取消
			select {
			case <-timeoutCtx.Done():
				return
			default:
				// 继续执行
			}
		}

		// 设置请求头 - 针对每个头部分别设置
		page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: randomUserAgent,
		})
		page.MustEval(`() => {
			Object.defineProperty(navigator, 'language', {
				get: () => 'zh-CN'
			});
		}`)

		// 设置 referrer
		page.MustEval(`() => {
			Object.defineProperty(document, 'referrer', {
				get: () => 'https://www.google.com/'
			});
		}`)

		// 使用map来进行URL去重
		visitedURLs := make(map[string]bool)
		allProductURLs := make(map[string]bool)

		// 访问目标网站
		if err := page.Navigate(websiteUrl); err != nil {
			logger.Errorf("导航到网站失败: %v\n", err)
			gerr = err
			return
		}
		visitedURLs[websiteUrl] = true

		fmt.Printf("导航到网站成功: %s\n", websiteUrl)

		// 等待页面加载完成
		page.MustWaitLoad()

		fmt.Printf("页面加载完成\n")

		// 模拟人类延迟和滚动行为
		humanLikeDelay()

		// 查找所有collections页面链接
		collectionURLs := make(map[string]bool)

		// 检查上下文是否已取消
		select {
		case <-timeoutCtx.Done():
			fmt.Printf("爬虫执行超时\n")
			gerr = timeoutCtx.Err()
			return
		default:
			// 继续执行
		}

		navLinks, err := page.Elements("a")
		if err != nil {
			logger.Warnf("查找导航链接失败: %v\n", err)
		} else {
			for _, link := range navLinks {
				href, err := link.Attribute("href")
				if err == nil && href != nil {
					linkHref := *href
					// 寻找包含"collections"关键词的链接
					if strings.Contains(strings.ToLower(linkHref), "/collections/") {
						crawlingProgress.LastUpdate = time.Now()
						collectionURL := linkHref
						if !strings.HasPrefix(collectionURL, "http") {
							if strings.HasPrefix(collectionURL, "/") {
								// 解析基础URL
								baseURL := websiteUrl
								baseURL = strings.TrimSuffix(baseURL, "/")
								collectionURL = baseURL + collectionURL
							} else {
								collectionURL = websiteUrl + "/" + collectionURL
							}
						}
						// 使用map直接去重
						collectionURLs[collectionURL] = true
						fmt.Printf("找到商品集合页面: %s\n", collectionURL)
					}
				}
			}
		}

		if len(collectionURLs) == 0 {
			logger.Warning("未找到任何collections页面，将尝试直接在首页寻找商品\n")
			// 如果没有找到collections链接，则将首页添加到collections列表
			collectionURLs[websiteUrl] = true
		}

		// 遍历每个collections页面
		collectionCount := len(collectionURLs)
		i := 0
		for collectionURL := range collectionURLs {
			i++
			// 检查上下文是否已取消
			select {
			case <-timeoutCtx.Done():
				logger.Error("爬虫执行超时\n")
				gerr = timeoutCtx.Err()
				return
			default:
				// 继续执行
			}

			fmt.Printf("%d/%d 正在爬取商品集合页面: %s\n", i, collectionCount, collectionURL)
			crawlingProgress.LastUpdate = time.Now()

			// 如果已经访问过这个URL，则跳过
			if visitedURLs[collectionURL] {
				fmt.Printf("%d/%d 页面已访问过1，跳过: %s\n", i, collectionCount, collectionURL)
				continue
			}

			// 导航到collections页面
			if err := page.Navigate(collectionURL); err != nil {
				logger.Warnf("%d/%d 导航到集合页面失败: %s, 错误: %v\n", i, collectionCount, collectionURL, err)
				continue
			}

			// 等待页面加载完成
			page.MustWaitLoad()

			humanLikeDelay()
			// humanLikeScroll()

			// 查找分页链接（如果存在）
			paginationURLs := make(map[string]bool)
			paginationLinks, err := page.Elements("a[href*='page=']")
			if err == nil && len(paginationLinks) > 0 {
				for _, pLink := range paginationLinks {
					pHref, err := pLink.Attribute("href")
					if err == nil && pHref != nil {
						paginationURL := *pHref
						if !strings.HasPrefix(paginationURL, "http") {
							if strings.HasPrefix(paginationURL, "/") {
								baseURL := websiteUrl
								baseURL = strings.TrimSuffix(baseURL, "/")
								paginationURL = baseURL + paginationURL
							} else {
								paginationURL = websiteUrl + "/" + paginationURL
							}
						}
						// 使用map直接去重
						paginationURLs[paginationURL] = true
					}
				}
			}

			// 将当前页面添加到分页列表
			if len(paginationURLs) == 0 {
				paginationURLs[collectionURL] = true
				// 再手动加几个分页
				// Shopline平台: page_num=1&page_size=32
				if strings.EqualFold(websitePlatform, "shopline") {
					for i := 1; i <= 10; i++ {
						paginationURLs[fmt.Sprintf("%v?page_num=%v&page_size=32", collectionURL, i)] = true
					}
				} else if strings.EqualFold(websitePlatform, "shoplazza") {
					for i := 1; i <= 10; i++ {
						paginationURLs[fmt.Sprintf("%v?page=%v", collectionURL, i)] = true
					}
				}
			}

			// 遍历每个分页
			for pageURL := range paginationURLs {
				crawlingProgress.LastUpdate = time.Now()
				// 检查上下文是否已取消
				select {
				case <-timeoutCtx.Done():
					fmt.Printf("爬虫执行超时\n")
					gerr = timeoutCtx.Err()
					return
				default:
					// 继续执行
				}

				// 如果已经访问过这个URL，则跳过
				if visitedURLs[pageURL] {
					fmt.Printf("%d/%d 页面已访问过2，跳过: %s\n", i, collectionCount, pageURL)
					continue
				}

				fmt.Printf("%d/%d 正在爬取分页: %s\n", i, collectionCount, pageURL)
				if err := page.Navigate(pageURL); err != nil {
					logger.Warnf("%d/%d 导航到分页失败: %s, 错误: %v\n", i, collectionCount, pageURL, err)
					continue
				}
				visitedURLs[pageURL] = true

				// 等待页面加载完成
				page.MustWaitLoad()

				humanLikeDelay()

				fmt.Printf("%d/%d 开始查找产品链接\n", i, collectionCount)
				// 查找当前页面上的所有产品链接

				var productLinks []*rod.Element
				var err error
				initialProductCount := len(allProductURLs)
				if strings.EqualFold(websitePlatform, "shoplazza") {
					productLinks, err = page.Elements("div[data-track-product-id]")
					if err != nil {
						logger.Warnf("%d/%d 查找产品链接失败: %v\n", i, collectionCount, err)
						continue
					}

					if len(productLinks) == 0 {
						fmt.Printf("%d/%d 页面未发现产品链接, 退出分页: %s\n", i, collectionCount, pageURL)
						break
					}

					for _, link := range productLinks {
						productId, _ := link.Attribute("data-track-product-id")
						productHandle, _ := link.Attribute("data-quick-shop")
						if productId == nil || productHandle == nil {
							continue
						}

						productURL := fmt.Sprintf("%s/products/%s\n", websiteUrl, *productHandle)

						// 使用map进行去重
						if !allProductURLs[productURL] {
							allProductURLs[productURL] = true
							fmt.Printf("%d/%d 找到产品链接: %s\n", i, collectionCount, productURL)

							// 发送到channel中
							select {
							case productChan <- &SpyProduct{ID: *productId, Handle: *productHandle, URL: productURL, Referrer: pageURL}:
								// 更新爬取进度
								crawlingProgress.ProductCount++
							case <-timeoutCtx.Done():
								fmt.Printf("爬虫执行超时，无法发送更多产品链接\n")
								gerr = timeoutCtx.Err()
								return
							}

						}

						crawlingProgress.LastUpdate = time.Now()
					}
				} else {
					productLinks, err = page.Elements("a[href*='/products/']")
					if err != nil {
						logger.Warnf("%d/%d 查找产品链接失败: %v\n", i, collectionCount, err)
						continue
					}

					if len(productLinks) == 0 {
						fmt.Printf("%d/%d 页面未发现产品链接, 退出分页: %s\n", i, collectionCount, pageURL)
						break
					}

					for _, link := range productLinks {
						href, err := link.Attribute("href")
						if err != nil || href == nil {
							continue
						}

						productURL := *href
						if !strings.HasPrefix(productURL, "http") {
							if strings.HasPrefix(productURL, "/") {
								baseURL := websiteUrl
								baseURL = strings.TrimSuffix(baseURL, "/")
								productURL = baseURL + productURL
							} else {
								productURL = websiteUrl + "/" + productURL
							}
						}
						// 所有的产品链接要去掉参数
						productURL = strings.Split(productURL, "?")[0]
						productHandle := ""
						pus := strings.Split(productURL, "/")
						if len(pus) > 0 {
							productHandle = pus[len(pus)-1]
						}

						// 使用map进行去重
						if !allProductURLs[productURL] {
							allProductURLs[productURL] = true
							fmt.Printf("%d/%d 找到产品链接: %s\n", i, collectionCount, productURL)

							// 发送到channel中
							select {
							case productChan <- &SpyProduct{Handle: productHandle, URL: productURL, Referrer: pageURL}:
								// 更新爬取进度
								crawlingProgress.ProductCount++
							case <-timeoutCtx.Done():
								fmt.Printf("爬虫执行超时，无法发送更多产品链接\n")
								gerr = timeoutCtx.Err()
								return
							}
						}

						crawlingProgress.LastUpdate = time.Now()
					}
				}

				// 如果这个页面没有发现新产品，记录日志
				if initialProductCount == len(allProductURLs) {
					fmt.Printf("%d/%d 页面未发现新产品: %s\n", i, collectionCount, pageURL)
				} else {
					fmt.Printf("%d/%d 页面发现 %d 个新产品: %s\n", i, collectionCount, len(allProductURLs)-initialProductCount, pageURL)
				}
			}
		}

		productCount = len(allProductURLs)
		logger.WithFields(logrus.Fields{
			"productCount": productCount,
			"collections":  len(collectionURLs),
			"visitedPages": len(visitedURLs),
		}).Info("爬取完成")
	}()

	// 等待爬虫完成或超时
	select {
	case <-done:
		// 爬虫正常完成
		fmt.Printf("爬取完成，发现 %d 个产品\n", productCount)
	case <-timeoutCtx.Done():
		gerr = timeoutCtx.Err()
		logger.Errorf("爬虫执行超时，最后更新时间: %s\n", crawlingProgress.LastUpdate.Format(time.RFC3339))
	}

	close(productChan)

	return
}
