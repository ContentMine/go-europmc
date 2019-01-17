//   Copyright 2018 Content Mine Ltd
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

// Package provides definitions for loading and reading open access XML papers from EuroPMC.
package europmc

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Base URL for EuroPMC REST API
const EUROPMC_API_URL string = "https://www.ebi.ac.uk/europepmc/webservices/rest"

type ContributorName struct {
	Surname    string `xml:"surname"`
	GivenNames string `xml:"given-names"`
}

type Contributor struct {
	Name ContributorName `xml:"name"`
}

type ContribGroup struct {
	XMLName      xml.Name      `xml:"contrib-group"`
	Contributors []Contributor `xml:"contrib"`
}

type ArticleTitleGroup struct {
	XMLName          xml.Name `xml:"title-group"`
	ArticleTitle     string   `xml:"article-title"`
	AlternativeTitle string   `xml:"alt-title"`
}

type JournalTitleGroup struct {
	XMLName      xml.Name `xml:"journal-title-group"`
	JournalTitle string   `xml:"journal-title"`
}

type JournalMeta struct {
	XMLName    xml.Name          `xml:"journal-meta"`
	TitleGroup JournalTitleGroup `xml:"journal-title-group"`
}

type ArticleID struct {
	XMLName xml.Name `xml:"article-id"`
	Type    string   `xml:"pub-id-type,attr"`
	ID      string   `xml:",chardata"`
}

type License struct {
	XMLName xml.Name `xml:"license"`
	Link    string   `xml:"href,attr"`
	Text    string   `xml:"license-p"`
}

type Permissions struct {
	XMLName            xml.Name `xml:"permissions"`
	CopyrightStatement string   `xml:"copyright-statement"`
	CopyrightYear      string   `xml:"copyright-year"`
	License            License  `xml:"license"`
}

type KeywordGroup struct {
	XMLName  xml.Name `xml:"kwd-group"`
	Title    string   `xml:"title"`
	Language string   `xml:"lang,attr"`
	Keywords []string `xml:"kwd"`
}

type ArticleMeta struct {
	XMLName           xml.Name          `xml:"article-meta"`
	IDs               []ArticleID       `xml:"article-id"`
	TitleGroup        ArticleTitleGroup `xml:"title-group"`
	ContributorGroups []ContribGroup    `xml:"contrib-group"`
	Permissions       Permissions       `xml:"permissions"`
	KeywordGroup      KeywordGroup      `xml:"kwd-group"`
}

type Front struct {
	XMLName     xml.Name    `xml:"front"`
	JournalMeta JournalMeta `xml:"journal-meta"`
	ArticleMeta ArticleMeta `xml:"article-meta"`
}

type OpenXMLPaper struct {
	XMLName xml.Name `xml:"article"`
	Front   Front    `xml:"front"`
}

// Parsing

func LoadPaperXMLFromFile(path string) (OpenXMLPaper, error) {
	var paper OpenXMLPaper

	f, err := os.Open(path)
	if err != nil {
		return OpenXMLPaper{}, err
	}

	err = xml.NewDecoder(f).Decode(&paper)
	return paper, err
}

// Getting things from EuroPMC

func FullTextURL(pmcid string) string {
	return fmt.Sprintf("%s/PMC%s/fullTextXML", EUROPMC_API_URL, pmcid)
}

func SupplementaryFilesURL(pmcid string) string {
	return fmt.Sprintf("%s/PMC%s/supplementaryFiles", EUROPMC_API_URL, pmcid)
}

func FetchFullText(pmcid string) (*OpenXMLPaper, error) {

	if len(pmcid) == 0 {
		return nil, fmt.Errorf("Empty PMCID provided")
	}

	resp, resp_err := http.Get(FullTextURL(pmcid))
	if resp_err != nil {
		return nil, resp_err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Status code %d", resp.StatusCode)
		} else {
			return nil, fmt.Errorf("Status code %d: %s", resp.StatusCode, body)
		}
	}

	paper := OpenXMLPaper{}
	err := xml.NewDecoder(resp.Body).Decode(&paper)
	if err != nil {
		return nil, err
	}

	return &paper, nil
}

// Convenience functions

func (author *ContributorName) String() string {
	if author == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s %s", author.GivenNames, author.Surname)
}

func (paper OpenXMLPaper) Title() string {
	return paper.Front.ArticleMeta.TitleGroup.ArticleTitle
}

func (paper OpenXMLPaper) JournalTitle() string {
	return paper.Front.JournalMeta.TitleGroup.JournalTitle
}

func (paper OpenXMLPaper) FirstAuthor() *ContributorName {
	contrib_groups := paper.Front.ArticleMeta.ContributorGroups
	if len(contrib_groups) > 0 {
		author_list := contrib_groups[0].Contributors
		if len(author_list) > 0 {
			return &(author_list[0].Name)
		}
	}
	return nil
}

func (paper OpenXMLPaper) ArticleID(id_type string) *string {
	for _, id := range paper.Front.ArticleMeta.IDs {
		if id.Type == id_type {
			return &(id.ID)
		}
	}
	return nil
}

func (paper OpenXMLPaper) PMCID() *string {
	return paper.ArticleID("pmcid")
}

func (paper OpenXMLPaper) PMID() *string {
	return paper.ArticleID("pmid")
}

func (paper OpenXMLPaper) LicenseURL() string {
	return paper.Front.ArticleMeta.Permissions.License.Link
}

func (paper OpenXMLPaper) Keywords() []string {
	return paper.Front.ArticleMeta.KeywordGroup.Keywords
}
