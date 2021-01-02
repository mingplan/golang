package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	accessToken = "xxxx" // 个人访问token
	projectID = 1111 // 共有产权项目
	subscriberID = 1111 // 个人订阅ID
  addressURL = "xxxxx" // 服务器地址
)

var houseEstate = flag.String("houseEstateName", "金隅凤栖家园小区619号地块", "地址(金隅凤栖家园小区618号地块/金隅凤栖家园小区619号地块/金隅凤栖家园小区624号地块/瑞泽家园小区)")
var group = flag.String("group", "", "京籍或者非京籍")
var block = flag.String("block", "", "地块")
var building = flag.String("building", "", "楼栋（1，2，3）")
var unit = flag.String("unit", "", "单元（1，2，3）")
var toward = flag.String("toward", "", "朝向")
var roomType = flag.String("roomType", "", "居室（一居，二居，三居）")
var roomTypeCode = flag.String("rootTypeCode", "", "户型")
var isSelected = flag.Int("isSelected", -1, "是否被选中，1：已被选，0：未被选")
var isValid = flag.Int("isValid", -1, "是否可选")

// 小区的详情， 使用了多少，总共多少
type AreaSummary struct {
	HouseEstateID   int    `json:"HouseEstateID"`
	HouseEstateName string `json:"HouseEstateName"`
	Total           int    `json:"Total"`
	Used            int    `json:"Used"`
}


type TotalInfo struct {
	RoomQuantityList []AreaSummary `json:"RoomQuantityList"`
	Code   int    `json:"Code"`
	ErrMsg string `json:"ErrMsg"`
}

// 本次工程的每个小区的情况
func totalDetail() TotalInfo {
	body := strings.NewReader(`ProjectId=76`)
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/api/GetRoomQuantityInfoByProject", addressURL), body)
	if err != nil {
		// handle err
	}
	req.Host = addressURL
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "HousingResourceNotarization/3.0 (iPhone; iOS 14.2; Scale/3.00)")
	req.Header.Set("Accept-Language", "zh-Hans-CN;q=1, en-CN;q=0.9")
	req.Header.Set("Access_token", accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var data TotalInfo
	doHttpRequest(req, &data)
	return data
}

func doHttpRequest(r *http.Request, data interface{}) error {
	response, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(data); err != nil {
		return err
	}
	return nil
}


type AreaBuildings struct {
	Buildings []string `json:"Buildings"`
	Code      int      `json:"Code"`
	ErrMsg    string   `json:"ErrMsg"`
}

func areaBuildings(houseEstateId int) AreaBuildings{
	s := fmt.Sprintf("HouseEstateId=%d&ProjectId=%d", houseEstateId, projectID)
	body := strings.NewReader(s)
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/api/GetHouseBuildings", addressURL), body)
	if err != nil {
		//xxx
	}
	req.Host = addressURL
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "HousingResourceNotarization/3.0 (iPhone; iOS 14.2; Scale/3.00)")
	req.Header.Set("Accept-Language", "zh-Hans-CN;q=1, en-CN;q=0.9")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var data AreaBuildings
	doHttpRequest(req, &data)
	return data
}


// 单个房间详情
type RoomInfo struct {
	HouseID             int         `json:"HouseID"`
	SerialNumber        int         `json:"SerialNumber"`
	HouseEstate         string      `json:"HouseEstate"`
	Group               string      `json:"Group"`
	Block               string      `json:"Block"`
	Building            string      `json:"Building"`
	Unit                string      `json:"Unit"`
	RoomNumber          string      `json:"RoomNumber"`
	Toward              string      `json:"Toward"`
	RoomType            string      `json:"RoomType"`
	RoomTypeCode        string      `json:"RoomTypeCode"`
	EstimateBuiltUpArea float64     `json:"EstimateBuiltUpArea"`
	EstimateLivingArea  float64     `json:"EstimateLivingArea"`
	AreaUnitPrice       float64     `json:"AreaUnitPrice"`
	TotalPrice          float64     `json:"TotalPrice"`
	SubscriberID        interface{} `json:"SubscriberID"`
	SubscriberName      interface{} `json:"SubscriberName"`
	IsSelected          int         `json:"IsSelected"`
	IsValid             int         `json:"IsValid"`
	IsAbandon           int         `json:"IsAbandon"`
	IsPreselect         int         `json:"IsPreselect"`
	OriginalHouseID     interface{} `json:"OriginalHouseID"`
}

type UnitRooms struct {
	Unit        string `json:"Unit"`
	HouseDetail []RoomInfo `json:"HouseDetail"`
}

type BuildingsRooms struct {
	Building string `json:"Building"`
	Houses   []UnitRooms `json:"Houses"`
	Code   int    `json:"Code"`
	ErrMsg string `json:"ErrMsg"`
}

func buildingsDetail(houseEstateId int, buildingNum string) BuildingsRooms {
	s := fmt.Sprintf("BuildingId=%s&HouseEstateId=%d&ProjectId=%d&SubscriberID=%d", buildingNum, houseEstateId, projectID, subscriberID)
	body := strings.NewReader(s)
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/api/GetHouseEstateDetailsWithPreSelect", addressURL), body)
	if err != nil {
		// handle err
	}
	req.Host = addressURL
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "HousingResourceNotarization/3.0 (iPhone; iOS 14.2; Scale/3.00)")
	req.Header.Set("Accept-Language", "zh-Hans-CN;q=1, en-CN;q=0.9")
	req.Header.Set("Access_token", accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var data BuildingsRooms
	doHttpRequest(req, &data)
	return data
}


func fetchData() []RoomInfo {
	datas := []RoomInfo{}
	totalData := totalDetail()
	for _, item := range totalData.RoomQuantityList {
		fmt.Println(item.HouseEstateID)
		areaData := areaBuildings(item.HouseEstateID)
		for _, buildingNum := range areaData.Buildings {
			buildingData := buildingsDetail(item.HouseEstateID, buildingNum)
			for _, unitsHouse := range buildingData.Houses {
				for _, room := range unitsHouse.HouseDetail {
					datas = append(datas, room)
				}
			}
		}
	}
	return datas
}

func writeData2file(fileName string, datas []RoomInfo) error {
	raw, err := json.Marshal(datas)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileName, raw, 0666)
	if err != nil {
		fmt.Println("write to file error", err.Error())
		return err
	}
	return nil
}

func readDataFromFile(fileName string) []RoomInfo {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return []RoomInfo{}
	}
	var r []RoomInfo
	if err := json.Unmarshal(data, &r); err != nil {
		return []RoomInfo{}
	}
	return r
}

type HouseEstate string
type Building string

func format(data []RoomInfo) map[HouseEstate]map[Building][]RoomInfo {
	exists := make(map[HouseEstate]map[Building][]RoomInfo)
	for _, item := range  data {
		if _, ok := exists[HouseEstate(item.HouseEstate)]; ok {
			cur := exists[HouseEstate(item.HouseEstate)]
			if _, ok := cur[Building(item.Building)]; ok {
				curList := cur[Building(item.Building)]
				curList = append(curList, item)
				cur[Building(item.Building)] = curList
			} else {
				curList := []RoomInfo{item}
				cur[Building(item.Building)] = curList
			}
			exists[HouseEstate(item.HouseEstate)] = cur
		} else {
			cur := map[Building][]RoomInfo{Building(item.Building):[]RoomInfo{}}
			exists[HouseEstate(item.HouseEstate)] = cur
			curList := cur[Building(item.Building)]
			curList = append(curList, item)
		}
	}
	return exists
}

func formatBuilding(data []RoomInfo, buildNum string, topLevel map[string]int) {
	units := map[string][]RoomInfo{}
	for _, item := range data {
		if _, ok := units[item.Unit]; ok {
			cur := units[item.Unit]
			cur = append(cur, item)
			units[item.Unit] = cur
		} else {
			cur := []RoomInfo{item}
			units[item.Unit] = cur
		}
	}
	unitNums := []string{}
	for k := range  units {
		unitNums = append(unitNums, k)
	}
	sort.Strings(unitNums)
	for _, unitName := range unitNums {
		fmt.Println(fmt.Sprintf("单元%s", unitName))
		unitTopLevel := fmt.Sprintf("%s-%s", buildNum, unitName)
		curTopLevel := topLevel[unitTopLevel]
		curUnit := units[unitName]
		level := 1
		for {
			s := levelItem(curUnit, level, level == curTopLevel)
			if len(s) == 0 {
				break
			}
			fmt.Println(s)
			level = level + 1
		}
	}
}

// 单元内level的数据展示
func levelItem(data []RoomInfo, level int, isTopLevel bool) string {
	rooms := map[string]RoomInfo{}
	for _, item := range data {
		rooms[item.RoomNumber] = item
	}
	i := 1
	s := ""
	for {
		roomNumber := fmt.Sprintf("%d0%d", level, i)
		if data, ok := rooms[roomNumber]; ok {
			if isTopLevel {
				s += "\033[31m"
			} else if data.IsSelected == 0 {
				s += "\033[32m"
			}
			s += fmt.Sprintf("\t%6d %4s %4s %1d %4s:%2s:%-6s", data.HouseID, data.RoomNumber, data.Group, data.IsSelected, data.RoomType, data.Toward, data.RoomTypeCode)
			if isTopLevel {
				s += "\033[0m"
			} else if data.IsSelected == 0 {
				s += "\033[0m"
			}

		} else if roomNumber == "101" {
			s += fmt.Sprintf("\t%s", strings.Repeat(" ", 35))
		} else if i > 1{
			break
		}
		i = i + 1
	}
	return s
}

func showAllDetail(datas []RoomInfo, topLevel map[HouseEstate]map[string]int) {
	formats := format(datas)
	for houseEstate, buildings := range formats {
		fmt.Println(houseEstate)
		buildingsNums := []string{}
		for item, _ := range buildings {
			buildingsNums = append(buildingsNums, string(item))
		}
		sort.Strings(buildingsNums)
		buildingsToplevels := topLevel[houseEstate]
		for _, buildingNum := range  buildingsNums {
			fmt.Println("#", buildingNum)
			cur := buildings[Building(buildingNum)]
			formatBuilding(cur, buildingNum, buildingsToplevels)
		}
	}
}

func SyncData() {
	timestamp := time.Now().Format(time.RFC3339)
	fileName := fmt.Sprintf("rooms-%s.txt", timestamp)
	data := fetchData()
	writeData2file(fileName, data)
}

func topLevel(data []RoomInfo) map[HouseEstate]map[string]int {
	x := make(map[HouseEstate]map[string]int)
	for _, item := range data {
		unitName := fmt.Sprintf("%s-%s", item.Building, item.Unit)
		level, err := strconv.Atoi(item.RoomNumber[0:len(item.RoomNumber)-2])
		if err != nil {
			fmt.Println("error ", err.Error())
		}
		if _, ok := x[HouseEstate(item.HouseEstate)]; ok {
			cur :=  x[HouseEstate(item.HouseEstate)]
			if _, ok := cur[unitName]; ok {
				if level > cur[unitName] {
					cur[unitName] = level
				}
			} else {
				cur[unitName] = level
			}
		} else {
			x[HouseEstate(item.HouseEstate)] = map[string]int{unitName:level}
		}
	}
	return x
}

type Filter struct {
	HouseEstate         *string      `json:"HouseEstate"`
	Group               *string      `json:"Group"`
	Block               *string      `json:"Block"`
	Building            *string      `json:"Building"`
	Unit                *string      `json:"Unit"`
	Toward              *string      `json:"Toward"`
	RoomType            *string      `json:"RoomType"`
	RoomTypeCode        *string      `json:"RoomTypeCode"`
	IsSelected          *int         `json:"IsSelected"`
	IsValid             *int         `json:"IsValid"`
	IsAbandon           *int         `json:"IsAbandon"`
	IsPreselect         *int         `json:"IsPreselect"`
}

func (f Filter) filter(data []RoomInfo) []RoomInfo {
	res := []RoomInfo{}
	for _, item := range data {
		if f.HouseEstate != nil && item.HouseEstate != *(f.HouseEstate) {
			continue
		}
		if f.Group != nil && item.Group != *(f.Group) {
			continue
		}
		if f.Block != nil && item.Block != *(f.Block) {
			continue
		}
		if f.Unit != nil && item.Unit != *(f.Unit) {
			continue
		}
		if f.Building != nil && item.Building != *(f.Building) {
			continue
		}
		if f.RoomType != nil && item.RoomType != *(f.RoomType) {
			continue
		}
		if f.RoomTypeCode != nil && item.RoomTypeCode != *(f.RoomTypeCode) {
			continue
		}
		if f.IsSelected != nil && item.IsSelected != *(f.IsSelected) {
			continue
		}
		if f.IsValid != nil && item.IsValid != *(f.IsValid) {
			continue
		}
		res = append(res, item)
	}
	return res
}

func (f Filter) str() string {
	data, _ := json.Marshal(f)
	return string(data)
}

func parseFilter() Filter {
	f := Filter{
		HouseEstate:  nil,
		Group:        nil,
		Block:        nil,
		Building:     nil,
		Unit:         nil,
		Toward:       nil,
		RoomType:     nil,
		RoomTypeCode: nil,
		IsSelected:   nil,
		IsValid:      nil,
		IsAbandon:    nil,
		IsPreselect:  nil,
	}
	if len(*houseEstate) > 0 {
		f.HouseEstate = houseEstate
	}
	if len(*group) > 0 {
		f.Group = group
	}
	if len(*block) > 0 {
		f.Block = block
	}
	if len(*building) > 0 {
		f.Building = building
	}
	if len(*unit) > 0 {
		f.Unit = unit
	}
	if len(*toward) > 0 {
		f.Toward = toward
	}
	if len(*roomType) > 0 {
		f.RoomType = roomType
	}
	if len(*roomTypeCode) > 0 {
		f.RoomTypeCode = roomTypeCode
	}
	if *isValid != -1 {
		f.IsValid = isValid
	}
	if *isSelected != -1 {
		f.IsSelected = isSelected
	}
	fmt.Println(f.str())
	return f
}


func summary(data []RoomInfo, group string) {
	// 京jiu的选房情况
	// 分楼盘的情况
	// 分户型的情况
	type Info struct {
		 Total int
		 Selected int
	}
	roomTypeMap := make(map[string]map[string]map[string]Info)
	for _, item := range data {
		if item.Group == group {
			if _, ok := roomTypeMap[item.HouseEstate]; !ok {
				roomTypeMap[item.HouseEstate] = make(map[string]map[string]Info)
			}
			curBuilding := roomTypeMap[item.HouseEstate]
			if _, ok := curBuilding[item.Building]; !ok {
				curBuilding[item.Building] = make(map[string]Info)
			}
			curRoomType := curBuilding[item.Building]
			if _, ok := curRoomType[item.RoomType]; !ok {
				curRoomType[item.RoomType] = Info{
					Total:    0,
					Selected: 0,
				}
			}
			curItem := curRoomType[item.RoomType]
			curItem.Total += 1
			curItem.Selected += item.IsSelected
			curRoomType[item.RoomType] = curItem
		}
	}

	for _, houseEs := range []string{"金隅凤栖家园小区619号地块", "金隅凤栖家园小区618号地块", "金隅凤栖家园小区624号地块", "瑞泽家园小区"} {
		if _, ok := roomTypeMap[houseEs]; !ok {
			continue
		}
		data := roomTypeMap[houseEs]
		fmt.Println(fmt.Sprintf("location: %s 分户型 %s情况", houseEs, group))
		hTotal := 0
		hSelected := 0
		hMap := map[string]Info{}
		for building, item := range data {
			selected := 0
			total := 0
			existed := map[string]string{}
			for roomType, inf := range item {
				existed[roomType] = fmt.Sprintf("%s: total = %d selected = %d", roomType, inf.Total, inf.Selected)
				total += inf.Total
				selected += inf.Selected
				if _, ok := hMap[roomType]; !ok {
					hMap[roomType] = Info{
						Total:    0,
						Selected: 0,
					}
				}
				cur := hMap[roomType]
				cur.Selected += inf.Selected
				cur.Total += inf.Total
				hMap[roomType] = cur
			}
			hTotal += total
			hSelected += selected
			fmt.Println(fmt.Sprintf("building #%s : total = %d selected = %d", building, total, selected))
			for _, roomType := range []string{"一居", "二居", "三居"} {
				if data, ok := existed[roomType]; ok {
					fmt.Println(data)
				}
			}
		}
		fmt.Println(fmt.Sprintf("location: %s 分户型 %s情况 total: %d, selected: %d", houseEs, group, hTotal, hSelected))
		for _, roomType := range []string{"一居", "二居", "三居"} {
			if _, ok := hMap[roomType]; ok {
				fmt.Println(fmt.Sprintf("location: %s %s情况 %s total: %d, selected: %d", houseEs, group, roomType, hMap[roomType].Total, hMap[roomType].Selected))
			}
		}
		fmt.Println()
	}
}

func main() {
	flag.Parse()
	f := parseFilter()
	data := fetchData()
	// timestamp := time.Now().Format(time.RFC3339)
	// fileName := fmt.Sprintf("rooms-%s.txt", timestamp)
	// writeData2file(fileName, data)
	summary(data, "非京籍")
	topLevel := topLevel(data)
	fmt.Println(fmt.Sprintf("total house amount %d", len(data)))
	cur := f.filter(data)
	showAllDetail(cur, topLevel)
	fmt.Println(fmt.Sprintf("%s filter: house amount %d", f.str(), len(cur)))
}

