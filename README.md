# GiftExchanger

Self-hostable Secret Santa for friends, families, and small groups.

Create an exchange, invite people with a link, collect wishlists, and assign gifters in private. Assignments stay on your own server — no third-party service needs to see who drew whom.

## Features

- Create a gift exchange and share an invite link
- Participants join, leave, and keep a simple wishlist
- Private Secret Santa assignments — you only see your own recipient
- Run the whole app locally with Docker

## Email

Login codes, magic links, and invites are sent with [Resend](https://resend.com). Without configuration, the API prints them to its log instead, which is all you need for local development.

To send real email:

1. Create an API key in Resend.
2. Add and verify a sending domain in Resend. It walks you through the SPF and DKIM DNS records.
3. Set these in `.env` (keep `MAIL_FROM` quoted):

   ```sh
   RESEND_API_KEY=re_...
   MAIL_FROM="GiftExchanger <hello@your-verified-domain.com>"
   ```

Without a verified domain, Resend only allows sending from `onboarding@resend.dev` to your own Resend account's address, which is fine for a quick test.

## Status

GiftExchanger is in early development. Local setup and a public demo are coming.

