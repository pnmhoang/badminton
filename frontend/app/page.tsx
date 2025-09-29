'use client';

import Layout from '../components/Layout';

export default function HomePage() {
  return (
    <Layout title="Badminton Tournament">
      <div className="container">
        <div className="hero-section">
          <h1>🏸 Welcome to Badminton Tournament Management</h1>
          <p>
            Organize and participate in badminton tournaments with ease. 
            Join tournaments, create teams, and track your progress!
          </p>
          <div className="hero-buttons">
            <a href="/auth/register" className="button button-primary">
              Get Started
            </a>
            <a href="/auth/login" className="button button-secondary">
              Sign In
            </a>
          </div>
        </div>

        <div className="features-section">
          <h2>Features</h2>
          <div className="features-grid">
            <div className="feature-card">
              <h3>🎯 Tournament Management</h3>
              <p>Create and manage singles, doubles, and mixed doubles tournaments</p>
            </div>
            <div className="feature-card">
              <h3>👥 Team Formation</h3>
              <p>Form teams for doubles tournaments and manage team registrations</p>
            </div>
            <div className="feature-card">
              <h3>🏆 Live Scoring</h3>
              <p>Track matches in real-time with comprehensive scoring system</p>
            </div>
            <div className="feature-card">
              <h3>📊 Player Rankings</h3>
              <p>Track player performance and maintain ranking systems</p>
            </div>
          </div>
        </div>
      </div>
    </Layout>
  );
}