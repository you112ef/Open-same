# Contributing to Open-Same

Thank you for your interest in contributing to Open-Same! This document provides guidelines and information for contributors.

## 🤝 How to Contribute

### Types of Contributions

We welcome various types of contributions:

- **Bug Reports** - Report bugs and issues
- **Feature Requests** - Suggest new features
- **Code Contributions** - Submit code improvements
- **Documentation** - Improve or add documentation
- **Testing** - Help test and improve quality
- **Community Support** - Help other users

### Getting Started

1. **Fork the Repository**
   ```bash
   git clone https://github.com/your-username/open-same.git
   cd open-same
   ```

2. **Set Up Development Environment**
   ```bash
   # Backend (Go)
   cd backend
   go mod download
   
   # Frontend (React)
   cd frontend
   npm install
   
   # SDK
   cd sdk
   npm install
   ```

3. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/your-bug-fix
   ```

## 🐛 Reporting Issues

### Before Reporting

- Check if the issue has already been reported
- Search existing issues and discussions
- Try to reproduce the issue with the latest version

### Issue Template

Use our issue template and include:

- **Clear Description** - What happened vs. what you expected
- **Steps to Reproduce** - Detailed steps to recreate the issue
- **Environment Details** - OS, browser, version, etc.
- **Screenshots/Logs** - Visual evidence if applicable
- **Additional Context** - Any other relevant information

## 💻 Code Contributions

### Code Style Guidelines

#### Go (Backend)
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting
- Run `go vet` and `golint` before submitting
- Write tests for new functionality

#### TypeScript/React (Frontend)
- Follow [Airbnb JavaScript Style Guide](https://github.com/airbnb/javascript)
- Use Prettier for formatting
- Run ESLint before submitting
- Write tests using Jest and React Testing Library

#### General
- Write clear, descriptive commit messages
- Keep commits focused and atomic
- Add comments for complex logic
- Follow existing naming conventions

### Pull Request Process

1. **Create a Pull Request**
   - Use a descriptive title
   - Reference related issues
   - Include a detailed description

2. **Code Review**
   - Address review comments
   - Make requested changes
   - Ensure all tests pass

3. **Merge Requirements**
   - All tests must pass
   - Code review approved
   - Documentation updated if needed

### Commit Message Format

Use conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Examples:
- `feat(auth): add OAuth 2.0 support`
- `fix(api): resolve rate limiting issue`
- `docs(readme): update installation instructions`

## 🧪 Testing

### Running Tests

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test

# SDK tests
cd sdk
npm test
```

### Writing Tests

- Write tests for new functionality
- Ensure good test coverage
- Use descriptive test names
- Mock external dependencies

## 📚 Documentation

### Documentation Standards

- Use clear, concise language
- Include code examples
- Keep documentation up-to-date
- Follow existing formatting patterns

### Documentation Areas

- API documentation
- User guides
- Developer guides
- Deployment instructions
- Troubleshooting guides

## 🔒 Security

### Security Issues

- **DO NOT** open public issues for security vulnerabilities
- Email security issues to: security@open-same.dev
- We'll acknowledge receipt within 48 hours
- We'll provide updates on progress

### Security Best Practices

- Never commit sensitive information
- Use environment variables for secrets
- Follow security guidelines in code
- Report potential security issues

## 🏷️ Labels and Milestones

### Issue Labels

- `bug` - Something isn't working
- `enhancement` - New feature or request
- `documentation` - Documentation improvements
- `good first issue` - Good for newcomers
- `help wanted` - Extra attention needed
- `priority: high/medium/low` - Issue priority

### Milestones

- `v1.0.0` - Initial release
- `v1.1.0` - Feature releases
- `v1.0.x` - Bug fix releases

## 🚀 Release Process

### Release Cycle

1. **Feature Freeze** - Stop adding new features
2. **Testing Phase** - Comprehensive testing
3. **Release Candidate** - Tag RC versions
4. **Final Release** - Tag stable version
5. **Documentation Update** - Update docs and changelog

### Versioning

We follow [Semantic Versioning](https://semver.org/):

- `MAJOR.MINOR.PATCH`
- `1.0.0` - Initial release
- `1.1.0` - New features, backward compatible
- `1.0.1` - Bug fixes, backward compatible

## 🤝 Community Guidelines

### Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Provide constructive feedback
- Help others learn and grow

### Communication Channels

- **GitHub Issues** - Bug reports and feature requests
- **GitHub Discussions** - General questions and discussions
- **Discord** - Real-time chat and community
- **Email** - Security and private matters

## 📋 Checklist for Contributors

Before submitting your contribution, ensure:

- [ ] Code follows style guidelines
- [ ] Tests are written and passing
- [ ] Documentation is updated
- [ ] Commit messages are clear
- [ ] Pull request description is detailed
- [ ] Related issues are referenced
- [ ] No sensitive information is included

## 🎉 Recognition

Contributors are recognized through:

- GitHub contributors list
- Release notes and changelog
- Contributor hall of fame
- Special acknowledgments for major contributions

## 📞 Getting Help

### Questions and Support

- Check existing documentation
- Search existing issues
- Ask in GitHub Discussions
- Join our Discord community

### Mentorship

- Look for issues labeled `good first issue`
- Ask questions in discussions
- Request help from maintainers
- Join community events

## 📄 License

By contributing to Open-Same, you agree that your contributions will be licensed under the [MIT License](LICENSE).

---

Thank you for contributing to Open-Same! Your contributions help make this platform better for everyone. 🚀