package integration

import "html/template"

type LoginPageData struct {
	// encoded lnurl string
	Encoded string
	// data uri string
	DataUri template.URL
	// cancel url
	CancelUrl string
}

var HomeTpl = template.Must(template.New("home").Parse(`
<!DOCTYPE html>
<html lang="en">
<body>
	{{ if . }}
		<p> You are authenticated. To logout, go <a href="/logout"> here </a>
	{{ else }}
		<p> You are not authenticated. To login, go <a href="/login"> here </a>
	{{ end }}
</body>
</html>
`))

var LoginPage = template.Must(template.New("login").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Login with lnurl-auth</title>
	<style>
		html, body {
			width: 100%;
			height: 100%;
			color: #fff;
			font-family: sans;
			background: #222;
			margin: 0;
		}
		h1 {
			margin: 0 0 .5rem 0;
		}
		p {
			font-weight: 100;
			line-height: 150%;
			margin: 0 0 1rem 0;
		}
		#content {
			display: flex;
			height: 100%;
			text-align: center;
		}
		.wrap {
			width: 80%;
			max-width: 24rem;
			align-self: center;
			margin: 0 auto;
		}
		#qrcode {
			display: block;
			width: 16rem;
			height: 16rem;
			margin: 0 auto;
		}
		#qrcode img {
			display: block;
			width: 100%;
		}
		.button {
			cursor: pointer;
			color: #fff;
			font-size: 1.1rem;
			line-height: 1.8rem;
			background: #c95b59;
			border: none;
			outline: none;
			border-radius: .4rem;
			opacity: .7;
			padding: .4rem .8rem;
		}
		a.button,
		a.button:hover {
			text-decoration: none;
		}
		.button:hover {
			opacity: 1;
		}
		#buttons {
			display: block;
			margin-top: 1rem;
		}
	</style>
	<meta http-equiv="refresh" content="8">
</head>
<body>

	<div id="content">
		<div class="wrap">
			<h1>Login with lnurl-auth</h1>
			<p>Scan the QR code with an app that supports lnurl-auth</p>
			<a id="qrcode" href="lightning:{{.Encoded}}"><img src="{{.DataUri}}"></a>
			<div id="buttons">
				{{if .CancelUrl }}
					<a href="{{.CancelUrl}}" class="button">Cancel</a>
				{{end}}
			</div>
		</div>
	</div>

</body>
</html>
`))
