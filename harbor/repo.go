package harbor

//
//url := fmt.Sprintf("https://your.harbor.domain/api/v2.0/projects/%s/repositories", projectName)
//
//req, err := http.NewRequest("GET", url, nil)
//if err != nil {
//fmt.Println("Error creating request:", err)
//return
//}
//
//// Set authentication header if needed
//req.SetBasicAuth("your_username", "your_password")
//
//client := &http.Client{}
//resp, err := client.Do(req)
//if err != nil {
//fmt.Println("Error making request:", err)
//return
//}
//defer resp.Body.Close()
//
//if resp.StatusCode != http.StatusOK {
//fmt.Println("Error response status:", resp.Status)
//return
//}
//
//body, err := ioutil.ReadAll(resp.Body)
//if err != nil {
//fmt.Println("Error reading response body:", err)
//return
//}
//
//var repositories []Repository
//if err := json.Unmarshal(body, &repositories); err != nil {
//fmt.Println("Error unmarshalling JSON:", err)
//return
//}
//
//for _, repo := range repositories {
//fmt.Println("Repository:", repo.Name)
//}
