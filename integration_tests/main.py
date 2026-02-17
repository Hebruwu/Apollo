import asyncio
import os
import sys
import traceback

import asyncpg
import httpx

BASE_URL = os.getenv("BASE_URL", "http://localhost:8080")
DATABASE_URL = os.getenv("TEST_DATABASE_URL", "postgresql://user:password@localhost:5432/relational")


async def clean_users():
    conn = await asyncpg.connect(DATABASE_URL)
    try:
        await conn.execute("DELETE FROM users")
    finally:
        await conn.close()


async def query_user(username: str) -> asyncpg.Record | None:
    conn = await asyncpg.connect(DATABASE_URL)
    try:
        return await conn.fetchrow("SELECT * FROM users WHERE username = $1", username)
    finally:
        await conn.close()


# --- Tests ---


async def test_register_returns_created_on_valid_request():
    await clean_users()
    async with httpx.AsyncClient() as client:
        resp = await client.post(f"{BASE_URL}/api/v1/users/register", json={
            "username": "newuser",
            "email": "new@example.com",
            "password": "securepassword",
        })
    assert resp.status_code == 201, f"expected 201, got {resp.status_code}"
    body = resp.json()
    assert body["success"] == "User created", f"unexpected body: {body}"


async def test_register_stores_user_in_database():
    await clean_users()
    async with httpx.AsyncClient() as client:
        await client.post(f"{BASE_URL}/api/v1/users/register", json={
            "username": "dbuser",
            "email": "db@example.com",
            "password": "securepassword",
        })
    row = await query_user("dbuser")
    assert row is not None, "user not found in database"
    assert row["username"] == "dbuser", f"expected 'dbuser', got '{row['username']}'"
    assert row["email"] == "db@example.com", f"expected 'db@example.com', got '{row['email']}'"


async def test_register_stores_hashed_password_not_plaintext():
    await clean_users()
    async with httpx.AsyncClient() as client:
        await client.post(f"{BASE_URL}/api/v1/users/register", json={
            "username": "hashuser",
            "email": "hash@example.com",
            "password": "myplaintextpassword",
        })
    row = await query_user("hashuser")
    assert row is not None, "user not found in database"
    password_hash = bytes(row["password_hash"])
    assert password_hash != b"myplaintextpassword", "password stored as plaintext"
    assert len(password_hash) == 32, f"expected 32-byte hash, got {len(password_hash)}"


async def test_register_stores_sixteen_byte_salt():
    await clean_users()
    async with httpx.AsyncClient() as client:
        await client.post(f"{BASE_URL}/api/v1/users/register", json={
            "username": "saltuser",
            "email": "salt@example.com",
            "password": "securepassword",
        })
    row = await query_user("saltuser")
    assert row is not None, "user not found in database"
    salt = bytes(row["salt"])
    assert len(salt) == 16, f"expected 16-byte salt, got {len(salt)}"


async def test_register_returns_conflict_on_duplicate_username():
    await clean_users()
    async with httpx.AsyncClient() as client:
        await client.post(f"{BASE_URL}/api/v1/users/register", json={
            "username": "duplicate",
            "email": "first@example.com",
            "password": "password1",
        })
        resp = await client.post(f"{BASE_URL}/api/v1/users/register", json={
            "username": "duplicate",
            "email": "second@example.com",
            "password": "password2",
        })
    assert resp.status_code == 409, f"expected 409, got {resp.status_code}"
    body = resp.json()
    assert body["error"] == "Username already exists", f"unexpected body: {body}"


async def test_register_does_not_overwrite_existing_user_on_duplicate():
    await clean_users()
    async with httpx.AsyncClient() as client:
        await client.post(f"{BASE_URL}/api/v1/users/register", json={
            "username": "keeper",
            "email": "original@example.com",
            "password": "password1",
        })
        await client.post(f"{BASE_URL}/api/v1/users/register", json={
            "username": "keeper",
            "email": "overwrite@example.com",
            "password": "password2",
        })
    row = await query_user("keeper")
    assert row is not None, "user not found in database"
    assert row["email"] == "original@example.com", f"expected 'original@example.com', got '{row['email']}'"


async def test_register_returns_bad_request_on_malformed_json():
    await clean_users()
    async with httpx.AsyncClient() as client:
        resp = await client.post(
            f"{BASE_URL}/api/v1/users/register",
            content=b"not valid json",
            headers={"Content-Type": "application/json"},
        )
    assert resp.status_code == 400, f"expected 400, got {resp.status_code}"


async def test_register_returns_json_content_type():
    await clean_users()
    async with httpx.AsyncClient() as client:
        resp = await client.post(f"{BASE_URL}/api/v1/users/register", json={
            "username": "contenttypeuser",
            "email": "ct@example.com",
            "password": "securepassword",
        })
    assert resp.status_code == 201, f"expected 201, got {resp.status_code}"
    assert "application/json" in resp.headers.get("content-type", ""), \
        f"expected application/json, got '{resp.headers.get('content-type')}'"


# --- Runner ---

TESTS = [
    test_register_returns_created_on_valid_request,
    test_register_stores_user_in_database,
    test_register_stores_hashed_password_not_plaintext,
    test_register_stores_sixteen_byte_salt,
    test_register_returns_conflict_on_duplicate_username,
    test_register_does_not_overwrite_existing_user_on_duplicate,
    test_register_returns_bad_request_on_malformed_json,
    test_register_returns_json_content_type,
]


async def main():
    passed = 0
    failed = 0

    for test in TESTS:
        name = test.__name__
        try:
            await test()
            print(f"  PASS  {name}")
            passed += 1
        except Exception:
            print(f"  FAIL  {name}")
            traceback.print_exc()
            failed += 1

    print(f"\n{passed} passed, {failed} failed, {len(TESTS)} total")

    if failed > 0:
        sys.exit(1)


if __name__ == "__main__":
    asyncio.run(main())
