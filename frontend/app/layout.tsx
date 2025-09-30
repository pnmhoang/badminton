import './globals.css'
import { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Badminton Tournament Manager',
  description: 'Manage badminton players, matches and tournaments',
}

interface RootLayoutProps {
  children: React.ReactNode
}

export default function RootLayout({ children }: RootLayoutProps) {
  return (
    <html lang="en">
      <body>
        {children}
      </body>
    </html>
  )
}