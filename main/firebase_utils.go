package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	firebase "firebase.google.com/go"
	db "firebase.google.com/go/db"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

type FireBaseStore struct {
	Database *db.Client
}

func (store *FireBaseStore) init(serviceAcctFile string) error {
	conf := &firebase.Config{DatabaseURL: "https://inventorysearch-9682b.firebaseio.com/"}
	opt := option.WithCredentialsFile(serviceAcctFile)
	app, err := firebase.NewApp(context.Background(), conf, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	} else {
		db, err := app.Database(context.Background())
		if err != nil {
			return err
		} else {
			store.Database = db
			return nil
		}
	}
}

func (store *FireBaseStore) updateVendorPrice(item Item) error {
	vendorPrices := store.Database.NewRef("vendorprices")
	queries, err := vendorPrices.OrderByChild("Code").EqualTo(item.Code).GetOrdered(context.Background())

	if err != nil {
		return err
	}

	for _, query := range queries {
		vendorPrices.Child(query.Key()).Update(
			context.Background(),
			map[string]interface{}{
				"retailPriceVAT": item.RetailPriceVAT,
				"extraPrice":     item.ExtraPrice,
			},
		)
	}

	return nil
}

func (store *FireBaseStore) updateItemPrice(item Item) error {
	items := store.Database.NewRef("items")
	queries, err := items.OrderByChild("Code").EqualTo(item.Code).GetOrdered(context.Background())
	if err != nil {
		return err
	}

	for _, query := range queries {
		items.Child(query.Key()).Update(
			context.Background(),
			map[string]interface{}{
				"retailPriceVAT": item.RetailPriceVAT,
				"extraPrice":     item.ExtraPrice,
			})
	}
	return nil
}

type Item struct {
	Code           string
	VmName         string
	RetailPriceVAT string
	ExtraPrice     string
}

func InitFirebase() (*firebase.App, error) {
	return InitFirebaseWithAccountFile("../serviceAccountKey.json")
}

func InitFirebaseWithAccountFile(filename string) (*firebase.App, error) {
	conf := &firebase.Config{DatabaseURL: "https://inventorysearch-9682b.firebaseio.com/"}
	opt := option.WithCredentialsFile(filename)
	app, err := firebase.NewApp(context.Background(), conf, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	return app, nil
}

func Database(App *firebase.App) (*db.Client, error) {
	return App.Database(context.Background())
}

func GetItem(Client *db.Client) ([]Item, error) {
	var items []Item

	err := Client.NewRef("items").OrderByKey().LimitToFirst(10).Get(context.Background(), &items)
	if err != nil {
		return nil, err
	} else {
		return items, nil
	}
}

func processItems(client *db.Client) error {
	cred, logonErr := logon()
	if logonErr != nil {
		return logonErr
	}

	itemRef := client.NewRef("items")
	qs, err := itemRef.OrderByKey().GetOrdered(context.Background())
	if err != nil {
		return err
	}

	for _, q := range qs {
		var v map[string]interface{}
		if err := q.Unmarshal(&v); err != nil {
			log.Printf("Error unmarshalling data %v\n", err)
		} else {
			if v["Code"] == nil {
				continue
			}
			item, err := getItem(fmt.Sprint(v["Code"]), cred)
			if err != nil {
				log.Printf("Error retrieving item from Anacle server: %v\n", err)
			} else {
				curPrice, err :=
					strconv.ParseFloat(strings.TrimSpace(strings.Replace(fmt.Sprint(v["retailPriceVAT"]), ",", "", -1)), 32)
				if err != nil {
					log.Printf("Unable to read price of item with code %v, error details %v\n", v["Code"], err.Error())
				} else {
					if curPrice != item.UnitPrice && item.UnitPrice == 0 {
						log.Printf("Item Key: %v, Current Price: %v, New Price: %v\n", q.Key(), curPrice, item.UnitPrice)
						if err := updatePrice(q.Key(), item.UnitPrice, itemRef); err != nil {
							log.Fatalf("Unable to update price for item with code %v", item.ObjectNumber)
						}
					} else {
						log.Printf("New price %v was not updated for item with Code %v\n", item.UnitPrice, item.ObjectNumber)
					}
				}
			}
		}
	}

	return nil
}

func updatePrice(key string, newPrice float64, ref *db.Ref) error {
	return ref.Child(key).Update(context.Background(), map[string]interface{}{"retailPriceVAT": newPrice})
}
