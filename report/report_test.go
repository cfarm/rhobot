package report

import (
	"reflect"
	"testing"

	log "github.com/Sirupsen/logrus"

	"github.com/cfpb/rhobot/config"
)

var conf *config.Config

func init() {
	conf = config.NewConfig()
	conf.SetLogLevel("info")
}

type SimpleRE struct {
	SimpleHeaders []string
}

func (sre SimpleRE) GetHeaders() []string {
	return sre.SimpleHeaders[0:]
}

func (sre SimpleRE) GetValue(key string) string {
	return "simple"
}

func TestJSONReport(t *testing.T) {
	var re Element
	var rs Set
	var jrr Runner
	var phr Handler

	re = SimpleRE{[]string{"Some", "Thing"}}
	jrr = JSONReportRunner{}
	phr = PrintHandler{}

	elements := []Element{re, re}
	metadata := map[string]interface{}{"test": "json"}
	rs = Set{elements, metadata}

	reader, err := jrr.ReportReader(rs)
	err = phr.HandleReport(reader)
	if err != nil {
		t.Fatalf("error writing report\n%s", err)
	}

}

func TestJSONReportToFile(t *testing.T) {
	var re Element
	var rs Set
	var jrr Runner
	var fhr Handler

	re = SimpleRE{[]string{"Some", "Thing"}}
	jrr = JSONReportRunner{}
	fhr = FileHandler{"something.json"}

	elements := []Element{re, re}
	metadata := map[string]interface{}{"test": "json"}
	rs = Set{elements, metadata}

	reader, err := jrr.ReportReader(rs)
	err = fhr.HandleReport(reader)
	if err != nil {
		t.Fatalf("error writing report\n%s", err)
	}

}

func TestPongo2Report(t *testing.T) {
	var re Element
	var rs Set
	var prr Runner
	var phr Handler

	re = SimpleRE{[]string{"Some", "Thing"}}
	prr = NewPongo2ReportRunnerFromString(TemplateSimple)
	phr = PrintHandler{}

	elements := []Element{re, re}
	metadata := map[string]interface{}{"test": "pongo2"}
	rs = Set{elements, metadata}

	reader, err := prr.ReportReader(rs)
	err = phr.HandleReport(reader)
	if err != nil {
		t.Fatalf("error writing report\n%s", err)
	}
}

func TestDistributionList(t *testing.T) {
	df, err := ReadDistributionFormatYAMLFromFile("distributionListTest.yml")
	if err != nil {
		t.Fatalf("Failed to read distribution format\n%s", err)
	}

	// using reflection can work for the general use case
	severityList := reflect.ValueOf(&df.Severity).Elem()
	severityType := severityList.Type()
	for i := 0; i < severityList.NumField(); i++ {
		f := severityList.Field(i)
		log.Debugf("%d: %s %s = %v\n", i,
			severityType.Field(i).Name, f.Type(), f.Interface())
	}

	df.Print()
}
