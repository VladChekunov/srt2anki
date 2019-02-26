package main

import (
	"fmt"
	"bufio"
	"os"
	"regexp"
	"strconv"
	gt "github.com/bas24/translategooglefree"
	"github.com/carrot/go-pinterest"
	"github.com/carrot/go-pinterest/controllers"
	_"github.com/carrot/go-pinterest/models"
)
	var wordcount int = 0;
	var wordList = []string{};
	var pinterestAT = "Pinterest token";

	func addWord(newWord string){
		var isNewword bool = true;

		for _, word := range wordList{
			if(newWord==word){
				isNewword = false;
				break;
			}
		}
		if(isNewword==true){
			fmt.Printf(newWord," ");
			wordList = append(wordList, newWord)
			wordcount++;

			client := pinterest.NewClient().RegisterAccessToken(pinterestAT);
			pins, _, _ := client.Me.Search.Pins.Fetch(newWord,
			    &controllers.MeSearchPinsFetchOptionals{Limit: 1,},
			);
			pin, _ := client.Pins.Fetch((*pins)[0].Id)
			fmt.Printf("%v", pin.Image.Original.Url);
			//Need to save
		}

	}

	func main(){
		if false{
		//Input file
		const inputFile string = "The.Man.In.The.High.Castle.S03E01.1080p.WEB.h264-SKGTV-HI.srt";
		//Test translate
		const text string = `Hello, World!`;
		result, _ := gt.Translate(text, "en", "ru");
		fmt.Println(result);

		file, _ := os.Open(inputFile)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		var lineStatus int = 0;

   		isText, _ := regexp.Compile("^(([0-9]{1,2}):([0-9]{1,2}):([0-9]{1,2}),([0-9]{1,3}) --> ([0-9]{1,2}):([0-9]{1,2}):([0-9]{1,2}),([0-9]{1,3}))$");
		
		for scanner.Scan() {
			lineText := scanner.Text();

			if(lineStatus==1 && lineText!=""){
				//fmt.Println("+",lineText);
				for _, match := range regexp.MustCompile(`[a-zA-Z]{1,100}`).FindAllString(lineText, -1) {
					addWord(match)
				}
			}else{
				//fmt.Println("-",lineText);
			}

			if(lineText=="" && lineStatus==1){
				lineStatus=0;
			}
			if(isText.MatchString(lineText)){
				lineStatus=1;
			}
		}
		fmt.Printf(strconv.Itoa(wordcount));
		}

		addWord("Test")
	}
