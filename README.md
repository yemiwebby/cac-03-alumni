# CAC-03 Alumni Birthday Reminder System

An automated WhatsApp birthday reminder system for the CAC-03 Alumni Association using Go and GitHub Actions.

## Features

- **Daily Birthday Reminders**: Automatically sends WhatsApp messages for upcoming birthdays
- **Monthly Birthday Reports**: Sends monthly summaries of all birthdays in the upcoming month
- **Duplicate Prevention**: Smart deduplication to prevent multiple messages for the same person
- **Flexible Scheduling**: GitHub Actions with manual triggers and automated scheduling
- **Target Month Override**: Ability to send monthly reports for specific months

## How It Works

The system reads birthday data from a CSV file and sends WhatsApp messages using the WhatsApp Business API when:
- **Daily**: Someone has a birthday tomorrow (configurable lookahead)
- **Monthly**: At the start of each month with all birthdays for that month

## GitHub Actions Usage

### Manual Trigger Instructions

1. **Navigate to Actions**:
   - Go to your repository: `https://github.com/yemiwebby/cac-03-alumni`
   - Click the **"Actions"** tab
   - Select **"Birthday Reminders"** workflow

2. **Run Workflow**:
   - Click **"Run workflow"** button (top right)
   - Configure the options in the dropdown form:

### Configuration Options

#### 📋 **Run Mode**
- `daily` - Send daily birthday reminders for tomorrow
- `monthly` - Send monthly birthday report

#### 🧪 **Dry Run** (Recommended for testing)
- `true` - **Test mode** (shows what would be sent without actually sending)
- `false` - **Live mode** (actually sends WhatsApp messages)

#### 📅 **Target Month** (For Monthly Reports)
- `0` - **Auto** (next month from current date) - *Default*
- `1` - January
- `2` - February
- `3` - March
- `4` - April
- `5` - May
- `6` - June
- `7` - July
- `8` - August
- `9` - September
- `10` - October
- `11` - November
- `12` - December

### Example Usage Scenarios

#### Send September Monthly Report
- **Run Mode**: `monthly`
- **Dry Run**: `true` (test first)
- **Target Month**: `9`

#### Send Daily Reminder (Test)
- **Run Mode**: `daily`
- **Dry Run**: `true`
- **Target Month**: `0` (ignored for daily)

#### Send Next Month Report (Default)
- **Run Mode**: `monthly`
- **Dry Run**: `false`
- **Target Month**: `0`

## Automated Scheduling

The system runs automatically:
- **Daily reminders**: Every day at 9:00 AM UTC (10:00 AM WAT)
- **Monthly reports**: 1st of every month at 8:00 AM UTC (9:00 AM WAT)

## Local Development

### Prerequisites
- Go 1.24+
- WhatsApp Business API credentials

### Setup
1. Clone the repository
2. Copy `.env.example` to `.env`
3. Configure your WhatsApp Business API credentials
4. Run: `go run cmd/main.go -dry` (test mode)

### Command Line Options
```bash
# Daily reminder for tomorrow (dry run)
go run cmd/main.go -dry

# Monthly report for next month (dry run)
go run cmd/main.go -monthly -dry

# Monthly report for specific month (dry run)
go run cmd/main.go -monthly -dry -target-month=9

# Live mode (actually sends messages)
go run cmd/main.go
go run cmd/main.go -monthly
```

## Configuration

### Environment Variables
- `WA_PHONE_ID`: WhatsApp Business Phone Number ID
- `WA_TOKEN`: WhatsApp Business API Access Token
- `WA_TEMPLATE`: WhatsApp template name (default: `cac_template_birthday`)
- `WA_LANG`: Language code (default: `en`)
- `WA_TO_LIST`: Comma-separated list of recipient phone numbers (E.164 format)
- `TIMEZONE`: Timezone (default: `Africa/Lagos`)

### GitHub Secrets Required
- `CSV_DATA`: Complete CSV file content (stored securely as encrypted secret)
- All WhatsApp API credentials (see environment variables above)

### CSV Data Security
🔐 **CSV data is stored as an encrypted GitHub Secret** (`CSV_DATA`) rather than in the repository files. This provides:
- **Complete privacy** - No alumni data visible in public repository
- **Secure access** - Only authorized workflows can access the data
- **Easy updates** - Update CSV by modifying the GitHub Secret
- **GitHub Actions compatibility** - File created dynamically during execution

### CSV Data Format
The system expects a CSV file with these columns:
- `FULL NAME (SCHOOL SURNAME FIRST)`
- `DATE OF BIRTH` (YYYY-MM-DD format)

## Security & Privacy

⚠️ **Important Security Notes**:
- Repository contains sensitive alumni data (names, birthdates, contact info)
- Manual workflow triggers are **restricted to repository owner only**
- WhatsApp API credentials are stored as encrypted GitHub Secrets
- All manual executions include authorization checks
- **Recommendation**: Consider making this repository private
- See [SECURITY.md](SECURITY.md) for detailed security guidelines

### Authorization
Only the repository owner can manually trigger workflows. Unauthorized attempts will be logged and blocked.

### Data Protection
- Personal data is processed according to privacy regulations
- CSV data includes full names, birth dates, and contact information
- Regular audits recommended for data accuracy and consent
- Consider implementing data retention policies

## Duplicate Prevention

The system automatically handles duplicate entries by:
- Creating unique keys based on name + birth month + birth day
- Case-insensitive name matching
- Skipping subsequent duplicate entries
- Ensuring only one message per person per birthday

## Support

For issues or questions related to the birthday reminder system, please contact the repository maintainer.

---

**CAC-03 Alumni Association** | Automated Birthday Reminder System | Powered by Go & WhatsApp Business API
