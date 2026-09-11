<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Package {{.ImportPath}}</title>
    <meta name="go-import" content="{{.ImportPath}} {{.VCS}} {{.RepositoryURL}}">
    {{if .SourceURL}}
        <meta name="go-source" content="{{.ImportPath}} {{.SourceURL}} {{.SourceDirURL}} {{.SourceFileURL}}">
    {{end}}
    <meta http-equiv="refresh" content="4; url={{.RepositoryURL}}" />
</head>
<body>
   Redirecting to <a href="{{.RepositoryURL}}">Source</a> in 3 seconds...
</body>
</html>