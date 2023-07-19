# Welcome to dynamo 👋

Project to easy connect on AWS Dynamo

## install

```bash
go get https://github.com/vcore8/dynamo.git
```

## how to use

```golang
cfg, _ := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))

db := dynamo.New(cfg)

// List By Hash and Sort key
err := db.Table("my-table").Get("hash", "123").SortBy("sort","321").All(ctx, &result)

// Get By Hash and Sort key
err := db.Table("my-table").Get("hash", "123").SortBy("sort","321").One(ctx, &result)

// Create item
err := db.Table("my-table").Create(ctx, data)
```

## Author

👤 **Eduardo Mello**

- Github: [@EduardoRMello](https://github.com/EduardoRMello)

## Show your support

Give a ⭐️ if this project helped you!

---

_This README was generated with ❤️ by [readme-md-generator](https://github.com/kefranabg/readme-md-generator)_
