# Authorization Ownership and Session Token Contract

Role authorization has one owner: the Role Authority in BC-Governance resolves a user's role assignments into a single authoritative role name, and BC-Identity asks for it via an internal call at login/SSO instead of reading role tables directly. The session token contract is locked to BC-Identity as the only issuer: the token carries user identity and the authoritative role name (user_id, username, ole, exp, iat, HS256); every other context only consumes it and must not define its own claim structure.

**Considered Options**

- Shared Go JWT package imported by both services (rejected: the two services are independent Go modules; a shared package would create a new shared code module to maintain for ~20 lines of validation, and map-based consumption already tolerates claim evolution).
- Per-service claim structs (rejected: the copies had already drifted - the level claim existed only in user-service).

**Consequences**

- admin-service token code shrinks to a consumer-only validator; its GenerateToken and Claims struct are removed.
- Level is removed from the token (display-only attribute) and the RequireLevel middleware is deleted.
- Once user-service no longer reads schema_admin, DB-FIX-03 can downgrade the user-service database account to schema_auth-only.
