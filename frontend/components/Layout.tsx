'use client';

import { ReactNode, useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import type { User } from '../types';

interface LayoutProps {
  children: ReactNode;
  title?: string;
}

export default function Layout({ children, title = 'Badminton Tournament' }: LayoutProps) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    // Get user from localStorage
    const userData = localStorage.getItem('user');
    if (userData) {
      try {
        const parsedUser = JSON.parse(userData);
        setUser(parsedUser);
      } catch (error) {
        console.error('Failed to parse user data:', error);
      }
    }
    setIsLoading(false);
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    router.push('/auth/login');
  };

  if (isLoading) {
    return <div>Loading...</div>;
  }

  return (
    <>
      <header className="header">
        <nav className="navbar">
          <div className="container">
            <div className="nav-brand">
              <h2>🏸 {title}</h2>
            </div>
            <div className="nav-menu">
              {user ? (
                <>
                  <a href="/dashboard">Dashboard</a>
                  <a href="/tournaments">Tournaments</a>
                  {user.role === 'admin' && (
                    <a href="/admin">Admin Panel</a>
                  )}
                  <div className="user-menu">
                    <span className="user-greeting">Hi, {user.full_name}</span>
                    <a href="/profile">Profile</a>
                    <button onClick={handleLogout} className="logout-btn">
                      Logout
                    </button>
                  </div>
                </>
              ) : (
                <div className="auth-menu">
                  <a href="/auth/login">Login</a>
                  <a href="/auth/register">Register</a>
                </div>
              )}
            </div>
          </div>
        </nav>
      </header>
      
      <main className="main-content">
        {children}
      </main>
      
      <footer className="footer">
        <div className="container">
          <p>&copy; 2025 Badminton Tournament Management</p>
        </div>
      </footer>
      
      <style jsx>{`
        .header {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          color: white;
          padding: 1rem 0;
          margin-bottom: 2rem;
          box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        }
        
        .navbar {
          display: flex;
          justify-content: space-between;
          align-items: center;
        }
        
        .nav-brand h2 {
          margin: 0;
          font-size: 1.5rem;
          font-weight: 700;
        }
        
        .nav-menu {
          display: flex;
          align-items: center;
          gap: 1rem;
        }
        
        .nav-menu a {
          color: white;
          text-decoration: none;
          padding: 0.5rem 1rem;
          border-radius: 6px;
          transition: all 0.2s;
          font-weight: 500;
        }
        
        .nav-menu a:hover {
          background: rgba(255, 255, 255, 0.1);
          transform: translateY(-1px);
        }
        
        .user-menu {
          display: flex;
          align-items: center;
          gap: 1rem;
          margin-left: 1rem;
          padding-left: 1rem;
          border-left: 1px solid rgba(255, 255, 255, 0.2);
        }
        
        .user-greeting {
          font-size: 0.875rem;
          opacity: 0.9;
          font-weight: 500;
        }
        
        .logout-btn {
          background: rgba(255, 255, 255, 0.1);
          color: white;
          border: 1px solid rgba(255, 255, 255, 0.3);
          padding: 0.5rem 1rem;
          border-radius: 6px;
          cursor: pointer;
          font-size: 0.875rem;
          font-weight: 500;
          transition: all 0.2s;
        }
        
        .logout-btn:hover {
          background: rgba(255, 255, 255, 0.2);
          transform: translateY(-1px);
        }
        
        .auth-menu {
          display: flex;
          gap: 1rem;
        }
        
        .main-content {
          min-height: calc(100vh - 200px);
          padding: 2rem 0;
        }
        
        .footer {
          background: #f8f9fa;
          text-align: center;
          padding: 1.5rem 0;
          margin-top: 3rem;
          border-top: 1px solid #e5e5e5;
          color: #666;
        }
        
        @media (max-width: 768px) {
          .navbar {
            flex-direction: column;
            gap: 1rem;
          }
          
          .nav-menu {
            flex-wrap: wrap;
            justify-content: center;
          }
          
          .user-menu {
            margin-left: 0;
            padding-left: 0;
            border-left: none;
            border-top: 1px solid rgba(255, 255, 255, 0.2);
            padding-top: 1rem;
            flex-wrap: wrap;
            justify-content: center;
          }
        }
      `}</style>
    </>
  );
}