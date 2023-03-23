package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"net/http"
	"time"
)

var recipes []Recipe
var chefs []Chef

func init() {
	recipes = make([]Recipe, 0)
	chefs = make([]Chef, 0)
}

type Chef struct {
	Id                string    `json:"id"`
	Name              string    `json:"name"`
	Country           string    `json:"country"`
	YearsOfExperience int       `json:"yearsOfExperience"`
	Recipes           []*Recipe `json:"recipes"`
}

type Recipe struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Keywords     []string  `json:"keywords"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
	ChefId       string    `json:"chefId"`
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("recipe-id")
	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].Id == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)

	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe deleted",
	})
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("recipe-id")

	var recipe Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].Id == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	// Find the chef for the recipe using the chef ID in the request body
	chefID := recipe.ChefId
	var chef *Chef
	for i := range chefs {
		if chefs[i].Id == chefID {
			chef = &chefs[i]
			break
		}
	}
	if chef == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Chef with ID %s not found", chefID),
		})
		return
	}

	// Remove the recipe from the old chef's list of recipes
	oldChefID := recipes[index].ChefId
	oldChef := getChefByID(oldChefID)
	if oldChef != nil {
		oldRecipes := oldChef.Recipes
		for i, r := range oldRecipes {
			if r.Id == id {
				oldChef.Recipes = append(oldRecipes[:i], oldRecipes[i+1:]...)
				break
			}
		}
	}

	// Add the recipe to the new chef's list of recipes
	chef.Recipes = append(chef.Recipes, &recipe)

	recipe.Id = id
	recipes[index] = recipe

	c.JSON(http.StatusOK, recipe)
}

func getChefByID(id string) *Chef {
	for i := range chefs {
		if chefs[i].Id == id {
			return &chefs[i]
		}
	}
	return nil
}

func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// function to create new recipe
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	recipe.Id = xid.New().String()
	recipe.PublishedAt = time.Now()

	// Find the chef for the recipe using the chef ID in the request body
	chefID := recipe.ChefId
	var chef *Chef
	for i := range chefs {
		if chefs[i].Id == chefID {
			chef = &chefs[i]
			break
		}
	}
	if chef == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Chef with ID %s not found", chefID),
		})
		return
	}

	// Add the recipe to the chef's list of recipes
	chef.Recipes = append(chef.Recipes, &recipe)

	// Append the recipe to the global list of recipes
	recipes = append(recipes, recipe)

	c.JSON(http.StatusOK, recipe)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("recipes/:recipe-id", UpdateRecipeHandler)
	router.DELETE("recipes/:recipe-id", DeleteRecipeHandler)
	err := router.Run()
	if err != nil {
		return
	}
}
