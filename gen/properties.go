package gen

import (
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (g *ApiGroup) parseProperty(raw []byte) Property {

	// log.Debugf("prop raw -> %s", raw)

	re1 := regexp.MustCompile(`[ \t]*?([A-Za-z0-9_}{]+) (.+?)( \[[^\]]*\])?\n((?s).+)`)
	matches2 := re1.FindAllSubmatch(raw, -1)

	// log.Debugf("m1 %s", matches2)

	if len(matches2[0]) == 1 {
		re1 = regexp.MustCompile(`[ \t]*?([A-Za-z0-9_}{]+) (.+?)\n((?s).+)`)
		matches2 = re1.FindAllSubmatch(raw, -1)
		// log.Debugf("m2 %s", matches2)
	}

	// log.Debugf("m2 %s", matches2[0][3])

	flags := []Flag{}
	flagListRaw := string(matches2[0][3])
	flagList := strings.Split(strings.Trim(flagListRaw, "[] "), ",")

	for _, f := range flagList {

		var flag Flag
		switch f {
		case "readonly":
			{
				flag = FlagReadOnly
			}
		case "readwrite":
			{
				flag = FlagReadWrite
			}
		case "experimental":
			{
				flag = FlagExperimental
			}
		}

		if flag > 0 {
			flags = append(flags, flag)
		}

	}

	docs := string(matches2[0][4])
	docs = strings.Replace(docs, " \t\n", "", -1)
	docs = strings.Trim(docs, " \t\n")

	name := string(matches2[0][2])
	name = strings.Replace(name, " \t\n", "", -1)

	p := Property{
		Type:  string(matches2[0][1]),
		Name:  name,
		Flags: flags,
		Docs:  docs,
	}

	if g.debug {
		log.Debugf("\t - %s %s %s", p.Type, p.Name, strings.Trim(flagListRaw, " "))
	}
	return p
}

func (g *ApiGroup) parseProperties(raw []byte) []Property {

	props := make([]Property, 0)
	slices := make([][]byte, 0)

	re := regexp.MustCompile(`(?s)\nProperties(.+)\n\n?`)
	matches1 := re.FindSubmatch(raw)

	if len(matches1) == 0 {
		return props
	}

	for _, propsRaw := range matches1[1:] {

		// string Modalias [readonly, optional]
		re1 := regexp.MustCompile(`(?s)[ \t]*(bool|byte|string|uint16|uint16_t|uint32|dict|array.*?) ([A-Za-z0-9_]+?).*?\n`)
		matches2 := re1.FindAllSubmatchIndex(propsRaw, -1)

		// log.Debugf("1*** %d", matches2)

		// if len(matches2) == 0 {
		// re1 := regexp.MustCompile(`[ \t]*(bool|byte|string|uint|dict|array.*) ([A-Za-z0-9_]+?)( ?) *\n`)
		// matches2 := re1.FindAllSubmatchIndex(propsRaw, -1)
		// }

		// log.Debugf("2 *** %d", matches2)

		if len(matches2) == 1 {
			if len(propsRaw) > 0 {
				// log.Debugf("ADD single *** %s", propsRaw)
				slices = append(slices, propsRaw)
			}
		} else {
			prevPos := -1
			for i := 0; i < len(matches2); i++ {

				if prevPos == -1 {
					prevPos = matches2[i][0]
					continue
				}

				nextPos := matches2[i][0]
				propRaw := propsRaw[prevPos:nextPos]
				prevPos = nextPos

				if len(propRaw) > 0 {
					slices = append(slices, propRaw)
				}

				// keep the last one
				lastItem := len(matches2) - 1
				if i == lastItem {
					propRaw = propsRaw[matches2[lastItem][0]:]
					if len(propRaw) > 0 {
						slices = append(slices, propRaw)
					}
				}
			}
		}
	}

	if g.debug {
		log.Debug("\tProperties:")
	}
	for _, propRaw := range slices {
		prop := g.parseProperty(propRaw)
		props = append(props, prop)
	}

	return props
}