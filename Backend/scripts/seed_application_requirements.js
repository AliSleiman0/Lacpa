const { MongoClient, ObjectId } = require('mongodb');
require('dotenv').config();

async function seedApplicationRequirements() {
    const client = new MongoClient(process.env.MONGO_URI);

    try {
        console.log('🔌 Connecting to MongoDB...');
        await client.connect();
        console.log('✅ Connected to MongoDB');

        const db = client.db('LACPA');
        const requirementsCollection = db.collection('application_requirements');

        // Clear existing requirements
        console.log('🗑️  Clearing existing requirements...');
        await requirementsCollection.deleteMany({});

        console.log('📝 Creating application requirements...\n');

        // ========================================
        // INDIVIDUAL MEMBERSHIP REQUIREMENTS
        // ========================================

        const individualRequirements = [
            {
                _id: new ObjectId(),
                title: 'Lebanese Residency',
                description: 'Must have been a Lebanese national for at least ten years.',
                icon: 'fas fa-id-card',
                application_type: 'Individual',
                is_required: true,
                order_index: 1,
                created_at: new Date(),
                updated_at: new Date()
            },
            {
                _id: new ObjectId(),
                title: 'Age Requirement',
                description: 'Must be at least 25 years of age.',
                icon: 'fas fa-calendar',
                application_type: 'Individual',
                is_required: true,
                order_index: 2,
                created_at: new Date(),
                updated_at: new Date()
            },
            {
                _id: new ObjectId(),
                title: 'Clean Legal Record',
                description: 'Must enjoy full civil rights and must not have been convicted of a felony or a disgraceful misdemeanor, in accordance with Paragraph 3 of Article 53 of Law No. 53/91 regulating the profession.',
                icon: 'fas fa-shield-alt',
                application_type: 'Individual',
                is_required: true,
                order_index: 3,
                created_at: new Date(),
                updated_at: new Date()
            },
            {
                _id: new ObjectId(),
                title: 'Educational Qualification',
                description: 'Must hold a Bachelor\'s degree in Business Administration or its equivalent, and the Lebanese Certified Public Accountant (CPA) degree issued by the Ministry of National Education, or a certificate considered equivalent issued by the Civil Service Council.',
                icon: 'fas fa-graduation-cap',
                application_type: 'Individual',
                is_required: true,
                order_index: 4,
                created_at: new Date(),
                updated_at: new Date()
            },
            {
                _id: new ObjectId(),
                title: 'Advanced Degree',
                description: 'The applicant must hold an advanced specialized degree in accounting, either: • a local degree recognized by the Association as equivalent, or • a foreign certificate at the level of.',
                icon: 'fas fa-user-graduate',
                application_type: 'Individual',
                is_required: true,
                order_index: 5,
                created_at: new Date(),
                updated_at: new Date()
            },
            {
                _id: new ObjectId(),
                title: 'Annual Duties',
                description: 'Must have completed annual accounting duties and evaluations.',
                icon: 'fas fa-clipboard-check',
                application_type: 'Individual',
                is_required: true,
                order_index: 6,
                created_at: new Date(),
                updated_at: new Date()
            },
            {
                _id: new ObjectId(),
                title: 'Training Completion',
                description: 'Must have completed the training period stipulated by law, and submit a certificate from the training supervisor proving successful completion (unless holding a certificate that grants exemption from the training period).',
                icon: 'fas fa-certificate',
                application_type: 'Individual',
                is_required: true,
                order_index: 7,
                created_at: new Date(),
                updated_at: new Date()
            },
            {
                _id: new ObjectId(),
                title: 'No Professional Conflict',
                description: 'Must not have practiced another profession in the public sector for reasons other than the internship, unless authorized by the competent authority.',
                icon: 'fas fa-users',
                application_type: 'Individual',
                is_required: true,
                order_index: 8,
                created_at: new Date(),
                updated_at: new Date()
            },
            {
                _id: new ObjectId(),
                title: 'No Disciplinary Actions',
                description: 'Must not have been removed from the register for disciplinary reasons or barred from practicing the profession.',
                icon: 'fas fa-ban',
                application_type: 'Individual',
                is_required: true,
                order_index: 9,
                created_at: new Date(),
                updated_at: new Date()
            }
        ];

        for (const req of individualRequirements) {
            await requirementsCollection.insertOne(req);
            console.log(`   ✓ Individual: ${req.title}`);
        }

        // ========================================
        // FIRM MEMBERSHIP REQUIREMENTS
        // ========================================

        const firmRequirements = [
            {
                _id: new ObjectId(),
                title: 'Business Registration',
                description: 'Firm must be legally registered with Lebanese authorities and hold a valid commercial registration certificate.',
                icon: 'fas fa-certificate',
                application_type: 'Firm',
                is_required: true,
                order_index: 1,
                created_at: new Date(),
                updated_at: new Date()
            },
            {
                _id: new ObjectId(),
                title: 'Professional Licensing',
                description: 'All partners and key personnel must hold valid professional accounting licenses or certifications.',
                icon: 'fas fa-file-contract',
                application_type: 'Firm',
                is_required: true,
                order_index: 2,
                created_at: new Date(),
                updated_at: new Date()
            },
            {
                _id: new ObjectId(),
                title: 'Minimum Partners',
                description: 'Firm must have at least one qualified partner who is a licensed accountant with minimum 5 years experience.',
                icon: 'fas fa-users',
                application_type: 'Firm',
                is_required: true,
                order_index: 3,
                created_at: new Date(),
                updated_at: new Date()
            },
            {
                _id: new ObjectId(),
                title: 'Professional Indemnity Insurance',
                description: 'Maintain adequate professional liability insurance coverage as per LACPA requirements.',
                icon: 'fas fa-shield-alt',
                application_type: 'Firm',
                is_required: true,
                order_index: 4,
                created_at: new Date(),
                updated_at: new Date()
            },
            {
                _id: new ObjectId(),
                title: 'Quality Control System',
                description: 'Implement and maintain quality control procedures in accordance with international auditing standards.',
                icon: 'fas fa-tasks',
                application_type: 'Firm',
                is_required: true,
                order_index: 5,
                created_at: new Date(),
                updated_at: new Date()
            },
            {
                _id: new ObjectId(),
                title: 'Office Infrastructure',
                description: 'Maintain a physical office in Lebanon with proper facilities and resources for professional practice.',
                icon: 'fas fa-building',
                application_type: 'Firm',
                is_required: true,
                order_index: 6,
                created_at: new Date(),
                updated_at: new Date()
            }
        ];

        for (const req of firmRequirements) {
            await requirementsCollection.insertOne(req);
            console.log(`   ✓ Firm: ${req.title}`);
        }

        console.log('\n🎉 Database seeding completed!');
        console.log('\n📊 Summary:');
        console.log(`   - Individual Requirements: ${individualRequirements.length} cards`);
        console.log(`   - Firm Requirements: ${firmRequirements.length} cards`);
        console.log(`   - Total Requirements: ${individualRequirements.length + firmRequirements.length} cards`);
        console.log('\n💡 You can now view the application page at: http://localhost:3000/membership/apply-now');

    } catch (error) {
        console.error('❌ Error seeding database:', error);
        process.exit(1);
    } finally {
        console.log('\n🔌 MongoDB connection closed');
        await client.close();
    }
}

seedApplicationRequirements();
