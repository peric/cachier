# Simple Cachier

## Install
`go get github.com/peric/cachier`

## Test
`go test -v -race`

## Usage
```go
import (
    "github.com/peric/cachier"
    "fmt"
)

func main() {
	// define sources 
	cachier.AddSource(&Source{Key: "kiwi", RefreshEverySeconds: 60, Type: cachier.TypeJson})
	cachier.AddSource(&Source{Key: "goavio", RefreshEverySeconds: 600, Type: cachier.TypeJson}) 
	
	// init cache and method that will refresh sources 
	cachier.Init()
	
	// get already initialized data
	value0 := cachier.Get("kiwi", "awesomekey")
	fmt.Println(value0)
	
	// set additional data if needed 
	cachier.Set("kiwi", "awesomekey3", "Best value") 
	
	// get data 
	value1 := cachier.Get("kiwi", "awesomekey3")
	fmt.Println(value1)
}
```

## Tips & Tricks

### Source types
Application is now fetching only from local JSON files. So if we would like to fetch the data from db or API, we would need additional implementations.
But in any case, received data would be structured the same as it is now.

### Refreshing functionality
It's a bit primitive and it currently works just with seconds.

## Keep in mind
This is my second GoLang app and I'm still reading my first book about it :)

![GoLang book](https://lh3.googleusercontent.com/HKLtAzHKJhHJFivw7XJAellfKLM0uNgzngaTVo3pT3DycxgJMbxE7GNhNifVEspYSOcnNWet24JwXz8rVJPS8-rNYn05gBqlTUN7WxNCnJEXHXJcGi2qCuPKA5rOVvXUkPEZC9M3SN4_Lvaj9V-jh5ADPYrhWO-XvusT-jzd_WVKjQawm5dqm2opUjbJKQjqC-BLmezs00unJXv5336bqTzFIWHhDKWR16ttfaA5v-erGtXgpcMgiId7EhKaf96M5-qHIfZ_CRbH6pzuxWK9bYPLH53J6t16WJBW0dE3SADWdkSzPiaHXFmHzZaqRt3oQVJAqB1bPwoYn8qwIiAV-6zEIDJLX6GJD-PqeSImbMx28hMggxaWegviAsRUfzwb97R-YIrW3fctPu3xCTubHPL4Lm0vzcpIJbEeUaIQcovtbLkjtBpymGXrJg6EGC0yeIyVeNngr-E3kJm13_piRItnv7VDhJ7eA76PJoWw4K6l1NrocvsE5hI9IsjtaqO_fuCbNgY_yhfrJOx6YDbxHa_0Xv3jP3Zl97F57TLk6ywkGb0kZ69UZRM2D70dBDq5FC9xXg_lbScbpwMcyM9gyfh_pIkGl_HBupTd_QBxAH2qkGdIx4I5TFDlVfso4OEzjE84i3V0yPeVb2WOTLdkkLgXAjGTQyVoHRrPyr5BtN91GoEncyG0teFB=w569-h758-no)
