'use client';

import { useState, ChangeEvent, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import { authAPI } from '../../../utils/api';
import { AuthForm, FormField } from '../../../components/auth/AuthForm';
import type { LoginRequest } from '../../../types';

export default function LoginPage() {
  const [formData, setFormData] = useState<LoginRequest>({
    username: '',
    password: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const router = useRouter();

  const handleInputChange = (field: keyof LoginRequest) => 
    (e: ChangeEvent<HTMLInputElement>) => {
      setFormData(prev => ({
        ...prev,
        [field]: e.target.value
      }));
    };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const response = await authAPI.login(formData);
      
      // Store token and user info
      localStorage.setItem('token', response.data.token);
      localStorage.setItem('user', JSON.stringify(response.data.user));
      
      // Redirect based on role
      if (response.data.user.role === 'admin') {
        router.push('/admin');
      } else {
        router.push('/dashboard');
      }
    } catch (error: any) {
      console.error('Login error:', error);
      setError(error.response?.data?.message || 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthForm
      title="Login"
      onSubmit={handleSubmit}
      loading={loading}
      error={error}
      submitText="Login"
      footerText="Don't have an account?"
      footerLink={{ text: 'Register here', href: '/auth/register' }}
    >
      <FormField
        label="Username or Email"
        type="text"
        value={formData.username}
        onChange={handleInputChange('username')}
        required
        disabled={loading}
      />
      
      <FormField
        label="Password"
        type="password"
        value={formData.password}
        onChange={handleInputChange('password')}
        required
        disabled={loading}
      />
    </AuthForm>
  );
}