package main

import (
	"fmt"
	"bufio"
	"os"
	"regexp"
	"strconv"
	"log"
	"strings"
	gt "github.com/bas24/translategooglefree"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	_"path/filepath"
	_"io/ioutil"
)
	var wordcount int = 0;
	var wordList = []string{};
	var db *sql.DB;

	func initDB(){
		os.MkdirAll("data", 0755);
		os.Create("./data/collection.anki2.sqlite3");
		db, _ = sql.Open("sqlite3", "./data/collection.anki2.sqlite3");

		db.Exec("CREATE TABLE revlog( id integer primary key, cid integer not null, usn integer not null, ease integer not null, ivl integer not null, lastIvl integer not null, factor integer not null, time integer not null, type integer not null);");
		db.Exec("CREATE TABLE notes ( id integer primary key, /* 0 */ guid text not null, /* 1 */ mid integer not null, /* 2 */ mod integer not null, /* 3 */ usn integer not null, /* 4 */ tags text not null, /* 5 */ flds text not null, /* 6 */ sfld integer not null, /* 7 */ csum integer not null, /* 8 */ flags integer not null, /* 9 */ data text not null /* 10 */ );");
		db.Exec("CREATE TABLE graves ( usn integer not null, oid integer not null, type integer not null );");
		db.Exec("CREATE TABLE col ( id integer primary key, crt integer not null, mod integer not null, scm integer not null, ver integer not null, dty integer not null, usn integer not null, ls integer not null, conf text not null, models text not null, decks text not null, dconf text not null, tags text not null );");
		db.Exec("CREATE TABLE cards ( id integer primary key, /* 0 */ nid integer not null, /* 1 */ did integer not null, /* 2 */ ord integer not null, /* 3 */ mod integer not null, /* 4 */ usn integer not null, /* 5 */ type integer not null, /* 6 */ queue integer not null, /* 7 */ due integer not null, /* 8 */ ivl integer not null, /* 9 */ factor integer not null, /* 10 */ reps integer not null, /* 11 */ lapses integer not null, /* 12 */ left integer not null, /* 13 */ odue integer not null, /* 14 */ odid integer not null, /* 15 */ flags integer not null, /* 16 */ data text not null /* 17 */ );");

		db.Exec("CREATE INDEX ix_revlog_usn on revlog (usn);");
		db.Exec("CREATE INDEX ix_revlog_cid on revlog (cid);");
		db.Exec("CREATE INDEX ix_notes_usn on notes (usn);");
		db.Exec("CREATE INDEX ix_notes_csum on notes (csum);");
		db.Exec("CREATE INDEX ix_cards_usn on cards (usn);");
		db.Exec("CREATE INDEX ix_cards_sched on cards (did, queue, due);");
		db.Exec("CREATE INDEX ix_cards_nid on cards (nid);");

		db.Close();
	}

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

			res, _ := http.Get("https://www.iconfinder.com/search/?q="+newWord+"&style=flat");
			defer res.Body.Close()
			if res.StatusCode != 200 {
				log.Fatalf("Status code error: %d %s", res.StatusCode, res.Status)
			}

			doc, _ := goquery.NewDocumentFromReader(res.Body)

			isFind := doc.Find("img.mr-4.float-left")
			if(len(isFind.Nodes)!=0){
				fmt.Printf("No images.\n");
			}else{
				imgSrc, isHere := doc.Find("img.d-block").First().Attr("src");
				if isHere {
					imgSrc = strings.Replace(imgSrc, "-128.png", "-512.png", 1)
					fmt.Printf(imgSrc)
					//save
				}
			}
		}

	}

	func main(){
		initDB();
		if false{
		//Input file
		const inputFile string = "in/The.Man.In.The.High.Castle.S03E01.1080p.WEB.h264-SKGTV-HI.srt";
		//Test translate
		const text string = `Hello, World!`;
		result, _ := gt.Translate(text, "en", "ru");
		fmt.Println(result);

		file, _ := os.Open(inputFile)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		var lineStatus int = 0;

		//ioutil.TempDir("card", "")

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

		addWord("Shit")
	}
