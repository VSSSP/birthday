import { useEffect, useState } from 'react';

interface ToastProps {
  message: string;
  type?: 'error' | 'success';
  onClose: () => void;
}

export default function Toast({ message, type = 'error', onClose }: ToastProps) {
  const [visible, setVisible] = useState(true);

  useEffect(() => {
    const timer = setTimeout(() => {
      setVisible(false);
      setTimeout(onClose, 300);
    }, 3000);
    return () => clearTimeout(timer);
  }, [onClose]);

  return (
    <div
      style={{
        position: 'fixed',
        bottom: '2rem',
        left: '50%',
        transform: `translateX(-50%) translateY(${visible ? '0' : '1rem'})`,
        background: type === 'error' ? 'var(--danger)' : 'var(--success)',
        color: 'white',
        padding: '0.75rem 1.5rem',
        borderRadius: 'var(--radius-sm)',
        fontSize: '0.9rem',
        zIndex: 100,
        opacity: visible ? 1 : 0,
        transition: 'opacity 0.3s, transform 0.3s',
        maxWidth: 'calc(100% - 2rem)',
        textAlign: 'center',
      }}
    >
      {message}
    </div>
  );
}
