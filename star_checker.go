package star_checker

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/urlfetch"

	"encoding/json"
	"io/ioutil"
)

type StarReport struct {
	Title      string
	Star_count int
	Uri        string
}

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/check", check)
}

func root(w http.ResponseWriter, r *http.Request) {
	report := new(StarReport)
	if err := starReportTemplate.Execute(w, report); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func check(w http.ResponseWriter, r *http.Request) {
	target_url := r.FormValue("target_url")

	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	resp, err := client.Get("http://s.hatena.ne.jp/blog.json?uri=" + target_url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	report := new(StarReport)
	err = report.JsonProc(string(contents))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := starReportTemplate.Execute(w, report); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p *StarReport) JsonProc(data string) (err error) {
	err = json.Unmarshal([]byte(data), &p)
	if err != nil {
		return err
	}
	return nil
}

var starReportTemplate = template.Must(template.New("report").Parse(starReportTemplateHTML))

const starReportTemplateHTML = `
<html>
  <body>
    <p>ページタイトル：<b>{{.Title}}</b></p>
    <p>ページURL：<b>{{.Uri}}</b></p>
    <p>スター総計：<b>{{.Star_count}}</b></p>
    <form action="/check" method="post">
      <div><input type="text" name="target_url" /></div>
      <div><input type="submit" value="Check Star Count"></div>
    </form>
  </body>
</html>
`
