# Security policy

`go-razdel` is a pure-Go tokenizer. It does not open network sockets or parse a wire protocol. The realistic issues are panics, hangs, or unexpected resource use on unusual Unicode or invalid UTF-8.

## Report a vulnerability

Use [GitHub private vulnerability reporting](https://docs.github.com/en/code-security/security-advisories/guidance-on-reporting-and-writing-information-about-vulnerabilities/privately-reporting-a-security-vulnerability) on this repository if it is enabled.

If that is not available, email the maintainer listed on [github.com/muonsoft](https://github.com/muonsoft) or open an issue **without** a crashing payload in the title.

Please include:

- Go version and OS
- a **minimal** input string (file attachment is fine)
- whether the failure is in `Tokenize`, `Sentenize`, or both

Do not attach production data.

Public API functions are not supposed to panic on user input, including invalid UTF-8. A panic is always in scope for this policy.
