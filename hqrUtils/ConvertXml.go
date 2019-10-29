package hqrUtils

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

//Store data elements in the xml file to this struct
type Submission struct {
	XMLName       xml.Name `xml:"submission"`
	ActionCode    string   `xml:"action-code,attr"`
	Data          string   `xml:"data,attr"`
	Type          string   `xml:"type,attr"`
	Version       string   `xml:"version,attr"`
	FileAuditData struct {
		CreateDate   string `xml:"create-date"`
		CreateTime   string `xml:"create-time"`
		CreateBy     string `xml:"create-by"`
		Version      string `xml:"version"`
		CreateByTool string `xml:"create-by-tool"`
	} `xml:"file-audit-data"`
	AbstractionAuditData struct {
		AbstractionDate      string `xml:"abstraction-date"`
		AbstractorID         string `xml:"abstractor-id"`
		TotalAbstractionTime string `xml:"total-abstraction-time"`
		Comment              string `xml:"comment"`
	} `xml:"abstraction-audit-data"`
	Provider struct {
		ProviderID string `xml:"provider-id"`
		Patient    struct {
			FirstName  string `xml:"first-name"`
			LastName   string `xml:"last-name"`
			Birthdate  string `xml:"birthdate"`
			Sex        string `xml:"sex"`
			Race       string `xml:"race"`
			Ethnic     string `xml:"ethnic"`
			PostalCode string `xml:"postal-code"`
			Encounter  struct {
				MeasureSet    string `xml:"measure-set,attr"`
				EncounterDate string `xml:"encounter-date"`
				ArrivalTime   string `xml:"arrival-time"`
				PatientID     string `xml:"patient-id"`
				Detail        []struct {
					AnswerCode string `xml:"answer-code,attr"`
					QuestionCd string `xml:"question-cd,attr"`
					RowNumber  string `xml:"row-number,attr"`
				} `xml:"detail"`
			} `xml:"encounter"`
		} `xml:"patient"`
	} `xml:"provider"`
}

// let's declare a global Submission array
// we initialize our Submission subCase
var subCase Submission

func ConvertXml(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Inside function ConvertXml")

	fmt.Println("Endpoint Hit: ConvertXml")

	for key, val := range r.URL.Query() {
		//fmt.Fprintf(w, key, val)
		log.Println(key, val)
		fmt.Fprintf(w, val[0], "\n")

		// Open our xmlFile op-18-subCase.xml
		xmlFile, err := os.Open(val[0])
		// if we os.Open returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Successfully Opened xml file")
		// defer the closing of our xmlFile so that we can parse it later on
		defer xmlFile.Close()

		//the following code is used to handle US-ASCII format of xml
		decoder := xml.NewDecoder(xmlFile)
		decoder.CharsetReader = func(label string, input io.Reader) (io.Reader, error) {
			return input, nil
		}
		err = decoder.Decode(&subCase)
		if err != nil {
			fmt.Println(err)
		}

		// test only
		// we iterate through every Details within our Detail array and
		// print out the QuestionCd, AnswerCode, and their RowNumber
		for i := 0; i < len(subCase.Provider.Patient.Encounter.Detail); i++ {
			fmt.Println("QuestionCd: " + subCase.Provider.Patient.Encounter.Detail[i].QuestionCd)
			fmt.Println("AnswerCode: " + subCase.Provider.Patient.Encounter.Detail[i].AnswerCode)
			fmt.Println("RowNumber: " + subCase.Provider.Patient.Encounter.Detail[i].RowNumber)
		}

		//store the data in json format
		json, err := json.Marshal(&subCase)
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Write(json)
	}

}
