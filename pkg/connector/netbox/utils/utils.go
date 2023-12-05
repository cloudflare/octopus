package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/cloudflare/octopus/pkg/model"
)

var (
	ifNameTagsRegexp = regexp.MustCompile(`^(?P<intf_name>.+?)\.((?P<outer_tag>\d+)\.)?(?P<inner_tag>\d+)$`)
)

func ParseUnitStr(unitStr string) (model.VLANTag, error) {
	if strings.Contains(unitStr, ".") {
		parts := strings.Split(unitStr, ".")
		if len(parts) != 2 {
			return model.VLANTag{}, fmt.Errorf("invalid unit string %q", unitStr)
		}

		outerTag, err := strconv.Atoi(parts[0])
		if err != nil {
			return model.VLANTag{}, fmt.Errorf("unable to convert %q to int: %v", parts[0], err)
		}

		innerTag, err := strconv.Atoi(parts[1])
		if err != nil {
			return model.VLANTag{}, fmt.Errorf("unable to convert %q to int: %v", parts[1], err)
		}

		return model.NewVLANTag(uint16(outerTag), uint16(innerTag)), nil
	}

	ctag, err := strconv.Atoi(unitStr)
	if err != nil {
		return model.VLANTag{}, fmt.Errorf("unable to convert %q to int: %v", unitStr, err)
	}

	return model.NewVLANTag(0, uint16(ctag)), nil
}

func SanitizeIPAddress(addr string) string {
	if strings.Contains(addr, "/") {
		return addr
	}

	if strings.Contains(addr, ".") {
		return addr + "/32"
	}

	return addr + "/128"
}

func GetMetaDataFromTags(tags []string) (*model.MetaData, error) {
	ret := model.NewMetaData()
	for _, tag := range tags {
		parts := strings.Split(tag, "=")

		// Semantic Tag
		if len(parts) == 2 {
			if _, exists := ret.SemanticTags[parts[0]]; exists {
				return nil, fmt.Errorf("key %q exists already: %q vs. %q", parts[0], ret.SemanticTags[parts[0]], parts[1])
			}

			ret.SemanticTags[parts[0]] = parts[1]

		} else {
			// Regular Tag
			ret.Tags = append(ret.Tags, tag)
		}
	}

	return ret, nil
}

func GetCustomFieldData(md *model.MetaData, customFieldData string) {
	if customFieldData == "" || customFieldData == "{}" {
		return
	}

	md.CustomFieldData = customFieldData
}

func GetInterfaceAndVLANTag(name string) (ifName string, vt model.VLANTag, err error) {
	if isLogicalInterface(name) {
		return extractInterfaceAndUnit(name)
	}

	return name, model.VLANTag{}, nil
}

func isLogicalInterface(name string) bool {
	return strings.Contains(name, ".")
}

func extractInterfaceAndUnit(name string) (string, model.VLANTag, error) {
	extractedVars := reSubMatchMap(ifNameTagsRegexp, name)
	if _, exists := extractedVars["intf_name"]; !exists {
		return "", model.VLANTag{}, fmt.Errorf("unable to extract interface name from %q", name)
	}

	if _, exists := extractedVars["inner_tag"]; !exists {
		return "", model.VLANTag{}, fmt.Errorf("unable to extract inner tag from %q", name)
	}

	innerTag, err := strconv.Atoi(extractedVars["inner_tag"])
	if err != nil {
		return "", model.VLANTag{}, fmt.Errorf("unable to convert inner tag from string %q to int (%q)", extractedVars["inner_tag"], name)
	}

	vt := model.VLANTag{
		InnerTag: uint16(innerTag),
	}

	outerTagStr := extractedVars["outer_tag"]
	if outerTagStr != "" {
		outerTag, err := strconv.Atoi(outerTagStr)
		if err != nil {
			return "", model.VLANTag{}, fmt.Errorf("unable to convert outer tag from string %q to int (%q)", outerTagStr, name)
		}

		vt.OuterTag = uint16(outerTag)
	}

	return extractedVars["intf_name"], vt, nil
}

func reSubMatchMap(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 && name != "" && i <= len(match) {
			subMatchMap[name] = match[i]
		}
	}

	return subMatchMap
}
