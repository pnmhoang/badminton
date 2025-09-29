import './globals.css'

export const metadata = {
  title: 'Badminton Tournament Manager',
  description: 'Manage badminton players, matches and tournaments',
}

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        {children}
      </body>
    </html>
  )
}