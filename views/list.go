package views

import "html/template"

var List = template.Must(template.New("list").Parse(list))

const list = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>tw-linkfeed</title>
    <style>
      body { font: 16px/1.3 serif; max-width: 40em; margin: 0 auto; }
      h1 { font-size: 1.7em; margin: 2.6rem 0; }
      h2 { font-size: 1.2em; }
      ul { list-style: none; padding-left: 0; }
      li { margin: 2.6rem 0; }
      blockquote { margin-left: 1.3rem; }
    </style>
  </head>
  <body>
    <header>
      <h1>tw-linkfeed</h1>
    </header>

    <ul>
      {{range .}}
        <li>
          <h2>
            {{with $url := index .Entities.Urls 0}}
              <a href="{{$url.ExpandedUrl}}">{{$url.DisplayUrl}}</a>
            {{end}}
          </h2>
          <blockquote>
            <p>{{.Text}}</p>
            <footer>
              &mdash;
              <cite><a href="{{.User.Url}}">{{.User.Name}}</a></cite>
            </footer>
          </blockquote>
        </li>
      {{end}}
    </ul>
  </body>
</html>`
