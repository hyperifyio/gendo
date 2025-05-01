GitHub Copilot Code Review Instructions

When reviewing code with GitHub Copilot, please ensure all code follows these fundamental principles. Every piece of code must be implemented as a small, independent module with clear boundaries and responsibilities. Each module should have a single, well-defined purpose and minimal dependencies on other modules.

All code must be accompanied by comprehensive unit tests. Test files should be written before implementing features, following test-driven development practices. Each test file must cover all code paths, including edge cases and error conditions. Use table-driven tests for multiple test cases and mock external dependencies appropriately.

Code organization is critical. Keep implementations modular and focused, following language-specific best practices and idioms. Use meaningful package names and maintain a clear file structure. Each module should be self-contained and easily testable in isolation.

Error handling must be thorough and explicit. Test cases should verify invalid inputs, malformed data, network failures, resource limitations, and edge cases in data processing. All error conditions must be properly documented and tested.

Documentation is essential. All code must include clear descriptions of its purpose, expected inputs and outputs, error conditions, and any special setup requirements. Comments should explain why code is written a certain way, not what it does.

When reviewing code, ensure that all dependencies are properly managed and that the code follows the project's established patterns and conventions. Code should be maintainable, readable, and follow the principle of least surprise.

Integration tests are required for any code that interacts with external systems or other modules. These tests should verify the correct behavior of the system as a whole, not just individual components.

Performance considerations should be taken into account for any code that processes data or handles user interactions. Include benchmarks for performance-critical operations and ensure that the code scales appropriately.

Security is paramount. All code must be reviewed for potential security vulnerabilities, especially when handling user input or interacting with external systems. Follow the principle of least privilege and implement proper input validation and sanitization.

Remember that code quality is not just about functionality but also about maintainability, readability, and reliability. Every line of code should be written with these principles in mind. 