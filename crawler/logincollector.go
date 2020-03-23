package crawler

// LoginCollector ...
// func LoginCollector(site string) {
// 	// const (
// 	// 	csrfTokenSelector = "#main-container > section.content > main > div > div.auth-form.sign-in-form > form > input[type=hidden]:nth-child(2)"
// 	// )
// 	// new collector
// 	c := colly.NewCollector()

// 	// detect csrf token field
// 	// usually in a form or input field

// 	//CSRF Token Authorization
// 	// c.OnHTML(csrfTokenSelector, func(e *colly.HTMLElement) {
// 	// 	//authenticity token
// 	// 	token := e.Attr("value")
// 	// 	//authenticate
// 	// 	err := c.Post("https://www.makeschool.com/login", map[string]string{"user[email]": "gary.frederick@smash.lpfi.org", "user[password]": "ouoyou12"})
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	}
// 	// 	fmt.Print(c.Cookies("https://www.makeschool.com/dashboard"))
// 	// })

// 	//detect username and password field usually in a form or input field
// 	err := c.Post("https://www.makeschool.com/login", map[string]string{"user[email]": "", "user[password]": ""})
// 	// uh oh panic! :(
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Print(c.Cookies("https://www.makeschool.com/dashboard"))
// 	// attach callbacks after login
// 	c.OnResponse(func(r *colly.Response) {
// 		fmt.Println(r.Body)
// 	})
// 	// start scraping
// 	c.Visit("https://www.makeschool.com/login")
// }
