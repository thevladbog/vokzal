/**
 * Session storage key for refresh token. Used to restore session on page reload.
 * Stored in sessionStorage (same-tab only) so it survives refresh but not new tabs.
 */
export const REFRESH_TOKEN_STORAGE_KEY = 'vokzal_refresh_token';
