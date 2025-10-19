# Database Seeding Scripts

This directory contains scripts to populate the LACPA database with dummy/test data.

## Prerequisites

Make sure you have:
- Node.js installed
- MongoDB running (locally or remote)
- The `.env` file configured in the Backend directory

## Seed Individual Members Data

To populate the `individual_members` collection with 50 dummy members:

### Step 1: Install dependencies (if not already installed)

```bash
cd Backend
npm install mongodb dotenv
```

### Step 2: Run the seed script

```bash
node scripts/seed_individuals.js
```

### Expected Output

```
🔌 Connecting to MongoDB...
✅ Connected to MongoDB

📝 Inserting 50 individual members...

✅ Successfully inserted 50 individual members!

📊 Breakdown:
   - Apprentices: 15
   - Practicing: 16
   - Non-Practicing: 9
   - Retired: 10

🎉 Database seeding completed!

💡 You can now view members at: http://localhost:3000/membership
```

## Seed Firms Data

To populate the `firm_members` collection with 10 dummy firms (4 Big 4, 2 Large, 2 Medium, 2 Small):

### Step 1: Install dependencies (if not already installed)

```bash
cd Backend
npm install mongodb dotenv
```

### Step 2: Run the seed script

```bash
node scripts/seed_firms.js
```

### Expected Output

```
🔌 Connecting to MongoDB...
✅ Connected to MongoDB

📝 Inserting firm members...

✅ Successfully inserted 10 firm members!

📊 Breakdown:
   - Big 4 firms: 4
   - Large firms: 2
   - Medium firms: 2
   - Small firms: 2

🎉 Database seeding completed!

💡 You can now view firms at: http://localhost:3000/membership/firms

🔌 Disconnected from MongoDB
```

## Firms Included

All firms include logo images stored in `LACPA_Web/assets/images/firms/`.

### Big 4 Firms
1. **Deloitte Middle East** - Full-service audit, tax, and consulting
   - Logo: `/assets/images/firms/deloitte.svg`
2. **PwC Lebanon** - Assurance, tax, advisory, and deals
   - Logo: `/assets/images/firms/pwc.svg`
3. **Ernst & Young (EY) Lebanon** - Assurance, tax, transactions, consulting
   - Logo: `/assets/images/firms/ey.svg`
4. **KPMG Lebanon** - Audit, tax, advisory, technology
   - Logo: `/assets/images/firms/kpmg.svg`

### Large Firms
5. **BDO Lebanon** - Mid-tier accounting and advisory
   - Logo: `/assets/images/firms/bdo.svg`
6. **Grant Thornton Lebanon** - Business advisory to dynamic organizations
   - Logo: `/assets/images/firms/grant-thornton.svg`

### Medium Firms
7. **Semaan, Gholam & Co.** - SME-focused accounting firm
   - Logo: `/assets/images/firms/semaan-gholam.png`
8. **Fiduciaire du Liban SAL** - Comprehensive accounting services
   - Logo: `/assets/images/firms/fiduciaire.png`

### Small Firms
9. **Khalil & Associates** - Boutique consultancy for startups
   - Logo: `/assets/images/firms/khalil.png`
10. **Najjar CPA Firm** - Local CPA serving individuals and small businesses
    - Logo: `/assets/images/firms/najjar.png`

## Customization

### Modifying Firm Data
To modify the data:
1. Open `seed_firms.js`
2. Edit the `firms` array
3. Run the script again

### Updating Firm Logos
Logos are stored in `LACPA_Web/assets/images/firms/`. To update a logo:

1. **Replace with your own logo:**
   - Add your logo file to `LACPA_Web/assets/images/firms/`
   - Supported formats: SVG (recommended), PNG, JPG
   - Recommended size: 200x80px or similar aspect ratio
   - Update the `logo_url` in the seed script

2. **Example:**
   ```javascript
   logo_url: "/assets/images/firms/your-firm-logo.svg"
   ```

3. **Logo display:**
   - Logos are displayed in a white container (128px × 96px)
   - Images are scaled to fit while maintaining aspect ratio
   - SVG format recommended for crisp rendering at any size

## Clearing Data

To clear existing firms before seeding, uncomment this line in `seed_firms.js`:

```javascript
// await firmsCollection.deleteMany({});
```

## Environment Variables

The script uses these environment variables from `.env`:

- `MONGO_URI` - MongoDB connection string (default: `mongodb://localhost:27017`)
- `MONGO_DATABASE` - Database name (default: `lacpa`)

## Troubleshooting

**Error: Cannot find module 'mongodb'**
```bash
npm install mongodb
```

**Error: Cannot connect to MongoDB**
- Check if MongoDB is running
- Verify `MONGO_URI` in your `.env` file
- Ensure firewall allows MongoDB connections

**Error: Authentication failed**
- Check MongoDB username/password in `MONGO_URI`
- Verify database permissions

## Testing

After seeding, verify the data:

1. **Via Browser:**
   - Navigate to: http://localhost:3000/membership/firms
   - You should see all 10 firms displayed

2. **Via MongoDB:**
   ```bash
   mongosh
   use lacpa
   db.firm_members.countDocuments()  // Should return 10
   db.firm_members.find({firm_size: "Big 4"}).count()  // Should return 4
   ```

## Next Steps

After seeding:
- Test pagination (firms page shows 12 per page by default)
- Test filtering by firm size
- Test firm profile views
- Add individual members and associate them with firms
