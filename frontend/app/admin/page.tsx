'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Layout from '../../components/Layout';
import { authAPI } from '../../utils/api';
import type { User } from '../../types';

interface UserRowProps {
  user: User;
  onRoleChange: (userId: number, newRole: 'player' | 'admin') => void;
  currentUser: User;
}

function UserRow({ user, onRoleChange, currentUser }: UserRowProps) {
  const [isUpdating, setIsUpdating] = useState(false);

  const handleRoleChange = async (newRole: 'player' | 'admin') => {
    if (user.id === currentUser.id) {
      alert("You cannot change your own role");
      return;
    }

    setIsUpdating(true);
    try {
      await onRoleChange(user.id, newRole);
    } finally {
      setIsUpdating(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  return (
    <tr>
      <td>{user.id}</td>
      <td>{user.full_name}</td>
      <td>{user.username}</td>
      <td>{user.email}</td>
      <td>
        <select 
          value={user.role} 
          onChange={(e) => handleRoleChange(e.target.value as 'player' | 'admin')}
          disabled={isUpdating || user.id === currentUser.id}
          className="role-select"
        >
          <option value="player">Player</option>
          <option value="admin">Admin</option>
        </select>
      </td>
      <td>
        <span className={`badge ${user.is_active ? 'badge-success' : 'badge-danger'}`}>
          {user.is_active ? 'Active' : 'Inactive'}
        </span>
      </td>
      <td>{user.ranking || '-'}</td>
      <td>{formatDate(user.created_at)}</td>
      
      <style jsx>{`
        .role-select {
          padding: 0.25rem 0.5rem;
          border: 1px solid #ddd;
          border-radius: 4px;
          font-size: 0.875rem;
        }
        
        .role-select:disabled {
          opacity: 0.6;
          cursor: not-allowed;
        }
      `}</style>
    </tr>
  );
}

interface UserStatsProps {
  users: User[];
}

function UserStats({ users }: UserStatsProps) {
  const totalUsers = users.length;
  const adminCount = users.filter(user => user.role === 'admin').length;
  const playerCount = users.filter(user => user.role === 'player').length;
  const activeUsers = users.filter(user => user.is_active).length;

  return (
    <div className="stats-grid">
      <div className="stat-card">
        <h3>Total Users</h3>
        <p className="stat-number">{totalUsers}</p>
      </div>
      <div className="stat-card">
        <h3>Admins</h3>
        <p className="stat-number admin">{adminCount}</p>
      </div>
      <div className="stat-card">
        <h3>Players</h3>
        <p className="stat-number player">{playerCount}</p>
      </div>
      <div className="stat-card">
        <h3>Active Users</h3>
        <p className="stat-number active">{activeUsers}</p>
      </div>
      
      <style jsx>{`
        .stats-grid {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
          gap: 1rem;
          margin-bottom: 2rem;
        }
        
        .stat-card {
          background: white;
          padding: 1.5rem;
          border-radius: 8px;
          box-shadow: 0 2px 4px rgba(0,0,0,0.1);
          text-align: center;
          border-left: 4px solid #667eea;
        }
        
        .stat-card h3 {
          margin: 0 0 0.5rem 0;
          font-size: 0.875rem;
          color: #666;
          text-transform: uppercase;
          letter-spacing: 0.05em;
        }
        
        .stat-number {
          margin: 0;
          font-size: 2rem;
          font-weight: bold;
          color: #333;
        }
        
        .stat-number.admin { color: #dc3545; }
        .stat-number.player { color: #28a745; }
        .stat-number.active { color: #007bff; }
      `}</style>
    </div>
  );
}

export default function AdminPage() {
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const router = useRouter();

  useEffect(() => {
    const initializeAdminPage = async () => {
      try {
        // Check if user is logged in and is admin
        const token = localStorage.getItem('token');
        if (!token) {
          router.push('/auth/login');
          return;
        }

        // Get current user profile
        const userProfile = await authAPI.getProfile();
        if (userProfile.role !== 'admin') {
          router.push('/dashboard');
          return;
        }

        setCurrentUser(userProfile);

        // Fetch all users (admin only endpoint)
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/v1/users`, {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        });

        if (!response.ok) {
          throw new Error('Failed to fetch users');
        }

        const usersData = await response.json();
        setUsers(usersData.data || []);
      } catch (error: any) {
        console.error('Admin page initialization error:', error);
        if (error.message.includes('401') || error.message.includes('unauthorized')) {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
          router.push('/auth/login');
        } else {
          setError('Failed to load admin data');
        }
      } finally {
        setLoading(false);
      }
    };

    initializeAdminPage();
  }, [router]);

  const handleRoleChange = async (userId: number, newRole: 'player' | 'admin') => {
    try {
      const token = localStorage.getItem('token');
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/v1/users/${userId}/role`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ role: newRole }),
      });

      if (!response.ok) {
        throw new Error('Failed to update user role');
      }

      const updatedUser = await response.json();
      
      // Update local state
      setUsers(prevUsers => 
        prevUsers.map(user => 
          user.id === userId 
            ? { ...user, role: newRole }
            : user
        )
      );

      alert(`User role updated to ${newRole} successfully!`);
    } catch (error: any) {
      console.error('Role change error:', error);
      alert(error.message || 'Failed to update user role');
    }
  };

  if (loading) {
    return (
      <Layout title="Admin Panel">
        <div className="container">
          <div style={{ textAlign: 'center', padding: '2rem' }}>
            <h2>Loading admin panel...</h2>
          </div>
        </div>
      </Layout>
    );
  }

  if (error) {
    return (
      <Layout title="Admin Panel">
        <div className="container">
          <div style={{ textAlign: 'center', padding: '2rem', color: '#dc3545' }}>
            <h2>Error: {error}</h2>
            <button 
              className="button button-primary" 
              onClick={() => window.location.reload()}
            >
              Retry
            </button>
          </div>
        </div>
      </Layout>
    );
  }

  if (!currentUser) {
    return null; // Will redirect
  }

  return (
    <Layout title="Admin Panel">
      <div className="container">
        <div className="admin-header">
          <h1>🔧 Admin Panel</h1>
          <p>Welcome, {currentUser.full_name}! Manage users and their roles.</p>
        </div>

        <UserStats users={users} />

        <div className="card">
          <div className="section-header">
            <h2>User Management</h2>
            <div className="header-actions">
              <button 
                className="button button-primary"
                onClick={() => router.push('/admin/tournaments')}
              >
                Manage Tournaments
              </button>
            </div>
          </div>

          <div className="table-container">
            <table className="table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Full Name</th>
                  <th>Username</th>
                  <th>Email</th>
                  <th>Role</th>
                  <th>Status</th>
                  <th>Ranking</th>
                  <th>Created</th>
                </tr>
              </thead>
              <tbody>
                {users.map(user => (
                  <UserRow
                    key={user.id}
                    user={user}
                    onRoleChange={handleRoleChange}
                    currentUser={currentUser}
                  />
                ))}
              </tbody>
            </table>
          </div>

          {users.length === 0 && (
            <div className="empty-state">
              <p>No users found.</p>
            </div>
          )}
        </div>

        <style jsx>{`
          .admin-header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 2rem;
            border-radius: 12px;
            margin-bottom: 2rem;
            text-align: center;
          }
          
          .admin-header h1 {
            margin: 0 0 0.5rem 0;
            font-size: 2.5rem;
          }
          
          .admin-header p {
            margin: 0;
            opacity: 0.9;
            font-size: 1.1rem;
          }
          
          .section-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 1.5rem;
          }
          
          .section-header h2 {
            margin: 0;
            color: #333;
          }
          
          .header-actions {
            display: flex;
            gap: 1rem;
          }
          
          .table-container {
            overflow-x: auto;
          }
          
          .empty-state {
            text-align: center;
            padding: 2rem;
            color: #666;
          }
        `}</style>
      </div>
    </Layout>
  );
}