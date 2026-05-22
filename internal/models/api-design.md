# High-Level Architecture

Our database currently models five major domains:

```text
Users
│
├── Refresh Tokens
├── Cart
│   └── Cart Items
└── Orders
    └── Order Items

Categories
└── Products
    └── Product Images
```

Visually:

```text
USER
 │
 ├──── CART
 │       │
 │       └──── CART_ITEMS
 │                 │
 │                 └──── PRODUCT
 │
 ├──── ORDERS
 │       │
 │       └──── ORDER_ITEMS
 │                 │
 │                 └──── PRODUCT
 │
 └──── REFRESH_TOKENS


CATEGORY
   │
   └──── PRODUCTS
             │
             └──── PRODUCT_IMAGES
```

This separation is important because every table represents a different business concept.

---

# 1. User

```go
type User struct
```

Represents a customer or administrator.

Database table:

```sql
users
```

Example row:

| id | email                                   | role     |
| -- | --------------------------------------- | -------- |
| 1  | [john@gmail.com](mailto:john@gmail.com) | customer |

---

## Why store role?

```go
Role UserRole
```

Because not everyone has the same permissions.

Example:

```text
Customer:
✓ Browse products
✓ Add to cart
✓ Place order

Admin:
✓ Create products
✓ Delete products
✓ Manage orders
✓ Manage categories
```

This becomes important later in middleware.

Example:

```go
if user.Role != UserRoleAdmin {
    return forbidden
}
```

---

## Why IsActive?

```go
IsActive bool
```

Instead of deleting users:

```sql
DELETE FROM users
```

we can simply disable them:

```sql
UPDATE users
SET is_active = false
```

Benefits:

* preserve order history
* preserve audit logs
* easier recovery

---

# 2. RefreshToken

```go
type RefreshToken struct
```

This table supports JWT authentication.

---

Without refresh tokens:

```text
Login
 ↓
JWT expires
 ↓
User must login again
```

With refresh tokens:

```text
Login
 ↓
Access Token (15 min)
Refresh Token (7 days)
 ↓
Access token expires
 ↓
Refresh token generates new access token
```

Example:

| id | user_id | token    |
| -- | ------- | -------- |
| 1  | 5       | asdkj123 |

---

Relationship:

```go
User
  └── RefreshTokens
```

One user can have many sessions:

```text
Laptop login
Phone login
Tablet login
```

Each session can have its own refresh token.

---

# 3. Category

```go
type Category struct
```

Represents product grouping.

Examples:

```text
Electronics
Books
Clothing
Shoes
Furniture
```

Table:

| id | name        |
| -- | ----------- |
| 1  | Electronics |
| 2  | Books       |

---

Why not store category name directly inside product?

Bad:

```text
iPhone -> Electronics
Samsung -> Electronics
Macbook -> Electronics
```

Repeated thousands of times.

Instead:

```sql
products.category_id
```

This is normalization.

---

# 4. Product

This is the center of the store.

Example:

```text
iPhone 17
₹99999
Stock = 20
SKU = IP17-BLK
```

---

Fields:

### CategoryID

```go
CategoryID uint
```

Foreign key.

Meaning:

```text
This product belongs to category X
```

---

### Stock

```go
Stock int
```

Inventory tracking.

Example:

```text
20 units available
```

Order placed:

```text
20 -> 19
```

---

### SKU

```go
SKU string
```

Stock Keeping Unit.

Example:

```text
IPHONE17-BLACK-256
```

Used internally.

Customers usually don't see it.

Warehouse systems do.

---

### IsActive

Useful when product is discontinued.

Instead of deleting:

```text
Set active=false
```

Product disappears from catalog.

Order history remains intact.

---

# 5. ProductImage

Products often have multiple images.

Example:

```text
Front view
Back view
Side view
Packaging
```

Therefore:

```go
Images []ProductImage
```

---

Database:

```text
product_images
```

Example:

| id | product_id | url       |
| -- | ---------- | --------- |
| 1  | 5          | front.jpg |
| 2  | 5          | back.jpg  |

---

Relationship:

```text
Product
 ├── Image1
 ├── Image2
 └── Image3
```

One-to-many.

---

Why separate table?

Because number of images varies.

Bad:

```sql
image1
image2
image3
image4
image5
```

Some products need 1 image.

Some need 20.

Separate table is flexible.

---

# 6. Cart

Represents shopping cart.

Example:

```text
John's Cart
```

Contains:

```text
iPhone
AirPods
Mouse
```

---

Why separate Cart table?

Because cart itself is an entity.

We may later store:

```go
LastViewedAt
CouponCode
Currency
ShippingEstimate
```

inside cart.

---

Relationship:

```text
User
 └── Cart
```

Typically one cart per user.

---

# 7. CartItem

This is where actual products inside the cart live.

Example:

```text
Cart:
    iPhone x2
    AirPods x1
```

Stored as:

| cart_id | product_id | quantity |
| ------- | ---------- | -------- |
| 1       | 5          | 2        |
| 1       | 9          | 1        |

---

Why not directly store products inside cart?

Because cart needs:

```text
multiple products
multiple quantities
```

That's a many-to-many relationship.

CartItem acts as a join table.

---

Visual:

```text
Cart
 ├── CartItem
 │     └── Product
 ├── CartItem
 │     └── Product
 └── CartItem
       └── Product
```

---

# 8. Order

Represents completed checkout.

Example:

```text
Order #1001
Customer John
Total ₹150000
Status shipped
```

---

Why separate from Cart?

Cart changes constantly:

```text
add item
remove item
change quantity
```

Order must never change after purchase.

Order is a snapshot.

---

Fields:

### TotalAmount

```go
TotalAmount float64
```

Stores final purchase amount.

Example:

```text
99999 + 49999
=
149998
```

---

### Status

```go
pending
confirmed
shipped
delivered
cancelled
```

Tracks lifecycle.

Example:

```text
pending
 ↓
confirmed
 ↓
shipped
 ↓
delivered
```

---

# 9. OrderItem

This is extremely important.

Example order:

```text
iPhone x2
AirPods x1
```

Stored as:

| order_id | product_id | quantity | price |
| -------- | ---------- | -------- | ----- |
| 1        | 5          | 2        | 99999 |
| 1        | 9          | 1        | 19999 |

---

Question:

Why store price?

We already have Product.Price.

Because prices change.

Example:

Today:

```text
iPhone = 1000
```

Tomorrow:

```text
iPhone = 1200
```

Customer bought yesterday.

Order history must remain:

```text
Paid 1000
```

not:

```text
Paid 1200
```

Therefore:

```go
Price float64
```

stores purchase-time price.

This is a very important real-world design decision.

---

# How Checkout Works

Imagine user buys products.

---

Step 1

Browse products:

```text
GET /products
```

---

Step 2

Add to cart:

```text
POST /cart/items
```

Creates:

```text
CartItem
```

---

Step 3

View cart:

```text
GET /cart
```

Loads:

```go
Cart
 -> CartItems
 -> Products
```

---

Step 4

Checkout

Backend:

```go
Create Order
Create OrderItems
Calculate TotalAmount
Reduce Product Stock
Clear Cart
```

---

Result:

```text
Cart
 ↓
converted into
 ↓
Order
```

---

# Why GORM Relationships Exist

Example:

```go
OrderItems []OrderItem
```

Without relationship:

```go
db.Find(&orders)
```

returns only order data.

With preload:

```go
db.Preload("OrderItems").Find(&orders)
```

GORM automatically loads:

```json
{
  "id": 1,
  "order_items": [...]
}
```

Similarly:

```go
db.Preload("Images").Find(&products)
```

loads all product images.

That's why relationship fields exist—they allow object graphs to be loaded easily instead of manually writing joins everywhere.

---

Overall, this schema models the entire lifecycle of an online purchase:

```text
User
 ↓
Browse Products
 ↓
Add to Cart
 ↓
Checkout
 ↓
Create Order
 ↓
Track Order Status
 ↓
Delivered
```