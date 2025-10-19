/**
 * Seed Individual Members Script
 * Run with: node scripts/seed_individuals.js
 * 
 * This script populates the individual_members collection with dummy data
 * for testing pagination and filtering.
 */

const { MongoClient, ObjectId } = require('mongodb');
require('dotenv').config();

// MongoDB connection URL from environment or default
const MONGO_URI = process.env.MONGO_URI || 'mongodb://localhost:27017';
const DB_NAME = process.env.MONGO_DATABASE || 'lacpa';

// Member types for testing
const MEMBER_TYPES = ['Apprentices', 'Practicing', 'Non-Practicing', 'Retired'];

// Lebanese first names
const FIRST_NAMES = [
    'Ali', 'Hassan', 'Fatima', 'Zeinab', 'Mohamed', 'Ahmad', 'Layla', 'Sara',
    'Karim', 'Nadia', 'Omar', 'Rania', 'Tarek', 'Maya', 'Walid', 'Dina',
    'Bilal', 'Lara', 'Sami', 'Hala', 'Fadi', 'Nour', 'Rami', 'Lea',
    'Georges', 'Marie', 'Pierre', 'Rita', 'Antoine', 'Carla', 'Elie', 'Joelle',
    'Michel', 'Nicole', 'Joseph', 'Mira', 'Tony', 'Nathalie', 'Charbel', 'Giselle'
];

// Lebanese last names
const LAST_NAMES = [
    'Hariri', 'Khoury', 'Haddad', 'Salameh', 'Nasrallah', 'Gemayel', 'Frangieh', 'Jumblatt',
    'El Khoury', 'Abboud', 'Gholam', 'Semaan', 'Najjar', 'Khalil', 'Hobeika', 'Geagea',
    'Aoun', 'Berri', 'Mikati', 'Salam', 'Karami', 'Hamadeh', 'Arslan', 'Murr',
    'Sleiman', 'Skaff', 'Tueini', 'Moawad', 'Chamoun', 'Eddé', 'Obeid', 'El Haj'
];

// Lebanese cities
const CITIES = ['Beirut', 'Tripoli', 'Sidon', 'Tyre', 'Zahle', 'Jounieh', 'Baabda', 'Nabatieh'];
const DISTRICTS = ['Achrafieh', 'Hamra', 'Verdun', 'Raouche', 'Hazmieh', 'Sin El Fil', 'Jdeideh'];

// Generate random member
function generateMember(index) {
    const firstName = FIRST_NAMES[Math.floor(Math.random() * FIRST_NAMES.length)];
    const lastName = LAST_NAMES[Math.floor(Math.random() * LAST_NAMES.length)];
    const memberType = MEMBER_TYPES[Math.floor(Math.random() * MEMBER_TYPES.length)];
    const city = CITIES[Math.floor(Math.random() * CITIES.length)];
    const district = DISTRICTS[Math.floor(Math.random() * DISTRICTS.length)];
    
    return {
        _id: new ObjectId(),
        lacpa_id: `IND-${new Date().getFullYear()}-${String(index + 1).padStart(4, '0')}`,
        first_name: firstName,
        last_name: lastName,
        middle_name: Math.random() > 0.5 ? FIRST_NAMES[Math.floor(Math.random() * FIRST_NAMES.length)] : "",
        member_type: memberType,
        membership_status: "Active",
        membership_tier: ["Bronze", "Silver", "Gold", "Platinum"][Math.floor(Math.random() * 4)],
        
        // Contact info
        phone: `+961 ${Math.floor(Math.random() * 90) + 10} ${Math.floor(Math.random() * 900000) + 100000}`,
        email: `${firstName.toLowerCase()}.${lastName.toLowerCase()}@example.com`,
        
        // Address
        city: city,
        district: district,
        governorate: ["Beirut", "Mount Lebanon", "North", "South", "Beqaa"][Math.floor(Math.random() * 5)],
        country: "Lebanon",
        postal_code: String(Math.floor(Math.random() * 9000) + 1000),
        
        // Professional info
        job_title: ["Accountant", "Auditor", "Tax Consultant", "Financial Analyst", "CFO"][Math.floor(Math.random() * 5)],
        company_name: Math.random() > 0.3 ? LAST_NAMES[Math.floor(Math.random() * LAST_NAMES.length)] + " & Associates" : "",
        years_of_experience: Math.floor(Math.random() * 30) + 1,
        
        // Education
        highest_degree: ["Bachelor", "Master", "PhD", "CPA"][Math.floor(Math.random() * 4)],
        university: ["AUB", "LAU", "USJ", "BAU", "NDU"][Math.floor(Math.random() * 5)],
        graduation_year: 2025 - Math.floor(Math.random() * 30),
        
        // Certifications
        certifications: Math.random() > 0.5 ? ["CPA", "CMA", "CIA"].slice(0, Math.floor(Math.random() * 3) + 1) : [],
        
        // Social media
        linkedin_url: `https://linkedin.com/in/${firstName.toLowerCase()}-${lastName.toLowerCase()}`,
        twitter_url: Math.random() > 0.7 ? `https://twitter.com/${firstName.toLowerCase()}${lastName.toLowerCase()}` : "",
        
        // Biography
        biography: `${firstName} ${lastName} is a dedicated ${memberType.toLowerCase()} member with expertise in accounting and financial management. ${Math.random() > 0.5 ? 'Specialized in audit and tax compliance.' : 'Focused on corporate finance and strategic planning.'}`,
        
        // Dates
        date_of_birth: new Date(1950 + Math.floor(Math.random() * 50), Math.floor(Math.random() * 12), Math.floor(Math.random() * 28) + 1),
        membership_start_date: new Date(2000 + Math.floor(Math.random() * 25), Math.floor(Math.random() * 12), 1),
        
        // Privacy settings
        show_phone: Math.random() > 0.3,
        show_email: Math.random() > 0.2,
        show_address: Math.random() > 0.4,
        show_social_media: Math.random() > 0.5,
        
        // Status
        is_active: true,
        is_verified: true,
        
        // Metadata
        created_at: new Date(),
        updated_at: new Date(),
    };
}

// Main seeding function
async function seedIndividuals() {
    const client = new MongoClient(MONGO_URI);

    try {
        // Connect to MongoDB
        console.log('🔌 Connecting to MongoDB...');
        await client.connect();
        console.log('✅ Connected to MongoDB');

        const db = client.db(DB_NAME);
        const individualsCollection = db.collection('individual_members');

        // Optional: Clear existing members (uncomment if you want to reset)
        console.log('🗑️  Clearing existing individual members...');
        await individualsCollection.deleteMany({});

        // Generate members
        const numberOfMembers = 50; // Generate 50 members for testing pagination
        const members = [];
        
        for (let i = 0; i < numberOfMembers; i++) {
            members.push(generateMember(i));
        }

        // Insert members
        console.log(`\n📝 Inserting ${numberOfMembers} individual members...`);
        const result = await individualsCollection.insertMany(members);

        console.log(`\n✅ Successfully inserted ${result.insertedCount} individual members!`);
        
        // Count by type
        const apprenticesCount = members.filter(m => m.member_type === 'Apprentices').length;
        const practicingCount = members.filter(m => m.member_type === 'Practicing').length;
        const nonPracticingCount = members.filter(m => m.member_type === 'Non-Practicing').length;
        const retiredCount = members.filter(m => m.member_type === 'Retired').length;
        
        console.log('\n📊 Breakdown:');
        console.log(`   - Apprentices: ${apprenticesCount}`);
        console.log(`   - Practicing: ${practicingCount}`);
        console.log(`   - Non-Practicing: ${nonPracticingCount}`);
        console.log(`   - Retired: ${retiredCount}`);
        
        console.log('\n🎉 Database seeding completed!');
        console.log('\n💡 You can now view members at: http://localhost:3000/membership');

    } catch (error) {
        console.error('❌ Error seeding database:', error);
        process.exit(1);
    } finally {
        // Close connection
        await client.close();
        console.log('\n🔌 Disconnected from MongoDB');
    }
}

// Run the seeding
seedIndividuals();
