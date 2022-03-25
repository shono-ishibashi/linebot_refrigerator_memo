package recipe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"
)

const RakutenAppId = "1051904742240404263"

type CategoryRoot struct {
	Result CategoryResult `json:"result"`
}
type CategoryResult struct {
	Small []Category `json:"small"`
}

type Category struct {
	CategoryId   string `json:"-"`
	CategoryName string `json:"categoryName"`
	CategoryUrl  string `json:"categoryUrl"`
}

type RecipeRood struct {
	Result []Recipe `json:"result"`
}

type Recipe struct {
	RecipeId          int    `json:"recipeId"`
	RecipeTitle       string `json:"recipeTitle"`
	RecipeUrl         string `json:"recipeUrl"`
	FoodImageUrl      string `json:"foodImageUrl"`
	RecipeIndication  string `json:"recipeIndication"`
	RecipeCost        string `json:"recipeCost"`
	Rank              string `json:"rank"`
	RecipeDescription string `json:"recipeDescription"`
}

func FetchCategoryList() ([]Category, error) {
	apiUrl := "https://app.rakuten.co.jp/services/api/Recipe/CategoryList/20170426?format=json&categoryType=small&applicationId=" + RakutenAppId

	res, requestGetErr := http.Get(apiUrl)
	defer res.Body.Close()

	if requestGetErr != nil {
		return nil, requestGetErr
	}

	body, ioutilReadAllErr := ioutil.ReadAll(res.Body)
	if ioutilReadAllErr != nil {
		return nil, ioutilReadAllErr
	}

	var resBody CategoryRoot
	jsonUnmarshalErr := json.Unmarshal(body, &resBody)
	if jsonUnmarshalErr != nil {
		return nil, jsonUnmarshalErr
	}

	categoryList := resBody.Result.Small

	for index, category := range resBody.Result.Small {
		categoryUrlRemovedDomain := strings.Replace(category.CategoryUrl, "https://recipe.rakuten.co.jp/category/", "", 1)
		categoryList[index].CategoryId = categoryUrlRemovedDomain[:utf8.RuneCount([]byte(categoryUrlRemovedDomain))-1]
	}
	return categoryList, nil
}

func SearchCategoryByFoodName(query string, categoryList []Category) []Category {
	var searchedCategoryList []Category
	for _, category := range categoryList {
		if strings.Contains(category.CategoryName, query) {
			searchedCategoryList = append(searchedCategoryList, category)
		}
	}
	return searchedCategoryList
}

func SearchRecipeByCategoryId(categoryId string) ([]Recipe, error) {
	apiUrlTemplate := "https://app.rakuten.co.jp/services/api/Recipe/CategoryRanking/20170426?format=json&categoryId=%s&applicationId=%s"
	apiUrl := fmt.Sprintf(apiUrlTemplate, categoryId, RakutenAppId)
	res, requestGetErr := http.Get(apiUrl)
	defer res.Body.Close()
	if requestGetErr != nil {
		return nil, requestGetErr
	}

	body, ioutilReadAllErr := ioutil.ReadAll(res.Body)
	if ioutilReadAllErr != nil {
		return nil, ioutilReadAllErr
	}

	var resBody RecipeRood
	jsonUnmarshalErr := json.Unmarshal(body, &resBody)
	if jsonUnmarshalErr != nil {
		return nil, jsonUnmarshalErr
	}
	fmt.Println(resBody.Result)
	return resBody.Result, nil
}

func main() {
	categoryList, fetchCategoryListErr := FetchCategoryList()
	if fetchCategoryListErr != nil {
		log.Fatalln(fetchCategoryListErr)
	}

	searchedCategoryList := SearchCategoryByFoodName("ポテト", categoryList)
	if len(searchedCategoryList) == 0 {
		log.Println("no result")
	}

	for _, searchedCategory := range searchedCategoryList {
		recipeList, searchRecipeErr := SearchRecipeByCategoryId(searchedCategory.CategoryId)
		if searchRecipeErr != nil {
			log.Fatalln(searchRecipeErr)
		}
		if len(recipeList) != 0 {
			for _, v := range recipeList {
				fmt.Println(v)
				fmt.Println(v.RecipeDescription)
			}
		}
	}
}
