package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Response struct {
	Books []struct {
		ID   string   `json:"id"`
		Tags []string `json:"tags"`
	} `json:"books"`
}

func GetBookInfo(keyword string) (string, []string, error) {
	url := fmt.Sprintf("https://booklog.jp/users/%s/all?category_id=all&status=all&sort=sort_desc&rank=all&tag=&page=1&keyword=%s&reviewed=&quoted=&json=true", config.Username, keyword)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Cookie", fmt.Sprintf("bid=%s%%3D", config.Cookie))
	req.Header.Set("Referer", fmt.Sprintf("https://booklog.jp/users/%s", config.Username))
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("error reading response: %v", err)
	}

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	if len(data.Books) > 0 {
		return data.Books[0].ID, data.Books[0].Tags, nil
	}

	return "", nil, fmt.Errorf("book not found")
}

func UpdateItemTag(itemId string, tags []string) error {
	apiURL := "https://booklog.jp/api/book/tag"

	tagString := strings.Join(tags, " ")

	data := url.Values{}
	data.Set("service_id", "1")
	data.Set("id", itemId)
	data.Set("tags", tagString)

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", fmt.Sprintf("bid=%s%%3D", config.Cookie))
	req.Header.Set("Referer", fmt.Sprintf("https://booklog.jp/users/%s", config.Username))
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func AddTagToItemByItemId(itemId, newTag string) error {
	// Get the current tags from the API
	_, currentTags, err := GetBookInfo(itemId)
	if err != nil {
		return fmt.Errorf("error getting current tags: %v", err)
	}

	// Check if the newTag is already in the currentTags
	for _, tag := range currentTags {
		if tag == newTag {
			fmt.Printf("Tag already exists: %v\n", itemId)
			return nil
		}
	}

	// Add the newTag to the currentTags
	newTags := append(currentTags, newTag)
	err = UpdateItemTag(itemId, newTags)
	if err != nil {
		return fmt.Errorf("error updating tags: %v", err)
	}

	fmt.Printf("Tag added successfully: %v\n", itemId)
	return nil

}

func AddTagToItemByIsbn(isbn, newTag string) error {
	itemId, _, err := GetBookInfo(isbn)
	if err != nil {
		return fmt.Errorf("error getting book info: %v", err)
	}

	err = AddTagToItemByItemId(itemId, newTag)
	if err != nil {
		return fmt.Errorf("error updating tags: %v", err)
	}

	return nil
}

func UpdateItemLocationByItemId(itemId, newLocation string) error {
	// Get the current tags from the API
	_, currentTags, err := GetBookInfo(itemId)
	if err != nil {
		return fmt.Errorf("error getting current tags: %v", err)
	}

	// Check if currentTags contain the newLocation tag or any tag starting with "loc_"
	found := false
	locationTagExists := false
	newTags := make([]string, 0)

	for _, tag := range currentTags {
		if tag == newLocation {
			locationTagExists = true
		} else if strings.HasPrefix(tag, "loc_") {
			found = true
		} else {
			newTags = append(newTags, tag)
		}
	}

	if locationTagExists {
		fmt.Println("Location tag already added")
		return nil
	}

	if found {
		// Replace the existing location tag with newLocation
		newTags = append(newTags, newLocation)
	} else {
		// Add newLocation to newTags
		newTags = append(currentTags, newLocation)
	}

	// Update the item's tags with newTags
	err = UpdateItemTag(itemId, newTags)
	if err != nil {
		return fmt.Errorf("error updating tags: %v", err)
	}

	fmt.Printf("Location updated successfully: %v\n", itemId)

	return nil
}

func UpdateItemLocationByIsbn(isbn, newLocation string) error {
	itemId, _, err := GetBookInfo(isbn)
	if err != nil {
		return fmt.Errorf("error getting book info: %v", err)
	}

	err = UpdateItemLocationByItemId(itemId, newLocation)
	if err != nil {
		return fmt.Errorf("error updating location: %v", err)
	}

	return nil
}
