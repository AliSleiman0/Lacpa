// Seed script for board members councils
require('dotenv').config();
const { MongoClient } = require('mongodb');

const uri = process.env.MONGODB_URI || 'mongodb://localhost:27017';
const dbName = process.env.DB_NAME || 'LACPA';

async function seedBoardMembers() {
    const client = new MongoClient(uri);

    try {
        console.log('🔌 Connecting to MongoDB...');
        await client.connect();
        console.log('✅ Connected to MongoDB');

        const db = client.db(dbName);
        const councilsCollection = db.collection('councils');

        // Clear existing councils
        console.log('🗑️  Clearing existing councils...');
        await councilsCollection.deleteMany({});

        console.log('📝 Inserting council data...\n');

        // Current Council (2025) - Twenty Third Council
        const currentCouncil = {
            title: 'Twenty Third Council',
            start_date: '19/06/2022',
            end_date: '21/06/2024',
            is_current: true,
            president: {
                name: 'Elie Georges Abboud',
                position: 'President',
                photo: '', // Add path if available
                biography: 'President of the Twenty Third Council of LACPA, serving from 2022 to 2024.'
            },
            other_members: [
                {
                    name: 'Mohamed Ali Blaik',
                    position: 'Vice President',
                    photo: '',
                    biography: 'Vice President with extensive experience in accounting and auditing.'
                },
                {
                    name: 'Jean Issa Rachid',
                    position: 'Board Treasurer',
                    photo: '',
                    biography: 'Board Treasurer responsible for financial oversight.'
                },
                {
                    name: 'Nabil Khoder Ghaith',
                    position: 'Board Secretary',
                    photo: '',
                    biography: 'Board Secretary managing administrative affairs.'
                },
                {
                    name: 'Fadi Haitham El Majari',
                    position: 'Board Members',
                    photo: '',
                    biography: 'Board member contributing to strategic decisions.'
                },
                {
                    name: 'Khalil Hussein Zeidan',
                    position: 'Board Members',
                    photo: '',
                    biography: 'Board member with expertise in financial regulations.'
                },
                {
                    name: 'Khalil Mohammad Ali Obeida',
                    position: 'Board Members',
                    photo: '',
                    biography: 'Board member focused on professional development.'
                },
                {
                    name: 'Graziella Antoine Hobeika',
                    position: 'Board Members',
                    photo: '',
                    biography: 'Board member advocating for member services.'
                },
                {
                    name: 'Simon Georges Staii',
                    position: 'Board Members',
                    photo: '',
                    biography: 'Board member supporting educational initiatives.'
                },
                {
                    name: 'Nadim Antoine Daher',
                    position: 'Board Members',
                    photo: '',
                    biography: 'Board member with expertise in audit standards.'
                }
            ],
            created_at: new Date(),
            updated_at: new Date()
        };

        // Twenty Second Council (2018-2022)
        const secondCouncil = {
            title: 'Twenty Second Council',
            start_date: '04/04/2018',
            end_date: '05/04/2022',
            is_current: false,
            president: {
                name: 'Hassan Kamal Ahmad',
                position: 'President',
                photo: '',
                biography: 'Former President of the Twenty Second Council.'
            },
            other_members: [
                {
                    name: 'Member Name 1',
                    position: 'Vice President',
                    photo: '',
                    biography: 'Board member of the Twenty Second Council.'
                },
                {
                    name: 'Member Name 2',
                    position: 'Board Treasurer',
                    photo: '',
                    biography: 'Board member of the Twenty Second Council.'
                },
                {
                    name: 'Member Name 3',
                    position: 'Board Secretary',
                    photo: '',
                    biography: 'Board member of the Twenty Second Council.'
                }
            ],
            created_at: new Date(),
            updated_at: new Date()
        };

        // Twenty First Council (2016-2018)
        const thirdCouncil = {
            title: 'Twenty First Council',
            start_date: '08/04/2016',
            end_date: '03/04/2018',
            is_current: false,
            president: {
                name: 'Ahmad Farid Tabbara',
                position: 'President',
                photo: '',
                biography: 'Former President of the Twenty First Council.'
            },
            other_members: [
                {
                    name: 'Member Name 1',
                    position: 'Vice President',
                    photo: '',
                    biography: 'Board member of the Twenty First Council.'
                },
                {
                    name: 'Member Name 2',
                    position: 'Board Treasurer',
                    photo: '',
                    biography: 'Board member of the Twenty First Council.'
                },
                {
                    name: 'Member Name 3',
                    position: 'Board Secretary',
                    photo: '',
                    biography: 'Board member of the Twenty First Council.'
                }
            ],
            created_at: new Date(),
            updated_at: new Date()
        };

        // Insert all councils
        const result = await councilsCollection.insertMany([
            currentCouncil,
            secondCouncil,
            thirdCouncil
        ]);

        console.log(`✅ Successfully inserted ${result.insertedCount} councils!`);
        console.log('\n📊 Council Summary:');
        console.log('   - Current Council: Twenty Third Council (2022-2024)');
        console.log('   - Past Council: Twenty Second Council (2018-2022)');
        console.log('   - Past Council: Twenty First Council (2016-2018)');
        console.log('\n🎉 Database seeding completed!');
        console.log('\n💡 You can now view board members at: http://localhost:3000/discover/board-of-directors');

    } catch (error) {
        console.error('❌ Error seeding database:', error);
    } finally {
        await client.close();
        console.log('\n🔌 MongoDB connection closed');
    }
}

// Run the seed function
seedBoardMembers();
