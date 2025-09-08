# 🔒 Security & Privacy Notice

## Data Protection

This repository contains sensitive alumni information. Please be aware of the following:

### Personal Data Included
- ✅ Full names
- ✅ Birth dates  
- ✅ Contact information (email, phone)
- ✅ Addresses

### Security Measures Implemented

#### 🛡️ **Access Control**
- Manual workflow triggers restricted to repository owner only
- Authorization checks prevent unauthorized execution
- GitHub Secrets protect WhatsApp API credentials

#### 🔐 **Workflow Security**
- Minimal required permissions (`contents: read`, `actions: read`)
- Explicit authorization verification for manual triggers
- Clear logging of authorized users in workflow runs

#### 📱 **WhatsApp API Security**
- All credentials stored as encrypted GitHub Secrets
- No API tokens exposed in code or logs
- Recipient lists controlled via environment variables

### Recommendations for Repository Owner

#### 🚨 **Immediate Actions**
1. **Consider making repository private** if possible
2. **Review collaborator access** regularly
3. **Monitor workflow execution logs** for unauthorized attempts
4. **Regularly rotate WhatsApp API tokens**

#### 📋 **Best Practices**
- Always use **dry run mode** for testing
- Review recipient lists before live execution
- Monitor WhatsApp message delivery for anomalies
- Keep CSV data updated with consent from alumni

#### 🔄 **Regular Maintenance**
- [ ] Review and update alumni data quarterly
- [ ] Audit GitHub repository access monthly
- [ ] Verify WhatsApp API token validity
- [ ] Check for duplicate entries in CSV data

### Emergency Procedures

#### 🚨 **In Case of Unauthorized Access**
1. Immediately revoke all GitHub Secrets
2. Generate new WhatsApp API tokens
3. Review repository access logs
4. Consider making repository private
5. Notify affected alumni if data was compromised

#### 📞 **Contact Information**
For security concerns or data protection issues, contact the repository owner immediately.

---

⚠️ **Important**: This system processes personal data. Ensure compliance with local data protection regulations and obtain proper consent from alumni before processing their information.
