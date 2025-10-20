// Test script for Authentication System
// Run with: node scripts/test_auth.js

const BASE_URL = 'http://localhost:3000/api/auth';

let testEmail = `test${Date.now()}@example.com`;
let testLACPAID = '';
let testToken = '';
let testResetToken = '';
let testOTP = ''; // You'll need to get this from console logs

const colors = {
    reset: '\x1b[0m',
    green: '\x1b[32m',
    red: '\x1b[31m',
    yellow: '\x1b[33m',
    blue: '\x1b[36m'
};

function log(message, color = colors.reset) {
    console.log(`${color}${message}${colors.reset}`);
}

async function makeRequest(endpoint, method = 'GET', body = null, token = null) {
    const options = {
        method,
        headers: {
            'Content-Type': 'application/json',
        }
    };

    if (token) {
        options.headers['Authorization'] = `Bearer ${token}`;
    }

    if (body) {
        options.body = JSON.stringify(body);
    }

    try {
        const response = await fetch(`${BASE_URL}${endpoint}`, options);
        const data = await response.json();
        return { status: response.status, data };
    } catch (error) {
        return { status: 0, error: error.message };
    }
}

async function test1_Signup() {
    log('\n1. Testing Signup...', colors.blue);
    
    const result = await makeRequest('/signup', 'POST', {
        full_name: 'Test User',
        email: testEmail,
        password: 'Test123!@#'
    });

    if (result.status === 201 && result.data.success) {
        testLACPAID = result.data.data.lacpa_id;
        log(`✓ Signup successful! LACPA ID: ${testLACPAID}`, colors.green);
        log(`  Email: ${testEmail}`, colors.green);
        log('  Check console for OTP!', colors.yellow);
        return true;
    } else {
        log(`✗ Signup failed: ${result.data.error || result.error}`, colors.red);
        return false;
    }
}

async function test2_LoginBeforeVerification() {
    log('\n2. Testing Login Before Email Verification...', colors.blue);
    
    const result = await makeRequest('/login', 'POST', {
        lacpa_id: testLACPAID,
        password: 'Test123!@#'
    });

    if (result.status === 403 && result.data.error.includes('verify')) {
        log('✓ Correctly blocked login before verification', colors.green);
        return true;
    } else {
        log('✗ Should not allow login before verification', colors.red);
        return false;
    }
}

async function test3_VerifyOTP() {
    log('\n3. Testing OTP Verification...', colors.blue);
    
    // Prompt user to enter OTP from console
    const readline = require('readline').createInterface({
        input: process.stdin,
        output: process.stdout
    });

    return new Promise((resolve) => {
        readline.question('Enter the 6-digit OTP from console: ', async (otp) => {
            readline.close();
            
            const result = await makeRequest('/verify-otp', 'POST', {
                email: testEmail,
                otp: otp.trim()
            });

            if (result.status === 200 && result.data.success) {
                testResetToken = result.data.data.reset_token;
                log('✓ OTP verified successfully', colors.green);
                resolve(true);
            } else {
                log(`✗ OTP verification failed: ${result.data.error || result.error}`, colors.red);
                resolve(false);
            }
        });
    });
}

async function test4_LoginAfterVerification() {
    log('\n4. Testing Login After Verification...', colors.blue);
    
    const result = await makeRequest('/login', 'POST', {
        lacpa_id: testLACPAID,
        password: 'Test123!@#'
    });

    if (result.status === 200 && result.data.success) {
        testToken = result.data.data.token;
        log('✓ Login successful!', colors.green);
        log(`  Token: ${testToken.substring(0, 30)}...`, colors.green);
        return true;
    } else {
        log(`✗ Login failed: ${result.data.error || result.error}`, colors.red);
        return false;
    }
}

async function test5_GetProfile() {
    log('\n5. Testing Get Profile (Protected Route)...', colors.blue);
    
    const result = await makeRequest('/profile', 'GET', null, testToken);

    if (result.status === 200 && result.data.success) {
        log('✓ Profile retrieved successfully', colors.green);
        log(`  Name: ${result.data.data.full_name}`, colors.green);
        log(`  Email: ${result.data.data.email}`, colors.green);
        log(`  Role: ${result.data.data.role}`, colors.green);
        return true;
    } else {
        log(`✗ Get profile failed: ${result.data.error || result.error}`, colors.red);
        return false;
    }
}

async function test6_GetProfileWithoutToken() {
    log('\n6. Testing Get Profile Without Token...', colors.blue);
    
    const result = await makeRequest('/profile', 'GET');

    if (result.status === 401) {
        log('✓ Correctly blocked unauthorized access', colors.green);
        return true;
    } else {
        log('✗ Should not allow access without token', colors.red);
        return false;
    }
}

async function test7_ForgotPassword() {
    log('\n7. Testing Forgot Password...', colors.blue);
    
    const result = await makeRequest('/forgot-password', 'POST', {
        email: testEmail
    });

    if (result.status === 200 && result.data.success) {
        log('✓ Forgot password request successful', colors.green);
        log('  Check console for new OTP!', colors.yellow);
        return true;
    } else {
        log(`✗ Forgot password failed: ${result.data.error || result.error}`, colors.red);
        return false;
    }
}

async function test8_ResendOTP() {
    log('\n8. Testing Resend OTP...', colors.blue);
    
    const result = await makeRequest('/resend-otp', 'POST', {
        email: testEmail
    });

    if (result.status === 200 && result.data.success) {
        log('✓ OTP resent successfully', colors.green);
        log('  Check console for new OTP!', colors.yellow);
        return true;
    } else {
        log(`✗ Resend OTP failed: ${result.data.error || result.error}`, colors.red);
        return false;
    }
}

async function test9_VerifyResetOTP() {
    log('\n9. Testing Verify Reset OTP...', colors.blue);
    
    const readline = require('readline').createInterface({
        input: process.stdin,
        output: process.stdout
    });

    return new Promise((resolve) => {
        readline.question('Enter the 6-digit OTP for password reset: ', async (otp) => {
            readline.close();
            
            const result = await makeRequest('/verify-otp', 'POST', {
                email: testEmail,
                otp: otp.trim()
            });

            if (result.status === 200 && result.data.success) {
                testResetToken = result.data.data.reset_token;
                log('✓ Reset OTP verified successfully', colors.green);
                log(`  Reset Token: ${testResetToken.substring(0, 30)}...`, colors.green);
                resolve(true);
            } else {
                log(`✗ Reset OTP verification failed: ${result.data.error || result.error}`, colors.red);
                resolve(false);
            }
        });
    });
}

async function test10_ResetPassword() {
    log('\n10. Testing Reset Password...', colors.blue);
    
    const result = await makeRequest('/reset-password', 'POST', {
        token: testResetToken,
        new_password: 'NewTest123!@#'
    });

    if (result.status === 200 && result.data.success) {
        log('✓ Password reset successful', colors.green);
        return true;
    } else {
        log(`✗ Password reset failed: ${result.data.error || result.error}`, colors.red);
        return false;
    }
}

async function test11_LoginWithNewPassword() {
    log('\n11. Testing Login With New Password...', colors.blue);
    
    const result = await makeRequest('/login', 'POST', {
        lacpa_id: testLACPAID,
        password: 'NewTest123!@#'
    });

    if (result.status === 200 && result.data.success) {
        log('✓ Login with new password successful!', colors.green);
        return true;
    } else {
        log(`✗ Login with new password failed: ${result.data.error || result.error}`, colors.red);
        return false;
    }
}

async function test12_Logout() {
    log('\n12. Testing Logout...', colors.blue);
    
    const result = await makeRequest('/logout', 'POST', null, testToken);

    if (result.status === 200 && result.data.success) {
        log('✓ Logout successful', colors.green);
        return true;
    } else {
        log(`✗ Logout failed: ${result.data.error || result.error}`, colors.red);
        return false;
    }
}

async function runTests() {
    log('='.repeat(60), colors.yellow);
    log('LACPA Authentication System Test Suite', colors.yellow);
    log('='.repeat(60), colors.yellow);
    log(`\nBase URL: ${BASE_URL}`, colors.blue);
    log(`Test Email: ${testEmail}`, colors.blue);

    const results = [];

    // Run all tests
    results.push(await test1_Signup());
    results.push(await test2_LoginBeforeVerification());
    results.push(await test3_VerifyOTP());
    results.push(await test4_LoginAfterVerification());
    results.push(await test5_GetProfile());
    results.push(await test6_GetProfileWithoutToken());
    results.push(await test7_ForgotPassword());
    results.push(await test8_ResendOTP());
    results.push(await test9_VerifyResetOTP());
    results.push(await test10_ResetPassword());
    results.push(await test11_LoginWithNewPassword());
    results.push(await test12_Logout());

    // Summary
    log('\n' + '='.repeat(60), colors.yellow);
    log('Test Summary', colors.yellow);
    log('='.repeat(60), colors.yellow);
    
    const passed = results.filter(r => r).length;
    const total = results.length;
    
    log(`\nTotal Tests: ${total}`, colors.blue);
    log(`Passed: ${passed}`, colors.green);
    log(`Failed: ${total - passed}`, colors.red);
    log(`Success Rate: ${((passed / total) * 100).toFixed(1)}%\n`, colors.blue);

    if (passed === total) {
        log('🎉 All tests passed!', colors.green);
    } else {
        log('❌ Some tests failed. Check the output above.', colors.red);
    }
}

// Check if server is running
async function checkServer() {
    try {
        const response = await fetch('http://localhost:3000/api/health');
        if (response.ok) {
            return true;
        }
    } catch (error) {
        return false;
    }
    return false;
}

// Main execution
(async () => {
    const serverRunning = await checkServer();
    
    if (!serverRunning) {
        log('❌ Server is not running on http://localhost:3000', colors.red);
        log('Please start the server with: cd Backend && air', colors.yellow);
        process.exit(1);
    }

    await runTests();
})();
