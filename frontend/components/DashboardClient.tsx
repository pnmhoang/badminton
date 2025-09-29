'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Layout from './Layout';
import { authAPI, tournamentAPI } from '../utils/api';
import type { Tournament, User } from '../types';

// Component for individual tournament card
interface TournamentCardProps {
  tournament: Tournament;
  onRegister: (tournamentId: number) => void;
}

function TournamentCard({ tournament, onRegister }: TournamentCardProps) {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  const getStatusColor = (status: Tournament['status']) => {
    const colors = {
      upcoming: '#28a745',
      ongoing: '#007bff', 
      completed: '#6c757d',
      cancelled: '#dc3545'
    };
    return colors[status] || '#6c757d';
  };

  return (
    <div className="tournament-card">
      <div className="tournament-header">
        <h3>{tournament.name}</h3>
        <span 
          className="tournament-status"
          style={{ backgroundColor: getStatusColor(tournament.status) }}
        >
          {tournament.status.charAt(0).toUpperCase() + tournament.status.slice(1)}
        </span>
      </div>
      
      <div className="tournament-info">
        <p><strong>Type:</strong> {tournament.tournament_type.replace('_', ' ').toUpperCase()}</p>
        <p><strong>Max Participants:</strong> {tournament.max_participants}</p>
        <p><strong>Start Date:</strong> {formatDate(tournament.start_date)}</p>
        <p><strong>Registration Deadline:</strong> {formatDate(tournament.registration_deadline)}</p>
        <p><strong>Entry Fee:</strong> ${tournament.entry_fee}</p>
        <p><strong>Prize Pool:</strong> ${tournament.prize_pool}</p>
      </div>
      
      {tournament.description && (
        <div className="tournament-description">
          <p>{tournament.description}</p>
        </div>
      )}
      
      {tournament.status === 'upcoming' && (
        <button 
          className="button button-primary"
          onClick={() => onRegister(tournament.id)}
          style={{ width: '100%', marginTop: '1rem' }}
        >
          Register
        </button>
      )}
      
      <style jsx>{`
        .tournament-card {
          background: white;
          border: 1px solid #ddd;
          border-radius: 8px;
          padding: 1.5rem;
          margin-bottom: 1rem;
          box-shadow: 0 2px 4px rgba(0,0,0,0.1);
          transition: transform 0.2s;
        }
        
        .tournament-card:hover {
          transform: translateY(-2px);
          box-shadow: 0 4px 8px rgba(0,0,0,0.15);
        }
        
        .tournament-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 1rem;
        }
        
        .tournament-header h3 {
          margin: 0;
          color: #333;
        }
        
        .tournament-status {
          color: white;
          padding: 0.25rem 0.75rem;
          border-radius: 12px;
          font-size: 0.875rem;
          font-weight: bold;
        }
        
        .tournament-info {
          margin-bottom: 1rem;
        }
        
        .tournament-info p {
          margin: 0.5rem 0;
          font-size: 0.875rem;
        }
        
        .tournament-description {
          background: #f8f9fa;
          padding: 1rem;
          border-radius: 4px;
          margin: 1rem 0;
        }
        
        .tournament-description p {
          margin: 0;
          font-style: italic;
          color: #666;
        }
      `}</style>
    </div>
  );
}

// Welcome section component
interface WelcomeSectionProps {
  user: User;
}

function WelcomeSection({ user }: WelcomeSectionProps) {
  return (
    <div className="welcome-section">
      <h1>Welcome back, {user.full_name}! 🏸</h1>
      <div className="user-stats">
        <div className="stat-card">
          <h3>Your Role</h3>
          <p>{user.role.charAt(0).toUpperCase() + user.role.slice(1)}</p>
        </div>
        {user.role === 'player' && user.ranking && (
          <div className="stat-card">
            <h3>Current Ranking</h3>
            <p>#{user.ranking}</p>
          </div>
        )}
        <div className="stat-card">
          <h3>Account Status</h3>
          <p>{user.is_active ? 'Active' : 'Inactive'}</p>
        </div>
      </div>
      
      <style jsx>{`
        .welcome-section {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          color: white;
          padding: 2rem;
          border-radius: 12px;
          margin-bottom: 2rem;
        }
        
        .welcome-section h1 {
          margin: 0 0 1.5rem 0;
          font-size: 2rem;
        }
        
        .user-stats {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
          gap: 1rem;
        }
        
        .stat-card {
          background: rgba(255, 255, 255, 0.1);
          padding: 1rem;
          border-radius: 8px;
          text-align: center;
        }
        
        .stat-card h3 {
          margin: 0 0 0.5rem 0;
          font-size: 0.875rem;
          opacity: 0.8;
        }
        
        .stat-card p {
          margin: 0;
          font-size: 1.25rem;
          font-weight: bold;
        }
      `}</style>
    </div>
  );
}

export default function DashboardPage() {
  const [user, setUser] = useState<User | null>(null);
  const [tournaments, setTournaments] = useState<Tournament[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const router = useRouter();

  useEffect(() => {
    const initializeDashboard = async () => {
      try {
        // Check if user is logged in
        const token = localStorage.getItem('token');
        if (!token) {
          router.push('/auth/login');
          return;
        }

        // Fetch user profile
        const userProfile = await authAPI.getProfile();
        setUser(userProfile);

        // Fetch tournaments
        const tournamentsResponse = await tournamentAPI.getTournaments();
        setTournaments(tournamentsResponse.data || []);
      } catch (error: any) {
        console.error('Dashboard initialization error:', error);
        if (error.response?.status === 401) {
          // Token expired, redirect to login
          localStorage.removeItem('token');
          localStorage.removeItem('user');
          router.push('/auth/login');
        } else {
          setError('Failed to load dashboard data');
        }
      } finally {
        setLoading(false);
      }
    };

    initializeDashboard();
  }, [router]);

  const handleTournamentRegistration = async (tournamentId: number) => {
    try {
      await tournamentAPI.registerForTournament(tournamentId);
      alert('Successfully registered for tournament!');
      
      // Refresh tournaments to update UI
      const tournamentsResponse = await tournamentAPI.getTournaments();
      setTournaments(tournamentsResponse.data || []);
    } catch (error: any) {
      console.error('Registration error:', error);
      alert(error.response?.data?.message || 'Failed to register for tournament');
    }
  };

  if (loading) {
    return (
      <Layout title="Dashboard">
        <div className="container">
          <div style={{ textAlign: 'center', padding: '2rem' }}>
            <h2>Loading your dashboard...</h2>
          </div>
        </div>
      </Layout>
    );
  }

  if (error) {
    return (
      <Layout title="Dashboard">
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

  if (!user) {
    return null; // Will redirect to login
  }

  return (
    <Layout title="Dashboard">
      <div className="container">
        <WelcomeSection user={user} />
        
        <section className="tournaments-section">
          <div className="section-header">
            <h2>Available Tournaments</h2>
            {user.role === 'admin' && (
              <button 
                className="button button-primary"
                onClick={() => router.push('/admin/tournaments/new')}
              >
                Create Tournament
              </button>
            )}
          </div>
          
          {tournaments.length === 0 ? (
            <div className="empty-state">
              <p>No tournaments available at the moment.</p>
              {user.role === 'admin' && (
                <p>Create the first tournament to get started!</p>
              )}
            </div>
          ) : (
            <div className="tournaments-grid">
              {tournaments.map(tournament => (
                <TournamentCard
                  key={tournament.id}
                  tournament={tournament}
                  onRegister={handleTournamentRegistration}
                />
              ))}
            </div>
          )}
        </section>
        
        <style jsx>{`
          .tournaments-section {
            margin-top: 2rem;
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
          
          .empty-state {
            text-align: center;
            padding: 3rem 1rem;
            background: #f8f9fa;
            border-radius: 8px;
            color: #6c757d;
          }
          
          .tournaments-grid {
            display: grid;
            gap: 1rem;
          }
          
          @media (min-width: 768px) {
            .tournaments-grid {
              grid-template-columns: repeat(2, 1fr);
            }
          }
          
          @media (min-width: 1200px) {
            .tournaments-grid {
              grid-template-columns: repeat(3, 1fr);
            }
          }
        `}</style>
      </div>
    </Layout>
  );
}