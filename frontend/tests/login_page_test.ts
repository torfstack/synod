import {expect, test} from '@playwright/test';

test('login page has button with "login"', async ({page}) => {
    await page.goto('/');
    await expect(page.getByRole('button', {
        name: 'Login',
        exact: true
    })).toBeVisible()
});

test('login page has button with "login with google"', async ({page}) => {
    await page.goto('/');
    await expect(page.getByRole('button', {
        name: 'login with google',
    })).toBeVisible()
});
