package main

// func main() {
// 	fName := os.Args[1]
// 	dat, err := os.Open(fName)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var stmts []scraper.Statement
// 	err = json.NewDecoder(dat).Decode(&stmts)
// 	if err != nil {
// 		panic(err)
// 	}

// 	out := make(map[string][]scraper.Statement)

// 	for sub, l := range bySubject(stmts) {
// 		out[sub.SubjectSlug] = l
// 	}

// 	noExt := strings.TrimSuffix(fName, filepath.Ext(fName))

// 	f, err := os.Create(noExt + "-sorted.json")
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer f.Close()

// 	err = json.NewEncoder(f).Encode(out)

// 	if err != nil {
// 		panic(err)
// 	}

// 	f.Sync()

// 	stat, err := f.Stat()
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(strconv.FormatInt(stat.Size(), 10))
// }
