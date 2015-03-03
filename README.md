# Go HTTP proxy ![](https://travis-ci.org/marinhero/goxy.svg?branch=master)

Go HTTP proxy or goxy (cheesy) is a reverse proxy that receives a test file what will be interpreted a a blacklist, all the outgoing requests for any given host inside the blacklist will return 404. 

This is particularly useful if you try to block ad providers or in a simple and quick way blacklist any distraction that the web can provide you.

## Install

*Go language (1.4) is required*

`$> git clone https://github.com/marinhero/goxy`

`$> cd goxy`

`$> go install`

`$> goxy --list=blacklist.txt`


## Demo
![](http://marinhero.com/media/goxy-demo.png)

## Achievements

* Full support over HTTP (POST, GET, PUT, etc)
* Cookies and Headers are preserved
* Productivity boost :)

## Possible improvements

* HTTPS Support
* Faster blacklisting
* Improve stability
* Serious testing

###Contact

Marin Alcaraz

marin.alcaraz@gmail.com

http://marinhero.com



