import { NextIntlClientProvider } from 'next-intl'
import { getMessages } from 'next-intl/server'
import { SiteHeader } from '@/components/layout/site-header'
import { Footer } from '@/components/layout/site-footer'
import { routing } from '@/i18n/routing'

export default async function UserLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  const messages = await getMessages()
  const locale = routing.defaultLocale

  return (
    <NextIntlClientProvider locale={locale} messages={messages}>
      <div className="relative flex min-h-screen flex-col">
        <SiteHeader />
        <main className="flex-1">{children}</main>
        <Footer />
      </div>
    </NextIntlClientProvider>
  )
}
