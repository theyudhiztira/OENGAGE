package pkg

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var adminRoleObjID, _ = primitive.ObjectIDFromHex("677a2e0bdcb48ff83fe4cf62")
var seederCreatedByObjID, _ = primitive.ObjectIDFromHex("677a2e0bdcb48ff83fe4cf61")

func GetAdminRoleObjID() primitive.ObjectID {
	return adminRoleObjID
}

var DBCollections = []CollectionIndex{
	{
		CollectionName: "system_configs",
		Indexes:        [][]Index{},
	},
	{
		CollectionName: "users",
		Indexes: [][]Index{
			{
				{
					Key:   "email",
					Value: 1,
					Options: map[string]interface{}{
						"unique": true,
					},
				},
			},
			{
				{
					Key:     "is_active",
					Value:   1,
					Options: map[string]interface{}{},
				},
				{
					Key:     "role_id",
					Value:   1,
					Options: map[string]interface{}{},
				},
			},
			{
				{
					Key:     "created_at",
					Value:   -1,
					Options: map[string]interface{}{},
				},
			},
			{
				{
					Key:     "name",
					Value:   "text",
					Options: map[string]interface{}{},
				},
			},
		},
	},
	{
		CollectionName: "roles",
		Indexes: [][]Index{
			{
				{
					Key:   "name",
					Value: 1,
					Options: map[string]interface{}{
						"unique": true,
					},
				},
			},
		},
	},
	{
		CollectionName: "modules",
		Indexes:        [][]Index{},
	},
	{
		CollectionName: "permissions",
		Indexes:        [][]Index{},
	},
	{
		CollectionName: "modules",
		Indexes: [][]Index{
			{
				{
					Key:     "name",
					Value:   1,
					Options: map[string]interface{}{},
				},
				{
					Key:     "path",
					Value:   "text",
					Options: map[string]interface{}{},
				},
			},
		},
	},
	{
		CollectionName: "whatsapp_configs",
		Indexes: [][]Index{
			{
				{
					Key:     "whatsapp_phone_number",
					Value:   1,
					Options: map[string]interface{}{},
				},
				{
					Key:     "createdAt",
					Value:   -1,
					Options: map[string]interface{}{},
				},
				{
					Key:     "updatedAt",
					Value:   -1,
					Options: map[string]interface{}{},
				},
			},
		},
	},
}
