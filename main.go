package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	b64 "encoding/base64"

	"github.com/dariubs/percent"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"gopkg.in/yaml.v2"
)

var (
	fileName, outputFile, logoFile                                      string
	companyName, companyAddress, companyIndex, companyPhone, companyWeb string
)

// Properties represents incoming data structure
type Properties struct {
	Properties []struct {
		PropertyName    string `yaml:"property_name"`
		PropertyRecords []struct {
			RecordName            string `yaml:"record_name"`
			RecordType            string `yaml:"record_type"`
			RecordValue           string `yaml:"record_value"`
			RecordValueIsAkamaiIP bool   `yaml:"record_value_is_akamai_ip,omitempty"`
		} `yaml:"property_records"`
	} `yaml:"properties"`
}

// Report represents data used to compose PDF
type Report struct {
	Total     int
	UnCovered int
	Percent   float64
	Table     [][]string
}

func getHeader() []string {
	return []string{"Property Name", "Domain", "DNS Record Type", "DNS Record Value", "Resolves to Akamai IP?"}
}

func readFile(fileName string) Properties {
	fmt.Println("Parsing YAML file: ", fileName)

	if fileName == "" {
		fmt.Println("Please provide yaml file by using -f option")
		os.Exit(1)
	}

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		os.Exit(1)
	}

	var properties Properties
	err = yaml.Unmarshal(yamlFile, &properties)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	return properties
}

func getChart(result float64) string {
	buffer := bytes.NewBuffer([]byte{})
	pie := chart.PieChart{
		Width:  200,
		Height: 200,
		Values: []chart.Value{
			{
				Value: result,
				Label: "NO",
				Style: chart.Style{
					FillColor:           drawing.ColorRed,
					TextHorizontalAlign: chart.TextHorizontalAlignCenter,
					TextVerticalAlign:   chart.TextVerticalAlignMiddle,
				},
			},
			{
				Value: 100 - result,
				Label: "YES",
				Style: chart.Style{
					FillColor:           drawing.ColorGreen,
					TextHorizontalAlign: chart.TextHorizontalAlignCenter,
					TextVerticalAlign:   chart.TextVerticalAlignMiddle,
				},
			},
		},
	}

	pie.Render(chart.PNG, buffer)

	return b64.StdEncoding.EncodeToString(buffer.Bytes())
}

func getContents(properties Properties) (result Report) {
	for _, property := range properties.Properties {
		for p, record := range property.PropertyRecords {
			result.Total++
			akamaiIP := "-"
			if record.RecordType == "A" {
				akamaiIP = strconv.FormatBool(record.RecordValueIsAkamaiIP)

				if !record.RecordValueIsAkamaiIP {
					result.UnCovered++
				}
			}
			if p == 0 {
				result.Table = append(result.Table, []string{property.PropertyName, record.RecordName, record.RecordType, record.RecordValue, akamaiIP})
			}
			result.Table = append(result.Table, []string{"", record.RecordName, record.RecordType, record.RecordValue, akamaiIP})
		}
	}

	result.Percent = percent.PercentOf(result.UnCovered, result.Total)

	return result
}

func getDarkGrayColor() color.Color {
	return color.Color{
		Red:   144,
		Green: 144,
		Blue:  144,
	}
}

func getGrayColor() color.Color {
	return color.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}

func main() {
	flag.StringVar(&fileName, "f", "", "YAML file to parse.")
	flag.StringVar(&outputFile, "o", "", "Output file name.")
	flag.StringVar(&companyName, "c", "", "Company name.")
	flag.StringVar(&companyPhone, "t", "", "Company phone.")
	flag.StringVar(&companyWeb, "s", "", "Company web site.")
	flag.StringVar(&logoFile, "l", "", "Company logo file.")
	flag.Parse()

	begin := time.Now()

	darkGrayColor := getDarkGrayColor()
	grayColor := getGrayColor()
	whiteColor := color.NewWhite()
	header := getHeader()
	properties := readFile(fileName)
	//percentage, total, unCovered := getOverview(properties)
	contents := getContents(properties)

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)
	//m.SetBorder(true)

	m.RegisterHeader(func() {
		m.Row(20, func() {
			m.Col(3, func() {
				_ = m.FileImage(logoFile, props.Rect{
					Center:  true,
					Percent: 80,
				})
			})

			m.ColSpace(6)

			m.Col(3, func() {
				m.Text(companyName, props.Text{
					Size:  8,
					Align: consts.Right,
				})
				m.Text(fmt.Sprintf("Tel: %s", companyPhone), props.Text{
					Top:   12,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Right,
				})
				m.Text(companyWeb, props.Text{
					Top:   15,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Right,
				})
			})
		})
	})

	m.RegisterFooter(func() {
		m.Row(20, func() {
			m.Col(12, func() {
				m.QrCode(fmt.Sprintf("https://%s", companyWeb), props.Rect{
					Top:     13,
					Center:  true,
					Percent: 60,
				})
				m.Text(companyWeb, props.Text{
					Top:   16,
					Style: consts.BoldItalic,
					Align: consts.Center,
					Size:  8,
				})
			})
		})
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Akamai Domain Coverage Report", props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})

	m.Row(50, func() {
		image := getChart(contents.Percent)
		m.Col(6, func() {
			_ = m.Base64Image(image, consts.Png, props.Rect{
				Center:  true,
				Percent: 80,
			})
		})
		m.Col(6, func() {
			m.Text(fmt.Sprintf("Percentage of domains not pointing to Akamai: %2.2f%%", contents.Percent), props.Text{
				Top:   13,
				Style: consts.Normal,
				Align: consts.Center,
			})
			m.Text(fmt.Sprintf("Total number of domains in Akamai: %d", contents.Total), props.Text{
				Top:   23,
				Style: consts.Normal,
				Align: consts.Center,
			})
			m.Text(fmt.Sprintf("Number of domains not pointing to Akamai: %d", contents.UnCovered), props.Text{
				Top:   33,
				Style: consts.Normal,
				Align: consts.Center,
			})
		})
	})

	m.SetBackgroundColor(darkGrayColor)

	m.Row(7, func() {
		m.Col(12, func() {
			m.Text("Properties", props.Text{
				Top:   1,
				Size:  10,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})

	m.SetBackgroundColor(whiteColor)

	m.TableList(header, contents.Table, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{3, 3, 1, 4, 1},
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{3, 3, 1, 4, 1},
		},
		Align:                consts.Center,
		AlternatedBackground: &grayColor,
		HeaderContentSpace:   1,
		Line:                 false,
	})

	m.Row(10, func() {})

	err := m.OutputFileAndClose(outputFile)
	if err != nil {
		fmt.Println("Could not save PDF:", err)
		os.Exit(1)
	}

	end := time.Now()
	fmt.Println("Generated in ", end.Sub(begin))
}
