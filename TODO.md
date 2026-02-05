# Todo List - Fix Signup Backend Connection Issue

## Issues Identified:
1. Frontend: `handleSignup()` doesn't capture email/password input values
2. Backend: CORS preflight (OPTIONS) requests aren't handled

## Fixes to Implement:
- [x] Fix Signuppage.tsx - Add React state for email/password
- [x] Fix main.go - Add CORS middleware to handle OPTIONS requests
- [x] Test the signup functionality

## Files to Edit:
1. frontend/src/Components/Login-Signup/Signuppage.tsx
2. backend/cmd/server/main.go

## Summary:
Both issues have been fixed:
1. Frontend now properly captures email/password using React state
2. Backend now handles CORS preflight (OPTIONS) requests properly

