/*
* Program: Weather Command
* Dev: MD
* Date: 5/25/2022
* filename: weather.go
* Purpose: for users to enter in at least city along with two
* extra commands for language and units to get a current weather report
* outputted as a discord message.
 */

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

//used to read in the temperature value inside the
//main json response from weather
type wMain struct {
	Temp float32 `json:"temp"`
}

//reads in the main and description from the weather api json response
type wWeather struct {
	Main        string `json:"main"`
	Description string `json:"description"`
}

//the main struct to handle all of the json response from weather api
type WeatherHTTP struct {
	Weather []wWeather `json:"weather"`
	Main    wMain      `json:"main"`
}

//mapped to the cmd command !weather
type WeatherCommander struct{}

/*
	Main functionality of the command which handles
	getting the weather, parsing the json response,
	then outputting the proper weather report with
	the apporiate units and associated weather emoji
*/
func (c WeatherCommander) Handle(s ApiNooter, m Message) {
	//checks for a nil discord command and only working for discord
	n, disCheck := s.(*DiscordNooter)
	if !disCheck {
		return
	}

	info, err := ioutil.ReadFile("key.txt")

	if err != nil {
		log.Fatal(err)
	}

	key := string(info)

	//args slice handles breaking up the user's command
	//url is for forming the proper api call
	args := []string{}
	url := ""

	//check the length of the command to insure at least
	//city exists or the remaining with city
	if len(m.Parsed.Postfix) != 0 {
		args = strings.Split(m.Parsed.Postfix, " ")
		url = fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", args[1], key)

		//struct that holds the json response from openweatherapi
		//emoji is used to pass the correct weather associated emoji
		res := &WeatherHTTP{Weather: []wWeather{{Main: ""}}, Main: wMain{Temp: 0.0}}
		emoji := ""

		//checks the length of the user's command to format the proper api call
		//along with getting the matching weather emoji
		//reports to the user if an error is found
		switch len(args) {
		case 2:
			//only city was entered
			GetJson(url, &res)
			emoji = weatherEmojiFinder(res.Weather[0].Main)
			n.NootMessage(fmt.Sprintf("weather: %.2f°K and %s %s with %s", res.Main.Temp, res.Weather[0].Main, emoji, res.Weather[0].Description))

		case 3:
			//city and country is entered
			url = url + "&lang=" + args[2]
			GetJson(url, &res)
			emoji = weatherEmojiFinder(res.Weather[0].Main)
			n.NootMessage(fmt.Sprintf("weather: %.2f°K and %s %s with %s", res.Main.Temp, res.Weather[0].Main, emoji, res.Weather[0].Description))

		case 4:
			//city, country, and units are entered
			url = url + "&lang=" + args[2] + "&units=" + args[3]
			GetJson(url, &res)
			emoji = weatherEmojiFinder(res.Weather[0].Main)
			if args[3] == "imperial" {
				n.NootMessage(fmt.Sprintf("weather: %.2f℉ and %s %s with %s", res.Main.Temp, res.Weather[0].Main, emoji, res.Weather[0].Description))
			} else if args[3] == "metric" {
				n.NootMessage(fmt.Sprintf("weather: %.2f℃ and %s %s with %s", res.Main.Temp, res.Weather[0].Main, emoji, res.Weather[0].Description))
			} else {
				n.NootMessage(fmt.Sprintf("weather: %.2f°K and %s %s with %s", res.Main.Temp, res.Weather[0].Main, emoji, res.Weather[0].Description))
			}

		default:
			//error was encountered
			n.NootMessage("Please provide at least the [city] argument to run this command otherwise the format is [city] [country code] [units]")
		}

	}
}

//simple map function that matches the weather api json:main value to the apporiate weather emoji
func weatherEmojiFinder(weatherType string) string {

	//creates map for weather emojis
	weatherEmojiMapper := make(map[string]string)
	weatherEmojiMapper["Rain"] = ":cloud_rain:"
	weatherEmojiMapper["Snow"] = ":cloud_snow:"
	weatherEmojiMapper["Thunderstorm"] = ":thunder_cloud_rain:"
	weatherEmojiMapper["Clear"] = ":sunny:"
	weatherEmojiMapper["Fog"] = ":fog:"
	weatherEmojiMapper["Clouds"] = ":cloud:"

	//takes the result of the api json main and matches with apporiate weather emoji
	emojiFound := weatherEmojiMapper[weatherType]

	return emojiFound
}