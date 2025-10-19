// Seed script for board members councils with proper structure
require('dotenv').config();
const { MongoClient, ObjectId } = require('mongodb');

const uri = process.env.MONGODB_URI || 'mongodb://localhost:27017';
const dbName = process.env.DB_NAME || 'LACPA';

async function seedCouncils() {
    const client = new MongoClient(uri);

    try {
        console.log('🔌 Connecting to MongoDB...');
        await client.connect();
        console.log('✅ Connected to MongoDB');

        const db = client.db(dbName);
        const councilsCollection = db.collection('councils');
        const positionsCollection = db.collection('council_positions');
        const membersCollection = db.collection('individual_members');

        // Clear existing data
        console.log('🗑️  Clearing existing councils and positions...');
        await councilsCollection.deleteMany({});
        await positionsCollection.deleteMany({});

        console.log('📝 Creating councils and board members...\n');

        // ========================================
        // CURRENT COUNCIL (2025) - Twenty Fifth Council
        // ========================================
        
        const currentCouncil = {
            _id: new ObjectId(),
            name: 'Twenty Fifth Council',
            start_date: new Date('2025-01-01'),
            end_date: new Date('2025-12-31'),
            is_active: true,
            description: 'Current serving council for 2025',
            created_at: new Date(),
            updated_at: new Date()
        };

        await councilsCollection.insertOne(currentCouncil);
        console.log('✅ Created current council: Twenty Fifth Council (2025)');

        // Create board members for current council (2025)
        const boardMembers = [
            // President
            { firstName: 'Elie', lastName: 'Georges Abboud', position: 'President', biography: 'President of the Twenty Fifth Council of LACPA for 2025, with extensive experience in accounting standards and professional development.' },
            
            // Second Row - VP, Treasurer, Secretary
            { firstName: 'Mohamed', lastName: 'Ali Blaik', position: 'Vice President', biography: 'Vice President for 2025 with over 20 years of experience in auditing and financial consulting.' },
            { firstName: 'Jean', lastName: 'Issa Rachid', position: 'Board Treasurer', biography: 'Board Treasurer for 2025 responsible for financial oversight and management of LACPA resources.' },
            { firstName: 'Nabil', lastName: 'Khoder Ghaith', position: 'Board Secretary', biography: 'Board Secretary for 2025 managing administrative affairs and council documentation.' },
            
            // Third Row - 6 Board Members
            { firstName: 'Fadi', lastName: 'Haitham El Majari', position: 'Board Member', biography: 'Board member contributing to strategic planning and member services.' },
            { firstName: 'Khalil', lastName: 'Hussein Zeidan', position: 'Board Member', biography: 'Board member with expertise in international accounting standards.' },
            { firstName: 'Khalil', lastName: 'Mohammad Ali Obeida', position: 'Board Member', biography: 'Board member focused on professional development and education.' },
            { firstName: 'Graziella', lastName: 'Antoine Hobeika', position: 'Board Member', biography: 'Board member advocating for transparency and ethics in the profession.' },
            { firstName: 'Simon', lastName: 'Georges Staii', position: 'Board Member', biography: 'Board member supporting educational initiatives and training programs.' },
            { firstName: 'Nadim', lastName: 'Antoine Daher', position: 'Board Member', biography: 'Board member with expertise in audit quality and compliance.' }
        ];

        console.log('📝 Creating individual members and assigning positions...');
        
        for (const memberData of boardMembers) {
            // Create the individual member
            const member = {
                _id: new ObjectId(),
                first_name: memberData.firstName,
                last_name: memberData.lastName,
                member_type: 'Practicing',
                lacpa_id: `LACPA-${Math.floor(10000 + Math.random() * 90000)}`,
                email: `${memberData.firstName.toLowerCase()}.${memberData.lastName.toLowerCase().replace(/ /g, '')}@lacpa.org.lb`,
                phone: `+961 ${Math.floor(1 + Math.random() * 9)} ${Math.floor(100000 + Math.random() * 900000)}`,
                city: 'Beirut',
                district: 'Beirut',
                biography: memberData.biography,
                show_phone: true,
                show_email: true,
                show_linkedin: true,
                show_address: true,
                linkedin_url: `https://linkedin.com/in/${memberData.firstName.toLowerCase()}-${memberData.lastName.toLowerCase().replace(/ /g, '-')}`,
                is_council_member: true,
                council_position: memberData.position,
                created_at: new Date(),
                updated_at: new Date()
            };

            await membersCollection.insertOne(member);

            // Create council position
            const position = {
                _id: new ObjectId(),
                member_id: member._id,
                council_id: currentCouncil._id,
                position: memberData.position,
                start_date: currentCouncil.start_date,
                end_date: currentCouncil.end_date,
                is_active: true,
                created_at: new Date(),
                updated_at: new Date()
            };

            await positionsCollection.insertOne(position);
            
            // Update member with position reference
            await membersCollection.updateOne(
                { _id: member._id },
                { 
                    $set: { 
                        current_council_position_id: position._id 
                    } 
                }
            );

            console.log(`   ✓ ${memberData.position}: ${memberData.firstName} ${memberData.lastName}`);
        }

        // ========================================
        // PAST COUNCIL 1 (2024) - Twenty Fourth Council
        // ========================================
        
        const pastCouncil1 = {
            _id: new ObjectId(),
            name: 'Twenty Fourth Council',
            start_date: new Date('2024-01-01'),
            end_date: new Date('2024-12-31'),
            is_active: false,
            description: 'Previous council term (2024)',
            created_at: new Date(),
            updated_at: new Date()
        };

        await councilsCollection.insertOne(pastCouncil1);
        console.log('\n✅ Created past council: Twenty Fourth Council (2024)');

        // Create full board for 2024 council (1 President + 1 VP + 1 Treasurer + 1 Secretary + 6 Board Members)
        const pastBoardMembers1 = [
            { firstName: 'Hassan', lastName: 'Kamal Ahmad', position: 'President', biography: 'Former President of the Twenty Fourth Council (2024), led significant reforms in professional standards.' },
            { firstName: 'Omar', lastName: 'Khalil Mansour', position: 'Vice President', biography: 'Former Vice President (2024) who contributed to modernizing LACPA operations.' },
            { firstName: 'Rania', lastName: 'Said Harb', position: 'Board Treasurer', biography: 'Former Board Treasurer (2024) who managed financial stability during challenging times.' },
            { firstName: 'Michel', lastName: 'Georges Khoury', position: 'Board Secretary', biography: 'Former Board Secretary (2024) who improved documentation and record-keeping systems.' },
            { firstName: 'Sami', lastName: 'Joseph Frem', position: 'Board Member', biography: 'Board Member (2024) specializing in audit quality assurance and regulatory compliance.' },
            { firstName: 'Diana', lastName: 'Elias Nader', position: 'Board Member', biography: 'Board Member (2024) focused on continuing education and professional development programs.' },
            { firstName: 'Marwan', lastName: 'Youssef Aoun', position: 'Board Member', biography: 'Board Member (2024) with expertise in international accounting standards implementation.' },
            { firstName: 'Rita', lastName: 'Antoine Saliba', position: 'Board Member', biography: 'Board Member (2024) dedicated to enhancing member services and communication.' },
            { firstName: 'George', lastName: 'Habib Karam', position: 'Board Member', biography: 'Board Member (2024) advocating for technological advancement in the profession.' },
            { firstName: 'Nadia', lastName: 'Ramez Moussa', position: 'Board Member', biography: 'Board Member (2024) promoting diversity and inclusion within the accounting profession.' }
        ];

        for (const memberData of pastBoardMembers1) {
            const member = {
                _id: new ObjectId(),
                first_name: memberData.firstName,
                last_name: memberData.lastName,
                member_type: 'Retired',
                lacpa_id: `LACPA-${Math.floor(10000 + Math.random() * 90000)}`,
                email: `${memberData.firstName.toLowerCase()}.${memberData.lastName.toLowerCase().replace(/ /g, '')}@lacpa.org.lb`,
                phone: `+961 ${Math.floor(1 + Math.random() * 9)} ${Math.floor(100000 + Math.random() * 900000)}`,
                city: 'Beirut',
                district: 'Beirut',
                biography: memberData.biography,
                show_phone: false,
                show_email: false,
                show_linkedin: false,
                show_address: true,
                is_council_member: false,
                created_at: new Date(),
                updated_at: new Date()
            };

            await membersCollection.insertOne(member);

            const position = {
                _id: new ObjectId(),
                member_id: member._id,
                council_id: pastCouncil1._id,
                position: memberData.position,
                start_date: pastCouncil1.start_date,
                end_date: pastCouncil1.end_date,
                is_active: false,
                created_at: new Date(),
                updated_at: new Date()
            };

            await positionsCollection.insertOne(position);
            console.log(`   ✓ ${memberData.position}: ${memberData.firstName} ${memberData.lastName}`);
        }

        // ========================================
        // PAST COUNCIL 2 (2023) - Twenty Third Council
        // ========================================
        
        const pastCouncil2 = {
            _id: new ObjectId(),
            name: 'Twenty Third Council',
            start_date: new Date('2023-01-01'),
            end_date: new Date('2023-12-31'),
            is_active: false,
            description: 'Previous council term (2023)',
            created_at: new Date(),
            updated_at: new Date()
        };

        await councilsCollection.insertOne(pastCouncil2);
        console.log('\n✅ Created past council: Twenty Third Council (2023)');

        // Create full board for 2023 council (1 President + 1 VP + 1 Treasurer + 1 Secretary + 6 Board Members)
        const pastBoardMembers2 = [
            { firstName: 'Ahmad', lastName: 'Farid Tabbara', position: 'President', biography: 'Former President (2023) who strengthened international partnerships and member engagement.' },
            { firstName: 'Laila', lastName: 'Hasan Fakhoury', position: 'Vice President', biography: 'Former Vice President (2023) known for promoting professional ethics.' },
            { firstName: 'Tony', lastName: 'Elias Chamoun', position: 'Board Treasurer', biography: 'Former Board Treasurer (2023) who ensured sound financial management.' },
            { firstName: 'Victor', lastName: 'Samir Yammine', position: 'Board Secretary', biography: 'Former Board Secretary (2023) who streamlined administrative processes.' },
            { firstName: 'Joelle', lastName: 'Kamil Farah', position: 'Board Member', biography: 'Board Member (2023) with focus on SME accounting standards and support.' },
            { firstName: 'Pierre', lastName: 'Nabil Hokayem', position: 'Board Member', biography: 'Board Member (2023) specializing in taxation and fiscal policy advocacy.' },
            { firstName: 'Maya', lastName: 'Rafic Abdo', position: 'Board Member', biography: 'Board Member (2023) championing digital transformation initiatives.' },
            { firstName: 'Elias', lastName: 'Boutros Maalouf', position: 'Board Member', biography: 'Board Member (2023) promoting sustainability reporting and ESG standards.' },
            { firstName: 'Carla', lastName: 'Elie Haddad', position: 'Board Member', biography: 'Board Member (2023) focused on public sector accounting reforms.' },
            { firstName: 'Fadi', lastName: 'Maroun Geagea', position: 'Board Member', biography: 'Board Member (2023) enhancing collaboration with international accounting bodies.' }
        ];

        for (const memberData of pastBoardMembers2) {
            const member = {
                _id: new ObjectId(),
                first_name: memberData.firstName,
                last_name: memberData.lastName,
                member_type: 'Retired',
                lacpa_id: `LACPA-${Math.floor(10000 + Math.random() * 90000)}`,
                email: `${memberData.firstName.toLowerCase()}.${memberData.lastName.toLowerCase().replace(/ /g, '')}@lacpa.org.lb`,
                phone: `+961 ${Math.floor(1 + Math.random() * 9)} ${Math.floor(100000 + Math.random() * 900000)}`,
                city: 'Beirut',
                district: 'Beirut',
                biography: memberData.biography,
                show_phone: false,
                show_email: false,
                show_linkedin: false,
                show_address: true,
                is_council_member: false,
                created_at: new Date(),
                updated_at: new Date()
            };

            await membersCollection.insertOne(member);

            const position = {
                _id: new ObjectId(),
                member_id: member._id,
                council_id: pastCouncil2._id,
                position: memberData.position,
                start_date: pastCouncil2.start_date,
                end_date: pastCouncil2.end_date,
                is_active: false,
                created_at: new Date(),
                updated_at: new Date()
            };

            await positionsCollection.insertOne(position);
            console.log(`   ✓ ${memberData.position}: ${memberData.firstName} ${memberData.lastName}`);
        }

        console.log('\n🎉 Database seeding completed!');
        console.log('\n📊 Summary:');
        console.log('   - Current Council (2025): 1 President + 1 VP + 1 Treasurer + 1 Secretary + 6 Board Members (10 total)');
        console.log('   - Past Council (2024): Twenty Fourth Council - Full roster (10 members)');
        console.log('   - Past Council (2023): Twenty Third Council - Full roster (10 members)');
        console.log('   - Each council serves a 1-year term');
        console.log('   - Total: 30 board members across 3 councils');
        console.log('\n💡 You can now view board members at: http://localhost:3000/discover/board-of-directors');

    } catch (error) {
        console.error('❌ Error seeding database:', error);
    } finally {
        await client.close();
        console.log('\n🔌 MongoDB connection closed');
    }
}

// Run the seed function
seedCouncils();
