package http

const errorTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>{{ .Code }}</title>
    <style>
    html {
        background-color: #2196f3;
        color: #fff;
        font-family: sans-serif;
    }
    code {
        background-color: rgba(0,0,0,0.1);
        border-radius: 5px;
        padding: 1em;
        display: block;
        box-sizing: border-box;
    }
    .center {
        max-width: 40em;
        margin: 2em auto 0;
    }
    a {
        text-decoration: none;
        color: #eee;
        font-weight: bold;
    }
	p {
		line-height: 1.3;
	}
    </style>
</head>

<body>
    <div class="center">
        <h1>Error {{ .Code }}</h1>
        <p>{{ .Message }}</p>
        <code>{{ .ID }}</code>
    </div>
</html>`
