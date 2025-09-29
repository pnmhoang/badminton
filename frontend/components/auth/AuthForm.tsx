import { ReactNode, FormEvent, ChangeEvent } from 'react';

interface FormFieldProps {
  label: string;
  type?: string;
  value: string;
  onChange: (e: ChangeEvent<HTMLInputElement>) => void;
  required?: boolean;
  disabled?: boolean;
  placeholder?: string;
  minLength?: number;
  error?: string;
}

export const FormField = ({
  label,
  type = 'text',
  value,
  onChange,
  required = false,
  disabled = false,
  placeholder,
  minLength,
  error
}: FormFieldProps) => {
  return (
    <div className="form-group">
      <label>{label}</label>
      <input
        type={type}
        value={value}
        onChange={onChange}
        required={required}
        disabled={disabled}
        placeholder={placeholder}
        minLength={minLength}
      />
      {error && (
        <div style={{ color: '#c33', fontSize: '0.875rem', marginTop: '0.25rem' }}>
          {error}
        </div>
      )}
    </div>
  );
};

interface AuthFormProps {
  title: string;
  onSubmit: (e: FormEvent<HTMLFormElement>) => void;
  loading: boolean;
  error?: string;
  children: ReactNode;
  submitText: string;
  footerText?: string;
  footerLink?: {
    text: string;
    href: string;
  };
}

export const AuthForm = ({
  title,
  onSubmit,
  loading,
  error,
  children,
  submitText,
  footerText,
  footerLink
}: AuthFormProps) => {
  return (
    <div className="container" style={{ maxWidth: '500px', marginTop: '2rem' }}>
      <div className="card">
        <h1 style={{ textAlign: 'center', marginBottom: '2rem' }}>
          🏸 {title}
        </h1>
        
        {error && (
          <div style={{ 
            background: '#fee', 
            color: '#c33', 
            padding: '1rem', 
            borderRadius: '4px',
            marginBottom: '1rem'
          }}>
            {error}
          </div>
        )}

        <form onSubmit={onSubmit}>
          {children}
          
          <button 
            type="submit" 
            className="button" 
            disabled={loading}
            style={{ width: '100%', marginBottom: '1rem' }}
          >
            {loading ? 'Loading...' : submitText}
          </button>
        </form>

        {footerText && footerLink && (
          <div style={{ textAlign: 'center' }}>
            <p>{footerText} <a href={footerLink.href}>{footerLink.text}</a></p>
          </div>
        )}
      </div>
    </div>
  );
};