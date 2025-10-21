package utils

// OTPEmailTemplate returns a beautifully designed HTML email template for OTP
func OTPEmailTemplate(otp, recipientName string) string {
	// If no name provided, use generic greeting
	if recipientName == "" {
		recipientName = "User"
	}

	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Your OTP Code</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            background-color: #f5f5f5;
            padding: 20px;
        }
        .email-container {
            max-width: 600px;
            margin: 0 auto;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border-radius: 16px;
            overflow: hidden;
            box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
        }
        .email-header {
            background: linear-gradient(135deg, #0ea5e9 0%, #0284c7 100%);
            padding: 40px 30px;
            text-align: center;
            border-bottom: 4px solid rgba(255, 255, 255, 0.2);
        }
        .logo-container {
            margin-bottom: 20px;
        }
        .logo-text {
            font-size: 32px;
            font-weight: bold;
            color: #ffffff;
            text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.2);
            letter-spacing: 1px;
        }
        .header-subtitle {
            color: rgba(255, 255, 255, 0.95);
            font-size: 16px;
            margin-top: 8px;
        }
        .email-body {
            background-color: #ffffff;
            padding: 40px 30px;
        }
        .greeting {
            font-size: 24px;
            color: #1e293b;
            margin-bottom: 20px;
            font-weight: 600;
        }
        .message {
            color: #475569;
            font-size: 16px;
            line-height: 1.8;
            margin-bottom: 30px;
        }
        .otp-container {
            background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
            border: 2px dashed #0ea5e9;
            border-radius: 12px;
            padding: 30px;
            text-align: center;
            margin: 30px 0;
        }
        .otp-label {
            font-size: 14px;
            color: #64748b;
            text-transform: uppercase;
            letter-spacing: 1px;
            margin-bottom: 12px;
            font-weight: 600;
        }
        .otp-code {
            font-size: 48px;
            font-weight: bold;
            color: #0ea5e9;
            letter-spacing: 8px;
            font-family: 'Courier New', monospace;
            text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.05);
            user-select: all;
            -webkit-user-select: all;
            -moz-user-select: all;
            -ms-user-select: all;
        }
        .otp-validity {
            margin-top: 15px;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 8px;
        }
        .clock-icon {
            font-size: 18px;
        }
        .validity-text {
            color: #ef4444;
            font-size: 14px;
            font-weight: 600;
        }
        .info-box {
            background-color: #fef3c7;
            border-left: 4px solid #f59e0b;
            padding: 16px;
            border-radius: 8px;
            margin: 25px 0;
        }
        .info-box-text {
            color: #92400e;
            font-size: 14px;
            line-height: 1.6;
        }
        .security-tips {
            margin-top: 30px;
            padding: 20px;
            background-color: #f8fafc;
            border-radius: 8px;
        }
        .security-title {
            font-size: 16px;
            color: #1e293b;
            font-weight: 600;
            margin-bottom: 12px;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .security-list {
            list-style: none;
            padding: 0;
        }
        .security-list li {
            color: #64748b;
            font-size: 14px;
            padding: 6px 0;
            padding-left: 24px;
            position: relative;
        }
        .security-list li:before {
            content: "✓";
            color: #10b981;
            font-weight: bold;
            position: absolute;
            left: 0;
        }
        .email-footer {
            background-color: #1e293b;
            padding: 30px;
            text-align: center;
        }
        .footer-text {
            color: #94a3b8;
            font-size: 14px;
            margin-bottom: 15px;
        }
        .footer-links {
            margin: 20px 0;
        }
        .footer-link {
            color: #60a5fa;
            text-decoration: none;
            margin: 0 15px;
            font-size: 14px;
        }
        .footer-link:hover {
            color: #93c5fd;
            text-decoration: underline;
        }
        .social-links {
            margin-top: 20px;
        }
        .social-link {
            display: inline-block;
            margin: 0 8px;
            color: #94a3b8;
            text-decoration: none;
            font-size: 20px;
        }
        .copyright {
            color: #64748b;
            font-size: 12px;
            margin-top: 15px;
        }
        @media only screen and (max-width: 600px) {
            .email-container {
                border-radius: 0;
            }
            .email-body {
                padding: 30px 20px;
            }
            .otp-code {
                font-size: 36px;
                letter-spacing: 4px;
            }
            .greeting {
                font-size: 20px;
            }
        }
    </style>
</head>
<body>
    <div class="email-container">
        <!-- Header -->
        <div class="email-header">
            <div class="logo-container">
                <div class="logo-text">LACPA</div>
            </div>
            <div class="header-subtitle">Lebanese Association of Certified Public Accountants</div>
        </div>

        <!-- Body -->
        <div class="email-body">
            <div class="greeting">Hello ` + recipientName + `,</div>
            
            <div class="message">
                We received a request to verify your email address. Please use the One-Time Password (OTP) below to complete your verification.
            </div>

            <div class="otp-container">
                <div class="otp-label">Your Verification Code</div>
                <div class="otp-code">` + otp + `</div>
                <div class="otp-validity">
                    <span class="clock-icon">⏰</span>
                    <span class="validity-text">Expires in 5 minutes</span>
                </div>
            </div>

            <div class="info-box">
                <div class="info-box-text">
                    <strong>⚠️ Important:</strong> If you didn't request this code, please ignore this email. Your account is secure, and no action is needed.
                </div>
            </div>

            <div class="security-tips">
                <div class="security-title">
                    🔒 Security Tips
                </div>
                <ul class="security-list">
                    <li>Never share your OTP with anyone, including LACPA staff</li>
                    <li>LACPA will never ask for your OTP via phone or email</li>
                    <li>This code is valid for one-time use only</li>
                    <li>Always verify the sender's email address</li>
                </ul>
            </div>

            <div class="message" style="margin-top: 30px;">
                If you have any questions or need assistance, please don't hesitate to contact our support team.
            </div>
        </div>

        <!-- Footer -->
        <div class="email-footer">
            <div class="footer-text">
                This is an automated message from LACPA. Please do not reply to this email.
            </div>
            
            <div class="footer-links">
                <a href="#" class="footer-link">Visit Website</a>
                <a href="#" class="footer-link">Contact Support</a>
                <a href="#" class="footer-link">Privacy Policy</a>
            </div>

            <div class="social-links">
                <a href="#" class="social-link">📘</a>
                <a href="#" class="social-link">🐦</a>
                <a href="#" class="social-link">💼</a>
                <a href="#" class="social-link">📷</a>
            </div>

            <div class="copyright">
                © 2025 Lebanese Association of Certified Public Accountants. All rights reserved.
            </div>
        </div>
    </div>
</body>
</html>`
}
