//   Copyright 2019 Content Mine Ltd
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

package europmc

import (
	"testing"
)

func TestLoadNonExistentFile(t *testing.T) {
	_, err := LoadPaperXMLFromFile("this/file/hopefully/does/not/exist.xml")
	if err == nil {
		t.Fatalf("Non existent file didn't cause error")
	}
}

func TestLoadTestXML(t *testing.T) {

	paper, err := LoadPaperXMLFromFile("testdata/test1.xml")
	if err != nil {
		t.Fatalf("Failed to load test paper: %v", err)
	}

    // Test the helper functions
	if paper.Title() != "Test title" {
	    t.Errorf("Unexpected title '%s'", paper.Title())
	}
	if paper.JournalTitle() != "International Proceedings of Test Data" {
	    t.Errorf("Unexpected journal title '%s'", paper.JournalTitle())
	}
	author := paper.FirstAuthor()
	if author == nil {
	    t.Errorf("Got nil first author")
	} else {
	    if author.Surname != "Dales" || author.GivenNames != "Michael W." {
	        t.Errorf("Got unexpected first author %v", *author)
	    }
	}
}
