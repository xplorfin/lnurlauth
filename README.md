[![Coverage Status](https://coveralls.io/repos/github/xplorfin/lnurlauth/badge.svg?branch=master)](https://coveralls.io/github/xplorfin/lnurlauth?branch=master)
[![Renovate enabled](https://img.shields.io/badge/renovate-enabled-brightgreen.svg)](https://app.renovatebot.com/dashboard#github/xplorfin/lnurlauth)
[![Build status](https://github.com/xplorfin/lnurlauth/workflows/test/badge.svg)](https://github.com/xplorfin/lnurlauth/actions?query=workflow%3Atest)
[![Build status](https://github.com/xplorfin/lnurlauth/workflows/goreleaser/badge.svg)](https://github.com/xplorfin/lnurlauth/actions?query=workflow%3Agoreleaser)
[![](https://godoc.org/github.com/xplorfin/lnurlauth?status.svg)](https://godoc.org/github.com/xplorfin/lnurlauth)
[![Go Report Card](https://goreportcard.com/badge/github.com/xplorfin/lnurlauth)](https://goreportcard.com/report/github.com/xplorfin/lnurlauth)

# LN Url

A golang lnurl example implementation. As per the [lnurl-rfc](https://github.com/fiatjaf/lnurl-rfc) this library also provides:

1. [auth](auth.go): A canonical way to authenticate users with lnurl 

This library is based on [passport-ln-url-auth](https://github.com/chill117/passport-lnurl-auth) and utilizes [go-lnurl](github.com/fiatjaf/go-lnurl).