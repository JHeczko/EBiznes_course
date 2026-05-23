// ============================================================
// Unit Tests for Handlers
// ============================================================

package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "zadanie4/models"

    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// ─────────────────────────────────────────────
// helpers
// ─────────────────────────────────────────────

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
       Logger: logger.Default.LogMode(logger.Silent),
    })
    require.NoError(t, err)

    err = db.AutoMigrate(
       &models.Category{},
       &models.Product{},
       &models.Users{},
       &models.Basket{},
       &models.Payments{},
    )
    require.NoError(t, err)
    return db
}

func seedData(t *testing.T, db *gorm.DB) {
    cat := models.Category{CategoryName: "Furnitures"}
    require.NoError(t, db.Create(&cat).Error)

    catID := cat.CategoryID
    p1 := models.Product{ProductName: "Couch", Price: 1500.00, CategoryID: &catID}
    p2 := models.Product{ProductName: "Chair", Price: 150.00, CategoryID: &catID}
    require.NoError(t, db.Create(&p1).Error)
    require.NoError(t, db.Create(&p2).Error)

    user := models.Users{UserName: "testuser", Email: "test@example.com"}
    require.NoError(t, db.Create(&user).Error)

    require.NoError(t, db.Create(&models.Basket{UserID: user.UserID, ProductID: p1.ProductID, Quantity: 2}).Error)
}

func echoCtx(method, path string) (echo.Context, *httptest.ResponseRecorder) {
    e := echo.New()
    req := httptest.NewRequest(method, path, nil)
    rec := httptest.NewRecorder()
    return e.NewContext(req, rec), rec
}

// ─────────────────────────────────────────────
// unit-01  getcategories – check 200 and not empty list
// ─────────────────────────────────────────────
func TestGetCategories_ReturnsOK(t *testing.T) {
    db := setupTestDB(t)
    seedData(t, db)

    c, rec := echoCtx(http.MethodGet, "/category")

    err := GetCategories(db)(c)
    // A1
    assert.NoError(t, err)
    // A2
    assert.Equal(t, http.StatusOK, rec.Code)

    var cats []models.Category
    err = json.Unmarshal(rec.Body.Bytes(), &cats)
    require.NoError(t, err)
    // A3
    assert.GreaterOrEqual(t, len(cats), 1)
    // A4
    assert.Equal(t, "Furnitures", cats[0].CategoryName)
}

// ─────────────────────────────────────────────
// unit-02  getcategory – test with good id
// ─────────────────────────────────────────────
func TestGetCategory_ExistingID(t *testing.T) {
    db := setupTestDB(t)
    seedData(t, db)

    var cat models.Category
    require.NoError(t, db.First(&cat).Error)

    c, rec := echoCtx(http.MethodGet, "/category/"+fmt.Sprintf("%d", cat.CategoryID))
    c.SetParamNames("id")
    c.SetParamValues(fmt.Sprintf("%d", cat.CategoryID))

    err := GetCategory(db)(c)
    // A5
    assert.NoError(t, err)
    // A6
    assert.Equal(t, http.StatusOK, rec.Code)

    var result models.Category
    err = json.Unmarshal(rec.Body.Bytes(), &result)
    require.NoError(t, err)
    // A7
    assert.Equal(t, cat.CategoryID, result.CategoryID)
    // A8
    assert.Equal(t, "Furnitures", result.CategoryName)
}

// ─────────────────────────────────────────────
// unit-03  getcategory – check 404 when id is missing
// ─────────────────────────────────────────────
func TestGetCategory_NotFound(t *testing.T) {
    db := setupTestDB(t)

    c, _ := echoCtx(http.MethodGet, "/category/9999")
    c.SetParamNames("id")
    c.SetParamValues("9999")

    err := GetCategory(db)(c)
    // A9
    assert.Error(t, err)
    he, ok := err.(*echo.HTTPError)
    // A10
    assert.True(t, ok)
    // A11
    assert.Equal(t, http.StatusNotFound, he.Code)
}

// ─────────────────────────────────────────────
// unit-04  getcategory – check 400 for bad id format
// ─────────────────────────────────────────────
func TestGetCategory_BadID(t *testing.T) {
    db := setupTestDB(t)

    c, _ := echoCtx(http.MethodGet, "/category/abc")
    c.SetParamNames("id")
    c.SetParamValues("abc")

    err := GetCategory(db)(c)
    he, ok := err.(*echo.HTTPError)
    // A12
    assert.True(t, ok)
    // A13
    assert.Equal(t, http.StatusBadRequest, he.Code)
}

// ─────────────────────────────────────────────
// unit-05  createcategory – test good data
// ─────────────────────────────────────────────
func TestCreateCategory_OK(t *testing.T) {
    db := setupTestDB(t)

    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/category?name=Sleepingroom", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    err := CreateCategory(db)(c)
    // A14
    assert.NoError(t, err)
    // A15
    assert.Equal(t, http.StatusOK, rec.Code)

    var cat models.Category
    err = json.Unmarshal(rec.Body.Bytes(), &cat)
    require.NoError(t, err)
    // A16
    assert.Equal(t, "Sleepingroom", cat.CategoryName)
    // A17
    assert.NotZero(t, cat.CategoryID)
}

// ─────────────────────────────────────────────
// unit-06  createcategory – 400 if name is gone
// ─────────────────────────────────────────────
func TestCreateCategory_MissingName(t *testing.T) {
    db := setupTestDB(t)

    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/category", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    err := CreateCategory(db)(c)
    he, ok := err.(*echo.HTTPError)
    // A18
    assert.True(t, ok)
    // A19
    assert.Equal(t, http.StatusBadRequest, he.Code)
}

// ─────────────────────────────────────────────
// unit-07  getproducts – check 200 and list
// ─────────────────────────────────────────────
func TestGetProducts_ReturnsOK(t *testing.T) {
    db := setupTestDB(t)
    seedData(t, db)

    c, rec := echoCtx(http.MethodGet, "/products")

    err := GetProducts(db)(c)
    // A20
    assert.NoError(t, err)
    // A21
    assert.Equal(t, http.StatusOK, rec.Code)

    var products []models.Product
    err = json.Unmarshal(rec.Body.Bytes(), &products)
    require.NoError(t, err)
    // A22
    assert.GreaterOrEqual(t, len(products), 2)
}

// ─────────────────────────────────────────────
// unit-08  getproduct – good id check
// ─────────────────────────────────────────────
func TestGetProduct_ExistingID(t *testing.T) {
    db := setupTestDB(t)
    seedData(t, db)

    var prod models.Product
    require.NoError(t, db.First(&prod).Error)

    c, rec := echoCtx(http.MethodGet, "/products/"+fmt.Sprintf("%d", prod.ProductID))
    c.SetParamNames("id")
    c.SetParamValues(fmt.Sprintf("%d", prod.ProductID))

    err := GetProduct(db)(c)
    // A23
    assert.NoError(t, err)
    // A24
    assert.Equal(t, http.StatusOK, rec.Code)

    var result models.Product
    err = json.Unmarshal(rec.Body.Bytes(), &result)
    require.NoError(t, err)
    // A25
    assert.Equal(t, prod.ProductID, result.ProductID)
    // A26
    assert.Equal(t, "Couch", result.ProductName)
    // A27
    assert.Equal(t, float32(1500.00), result.Price)
}

// ─────────────────────────────────────────────
// unit-09  getproduct – 404 for missing id
// ─────────────────────────────────────────────
func TestGetProduct_NotFound(t *testing.T) {
    db := setupTestDB(t)

    c, _ := echoCtx(http.MethodGet, "/products/9999")
    c.SetParamNames("id")
    c.SetParamValues("9999")

    err := GetProduct(db)(c)
    he, ok := err.(*echo.HTTPError)
    // A28
    assert.True(t, ok)
    // A29
    assert.Equal(t, http.StatusNotFound, he.Code)
}

// ─────────────────────────────────────────────
// unit-10  getproduct – 400 for bad id
// ─────────────────────────────────────────────
func TestGetProduct_BadID(t *testing.T) {
    db := setupTestDB(t)

    c, _ := echoCtx(http.MethodGet, "/products/xyz")
    c.SetParamNames("id")
    c.SetParamValues("xyz")

    err := GetProduct(db)(c)
    he, ok := err.(*echo.HTTPError)
    // A30
    assert.True(t, ok)
    // A31
    assert.Equal(t, http.StatusBadRequest, he.Code)
}

// ─────────────────────────────────────────────
// unit-11  getitems – test cart for user
// ─────────────────────────────────────────────
func TestGetItems_ExistingUser(t *testing.T) {
    db := setupTestDB(t)
    seedData(t, db)

    var user models.Users
    require.NoError(t, db.First(&user).Error)

    c, rec := echoCtx(http.MethodGet, "/cart/"+fmt.Sprintf("%d", user.UserID))
    c.SetParamNames("user_id")
    c.SetParamValues(fmt.Sprintf("%d", user.UserID))

    err := GetItems(db)(c)
    // A32
    assert.NoError(t, err)
    // A33
    assert.Equal(t, http.StatusOK, rec.Code)

    var items []models.Basket
    err = json.Unmarshal(rec.Body.Bytes(), &items)
    require.NoError(t, err)
    // A34
    assert.Equal(t, 1, len(items))
    // A35
    assert.Equal(t, uint(2), items[0].Quantity)
}

// ─────────────────────────────────────────────
// unit-12  getitems – 400 for bad user_id
// ─────────────────────────────────────────────
func TestGetItems_BadUserID(t *testing.T) {
    db := setupTestDB(t)

    c, _ := echoCtx(http.MethodGet, "/cart/abc")
    c.SetParamNames("user_id")
    c.SetParamValues("abc")

    err := GetItems(db)(c)
    he, ok := err.(*echo.HTTPError)
    // A36
    assert.True(t, ok)
    // A37
    assert.Equal(t, http.StatusBadRequest, he.Code)
}

// ─────────────────────────────────────────────
// unit-13  createitem – add thing to cart
// ─────────────────────────────────────────────
func TestCreateItem_NewEntry(t *testing.T) {
    db := setupTestDB(t)
    seedData(t, db)

    var user models.Users
    require.NoError(t, db.First(&user).Error)
    var prod models.Product
    require.NoError(t, db.Where("ProductName = ?", "Chair").First(&prod).Error)

    e := echo.New()
    req := httptest.NewRequest(http.MethodPost,
       fmt.Sprintf("/cart/%d?prod_id=%d&quantity=3", user.UserID, prod.ProductID), nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("user_id")
    c.SetParamValues(fmt.Sprintf("%d", user.UserID))

    err := CreateItem(db)(c)
    // A38
    assert.NoError(t, err)
    // A39
    assert.Equal(t, http.StatusOK, rec.Code)

    var item models.Basket
    err = json.Unmarshal(rec.Body.Bytes(), &item)
    require.NoError(t, err)
    // A40
    assert.Equal(t, uint(3), item.Quantity)
    // A41
    assert.Equal(t, prod.ProductID, item.ProductID)
}

// ─────────────────────────────────────────────
// unit-14  createitem – 400 if prod_id is wrong
// ─────────────────────────────────────────────
func TestCreateItem_BadProdID(t *testing.T) {
    db := setupTestDB(t)
    seedData(t, db)

    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/cart/1?prod_id=abc", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("user_id")
    c.SetParamValues("1")

    err := CreateItem(db)(c)
    he, ok := err.(*echo.HTTPError)
    // A42
    assert.True(t, ok)
    // A43
    assert.Equal(t, http.StatusBadRequest, he.Code)
}

// ─────────────────────────────────────────────
// unit-15  deleteitem – remove all if no prod_id
// ─────────────────────────────────────────────
func TestDeleteItem_AllItems(t *testing.T) {
    db := setupTestDB(t)
    seedData(t, db)

    var user models.Users
    require.NoError(t, db.First(&user).Error)

    c, rec := echoCtx(http.MethodDelete, "/cart/"+fmt.Sprintf("%d", user.UserID))
    c.SetParamNames("user_id")
    c.SetParamValues(fmt.Sprintf("%d", user.UserID))

    err := DeleteItem(db)(c)
    // A44
    assert.NoError(t, err)
    // A45
    assert.Equal(t, http.StatusOK, rec.Code)

    var response map[string]interface{}
    err = json.Unmarshal(rec.Body.Bytes(), &response)
    require.NoError(t, err)
    // A46
    assert.Equal(t, "all", response["type"])

    var remaining []models.Basket
    require.NoError(t, db.Where("UserID = ?", user.UserID).Find(&remaining).Error)
    // A47
    assert.Equal(t, 0, len(remaining))
}

// ─────────────────────────────────────────────
// unit-16  getpayments – check 200 ok
// ─────────────────────────────────────────────
func TestGetPayments_ReturnsOK(t *testing.T) {
    db := setupTestDB(t)
    seedData(t, db)

    var user models.Users
    require.NoError(t, db.First(&user).Error)

    c, rec := echoCtx(http.MethodGet, "/payments/"+fmt.Sprintf("%d", user.UserID))
    c.SetParamNames("user_id")
    c.SetParamValues(fmt.Sprintf("%d", user.UserID))

    err := GetPayments(db)(c)
    // A48
    assert.NoError(t, err)
    // A49
    assert.Equal(t, http.StatusOK, rec.Code)

    var payments []models.Payments
    err = json.Unmarshal(rec.Body.Bytes(), &payments)
    require.NoError(t, err)
    // A50
    assert.IsType(t, []models.Payments{}, payments)
}

// ─────────────────────────────────────────────
// unit-17  getpayments – 400 for bad user_id format
// ─────────────────────────────────────────────
func TestGetPayments_BadUserID(t *testing.T) {
    db := setupTestDB(t)

    c, _ := echoCtx(http.MethodGet, "/payments/xyz")
    c.SetParamNames("user_id")
    c.SetParamValues("xyz")

    err := GetPayments(db)(c)
    he, ok := err.(*echo.HTTPError)
    // A51
    assert.True(t, ok)
    // A52
    assert.Equal(t, http.StatusBadRequest, he.Code)
}

// ─────────────────────────────────────────────
// unit-18  addpayment – 500 if basket is empty
// ─────────────────────────────────────────────
func TestAddPayment_EmptyBasket(t *testing.T) {
    db := setupTestDB(t)
    seedData(t, db)

    c, _ := echoCtx(http.MethodPost, "/payments/9999")
    c.SetParamNames("user_id")
    c.SetParamValues("9999")

    err := AddPayment(db)(c)
    he, ok := err.(*echo.HTTPError)
    // A53
    assert.True(t, ok)
    // A54
    assert.Equal(t, http.StatusInternalServerError, he.Code)
}

// ─────────────────────────────────────────────
// unit-19  updatecategory – change name test
// ─────────────────────────────────────────────
func TestUpdateCategory_OK(t *testing.T) {
    db := setupTestDB(t)
    seedData(t, db)

    var cat models.Category
    require.NoError(t, db.First(&cat).Error)

    e := echo.New()
    req := httptest.NewRequest(http.MethodPatch,
       fmt.Sprintf("/category/%d?name=NewCat", cat.CategoryID), nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(fmt.Sprintf("%d", cat.CategoryID))

    err := UpdateCategory(db)(c)
    // A55
    assert.NoError(t, err)
    // A56
    assert.Equal(t, http.StatusOK, rec.Code)

    var updated models.Category
    require.NoError(t, db.First(&updated, cat.CategoryID).Error)
    // A57
    assert.Equal(t, "NewCat", updated.CategoryName)
}

// ─────────────────────────────────────────────
// unit-20  deletecategory – good id check
// ─────────────────────────────────────────────
func TestDeleteCategory_OK(t *testing.T) {
    db := setupTestDB(t)
    cat := models.Category{CategoryName: "ToDelete"}
    require.NoError(t, db.Create(&cat).Error)

    c, rec := echoCtx(http.MethodDelete, "/category/"+fmt.Sprintf("%d", cat.CategoryID))
    c.SetParamNames("id")
    c.SetParamValues(fmt.Sprintf("%d", cat.CategoryID))

    err := DeleteCategory(db)(c)
    // A58
    assert.NoError(t, err)
    // A59
    assert.Equal(t, http.StatusOK, rec.Code)

    var check models.Category
    result := db.First(&check, cat.CategoryID)
    // A60
    assert.Error(t, result.Error)
}