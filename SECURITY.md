# Security Policy

## Reporting a Vulnerability

Use [the Vulnerability Disclosure issue template](https://gitlab.com/l0nax/typact/-/issues/new?issuable_template=Vulnerability%20Disclosure) to report a new security vulnerability.

New security issue should follow these guidelines:

- Create new issues as `confidential` if unsure whether issue a potential
vulnerability or not. It is easier to make an issue that should have been
public open than to remediate an issue that should have been confidential.
Consider adding the `/confidential` quick action to a project issue template.

- Always label as ``~security`` at a minimum. If you're reporting a vulnerability (or something you suspect may possibly be one) please use the [Vulnerability Disclosure](https://gitlab.com/l0nax/typact/-/issues/new?issuable_template=Vulnerability%20Disclosure) template while creating the issue. Otherwise, follow the steps here (with a security label).

- Add any additional labels you know apply. Additional labels will be applied
by the security team and other engineering personnel, but it will help with
the triage process:

  - [`~"type::bug"`, or `~"type::feature"` if appropriate](https://about.gitlab.com/handbook/engineering/security/security-engineering-and-research/application-security/vulnerability-management.html#vulnerability-vs-feature-vs-bug)

- Issues that contain customer specific data, such as private repository contents,
should be assigned `~keep confidential`.

Occasionally, data that should remain confidential, such as the private
project contents of a user that reported an issue, may get included in an
issue. If necessary, a sanitized issue may need to be created with more
general discussion and examples appropriate for public disclosure prior to
release.
