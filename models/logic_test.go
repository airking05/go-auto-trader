package models

import (
	"testing"
)

var testClose = []float64{201.28, 197.64, 195.78, 198.22, 201.74, 200.12, 198.55, 197.99, 196.8, 195.0, 197.55, 197.97, 198.97, 201.93, 200.83, 201.3, 198.64, 196.09, 197.91, 195.42, 197.84, 200.7, 199.93, 201.95, 201.39, 200.49, 202.63, 202.75, 204.7, 205.54, 205.86, 205.88, 205.73, 206.97, 206.94, 207.53, 207.35, 207.11, 206.4, 207.7, 206.85, 205.98, 206.2, 203.3, 204.15, 200.84, 200.37, 202.91, 201.67, 204.36, 203.76, 206.2, 205.26, 207.08, 206.67, 205.51, 202.5, 202.02, 202.48, 204.95, 203.16, 202.44, 203.17, 204.54, 204.0, 204.68, 205.59, 206.71, 205.78, 206.17, 207.1, 207.04, 204.66, 206.52, 206.28, 207.29, 207.81, 208.3, 207.43, 208.09, 207.23, 205.16, 207.38, 207.97, 205.59, 204.74, 205.56, 208.27, 207.27, 206.65, 206.69, 208.85, 209.07, 209.72, 209.65, 209.51, 210.12, 209.62, 207.36, 209.33, 209.09, 207.79, 208.22, 208.01, 208.56, 206.8, 206.45, 205.18, 205.15, 207.62, 208.28, 206.68, 205.79, 206.92, 207.25, 209.41, 208.48, 209.55, 209.71, 208.18, 207.54, 207.5, 203.15, 203.57, 205.21, 205.03, 204.43, 205.71, 202.27, 202.63, 205.19, 207.44, 208.35, 208.28, 209.95, 210.12, 210.24, 209.42, 209.03, 207.86, 205.7, 204.5, 207.02, 208.44, 208.49, 208.17, 207.47, 207.06, 207.75, 206.05, 205.65, 208.24, 206.35, 206.61, 206.35, 207.1, 208.26, 207.66, 206.02, 201.71, 195.64, 187.4, 185.2, 192.31, 197.07, 197.08, 195.48, 189.65, 193.25, 193.39, 190.46, 195.25, 192.64, 193.68, 194.56, 193.84, 196.26, 197.97, 197.52, 194.29, 195.3, 192.75, 192.45, 191.76, 191.73, 186.9, 187.01, 190.5, 190.99, 193.85, 197.3, 196.62, 198.23, 200.02, 200.14, 200.33, 199.07, 198.11, 201.15, 202.07, 202.17, 201.91, 200.66, 204.05, 206.28, 205.78, 205.38, 207.71, 207.59, 206.7, 209.15, 209.75, 209.12, 208.91, 208.8, 206.85, 207.33, 206.51, 203.63, 201.34, 204.4, 204.25, 207.5, 207.32, 208.07, 207.83, 208.11, 208.08, 208.32, 207.46, 209.43, 207.3, 204.39, 208.38, 207.12, 205.73, 204.13, 204.65, 200.69, 201.7, 203.82, 206.8, 203.65, 200.02, 201.67, 203.5, 200.02, 199.68, 198.21, 195.4, 192.93, 185.87}

var testVolume = []float64{121465900, 169632600, 209151400, 125346700, 147217800, 158567300, 144396100, 214553300, 192991100, 176613900, 211879600, 130991100, 122942700, 174356000, 117516800, 92009700, 134044600, 168514300, 173585400, 197729700, 163107000, 124212900, 134306700, 97953200, 125672000, 87219000, 96164200, 91087800, 97545900, 93670400, 76968200, 80652900, 91462500, 140896400, 74411100, 72472300, 73061700, 72697900, 108076000, 87491400, 110325800, 114497200, 76873000, 188128000, 89818900, 157121300, 110145700, 93993500, 162410900, 136099200, 94510400, 228808500, 117917300, 177715100, 71784500, 77805300, 159521700, 153067200, 118939000, 96180400, 126768700, 137303600, 86900900, 114368200, 81236300, 89351900, 85548900, 72722900, 74436600, 75099900, 99529300, 68934900, 191113200, 92189500, 72559800, 78264600, 102585900, 61327400, 79358100, 86863500, 125684900, 161304900, 103399700, 70927200, 113326200, 135060200, 88244900, 155877300, 75708100, 119727600, 94667900, 95934000, 76510100, 74549700, 72114600, 76857500, 64764600, 57433500, 124308600, 93214000, 74974600, 124919600, 93338800, 91531000, 87820900, 151882800, 121704700, 89063300, 105034700, 134551300, 73876400, 135382400, 124384200, 85308200, 126708600, 165867900, 130478700, 70696000, 68476800, 92307300, 97107400, 104174800, 202621300, 182925100, 135979900, 104373700, 117975400, 173820200, 164020100, 144113100, 129456900, 106069400, 81709600, 97914100, 106683300, 89030000, 70446800, 77965000, 88667900, 90509100, 117755000, 132361100, 123544800, 105791300, 91304400, 103266900, 113965700, 81820800, 85786800, 116030800, 117858000, 80270700, 126081400, 172123700, 89383300, 72786500, 79072600, 71692700, 172946000, 194327900, 346588500, 507244300, 369833100, 339257000, 274143900, 160414400, 163298800, 256000400, 160269300, 152087800, 207081000, 116025700, 149347700, 158611100, 119691200, 79452000, 113806200, 99581600, 276046600, 223657500, 105726200, 153890900, 92790600, 159378800, 155054800, 178515900, 159045600, 163452000, 131079000, 211003300, 126320800, 110274500, 124307300, 153055200, 107069200, 56395600, 88038700, 99106200, 134142200, 109692900, 76523900, 78448500, 102038000, 174911700, 144442300, 69033000, 77905800, 135906700, 90525500, 131076900, 86270800, 95246100, 96224500, 78408700, 110471500, 131008700, 75874600, 67846000, 121315200, 153577100, 117645200, 121123700, 121342500, 88220500, 94011500, 64931200, 98874400, 51980100, 37317800, 112822700, 97858400, 108441300, 166224200, 192913900, 102027100, 103372400, 162401500, 116128900, 211173300, 182385200, 154069600, 197017000, 173092500, 251393500, 99094300, 111026200, 110987200, 48542200, 65899900, 92640700, 63317700, 114877900}

func TestSmaLineCross(t *testing.T) {
	positionType := Long
	longPeriod := 50
	shortPeriod := 15
	keepPeriod := 3

	smaLineCross := NewSmaLineCross(positionType, shortPeriod, longPeriod, keepPeriod)
	charts := make([]Chart, 0)

	for _, close := range testClose {
		charts = append(charts, Chart{
			Last: close,
		})
	}
	judge := smaLineCross.Execute(charts)
	if judge {
		t.FailNow()
	}
	judge = smaLineCross.Execute([]Chart{})
	if judge {
		t.FailNow()
	}
}

func TestRSIContrarian(t *testing.T) {
	positionType := Long
	period := 15

	param := 30.0
	rsicontrarian := NewRSIContrarian(positionType, period, param)
	charts := make([]Chart, 0)

	for _, close := range testClose {
		charts = append(charts, Chart{
			Last: close,
		})
	}
	judge := rsicontrarian.Execute(charts)
	if !judge {
		t.FailNow()
	}
	judge = rsicontrarian.Execute([]Chart{})
	if judge {
		t.FailNow()
	}
}
func TestRSIFollow(t *testing.T) {
	positionType := Long
	period := 15

	param := 20.0
	rsifollow := NewRSIFollow(positionType, period, param)
	charts := make([]Chart, 0)

	for _, close := range testClose {
		charts = append(charts, Chart{
			Last: close,
		})
	}
	judge := rsifollow.Execute(charts)
	if !judge {
		t.FailNow()
	}
	judge = rsifollow.Execute([]Chart{})
	if judge {
		t.FailNow()
	}
}
func TestWMADif(t *testing.T) {
	period := 30

	param := 0.03
	wmadif := NewWMADif(period, param)
	charts := make([]Chart, 0)
	for _, close := range testClose {
		charts = append(charts, Chart{
			Last: close,
		})
	}
	judge := wmadif.Execute(charts)
	if !judge {
		t.FailNow()
	}
	judge = wmadif.Execute([]Chart{})
	if judge {
		t.FailNow()
	}
}

func TestEMADif(t *testing.T) {
	period := 30

	param := 0.03
	emadif := NewEMADif(period, param)
	charts := make([]Chart, 0)

	for _, close := range testClose {
		charts = append(charts, Chart{
			Last: close,
		})
	}
	judge := emadif.Execute(charts)
	if !judge {
		t.FailNow()
	}
	judge = emadif.Execute([]Chart{})
	if judge {
		t.FailNow()
	}
}

func TestSMADif(t *testing.T) {
	period := 30
	param := 0.03
	smadif := NewSMADif(period, param)
	charts := make([]Chart, 0)

	for _, close := range testClose {
		charts = append(charts, Chart{
			Last: close,
		})
	}
	judge := smadif.Execute(charts)
	if !judge {
		t.FailNow()
	}
	judge = smadif.Execute([]Chart{})
	if judge {
		t.FailNow()
	}
}

func TestOBV(t *testing.T) {
	positionType := Long
	period := 30
	param := 0.03
	obv := NewOBV(positionType, period, param)
	charts := make([]Chart, 0)

	for n := range testClose {
		charts = append(charts, Chart{
			Last:   testClose[n],
			Volume: testVolume[n],
		})
	}
	judge := obv.Execute(charts)
	if !judge {
		t.FailNow()
	}
	judge = obv.Execute([]Chart{})
	if judge {
		t.FailNow()
	}
}

func TestMACD(t *testing.T) {
	positionType := Long
	period := 30
	param := 0.03
	goldencross := NewGoldenCross(positionType, period, param)
	charts := make([]Chart, 0)

	for _, close := range testClose {
		charts = append(charts, Chart{
			Last: close,
		})
	}
	judge := goldencross.Execute(charts)
	if judge {
		t.FailNow()
	}
	judge = goldencross.Execute([]Chart{})
	if judge {
		t.FailNow()
	}
}
