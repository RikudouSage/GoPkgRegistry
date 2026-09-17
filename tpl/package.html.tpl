<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Package {{.ImportPath}}</title>
    <meta name="go-import" content="{{.ImportPath}} {{.VCS}} {{.RepositoryURL}}">
    {{if .SourceURL}}
        <meta name="go-source" content="{{.ImportPath}} {{.SourceURL}} {{.SourceDirURL}} {{.SourceFileURL}}">
    {{end}}
    {{if eq .RedirectTo "source"}}
        <meta http-equiv="refresh" content="4; url={{.RepositoryURL}}" />
    {{else}}
        <meta http-equiv="refresh" content="4; url={{.GoPkgURL}}" />
    {{end}}
</head>
<body>
   {{if eq .RedirectTo "source"}}
       Redirecting to <a href="{{.RepositoryURL}}">source</a> in 3 seconds...
   {{else}}
       Redirecting to <a href="{{.GoPkgURL}}">pkg.go.dev</a> in 3 seconds...
   {{end}}
</body>
</html>