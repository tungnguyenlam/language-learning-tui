## 2026-05-10: Fix E2E Test Failure - New Content Visibility

- Fixed the `test_new_decks_visibility` E2E test that was failing due to a mismatch in deck name.
- The test was searching for "German Comprehensive" but the actual deck name is "German Comprehensive (A1-B2)".
- Updated the test expectation to match the correct deck name in `e2e_tests/test_new_content_visibility.py`.
- After the fix, all 149 E2E tests pass and the full verification suite (`./scripts/verify.sh`) executes successfully.
- The application now launches without errors, all views render correctly, core user interactions respond as expected, and state is persisted to SQLite.