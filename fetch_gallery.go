package main

import (
	   "bufio"
	   "fmt"
	   "http"
	   "io"
	   "io/ioutil"
	   "log"
	   "os"
	   "strings"
	   "xml"
)

// Container for photo metadata
type Photo struct {
	 Credit				string
	 Description		string
	 Link				string
}

const guardian_url = "http://www.guardian.co.uk/news/series/24hoursinpictures/rss"

func main() {
	var r *strings.Reader
	if len(os.Args) > 2 {
		fmt.Println("Usage: ", os.Args[0], "[<XML file>]")
		os.Exit(1)
	} else if len(os.Args) == 2 {
	    file := os.Args[1]
	    bytes, err := ioutil.ReadFile(file)
	    checkError(err)
	    r = strings.NewReader(string(bytes))
	} else {
		feed_data := fetchXMLFeed(guardian_url)
		r = strings.NewReader(feed_data)
    }

	parser := xml.NewParser(r)
	item_tags_seen := 0
	var photos []Photo
	var photo Photo
	for {
		// Only read the top (most recent) "item" in the feed
		if item_tags_seen == 2 {
		   break
		}
		token, err := parser.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
			case xml.StartElement:
			   elmt := xml.StartElement(t)
			   name := elmt.Name.Local

			switch name {
			   case "item":
			      item_tags_seen++
			   case "content":
			   	  if photo.Credit != "" && photo.Description != "" && photo.Link != "" {
			      	 fmt.Println("Appending photo")
			   	  	 photos = append(photos, photo)
			   	  } else {
			   	     photo = Photo{Link: getAttributeFromTag("url", elmt)}
			   	  }
			   case "credit":
					 photo.Credit = getTagContents(parser)
			   case "description":
					 photo.Description = getTagContents(parser)
			}
/*
  Unused token types:
		case xml.EndElement:
			continue
		case xml.CharData:
			bytes := xml.CharData(t)
			printElmt("\""+string([]byte(bytes))+"\"", depth)
		case xml.Comment:
			printElmt("Comment", depth)
		case xml.ProcInst:
			printElmt("ProcInst", depth)
		case xml.Directive:
			printElmt("Directive", depth)
*/
		}
	}
	fmt.Println("Number of photos: ", len(photos))
	for i, p := range(photos) {
		fmt.Println("Photo #", i, ":")
		fmt.Println("\tDesc: ", p.Description)
		fmt.Println("\tCred: ", p.Credit)
		fmt.Println("\tLink: ", p.Link)
	} 
}

func printElmt(s string, depth int) {
	for n := 0; n < depth; n++ {
		fmt.Print("  ")
	}
	fmt.Println(s)
}

func fetchXMLFeed(url string) string {
	client := http.Client{}
	response, err := client.Get(url)
	if err != nil {
	   return ""
	}
	defer response.Body.Close()
	return readAllContents(response.Body)
}

func readAllContents(rd io.Reader) string {
    reader := bufio.NewReader(rd)
    buffer := make([]byte, 4096)
    var content string
    for {
        n, err := reader.Read(buffer)
        if n > 0 {
            content += string(buffer[:n])
        }
        if err == os.EOF {
            break
        } else if err != nil {
            log.Fatal(err)
            break
        }
    }
    return content;
}

func getAttributeFromTag(attrib string, el xml.StartElement) string {
	 for _, attr := range el.Attr {
	 	 if attr.Name.Local == attrib {
		 	return attr.Value
		 }
     }
	 return ""
}

func getTagContents(parser *xml.Parser) string {
	 tag, err := parser.Token()
	 if err != nil { return "" }
	 switch dtype := tag.(type) {
	    case xml.CharData:
		    bytes := xml.CharData(dtype)
			text := string([]byte(bytes))
			return text
     }
	 return ""
}

func checkError(err os.Error) {
	if err != nil {
		fmt.Println("Fatal error ", err.String())
		os.Exit(1)
	}
}