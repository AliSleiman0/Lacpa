/**
 * Clear Test Members Script
 * Run with: node scripts/clear_test_members.js
 * 
 * This script removes all test/seed data from the members collections
 */

const { MongoClient } = require('mongodb');
require('dotenv').config();

// MongoDB connection URL from environment or default
const MONGO_URI = process.env.MONGO_URI || 'mongodb://localhost:27017';
const DB_NAME = process.env.MONGO_DATABASE || 'lacpa';

async function clearTestMembers() {
    const client = new MongoClient(MONGO_URI);

    try {
        // Connect to MongoDB
        console.log('🔌 Connecting to MongoDB...');
        await client.connect();
        console.log('✅ Connected to MongoDB');

        const db = client.db(DB_NAME);
        const individualsCollection = db.collection('individual_members');
        const firmsCollection = db.collection('firm_members');

        // Count before deletion
        const individualsCount = await individualsCollection.countDocuments({});
        const firmsCount = await firmsCollection.countDocuments({});

        console.log('\n📊 Current counts:');
        console.log(`   - Individual members: ${individualsCount}`);
        console.log(`   - Firm members: ${firmsCount}`);

        // Ask for confirmation
        console.log('\n⚠️  This will DELETE ALL members from the database!');
        console.log('   Press Ctrl+C to cancel, or wait 3 seconds to continue...\n');

        await new Promise(resolve => setTimeout(resolve, 3000));

        // Delete all members
        console.log('🗑️  Deleting individual members...');
        const individualsResult = await individualsCollection.deleteMany({});
        console.log(`✅ Deleted ${individualsResult.deletedCount} individual members`);

        console.log('\n🗑️  Deleting firm members...');
        const firmsResult = await firmsCollection.deleteMany({});
        console.log(`✅ Deleted ${firmsResult.deletedCount} firm members`);

        console.log('\n🎉 Database cleanup completed!');
        console.log('\n💡 You can now run the seed scripts to add fresh test data:');
        console.log('   node scripts/seed_individuals.js');
        console.log('   node scripts/seed_firms.js');

    } catch (error) {
        console.error('❌ Error clearing database:', error);
        process.exit(1);
    } finally {
        // Close connection
        await client.close();
        console.log('\n🔌 Disconnected from MongoDB');
    }
}

// Run the cleanup
clearTestMembers();
