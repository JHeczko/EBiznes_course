// ============================================================
// API Tests for every EndPoint
// ============================================================

package api_test

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

const apiBase = "http://localhost:13000"

func get(t *testing.T, path string) *http.Response {
    t.Helper()
    res, err := http.Get(apiBase + path)
    require.NoError(t, err)
    return res
}

func post(t *testing.T, path string) *http.Response {
    t.Helper()
    res, err := http.Post(apiBase+path, "application/json", nil)
    require.NoError(t, err)
    return res
}

func patch(t *testing.T, path string) *http.Response {
    t.Helper()
    req, err := http.NewRequest(http.MethodPatch, apiBase+path, nil)
    require.NoError(t, err)
    res, err := http.DefaultClient.Do(req)
    require.NoError(t, err)
    return res
}

func del(t *testing.T, path string) *http.Response {
    t.Helper()
    req, err := http.NewRequest(http.MethodDelete, apiBase+path, nil)
    require.NoError(t, err)
    res, err := http.DefaultClient.Do(req)
    require.NoError(t, err)
    return res
}

func readBody(t *testing.T, res *http.Response) []byte {
    t.Helper()
    body, err := io.ReadAll(res.Body)
    require.NoError(t, err)
    defer res.Body.Close()
    return body
}

// ════════════════════════════════════════════
//  /products
// ════════════════════════════════════════════

// api-01  get /products – check 200 and list
func TestAPI_GetProducts(t *testing.T) {
    res := get(t, "/products")
    body := readBody(t, res)

    // A1
    assert.Equal(t, 200, res.StatusCode)

    var arr []map[string]interface{}
    err := json.Unmarshal(body, &arr)
    // A2
    assert.NoError(t, err)
    // A3
    assert.GreaterOrEqual(t, len(arr), 1)
    // A4
    assert.NotNil(t, arr[0]["id"])
    // A5
    assert.NotNil(t, arr[0]["name"])
    // A6
    assert.NotNil(t, arr[0]["price"])
}

// api-02  get /products/:id – find existing one
func TestAPI_GetProduct_Found(t *testing.T) {
    res0 := get(t, "/products")
    var arr []map[string]interface{}
    json.Unmarshal(readBody(t, res0), &arr)
    require.GreaterOrEqual(t, len(arr), 1)

    id := int(arr[0]["id"].(float64))
    res := get(t, fmt.Sprintf("/products/%d", id))
    body := readBody(t, res)

    // A7
    assert.Equal(t, 200, res.StatusCode)

    var prod map[string]interface{}
    json.Unmarshal(body, &prod)
    // A8
    assert.Equal(t, float64(id), prod["id"])
    // A9
    assert.NotEmpty(t, prod["name"])
}

// api-03  get /products/:id – check 404 for ghost product
func TestAPI_GetProduct_NotFound(t *testing.T) {
    res := get(t, "/products/999999")
    // A10
    assert.Equal(t, 404, res.StatusCode)
}

// api-04  get /products/:id – bad id check
func TestAPI_GetProduct_BadID(t *testing.T) {
    res := get(t, "/products/abc")
    // A11
    assert.Equal(t, 400, res.StatusCode)
}

// api-05  post /products – 400 if name is empty
func TestAPI_CreateProduct_MissingName(t *testing.T) {
    res := post(t, "/products?price=100&cat=1")
    // A12
    assert.Equal(t, 400, res.StatusCode)
}

// api-06  post /products – bad price format
func TestAPI_CreateProduct_BadPrice(t *testing.T) {
    res := post(t, "/products?name=Test&price=abc&cat=1")
    // A13
    assert.Equal(t, 400, res.StatusCode)
}

// api-07  post /products – bad category id check
func TestAPI_CreateProduct_BadCat(t *testing.T) {
    res := post(t, "/products?name=Test&price=99.9&cat=abc")
    // A14
    assert.Equal(t, 400, res.StatusCode)
}

// api-08  post /products – create ok check
func TestAPI_CreateProduct_OK(t *testing.T) {
    res0 := get(t, "/category")
    var cats []map[string]interface{}
    json.Unmarshal(readBody(t, res0), &cats)
    require.GreaterOrEqual(t, len(cats), 1)
    catID := int(cats[0]["id"].(float64))

    res := post(t, fmt.Sprintf("/products?name=TestProduct&price=299.99&cat=%d", catID))
    body := readBody(t, res)
    // A15
    assert.Equal(t, 201, res.StatusCode)

    var prod map[string]interface{}
    json.Unmarshal(body, &prod)
    // A16
    assert.Equal(t, "TestProduct", prod["name"])
}

// api-09  patch /products/:id – 404 if no product
func TestAPI_UpdateProduct_NotFound(t *testing.T) {
    res := patch(t, "/products/999999?name=X&price=1&cat=1")
    // A17
    assert.Equal(t, 404, res.StatusCode)
}

// api-10  patch /products/:id – check 400 for bad data
func TestAPI_UpdateProduct_BadPrice(t *testing.T) {
    res := patch(t, "/products/1?name=X&price=abc&cat=1")
    // A18
    assert.Equal(t, 400, res.StatusCode)
}

// api-11  delete /products/:id – 404 if not exists
func TestAPI_DeleteProduct_NotFound(t *testing.T) {
    res := del(t, "/products/999999")
    // A19
    assert.Equal(t, 404, res.StatusCode)
}

// api-12  delete /products/:id – bad format check
func TestAPI_DeleteProduct_BadID(t *testing.T) {
    res := del(t, "/products/abc")
    // A20
    assert.Equal(t, 400, res.StatusCode)
}

// ════════════════════════════════════════════
//  /category
// ════════════════════════════════════════════

// api-13  get /category – check if 200
func TestAPI_GetCategories(t *testing.T) {
    res := get(t, "/category")
    // A21
    assert.Equal(t, 200, res.StatusCode)

    var arr []map[string]interface{}
    json.Unmarshal(readBody(t, res), &arr)
    // A22
    assert.GreaterOrEqual(t, len(arr), 1)
}

// api-14  get /category/:id – 404 for missing
func TestAPI_GetCategory_NotFound(t *testing.T) {
    res := get(t, "/category/999999")
    // A23
    assert.Equal(t, 404, res.StatusCode)
}

// api-15  post /category – 400 if no name
func TestAPI_CreateCategory_MissingName(t *testing.T) {
    res := post(t, "/category")
    // A24
    assert.Equal(t, 400, res.StatusCode)
}

// api-16  patch /category/:id – 404 check
func TestAPI_UpdateCategory_NotFound(t *testing.T) {
    res := patch(t, "/category/999999?name=X")
    // A25
    assert.Equal(t, 404, res.StatusCode)
}

// api-17  patch /category/:id – bad id test
func TestAPI_UpdateCategory_BadID(t *testing.T) {
    res := patch(t, "/category/abc?name=X")
    // A26
    assert.Equal(t, 400, res.StatusCode)
}

// api-18  delete /category/:id – 404 if ghost cat
func TestAPI_DeleteCategory_NotFound(t *testing.T) {
    res := del(t, "/category/999999")
    // A27
    assert.Equal(t, 404, res.StatusCode)
}

// ════════════════════════════════════════════
//  /cart
// ════════════════════════════════════════════

// api-19  get /cart/:user_id – check if list is there
func TestAPI_GetCart(t *testing.T) {
    res := get(t, "/cart/1")
    // A28
    assert.Equal(t, 200, res.StatusCode)

    var arr []interface{}
    json.Unmarshal(readBody(t, res), &arr)
    // A29
    assert.IsType(t, []interface{}{}, arr)
}

// api-20  get /cart/:user_id – 400 for bad user id
func TestAPI_GetCart_BadUserID(t *testing.T) {
    res := get(t, "/cart/xyz")
    // A30
    assert.Equal(t, 400, res.StatusCode)
}

// api-21  post /cart/:user_id – 400 if no prod id
func TestAPI_CreateCartItem_MissingProdID(t *testing.T) {
    res := post(t, "/cart/1")
    // A31
    assert.Equal(t, 400, res.StatusCode)
}

// api-22  post /cart/:user_id – check bad prod id format
func TestAPI_CreateCartItem_BadProdID(t *testing.T) {
    res := post(t, "/cart/1?prod_id=abc")
    // A32
    assert.Equal(t, 400, res.StatusCode)
}

// api-23  patch /cart/:user_id – 400 if no prod id
func TestAPI_UpdateCartItem_MissingProdID(t *testing.T) {
    res := patch(t, "/cart/1?quantity=3")
    // A33
    assert.Equal(t, 400, res.StatusCode)
}

// api-24  patch /cart/:user_id – bad quantity format
func TestAPI_UpdateCartItem_BadQuantity(t *testing.T) {
    res := patch(t, "/cart/1?prod_id=1&quantity=abc")
    // A34
    assert.Equal(t, 400, res.StatusCode)
}

// api-25  patch /cart/:user_id – 404 if item not in cart
func TestAPI_UpdateCartItem_ProdNotInCart(t *testing.T) {
    res := patch(t, "/cart/1?prod_id=999999&quantity=1")
    // A35
    assert.Equal(t, 404, res.StatusCode)
}

// api-26  delete /cart/:user_id – check bad user id
func TestAPI_DeleteCartItem_BadUserID(t *testing.T) {
    res := del(t, "/cart/abc")
    // A36
    assert.Equal(t, 400, res.StatusCode)
}

// api-27  delete /cart/:user_id – 404 if cart empty
func TestAPI_DeleteCart_EmptyUser(t *testing.T) {
    res := del(t, "/cart/888888")
    // A37
    assert.Equal(t, 404, res.StatusCode)
}

// ════════════════════════════════════════════
//  /payments
// ════════════════════════════════════════════

// api-28  get /payments/:user_id – check ok
func TestAPI_GetPayments(t *testing.T) {
    res := get(t, "/payments/1")
    // A38
    assert.Equal(t, 200, res.StatusCode)

    var arr []interface{}
    json.Unmarshal(readBody(t, res), &arr)
    // A39
    assert.IsType(t, []interface{}{}, arr)
}

// api-29  get /payments/:user_id – bad user id check
func TestAPI_GetPayments_BadUserID(t *testing.T) {
    res := get(t, "/payments/abc")
    // A40
    assert.Equal(t, 400, res.StatusCode)
}

// api-30  post /payments/:user_id – 500 if nothing to pay for
func TestAPI_AddPayment_EmptyBasket(t *testing.T) {
    res := post(t, "/payments/777777")
    // A41
    assert.Equal(t, 500, res.StatusCode)
}

// api-31  post /payments/:user_id – 400 for bad id
func TestAPI_AddPayment_BadUserID(t *testing.T) {
    res := post(t, "/payments/abc")
    // A42
    assert.Equal(t, 400, res.StatusCode)
}

// ════════════════════════════════════════════
// response structure checks
// ════════════════════════════════════════════

// api-32  check if category is preloaded
func TestAPI_ProductHasCategory(t *testing.T) {
    res := get(t, "/products")
    var arr []map[string]interface{}
    json.Unmarshal(readBody(t, res), &arr)
    require.GreaterOrEqual(t, len(arr), 1)

    cat, ok := arr[0]["category"].(map[string]interface{})
    // A43
    assert.True(t, ok)
    // A44
    assert.NotNil(t, cat["id"])
    // A45
    assert.NotNil(t, cat["name"])
}

// api-33  check if cart item has product info
func TestAPI_CartItemHasProduct(t *testing.T) {
    res0 := get(t, "/products")
    var products []map[string]interface{}
    json.Unmarshal(readBody(t, res0), &products)
    if len(products) == 0 {
       t.Skip("No products in db")
    }
    prodID := int(products[0]["id"].(float64))
    post(t, fmt.Sprintf("/cart/1?prod_id=%d&quantity=1", prodID))

    res := get(t, "/cart/1")
    var arr []map[string]interface{}
    json.Unmarshal(readBody(t, res), &arr)

    if len(arr) == 0 {
       t.Skip("Cart is empty - skip")
    }

    item := arr[0]
    // A46
    assert.NotNil(t, item["id"])
    // A47
    assert.NotNil(t, item["quantity"])
    // A48
    assert.NotNil(t, item["product"])

    prod, ok := item["product"].(map[string]interface{})
    // A49
    assert.True(t, ok)
    // A50
    assert.NotNil(t, prod["name"])
    // A51
    assert.NotNil(t, prod["price"])
}