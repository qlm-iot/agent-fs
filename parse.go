package main

import (
	"os"
	"io/ioutil"
	"bufio"
	"strconv"
	"../qlm/df"
)

func getline(filename string) string {
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	file.Close()
	return line
}

func ParseFs(root string) (df.Objects, error) {
	var qlmobjects df.Objects

	entries, err := ioutil.ReadDir(root)
	if err != nil {
		return qlmobjects, err
	}

	for _, entry := range entries {
		if (entry.IsDir() && entry.Name() == "Objects") {
			qlmobjects = ParseObjects(root + "/" + entry.Name());
		}
	}

	return qlmobjects, err
}

func ParseObjects(path string) df.Objects {
	var qlmobjects df.Objects
	var fullpath string

	objects := make([]df.Object, 0)
	entries, _ := ioutil.ReadDir(path)

	for _, entry := range entries {
		fullpath = path + "/" + entry.Name()
		if entry.IsDir() {
			if ((len(entry.Name()) > 7) && (entry.Name()[:7] == "Object.")) {
				objects = append(objects, ParseObject(fullpath))
			}
		} else {
			if (entry.Name() == "XmlnsXsi") {
				qlmobjects.XmlnsXsi = getline(fullpath)
			} else if (entry.Name() == "NoNamespaceSchemaLocation") {
				qlmobjects.NoNamespaceSchemaLocation = getline(fullpath)
			} else if (entry.Name() == "Version") {
				qlmobjects.Version = getline(fullpath)
			}
		}
	}

	qlmobjects.Objects = objects

	return qlmobjects
}

func ParseObject(path string) df.Object {
	var qlmobject df.Object
	var fullpath string

	objects := make([]df.Object, 0)
	infoitems := make([]df.InfoItem, 0)
	entries, _ := ioutil.ReadDir(path)

	for _, entry := range entries {
		fullpath = path + "/" + entry.Name()
		if entry.IsDir() {
			if (entry.Name() == "Id") {
				qlmobject.Id = ParseQLMID(fullpath)
			} else if (entry.Name() == "Description") {
				qlmobject.Description = ParseDescription(fullpath)
			} else if ((len(entry.Name()) > 7) && (entry.Name()[:7] == "Object.")) {
				objects = append(objects, ParseObject(fullpath))
			} else if ((len(entry.Name()) > 9) && (entry.Name()[:9] == "InfoItem.")) {
				infoitems = append(infoitems, ParseInfoItem(fullpath))
			}
		} else {
			if entry.Name() == "Type" {
				qlmobject.Type = getline(fullpath)
			} else if entry.Name() == "Udef" {
				qlmobject.Udef = getline(fullpath)
			}
		}
	}

	qlmobject.InfoItems = infoitems
	qlmobject.Objects = objects

	return qlmobject
}

func ParseQLMID(path string) *df.QLMID {
	var qlmid df.QLMID
	var fullpath string

	entries, _ := ioutil.ReadDir(path)

	for _, entry := range entries {
		fullpath = path + "/" + entry.Name()
		if (entry.IsDir() == false) {
			if (entry.Name() == "IdType") {
				qlmid.IdType = getline(fullpath)
			} else if (entry.Name() == "TagType") {
				qlmid.TagType = getline(fullpath)
			} else if (entry.Name() == "StartDate") {
				qlmid.StartDate = getline(fullpath)
			} else if (entry.Name() == "EndDate") {
				qlmid.EndDate = getline(fullpath)
			} else if (entry.Name() == "Udef") {
				qlmid.Udef = getline(fullpath)
			} else if (entry.Name() == "Text") {
				qlmid.Text = getline(fullpath)
			}
		}
	}

	return &qlmid;
}

func ParseDescription(path string) *df.Description {
	var qlmdescription df.Description
	var fullpath string

	entries, _ := ioutil.ReadDir(path)

	for _, entry := range entries {
		fullpath = path + "/" + entry.Name()
		if (entry.IsDir() == false) {
			if (entry.Name() == "Lang") {
				qlmdescription.Lang = getline(fullpath)
			} else if (entry.Name() == "Udef") {
				qlmdescription.Udef = getline(fullpath)
			} else if (entry.Name() == "Text") {
				qlmdescription.Text = getline(fullpath)
			}
		}
	}

	return &qlmdescription;
}

func ParseInfoItem(path string) df.InfoItem {
	var qlminfoitem df.InfoItem
	var fullpath string

	values := make([]df.Value, 0)
	entries, _ := ioutil.ReadDir(path)

	for _, entry := range entries {
		fullpath = path + "/" + entry.Name()
		if entry.IsDir() {
			if (entry.Name() == "Description") {
				qlminfoitem.Description = ParseDescription(fullpath)
			} else if (entry.Name() == "MetaData") {
				qlminfoitem.MetaData = ParseMetaData(fullpath)
			} else if ((len(entry.Name()) > 6) && (entry.Name()[:6] == "Value.")) {
				values = append(values, ParseValue(fullpath))
			}
		} else {
			if entry.Name() == "Udef" {
				qlminfoitem.Udef = getline(fullpath)
			} else if entry.Name() == "Name" {
				qlminfoitem.Name = getline(fullpath)
			} else if entry.Name() == "OtherNames" {
				qlminfoitem.OtherNames = ParseOtherNames(fullpath)
			}
		}
	}

	qlminfoitem.Values = values

	return qlminfoitem
}

func ParseMetaData(path string) *df.MetaData {
	var qlmmetadata df.MetaData
	var fullpath string

	infoitems := make([]df.InfoItem, 0)
	entries, _ := ioutil.ReadDir(path)

	for _, entry := range entries {
		fullpath = path + "/" + entry.Name()
		if entry.IsDir() {
			if ((len(entry.Name()) > 9) && (entry.Name()[:9] == "InfoItem.")) {
				infoitems = append(infoitems, ParseInfoItem(fullpath))
			}
		}
	}

	qlmmetadata.InfoItems = infoitems

	return &qlmmetadata
}

func ParseValue(path string) df.Value {
	var qlmvalue df.Value
	var fullpath string

	entries, _ := ioutil.ReadDir(path)

	for _, entry := range entries {
		fullpath = path + "/" + entry.Name()
		if (entry.IsDir() == false) {
			if (entry.Name() == "Text") {
				qlmvalue.Text = getline(fullpath)
			} else if (entry.Name() == "Type") {
				qlmvalue.Type = getline(fullpath)
			} else if (entry.Name() == "DateTime") {
				qlmvalue.DateTime = getline(fullpath)
			} else if (entry.Name() == "UnixTime") {
				qlmvalue.UnixTime, _ = strconv.ParseInt(getline(fullpath), 10, 64)
			}
		}
	}

	return qlmvalue
}

func ParseOtherNames(path string) []string {
	var name string
	names := make([]string, 0)

	file, _ := os.Open(path)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name = scanner.Text()
		if (len(name) > 0) {
			names = append(names, name)
		}
	}
	file.Close()

	return names
}
