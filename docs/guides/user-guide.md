# ShopJoy Admin User Guide

> Step-by-step tutorials for managing your e-commerce store

## Table of Contents

1. [Getting Started](#getting-started)
2. [Product Management](#product-management)
3. [Category Management](#category-management)
4. [Brand Management](#brand-management)
5. [Market Configuration](#market-configuration)
6. [Inventory Management](#inventory-management)
7. [Best Practices](#best-practices)

---

## Getting Started

### Logging In

1. Navigate to your admin dashboard URL (e.g., `https://admin.shopjoy.com`)
2. Enter your email address and password
3. Click **Sign In**

> 💡 **Tip**: Your first login credentials are provided by your store administrator.

### Dashboard Overview

After logging in, you'll see the main dashboard with:

| Section | Description |
|---------|-------------|
| **Products** | Manage your product catalog |
| **Categories** | Organize products into categories |
| **Brands** | Manage brand information |
| **Markets** | Configure multi-market selling |
| **Inventory** | Stock and warehouse management |

---

## Product Management

### Creating Your First Product

1. **Navigate to Products**
   - Click **Products** in the left sidebar
   - Click **Create Product** button

2. **Fill in Basic Information**
   - **Name**: Enter product name (required)
   - **Description**: Add product description
   - **Category**: Select from dropdown
   - **Brand**: Optional brand association

3. **Set Pricing**
   - **Price**: Enter price in cents (e.g., 9900 = $99.00)
   - **Currency**: Select currency (CNY, USD, EUR, etc.)
   - **Cost Price**: Optional, for profit calculations

4. **Add Compliance Information** (for cross-border)
   - **HS Code**: Harmonized System code for customs
   - **Country of Origin**: e.g., CN for China
   - **Weight**: Product weight
   - **Dimensions**: Length × Width × Height

5. **Add Images**
   - Click **Upload Images**
   - Select images from your computer
   - Drag to reorder

6. **Save as Draft**
   - Click **Save Draft**
   - Product status will be "Draft"

### Putting Products On Sale

1. Navigate to **Products** > select product
2. Ensure stock is available (Stock > 0)
3. Click **Put On Sale** button
4. Product is now visible to customers

### Taking Products Off Sale

1. Navigate to **Products** > select product
2. Click **Take Off Sale** button
3. Product will be hidden from customers but not deleted

### Managing Product Variants (SKUs)

Products with multiple sizes/colors need variants:

1. **Enable Variants**
   - Check **This product has variants**
   - Save the product

2. **Create Variants**
   - Click **Variants** tab
   - Click **Add Variant**
   - Fill in:
     - **Code**: Unique SKU code (e.g., `SKU-001-BLK-42`)
     - **Price**: Variant price
     - **Stock**: Initial stock
     - **Attributes**: Color, Size, etc.

3. **Example Attributes**
   ```json
   {
     "Color": "Black",
     "Size": "42"
   }
   ```

### Adding Multi-Language Content

For international markets, add translations:

1. Navigate to **Products** > select product
2. Click **Localizations** tab
3. Click **Add Localization**
4. Select language (en, zh-CN, ja, de, fr)
5. Enter translated name and description

---

## Category Management

### Creating Categories

1. **Navigate to Categories**
   - Click **Categories** in sidebar

2. **Create Root Category**
   - Click **Create Category**
   - Enter **Name** (e.g., "Electronics")
   - Leave **Parent** as "None"
   - Click **Save**

3. **Create Sub-category**
   - Click **Create Category**
   - Enter **Name** (e.g., "Phones")
   - Select **Parent** category
   - Click **Save**

### Category Hierarchy Limits

```
Level 1 (Root)     Electronics
Level 2            ├── Phones
Level 3            │   ├── iPhone Cases
                   │   └── Android Cases
                   └── Computers
                       ├── Laptops
                       └── Desktops
```

> ⚠️ **Maximum 3 levels** of category depth

### Setting Category Visibility

Control which markets see each category:

1. Select category from list
2. Click **Market Visibility** tab
3. Check/uncheck markets
4. Click **Save**

### Reordering Categories

1. Navigate to **Categories**
2. Drag categories to reorder
3. Changes are saved automatically

---

## Brand Management

### Creating a Brand

1. **Navigate to Brands**
   - Click **Brands** in sidebar
   - Click **Create Brand**

2. **Fill in Information**
   - **Name**: Brand name (required)
   - **Logo**: Upload brand logo
   - **Description**: Brand story
   - **Website**: Official website URL

3. **Trademark Information** (optional)
   - **Trademark Number**: Registration number
   - **Trademark Country**: Country of registration

4. **Enable Brand Page**
   - Check **Enable Brand Page**
   - Creates a dedicated brand landing page

5. Click **Save**

### Managing Brand Market Visibility

1. Select brand from list
2. Click **Market Visibility** tab
3. Configure visibility per market
4. Click **Save**

---

## Market Configuration

### Understanding Markets

Markets represent regions where you sell products:

| Market | Code | Currency | Tax System |
|--------|------|----------|------------|
| China | CN | CNY | VAT 13% |
| United States | US | USD | No VAT |
| United Kingdom | UK | GBP | VAT 20%, IOSS |
| Germany | DE | EUR | VAT 19%, IOSS |
| France | FR | EUR | VAT 20%, IOSS |
| Australia | AU | AUD | GST 10% |

### Adding a New Market

1. **Navigate to Markets**
   - Click **Markets** in sidebar
   - Click **Create Market**

2. **Configure Market**
   - **Code**: Market code (US, UK, DE, etc.)
   - **Name**: Display name
   - **Currency**: Local currency
   - **Default Language**: Primary language

3. **Configure Tax Rules**
   ```json
   {
     "vat_rate": "0.20",
     "gst_rate": "0",
     "ioss_enabled": true,
     "include_tax": true
   }
   ```

4. Click **Save**

### Setting Default Market

1. Navigate to **Markets**
2. Click **Set as Default** on desired market
3. Default market is used for primary pricing

---

## Inventory Management

### Setting Up Warehouses

1. **Navigate to Warehouses**
   - Click **Warehouses** in sidebar
   - Click **Create Warehouse**

2. **Enter Details**
   - **Code**: Unique identifier (e.g., `WH-SH-01`)
   - **Name**: Warehouse name
   - **Country**: Location country
   - **Address**: Full address

3. **Set as Default**
   - Check **Default Warehouse**
   - Used for primary stock allocation

4. Click **Save**

### Updating Stock

1. **Navigate to Inventory**
   - Click **Inventory** in sidebar

2. **Find Product**
   - Search by SKU code or product name

3. **Update Stock**
   - Click **Update Stock**
   - Enter new available stock
   - Add remark (optional)
   - Click **Confirm**

### Adjusting Stock

For corrections (damage, loss, found items):

1. Navigate to **Inventory**
2. Click **Adjust Stock**
3. Enter:
   - **SKU Code**: Product variant
   - **Warehouse**: Location
   - **Quantity**: Positive (add) or negative (remove)
   - **Remark**: Reason for adjustment
4. Click **Confirm**

### Monitoring Low Stock

1. Navigate to **Inventory** > **Low Stock**
2. View SKUs below safety threshold
3. Update safety stock if needed

### Setting Safety Stock

Safety stock triggers low stock alerts:

1. Navigate to **Inventory**
2. Click **Safety Stock**
3. Update thresholds for each SKU
4. Click **Save**

### Viewing Inventory Logs

Track all stock changes:

1. Navigate to **Inventory** > **Logs**
2. Filter by:
   - SKU Code
   - Product
   - Change Type (manual, order, return, adjustment)
3. View:
   - Before/After stock
   - Operator
   - Timestamp
   - Remark

---

## Best Practices

### Product Naming

- Use clear, descriptive names
- Include key attributes (color, size) in variants
- Keep names under 200 characters

### Pricing

- Enter prices in cents (multiply by 100)
- Always verify currency
- Set cost prices for profit tracking

### Compliance

- Always fill HS codes for cross-border products
- Enter accurate weights and dimensions
- Declare dangerous goods properly

### Stock Management

- Set safety stock for popular items
- Regular stock audits
- Use multiple warehouses strategically

### Categories

- Don't exceed 3 levels
- Use descriptive names
- Keep categories balanced (avoid too many or too few products)

### Markets

- Configure tax rules correctly
- Test pricing in each market
- Localize product content

---

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `Ctrl + S` | Save current form |
| `Ctrl + N` | Create new item |
| `Esc` | Cancel/close modal |
| `?` | Show keyboard shortcuts |

---

## Getting Help

- **Documentation**: [docs.shopjoy.com](https://docs.shopjoy.com)
- **Support Email**: support@shopjoy.com
- **In-app Help**: Click `?` icon in top right