package main

import (
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type poolItem struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Used      bool   `json:"used"`
}

type slotItem struct {
	Slot   uint       `json:"slot"`
	Item   poolItem   `json:"item"`
	Player playerItem `json:"player"`
}

type playerItem struct {
	Playing bool    `json:"playing"`
	Volume  float64 `json:"volume"`
	Loop    bool    `json:"loop"`
	Stop    bool    `json:"stop"` //only used in PlayerInput type
}

func initGraphql() {
	poolType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "PoolItem",
			Description: `A PoolItem contains information about one item in the pool of files.`,
			Fields: graphql.Fields{
				"name": &graphql.Field{
					Type:        graphql.String,
					Description: `The name of the pool item. Is the same as the filename of the uploaded file without extension`,
				},
				"extension": &graphql.Field{
					Type:        graphql.String,
					Description: `The file type of the pool item. Can be "wav" or "mp3". Nothing more is supported at the moment.`,
				},
				"used": &graphql.Field{
					Type:        graphql.Boolean,
					Description: `This is true if the item is in use in one or more slots. If it is used, it can not be deleted.`,
				},
			},
		},
	)

	playerType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Player",
			Description: `Player stores the state of a Playback on a slot.`,
			Fields: graphql.Fields{
				"playing": &graphql.Field{
					Type:        graphql.Boolean,
					Description: `This is true if the playback is currently running.`,
				},
				"volume": &graphql.Field{
					Type: graphql.Float,
					Description: `The volume value as decimal number. 0 means no change to the original. 
					Positive values will increase the loudness, but increase the risk for clipping!
					Typical values are +3. Negative Values will decrease the loudness. The rating of this
					value will may change in the future to a dB scale.`,
				},
				"loop": &graphql.Field{
					Type:        graphql.Boolean,
					Description: `If this is true, the playback will be looped forever.`,
				},
			},
		},
	)

	playerInputType := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name:        "PlayerInput",
			Description: `The Player input type is used to change the playing state of a slot.`,
			Fields: graphql.InputObjectConfigFieldMap{
				"playing": &graphql.InputObjectFieldConfig{
					Type:        graphql.Boolean,
					Description: `This is true if the playback is currently running.`,
				},
				"volume": &graphql.InputObjectFieldConfig{
					Type: graphql.Float,
					Description: `The volume value as decimal number. 0 means no change to the original. 
					Positive values will increase the loudness, but increase the risk for clipping!
					Typical values are +3. Negative Values will decrease the loudness. The rating of this
					value will may change in the future to a dB scale.`,
				},
				"loop": &graphql.InputObjectFieldConfig{
					Type:        graphql.Boolean,
					Description: `If this is true, the playback will be looped forever.`,
				},
				"stop": &graphql.InputObjectFieldConfig{
					Type: graphql.Boolean,
					Description: `This is the field with the highest priority. If this is true, the playback will
					bew rewinded at the begin of the execution of the other values.`,
				},
			},
		},
	)

	slotType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "SlotItem",
			Description: `A SlotItem contains the information about one slot.`,
			Fields: graphql.Fields{
				"slot": &graphql.Field{
					Type:        graphql.Int,
					Description: `The slot number to which the PoolItem belongs.`,
				},
				"item": &graphql.Field{
					Type:        poolType,
					Description: `The PoolItem that is currently mapped to the slot number.`,
				},
				"player": &graphql.Field{
					Type:        playerType,
					Description: `The player information for the slot.`,
				},
			},
		},
	)

	rootQuery := graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"pool": &graphql.Field{
				Type:        graphql.NewList(poolType),
				Description: `A list with all available pool items.`,
				Resolve:     queryPoolItems,
			},
			"slots": &graphql.Field{
				Type:        graphql.NewList(slotType),
				Description: `A list with all populated slots.`,
				Resolve:     querySlotItems,
				Args: graphql.FieldConfigArgument{
					"slot": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: `The slot from which the information should come. Optional.`,
					},
				},
			},
		},
	}

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"setSlotMapping": &graphql.Field{
				Type: graphql.Boolean,
				Description: `Set the mapping of a slot to the given filename in the pool. 
				If the slot already exists, it will be overwritten with the new data.`,
				Args: graphql.FieldConfigArgument{
					"slot": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.Int),
						Description: `The slot to which the mapping will be set.`,
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
						Description: `The name of the poolItem that should be mapped to the slot.
						The format has to be <name>.<extension>.`,
					},
				},
				Resolve: mutateSetSlotMap,
			},
			"removeSlotMapping": &graphql.Field{
				Type:        graphql.Boolean,
				Description: `Clears a slot from its mapping. If something was cleared it returns true.`,
				Args: graphql.FieldConfigArgument{
					"slot": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.Int),
						Description: `The slot that should be cleared.`,
					},
				},
				Resolve: mutateClearSlot,
			},
			"setPlaying": &graphql.Field{
				Type:        playerType,
				Description: `Set the playing state for a slot.`,
				Args: graphql.FieldConfigArgument{
					"slot": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.Int),
						Description: `The slot that should be changed in its playing state.`,
					},
					"player": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(playerInputType),
						Description: `The player input type with the data to set. 
						Any empty fields will not change the current state.`,
					},
				},
				Resolve: mutatePlayingState,
			},
			"removeFile": &graphql.Field{
				Type: graphql.Boolean,
				Description: `Delete the file from the pool. It will not be deleted when it is in use. 
				The return value is true if the file could be deleted.`,
				Args: graphql.FieldConfigArgument{
					"file": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: `The file that should be deleted. Format: <name>.<extension>`,
					},
				},
				Resolve: mutateDeleteFile,
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    graphql.NewObject(rootQuery),
		Mutation: rootMutation,
	})
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true, //hosts a graphiql instance under the entrypoint
	})
	// serve HTTP
	http.Handle("/graphql", h)
}
