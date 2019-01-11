package main

import (
	"bytes"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"path"

	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/russross/blackfriday"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Page struct {
	Id        int64
	Title     string
	Body      string
	Note      string `sql:"type:blob"`
	Source    string `sql:"type:blob"` // include URL
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	WasDisplayedMark bool
	FirstDisplayedAt time.Time

	Image   Binary
	ImageId sql.NullInt64
}

type Binary struct {
	Id        int64
	Header    string
	Body      []byte `sql:"type:blob"`
	PageId    int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

var (
	addr = flag.Bool("addr", false, "find open address and print to final-port.txt")

	db    gorm.DB
	dbErr error

	displayDuration       = 86400 // sec
	queryPagesDisplayable = "(strftime('%s', 'now') - strftime('%s', first_displayed_at)) < " + strconv.Itoa(displayDuration)

	managePath = regexp.MustCompile("^/(edit|save|add)/([[:digit:]]+)$")
	addPath    = regexp.MustCompile("^/add/$")
	viewPath   = regexp.MustCompile("^/([[:digit:]]+)?$")
	imagePath  = regexp.MustCompile("^/image/([[:digit:]]+)?$")

	commonTitle string = "go-digital-signage"

	staticRoot    = "assets"
	currentDir, _ = os.Getwd()
)

func _init() {
	if len(os.Args) > 1 {
		displayDuration, _ = strconv.Atoi(os.Args[1])
		queryPagesDisplayable = "(strftime('%s', 'now') - strftime('%s', first_displayed_at)) < " + strconv.Itoa(displayDuration)
	}
}

func dbconnect(dbfile string) {
	// init db
	db, dbErr = gorm.Open("sqlite3", path.Join(currentDir, dbfile))
	if dbErr != nil {
		log.Fatal(dbErr)
		return
	}

	db.DB()
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.AutoMigrate(&Page{})
	db.AutoMigrate(&Binary{})
}

func truncateTables() {
	db.DropTable(&Page{})
	db.DropTable(&Binary{})
	db.AutoMigrate(&Page{})
	db.AutoMigrate(&Binary{})
}

func makeTitleForHeader(title string) string {
	if title == "" {
		return commonTitle
	}
	return commonTitle + " | " + title
}
func (page *Page) makeTitleForHeader() string {
	return makeTitleForHeader(page.Title)
}

func (page *Page) create() error {
	db.Create(page)
	return nil
}

func (page *Page) save() error {
	db.Save(page)
	return nil
}

func pickupUndisplayedPageRandom() (*Page, error) {
	var page Page
	//if err := db.Model(Page{}).Related(&Binary{}).Where("was_displayed_mark = ?", 0).Order("RANDOM()").First(&page).Error; err != nil {
	if err := db.Model(Page{}).Where("was_displayed_mark = ?", 0).Order("RANDOM()").First(&page).Error; err != nil {
		return &Page{}, err
	}
	return &page, nil
}

func pickupWasDisplayedMarkPage() (*Page, error) {
	var page Page
	//var image Binary
	//if err := db.Model(Page{}).Related(&Binary{}).Where(queryPagesDisplayable).Where(&Page{
	if err := db.Model(Page{}).Where(queryPagesDisplayable).Where(&Page{
		WasDisplayedMark: true,
	}).First(&page).Error; err != nil {
		//return &Page{}, errors.New("record not found")
		return &Page{}, err
	}
	return &page, nil
}

func (page *Page) updateWasDisplayedMarkAndFirstDisplayedTime() {
	tx := db.Begin()
	page.WasDisplayedMark = true
	page.FirstDisplayedAt = time.Now()
	if page.save() != nil {
		tx.Rollback()
	}
	tx.Commit()
}

func resetWasDisplayedMarks() error {
	tx := db.Begin()
	if db.Model(Page{}).Updates(map[string]interface{}{"was_displayed_mark": 0}).Error != nil {
		tx.Rollback()
		return errors.New("Couldn't reset WasDisplayedMarks.")
	}
	tx.Commit()
	return nil
}

func loadPage(Id int64) (*Page, error) {
	var page Page
	//var image Binary
	//if db.Model(Page{}).Related(&Binary{}).First(&page, Id).Error != nil {
	if db.Model(Page{}).First(&page, Id).Error != nil {
		return &Page{}, errors.New("Page not found.")
	}
	return &page, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, page *Page, endpoint string) {
	funcMap := template.FuncMap{
		"safehtml":     func(text string) template.HTML { return template.HTML(text) },
		"safejs":       func(text string) template.JS { return template.JS(text) },
		"safehtmlattr": func(text string) template.HTMLAttr { return template.HTMLAttr(text) },
		"maketitle":    func(text string) string { return makeTitleForHeader(text) },
		"markdowning": func(text string) string {
			output := blackfriday.MarkdownCommon([]byte(text))
			return string(output)
		},
	}
	templates := template.Must(template.New("").Funcs(funcMap).ParseFiles("tmpl/base.html",
		"tmpl/"+tmpl+".html"))
	dat := map[string]interface{}{
		"Page":            page,
		"Endpoint":        endpoint,
		"DisplayDuration": displayDuration + 1,
	}
	err := templates.ExecuteTemplate(w, "base", dat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, int64)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// is URL path valid
		m := managePath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		if m[2] == "" {
			m[2] = "0"
		}
		requestId, err := strconv.Atoi(m[2])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fn(w, r, int64(requestId))
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	// is URL path valid
	m := viewPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	// if access "/"
	if m[1] == "" {
		p, err := pickupWasDisplayedMarkPage()
		if err != nil {
			p, err = pickupUndisplayedPageRandom()
			if err != nil {
				resetWasDisplayedMarks()
				p, _ = pickupUndisplayedPageRandom()
			}
			p.updateWasDisplayedMarkAndFirstDisplayedTime()
		}
		renderTemplate(w, "view", p, "root")
		return
	}

	// else if access "/<number>"
	requestId, err := strconv.Atoi(m[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Id := int64(requestId)
	p, err := loadPage(Id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "view", p, "page")
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	// is URL path valid
	m := imagePath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	// else if access "/<number>"
	requestId, err := strconv.Atoi(m[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Id := int64(requestId)
	var binary Binary
	db.First(&binary, &Binary{Id: Id})
	image := bytes.NewBuffer(binary.Body)
	r.Header.Set("Content-Type", binary.Header)
	fmt.Fprint(w, image)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	// is URL path valid
	m := addPath.MatchString(r.URL.Path)
	if m == false {
		http.NotFound(w, r)
		return
	}
	page := &Page{}
	renderTemplate(w, "edit", page, "add")
}

func editHandler(w http.ResponseWriter, r *http.Request, Id int64) {
	p, err := loadPage(Id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "edit", p, "edit")
}

func saveHandler(w http.ResponseWriter, r *http.Request, Id int64) {
	// Page
	page := Page{
		Title:  r.FormValue("title"),
		Body:   r.FormValue("body"),
		Note:   r.FormValue("note"),
		Source: r.FormValue("source"),
		//CurrentDisplayMark: false,
		WasDisplayedMark: false,
		FirstDisplayedAt: time.Now(),
		Image:            Binary{},
	}
	if Id == 0 {
		page.create()
	} else {
		page.Id = Id
		page.save()
	}

	// Binary
	r.ParseMultipartForm(32 << 20) // 32MiB
	file, handler, err := r.FormFile("image")
	//fmt.Printf("file: %v, handler: %v, err: %v\n", file, handler, err)
	if err == nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//return
		defer file.Close()
		fileContent, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		image := Binary{
			Id:     page.Id,
			Header: handler.Header.Get("Content-Type"),
			Body:   fileContent,
			PageId: page.Id,
		}
		page.Image = image
	}

	if r.FormValue("source") == "1" {
		page.Image = Binary{}
	}
	//fmt.Printf("file: %v\n", file)

	db.Save(&page)

	// Redirect
	idString := strconv.FormatInt(int64(page.Id), 10)
	http.Redirect(w, r, "/edit/"+idString, http.StatusFound)
}
