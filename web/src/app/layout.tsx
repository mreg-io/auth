import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Authentication | My Registry',
  description:
    'Access and manage your software supply chain efficiently with My Registry. Log in to your account to track, control, and secure your software assets. New to My Registry? Sign up today to leverage our powerful SaaS platform designed to streamline your software supply chain operations. Enjoy seamless integration, advanced security features, and real-time insights. Join now and take control of your software supply chain with My Registry.',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>{children}</body>
    </html>
  );
}
