package line_utils

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"linebot/config"
	"linebot/models"
	racipe "linebot/recipe"
	"log"
	"strconv"
	"time"
)

const DateFormat = "2006-01-02"

func GenerateAddFoodConfirmationTemplate(foodName string) *linebot.BubbleContainer {
	container := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "冷蔵庫への追加",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeXl,
				},
				&linebot.TextComponent{
					Type:   "text",
					Text:   fmt.Sprintf("「%s」をLINE冷蔵庫に追加します。賞味期限または、消費期限を入力してください", foodName),
					Weight: "bold",
					Wrap:   true,
				},
				&linebot.ButtonComponent{
					Type: linebot.FlexComponentTypeButton,
					Action: linebot.NewDatetimePickerAction(
						"日付を入力してください",
						fmt.Sprintf("{\"foodName\": \"%s\",\"type\": \"%s\"}", foodName, "add"),
						"date",
						time.Now().Format(DateFormat),
						"2030-01-01",
						time.Now().Format(DateFormat),
					),
				},
			},
		},
	}
	return container
}

func GenerateListTemplate(foods []models.Food) *linebot.BubbleContainer {
	var containerContents []linebot.FlexComponent

	containerContents = append(containerContents,
		&linebot.TextComponent{
			Type: linebot.FlexComponentTypeText,
			Text: "食品をタップして詳細へ",
			Size: linebot.FlexTextSizeTypeLg,
		})

	containerContents = append(containerContents,
		&linebot.SeparatorComponent{
			Type: linebot.FlexComponentTypeSeparator,
		})

	for _, food := range foods {
		containerContents = append(containerContents,
			&linebot.BoxComponent{
				Type:   linebot.FlexComponentTypeBox,
				Layout: linebot.FlexBoxLayoutTypeHorizontal,
				Contents: []linebot.FlexComponent{
					&linebot.TextComponent{
						Type:    linebot.FlexComponentTypeText,
						Text:    food.Name,
						Weight:  linebot.FlexTextWeightTypeBold,
						Size:    linebot.FlexTextSizeTypeXl,
						Align:   linebot.FlexComponentAlignTypeStart,
						Gravity: linebot.FlexComponentGravityTypeCenter,
					},
					&linebot.TextComponent{
						Type:    linebot.FlexComponentTypeText,
						Text:    fmt.Sprintf("~%s", food.ExpirationDate.Format(DateFormat)),
						Align:   linebot.FlexComponentAlignTypeEnd,
						Gravity: linebot.FlexComponentGravityTypeCenter,
					},
				},
				Action: &linebot.PostbackAction{
					Data: fmt.Sprintf("{\"type\": \"%s\", \"foodId\": \"%d\"}", "detail", food.ID),
				},
			},
			&linebot.SeparatorComponent{
				Type: linebot.FlexComponentTypeSeparator,
			},
		)
	}

	container := &linebot.BubbleContainer{
		Type:      linebot.FlexContainerTypeBubble,
		Direction: linebot.FlexBubbleDirectionTypeLTR,
		Body: &linebot.BoxComponent{
			Type:     linebot.FlexComponentTypeBox,
			Layout:   linebot.FlexBoxLayoutTypeVertical,
			Contents: containerContents,
		},
	}
	return container
}

func GenerateDetailTemplate(food models.Food) *linebot.BubbleContainer {
	titleComponent := &linebot.BoxComponent{
		Type:   linebot.FlexComponentTypeBox,
		Layout: linebot.FlexBoxLayoutTypeHorizontal,
		Contents: []linebot.FlexComponent{
			&linebot.TextComponent{
				Type:   linebot.FlexComponentTypeText,
				Text:   food.Name,
				Weight: linebot.FlexTextWeightTypeBold,
				Size:   linebot.FlexTextSizeTypeXl,
				Align:  linebot.FlexComponentAlignTypeStart,
				Margin: linebot.FlexComponentMarginTypeSm,
			},
			&linebot.TextComponent{
				Type:   linebot.FlexComponentTypeText,
				Text:   fmt.Sprintf("~%s", food.ExpirationDate.Format(DateFormat)),
				Weight: linebot.FlexTextWeightTypeRegular,
				Align:  linebot.FlexComponentAlignTypeEnd,
				Margin: linebot.FlexComponentMarginTypeSm,
			},
		},
	}

	eatButtonComponent := &linebot.ButtonComponent{
		Type: linebot.FlexComponentTypeButton,
		Action: &linebot.PostbackAction{
			Data:  fmt.Sprintf("{\"type\": \"%s\", \"status\": \"%s\", \"foodId\": \"%d\"}", "eat", "1", food.ID),
			Label: "食べた！",
		},
	}

	discardButtonComponent := &linebot.ButtonComponent{
		Type: linebot.FlexComponentTypeButton,
		Action: &linebot.PostbackAction{
			Data:  fmt.Sprintf("{\"type\": \"%s\", \"status\": \"%s\",\"foodId\": \"%d\"}", "discarded", "2", food.ID),
			Label: "ダメにしてしまった、、、。",
		},
	}

	deleteButtonComponent := &linebot.ButtonComponent{
		Type: linebot.FlexComponentTypeButton,
		Action: &linebot.PostbackAction{
			Data:  fmt.Sprintf("{\"type\": \"%s\", \"foodId\":\"%d\"}", "delete", food.ID),
			Label: "削除",
		},
	}

	containerContents := []linebot.FlexComponent{
		titleComponent,
		eatButtonComponent,
		discardButtonComponent,
		deleteButtonComponent,
	}

	container := &linebot.BubbleContainer{
		Type:      linebot.FlexContainerTypeBubble,
		Direction: linebot.FlexBubbleDirectionTypeLTR,
		Body: &linebot.BoxComponent{
			Type:     linebot.FlexComponentTypeBox,
			Layout:   linebot.FlexBoxLayoutTypeVertical,
			Contents: containerContents,
		},
	}
	return container
}

func GenerateFoodsEatenRateTemplate(foodRate models.FoodRate) *linebot.BubbleContainer {
	titleComponent := &linebot.BoxComponent{
		Type:   linebot.FlexComponentTypeBox,
		Layout: linebot.FlexBoxLayoutTypeHorizontal,
		Contents: []linebot.FlexComponent{
			&linebot.TextComponent{
				Type:   linebot.FlexComponentTypeText,
				Text:   "食べた割合",
				Weight: linebot.FlexTextWeightTypeBold,
				Align:  linebot.FlexComponentAlignTypeStart,
			},
			&linebot.TextComponent{
				Type:   linebot.FlexComponentTypeText,
				Text:   "破棄した割合",
				Weight: linebot.FlexTextWeightTypeBold,
				Align:  linebot.FlexComponentAlignTypeEnd,
			},
		},
	}

	barHeight := "20px"

	rateBarComponent := &linebot.BoxComponent{
		Type:            linebot.FlexComponentTypeBox,
		Layout:          linebot.FlexBoxLayoutTypeVertical,
		BackgroundColor: "#63FF68FF",
		Height:          barHeight,
		Width:           "100%",
		Contents: []linebot.FlexComponent{
			&linebot.BoxComponent{
				Type:            linebot.FlexComponentTypeBox,
				Layout:          linebot.FlexBoxLayoutTypeVertical,
				BackgroundColor: "#00CC06FF",
				Width:           strconv.Itoa(foodRate.AteRate) + "%",
				Height:          barHeight,
				Contents: []linebot.FlexComponent{
					&linebot.FillerComponent{
						Type: linebot.FlexComponentTypeFiller,
					},
				},
			},
		},
	}

	rateTextComponent := &linebot.BoxComponent{
		Type:   linebot.FlexComponentTypeBox,
		Layout: linebot.FlexBoxLayoutTypeHorizontal,
		Contents: []linebot.FlexComponent{
			&linebot.TextComponent{
				Type:  linebot.FlexComponentTypeText,
				Text:  strconv.Itoa(foodRate.AteRate) + "%",
				Align: linebot.FlexComponentAlignTypeStart,
			},
			&linebot.TextComponent{
				Type:  linebot.FlexComponentTypeText,
				Text:  strconv.Itoa(100-foodRate.AteRate) + "%",
				Align: linebot.FlexComponentAlignTypeEnd,
			},
		},
	}

	containerContents := []linebot.FlexComponent{
		titleComponent,
		rateBarComponent,
		rateTextComponent,
	}
	container := &linebot.BubbleContainer{
		Type:      linebot.FlexContainerTypeBubble,
		Direction: linebot.FlexBubbleDirectionTypeLTR,
		Body: &linebot.BoxComponent{
			Type:     linebot.FlexComponentTypeBox,
			Layout:   linebot.FlexBoxLayoutTypeVertical,
			Contents: containerContents,
		},
	}
	return container
}

func GenerateRecipeTemplate(recipe racipe.Recipe) *linebot.BubbleContainer {

	imageComponet := &linebot.ImageComponent{
		Type:        linebot.FlexComponentTypeImage,
		URL:         recipe.FoodImageUrl,
		Size:        linebot.FlexImageSizeTypeFull,
		AspectRatio: "20:13",
		AspectMode:  linebot.FlexImageAspectModeTypeCover,
		Action: &linebot.URIAction{
			URI:   recipe.RecipeUrl,
			Label: "action",
		},
	}

	titleComponent := &linebot.BoxComponent{
		Type:   linebot.FlexComponentTypeBox,
		Layout: linebot.FlexBoxLayoutTypeVertical,
		Contents: []linebot.FlexComponent{
			&linebot.TextComponent{
				Type:   linebot.FlexComponentTypeText,
				Text:   "食べた割合",
				Weight: linebot.FlexTextWeightTypeBold,
				Align:  linebot.FlexComponentAlignTypeStart,
				Size:   linebot.FlexTextSizeTypeXl,
			},
		},
	}

	contentComponent := &linebot.BoxComponent{
		Type:    linebot.FlexComponentTypeBox,
		Layout:  linebot.FlexBoxLayoutTypeVertical,
		Spacing: linebot.FlexComponentSpacingTypeSm,
		Contents: []linebot.FlexComponent{
			// 調理時間コンポーネント
			&linebot.BoxComponent{
				Type:   linebot.FlexComponentTypeBox,
				Layout: linebot.FlexBoxLayoutTypeBaseline,
				Contents: []linebot.FlexComponent{
					&linebot.TextComponent{
						Type:   linebot.FlexComponentTypeText,
						Weight: linebot.FlexTextWeightTypeBold,
						Margin: linebot.FlexComponentMarginTypeSm,
						Text:   "調理時間",
						Align:  linebot.FlexComponentAlignTypeStart,
					},
					&linebot.TextComponent{
						Type:   linebot.FlexComponentTypeText,
						Margin: linebot.FlexComponentMarginTypeSm,
						Text:   recipe.RecipeIndication,
						Align:  linebot.FlexComponentAlignTypeEnd,
					},
				},
			},
			// 費用コンポーネント
			&linebot.BoxComponent{
				Type:   linebot.FlexComponentTypeBox,
				Layout: linebot.FlexBoxLayoutTypeBaseline,
				Contents: []linebot.FlexComponent{
					&linebot.TextComponent{
						Type:   linebot.FlexComponentTypeText,
						Weight: linebot.FlexTextWeightTypeBold,
						Margin: linebot.FlexComponentMarginTypeSm,
						Text:   "費用",
						Align:  linebot.FlexComponentAlignTypeStart,
					},
					&linebot.TextComponent{
						Type:   linebot.FlexComponentTypeText,
						Margin: linebot.FlexComponentMarginTypeSm,
						Text:   recipe.RecipeCost,
						Align:  linebot.FlexComponentAlignTypeEnd,
					},
				},
			},
			// 説明コンポーネント
			&linebot.TextComponent{
				Type:  linebot.FlexComponentTypeText,
				Text:  recipe.RecipeDescription,
				Size:  linebot.FlexTextSizeTypeXs,
				Color: "#AAAAAA",
				Wrap:  true,
			},
		},
	}

	footerComponent := &linebot.BoxComponent{
		Type:   linebot.FlexComponentTypeBox,
		Layout: linebot.FlexBoxLayoutTypeVertical,
		Contents: []linebot.FlexComponent{
			&linebot.SpacerComponent{
				Size: linebot.FlexSpacerSizeTypeXxl,
			},
			&linebot.ButtonComponent{
				Type: linebot.FlexComponentTypeButton,
				Action: &linebot.URIAction{
					URI:   recipe.RecipeUrl,
					Label: "レシピを見る",
				},
				Color: "#0048FFFF",
				Style: linebot.FlexButtonStyleTypePrimary,
			},
		},
	}

	container := &linebot.BubbleContainer{
		Type:      linebot.FlexContainerTypeBubble,
		Direction: linebot.FlexBubbleDirectionTypeLTR,
		Hero:      imageComponet,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				titleComponent,
				contentComponent,
			},
		},
		Footer: footerComponent,
	}
	return container
}

func GenerateRecipeCarousel(bubble []*linebot.BubbleContainer) *linebot.CarouselContainer {
	carouselContainer := &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: bubble,
	}
	return carouselContainer
}

var Bot *linebot.Client
var err error

func initBot() {
	Bot, err = linebot.New(
		config.Config.ChannelSecret,
		config.Config.ChannelToken,
	)

	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	initBot()
}
