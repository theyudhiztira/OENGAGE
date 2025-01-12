package pkg

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Role        string             `bson:"role"`
	Permissions []Permission       `bson:"permissions"`
	CreatedBy   primitive.ObjectID `bson:"created_by"`
	CreeatedAt  time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type Permission struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty"`
	Module         string              `bson:"module"`
	PermissionRule ReadWritePermission `bson:"permission_rule"`
}

type ReadWritePermission struct {
	Read  bool `bson:"read"`
	Write bool `bson:"write"`
}

type Module struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Path        string             `bson:"path"`
	Description string             `bson:"description"`
}

type SystemConfig struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	IsSeeded       bool               `bson:"is_seeded"`
	IsMigrated     bool               `bson:"is_migrated"`
	WhatsappToken  string             `bson:"whatsapp_token"`
	WhatsappWabaID string             `bson:"whatsapp_waba_id"`
	InstalledAt    time.Time          `bson:"installed_at"`
}

var RolesSeeder = []Role{
	{
		ID:   GetAdminRoleObjID(),
		Role: "admin",
		Permissions: []Permission{
			{
				ID:     primitive.NewObjectID(),
				Module: "users",
				PermissionRule: ReadWritePermission{
					Read:  true,
					Write: true,
				},
			},
			{
				ID:     primitive.NewObjectID(),
				Module: "whatsapp",
				PermissionRule: ReadWritePermission{
					Read:  true,
					Write: true,
				},
			},
			{
				ID:     primitive.NewObjectID(),
				Module: "template",
				PermissionRule: ReadWritePermission{
					Read:  true,
					Write: true,
				},
			},
		},
		CreatedBy:  seederCreatedByObjID,
		CreeatedAt: time.Now(),
		UpdatedAt:  time.Now(),
	},
}

var ModulesSeeder = []Module{
	{
		Name:        "Users",
		Description: "Module for managing users",
		Path:        "/users",
	},
	{
		Name:        "Whatsapp",
		Description: "Module for managing whatsapp",
		Path:        "/whatsapp",
	},
	{
		Name:        "Template",
		Description: "Module for managing whatsapp templates",
		Path:        "/template",
	},
}

func SeedDB(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	sysConfigC := db.Collection(DBCollections[0].CollectionName)
	var systemConfig SystemConfig
	sysConfigC.FindOne(ctx, bson.D{}).Decode(&systemConfig)

	if systemConfig.IsMigrated {
		println("Database already migrated.")
		return nil
	}

	sysCErr := SeedSystemConfig(sysConfigC, ctx)
	if sysCErr != nil {
		panic(sysCErr.Error())
	}

	roleErr := SeedRoles(db.Collection("roles"), ctx)
	if roleErr != nil {
		panic(roleErr)
	}

	moduleErr := SeedModules(db.Collection("modules"), ctx)
	if moduleErr != nil {
		panic(moduleErr)
	}

	return nil
}

func SeedSystemConfig(c *mongo.Collection, ctx context.Context) error {
	log.Print("Seeding system config...")
	_, err := c.InsertOne(ctx, bson.M{
		"is_seeded":    true,
		"is_migrated":  true,
		"installed_at": time.Now(),
	})
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Print("Seeding system config done.")
	return nil
}

func SeedRoles(c *mongo.Collection, ctx context.Context) error {
	log.Print("Seeding roles...")
	documents := make([]interface{}, len(RolesSeeder))
	for i, role := range RolesSeeder {
		documents[i] = role
	}

	_, err := c.InsertMany(ctx, documents)
	if err != nil {
		return err
	}

	log.Print("Seeding roles done.")
	return nil
}

func SeedModules(c *mongo.Collection, ctx context.Context) error {
	log.Print("Seeding modules...")
	documents := make([]interface{}, len(ModulesSeeder))
	for i, module := range ModulesSeeder {
		documents[i] = module
	}

	_, err := c.InsertMany(ctx, documents)
	if err != nil {
		return err
	}

	log.Print("Seeding modules done.")
	return nil
}
